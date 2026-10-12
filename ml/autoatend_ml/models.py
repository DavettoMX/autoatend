"""Model families, imbalance strategies and cross-validation (TE §5.1.3, §6.4.2).

Families:
  nbm  Multinomial Naive Bayes, Laplace/Lidstone smoothing alpha in {0.01, 0.1, 0.5, 1.0}.
  rls  Multinomial (softmax) logistic regression, L2 penalty with parameter C in {0.1, 1, 10},
       solver='saga' (stochastic average gradient, a stochastic-gradient-family solver).

Strategies for the class imbalance (30/20/20/20/10 %):
  pesos_balanceados  class_weight='balanced' (NB: sample_weight from compute_sample_weight('balanced')).
  sobremuestreo      random oversampling of the minority class (Reconocimientos), done here with
                     numpy, deterministic, applied only to the training part of each fold.
  linea_base         no adjustment.
"""

from __future__ import annotations

import warnings
from dataclasses import dataclass

import numpy as np
from scipy import sparse
from sklearn.exceptions import ConvergenceWarning
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import f1_score
from sklearn.model_selection import StratifiedKFold
from sklearn.naive_bayes import MultinomialNB
from sklearn.utils.class_weight import compute_sample_weight

from .data import RECONOCIMIENTOS

SEED = 10433
N_FOLDS = 5
FAMILIES = ("nbm", "rls")
STRATEGIES = ("pesos_balanceados", "sobremuestreo", "linea_base")
GRIDS = {
    "nbm": [{"alpha": a} for a in (0.01, 0.1, 0.5, 1.0)],
    "rls": [{"C": c} for c in (0.1, 1.0, 10.0)],
}
# saga settings: tight enough to converge on the TF-IDF features, fixed for determinism.
LR_MAX_ITER = 300
LR_TOL = 1e-4


def oversample_indices(y: np.ndarray, seed: int, minority: int = RECONOCIMIENTOS) -> np.ndarray:
    """Row indices of the training set after oversampling the minority class.

    All original rows are kept (in order); rows of the minority class are drawn
    with replacement and appended until the class reaches the median class size
    (the 20 % classes). Deterministic for a given seed. Only ever called on the
    training part of a fold, never on validation or test data.
    """
    classes, counts = np.unique(y, return_counts=True)
    target = int(np.median(counts))
    pool = np.flatnonzero(y == minority)
    missing = target - len(pool)
    if missing <= 0:
        return np.arange(len(y))
    rng = np.random.default_rng(seed)
    extra = rng.choice(pool, size=missing, replace=True)
    return np.concatenate([np.arange(len(y)), extra])


@dataclass
class Fitted:
    familia: str
    estrategia: str
    params: dict
    estimator: object
    n_iter: int | None
    converged: bool | None


def fit_model(familia: str, estrategia: str, params: dict, X: sparse.csr_matrix, y: np.ndarray, seed: int) -> Fitted:
    """Fits one (family, strategy, hyperparameters) on the given training data."""
    sample_weight = None
    class_weight = None
    if estrategia == "sobremuestreo":
        idx = oversample_indices(y, seed)
        X, y = X[idx], y[idx]
    elif estrategia == "pesos_balanceados":
        if familia == "nbm":
            sample_weight = compute_sample_weight("balanced", y)
        else:
            class_weight = "balanced"
    elif estrategia != "linea_base":
        raise ValueError(estrategia)

    if familia == "nbm":
        est = MultinomialNB(alpha=params["alpha"])
        est.fit(X, y, sample_weight=sample_weight)
        return Fitted(familia, estrategia, params, est, None, None)

    est = LogisticRegression(
        solver="saga",
        C=params["C"],
        class_weight=class_weight,
        max_iter=LR_MAX_ITER,
        tol=LR_TOL,
        random_state=seed,
    )
    with warnings.catch_warnings():
        warnings.simplefilter("ignore", ConvergenceWarning)
        est.fit(X, y)
    n_iter = int(est.n_iter_.max())
    return Fitted(familia, estrategia, params, est, n_iter, n_iter < LR_MAX_ITER)


def cross_validate(familia: str, estrategia: str, params: dict, X: sparse.csr_matrix, y: np.ndarray) -> dict:
    """5-fold stratified CV (shuffle, random_state=10433); returns Macro-F1 per fold, mean and std."""
    skf = StratifiedKFold(n_splits=N_FOLDS, shuffle=True, random_state=SEED)
    f1s: list[float] = []
    converged = True
    for fold, (tr, va) in enumerate(skf.split(X, y)):
        fitted = fit_model(familia, estrategia, params, X[tr], y[tr], SEED + fold)
        pred = fitted.estimator.predict(X[va])
        f1s.append(float(f1_score(y[va], pred, average="macro")))
        if fitted.converged is False:
            converged = False
    arr = np.asarray(f1s)
    return {
        "familia": familia,
        "estrategia": estrategia,
        "hiperparametros": params,
        "f1_folds": f1s,
        "f1_media": float(arr.mean()),
        "f1_desv": float(arr.std(ddof=1)),
        "convergio": converged,
    }
