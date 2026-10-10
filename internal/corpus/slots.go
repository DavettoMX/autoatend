package corpus

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"autoatend/internal/corpus/templates"
)

// slotFiller expands {SLOT} placeholders using one rng stream. It remembers the
// last two amounts drawn so that {MONTO_B} is guaranteed to differ from the
// amount before it (e.g. "me cobraron {MONTO} y debía ser {MONTO_B}").
type slotFiller struct {
	rng        *rand.Rand
	amount     int
	prevAmount int
}

func newSlotFiller(rng *rand.Rand) *slotFiller { return &slotFiller{rng: rng} }

func (f *slotFiller) fill(tpl string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(tpl); {
		c := tpl[i]
		if c != '{' {
			b.WriteByte(c)
			i++
			continue
		}
		end := strings.IndexByte(tpl[i:], '}')
		if end < 0 {
			return "", fmt.Errorf("gencorpus: unterminated slot in template %q", tpl)
		}
		name := tpl[i+1 : i+end]
		v, ok := f.value(name)
		if !ok {
			return "", fmt.Errorf("gencorpus: unknown slot {%s} in template %q", name, tpl)
		}
		b.WriteString(v)
		i += end + 1
	}
	return b.String(), nil
}

func (f *slotFiller) value(name string) (string, bool) {
	switch name {
	case "MONTO":
		return f.drawAmount(false), true
	case "MONTO_B":
		return f.drawAmount(true), true
	case "FOLIO":
		return f.folio(), true
	}
	vals, ok := templates.SlotVocab[name]
	if !ok {
		return "", false
	}
	return vals[f.rng.Intn(len(vals))], true
}

func (f *slotFiller) drawAmount(differ bool) string {
	n := templates.AmountPool[f.rng.Intn(len(templates.AmountPool))]
	for differ && n == f.amount {
		n = templates.AmountPool[f.rng.Intn(len(templates.AmountPool))]
	}
	f.prevAmount, f.amount = f.amount, n
	return formatAmount(n, f.rng.Intn(amountStyles))
}

// Surface variety of the amount and folio slots.
const (
	amountStyles   = 6 // surface formats of formatAmount
	folioMinDigits = 5
	folioMaxDigits = 7
)

// formatAmount renders n pesos in one of amountStyles surface formats ("variable format").
func formatAmount(n, style int) string {
	s := withThousands(n)
	switch style {
	case 0:
		return "$" + s
	case 1:
		return "$" + s + ".00"
	case 2:
		return s + " pesos"
	case 3:
		return "MXN " + s
	case 4:
		return "$" + s + " MXN"
	default:
		return "$" + s + ".50"
	}
}

func withThousands(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	return s[:len(s)-3] + "," + s[len(s)-3:]
}

func (f *slotFiller) folio() string {
	prefix := templates.FolioPrefixes[f.rng.Intn(len(templates.FolioPrefixes))]
	digits := folioMinDigits + f.rng.Intn(folioMaxDigits-folioMinDigits+1)
	lo := 1
	for i := 1; i < digits; i++ {
		lo *= 10
	}
	return prefix + strconv.Itoa(lo+f.rng.Intn(9*lo))
}
