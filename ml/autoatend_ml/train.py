"""Offline training pipeline of AutoAtend (Bloque C, TE §5.1.3, §5.2.2, §6.4.2).

    uv run --project ml python -m autoatend_ml.train

Reads data/corpus.db (read-only), vectorizes with the Go tokenizer, compares two
families under three imbalance strategies with 5-fold stratified CV on train,
refits the winner of each family on the whole train, evaluates ONCE on the
held-out test set and writes data/modelo.json, data/modelo_alternativo.json,
data/gold5.json, testdata/paridad_1000.json and the report. Nothing is tuned on
test, and every output is deterministic (fixed seeds, n_jobs=1).
"""

from __future__ import annotations

import argparse
import json
import sys
import time
from pathlib import Path

import numpy as np

from . import data, export, features, models
from .evaluate import evaluate
from .report import build_report

PARITY_N = 1000
GOLD_MIN_CONF = 0.80


def log(msg: str) -> None:
    print(f"[train] {msg}", file=sys.stderr, flush=True)


def pick_best(cells: list[dict]) -> dict:
    """Highest mean CV Macro-F1; ties go to the smaller std, then to the earlier grid position."""
    return min(cells, key=lambda c: (-c["f1_media"], c["f1_desv"], c["orden"]))


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--corpus", default="data/corpus.db")
    ap.add_argument("--catalogo", default="data/resultados.db")
    ap.add_argument("--df", default="data/df_train.json")
    ap.add_argument("--tokenizar", default="dist/tokenizar", help="binary built from ./cmd/tokenizar")
    ap.add_argument("--out", default="data", help="directory for modelo.json, modelo_alternativo.json, gold5.json")
    ap.add_argument("--testdata", default="testdata", help="directory for paridad_1000.json")
    ap.add_argument("--report", default="docs/reports/evaluacion_modelo.md")
    args = ap.parse_args(argv)
    t0 = time.time()

    # ---- data, vocabulary, features ---------------------------------------------------------
    corpus = data.load_corpus(args.corpus)
    vocab = features.load_vocabulary(args.df)
    if corpus.cutoff != vocab.cutoff or sum(corpus.is_train) != vocab.n_train:
        raise SystemExit(
            f"split mismatch: corpus cutoff {corpus.cutoff}/{sum(corpus.is_train)} train vs "
            f"df_train.json {vocab.cutoff}/{vocab.n_train}"
        )
    log(f"split: corte {corpus.cutoff}, train {vocab.n_train}, test {len(corpus.ids) - vocab.n_train}")
    toks = features.tokenize_batch(args.tokenizar, args.catalogo, corpus.ids, corpus.texts)
    n_catalogue = data.catalogue_size(args.catalogo)
    df_check = features.document_frequency(toks, corpus.is_train)
    mism = sum(1 for lx in vocab.lexemes if df_check.get(lx, 0) != vocab.df[vocab.index[lx]])
    extra = sum(1 for lx in df_check if lx not in vocab.index)
    if mism or extra:
        raise SystemExit(f"df_train.json disagrees with the tokenizer output: {mism} mismatches, {extra} extra")
    df_cero = n_catalogue - len(vocab.lexemes)
    log(f"vocabulary {len(vocab.lexemes)} lexemes, catalogue {n_catalogue}, df=0 excluded: {df_cero}")

    docs = [features.vectorize(t, vocab) for t in toks]
    X = features.to_matrix(docs, len(vocab.lexemes))
    y = np.asarray(corpus.labels)
    tr = np.asarray(corpus.is_train)
    X_tr, y_tr, X_te, y_te = X[tr], y[tr], X[~tr], y[~tr]
    te_pos = np.flatnonzero(~tr)  # positions in the corpus, in id order
    te_docs = [docs[p] for p in te_pos]
    te_texts = [corpus.texts[p] for p in te_pos]
    te_ids = [corpus.ids[p] for p in te_pos]
    train_texts = {corpus.texts[p] for p in np.flatnonzero(tr)}
    nodup = np.asarray([t not in train_texts for t in te_texts])
    n_empty = int(sum(1 for idxs, _ in te_docs if not idxs))
    log(f"test duplicates of train: {int((~nodup).sum())} of {len(nodup)}; test with empty vector: {n_empty}")

    # ---- 5-fold stratified CV on train ------------------------------------------------------
    cells: list[dict] = []
    grid: list[dict] = []
    order = 0
    for fam in models.FAMILIES:
        for strat in models.STRATEGIES:
            results = []
            for params in models.GRIDS[fam]:
                r = models.cross_validate(fam, strat, params, X_tr, y_tr)
                r["orden"] = order
                order += 1
                results.append(r)
                grid.append(r)
                log(f"cv {fam}/{strat} {params}: {r['f1_media']:.5f} +- {r['f1_desv']:.5f}")
            cells.append(pick_best(results))
    winner_cell = pick_best(cells)
    other = [c for c in cells if c["familia"] != winner_cell["familia"]]
    alt_cell = pick_best(other)
    log(f"winner {winner_cell['familia']}/{winner_cell['estrategia']} {winner_cell['hiperparametros']}")
    log(f"alternative {alt_cell['familia']}/{alt_cell['estrategia']} {alt_cell['hiperparametros']}")

    def cv_summary(cell: dict) -> dict:
        return {k: cell[k] for k in ("familia", "estrategia", "hiperparametros", "f1_media", "f1_desv", "f1_folds", "convergio")}

    cv_table = [cv_summary(c) for c in cells]

    # ---- refit on the whole train, evaluate once on test ------------------------------------
    finals = {}
    for tag, cell in (("principal", winner_cell), ("alternativo", alt_cell)):
        fitted = models.fit_model(cell["familia"], cell["estrategia"], cell["hiperparametros"], X_tr, y_tr, models.SEED)
        pred = fitted.estimator.predict(X_te)
        ev = evaluate(y_te, pred)
        ev_nd = evaluate(y_te[nodup], pred[nodup])
        metricas = {
            "cv_elegido": cv_summary(cell),
            "cv": cv_table,
            "test": ev,
            "test_sin_duplicados": ev_nd,
            "duplicados_test": {"n_test": int(len(nodup)), "n_duplicados": int((~nodup).sum()), "n_sin_duplicados": int(nodup.sum())},
        }
        mdl = export.model_dict(fitted, vocab, metricas, df_cero)
        scorer = export.Scorer(mdl)
        probs = [scorer.proba(i, v) for i, v in te_docs]
        sk_probs = fitted.estimator.predict_proba(X_te)
        mirror_cls = np.asarray([scorer.argmax(p) for p in probs])
        max_diff = float(np.abs(np.asarray(probs) - sk_probs).max())
        agree = float((mirror_cls == pred).mean())
        log(f"{tag}: test macro-F1 {ev['macro_f1']:.5f}, recall Reconocimientos {ev['recall_reconocimientos']:.5f}, "
            f"sin duplicados {ev_nd['macro_f1']:.5f}/{ev_nd['recall_reconocimientos']:.5f}; "
            f"linear scorer vs predict_proba max diff {max_diff:.2e}, class agreement {agree:.6f}")
        finals[tag] = {
            "fitted": fitted, "model": mdl, "scorer": scorer, "probs": probs, "sk_probs": sk_probs,
            "pred_sklearn": pred, "pred_mirror": mirror_cls, "eval": ev, "eval_nd": ev_nd,
            "linealidad": {"max_dif_prob": max_diff, "acuerdo_clase": agree},
        }

    # ---- artifacts ----------------------------------------------------------------------------
    out = Path(args.out)
    export.write_text(out / "modelo.json", export.dumps_model(finals["principal"]["model"]))
    export.write_text(out / "modelo_alternativo.json", export.dumps_model(finals["alternativo"]["model"]))

    rng = np.random.default_rng(models.SEED)
    sel = np.sort(rng.choice(len(te_ids), size=PARITY_N, replace=False))
    casos = []
    for j in sel:
        idxs, vals = te_docs[j]
        rec = {"id": te_ids[j], "texto": te_texts[j], "vector": [[i, v] for i, v in zip(idxs, vals)]}
        for tag in ("principal", "alternativo"):
            f = finals[tag]
            rec[tag] = {
                "prob": f["probs"][j],
                "prob_sklearn": [float(v) for v in f["sk_probs"][j]],
                "clase": int(f["pred_mirror"][j]),
            }
        casos.append(rec)
    parity = {"semilla": models.SEED, "n": PARITY_N, "modelo": finals["principal"]["model"]["familia"], "casos": casos}
    export.write_text(
        Path(args.testdata) / "paridad_1000.json",
        '{"semilla":%d,"n":%d,"modelo":"%s","casos":[\n' % (models.SEED, PARITY_N, parity["modelo"])
        + ",\n".join(json.dumps(c, ensure_ascii=False, separators=(",", ":"), allow_nan=False) for c in casos)
        + "\n]}\n",
    )

    gold, demo = pick_gold_and_demo(finals["principal"], te_ids, te_texts, y_te, nodup, vocab, te_docs)
    export.write_text(out / "gold5.json", json.dumps(gold, ensure_ascii=False, indent=1, allow_nan=False) + "\n")

    ctx = {
        "corpus": corpus, "vocab": vocab, "n_catalogue": n_catalogue, "df_cero": df_cero,
        "n_test": int(len(te_ids)), "n_dup": int((~nodup).sum()), "n_nodup": int(nodup.sum()), "n_empty_test": n_empty,
        "grid": grid, "cells": cells, "winner": winner_cell, "alt": alt_cell, "finals": finals,
        "gold": gold, "demo": demo, "parity_n": PARITY_N, "df_check": (mism, extra),
    }
    export.write_text(args.report, build_report(ctx))
    log(f"done in {time.time() - t0:.1f} s")
    return 0


