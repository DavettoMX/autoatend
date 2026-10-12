# AutoAtend — CENEVAL caso 10433

Sistema de clasificación de comentarios de clientes de **TelArg** en cinco categorías (Falta de Servicio, Cobros Indebidos, Problemas de Facturación, Quejas, Reconocimientos), desarrollado para el proceso de acreditación por Acuerdo 286, Licenciatura en Ciencias Computacionales.

Atiende las cinco consignas del caso:

| Consigna | Entregable | Estado |
|---|---|---|
| a) Base de datos del corpus | `data/corpus.db` (SQLite 3FN) generado por `cmd/gencorpus` | ✅ Fase 1 |
| b) Programa analizador | Analizador léxico propio (parser de ER → Thompson → AFN-ε) y `cmd/simulador` → `resultados.db` | 🚧 en curso (Fase 2: lexer y barrido listos; falta el motor de ER) |
| c) Plan de pruebas | TE §6.3 y Anexo II; pruebas en el repo y matriz `docs/CASOS_PRUEBA.md` | ⏳ Fases 2–5 |
| d) Justificación del plan | TE §6.3 (ISO/IEC 25010, invariantes I-1…I-4); cifras medidas en `PLAN.md` §5.2 | ⏳ Fase 5 |
| e) Servicio web clasificador | `cmd/autoatend`: HTTP en `127.0.0.1:8471`, inferencia nativa sobre `modelo.json` | ✅ Fase 4 implementada (servicio, UI y `--selftest` listos en Linux; falta la verificación nativa en Windows y Mac, Bloque E) |

El plan de ejecución hacia la demo, con su checklist de aceptación, está en [`PLAN.md`](PLAN.md).

## Requisitos

