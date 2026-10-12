package corpus

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"autoatend/internal/corpus/templates"
)

// Coherence regression tests (TE §6.1.4: a manual review of the sample must find the corpus plausible).
// Every check works on generated text with its own parsing; none reuses the slot filler's internal state.

// planPrice is the independent oracle of plan prices: the number in the plan name, and the fixed
// price documented for the plans without a number.
var planPrice = map[string]float64{
	"Libre 199": 199, "Amigo Sin Límite": 299, "Total 399": 399, "Ahorro 99": 99, "Max 699": 699,
	"Familiar 549": 549, "Prepago Básico": 99, "Joven 149": 149, "Empresarial 899": 899,
	"Premium 999": 999, "Navega 249": 249, "Nocturno 79": 79, "Social 299": 299, "Hogar 449": 449,
}

var monthNames = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

// cleanTexts composes n noise-free comments of one category with a dedicated rng.
func cleanTexts(t *testing.T, cat, n int, seed int64) []string {
	t.Helper()
	rng := rand.New(rand.NewSource(seed))
	out := make([]string, n)
	for i := range out {
		s, _, err := compose(cat, rng)
		if err != nil {
			t.Fatal(err)
		}
		out[i] = s
	}
	return out
}

// firstAmount parses the first number of s ("$1,499.00", "MXN 199", "199 pesos").
func firstAmount(s string) (float64, bool) {
	i := strings.IndexAny(s, "0123456789")
	if i < 0 {
		return 0, false
	}
	j := i
	for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == ',' || s[j] == '.') {
		j++
	}
	tok := strings.ReplaceAll(strings.TrimRight(s[i:j], ".,"), ",", "")
	v, err := strconv.ParseFloat(tok, 64)
	return v, err == nil
}

// amountAfter parses the first amount that follows marker in s.
func amountAfter(s, marker string) (float64, bool) {
	i := strings.Index(s, marker)
	if i < 0 {
		return 0, false
	}
	return firstAmount(s[i+len(marker):])
}

// literalSegments splits a template on its {SLOT}s and keeps the lowercased literal pieces.
func literalSegments(tpl string) []string {
	var segs []string
	for len(tpl) > 0 {
		i := strings.IndexByte(tpl, '{')
		if i < 0 {
			segs = append(segs, tpl)
			break
		}
		segs = append(segs, tpl[:i])
		j := strings.IndexByte(tpl[i:], '}')
		if j < 0 {
			break
		}
		tpl = tpl[i+j+1:]
	}
	out := segs[:0]
	for _, s := range segs {
		if len(strings.TrimSpace(s)) >= 3 {
			out = append(out, strings.ToLower(s))
		}
	}
	return out
}

// matchesTemplate reports whether lowered contains, in order, every literal piece of tpl.
func matchesTemplate(lowered, tpl string) bool {
	pos := 0
	for _, seg := range literalSegments(tpl) {
		i := strings.Index(lowered[pos:], seg)
		if i < 0 {
			return false
		}
		pos += i + len(seg)
	}
	return true
}

func TestCoherence_BaseTemplatesAreIdentifiable(t *testing.T) {
	for ci, ts := range templates.Templates {
		for _, tpl := range ts {
			n := 0
			for _, s := range literalSegments(tpl) {
				n += len(s)
			}
			if n < 15 {
				t.Errorf("category %d: template %q has only %d literal characters, too few to identify it", ci+1, tpl, n)
			}
		}
	}
}

func TestCoherence_NoCommentHasTwoBaseTemplates(t *testing.T) {
	for cat := 1; cat <= 5; cat++ {
		for _, text := range cleanTexts(t, cat, 4000, int64(100+cat)) {
			lowered := strings.ToLower(text)
			var hits []string
			for _, tpl := range templates.Templates[cat-1] {
				if matchesTemplate(lowered, tpl) {
					hits = append(hits, tpl)
				}
			}
			if len(hits) != 1 {
				t.Fatalf("category %d: comment %q matches %d base templates, want exactly 1: %q", cat, text, len(hits), hits)
			}
		}
	}
}

