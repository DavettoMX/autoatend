"""Fills the Python <-> Go parity section of the report from the output of
`go test -v -run TestCPA06 ./internal/ml` (the Go side is the one that measures it)."""

from __future__ import annotations

import re
import sys
from pathlib import Path

from .report import PARITY_MARKER

LINE = re.compile(r"PARIDAD n=(\d+) L1=(\S+) support_diff=(\d+) class_agree=(\d+) max_dprob=(\S+)")


def section(go_output: str) -> str:
    m = LINE.search(go_output)
    if not m:
        raise SystemExit("paridad: the PARIDAD line was not found in the Go test output")
    n, l1, sdiff, agree, dprob = m.groups()
    ok = "--- PASS: TestCPA06" in go_output and "FAIL" not in go_output
    return (
        "Medida por la prueba Go `TestCPA06_ParidadPythonGo` (`go test -run TestCPA06 ./internal/ml`) sobre "
        "`testdata/paridad_1000.json`: 1,000 comentarios de test elegidos con la semilla 10433. Python vectoriza con "
        "`cmd/tokenizar`; Go vuelve a tokenizar y vectorizar el texto con `internal/ml`.\n\n"
        "| Medida | Resultado |\n|---|---|\n"
        f"| Comentarios comparados | {int(n):,} |\n"
        f"| Soportes distintos (índices) | {sdiff} |\n"
        f"| L1 entre vectores TF-IDF (Σ\\|Δ\\|, igualdad de bits) | {l1} |\n"
        f"| Clase predicha idéntica | {int(agree):,} de {int(n):,} |\n"
        f"| Máx. \\|Δprob\\| (Go vs Python) | {dprob} |\n"
        f"| Resultado de la prueba | {'PASS' if ok else 'FAIL'} |\n\n"
        "Las probabilidades pasan por `exp`: Python usa la libm del sistema y Go `math.Exp`, que pueden diferir en el último bit, "
        "por eso se comparan con tolerancia 1e-12 y se reporta el máximo. Los vectores, en cambio, se comparan bit a bit.\n"
    )


def main(argv: list[str] | None = None) -> int:
    argv = sys.argv[1:] if argv is None else argv
    if len(argv) != 2:
        raise SystemExit("usage: python -m autoatend_ml.paridad REPORT.md GO_TEST_OUTPUT.txt")
    rep, log = Path(argv[0]), Path(argv[1])
    text = rep.read_text(encoding="utf-8")
    if PARITY_MARKER not in text:
        raise SystemExit(f"paridad: marker {PARITY_MARKER} not found in {rep}")
    new = text.replace(PARITY_MARKER, section(log.read_text(encoding="utf-8")).rstrip("\n"))
    rep.write_text(new, encoding="utf-8", newline="\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