def pick_gold_and_demo(final: dict, ids, texts, y_te, nodup, vocab, te_docs):
    """gold5: one test comment per category, correct, confidence >= 0.80, not a train duplicate,
    readable length; the choice among eligible ones is random with the fixed seed.
    Demo candidates (NOT in gold5): a clear Falta de Servicio and the borderline Reconocimiento."""
    probs = final["probs"]
    pred = final["pred_mirror"]
    conf = np.asarray([p[int(c) - 1] for p, c in zip(probs, pred)])
    readable = np.asarray([30 <= len(t) <= 120 for t in texts])
    rng = np.random.default_rng(models.SEED)
    gold = []
    chosen = set()
    for cid, name in data.CATEGORIES:
        ok = np.flatnonzero((y_te == cid) & (pred == cid) & (conf >= GOLD_MIN_CONF) & nodup & readable)
        if len(ok) == 0:
            raise SystemExit(f"no eligible gold comment for category {cid}")
        j = int(rng.choice(ok))
        chosen.add(j)
        gold.append({"id": ids[j], "texto": texts[j], "categoria_id": cid, "categoria": name, "confianza_python": float(conf[j])})

    def describe(j: int) -> dict:
        p = probs[j]
        order = sorted(range(5), key=lambda k: (-p[k], k))
        return {
            "id": ids[j], "texto": texts[j], "categoria_real": data.CATEGORIES[int(y_te[j]) - 1][1],
            "prediccion": data.CATEGORIES[int(pred[j]) - 1][1], "confianza": float(p[order[0]]),
            "segunda": data.CATEGORIES[order[1]][1], "confianza_segunda": float(p[order[1]]),
            "lexemas": [vocab.lexemes[i] for i in te_docs[j][0]],
        }

    clear = np.flatnonzero((y_te == 1) & (pred == 1) & (conf >= 0.97) & nodup & readable)
    clear = [j for j in clear if j not in chosen]
    j_clear = int(rng.choice(clear))
    rec = np.flatnonzero((y_te == data.RECONOCIMIENTOS) & (pred == data.RECONOCIMIENTOS) & nodup & readable)
    rec = [j for j in rec if j not in chosen]
    j_border = int(min(rec, key=lambda j: (conf[j], j)))
    return gold, {"claro_falta_de_servicio": describe(j_clear), "frontera_reconocimiento": describe(j_border)}


if __name__ == "__main__":
    raise SystemExit(main())
