package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Options configures one gencorpus run.
type Options struct {
	Config     Config
	OutPath    string // corpus.db to create (must not exist)
	ReportPath string // validation report; empty = do not write
	generate   func(Config) ([]Comment, error)
}

// Result is everything a run measured.
type Result struct {
	Validation Validation
	Persisted  PersistedStats
}

// ValidationError is returned, and nothing is persisted, when the battery fails.
type ValidationError struct{ Failed []Check }

func (e *ValidationError) Error() string {
	var parts []string
	for _, c := range e.Failed {
		parts = append(parts, fmt.Sprintf("%s (%s)", c.Name, c.Detail))
	}
	return "gencorpus: validation battery failed: " + strings.Join(parts, "; ")
}

// run is the whole pipeline: generate -> validate -> persist -> verify -> report.
func run(opts Options) (Result, error) {
	gen := opts.generate
	if gen == nil {
		gen = Generate
	}
	var res Result
	cs, err := gen(opts.Config)
	if err != nil {
		return res, err
	}
	res.Validation = Validate(cs, opts.Config)
	if !res.Validation.OK() {
		var failed []Check
		for _, c := range res.Validation.Checks {
			if !c.Pass {
				failed = append(failed, c)
			}
		}
		return res, &ValidationError{Failed: failed}
	}
	if err := Persist(opts.OutPath, cs); err != nil {
		return res, err
	}
	if res.Persisted, err = verifyCorpus(opts.OutPath); err != nil {
		return res, err
	}
	if opts.ReportPath != "" {
		if err := os.MkdirAll(filepath.Dir(opts.ReportPath), 0o755); err != nil {
			return res, err
		}
		f, err := os.Create(opts.ReportPath)
		if err != nil {
			return res, err
		}
		defer f.Close()
		if err := writeReport(f, opts.Config, res); err != nil {
			return res, err
		}
	}
	return res, nil
}
