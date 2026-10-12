# CENEVAL Case 10433 - Contexto Arquitectónico (AutoAtend)

Proyecto de examen. Las fuentes de verdad, en orden de precedencia:
1. `../documentos/2907_10433.pdf`: enunciado del caso (consignas a–e).
2. `../documentos/TrabajoEscrito_Consolidado.pdf`: diseño ya entregado a los sinodales, por lo que se encuentra congelado y no se modifica. La implementación lo respeta tanto como sea posible; cualquier diferencia requiere la aprobación previa del usuario (ver "Apego al Trabajo Escrito").
3. `PLAN.md`: plan hacia la demo del 13-oct-2026. Su §5 es la checklist viva: marca `[x]` solo lo que verificaste ejecutándolo.

## Reglas duras (no negociables)

- **Prohibido `import "regexp"`** y cualquier biblioteca de NLP o generador de parsers. El analizador léxico (parser de ER, Thompson, simulación del AFN-ε) se escribe a mano. Verifícalo con `make lint-regexp`.
- **Prohibido CGo.** Todo debe compilar con `CGO_ENABLED=0` para `windows/amd64` y `darwin/arm64` (`make cross`). Solo se permiten dependencias en Go puro; hoy la única es `modernc.org/sqlite`.
- **El único `runtime.GOOS` permitido** es la apertura opcional del navegador en la Fase 4.
- **`internal/db` (tiempo de ejecución) nunca importa `internal/db/build`.** La creación de archivos, la DDL y el `VACUUM` viven solo en `build`; `internal/db` solo lee (`OpenReadOnly`, `SplitCutoff`).
- **Determinismo (invariante I-3):** nada de timestamps, nada de iterar mapas al escribir y PRAGMAs fijos. Dos ejecuciones con la misma semilla deben producir archivos byte a byte idénticos.
- **No dar por hecho nada que no se ejecutó.** Antes de afirmar que algo pasa, córrelo.
- **Apego al Trabajo Escrito.** Se implementa lo que dice el TE. Si no se puede cumplir tal cual, o si parece haber una opción mejor, **no se implementa ni se decide por cuenta propia: se valida primero con el usuario**, con el procedimiento de "Apego al Trabajo Escrito".

## Regresión del corpus

`data/corpus.db` está versionado. Cualquier refactor de `internal/corpus`, `internal/corpus/templates`, `internal/db/build` o `internal/db/schema` debe dejar intacto su SHA-256:

```bash
d=$(mktemp -d); go run ./cmd/gencorpus -out $d/c.db -report $d/r.md -ambiguos $d/a.json
sha256sum $d/c.db   # debe ser 37f285ff1da70f4bd5205c51e255a612d1a72ad45b70fc9fda35dbac1b3e6e63
diff $d/r.md docs/reports/generator_validation_report.md   # debe salir vacío
cmp $d/a.json data/corpus_ambiguos.json                    # debe salir vacío
```

Cosas que cambian el SHA sin que se note:
- **El texto de la DDL en `internal/db/schema/*.sql`.** SQLite guarda literalmente cada sentencia `CREATE` en `sqlite_master`, así que un espacio, una sangría o un comentario de más lo alteran. Los `.sql` no llevan comentarios; la explicación va en `schema.go`.
- **El orden o la longitud de cualquier slice de `templates`** y el orden de las llamadas al rng en `internal/corpus`.

Si un cambio *debe* alterar el corpus, es una decisión de diseño: regenera `data/corpus.db` y el reporte, actualiza el SHA aquí, en el `Makefile`, en `README.md` y en `PLAN.md`. Si el cambio se aparta del Trabajo Escrito, primero se valida con el usuario (ver "Apego al Trabajo Escrito").

## Metodología: TDD estricto

El Trabajo Escrito (§5.2.1) compromete el ciclo rojo-verde-refactor y los sinodales pueden pedir evidencia, que es el historial de git:
- Commits separados: `test: … (red)` → `feat: … (green)` → `refactor: …`.
- Nombra las pruebas con el ID del plan cuando exista, por ejemplo `TestCPU01_ParserPrecedencia` o `TestCPP02_ConservacionFrecuencias`.
- La cobertura objetivo es **≥ 80 % por paquete** (se excluyen `cmd/...` e `internal/testutil`). Hasta la demo se mide con `make cover` y no bloquea (`PLAN.md` §1.3); `make cover-gate` sigue disponible como gate estricto.
- Las pruebas de integración pesadas se saltan con `testing.Short()`; `make test` usa `-short`.
- Las aserciones se hacen contra oráculos independientes o consultando la BD directamente, no con las estadísticas que reporta el propio código bajo prueba.

