# AutoAtend — CENEVAL caso 10433

Sistema de clasificación de comentarios de clientes de **TelArg** en cinco categorías (Falta de Servicio, Cobros Indebidos, Problemas de Facturación, Quejas, Reconocimientos), desarrollado para el proceso de acreditación por Acuerdo 286, Licenciatura en Ciencias Computacionales.

Atiende las cinco consignas del caso:

| Consigna | Entregable | Estado |
|---|---|---|
| a) Base de datos del corpus | `data/corpus.db` (SQLite 3FN) generado por `cmd/gencorpus` | ✅ Fase 1 |
| b) Programa analizador | Analizador léxico propio (parser de ER → Thompson → AFN-ε) y `cmd/simulador` → `resultados.db` | 🚧 en curso (Fase 2: lexer y barrido listos; falta el motor de ER) |
| c) Plan de pruebas | TE §6.3 y Anexo II; pruebas en el repo y matriz `docs/CASOS_PRUEBA.md` | ⏳ Fases 2–5 |
| d) Justificación del plan | TE §6.3 (ISO/IEC 25010, invariantes I-1…I-4); cifras medidas en `PLAN.md` §5.2 | ⏳ Fase 5 |
| e) Servicio web clasificador | `cmd/autoatend`: HTTP en `127.0.0.1:8471`, inferencia nativa sobre `modelo.json` | ⏳ Fases 3–4 |

El plan de ejecución hacia la demo, con su checklist de aceptación, está en [`PLAN.md`](PLAN.md).

## Requisitos

- Go 1.26 o superior. No hace falta compilador de C: la única dependencia es `modernc.org/sqlite`, escrita en Go puro.
- `make` (opcional; los comandos equivalentes aparecen abajo).
- Python 3.11 y scikit-learn, solo para el entrenamiento offline de la Fase 3 (no forma parte del binario).

## Uso rápido

```bash
make test               # pruebas unitarias (omite la integración de 73,000 filas)
make test-integration   # suite completa, incluida la integración Level 4a
make cover              # informa la cobertura por paquete (no falla)
make cover-gate         # falla si algún paquete queda por debajo del 80 % de cobertura
make lint-regexp        # falla si aparece import "regexp" o cgo
make sweep              # barre data/corpus.db y genera data/resultados.db y data/df_train.json
make sweep-sha          # invariante I-3: dos barridos dan el mismo SHA-256
make cross              # binarios estáticos en dist/ para windows/amd64 y darwin/arm64
```

### Regenerar el catálogo léxico

`data/resultados.db` y `data/df_train.json` están versionados. `cmd/simulador` los reconstruye a partir de `data/corpus.db` con el lexer propio (`internal/lexer`, ADR-003): 6,720 lexemas con su ER, frecuencias global y por categoría, y variantes.

```bash
make clean-results && make sweep
sha256sum data/resultados.db
# d93d33a6b7f28395d20408f5d862d4034985d48e8eefa756405a3b75e266c09f
```

El simulador no sobrescribe archivos existentes. Si falla algún chequeo posterior (la consulta I-2 del TE §6.3.3, una sola corrida, lexemas sin variantes o sin frecuencias, integridad de `categoria_id`), borra lo que creó. `df_train.json` es la frecuencia de documento de train para el entrenamiento en Python.

### Regenerar el corpus

`data/corpus.db` está versionado y es reproducible: semilla **10433**, 73,000 comentarios (100 por día durante 730 días, del 2024-08-30 al 2026-08-29) y distribución 30/20/20/20/10.

```bash
make clean-data && make gen
sha256sum data/corpus.db
# ee3f7d55420b455f914359d248facde22c14227b55ce7c838f424fd57ac8d219
```

`gencorpus` no sobrescribe un archivo existente. Antes de escribir, ejecuta una batería de validación (conteos por categoría, textos vacíos, fechas, tasa de ruido, solapamiento léxico y una muestra de 50 comentarios por categoría); si algo falla, no persiste nada. El resultado queda en [`docs/reports/generator_validation_report.md`](docs/reports/generator_validation_report.md).

Para generar en otra ruta: `go run ./cmd/gencorpus -out /tmp/c.db -report /tmp/r.md`.

## Estructura

```
cmd/gencorpus/          CLI del generador (solo flags)
cmd/simulador/          CLI del barrido léxico (solo flags)
internal/corpus/        generación, ruido, validación, persistencia y reporte
internal/corpus/templates/  plantillas y vocabularios del dominio (solo datos)
internal/db/            acceso en tiempo de ejecución: categorías, OpenReadOnly, DSN
internal/db/build/      solo pipeline de generación: crear, finalizar, split 80/20, ATTACH
internal/db/schema/     DDL de corpus.db y resultados.db (.sql embebido, Anexo I)
internal/lexer/         normalizador, tokenizador AFD, stoplist, stemmer, tabla de símbolos, generador de ER
internal/regex/         motor de ER propio: parser, Thompson y simulación del AFN-ε
internal/sweep/         barrido del corpus: construye resultados.db y df_train.json
internal/testutil/      helpers de prueba compartidos
data/                   artefactos generados (corpus.db, resultados.db, df_train.json)
docs/adr/               decisiones de arquitectura (ADR)
docs/reports/           reportes generados
```

## Bases de datos

- **`corpus.db`**: `categorias` (CHECK con las 5 etiquetas del caso) y `comentarios` (FK `ON DELETE RESTRICT`, fecha `YYYY-MM-DD`), como en el Anexo I. El split 80/20 no se guarda: es la consulta `db.SplitCutoff`, un corte temporal global (TE §5.2.2). Todo comentario de train es anterior a todo comentario de test; con este corpus el corte es 2026-04-06 (58,400 / 14,600).
- **`resultados.db`**: `corridas`, `lexemas`, `frecuencias` (`categoria_id NULL` = frecuencia global) y `variantes`.

Como indica el Anexo I, `corpus.db` usa `journal_mode=WAL` y `resultados.db` `journal_mode=DELETE`. Los dos se compactan con `VACUUM` y quedan como un único archivo, sin `-wal` ni `-shm`, para poder embeberse en el binario. La lectura usa `mode=ro&immutable=1` (TE §5.1.4).

## Decisiones de arquitectura

Solo hay ADR para las decisiones que el Trabajo Escrito promete documentar (§5.2.1, §6.2.3 y §6.4).

| ADR | Tema |
|---|---|
| [ADR-001](docs/adr/ADR-001-sqlite-modernc.md) | SQLite mediante `modernc.org/sqlite` (Go puro, sin CGo) |
| [ADR-002](docs/adr/ADR-002-simulacion-afn.md) | Simulación directa del AFN-ε por conjuntos de estados, sin construir un AFD (O(n·m)) |
| [ADR-003](docs/adr/ADR-003-stoplist.md) | Stoplist, homógrafos, números y criterio de agrupación de variantes |

Pendientes: ADR-004 (modelo) y ADR-005 (despliegue: solo loopback, sin TLS, autenticación ni CORS).

## Objetivos de despliegue

Un solo ejecutable, sin instalación, sin permisos de administrador y sin red:

- **Principal:** Windows x86-64 (Core i5, 8 GB) → `dist/autoatend.exe`
- **Respaldo:** macOS ARM64 (Apple M2) → `dist/autoatend-macos-arm64`
