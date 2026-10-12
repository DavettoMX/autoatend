"""Test-set metrics (TE §5.2.2): Macro-F1 (CPM01), recall of Reconocimientos (CPM02),
per-class precision/recall/F1 and the row-normalized confusion matrix."""

from __future__ import annotations

import numpy as np
from sklearn.metrics import confusion_matrix, f1_score, precision_recall_fscore_support

from .data import CATEGORIES, RECONOCIMIENTOS

LABELS = [c[0] for c in CATEGORIES]


def evaluate(y_true: np.ndarray, y_pred: np.ndarray) -> dict:
    prec, rec, f1, sup = precision_recall_fscore_support(y_true, y_pred, labels=LABELS, zero_division=0)
    cm = confusion_matrix(y_true, y_pred, labels=LABELS)
    norm = confusion_matrix(y_true, y_pred, labels=LABELS, normalize="true")
    por_clase = []
    for i, (cid, nombre) in enumerate(CATEGORIES):
        por_clase.append(
            {
                "id": cid,
                "nombre": nombre,
                "precision": float(prec[i]),
                "recall": float(rec[i]),
                "f1": float(f1[i]),
                "soporte": int(sup[i]),
            }
        )
    return {
        "n": int(len(y_true)),
        "accuracy": float((y_true == y_pred).mean()),
        "macro_f1": float(f1_score(y_true, y_pred, labels=LABELS, average="macro")),
        "recall_reconocimientos": float(rec[LABELS.index(RECONOCIMIENTOS)]),
        "por_clase": por_clase,
        "confusion_conteos": [[int(v) for v in row] for row in cm],
        "confusion_normalizada": [[float(v) for v in row] for row in norm],
    }