## Flujo de git

Los comandos `make` y `go` se ejecutan desde `proyect/`; la raíz del repositorio git es la carpeta padre.
- Cada tarea se hace en su propia rama desde `main`: `faseN/<tema>`, `fix/<tema>` o `docs/<tema>`. No se hacen commits directos en `main`.
- Integración con `git merge --no-ff`; **nunca squash**, para conservar la evidencia rojo-verde-refactor.
- Antes de pedir el merge debe pasar `make verify`.
- El merge a `main` y el push solo se hacen con aprobación explícita del usuario.

## Flujo por tarea

1. Crear la rama (ver arriba).
2. Antes de diseñar, lanzar el agente `mapeo-reuso` con la descripción de la tarea.
3. El plan cita la sección del TE que implementa y qué reutiliza (`archivo:línea`), y justifica cualquier código nuevo que se parezca a algo existente. Si el plan se aparta del TE, se detiene aquí y se valida con el usuario.
4. Ciclo rojo → verde → refactor, con un commit por paso.
5. `make verify` en verde.
6. Pedir aprobación para el merge.

## Criterios de éxito

### Por tarea (Definition of Done)
- [ ] `make verify` en verde: `fmt-check`, `vet`, `lint-regexp`, `test-integration`, `cover` (informa la cobertura; no bloquea hasta la demo), `corpus-sha`, `sweep-sha` y `cross`.
- [ ] Commits rojo → verde → refactor visibles en la rama.
- [ ] Pruebas nombradas con su ID del plan cuando exista.
- [ ] Sin código duplicado respecto a lo que reportó `mapeo-reuso`.
- [ ] Ninguna desviación del Trabajo Escrito sin aprobación explícita del usuario; las aprobadas están registradas en el ADR que el TE promete para ese tema.
- [ ] `PLAN.md` §5 y el "Mapa del código" de este archivo actualizados si cambiaron.
- [ ] Ningún `[x]` marcado sin haberlo ejecutado.

**Estado actual:** `make verify` en verde desde `fase1/apego-te` (cobertura: corpus 91.5 %, db 87.5 %, db/build 90.8 %, db/schema 100 %).

### Por fase (umbrales de `PLAN.md`)

| Fase | Criterios de aceptación |
|---|---|
| 1 — Corpus | SHA de `data/corpus.db` = `CORPUS_SHA`; conteos exactos 21,900/14,600/14,600/14,600/7,300; distribución ±1 %; 0 textos vacíos; sin fuga temporal en el split |
| 2 — Analizador | I-1 round-trip = 0 fallos; I-2 = 0 filas (consulta literal del TE §6.3.3; una sola corrida por archivo); I-3 mismo SHA-256 de `resultados.db` en dos corridas; I-4 R² ≥ 0.99 y T(73k) < 60 s; CPU01–CPU14, CPP01–CPP04 y CPE01–CPE04 en verde; L-1 y L-2 medidos (`PLAN.md` §5.2) |
| 3 — Modelo | CPM01 Macro-F1 ≥ 0.85; CPM02 recall de Reconocimientos ≥ 0.70; CPA06 paridad Python↔Go L1 = 0 en 1,000 comentarios; métricas también sobre test sin duplicados en train |
| 4 — Servicio | CPA01–CPA05 en verde; P99 < 50 ms en 100 llamadas a `/api/clasificar`; `--selftest` < 3 s con exit 0/1; fallback 8471 → 8472 |
| 5 — Verificación | Ejecución nativa en Windows x86-64 y macOS ARM64 (sin CI ni Wine, desviación de `PLAN.md` §1.3): `--selftest` con salida 0 y guion de 5 pasos; ADR-001 a ADR-005 presentes; `docs/CASOS_PRUEBA.md` (matriz Anexo II → prueba) |

## Comandos

```bash
make verify       # criterio de terminado: todos los gates en orden
make test | test-integration | cover | cover-gate | lint-regexp | corpus-sha | sweep-sha | cross
make gen          # se niega a sobrescribir data/corpus.db; usa make clean-data antes
make sweep        # genera data/resultados.db y data/df_train.json; usa make clean-results antes
make train        # entrena (uv + scikit-learn) y regenera modelo.json, modelo_alternativo.json, gold5.json, paridad_1000.json y el reporte; make train-sha comprueba el determinismo
go test -count=1 -run NombrePrueba ./internal/...
```

