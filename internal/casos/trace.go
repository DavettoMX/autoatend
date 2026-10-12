package casos

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// caseIDOf returns the case ID encoded in a test name of the form TestCP<L><NN>_<Name>, where
// L is one of U, P, E, A, M (the five families of the plan) and NN two digits: "CPU01" for
// TestCPU01_ParserNominal. It reports false for any other name.
func caseIDOf(name string) (string, bool) {
	rest, ok := strings.CutPrefix(name, "TestCP")
	if !ok || len(rest) < 4 || rest[3] != '_' {
		return "", false
	}
	if !strings.ContainsRune("UPEAM", rune(rest[0])) || rest[1] < '0' || rest[1] > '9' || rest[2] < '0' || rest[2] > '9' {
		return "", false
	}
	return "CP" + rest[:3], true
}

// testFunctions maps every top-level TestCP<L><NN>_* function under root to the files (relative
// to root, slash-separated) that declare it. Hidden directories, dist and the Python project are
// skipped.
func testFunctions(root string) (map[string][]string, error) {
	found := map[string][]string{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != root && (strings.HasPrefix(name, ".") || path == filepath.Join(root, "dist") || path == filepath.Join(root, "ml")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil {
				if _, isCase := caseIDOf(fn.Name.Name); isCase {
					found[fn.Name.Name] = append(found[fn.Name.Name], rel)
				}
			}
		}
		return nil
	})
	return found, err
}

// Trazabilidad checks the two-way traceability between the matrix and the *_test.go files under
// root (the module root): every listed test is declared in the listed file and under the case
// its name says, and every TestCP<L><NN>_* function of the repository is listed in its case. The
// returned problems are sorted; none means the matrix and the code agree.
func Trazabilidad(m *Matriz, root string) ([]string, error) {
	found, err := testFunctions(root)
	if err != nil {
		return nil, err
	}
	var problems []string
	listed := map[string]bool{}
	for _, c := range m.Casos {
		for _, p := range c.Pruebas {
			listed[c.ID+"/"+p.Funcion] = true
			if id, ok := caseIDOf(p.Funcion); !ok || id != c.ID {
				problems = append(problems, fmt.Sprintf("%s lists %s, whose name does not belong to that case", c.ID, p.Funcion))
			}
			files, declared := found[p.Funcion]
			if !declared {
				problems = append(problems, fmt.Sprintf("%s lists %s, which no test file declares", c.ID, p.Funcion))
				continue
			}
			if !contains(files, p.Archivo) {
				problems = append(problems, fmt.Sprintf("%s lists %s in %s, but it is declared in %s", c.ID, p.Funcion, p.Archivo, strings.Join(files, ", ")))
			}
		}
	}
	for name, files := range found {
		id, _ := caseIDOf(name)
		if !listed[id+"/"+name] {
			problems = append(problems, fmt.Sprintf("%s (%s) is not listed under %s in the matrix", name, files[0], id))
		}
	}
	sort.Strings(problems)
	return problems, nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