func TestCoherence_PlanPriceAndChargeAmounts(t *testing.T) {
	var cuesta, cobrando, factura, extra, cargo int
	for _, cat := range []int{2, 3} {
		for _, text := range cleanTexts(t, cat, 20000, int64(200+cat)) {
			for name, price := range planPrice {
				if i := strings.Index(text, "plan "+name+" cuesta "); i >= 0 {
					cuesta++
					rest := text[i+len("plan "+name+" cuesta "):]
					got, ok := firstAmount(rest)
					if !ok || got != price {
						t.Fatalf("plan %q costs %v in %q, want its price %v", name, got, text, price)
					}
					if j := strings.Index(rest, "cobrando "); j >= 0 {
						cobrando++
						charge, ok := firstAmount(rest[j+len("cobrando "):])
						if !ok || charge <= price || charge > 2*price {
							t.Fatalf("charge %v for plan %q (price %v) must be in (price, 2*price]: %q", charge, name, price, text)
						}
					}
					if j := strings.Index(rest, "monto total es de "); j >= 0 {
						factura++
						total, ok := firstAmount(rest[j+len("monto total es de "):])
						if !ok || total <= price {
							t.Fatalf("invoice total %v must exceed plan price %v: %q", total, price, text)
						}
					}
				}
				if strings.Contains(text, "de más desde que cambié de plan a "+name) || strings.Contains(text, "cambié de plan a "+name+" me cobran ") {
					extra++
					var v float64
					var ok bool
					if i := strings.Index(text, "cambié de plan a "+name+" me cobran "); i >= 0 {
						v, ok = amountAfter(text[i:], "me cobran ")
					} else {
						v, ok = amountAfter(text, "Me cobran ")
					}
					if !ok || v <= 0 || v >= price {
						t.Fatalf("extra charge %v after switching to %q (price %v) must be below the price: %q", v, name, price, text)
					}
				}
				if strings.Contains(text, "cancelar mi plan "+name) {
					cargo++
					v, ok := amountAfter(text, "cargaron ")
					if ok && v != price {
						t.Fatalf("charge %v after cancelling %q must equal its price %v: %q", v, name, price, text)
					}
				}
			}
		}
	}
	for what, n := range map[string]int{"cuesta": cuesta, "cobrando": cobrando, "factura": factura, "extra": extra, "cargo": cargo} {
		if n < 50 {
			t.Errorf("only %d %q occurrences in 40000 comments: the check would be vacuous", n, what)
		}
	}
}

func TestCoherence_GeneratedCorpusPriceBelowCharge(t *testing.T) {
	cs, err := Generate(smallConfig(10433))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, c := range cs {
		i := strings.Index(c.Texto, " cuesta ")
		if i < 0 || !strings.Contains(c.Texto[i:], "cobrando ") {
			continue
		}
		n++
		price, ok1 := amountAfter(c.Texto[i:], " cuesta ")
		charge, ok2 := amountAfter(c.Texto[i:], "cobrando ")
		if !ok1 || !ok2 || charge <= price {
			t.Fatalf("category %d: charge %v is not above price %v in %q", c.CategoriaID, charge, price, c.Texto)
		}
	}
	if n == 0 {
		t.Error("no cuesta/cobrando comment in the generated corpus: the check is vacuous")
	}
}

// fillOf returns the text between marker and the first occurrence of stop after it.
func fillOf(text, marker, stop string) (string, bool) {
	i := strings.Index(text, marker)
	if i < 0 {
		return "", false
	}
	rest := text[i+len(marker):]
	j := strings.Index(rest, stop)
	if j < 0 {
		return "", false
	}
	return rest[:j], true
}

func inSet(s string, set ...string) bool {
	for _, v := range set {
		if s == v {
			return true
		}
	}
	return false
}

