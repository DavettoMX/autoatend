package sweep

import "testing"

// TE Etapa 6: a generated ER is validated by round-trip (compile it with the
// own engine and require that it accepts every variant) before it is stored.
func TestValidatedER_KeepsAGoodER(t *testing.T) {
	forms := []string{"factura", "facturas", "facturación"}
	er, fallback := validatedER("factura", forms, "factur(a|ación|as)")
	if fallback || er != "factur(a|ación|as)" {
		t.Errorf("good ER: got (%q, %v), want it kept", er, fallback)
	}
}

func TestValidatedER_FallsBackToTheCanonicalForm(t *testing.T) {
	forms := []string{"factura", "facturas"}
	cases := map[string]string{
		"syntax error":        "factur(a|as",
		"rejects a variant":   "factur(a|ación)",
		"accepts only a part": "factur",
	}
	for name, bad := range cases {
		er, fallback := validatedER("factura", forms, bad)
		if !fallback || er != "factura" {
			t.Errorf("%s: got (%q, %v), want the canonical form and fallback", name, er, fallback)
		}
	}
}

func TestValidatedER_FallbackEscapesTheCanonicalForm(t *testing.T) {
	er, fallback := validatedER("a+b", []string{"a+b", "a+bb"}, "(")
	if !fallback || er != `a\+b` {
		t.Errorf("got (%q, %v), want a\\+b", er, fallback)
	}
}
