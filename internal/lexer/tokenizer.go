package lexer

// Kind is the lexical category a token is emitted with.
type Kind int

const (
	// Palabra is a maximal run of letters.
	Palabra Kind = iota
	// Numero is a maximal run of digits. The DFA recognizes numbers (TE §6.2.2)
	// but Analyze does not catalogue them.
	Numero
)

func (k Kind) String() string {
	if k == Numero {
		return "numero"
	}
	return "palabra"
}

// Token is one lexeme occurrence: its surface text and its kind.
type Token struct {
	Text string
	Kind Kind
}

// state is a state of the tokenizer DFA.
type state uint8

const (
	stInicio state = iota
	stPalabra
	stNumero
	stSeparador
	numStates
)

// class is a class of input characters.
type class uint8

const (
	clLetra class = iota
	clDigito
	clSep
	numClasses
)

// transition is one cell of the DFA table: the next state, and whether the
// token being built (the state we leave) must be emitted before moving on.
type transition struct {
	next state
	emit bool
}

// delta is the explicit transition table [state][class]. The DFA reads the
// input once, left to right, with no backtracking:
//
//	INICIO     --LETRA--> PALABRA     --DIGITO--> NUMERO     --SEP--> SEPARADOR
//	PALABRA    --LETRA--> PALABRA     --DIGITO--> NUMERO (emits the word)  --SEP--> SEPARADOR (emits the word)
//	NUMERO     --DIGITO-> NUMERO      --LETRA---> PALABRA (emits the number) --SEP--> SEPARADOR (emits the number)
//	SEPARADOR  --LETRA--> PALABRA     --DIGITO--> NUMERO     --SEP--> SEPARADOR
//
// A change between a letter and a digit closes the token in progress and opens
// another one, so "abc123" yields the word "abc" and the number "123". The
// end of input emits the token in progress, like a SEP.
var delta = [numStates][numClasses]transition{
	stInicio:    {clLetra: {stPalabra, false}, clDigito: {stNumero, false}, clSep: {stSeparador, false}},
	stPalabra:   {clLetra: {stPalabra, false}, clDigito: {stNumero, true}, clSep: {stSeparador, true}},
	stNumero:    {clLetra: {stPalabra, true}, clDigito: {stNumero, false}, clSep: {stSeparador, true}},
	stSeparador: {clLetra: {stPalabra, false}, clDigito: {stNumero, false}, clSep: {stSeparador, false}},
}

// classify maps a character to its input class. Anything outside the
// lowercase alphabet and the ASCII digits is a separator, so the tokenizer is
// total; the text is expected to come out of Normalize.
func classify(r rune) class {
	switch {
	case isLetter(r):
		return clLetra
	case isDigit(r):
		return clDigito
	}
	return clSep
}

// kindOf is the kind of the token the DFA is building while in state s.
func kindOf(s state) Kind {
	if s == stNumero {
		return Numero
	}
	return Palabra
}

// Tokenize runs the DFA over s and returns its tokens in order of appearance.
func Tokenize(s string) []Token {
	var toks []Token
	st := stInicio
	start := 0
	for i, r := range s {
		tr := delta[st][classify(r)]
		if tr.emit {
			toks = append(toks, Token{Text: s[start:i], Kind: kindOf(st)})
		}
		if tr.next != st && (tr.next == stPalabra || tr.next == stNumero) {
			start = i
		}
		st = tr.next
	}
	if st == stPalabra || st == stNumero {
		toks = append(toks, Token{Text: s[start:], Kind: kindOf(st)})
	}
	return toks
}
