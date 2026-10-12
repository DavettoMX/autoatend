package regex

import (
	"errors"
	"strings"
	"testing"
)

// mustParse parses src and fails the test on error.
func mustParse(t *testing.T, src string) Node {
	t.Helper()
	n, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse(%q): %v", src, err)
	}
	return n
}

// CPU01 (TE §6.2.3): the nominal expression the generator emits parses into
// the expected typed AST.
func TestCPU01_ParserNominal(t *testing.T) {
	n := mustParse(t, "factur(a|as|ación|ado|ando)")
	want := "Cat(f,a,c,t,u,r,Group(Alt(a,Cat(a,s),Cat(a,c,i,ó,n),Cat(a,d,o),Cat(a,n,d,o))))"
	if got := n.String(); got != want {
		t.Fatalf("AST = %s\nwant  %s", got, want)
	}

	// The node kinds are the typed ones of the TE, not just a rendering.
	cat, ok := n.(*Concatenation)
	if !ok || len(cat.Parts) != 7 {
		t.Fatalf("root = %T, want *Concatenation with 7 parts", n)
	}
	if c, ok := cat.Parts[0].(*Char); !ok || c.R != 'f' {
		t.Errorf("first part = %#v, want Char f", cat.Parts[0])
	}
	grp, ok := cat.Parts[6].(*Group)
	if !ok {
		t.Fatalf("last part = %T, want *Group", cat.Parts[6])
	}
	alt, ok := grp.Expr.(*Alternation)
	if !ok || len(alt.Alts) != 5 {
		t.Fatalf("group content = %T, want *Alternation with 5 alternatives", grp.Expr)
	}
}

func TestCPU01_ParserSingleChar(t *testing.T) {
	// A one-element alternation or concatenation collapses to its element.
	if c, ok := mustParse(t, "a").(*Char); !ok || c.R != 'a' {
		t.Errorf("Parse(a) is not Char a")
	}
	if _, ok := mustParse(t, "(a)").(*Group); !ok {
		t.Errorf("Parse((a)) is not a Group")
	}
}

// Closure binds tighter than concatenation, which binds tighter than alternation.
func TestCPU03_ParserPrecedence(t *testing.T) {
	cases := []struct{ src, want string }{
		{"ab*", "Cat(a,Star(b))"},
		{"a|bc", "Alt(a,Cat(b,c))"},
		{"ab|cd", "Alt(Cat(a,b),Cat(c,d))"},
		{"(ab)*", "Star(Group(Cat(a,b)))"},
		{"ab+|c?", "Alt(Cat(a,Plus(b)),Opt(c))"},
		{"a|b|c", "Alt(a,b,c)"},
		{"x(a|b)+y", "Cat(x,Plus(Group(Alt(a,b))),y)"},
		{"[a-c]*d", "Cat(Star(Class[a-c]),d)"},
		{"[a-cxz]", "Class[a-c,x,z]"},
		{"[ñü]", "Class[ñ,ü]"},
		{"7a", "Cat(7,a)"},
		{"((a))", "Group(Group(a))"},
	}
	for _, c := range cases {
		if got := mustParse(t, c.src).String(); got != c.want {
			t.Errorf("Parse(%q) = %s, want %s", c.src, got, c.want)
		}
	}
}

