# ADR-006 — Uniqueness of the global `frecuencias` row

**Status:** Accepted (Phase 1). Deviation from baseline Annex I DDL.

## Context

`frecuencias` is declared with `PRIMARY KEY (lexema_id, corrida_id, categoria_id)`, and `categoria_id IS NULL`
denotes the *global* frequency of a lexeme in a run. SQLite treats NULLs as distinct inside a PRIMARY KEY
and UNIQUE constraint, so the key does not stop two global rows for the same `(lexema_id, corrida_id)`.
That would break invariant I-2 (`global = Σ per-category`, queried with
`GROUP BY l.id, f.corrida_id … HAVING global != suma_categorias`), which presumes exactly one global row.
Verified by mutation: with the index removed, a second `(1, 1, NULL, 3)` insert is accepted and the test fails.

## Decision

Add a partial unique index:

```sql
CREATE UNIQUE INDEX uq_frecuencias_global ON frecuencias(lexema_id, corrida_id) WHERE categoria_id IS NULL;
```

The composite primary key stays unchanged and keeps protecting the per-category rows.

## Alternatives rejected

- **Sentinel value** (e.g. `categoria_id = 0` for global): restores a plain composite key, but abandons
  the NULL-as-global semantics that the baseline dictionary and the I-2 query rely on, and puts a value in a
  column that is logically a foreign key to `corpus.categorias(id)` that does not exist there.
- **`COALESCE` expression index over all rows:** equivalent for the global row but also re-expresses
  the per-category uniqueness the primary key already provides.
- **Separate table for global frequencies:** splits one relation in two and complicates the I-2 query.

## Consequences

- The schema keeps NULL-as-global and now rejects a second global row with a `UNIQUE constraint failed` error.
- The Phase 2 sweep must write at most one global row per `(lexema, corrida)`; it already must, by I-2.
- Verified in `internal/db/schema_test.go`: `TestFrecuenciasGlobalRowIsUnique` (duplicate global rejected;
  global rows for another lexeme or another run accepted) and `TestFrecuenciasPerCategoryRowsKeepCompositePK`
  (per-category duplicates still rejected, a global row coexists with per-category rows). Removing the index
  makes the first test fail, confirming that it exercises the fix.