## Mapa del código

- `cmd/<herramienta>/main.go`: solo flags; la lógica vive en `internal/`.
- `internal/corpus`: generador del corpus sintético (semilla 10433, 73,000 filas, 30/20/20/20/10). La batería incluye el criterio de separabilidad registrado antes de cambiar las plantillas (K1 Jaccard después de la stoplist, K2 colisión de la oración base, K3 sonda Naive Bayes; `PLAN.md` §3). `templates.BaseCount`/`CrossCount` separan plantillas base, de contexto cruzado y compartidas; `gencorpus` escribe también `data/corpus_ambiguos.json` (ids del subconjunto de colisión, fuera de la BD).
- `internal/db`: acceso en tiempo de ejecución, solo lectura (`Categories`, `OpenReadOnly`, `FileURI`, `ReadOnlyParams`, `SplitCutoff`).
- `internal/db/build`: solo generación (`CreateCorpus`, `CreateResultados`, `InsertCategorias`, `Finalize`, `CheckCategoriaIntegrity`).
- `internal/db/schema`: DDL embebida, literal del Anexo I.
- `internal/testutil`: helpers de prueba compartidos (`MustExec`, `ExpectErr`, `Scalar`).
- `internal/lexer`: front end léxico (ADR-003). `Normalize` (NFC restringida al alfabeto), `Tokenize` (AFD de 4 estados), `Lookup`/`Version` (stoplist versionada en `stoplist/v1.json`), `Stem`, `Analyze` (única definición léxica compartida), `Table` (tabla de símbolos), `GenerateER` y `PrefixIndex`.
- `internal/regex`: motor de ER propio (ADR-002). `Parse` (descenso recursivo sobre runas; AST tipado `Alternation`, `Concatenation`, `Closure`, `Group`, `Class`, `Char`; `ParseError` con posición; `MaxDepth` = 200), `Thompson`/`Compile` (AFN-ε de a lo más 2·|r| estados) y `(*NFA).Match` (simulación por conjuntos de estados, O(n·m), cadena completa). No importa `regexp`.
- `internal/sweep`: barrido del corpus (`Run`): construye `resultados.db` (corridas, lexemas, frecuencias, variantes) y `df_train.json`, y verifica I-2 reabriendo el archivo. Valida cada ER por ida y vuelta con `internal/regex` antes de persistirla (`validatedER`; respaldo a la forma canónica; `Result.ERFallbacks`, esperado 0). Puede importar `internal/db/build`.
- `cmd/simulador`: flags de `sweep.Run` (`make sweep`; `make sweep-sha` comprueba I-3).
- `internal/ml`: lado de ejecución del clasificador (ADR-004). `Lexicon`/`NewLexicon`/`LoadCatalogue` (forma → lexema canónico del catálogo por `lexer.Stem`), `RunTokenizar` (lógica de `cmd/tokenizar`: JSON lines `{id,texto}` → `{id,lexemas}`), `Load` (lee y valida `modelo.json`), `(*Model).Vectorize` (TF conteo × IDF almacenado, L2, suma en orden ascendente sin FMA), `Predict`/`PredictSparse` (`z = b + Σ w·x` sobre el soporte, softmax con resta del máximo, `LexemasClave` top-5 por `|w·x|`; `Vacio` si ningún token está en el vocabulario). Go nunca recalcula el IDF. Importa `internal/db`, nunca `internal/db/build`.
- `cmd/tokenizar`: flags de `ml.RunTokenizar`; Python lo invoca por lotes (CPA06). `make train` lo compila a `dist/tokenizar`.
- `ml/` (Python 3.11 + scikit-learn 1.9.1 con `uv`, solo desarrollo, nunca en el binario): `autoatend_ml.train` (split con la misma consulta SQL de `db.SplitCutoff`, TF-IDF de train, NBM y RLS × 3 estrategias de desbalance, CV estratificada de 5, evaluación en test con y sin duplicados, exportación), `export` (JSON determinista y escalador de referencia), `report`, `paridad`. `make train` regenera `data/modelo.json`, `data/modelo_alternativo.json`, `data/gold5.json`, `testdata/paridad_1000.json` y `docs/reports/evaluacion_modelo.md`; `make train-sha` comprueba que dos entrenamientos dan los mismos bytes. `make verify` no requiere Python: las pruebas Go leen los JSON versionados (CPA06, equivalencia lineal, `gold5`, CPM01/CPM02 medidos).
- Por crear según `PLAN.md`: `internal/web`, `internal/assets` (`go:embed` no admite `..`, así que los artefactos embebidos van ahí), `cmd/autoatend`.

