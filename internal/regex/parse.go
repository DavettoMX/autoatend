package regex

import (
	"fmt"
	"unicode"
)

// MaxDepth is the deepest nesting of parenthesized groups the parser accepts.
// Each level costs a few Go stack frames (expresion → alternancia →
// concatenacion → factor → primario), and Thompson's construction recurses
// over the same tree, so the limit keeps both bounded for any input. The ERs
// the generator emits nest one level; the TE's stress case has 50 levels.
const MaxDepth = 200

// metachars are the characters that must be escaped with '\' to be literal:
// the operators of the grammar, the escape itself and '-' (the range
// operator). It is the escape set of lexer.GenerateER.
const metachars = `\|*+?()[]-`

// Parse parses src with the TE §6.2.3 grammar:
//
//	expresion     → alternancia
//	alternancia   → concatenacion ( '|' concatenacion )*
//	concatenacion → factor factor*
//	factor        → primario cuantificador?
//	cuantificador → '*' | '+' | '?'
//	primario      → '(' expresion ')' | clase_char | caracter
//	clase_char    → '[' rango+ ']'
//	rango         → caracter '-' caracter | caracter
//	caracter      → LETRA | DIGITO | '\' metacaracter
//
// LETRA is any Unicode letter (the lexer's alphabet is a-z and áéíóúüñ) and
// DIGITO a decimal digit. The grammar has no empty alternative: "a|", "()"
// and "" are errors. Closures do not stack: "a**" is an error. Errors are
// *ParseError values carrying the rune position.
func Parse(src string) (Node, error) {
	p := &parser{src: []rune(src)}
	if len(p.src) == 0 {
		return nil, &ParseError{0, "empty expression"}
	}
	n, err := p.alternation(0)
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.src) { // alternation only stops early at an unmatched ')'
		return nil, &ParseError{p.pos, "unbalanced ')': no matching '('"}
	}
	return n, nil
}

type parser struct {
	src []rune
	pos int
}

func (p *parser) eof() bool { return p.pos >= len(p.src) }

func (p *parser) peek() rune { return p.src[p.pos] }

func (p *parser) errf(pos int, format string, args ...any) error {
	return &ParseError{pos, fmt.Sprintf(format, args...)}
}

// alternation → concatenation ( '|' concatenation )*. depth is the number of
// enclosing groups. It stops at EOF or at a ')' for the caller to handle.
func (p *parser) alternation(depth int) (Node, error) {
	var alts []Node
	for {
		c, err := p.concatenation(depth)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, p.errf(p.pos, "empty alternative: expected a character, group or class")
		}
		alts = append(alts, c)
		if p.eof() || p.peek() != '|' {
			break
		}
		p.pos++
	}
	if len(alts) == 1 {
		return alts[0], nil
	}
	return &Alternation{alts}, nil
}

// concatenation → factor factor*. It returns nil, nil when no factor starts
// at the current position (end of input, '|' or ')').
func (p *parser) concatenation(depth int) (Node, error) {
	var parts []Node
	for !p.eof() && p.peek() != '|' && p.peek() != ')' {
		f, err := p.factor(depth)
		if err != nil {
			return nil, err
		}
		parts = append(parts, f)
	}
	switch len(parts) {
	case 0:
		return nil, nil
	case 1:
		return parts[0], nil
	}
	return &Concatenation{parts}, nil
}

// factor → primary quantifier?
func (p *parser) factor(depth int) (Node, error) {
	n, err := p.primary(depth)
	if err != nil {
		return nil, err
	}
	if p.eof() {
		return n, nil
	}
	var op Op
	switch p.peek() {
	case '*':
		op = OpStar
	case '+':
		op = OpPlus
	case '?':
		op = OpOptional
	default:
		return n, nil
	}
	p.pos++
	if !p.eof() && isQuantifier(p.peek()) {
		return nil, p.errf(p.pos, "dangling quantifier %q: it follows another quantifier", p.peek())
	}
	return &Closure{op, n}, nil
}

func isQuantifier(r rune) bool { return r == '*' || r == '+' || r == '?' }

// primary → '(' expression ')' | class | character
func (p *parser) primary(depth int) (Node, error) {
	r := p.peek()
	switch {
	case r == '(':
		return p.group(depth)
	case r == '[':
		return p.class()
	case isQuantifier(r):
		return nil, p.errf(p.pos, "dangling quantifier %q: nothing to repeat", r)
	}
	c, err := p.character("")
	if err != nil {
		return nil, err
	}
	return &Char{c}, nil
}

func (p *parser) group(depth int) (Node, error) {
	open := p.pos
	if depth+1 > MaxDepth {
		return nil, p.errf(open, "nesting too deep: more than %d nested groups", MaxDepth)
	}
	p.pos++
	n, err := p.alternation(depth + 1)
	if err != nil {
		return nil, err
	}
	if p.eof() {
		return nil, p.errf(open, "unbalanced '(': missing ')' before the end of the expression")
	}
	p.pos++ // the ')' that stopped the alternation
	return &Group{n}, nil
}

// class → '[' range+ ']'. The position of every class error is the '[' except
// for a bad range, which points at the offending range or '-'.
func (p *parser) class() (Node, error) {
	open := p.pos
	p.pos++
	var ranges []Range
	for {
		if p.eof() {
			return nil, p.errf(open, "unterminated class: missing ']'")
		}
		if p.peek() == ']' {
			p.pos++
			break
		}
		start := p.pos
		lo, err := p.character("class")
		if err != nil {
			return nil, err
		}
		hi := lo
		if !p.eof() && p.peek() == '-' {
			dash := p.pos
			p.pos++
			if p.eof() {
				return nil, p.errf(open, "unterminated class: missing ']'")
			}
			if p.peek() == ']' {
				return nil, p.errf(dash, "bad range: '-' without an upper bound (escape it as \\- to match a hyphen)")
			}
			if hi, err = p.character("class"); err != nil {
				return nil, err
			}
			if hi < lo {
				return nil, p.errf(start, "bad range %c-%c: lower bound is greater than upper bound", lo, hi)
			}
		}
		ranges = append(ranges, Range{lo, hi})
	}
	if len(ranges) == 0 {
		return nil, p.errf(open, "empty class: '[]' needs at least one range")
	}
	return &Class{ranges}, nil
}

// character → LETTER | DIGIT | '\' metacharacter. where is "class" inside a
// class, where an unescaped '-' is a bad range instead of an unexpected character.
func (p *parser) character(where string) (rune, error) {
	r := p.peek()
	switch {
	case r == '\\':
		if p.pos+1 >= len(p.src) {
			return 0, p.errf(p.pos, "trailing backslash: '\\' must be followed by a metacharacter")
		}
		m := p.src[p.pos+1]
		if !isMeta(m) {
			return 0, p.errf(p.pos, "invalid escape '\\%c': only metacharacters can be escaped (%s)", m, metachars)
		}
		p.pos += 2
		return m, nil
	case unicode.IsLetter(r) || unicode.IsDigit(r):
		p.pos++
		return r, nil
	case r == '-' && where == "class":
		return 0, p.errf(p.pos, "bad range: '-' without a lower bound (escape it as \\- to match a hyphen)")
	}
	return 0, p.errf(p.pos, "unexpected character %q: escape it with '\\' to match it literally", r)
}

func isMeta(r rune) bool {
	for _, m := range metachars {
		if r == m {
			return true
		}
	}
	return false
}
