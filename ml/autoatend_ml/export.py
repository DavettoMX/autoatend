"""Export of a trained model to modelo.json, and the pure-Python reference scorer.

Both families are exported as one LINEAR scorer so that Go has a single
inference path (internal/ml):

    z_k = b_k + sum_{i in supp(x)} w_{k,i} * x_i        softmax(z)

  nbm  w = feature_log_prob_, b = class_log_prior_ (softmax of the joint
       log-likelihood is exactly the Naive Bayes posterior, i.e. predict_proba).
  rls  w = coef_, b = intercept_.

The JSON is deterministic: fixed key order, floats written with repr (shortest
round-trip), no timestamps. Weights are sparse per class: [[index, value], ...]
with the non-zero entries only, in ascending index order.
"""

from __future__ import annotations

import json
import math
from pathlib import Path

import numpy as np
import sklearn

from .data import CATEGORIES
from .features import IDF_FORMULA, Vocabulary
from .models import LR_MAX_ITER, LR_TOL, SEED, Fitted

FORMAT_VERSION = 1


def linear_parameters(fitted: Fitted) -> tuple[np.ndarray, np.ndarray]:
    est = fitted.estimator
    if list(est.classes_) != [c[0] for c in CATEGORIES]:
        raise ValueError(f"unexpected classes_ {est.classes_}")
    if fitted.familia == "nbm":
        return est.feature_log_prob_, est.class_log_prior_
    return est.coef_, est.intercept_


def sparse_weights(W: np.ndarray) -> list[list[list]]:
    out = []
    for row in W:
        nz = np.flatnonzero(row)
        out.append([[int(i), float(row[i])] for i in nz])
    return out


def hyperparameters(fitted: Fitted) -> dict:
    if fitted.familia == "nbm":
        return {"alpha": float(fitted.params["alpha"]), "suavizado": "Laplace/Lidstone", "fit_prior": True}
    return {
        "C": float(fitted.params["C"]),
        "penalizacion": "l2",
        "solver": "saga",
        "multinomial": True,
        "max_iter": LR_MAX_ITER,
        "tol": LR_TOL,
        "iteraciones": fitted.n_iter,
        "convergio": fitted.converged,
        "random_state": SEED,
    }


def model_dict(fitted: Fitted, vocab: Vocabulary, metricas: dict, df_cero: int) -> dict:
    W, b = linear_parameters(fitted)
    return {
        "version": FORMAT_VERSION,
        "familia": fitted.familia,
        "estrategia": fitted.estrategia,
        "hiperparametros": hyperparameters(fitted),
        "categorias": [{"id": cid, "nombre": name} for cid, name in CATEGORIES],
        "tf": "conteo crudo del lexema en el comentario",
        "idf_formula": IDF_FORMULA,
        "normalizacion": "l2",
        "documentos_train": vocab.n_train,
        "corte_train": vocab.cutoff,
        "lexemas_df_cero_excluidos": df_cero,
        "scorer": "z_k = sesgo_k + sum_i peso_{k,i} * x_i; softmax con resta del maximo",
        "sklearn": sklearn.__version__,
        "vocabulario": vocab.lexemes,
        "idf": vocab.idf,
        "pesos": sparse_weights(W),
        "sesgo": [float(v) for v in b],
        "metricas": metricas,
    }


def dumps_model(m: dict) -> str:
    """Deterministic text: one top-level key per line, weights one class per line."""

    def enc(v) -> str:
        return json.dumps(v, ensure_ascii=False, separators=(",", ":"), allow_nan=False)

    lines = []
    for key, val in m.items():
        if key == "pesos":
            body = "[\n" + ",\n".join(enc(row) for row in val) + "\n]"
        else:
            body = enc(val)
        lines.append(f"{enc(key)}:{body}")
    return "{\n" + ",\n".join(lines) + "\n}\n"


def write_text(path: str | Path, text: str) -> None:
    Path(path).parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(text)


class Scorer:
    """Pure-Python mirror of Model.PredictSparse (internal/ml/predict.go), built from the
    exported dict so that it tests what was exported. Plain float64 operations, ascending
    index order, softmax with the maximum subtracted."""

    def __init__(self, model: dict):
        v = len(model["vocabulario"])
        self.bias = list(model["sesgo"])
        self.w = [[0.0] * v for _ in model["pesos"]]
        for k, row in enumerate(model["pesos"]):
            for i, val in row:
                self.w[k][i] = val

    def proba(self, idxs: list[int], vals: list[float]) -> list[float]:
        z = []
        for k, b in enumerate(self.bias):
            s = b
            wk = self.w[k]
            for i, x in zip(idxs, vals):
                s += wk[i] * x
            z.append(s)
        mx = max(z)
        e = [math.exp(s - mx) for s in z]
        tot = 0.0
        for v in e:
            tot += v
        return [v / tot for v in e]

    @staticmethod
    def argmax(p: list[float]) -> int:
        best = 0
        for k, v in enumerate(p):
            if v > p[best]:
                best = k
        return best + 1  # categoria_id
