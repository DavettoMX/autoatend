// Package regex is the hand-written regular-expression engine of the lexical
// analyzer (TE §6.2.3 and §6.2.4): a recursive-descent parser over runes that
// builds a typed AST, Thompson's construction of an NFA with ε-transitions
// (one start and one accept state per fragment), and a direct simulation of
// that NFA by state sets. There is no backtracking and no use of the standard
// library's regexp package: matching a string of n runes against an NFA of m
// states takes O(n·m) time and O(m) space.
package regex

import (
	"strconv"
	"strings"
)

// Node is a node of the AST. The concrete kinds map onto the TE names:
//
//	Alternation   NodoAlternancia
//	Concatenation NodoConcatenacion
//	Closure       NodoCerradura
//	Group, Class  NodoPrimario
//	Char          NodoCaracter
type Node interface {
	// String renders the tree as a prefix expression (Alt(...), Cat(...),
	// Star(...), Group(...), Class[...], or the bare character). It is meant
	// for tests and diagnostics, not as a serialization.
	String() string
	node()
}

// Alternation is NodoAlternancia: r1|r2|…|rk with k >= 2. A single
// alternative is never wrapped.
type Alternation struct{ Alts []Node }

// Concatenation is NodoConcatenacion: r1 r2 … rk with k >= 2. A single factor
// is never wrapped.
type Concatenation struct{ Parts []Node }

// Op is the operator of a Closure.
type Op int

// Closure operators (TE cuantificador: '*', '+', '?').
const (
	OpStar     Op = iota // zero or more
	OpPlus               // one or more
	OpOptional           // zero or one
)

// Closure is NodoCerradura: Operand followed by '*', '+' or '?'.
type Closure struct {
	Op      Op
	Operand Node
}

// Group is NodoPrimario in its parenthesized form: ( expresion ).
type Group struct{ Expr Node }

// Range is an inclusive range of runes, Lo <= Hi. A single character c is the
// range c-c.
type Range struct{ Lo, Hi rune }

// Class is NodoPrimario in its character-class form: [ rango+ ]. It matches
// one rune that falls in any of its ranges.
type Class struct{ Ranges []Range }

// Char is NodoCaracter: one literal rune.
type Char struct{ R rune }

func (*Alternation) node()   {}
func (*Concatenation) node() {}
func (*Closure) node()       {}
func (*Group) node()         {}
func (*Class) node()         {}
func (*Char) node()          {}

func (n *Alternation) String() string   { return "Alt(" + joinNodes(n.Alts) + ")" }
func (n *Concatenation) String() string { return "Cat(" + joinNodes(n.Parts) + ")" }
func (n *Group) String() string         { return "Group(" + n.Expr.String() + ")" }
func (n *Char) String() string          { return string(n.R) }

func (n *Closure) String() string {
	name := [...]string{OpStar: "Star", OpPlus: "Plus", OpOptional: "Opt"}[n.Op]
	return name + "(" + n.Operand.String() + ")"
}

func (n *Class) String() string {
	parts := make([]string, len(n.Ranges))
	for i, r := range n.Ranges {
		if r.Lo == r.Hi {
			parts[i] = string(r.Lo)
		} else {
			parts[i] = string(r.Lo) + "-" + string(r.Hi)
		}
	}
	return "Class[" + strings.Join(parts, ",") + "]"
}

func joinNodes(ns []Node) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = n.String()
	}
	return strings.Join(parts, ",")
}

// ParseError is a syntax error. Pos is the 0-based offset, in runes, of the
// offending character in the source expression.
type ParseError struct {
	Pos int
	Msg string
}

func (e *ParseError) Error() string {
	return "regex: " + e.Msg + " (at position " + strconv.Itoa(e.Pos) + ")"
}