- Go 1.26 o superior. No hace falta compilador de C: la única dependencia es `modernc.org/sqlite`, escrita en Go puro.
- `make` (opcional; los comandos equivalentes aparecen abajo).
- Solo para reentrenar (`make train`): [`uv`](https://docs.astral.sh/uv/) con Python 3.11 y scikit-learn 1.9.1 (fijados en `ml/pyproject.toml` y `ml/uv.lock`). No forman parte del binario ni de `make verify`.

## Uso rápido

```bash
make test               # pruebas unitarias (omite la integración de 73,000 filas)
make test-integration   # suite completa, incluida la integración Level 4a
make cover              # informa la cobertura por paquete (no falla)
make cover-gate         # falla si algún paquete queda por debajo del 80 % de cobertura
make lint-regexp        # falla si aparece import "regexp" o cgo
make sweep              # barre data/corpus.db y genera data/resultados.db y data/df_train.json
make sweep-sha          # invariante I-3: dos barridos dan el mismo SHA-256
make train              # entrena y regenera modelo.json, modelo_alternativo.json, gold5.json, paridad_1000.json y el reporte (requiere uv)
make train-sha          # dos entrenamientos independientes producen los mismos bytes
make assets             # copia los artefactos de data/ a internal/assets/data/ (los embebe el binario)
make build              # dist/autoatend para este equipo, con los artefactos embebidos
make cross              # binarios estáticos en dist/ para windows/amd64 y darwin/arm64
```

Los objetivos que compilan Go dependen de `make assets`. Con `go test` o `go build` directos, ejecute `make assets` antes (las copias en `internal/assets/data/` no se versionan).

### Regenerar el catálogo léxico

`data/resultados.db` y `data/df_train.json` están versionados. `cmd/simulador` los reconstruye a partir de `data/corpus.db` con el lexer propio (`internal/lexer`, ADR-003): 6,196 lexemas con su ER, frecuencias global y por categoría, y variantes.

```bash
make clean-results && make sweep
sha256sum data/resultados.db
# 9620196ffdd4e10e231ca4b7dc66b5a68bfa6e67057fbf4447ecaed8b24a79f3
```

El simulador no sobrescribe archivos existentes. Si falla algún chequeo posterior (la consulta I-2 del TE §6.3.3, una sola corrida, lexemas sin variantes o sin frecuencias, integridad de `categoria_id`), borra lo que creó. `df_train.json` es la frecuencia de documento de train para el entrenamiento en Python.

### Entrenar el modelo (consigna e, Fase 3)

`data/modelo.json` (ganador), `data/modelo_alternativo.json` (mejor de la otra familia), `data/gold5.json` y `testdata/paridad_1000.json` están versionados; `make verify` solo los lee. `make train` los reproduce: compila `cmd/tokenizar` (la tokenización la hace Go, no Python), entrena en `ml/` con TF-IDF de train, compara Naive Bayes multinomial y regresión logística softmax bajo tres estrategias de desbalance con validación cruzada estratificada de 5 particiones, evalúa una sola vez en test (también sin los duplicados exactos de train y sin el subconjunto ambiguo de `data/corpus_ambiguos.json`) y mide en Go la paridad Python↔Go (CPA06). El resultado y la evidencia están en [`docs/reports/evaluacion_modelo.md`](docs/reports/evaluacion_modelo.md) y en [`ADR-004`](docs/adr/ADR-004-modelo.md): ganador regresión logística (C = 1, línea base), Macro-F1 de test 0.9795 y recall de Reconocimientos 0.9819 (Naive Bayes: 0.9708). Las métricas se reportan también sin duplicados de train, sin el subconjunto ambiguo (0.9997) y, solo como reporte, con plantillas no vistas (0.8350 de media, [`plantillas_no_vistas.md`](docs/reports/plantillas_no_vistas.md)).

### Servicio web (consigna e, Fase 4)

`cmd/autoatend` es un único ejecutable con el corpus, el catálogo, los dos modelos, `gold5` y la interfaz embebidos (`internal/assets`). Al iniciar escribe las bases en un directorio temporal propio, las abre con `mode=ro&immutable=1`, verifica la integridad de `categoria_id` por `ATTACH` y escucha solo en `127.0.0.1:8471` (fallback 8472–8480). Detalles y la desviación aprobada "sin CI ni Wine" en [`ADR-005`](docs/adr/ADR-005-despliegue.md).

```bash
make assets build                      # o make cross para Windows y macOS
./dist/autoatend                       # banner con la URL y abre el navegador
./dist/autoatend -no-browser           # sin abrir el navegador
./dist/autoatend -modelo alternativo   # sirve Naive Bayes (modelo_alternativo.json)
./dist/autoatend --selftest            # verificaciones headless: salida 0 (éxito) o 1 (fallo), < 3 s
```

| Bandera | Efecto |
|---|---|
| `-modelo principal` o `alternativo` | Modelo servido; por omisión `principal` (regresión logística) |
| `-no-browser` | No abre el navegador |
| `--selftest` | Corre las verificaciones (corpus de 73,000 filas y 5 categorías, integridad, ER del catálogo compiladas con `internal/regex` que aceptan sus variantes, `gold5`, puerto) y sale con 0 o 1 |

La versión se fija al compilar: `go build -ldflags "-X main.version=1.0.0" ./cmd/autoatend` (`make` usa `VERSION`).

| Endpoint | Respuesta |
|---|---|
| `GET /` | Interfaz web (un HTML, fuentes del sistema, sin CDN) |
| `POST /api/clasificar` con `{"texto": "..."}` | `{"categoria", "confianzas": {5 categorías, suman 1}, "lexemas_clave": [{"lexema", "expresion_regular", "contribucion"}], "vacio"}`; texto vacío o JSON inválido → 400; cuerpo mayor a 1 MB → 413 |
| `GET /api/lexemas?pagina=1&tamanio=100` | Catálogo paginado (máx. 500 por página): `{total, pagina, tamanio, lexemas: [{forma_canonica, expresion_regular, frecuencia_global, frecuencias, variantes}]}`; filtros opcionales `q=prefijo` y `orden=frecuencia` o `orden=variantes` |
| `GET /api/categorias` | Las 5 categorías con su `id` |
| `GET /api/health` | `{"ok": true, "version", "lexemas", "puerto"}` |

### Regenerar el corpus

`data/corpus.db` está versionado y es reproducible: semilla **10433**, 73,000 comentarios (100 por día durante 730 días, del 2024-08-30 al 2026-08-29) y distribución 30/20/20/20/10.

```bash
make clean-data && make gen
sha256sum data/corpus.db
# 37f285ff1da70f4bd5205c51e255a612d1a72ad45b70fc9fda35dbac1b3e6e63
```

`gencorpus` no sobrescribe un archivo existente. Antes de escribir, ejecuta una batería de validación (conteos por categoría, textos vacíos, fechas, tasa de ruido, solapamiento léxico después de la stoplist, colisión de la oración base entre categorías, sonda Naive Bayes y una muestra de 50 comentarios por categoría); si algo falla, no persiste nada. Los tres criterios de separabilidad (K1–K3) se registraron en `PLAN.md` §3 antes de cambiar las plantillas. `make gen` escribe también `data/corpus_ambiguos.json`, con los ids del subconjunto ambiguo. El resultado queda en [`docs/reports/generator_validation_report.md`](docs/reports/generator_validation_report.md).

Para generar en otra ruta: `go run ./cmd/gencorpus -out /tmp/c.db -report /tmp/r.md`.

## Estructura

```
cmd/gencorpus/          CLI del generador (solo flags)
cmd/simulador/          CLI del barrido léxico (solo flags)
cmd/tokenizar/          CLI del tokenizador por lotes para Python (solo flags)
cmd/autoatend/          servicio web (consigna e): flags, señales y versión
internal/corpus/        generación, ruido, validación, persistencia y reporte
internal/corpus/templates/  plantillas y vocabularios del dominio (solo datos)
internal/db/            acceso en tiempo de ejecución: categorías, OpenReadOnly, DSN
internal/db/build/      solo pipeline de generación: crear, finalizar, split 80/20, ATTACH
internal/db/schema/     DDL de corpus.db y resultados.db (.sql embebido, Anexo I)
internal/lexer/         normalizador, tokenizador AFD, stoplist, stemmer, tabla de símbolos, generador de ER
internal/regex/         motor de ER propio: parser, Thompson y simulación del AFN-ε
internal/sweep/         barrido del corpus: construye resultados.db y df_train.json
internal/ml/            carga de modelo.json, vectorización TF-IDF e inferencia lineal + softmax
internal/web/           servicio: arranque a temporal, endpoints, puerto, navegador y --selftest
internal/assets/        go:embed de las bases, modelos, gold5 (copias de data/ por make assets) y ui/index.html
ml/                     entrenamiento offline en Python (uv, scikit-learn); nunca forma parte del binario
internal/testutil/      helpers de prueba compartidos
data/                   artefactos generados (corpus.db, resultados.db, df_train.json, modelo.json, modelo_alternativo.json, gold5.json)
testdata/               paridad_1000.json (fixture de CPA06)
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
| [ADR-004](docs/adr/ADR-004-modelo.md) | Elección del modelo y formato JSON |
| [ADR-005](docs/adr/ADR-005-despliegue.md) | Despliegue: solo loopback, sin TLS, autenticación ni CORS; un binario; sin CI ni Wine |

## Objetivos de despliegue

Un solo ejecutable, sin instalación, sin permisos de administrador y sin red:

- **Principal:** Windows x86-64 (Core i5, 8 GB) → `dist/autoatend.exe`
- **Respaldo:** macOS ARM64 (Apple M2) → `dist/autoatend-macos-arm64`
