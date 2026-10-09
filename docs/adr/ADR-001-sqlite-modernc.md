# ADR-001 — SQLite through `modernc.org/sqlite` (pure Go, no CGo)

**Status:** Accepted (Phase 1)

## Context

AutoAtend must run on the exam machine (Windows x86-64, Core i5, 8 GB) with no installation, no
administrator rights, no network and no assumed runtime. A macOS ARM64 laptop is the backup target.
Consigna a) requires a relational database for the corpus (`corpus.db`) and another for the lexical
catalog (`resultados.db`), both in third normal form.

## Decision

- Store both databases as SQLite files, accessed through `modernc.org/sqlite` **v1.60.1** (pinned in `go.mod`).
- Build with `CGO_ENABLED=0`. The driver is a transpiled-to-Go SQLite, so no C toolchain, DLL or ABI
  concern exists on either target and one module cross-compiles with plain `go build`.
- Two files, not one: `corpus.db` is the historical input and is never modified by the analyzer;
  `resultados.db` is the simulator output and carries run metadata (baseline §6.1.1).
- Fixed build-time PRAGMAs (`internal/db/build/open.go`): `page_size=4096`, `encoding='UTF-8'`,
  `journal_mode=DELETE`, `foreign_keys=ON`; the build pool is pinned to one connection so per-connection
  PRAGMAs always hold. `Finalize` re-asserts `journal_mode=DELETE` and then runs `VACUUM`.
- `journal_mode=DELETE` replaces the baseline Annex I `journal_mode=WAL` for `corpus.db`
  (deviation recorded in ADR-007): WAL leaves `-wal`/`-shm` sidecars and a WAL-flagged header, which is
  incompatible with embedding the database file as a single image.

## Alternatives rejected

- `mattn/go-sqlite3`: requires CGo and a C compiler on every build host, which breaks the single-static-binary constraint.
- A client/server database (PostgreSQL, MySQL): needs an installation and a network service on the exam machine.
- Flat files / JSON: no declarative constraints; the schema's CHECK, FOREIGN KEY and UNIQUE rules are the
  first line of defence against mislabeled data and are tested directly (`internal/db/build/schema_test.go`).

## Consequences

- Pure-Go SQLite is slower than the C library, but the workload here is small: Phase 1 generates, builds,
  splits and VACUUMs a 73,000-row `corpus.db` (≈16 MB). Measured on the development machine, the full-size
  integration test (two complete builds plus read-back verification) runs in ≈1.7 s.
- Reproducibility depends on determinism of the write sequence (no timestamps, fixed PRAGMAs, sorted
  iteration, fixed insertion order). Phase 1 already verifies that two builds with the same seed yield a
  byte-identical `corpus.db` (Level 4a test), which is the groundwork for invariant I-3.
- Foreign keys are enforced per connection, so every connection opened by this project sets
  `foreign_keys=ON` through the DSN rather than relying on a one-off statement.