func TestCPU01_ParserEscapes(t *testing.T) {
	// Every metacharacter of the generator's escape set is a literal when escaped.
	for _, m := range `|*+?()[]\-` {
		src := `a\` + string(m) + `b`
		n := mustParse(t, src)
		cat, ok := n.(*Concatenation)
		if !ok || len(cat.Parts) != 3 {
			t.Fatalf("Parse(%q) = %s, want a 3-part concatenation", src, n)
		}
		if c, ok := cat.Parts[1].(*Char); !ok || c.R != m {
			t.Errorf("Parse(%q): middle = %#v, want Char %q", src, cat.Parts[1], m)
		}
	}
	// Escaped characters also work inside a class.
	if got := mustParse(t, `[\-a\]]`).String(); got != `Class[-,a,]]` {
		t.Errorf("class with escapes = %s", got)
	}
}

// CPU02 (TE §6.2.3): invalid expressions give a descriptive error with the position.
func TestCPU02_ParserErrors(t *testing.T) {
	cases := []struct {
		src  string
		pos  int
		want string // substring of the message
	}{
		{"factur(a|as", 6, "unbalanced '('"},
		{"a)", 1, "unbalanced ')'"},
		{"((a)", 0, "unbalanced '('"},
		{"", 0, "empty"},
		{"a|", 2, "empty alternative"},
		{"|a", 0, "empty alternative"},
		{"a||b", 2, "empty alternative"},
		{"()", 1, "empty alternative"},
		{"(|a)", 1, "empty alternative"},
		{"*a", 0, "dangling quantifier"},
		{"+", 0, "dangling quantifier"},
		{"a**", 2, "dangling quantifier"},
		{"a|*", 2, "dangling quantifier"},
		{"(?a)", 1, "dangling quantifier"},
		{"[z-a]", 1, "bad range"},
		{"[a-]", 2, "bad range"},
		{"[-a]", 1, "bad range"},
		{"[a-c-e]", 4, "bad range"},
		{"[]", 0, "empty class"},
		{"[a", 0, "unterminated class"},
		{"x[a-c", 1, "unterminated class"},
		{"a\\", 1, "trailing backslash"},
		{`\a`, 0, "invalid escape"},
		{"a b", 1, "unexpected character"},
		{"a-b", 1, "unexpected character"},
		{"]", 0, "unexpected character"},
		{"a.b", 1, "unexpected character"},
	}
	for _, c := range cases {
		_, err := Parse(c.src)
		if err == nil {
			t.Errorf("Parse(%q): expected an error", c.src)
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error %T is not *ParseError", c.src, err)
			continue
		}
		if pe.Pos != c.pos {
			t.Errorf("Parse(%q): position %d, want %d (%v)", c.src, pe.Pos, c.pos, err)
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("Parse(%q): message %q does not contain %q", c.src, err, c.want)
		}
		if !strings.Contains(err.Error(), "position") {
			t.Errorf("Parse(%q): message %q does not mention the position", c.src, err)
		}
	}
}

// nested builds an expression of `levels` nested alternatives,
// (abcdef|(abcdef|(...x))), padding the innermost literal to total runes.
func nested(levels, total int) string {
	var b strings.Builder
	for i := 0; i < levels; i++ {
		b.WriteString("(abcdef|")
	}
	b.WriteString(strings.Repeat("x", total-b.Len()-levels))
	b.WriteString(strings.Repeat(")", levels))
	return b.String()
}

// depth is the deepest nesting of Group nodes.
func depth(n Node) int {
	switch v := n.(type) {
	case *Group:
		return 1 + depth(v.Expr)
	case *Alternation:
		d := 0
		for _, a := range v.Alts {
			d = max(d, depth(a))
		}
		return d
	case *Concatenation:
		d := 0
		for _, a := range v.Parts {
			d = max(d, depth(a))
		}
		return d
	case *Closure:
		return depth(v.Operand)
	}
	return 0
}

// CPU03 (TE §6.2.3): a 500-character expression with 50 nested alternatives
// parses into a valid AST without exhausting the stack.
func TestCPU03_ParserDeepButValid(t *testing.T) {
	src := nested(50, 500)
	if len([]rune(src)) != 500 {
		t.Fatalf("test input has %d runes, want 500", len([]rune(src)))
	}
	n := mustParse(t, src)
	if d := depth(n); d != 50 {
		t.Errorf("AST nesting = %d, want 50", d)
	}
	// 50 flat alternatives are valid too.
	flat := strings.Repeat("ab|", 49) + "ab"
	if alt, ok := mustParse(t, flat).(*Alternation); !ok || len(alt.Alts) != 50 {
		t.Errorf("flat alternation did not yield 50 alternatives")
	}
}

func TestCPU03_ParserDepthLimit(t *testing.T) {
	ok := strings.Repeat("(", MaxDepth) + "a" + strings.Repeat(")", MaxDepth)
	if d := depth(mustParse(t, ok)); d != MaxDepth {
		t.Errorf("nesting at the limit = %d, want %d", d, MaxDepth)
	}
	over := strings.Repeat("(", MaxDepth+1) + "a" + strings.Repeat(")", MaxDepth+1)
	_, err := Parse(over)
	var pe *ParseError
	if !errors.As(err, &pe) || !strings.Contains(err.Error(), "nesting too deep") {
		t.Fatalf("beyond the limit: err = %v, want a nesting error", err)
	}
	if pe.Pos != MaxDepth {
		t.Errorf("limit error at position %d, want %d (the offending '(')", pe.Pos, MaxDepth)
	}
	// A very deep input must be rejected without running out of stack.
	if _, err := Parse(strings.Repeat("(", 2_000_000)); err == nil {
		t.Error("2M nested parentheses were accepted")
	}
}
