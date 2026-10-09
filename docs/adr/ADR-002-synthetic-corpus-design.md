# ADR-002 — Design of the synthetic corpus

**Status:** Accepted (Phase 1)

## Context

TelArg's real history (≈73,000 labeled comments: 100/day × 365 × 2 years) is not available, so the
corpus must be synthesized. It must be reproducible (invariant I-3), statistically distinguishable per
category (so a classifier can learn) but not trivially separable (so the metrics are not inflated), and
noisy enough to justify the per-lexeme regular expression of Reading B.

## Decision (`cmd/gencorpus`)

- **Seed:** `rand.New(rand.NewSource(10433))`. The stream is consumed in a fixed order: one shuffle of the
  category sequence, then one composition + noise pass per comment in chronological order.
- **Volume and classes:** 73,000 comments; Falta de Servicio 30 % (21,900), Cobros Indebidos / Problemas de
  Facturación / Quejas 20 % each (14,600), Reconocimientos 10 % (7,300). The imbalance is deliberate (it
  motivates Macro-F1 and class-weighting in Phase 3). Counts are exact, not sampled.
- **Time axis:** exactly 100 comments per day over 730 days, 2024-08-30 to 2026-08-29; ids follow time.
  This makes the temporal split meaningful and leaves ties within a day, broken by id.
- **Templates:** slot-typed templates in Mexican Spanish: 36 / 34 / 34 / 34 / 32 per category (the baseline
  requires 30–60), plus 10 neutral openers and 10 neutral closers shared by all categories. A comment is an
  optional opener (p = 0.25), one template, a second template of the same category (p = 0.40), and an
  optional closer (p = 0.30).
- **Slots and vocabularies:** plans, amounts (24 base values, 6 surface formats, `{MONTO_B}` guaranteed
  different from the preceding amount), ticket/folio numbers (6 prefixes, 5–7 digits), cities, charge
  concepts, months, service areas, agent names, devices, colloquialisms, channels, waiting times, praise
  phrases. Unknown or unterminated slots are errors, never silently emitted.
- **Calibrated noise** (`DefaultNoise`): 30 % of comments are selected; a selected comment is guaranteed to
  change. Within them: tilde dropping (0.80 per accented letter, `ñ` preserved), whole-word abbreviations
  `q`/`x`/`tel`/`fact` (0.50 per word), QWERTY-neighbour/transposition/deletion/duplication typos (0.04 per
  word of 4+ letters), case jitter (0.25), irregular punctuation (0.30).
- **Pre-persist validation battery** (the corpus is written only if all pass):
  1. per-category counts equal the targets exactly;
  2. no empty or whitespace-only text; 3. no unfilled `{SLOT}`;
  4. well-formed, chronological dates with exactly 100 comments per day;
  5. realized noise rate within ±0.015 of the configured 0.30 (measured: 0.2994);
  6. bounded lexical overlap: Jaccard of the 100 most frequent tokens per category, for every pair of
     distinct categories, within [0.15, 0.80] (below: trivially separable; above: indistinguishable).
     These bounds were fixed before the first run; the measured range is 0.220–0.342;
  7. sample review: 50 comments per category (own rng stream, seed + 1) must pass sanity checks (10–800
     runes, ≥ 2 words, no slot residue) and are printed in the report for manual reading.
- **Split:** the 80/20 temporal-stratified assignment is materialized once in `corpus.db` (`particiones`,
  see `internal/db/build/split.go`) and never recomputed downstream.

## Consequences and known limitations

- The vocabulary is generated, not observed: reported metrics measure the classifier against this
  generator, not against real TelArg traffic.
- Sentences inside one comment are independent templates of the same category, so inter-sentence
  coherence is loose (visible in the report's sample review). It does not affect bag-of-words features.
- **Duplicates:** 16.25 % of the 73,000 texts are exact duplicates of another comment, and 19.5 % of the
  14,600 test comments (2,848) have an identical same-category text in the train partition. A
  temporal split does not prevent this, so test-set metrics in Phase 3 can be optimistic. The design is
  left unchanged here; Phase 3 should report this figure next to CPM01/CPM02.
- Overlap is checked on the top-100 tokens only; the discriminative power of rarer tokens is not bounded.
