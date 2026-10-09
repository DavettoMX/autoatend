GO     ?= go
CORPUS := data/corpus.db
DIST   := dist

.PHONY: gen clean-data test test-integration cover-gate lint-regexp cross sweep train assets

# Generates data/corpus.db and docs/reports/generator_validation_report.md.
# Fails if data/corpus.db already exists (the corpus is built once).
gen:
	mkdir -p data docs/reports
	$(GO) run ./cmd/gencorpus

clean-data:
	rm -f $(CORPUS)

test:
	$(GO) test -short -count=1 ./...

test-integration:
	$(GO) test -count=1 ./...

# Fails if any package with statements is below 80.0% coverage.
# cmd/... (main packages) and internal/testutil are excluded.
cover-gate:
	@out=$$($(GO) test -count=1 -cover ./...) || { echo "$$out"; echo "cover-gate: tests failed"; exit 1; }; \
	echo "$$out"; \
	echo "$$out" | awk ' \
	  /coverage: [0-9.]+% of statements/ { \
	    for (i = 1; i <= NF; i++) if ($$i ~ /^autoatend\//) pkg = $$i; \
	    if (pkg ~ /\/cmd\// || pkg ~ /\/internal\/testutil$$/) next; \
	    for (i = 1; i <= NF; i++) if ($$i ~ /%$$/) pct = $$i; \
	    sub(/%/, "", pct); \
	    n++; \
	    if (pct + 0 < 80.0) { printf "cover-gate: FAIL %s %.1f%% (< 80.0%%)\n", pkg, pct + 0; fail = 1 } \
	    else printf "cover-gate: ok   %s %.1f%%\n", pkg, pct + 0 \
	  } \
	  END { if (n == 0) { print "cover-gate: no packages with statements"; exit 1 }; if (fail) { exit 1 }; print "cover-gate: PASS" }'

# Forbids the stdlib regexp package and cgo: the matcher is hand-built, and the binary must stay CGO-free.
lint-regexp:
	@hits=$$(grep -rln '"regexp"' --include='*.go' . || true); \
	if [ -n "$$hits" ]; then echo 'lint-regexp: forbidden import "regexp" in:'; echo "$$hits"; exit 1; fi; \
	hits=$$(grep -rln 'import "C"' --include='*.go' . || true); \
	if [ -n "$$hits" ]; then echo 'lint-regexp: forbidden cgo (import "C") in:'; echo "$$hits"; exit 1; fi; \
	echo "lint-regexp: ok"

# Static, CGO-free cross-builds of the offline generator.
cross:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/gencorpus.exe ./cmd/gencorpus
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/gencorpus-macos-arm64 ./cmd/gencorpus
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/autoatend.exe ./cmd/autoatend
#	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags="-s -w" -o $(DIST)/autoatend-macos-arm64 ./cmd/autoatend

sweep:
	@echo "pendiente: Fase 2"; exit 1

train:
	@echo "pendiente: Fase 3"; exit 1

assets:
	@echo "pendiente: Fase 4"; exit 1
