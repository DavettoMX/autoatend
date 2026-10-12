"""TF-IDF features with vocabulary and IDF computed from the TRAIN partition only.

Tokenization is done by the Go tokenizer (cmd/tokenizar) invoked as a
subprocess, once, in batch (TE §6.3.6, CPA06). Everything numeric that must be
bit-identical in Go (term frequency x idf, L2 norm) is done here with plain
float64 operations in pure Python, summing in ascending index order. No numpy
reductions are used for those, because numpy sums pairwise.

IDF formula (smooth): idf_i = ln((1 + N) / (1 + df_i)) + 1, with N the number of
training documents and df_i the number of training documents that contain
lexeme i (both read from df_train.json). The values are exported as computed
here; Go never recomputes them.
"""

from __future__ import annotations

import json
import math
import subprocess
from dataclasses import dataclass
from pathlib import Path

import numpy as np
from scipy import sparse

IDF_FORMULA = "ln((1+N)/(1+df))+1"


def tokenize_batch(binary: str | Path, catalogue: str | Path, ids: list[int], texts: list[str]) -> list[list[str]]:
    """Runs cmd/tokenizar once over all texts; returns the lexemes of each text, in input order."""
    payload = "".join(
        json.dumps({"id": i, "texto": t}, ensure_ascii=False) + "\n" for i, t in zip(ids, texts)
    )
    proc = subprocess.run(
        [str(binary), "-catalogo", str(catalogue)],
        input=payload.encode("utf-8"),
        capture_output=True,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(f"tokenizar failed ({proc.returncode}): {proc.stderr.decode('utf-8', 'replace')}")
    out: list[list[str]] = []
    for n, line in enumerate(proc.stdout.decode("utf-8").splitlines()):
        rec = json.loads(line)
        if rec["id"] != ids[n]:
            raise RuntimeError(f"tokenizar out of order at line {n + 1}: id {rec['id']} != {ids[n]}")
        out.append(rec["lexemas"])
    if len(out) != len(ids):
        raise RuntimeError(f"tokenizar returned {len(out)} lines for {len(ids)} texts")
    return out


@dataclass
class Vocabulary:
    lexemes: list[str]  # canonical lexemes, index order
    idf: list[float]
    df: list[int]
    n_train: int
    cutoff: str
    index: dict[str, int]


def load_vocabulary(df_path: str | Path) -> Vocabulary:
    with open(df_path, encoding="utf-8") as f:
        d = json.load(f)
    n = d["documentos_train"]
    lexemes, df = [], []
    for e in d["df"]:
        if e["df"] <= 0:
            raise ValueError(f"df_train.json lists {e['lexema']!r} with df={e['df']}")
        lexemes.append(e["lexema"])
        df.append(e["df"])
    idf = [math.log((1 + n) / (1 + c)) + 1.0 for c in df]
    return Vocabulary(lexemes, idf, df, n, d["corte"], {lx: i for i, lx in enumerate(lexemes)})


def vectorize(lexemes: list[str], vocab: Vocabulary) -> tuple[list[int], list[float]]:
    """Raw term counts x idf, L2-normalized. Mirrors Model.Vectorize in internal/ml/predict.go."""
    counts: dict[int, int] = {}
    for lx in lexemes:
        i = vocab.index.get(lx)
        if i is not None:
            counts[i] = counts.get(i, 0) + 1
    idxs = sorted(counts)
    vals = [float(counts[i]) * vocab.idf[i] for i in idxs]
    s = 0.0
    for v in vals:
        s += v * v
    if not idxs:
        return [], []
    norm = math.sqrt(s)
    return idxs, [v / norm for v in vals]


def to_matrix(docs: list[tuple[list[int], list[float]]], n_features: int) -> sparse.csr_matrix:
    indptr = [0]
    indices: list[int] = []
    data: list[float] = []
    for idxs, vals in docs:
        indices.extend(idxs)
        data.extend(vals)
        indptr.append(len(indices))
    return sparse.csr_matrix(
        (np.asarray(data, dtype=np.float64), np.asarray(indices, dtype=np.int32), np.asarray(indptr, dtype=np.int64)),
        shape=(len(docs), n_features),
    )


def document_frequency(token_lists: list[list[str]], mask: list[bool]) -> dict[str, int]:
    """Document frequency over the selected documents (cross-check of df_train.json)."""
    df: dict[str, int] = {}
    for toks, keep in zip(token_lists, mask):
        if keep:
            for lx in set(toks):
                df[lx] = df.get(lx, 0) + 1
    return df