## Apego al Trabajo Escrito

Precedencia: el caso (`2907_10433.pdf`) manda; después, el Trabajo Escrito tanto como sea posible. Los sinodales ya tienen el TE, así que cada diferencia con lo implementado es una pregunta potencial en la defensa oral.

Antes de apartarte del TE, o si crees que hay una opción mejor:
1. Busca una alternativa fiel al TE que cumpla el caso y, si es posible, verifícala ejecutándola. Si existe, úsala; no hace falta preguntar.
2. Si no existe, o si la opción distinta parece mejor, **detente y pregúntale al usuario antes de implementar**. Incluye:
   - la sección del TE afectada;
   - el motivo: qué consigna del caso la exige, o por qué la otra opción es mejor;
   - las alternativas fieles al TE que evaluaste y por qué no bastan;
   - la evidencia (pruebas o mediciones ejecutadas);
   - tu recomendación.
3. Solo con la aprobación explícita del usuario se implementa, y se registra en el ADR que el TE promete para ese tema. Si no hay un ADR prometido para ese tema, también se pregunta dónde registrarlo.

Esto también aplica a las "mejoras": robustez extra, índices, tablas, validaciones o pasos que el TE no menciona. Contradicciones internas del TE: no elijas por tu cuenta; plantea las opciones al usuario.

Solo hay ADR para lo que el TE promete documentar:
- ADR-001: SQLite.
- ADR-002: simulación directa del AFN-ε.
- ADR-003: stoplist y homógrafos.
- ADR-004: modelo.
- ADR-005: despliegue.

Ya implementado como lo dice el TE (no lo "mejores"):
- `corpus.db` en WAL, `resultados.db` en DELETE (Anexo I). Lectura con `mode=ro&immutable=1` (§5.1.4). Integridad cruzada con `ATTACH` (§6.1.3).
- DDL literal del Anexo I: sin índices ni tablas extra. La fila global duplicada que admite la PK la detecta la consulta I-2 del TE.
- Split temporal global (§5.2.2) como consulta (`db.SplitCutoff`), sin tabla.

Decisiones para las fases siguientes, ya aprobadas por el usuario (no se vuelven a preguntar):
- **I-2**: consulta literal del TE §6.3.3; cada `resultados.db` contiene una sola corrida (lo implican `forma_canonica UNIQUE` y el I-3).
- **IDF**: sale solo de train ("precalculados durante el entrenamiento", §6.4.2); el catálogo de `resultados.db` usa el corpus completo.
- **Fase 4**: el embed se escribe en el directorio temporal del usuario y se abre con `mode=ro&immutable=1`; `ATTACH` en tiempo de ejecución; sin `Deserialize`.
- **Stoplist (ADR-003)**: resuelve las contradicciones internas del TE a favor del caso.
  - Se filtran artículos, preposiciones (incluida `hacia`), conjunciones e interjecciones.
  - Se conservan los homógrafos del TE: `bajo`, `para`, `sobre`, `entre`, `como`, `más`, `solo`, `aun`.
  - `se`, `no`, `este` y `eso` no se filtran, porque no pertenecen a las 4 clases.
  - Variantes de superficie de formas filtradas, hueco detectado en el desarrollo: `segun`, `ole` (entrada libre en la consigna e) y las abreviaturas de ruido `q` y `x`.
  - No hay plegado global de acentos.

## Pendientes conocidos que afectan el diseño

Consulta `PLAN.md` §3 antes de empezar la Fase 2:
- El 24.4 % de los comentarios de test (3,560 de 14,600) tiene un duplicado exacto en train, lo que infla CPM01.

## Idioma

Código, identificadores y comentarios de Go en inglés; los nombres de dominio y de esquema van en español (`comentarios`, `categoria_id`, `corridas`). Los documentos para los sinodales (`docs/CASOS_PRUEBA.md`, reportes) van en español.
