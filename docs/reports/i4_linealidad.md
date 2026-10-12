# I-4: linealidad del barrido (CPP04)

Invariante I-4 del TE §6.3.3 y §6.3.8: el tiempo de barrido crece linealmente con el número de comentarios, T(k·N) ≈ k·T(N). Umbrales: pendiente positiva, R² ≥ 0.99 y T(73k) < 60 s.

- Equipo: Linux x86_64, 24 CPU, AMD Ryzen 9 9900X 12-Core Processor
- Fecha: 2026-10-12
- Método: subcorpus de los primeros N comentarios de `data/corpus.db` (por `id`), copiados a una base temporal; se mide `sweep.Run` completo (lectura, tabla de símbolos, persistencia y verificación de `resultados.db`); mediana de 3 corridas por tamaño.

| N | Mediana (s) | Mínimo (s) | Máximo (s) |
|---:|---:|---:|---:|
| 1000 | 0.018 | 0.018 | 0.019 |
| 10000 | 0.070 | 0.069 | 0.071 |
| 30000 | 0.179 | 0.178 | 0.179 |
| 50000 | 0.282 | 0.259 | 0.285 |
| 73000 | 0.368 | 0.361 | 0.371 |

Recta de mínimos cuadrados sobre las medianas: T(N) = 0.0219 + 4.925e-06·N segundos.

| Criterio | Umbral | Medido | Resultado |
|---|---|---|---|
| Pendiente | > 0 | 4.925e-06 s por comentario | cumple |
| R² | ≥ 0.99 | 0.99376 | cumple |
| T(73000) medido | < 60 s | 0.368 s (ajustado por la recta: 0.381 s) | cumple |
