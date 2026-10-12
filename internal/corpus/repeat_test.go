package corpus

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
	"unicode"

	"autoatend/internal/corpus/templates"
)

// repeatIgnore are words that may legitimately repeat inside one sentence.
var repeatIgnore = map[string]bool{
	"mensajes": true, "llamadas": true, "internet": true, "senal": true, "servicio": true,
	"pesos": true, "total": true, // the unit of two amounts, and "Total 399" next to "monto total"
}

// contentWords are the lowercased words of s with at least 5 letters, minus the ignored ones.
func contentWords(s string) []string {
	var out []string
	s = strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u").Replace(strings.ToLower(s))
	for _, w := range strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) }) {
		if len([]rune(w)) >= 5 && !repeatIgnore[w] {
			out = append(out, w)
		}
	}
	return out
}

// slotValues expands tpl like the slot filler, returning the literal text and each slot's value.
func slotValues(t *testing.T, tpl string, rng *rand.Rand) (literal string, vals []string) {
	t.Helper()
	f := newSlotFiller(rng)
	var lit strings.Builder
	for i := 0; i < len(tpl); {
		if tpl[i] != '{' {
			lit.WriteByte(tpl[i])
			i++
			continue
		}
		end := strings.IndexByte(tpl[i:], '}')
		v, ok := f.value(tpl[i+1 : i+end])
		if !ok {
			t.Fatalf("template %q: unknown slot %s", tpl, tpl[i:i+end+1])
		}
		vals = append(vals, v)
		i += end + 1
	}
	return lit.String(), vals
}

// TestCoherence_NoWordRepeatedBetweenSlotsOrLiteral: inside one sentence, two slots (or a slot and the fixed
// text) must not say the same place or concept twice ("Un diez para la sucursal: resolvieron mi duda en la sucursal").
func TestCoherence_NoWordRepeatedBetweenSlotsOrLiteral(t *testing.T) {
	bad := map[string]bool{}
	for ci, ts := range templates.Templates {
		for _, tpl := range ts {
			rng := rand.New(rand.NewSource(int64(900 + ci)))
			for k := 0; k < 600; k++ {
				lit, vals := slotValues(t, tpl, rng)
				litWords := map[string]bool{}
				for _, w := range contentWords(lit) {
					litWords[w] = true
				}
				seen := map[string]int{}
				for i, v := range vals {
					for _, w := range contentWords(v) {
						if litWords[w] || (seen[w] > 0 && seen[w] != i+1) {
							bad["cat "+string(rune('1'+ci))+": "+tpl+" -> "+w] = true
						}
						seen[w] = i + 1
					}
				}
			}
		}
	}
	var msgs []string
	for m := range bad {
		msgs = append(msgs, m)
	}
	sort.Strings(msgs)
	for _, m := range msgs {
		t.Errorf("repeated word: %s", m)
	}
}

// TestCoherence_NumberAgreementAndSpecificFrames covers the residual defects found in the second manual review.
func TestCoherence_NumberAgreementAndSpecificFrames(t *testing.T) {
	// "ya van/pasaron una semana" disagrees in number; the verb must work for every duration.
	for ci, ts := range templates.Templates {
		for _, tpl := range ts {
			if strings.Contains(tpl, "{DURACION}") || strings.Contains(tpl, "{DURACION_DIAS}") {
				for _, bad := range []string{"ya van {DURACION", "ya pasaron {DURACION", "van {DURACION", "pasaron {DURACION", "pasó {DURACION"} {
					if strings.Contains(tpl, bad) {
						t.Errorf("category %d: %q forces a verb number on a singular/plural duration", ci+1, tpl)
					}
				}
			}
			// plural slot values need a plural verb and a neutral pronoun
			if strings.Contains(tpl, "Se pierde {SERVICIO}") || strings.Contains(tpl, "y no lo he recuperado") {
				t.Errorf("category %d: %q does not agree with plural services (datos móviles)", ci+1, tpl)
			}
		}
	}
	for _, text := range cleanTexts(t, 1, 30000, 1100) {
		for _, w := range []string{"ya van una semana", "ya pasaron una semana", "ya pasaron más de 48 horas", "ya van más de 48 horas", "ya pasaron varias horas", "se pierde datos móviles"} {
			if strings.Contains(strings.ToLower(text), w) {
				t.Fatalf("number disagreement %q in %q", w, text)
			}
		}
	}
	// "Me suscribieron a suscripción de música"
	for _, tpl := range templates.Templates[1] {
		if strings.Contains(tpl, "suscribieron a {") {
			t.Errorf("%q: 'suscribieron a' + a subscription noun reads as a subscription to a subscription", tpl)
		}
	}
	for _, text := range cleanTexts(t, 2, 30000, 1101) {
		if strings.Contains(text, "suscribieron a suscripción") || strings.Contains(text, "suscribieron a membresía") {
			t.Fatalf("a subscription to a subscription: %q", text)
		}
	}
	// "dos recibos casi juntos, uno de X" leaves the other receipt unsaid
	seen := 0
	for _, text := range cleanTexts(t, 3, 30000, 1102) {
		if i := strings.Index(text, "dos recibos casi juntos"); i >= 0 {
			seen++
			rest := text[i+len("dos recibos casi juntos"):]
			if strings.HasPrefix(rest, ", uno de") || !strings.Contains(rest, "cada uno") {
				t.Fatalf("incomplete thought about the two receipts: %q", text)
			}
		}
	}
	if seen < 20 {
		t.Errorf("the two-receipts frame was seen %d times: the check is vacuous", seen)
	}
}

// TestCoherence_ComplementsFitEveryBaseTemplateOfTheirCategory: a complement is attached to any base template of its
// category, so it must not assume a channel (a phone call, a visit), a person, a document, a recurrence, a promised
// date or a kind of service that only some templates have ("Es un gusto que existan personas así" after
// "Excelente cobertura en Villahermosa…").
func TestCoherence_ComplementsFitEveryBaseTemplateOfTheirCategory(t *testing.T) {
	assumptions := map[string][]string{
		"channel":          {"llam", "teléfono", "chat", "sucursal", "acudir", "visita", "correo"},
		"person":           {"persona", "atendi", "trato", "amable", "agente"},
		"document":         {"documento", "recibo", "factura", "datos"},
		"recurrence":       {"cada mes", "cada vez", "periodo", "otra vez", "siempre"},
		"promised date":    {"fecha", "hoy mismo"},
		"specific service": {"fácil", "rápid", "cobertura", "señal", "aplicación", "promoci"},
		"own purchase":     {"no pedí", "contraté"},
	}
	for ci, cs := range templates.Complements {
		for _, c := range cs {
			lowered := strings.ToLower(c)
			for kind, words := range assumptions {
				for _, w := range words {
					if strings.Contains(lowered, w) {
						t.Errorf("category %d: complement %q assumes a %s (%q), which not every base template has", ci+1, c, kind, w)
					}
				}
			}
		}
	}
}
