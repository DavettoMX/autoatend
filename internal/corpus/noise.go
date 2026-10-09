package corpus

import (
	"math/rand"
	"strings"
	"unicode"
)

// NoiseConfig calibrates the noise injector. Rate is the probability that a
// comment is selected; once selected, every other field is the probability of
// its own operation (per accented letter, per word, or per comment).
type NoiseConfig struct {
	Rate       float64 // comment selected for noise
	TildeDrop  float64 // per accented letter
	Abbrev     float64 // per abbreviable word
	Typo       float64 // per word of 4+ letters
	CaseJitter float64 // per comment
	Punct      float64 // per comment
}

// DefaultNoise is the calibrated configuration used for the case corpus.
var DefaultNoise = NoiseConfig{
	Rate:       0.30,
	TildeDrop:  0.80,
	Abbrev:     0.50,
	Typo:       0.04,
	CaseJitter: 0.25,
	Punct:      0.30,
}

var tildeFold = map[rune]rune{
	'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u', 'ü': 'u',
	'Á': 'A', 'É': 'E', 'Í': 'I', 'Ó': 'O', 'Ú': 'U', 'Ü': 'U',
}

// abbreviations are the colloquial short forms of the case (q, x, tel, fact).
var abbreviations = map[string]string{
	"que": "q", "qué": "q",
	"por":      "x",
	"teléfono": "tel", "telefono": "tel",
	"factura": "fact", "facturación": "fact", "facturacion": "fact",
}

// qwertyNeighbors lists adjacent keys for the typo injector.
var qwertyNeighbors = map[rune]string{
	'a': "qwsz", 'b': "vghn", 'c': "xdfv", 'd': "serfcx", 'e': "wsdr", 'f': "drtgvc",
	'g': "ftyhbv", 'h': "gyujnb", 'i': "ujko", 'j': "huikmn", 'k': "jiolm", 'l': "kop",
	'm': "njk", 'n': "bhjm", 'o': "iklp", 'p': "ol", 'q': "wa", 'r': "edft", 's': "awedxz",
	't': "rfgy", 'u': "yhji", 'v': "cfgb", 'w': "qase", 'x': "zsdc", 'y': "tghu", 'z': "asx",
}

// injectNoise applies the calibrated lexical/orthographic noise to text.
// It returns the (possibly) altered text and whether the comment was selected;
// a selected comment is guaranteed to differ from the input.
func injectNoise(text string, cfg NoiseConfig, rng *rand.Rand) (string, bool) {
	if rng.Float64() >= cfg.Rate {
		return text, false
	}
	out := mapWords(text, func(w string) string {
		if short, ok := abbreviations[strings.ToLower(w)]; ok && rng.Float64() < cfg.Abbrev {
			return short
		}
		return w
	})
	out = mapWords(out, func(w string) string {
		if len([]rune(w)) >= 4 && rng.Float64() < cfg.Typo {
			return typo(w, rng)
		}
		return w
	})
	out = dropTildes(out, cfg.TildeDrop, rng)
	if rng.Float64() < cfg.Punct {
		out = irregularPunct(out, rng)
	}
	if rng.Float64() < cfg.CaseJitter {
		out = jitterCase(out, rng)
	}
	if out == text {
		out = forceChange(text)
	}
	return out, true
}

// mapWords applies f to every maximal run of letters, keeping everything else intact.
func mapWords(s string, f func(string) string) string {
	var b strings.Builder
	var word []rune
	flush := func() {
		if len(word) > 0 {
			b.WriteString(f(string(word)))
			word = word[:0]
		}
	}
	for _, r := range s {
		if unicode.IsLetter(r) {
			word = append(word, r)
			continue
		}
		flush()
		b.WriteRune(r)
	}
	flush()
	return b.String()
}

func dropTildes(s string, p float64, rng *rand.Rand) string {
	rs := []rune(s)
	for i, r := range rs {
		if base, ok := tildeFold[r]; ok && rng.Float64() < p {
			rs[i] = base
		}
	}
	return string(rs)
}

// typo applies one keyboard-level error to a word: neighbouring key,
// transposition, deletion or duplication. The word always changes.
func typo(w string, rng *rand.Rand) string {
	rs := []rune(w)
	pos := rng.Intn(len(rs) - 1) // leave room for the transposition partner
	switch rng.Intn(4) {
	case 0:
		if nb, ok := qwertyNeighbors[rs[pos]]; ok {
			rs[pos] = rune(nb[rng.Intn(len(nb))])
			return string(rs)
		}
	case 1:
		if rs[pos] != rs[pos+1] {
			rs[pos], rs[pos+1] = rs[pos+1], rs[pos]
			return string(rs)
		}
	case 2:
		if len(rs) > 2 {
			out := append(append([]rune{}, rs[:pos]...), rs[pos+1:]...)
			if string(out) != w {
				return string(out)
			}
		}
	}
	// duplication: also the fallback when the chosen operation was a no-op
	out := append(append(append([]rune{}, rs[:pos+1]...), rs[pos]), rs[pos+1:]...)
	return string(out)
}

func splitTerminal(s string) (core, term string) {
	core = strings.TrimRight(s, ".!?")
	return core, s[len(core):]
}

func irregularPunct(s string, rng *rand.Rand) string {
	core, term := splitTerminal(s)
	var out string
	switch rng.Intn(4) {
	case 0:
		out = core
		if term == "" {
			out = core + "..."
		}
	case 1:
		out = core + "!!!"
	case 2:
		out = core + "..."
	default:
		out = strings.ReplaceAll(s, ",", "")
		if out == s {
			out = core + "..."
		}
	}
	if out == s {
		out = core + "!!!"
	}
	return out
}

func jitterCase(s string, rng *rand.Rand) string {
	switch rng.Intn(3) {
	case 0:
		return strings.ToUpper(s)
	case 1:
		return strings.ToLower(s)
	}
	toggled := false
	out := mapWords(s, func(w string) string {
		if rng.Float64() < 0.5 {
			rs := []rune(w)
			if unicode.IsUpper(rs[0]) {
				rs[0] = unicode.ToLower(rs[0])
			} else {
				rs[0] = unicode.ToUpper(rs[0])
			}
			toggled = true
			return string(rs)
		}
		return w
	})
	if !toggled {
		return strings.ToLower(s)
	}
	return out
}

// forceChange guarantees a selected comment differs from its input.
func forceChange(s string) string {
	rs := []rune(s)
	for i, r := range rs {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				rs[i] = unicode.ToLower(r)
			} else {
				rs[i] = unicode.ToUpper(r)
			}
			if string(rs) != s {
				return string(rs)
			}
			break
		}
	}
	return s + "..."
}
