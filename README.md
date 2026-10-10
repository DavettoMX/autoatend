# AutoAtend — CENEVAL caso 10433

Sistema de clasificación de comentarios de clientes de **TelArg** en cinco categorías (Falta de Servicio, Cobros Indebidos, Problemas de Facturación, Quejas, Reconocimientos), desarrollado para el proceso de acreditación por Acuerdo 286, Licenciatura en Ciencias Computacionales.

Atiende las cinco consignas del caso:

| Consigna | Entregable | Estado |
|---|---|---|
| a) Base de datos del corpus | `data/corpus.db` (SQLite 3FN) generado por `cmd/gencorpus` | ✅ Fase 1 |
| b) Programa analizador | Analizador léxico propio (parser de ER → Thompson → AFN-ε) y `cmd/simulador` → `resultados.db` | ⏳ Fase 2 |
| c) Plan de pruebas | `docs/PLAN_PRUEBAS.md` | ⏳ Fase 2 |
| d) Justificación del plan | `docs/PLAN_PRUEBAS.md` (ISO/IEC 25010, invariantes I-1…I-4) | ⏳ Fase 2 |
| e) Servicio web clasificador | `cmd/autoatend`: HTTP en `127.0.0.1:8471`, inferencia nativa sobre `modelo.json` | ⏳ Fases 3–4 |

El plan de ejecución completo está en [`PLAN.md`](PLAN.md) y la checklist de pendientes derivada de la auditoría está en [`OPTIMIZATIONS.md`](OPTIMIZATIONS.md).

## Requisitos

- Go 1.26 o superior. No hace falta compilador de C: la única dependencia es `modernc.org/sqlite`, escrita en Go puro.
- `make` (opcional; los comandos equivalentes aparecen abajo).
- Python 3.11 y scikit-learn, solo para el entrenamiento offline de la Fase 3 (no forma parte del binario).

## Uso rápido

```bash
make test               # pruebas unitarias (omite la integración de 73,000 filas)
make test-integration   # suite completa, incluida la integración Level 4a
make cover-gate         # falla si algún paquete queda por debajo del 80 % de cobertura
make lint-regexp        # falla si aparece import "regexp" o cgo
make cross              # binarios estáticos en dist/ para windows/amd64 y darwin/arm64
```

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
internal/corpus/        generación, ruido, validación, persistencia y reporte
internal/corpus/templates/  plantillas y vocabularios del dominio (solo datos)
internal/db/            acceso en tiempo de ejecución: categorías, OpenReadOnly, DSN
internal/db/build/      solo pipeline de generación: crear, finalizar, split 80/20, ATTACH
internal/db/schema/     DDL de corpus.db y resultados.db (.sql embebido, Anexo I)
internal/testutil/      helpers de prueba compartidos
data/                   artefactos generados (corpus.db)
docs/adr/               decisiones de arquitectura (ADR)
docs/reports/           reportes generados
```

## Bases de datos

- **`corpus.db`**: `categorias` (CHECK con las 5 etiquetas del caso) y `comentarios` (FK `ON DELETE RESTRICT`, fecha `YYYY-MM-DD`), como en el Anexo I. El split 80/20 no se guarda: es la consulta `db.SplitCutoff`, un corte temporal global (TE §5.2.2). Todo comentario de train es anterior a todo comentario de test; con este corpus el corte es 2026-04-06 (58,400 / 14,600).
- **`resultados.db`** (Fase 2): `corridas`, `lexemas`, `frecuencias` (`categoria_id NULL` = frecuencia global) y `variantes`.

Como indica el Anexo I, `corpus.db` usa `journal_mode=WAL` y `resultados.db` `journal_mode=DELETE`. Los dos se compactan con `VACUUM` y quedan como un único archivo, sin `-wal` ni `-shm`, para poder embeberse en el binario. La lectura usa `mode=ro&immutable=1` (TE §5.1.4).

## Decisiones de arquitectura

Solo hay ADR para las decisiones que el Trabajo Escrito promete documentar (§5.2.1, §6.2.3 y §6.4).

| ADR | Tema |
|---|---|
| [ADR-001](docs/adr/ADR-001-sqlite-modernc.md) | SQLite mediante `modernc.org/sqlite` (Go puro, sin CGo) |

Pendientes: ADR-002 (simulación directa del AFN-ε), ADR-003 (stoplist y homógrafos), ADR-004 (modelo) y ADR-005 (despliegue: solo loopback, sin TLS, autenticación ni CORS).

## Objetivos de despliegue

Un solo ejecutable, sin instalación, sin permisos de administrador y sin red:

- **Principal:** Windows x86-64 (Core i5, 8 GB) → `dist/autoatend.exe`
- **Respaldo:** macOS ARM64 (Apple M2) → `dist/autoatend-macos-arm64`
