package lexer

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func TestAnalyze_ReturnsSurfaceForms(t *testing.T) {
	got := Analyze("Me cobraron de más por las Facturas del mes")
	want := []string{"me", "cobraron", "más", "facturas", "mes"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAnalyze_FiltersDe(t *testing.T) {
	if got := Analyze("el cobro de la factura"); contains(got, "de") {
		t.Errorf("Analyze kept %q: %v", "de", got)
	}
}

func TestAnalyze_KeepsHomographsAndContentWords(t *testing.T) {
	got := Analyze("Estoy bajo contrato, para el servicio como sobre entre más solo aun")
	for _, w := range []string{"bajo", "para", "servicio", "como", "sobre", "entre", "más", "solo", "aun"} {
		if !contains(got, w) {
			t.Errorf("Analyze dropped %q: %v", w, got)
		}
	}
}

func TestAnalyze_DropsNumbers(t *testing.T) {
	got := Analyze("Plan Max 699 cuesta 1,499.00 pesos, ticket 254463")
	want := []string{"plan", "max", "cuesta", "pesos", "ticket"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAnalyze_NoiseSpellings(t *testing.T) {
	// "q" (que) and "x" (por) are filtered; the unaccented "segun" too.
	got := Analyze("no se q paso x favor segun el contrato")
	want := []string{"no", "se", "paso", "favor", "contrato"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCPE01_OnlyStopwordsGivesEmpty(t *testing.T) {
	if got := Analyze("de la y el que en"); len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
	if got := Analyze("¡Ah, eh! ¿Y el de la?"); len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestCPE02_EmptyTextGivesEmptyNoError(t *testing.T) {
	for _, in := range []string{"", "   ", "\n\t", "...!!!", "123 456"} {
		if got := Analyze(in); len(got) != 0 {
			t.Errorf("Analyze(%q) = %v, want empty", in, got)
		}
	}
}

func TestCPE03_EmojiAndURL(t *testing.T) {
	got := Analyze("Pésimo servicio \U0001F621\U0001F621 https://telarg.mx/quejas?id=99 llevo días")
	want := []string{"pésimo", "servicio", "https", "telarg", "mx", "quejas", "id", "llevo", "días"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCPE04_LongTextUnder500ms(t *testing.T) {
	base := "Me cobraron de más en la factura del mes pasado y nadie me resuelve. "
	text := strings.Repeat(base, 5000/len(base)+1)[:5000]
	start := time.Now()
	got := Analyze(text)
	if d := time.Since(start); d > 500*time.Millisecond {
		t.Errorf("5,000 characters took %v, want < 500ms", d)
	}
	if len(got) == 0 {
		t.Error("no tokens from a 5,000-character text")
	}
}

func TestAnalyze_AcceptsDecomposedInput(t *testing.T) {
	got := Analyze("Facturacio\u0301n Sen\u0303al")
	want := []string{"facturación", "señal"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
