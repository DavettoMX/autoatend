package ml_test

import (
	"bytes"
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/ml"
)

const resultadosPath = "../../data/resultados.db"

func TestLexicon_MapsFormsToCanonical(t *testing.T) {
	lex, err := ml.NewLexicon([]string{"factura", "cobraron", "excelente"})
	if err != nil {
		t.Fatal(err)
	}
	if lex.Size() != 3 {
		t.Errorf("Size = %d, want 3", lex.Size())
	}
	// "facturación" and "facturas" share the stem "factur" with "factura"; "cobro" has
	// stem "cobr" like "cobraron"; numbers, stopwords and unknown stems disappear.
	got := lex.Lexemes("Me cobro la facturación, facturas 2 veces; excelente xyzzy 123")
	want := []string{"cobraron", "factura", "factura", "excelente"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Lexemes = %v, want %v", got, want)
	}
	if got := lex.Lexemes(""); got == nil || len(got) != 0 {
		t.Errorf("Lexemes(\"\") = %#v, want empty non-nil slice", got)
	}
}

func TestLexicon_RejectsStemCollision(t *testing.T) {
	// "factura" and "facturas" have the same stem, so a catalogue cannot hold both.
	if _, err := ml.NewLexicon([]string{"factura", "facturas"}); err == nil {
		t.Error("NewLexicon accepted two canonical forms with the same stem")
	}
}

func TestLoadCatalogue_MatchesResultadosDB(t *testing.T) {
	lex, err := ml.LoadCatalogue(resultadosPath)
	if err != nil {
		t.Fatal(err)
	}
	d, err := db.OpenReadOnly(resultadosPath)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var n int
	if err := d.QueryRow(`SELECT COUNT(*) FROM lexemas`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if lex.Size() != n || n == 0 {
		t.Errorf("Size = %d, lexemas rows = %d", lex.Size(), n)
	}
	// Every canonical form of the catalogue maps to itself.
	rows, err := d.Query(`SELECT forma_canonica FROM lexemas`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	checkSelfMapping(t, lex, rows)
}

func checkSelfMapping(t *testing.T, lex *ml.Lexicon, rows *sql.Rows) {
	t.Helper()
	bad := 0
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			t.Fatal(err)
		}
		got := lex.Lexemes(c)
		if len(got) != 1 || got[0] != c {
			// A canonical form that is a stopword cannot occur: the sweep drops those.
			bad++
			if bad <= 3 {
				t.Errorf("canonical %q maps to %v", c, got)
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if bad > 0 {
		t.Errorf("%d canonical forms do not map to themselves", bad)
	}
}

func TestRunTokenizar_JSONLines(t *testing.T) {
	lex, err := ml.NewLexicon([]string{"factura", "cobraron", "excelente"})
	if err != nil {
		t.Fatal(err)
	}
	in := strings.Join([]string{
		`{"id":7,"texto":"Me cobraron la factura 2 veces"}`,
		`{"id":"a-1","texto":"el de la y"}`,
		``,
		`{"id":9,"texto":"¡Excelente! excelente"}`,
	}, "\n") + "\n"
	var out bytes.Buffer
	if err := ml.RunTokenizar(strings.NewReader(in), &out, lex); err != nil {
		t.Fatal(err)
	}
	want := `{"id":7,"lexemas":["cobraron","factura"]}` + "\n" +
		`{"id":"a-1","lexemas":[]}` + "\n" +
		`{"id":9,"lexemas":["excelente","excelente"]}` + "\n"
	if out.String() != want {
		t.Errorf("output:\n%s\nwant:\n%s", out.String(), want)
	}
}

func TestRunTokenizar_BadLine(t *testing.T) {
	lex, _ := ml.NewLexicon([]string{"factura"})
	var out bytes.Buffer
	err := ml.RunTokenizar(strings.NewReader("{\"id\":1,\"texto\":\"factura\"}\nno es json\n"), &out, lex)
	if err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Errorf("err = %v, want an error naming line 2", err)
	}
}
