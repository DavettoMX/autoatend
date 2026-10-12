package regex

import "sort"

// edge is a labelled transition: it consumes one rune in [lo, hi].
type edge struct {
	lo, hi rune
	to     int
}

// state is one NFA state: its ε-transitions and its labelled transitions.
// Thompson's construction gives every state at most two outgoing transitions,
// except the start of a character class, which has one per range.
type state struct {
	eps   []int
	edges []edge
}

// NFA is a Thompson automaton with ε-transitions (TE §6.2.4). The whole
// automaton has one start state and one accept state. It is immutable once
// built, so a single NFA can be matched from several goroutines.
type NFA struct {
	states        []state
	start, accept int
}

// frag is a fragment under construction: a sub-automaton with one start and
// one accept state.
type frag struct{ start, accept int }

type builder struct{ states []state }

func (b *builder) newState() int {
	b.states = append(b.states, state{})
	return len(b.states) - 1
}

func (b *builder) eps(from, to int) {
	b.states[from].eps = append(b.states[from].eps, to)
}

func (b *builder) edge(from int, lo, hi rune, to int) {
	b.states[from].edges = append(b.states[from].edges, edge{lo, hi, to})
}

// Compile parses src (TE §6.2.3) and builds its NFA (TE §6.2.4).
func Compile(src string) (*NFA, error) {
	n, err := Parse(src)
	if err != nil {
		return nil, err
	}
	return Thompson(n), nil
}

// Thompson builds the NFA of an AST with Thompson's construction. Each
// fragment has one start and one accept state, and the rules compose them:
//
//   - character: two states and one labelled transition;
//   - class: two states and one labelled transition per range;
//   - concatenation N1 N2: an ε from accept(N1) to start(N2);
//   - alternation N1|N2: a new start with ε to both starts and a new accept
//     with ε from both accepts (n alternatives fold pairwise, left to right);
//   - N*: a new start and accept; start ε→ N.start and ε→ accept, and
//     N.accept ε→ N.start and ε→ accept;
//   - N+: like N* without the ε from the new start to the new accept;
//   - N?: like N* without the ε from N.accept back to N.start.
//
// Every operator of the expression adds at most two states, so the NFA has no
// more than 2·|r| states.
func Thompson(n Node) *NFA {
	b := &builder{}
	f := b.build(n)
	return &NFA{states: b.states, start: f.start, accept: f.accept}
}

func (b *builder) build(n Node) frag {
	switch v := n.(type) {
	case *Char:
		s, a := b.newState(), b.newState()
		b.edge(s, v.R, v.R, a)
		return frag{s, a}
	case *Class:
		s, a := b.newState(), b.newState()
		for _, r := range v.Ranges {
			b.edge(s, r.Lo, r.Hi, a)
		}
		return frag{s, a}
	case *Group:
		return b.build(v.Expr)
	case *Concatenation:
		f := b.build(v.Parts[0])
		for _, p := range v.Parts[1:] {
			g := b.build(p)
			b.eps(f.accept, g.start)
			f.accept = g.accept
		}
		return f
	case *Alternation:
		f := b.build(v.Alts[0])
		for _, alt := range v.Alts[1:] {
			g := b.build(alt)
			s, a := b.newState(), b.newState()
			b.eps(s, f.start)
			b.eps(s, g.start)
			b.eps(f.accept, a)
			b.eps(g.accept, a)
			f = frag{s, a}
		}
		return f
	case *Closure:
		inner := b.build(v.Operand)
		s, a := b.newState(), b.newState()
		b.eps(s, inner.start)
		if v.Op != OpPlus {
			b.eps(s, a)
		}
		if v.Op != OpOptional {
			b.eps(inner.accept, inner.start)
		}
		b.eps(inner.accept, a)
		return frag{s, a}
	}
	panic("regex: unknown AST node")
}

// NumStates is the number of states of the automaton.
func (m *NFA) NumStates() int { return len(m.states) }

// NumEpsilon is the number of ε-transitions.
func (m *NFA) NumEpsilon() int {
	n := 0
	for i := range m.states {
		n += len(m.states[i].eps)
	}
	return n
}

// NumLabelled is the number of labelled (non-ε) transitions.
func (m *NFA) NumLabelled() int {
	n := 0
	for i := range m.states {
		n += len(m.states[i].edges)
	}
	return n
}

// Start is the start state.
func (m *NFA) Start() int { return m.start }

// Accept is the single accept state.
func (m *NFA) Accept() int { return m.accept }

// Match reports whether the whole of s belongs to the language of the
// expression, by simulating the NFA with sets of states (no backtracking):
// S' = ε-closure(move(S, c)) for every rune c of s, and s is accepted iff the
// final set contains the accept state. For n runes and m states it takes
// O(n·m) time and O(m) space.
func (m *NFA) Match(s string) bool {
	ok, _ := m.simulate(s)
	return ok
}

// EpsilonClosure returns the states reachable from the given ones through
// ε-transitions alone (themselves included), in ascending order.
func (m *NFA) EpsilonClosure(states ...int) []int {
	sim := m.newSim()
	var set []int
	for _, s := range states {
		set = sim.add(set, s)
	}
	sort.Ints(set)
	return set
}

// step is ε-closure(move(set, r)) for a set that is already closed.
func (m *NFA) step(set []int, r rune) []int {
	out := m.newSim().move(nil, set, r)
	sort.Ints(out)
	return out
}

// simulate is Match plus a deterministic cost: the number of states entered
// plus the number of transitions examined, which the O(n·m) tests bound.
func (m *NFA) simulate(s string) (bool, int) {
	sim := m.newSim()
	cur := sim.add(nil, m.start)
	var next []int
	for _, r := range s {
		next = sim.move(next[:0], cur, r)
		cur, next = next, cur
		if len(cur) == 0 {
			return false, sim.work
		}
	}
	for _, st := range cur {
		if st == m.accept {
			return true, sim.work
		}
	}
	return false, sim.work
}

// sim holds the scratch memory of one simulation: O(m) in total.
type sim struct {
	m     *NFA
	mark  []int // mark[s] == gen when s is already in the set being built
	gen   int
	stack []int
	work  int
}

func (m *NFA) newSim() *sim {
	return &sim{m: m, mark: make([]int, len(m.states))}
}

// add puts state s and everything reachable from it by ε-transitions into
// set, skipping the states the current generation already holds. The walk is
// iterative, so an ε-cycle (as in (a*)*) cannot loop or overflow the stack.
func (x *sim) add(set []int, s int) []int {
	if x.gen == 0 {
		x.gen = 1
	}
	x.stack = append(x.stack[:0], s)
	for len(x.stack) > 0 {
		top := x.stack[len(x.stack)-1]
		x.stack = x.stack[:len(x.stack)-1]
		if x.mark[top] == x.gen {
			continue
		}
		x.mark[top] = x.gen
		x.work++
		set = append(set, top)
		for _, t := range x.m.states[top].eps {
			x.work++
			if x.mark[t] != x.gen {
				x.stack = append(x.stack, t)
			}
		}
	}
	return set
}

// move appends to dst the ε-closure of the states reachable from set by a
// transition labelled with r.
func (x *sim) move(dst, set []int, r rune) []int {
	x.gen++
	for _, s := range set {
		for _, e := range x.m.states[s].edges {
			x.work++
			if e.lo <= r && r <= e.hi {
				dst = x.add(dst, e.to)
			}
		}
	}
	return dst
}