// TestCoherence_SlotsAgreeInsideOneSentence checks that each constrained slot only takes values that read
// naturally in its frame. Every frame must occur (otherwise the check is vacuous).
func TestCoherence_SlotsAgreeInsideOneSentence(t *testing.T) {
	type frame struct {
		cat          int
		marker, stop string
		ok           func(fill string) bool
	}
	numUpTo := func(max int) func(string) bool {
		return func(f string) bool { n, err := strconv.Atoi(f); return err == nil && n <= max }
	}
	frames := map[string]frame{
		"internet lento": {1, "como no tener ", ".", func(f string) bool { return inSet(f, "internet", "datos móviles") }},
		"equipo con pantalla": {1, "", " dice sin servicio", func(f string) bool {
			return strings.HasSuffix(f, "mi celular") || strings.HasSuffix(f, "mi teléfono")
		}},
		"espera corta": {5, "en menos de ", " me solucionaron", func(f string) bool {
			return !strings.Contains(f, "más de") && !strings.Contains(f, "casi") && f != "una hora"
		}},
		"traslados":           {4, "Me transfirieron ", " veces", numUpTo(5)},
		"respuestas":          {4, "me dieron ", " respuestas", numUpTo(5)},
		"viajes en balde":     {4, "ya van ", " veces que voy", numUpTo(5)},
		"meses pagando":       {2, "ya llevo ", " meses", numUpTo(12)},
		"factura tarde":       {3, "Recibo mi factura ", " días tarde", numUpTo(10)},
		"factura electronica": {3, "Solicité factura electrónica ", " y me siguen", func(f string) bool { return f != "en papel" }},
		"cancelacion":         {2, "Cancelé ", " hace meses", func(f string) bool { return f != "reconexión" && f != "roaming internacional" }},
		"documentos":          {3, "mandé los documentos ", ".", func(f string) bool { return f != "por teléfono" }},
		"pago":                {3, "No aparece el pago que hice ", " en mi estado", func(f string) bool { return f != "por correo" && f != "por chat" }},
	}
	seen := map[string]int{}
	for _, cat := range []int{1, 2, 3, 4, 5} {
		for _, text := range cleanTexts(t, cat, 20000, int64(300+cat)) {
			for name, fr := range frames {
				if fr.cat != cat {
					continue
				}
				var fill string
				var ok bool
				if fr.marker == "" {
					// subject of "dice sin servicio": the text before it up to the previous sentence break
					i := strings.Index(text, fr.stop)
					if i < 0 {
						continue
					}
					head := text[:i]
					if k := strings.LastIndexAny(head, ".!?"); k >= 0 {
						head = strings.TrimSpace(head[k+1:])
					}
					fill, ok = strings.ToLower(head), true
				} else {
					fill, ok = fillOf(text, fr.marker, fr.stop)
				}
				if !ok {
					continue
				}
				seen[name]++
				if !fr.ok(fill) {
					t.Fatalf("frame %q: unnatural fill %q in %q", name, fill, text)
				}
			}
		}
	}
	for name := range frames {
		if seen[name] < 20 {
			t.Errorf("frame %q seen %d times: the check is vacuous (template wording changed?)", name, seen[name])
		}
	}
}

func TestCoherence_ZoneFramesAgreeWithTheSentence(t *testing.T) {
	indoor := []string{"mi casa", "mi oficina", "mi trabajo", "mi domicilio"}
	var antenna, street int
	for _, text := range cleanTexts(t, 1, 30000, 410) {
		if fill, ok := fillOf(text, "cambiaron la antena en ", " perdí"); ok {
			antenna++
			if inSet(fill, indoor...) {
				t.Fatalf("an antenna is not changed inside %q: %q", fill, text)
			}
		}
		if i := strings.Index(text, " nunca hay señal, tengo que salir a la calle"); i >= 0 {
			street++
			head := text[:i]
			if k := strings.LastIndex(head, "En "); k >= 0 {
				head = head[k+3:]
			}
			if inSet(strings.ToLower(head), "toda mi calle", "el centro", "mi barrio", "toda la zona", "mi colonia") {
				t.Fatalf("cannot step out to the street when the whole area has no signal: %q", text)
			}
		}
		if i := strings.Index(text, "no registra en la red desde "); i >= 0 {
			rest := text[i+len("no registra en la red desde "):]
			for _, m := range monthNames {
				if strings.HasPrefix(rest, m) {
					t.Fatalf("a month followed by a duration in hours/days contradicts itself: %q", text)
				}
			}
		}
	}
	if antenna < 20 || street < 20 {
		t.Errorf("frames seen antenna=%d street=%d: the check is vacuous", antenna, street)
	}
}

