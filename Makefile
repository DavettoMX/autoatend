GO     ?= go
CORPUS := data/corpus.db
RESULTS := data/resultados.db
DF     := data/df_train.json
CORPUS_SHA ?= b810bad233f9180484082e056b0aefd7676a63e5c4c87441c3ffc2b309760d71
REPORT := docs/reports/generator_validation_report.md
DIST   := dist
VERSION ?= 1.0.0
ASSETS_DIR := internal/assets/data
ASSET_FILES := corpus.db resultados.db modelo.json modelo_alternativo.json gold5.json

.PHONY: gen clean-data clean-results sweep-sha train-sha test test-integration cover cover-gate lint-regexp cross sweep train assets build fmt-check vet corpus-sha verify

# Generates data/corpus.db and docs/reports/generator_validation_report.md.
# Fails if data/corpus.db already exists (the corpus is built once).
gen:
	$(GO) run ./cmd/gencorpus

clean-data:
	rm -f $(CORPUS)

clean-results:
	rm -f $(RESULTS) $(DF)

test: assets
	$(GO) test -short -count=1 ./...

test-integration: assets
	$(GO) test -count=1 ./...

# Fails if any package with statements is below 80.0% coverage (COVER_STRICT=0 only reports).
# cmd/... (main packages) and internal/testutil are excluded.
COVER_STRICT ?= 1
cover-gate: assets
	@out=$$($(GO) test -count=1 -cover ./...) || { echo "$$out"; echo "cover-gate: tests failed"; exit 1; }; \
	echo "$$out"; \
	echo "$$out" | awk -v strict=$(COVER_STRICT) ' \
	  /coverage: [0-9.]+% of statements/ { \
	    for (i = 1; i <= NF; i++) if ($$i ~ /^autoatend\//) pkg = $$i; \
	    if (pkg ~ /\/cmd\// || pkg ~ /\/internal\/testutil$$/) next; \
	    for (i = 1; i <= NF; i++) if ($$i ~ /%$$/) pct = $$i; \
	    sub(/%/, "", pct); \
	    n++; \
	    if (pct + 0 < 80.0) { printf "cover-gate: FAIL %s %.1f%% (< 80.0%%)\n", pkg, pct + 0; fail = 1 } \
	    else printf "cover-gate: ok   %s %.1f%%\n", pkg, pct + 0 \
	  } \
	  END { if (n == 0) { print "cover-gate: no packages with statements"; exit 1 }; if (fail && strict) { exit 1 }; print (fail ? "cover-gate: below 80% (informational)" : "cover-gate: PASS") }'

# Reports per-package coverage without failing below 80% (demo sprint, PLAN.md §1.3).
cover:
	@$(MAKE) --no-print-directory cover-gate COVER_STRICT=0

# Forbids the stdlib regexp package and cgo: the matcher is hand-built, and the binary must stay CGO-free.
lint-regexp:
	@hits=$$(grep -rln '"regexp"' --include='*.go' . || true); \
	if [ -n "$$hits" ]; then echo 'lint-regexp: forbidden import "regexp" in:'; echo "$$hits"; exit 1; fi; \
	hits=$$(grep -rln 'import "C"' --include='*.go' . || true); \
	if [ -n "$$hits" ]; then echo 'lint-regexp: forbidden cgo (import "C") in:'; echo "$$hits"; exit 1; fi; \
	echo "lint-regexp: ok"

# Static, CGO-free cross-builds of the offline generator.
cross: assets
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/gencorpus.exe ./cmd/gencorpus
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/gencorpus-macos-arm64 ./cmd/gencorpus
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/simulador.exe ./cmd/simulador
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/simulador-macos-arm64 ./cmd/simulador
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/tokenizar.exe ./cmd/tokenizar
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/tokenizar-macos-arm64 ./cmd/tokenizar
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST)/autoatend.exe ./cmd/autoatend
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST)/autoatend-macos-arm64 ./cmd/autoatend

# Sweeps data/corpus.db into data/resultados.db and data/df_train.json.
# Fails if either exists (the catalogue is built once); use make clean-results first.
sweep:
	$(GO) run ./cmd/simulador

