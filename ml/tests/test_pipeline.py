import math
import re
from pathlib import Path

import numpy as np

from autoatend_ml import data, features, models

ROOT = Path(__file__).resolve().parents[2]


def test_split_sql_is_a_copy_of_the_go_query():
    go = (ROOT / "internal/db/split.go").read_text(encoding="utf-8")
    m = re.search(r"const splitCutoffSQL = `(.*?)`", go, re.S)
    assert m, "splitCutoffSQL not found in split.go"
    assert m.group(1) == data.SPLIT_CUTOFF_SQL
    assert int(re.search(r"const TrainPercent = (\d+)", go).group(1)) == data.TRAIN_PERCENT


def test_split_matches_df_train_cutoff():
    corpus = data.load_corpus(ROOT / "data/corpus.db")
    vocab = features.load_vocabulary(ROOT / "data/df_train.json")
    assert corpus.cutoff == vocab.cutoff
    assert sum(corpus.is_train) == vocab.n_train


def test_oversample_only_adds_minority_rows_deterministically():
    y = np.array([1] * 30 + [2] * 20 + [3] * 20 + [4] * 20 + [5] * 10)
    a = models.oversample_indices(y, 7)
    b = models.oversample_indices(y, 7)
    assert (a == b).all()
    assert (a[: len(y)] == np.arange(len(y))).all()
    extra = a[len(y):]
    assert len(extra) == 10 and (y[extra] == 5).all()
    assert np.unique(y[a], return_counts=True)[1].tolist() == [30, 20, 20, 20, 20]


def test_vectorize_is_l2_normalized_and_ordered():
    vocab = features.Vocabulary(["a", "b", "c"], [1.5, 2.0, 1.0], [1, 1, 1], 10, "x", {"a": 0, "b": 1, "c": 2})
    idxs, vals = features.vectorize(["b", "a", "a", "zzz"], vocab)
    assert idxs == [0, 1]
    n = math.sqrt(3.0 * 3.0 + 2.0 * 2.0)
    assert vals == [3.0 / n, 2.0 / n]
    assert features.vectorize(["zzz"], vocab) == ([], [])


def test_idf_formula():
    vocab = features.load_vocabulary(ROOT / "data/df_train.json")
    i = 0
    assert vocab.idf[i] == math.log((1 + vocab.n_train) / (1 + vocab.df[i])) + 1.0