func TestCoherence_OddWordingIsGone(t *testing.T) {
	for _, tpl := range templates.Templates[1] {
		if strings.Contains(tpl, "nunca me avisaron que se activaba") {
			t.Errorf("template still has the odd wording: %q", tpl)
		}
	}
	// The same-day double charge names a single amount, not two different ones for the same thing.
	seen := 0
	for _, text := range cleanTexts(t, 2, 20000, 420) {
		seg, ok := fillOf(text, "Cobro indebido: me descontaron ", " por lo mismo")
		if !ok {
			continue
		}
		seen++
		runs := 0
		for i := 0; i < len(seg); i++ {
			if seg[i] >= '0' && seg[i] <= '9' && (i == 0 || !(seg[i-1] >= '0' && seg[i-1] <= '9' || seg[i-1] == ',' || seg[i-1] == '.')) {
				runs++
			}
		}
		if runs != 1 {
			t.Fatalf("the double charge names %d amounts, want 1: %q", runs, text)
		}
	}
	if seen < 20 {
		t.Errorf("double-charge template seen %d times: the check is vacuous", seen)
	}
}

func TestCoherence_RecognitionsAreAlwaysPositive(t *testing.T) {
	negative := []string{
		"necesito ayuda", "exijo", "queja", "favor", "avísenme", "avisenme", "novedad", "pronta respuesta",
		"respuesta pronto", "atento a su respuesta", "quedo al pendiente", "ya basta", "ay no", "qué barbaridad", "chale",
		"híjole", "molest", "frustr", "inaceptable", "gracias de antemano", "tengo un problema", "problema con",
	}
	for _, text := range cleanTexts(t, 5, 20000, 500) {
		lowered := strings.ToLower(text)
		for _, w := range negative {
			if strings.Contains(lowered, w) {
				t.Fatalf("a Reconocimientos comment contains the non-positive %q: %q", w, text)
			}
		}
		if strings.Contains(lowered, "fue muy atento y profesional") || strings.Contains(lowered, "un verdadero profesional") {
			t.Fatalf("masculine adjective bound to a name of any gender: %q", text)
		}
	}
}

func TestCoherence_ComplaintCategoriesNeverPraise(t *testing.T) {
	positive := []string{"felicit", "excelente", "comentario positivo", "buena noticia", "los recomiendo", "gracias por resolver", "me encantó"}
	for cat := 1; cat <= 4; cat++ {
		for _, text := range cleanTexts(t, cat, 8000, int64(600+cat)) {
			lowered := strings.ToLower(text)
			for _, w := range positive {
				if strings.Contains(lowered, w) {
					t.Fatalf("category %d comment contains the praise %q: %q", cat, w, text)
				}
			}
		}
	}
}

func TestCoherence_ComplementsElaborateTheSameProblem(t *testing.T) {
	for ci, cs := range templates.Complements {
		if len(cs) < 10 || len(cs) > 15 {
			t.Errorf("category %d has %d complements, want 10..15", ci+1, len(cs))
		}
		seen := map[string]bool{}
		for _, c := range cs {
			if seen[c] {
				t.Errorf("category %d: duplicate complement %q", ci+1, c)
			}
			seen[c] = true
			if strings.ContainsAny(c, "{}$0123456789") {
				t.Errorf("category %d: complement %q introduces a slot, an amount or a figure", ci+1, c)
			}
			if !strings.HasSuffix(c, ".") {
				t.Errorf("category %d: complement %q is not a closed sentence", ci+1, c)
			}
			// closed-class homographs retained by ADR-003 would create new L-1 contexts
			for _, w := range strings.FieldsFunc(strings.ToLower(c), func(r rune) bool { return r == ' ' || r == ',' || r == '.' }) {
				if inSet(w, "para", "como", "cómo", "entre", "bajo", "sobre") {
					t.Errorf("category %d: complement %q uses the L-1 homograph %q", ci+1, c, w)
				}
			}
			// a complement is never a second incident: no concrete topic appears in it
			lowered := strings.ToLower(c)
			for _, w := range []string{"roaming", "premium", "seguro", "suscripción", "antena", "sms"} {
				if strings.Contains(lowered, w) {
					t.Errorf("category %d: complement %q names a concrete topic (%q)", ci+1, c, w)
				}
			}
		}
	}
}