# Offline training (Phase 3, consigna e). Python is dev-machine only and never part of the binary.
# Builds the Go tokenizer, trains (uv run, project venv in ml/.venv), writes data/modelo.json,
# data/modelo_alternativo.json, data/gold5.json, testdata/paridad_1000.json and the report, then
# measures the Python<->Go parity (CPA06) in Go and writes it into the report.
ML_ENV := OMP_NUM_THREADS=1 OPENBLAS_NUM_THREADS=1 MKL_NUM_THREADS=1
EVAL_REPORT := docs/reports/evaluacion_modelo.md
train:
	mkdir -p $(DIST) testdata
	$(GO) build -o $(DIST)/tokenizar ./cmd/tokenizar
	$(ML_ENV) uv run --project ml python -m autoatend_ml.train
	$(GO) test -count=1 -v -run TestCPA06 ./internal/ml > $(DIST)/paridad.log || { cat $(DIST)/paridad.log; exit 1; }
	uv run --project ml python -m autoatend_ml.paridad $(EVAL_REPORT) $(DIST)/paridad.log

# Determinism of the training: two independent trainings must write byte-identical artifacts,
# identical to the committed data/modelo.json.
train-sha:
	@d=$$(mktemp -d); trap 'rm -rf "$$d"' EXIT; \
	mkdir -p $(DIST) "$$d/a" "$$d/b"; \
	$(GO) build -o $(DIST)/tokenizar ./cmd/tokenizar || exit 1; \
	for r in a b; do \
	  $(ML_ENV) uv run --project ml python -m autoatend_ml.train --out "$$d/$$r" --testdata "$$d/$$r" --report "$$d/$$r/r.md" 2>/dev/null || exit 1; \
	done; \
	for f in modelo.json modelo_alternativo.json gold5.json paridad_1000.json r.md; do \
	  cmp -s "$$d/a/$$f" "$$d/b/$$f" || { echo "train-sha: FAIL $$f differs between two trainings"; exit 1; }; \
	done; \
	cmp -s "$$d/a/modelo.json" data/modelo.json || { echo "train-sha: FAIL data/modelo.json differs from a fresh training"; exit 1; }; \
	echo "train-sha: ok $$(sha256sum data/modelo.json | cut -d' ' -f1)"

# Copies the artifacts of data/ that the binary embeds into internal/assets/data/ (go:embed cannot
# use "..", TE §7). The copies are never edited by hand and are not versioned (.gitignore); they
# refresh whenever the source in data/ is newer. internal/assets/ui/index.html is the versioned UI source.
$(ASSETS_DIR)/%: data/%
	@mkdir -p $(ASSETS_DIR)
	cp $< $@

assets: $(addprefix $(ASSETS_DIR)/,$(ASSET_FILES))

# The single binary for the host (Linux/macOS/Windows) with the artifacts embedded.
build: assets
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST)/autoatend ./cmd/autoatend

fmt-check:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "fmt-check: unformatted files:"; echo "$$out"; exit 1; fi; echo "fmt-check: ok"

vet: assets
	$(GO) vet ./...

# Regenerates the corpus in a temp dir and requires the baseline SHA-256 and an identical report (invariant I-3 groundwork).
corpus-sha:
	@d=$$(mktemp -d); trap 'rm -rf "$$d"' EXIT; \
	$(GO) run ./cmd/gencorpus -out "$$d/c.db" -report "$$d/r.md" >/dev/null || exit 1; \
	got=$$(sha256sum "$$d/c.db" | cut -d' ' -f1); \
	if [ "$$got" != "$(CORPUS_SHA)" ]; then echo "corpus-sha: FAIL got $$got, want $(CORPUS_SHA)"; exit 1; fi; \
	diff -q "$$d/r.md" $(REPORT) >/dev/null || { echo "corpus-sha: FAIL report differs from $(REPORT)"; exit 1; }; \
	echo "corpus-sha: ok $$got"

# Invariant I-3: two sweeps of the corpus must produce byte-identical files.
sweep-sha:
	@d=$$(mktemp -d); trap 'rm -rf "$$d"' EXIT; \
	$(GO) run ./cmd/simulador -corpus $(CORPUS) -out "$$d/a.db" -df "$$d/a.json" >/dev/null || exit 1; \
	$(GO) run ./cmd/simulador -corpus $(CORPUS) -out "$$d/b.db" -df "$$d/b.json" >/dev/null || exit 1; \
	a=$$(sha256sum "$$d/a.db" | cut -d' ' -f1); b=$$(sha256sum "$$d/b.db" | cut -d' ' -f1); \
	if [ "$$a" != "$$b" ]; then echo "sweep-sha: FAIL $$a != $$b"; exit 1; fi; \
	cmp -s "$$d/a.json" "$$d/b.json" || { echo "sweep-sha: FAIL df_train.json differs"; exit 1; }; \
	echo "sweep-sha: ok $$a"

# Definition of done for a task: everything must pass before asking for a merge.
verify: fmt-check vet lint-regexp test-integration cover corpus-sha sweep-sha cross
	@echo "verify: PASS"
