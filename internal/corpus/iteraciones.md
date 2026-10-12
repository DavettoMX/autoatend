## Iteraciones del generador

El criterio K1–K3 quedó registrado (`PLAN.md` §3, commit `test: pre-registered separability criterion K1-K3 … (red)`) antes de cambiar cualquier plantilla. Cada iteración se ajustó mirando **solo** la salida de la batería; ningún modelo se entrenó ni se evaluó hasta que la batería pasó.

| # | Cambio | K1 (Jaccard mín.–máx.) | K2 (global; Reconocimientos) | K3 (Macro-F1 sonda) | Veredicto |
|---|---|---|---|---|---|
| 0 | Plantillas originales (corpus `b810bad2…`) | 0.111–0.282 FAIL | 0.0000; 0.0000 FAIL | 1.0000 FAIL | FAIL |
| 1 | +6 plantillas de contexto cruzado por categoría (a partir de `docs/data/grok_examples.txt` y de los ejemplos del usuario), +10 plantillas compartidas en los pares del plan, marco "(ref. {FOLIO})" en las 5 categorías | 0.099–0.316 FAIL (Facturación–Reconocimientos) | 0.0925; 0.0500 PASS | 0.9710 PASS | FAIL |
| 2 | Vocabulario de facturación y cobros en 3 plantillas de contexto cruzado de Reconocimientos y "atención"/"todo" en 2 de Facturación | 0.111–0.316 FAIL | 0.0926; 0.0518 FAIL (Reconocimientos > 0.05) | 0.9702 PASS | FAIL |
| 3 | +6 plantillas base de Reconocimientos sobre pagos y facturación (diluyen las 2 compartidas de Reconocimientos y comparten vocabulario con Cobros y Facturación) | 0.156–0.307 PASS | 0.0895; 0.0381 PASS | 0.9726 PASS | PASS |
| 4 | Coherencia: la plantilla compartida Quejas↔Reconocimientos decía "tuve que llamar" tras un canal escrito; ahora "tuve que insistir" | 0.156–0.299 PASS | 0.0895; 0.0381 PASS | 0.9726 PASS | PASS (corpus persistido) |

Antes del registro, el análisis exploratorio de `/tmp/separabilidad/` (variantes F y FA) sí entrenó modelos de prueba sobre corpus alternativos; esos resultados motivaron el criterio, pero no se usaron para fijar sus umbrales ni para ajustar estas iteraciones.
