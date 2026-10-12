"""Reads corpus.db (read-only) and applies the train/test split of the TE (§5.2.2).

The split is NOT recomputed here with any new rule: SPLIT_CUTOFF_SQL is a copy of
the query text of ``internal/db/split.go`` (``splitCutoffSQL``), with the same
parameter (TrainPercent = 80). Train is ``fecha < corte``; test is
``fecha >= corte``, exactly like ``internal/sweep`` does when it builds
``df_train.json``.
"""

from __future__ import annotations

import sqlite3
from dataclasses import dataclass
from pathlib import Path

# Copy of splitCutoffSQL in internal/db/split.go. Keep it identical.
SPLIT_CUTOFF_SQL = """
	SELECT fecha FROM comentarios
	ORDER BY fecha
	LIMIT 1 OFFSET (SELECT COUNT(*) * ? / 100 FROM comentarios)"""

# TrainPercent in internal/db/split.go.
TRAIN_PERCENT = 80

CATEGORIES = [
    (1, "Falta de Servicio"),
    (2, "Cobros Indebidos"),
    (3, "Problemas de Facturación"),
    (4, "Quejas"),
    (5, "Reconocimientos"),
]
RECONOCIMIENTOS = 5


@dataclass
class Corpus:
    cutoff: str
    ids: list[int]
    texts: list[str]
    labels: list[int]  # categoria_id, 1..5
    is_train: list[bool]


def open_readonly(path: str | Path) -> sqlite3.Connection:
    """Opens a SQLite file with mode=ro&immutable=1 (TE §5.1.4)."""
    uri = Path(path).resolve().as_uri() + "?mode=ro&immutable=1"
    return sqlite3.connect(uri, uri=True)


def load_corpus(path: str | Path) -> Corpus:
    con = open_readonly(path)
    try:
        (cutoff,) = con.execute(SPLIT_CUTOFF_SQL, (TRAIN_PERCENT,)).fetchone()
        rows = con.execute(
            "SELECT id, texto, categoria_id, fecha < ? FROM comentarios ORDER BY id",
            (cutoff,),
        ).fetchall()
    finally:
        con.close()
    return Corpus(
        cutoff=cutoff,
        ids=[r[0] for r in rows],
        texts=[r[1] for r in rows],
        labels=[r[2] for r in rows],
        is_train=[bool(r[3]) for r in rows],
    )


def catalogue_size(path: str | Path) -> int:
    con = open_readonly(path)
    try:
        return con.execute("SELECT COUNT(*) FROM lexemas").fetchone()[0]
    finally:
        con.close()
