# CLAUDE.md — AutoAtend (CENEVAL caso 10433)

Proyecto de examen. Las fuentes de verdad, en orden de precedencia:
1. `../documentos/2907_10433.pdf`: enunciado del caso (consignas a–e).
2. `../documentos/TrabajoEscrito_Consolidado.pdf`: diseño ya entregado a los sinodales, por lo que se encuentra congelado y no se modifica. Toda diferencia entre lo implementado y lo escrito se justifica en un ADR (`docs/adr/`) para la defensa oral.
3. `PLAN.md`: plan por fases. `OPTIMIZATIONS.md` es la checklist viva: marca `[x]` solo lo que verificaste ejecutándolo.

## Reglas duras (no negociables)

- **Prohibido `import "regexp"`** y cualquier biblioteca de NLP o generador de parsers. El analizador léxico (parser de ER, Thompson, simulación del AFN-ε) se escribe a mano. Verifícalo con `make lint-regexp`.
- **Prohibido CGo.** Todo debe compilar con `CGO_ENABLED=0` para `windows/amd64` y `darwin/arm64` (`make cross`). Solo se permiten dependencias en Go puro; hoy la única es `modernc.org/sqlite`.
- **El único `runtime.GOOS` permitido** es la apertura opcional del navegador en la Fase 4.
- **`internal/db` (tiempo de ejecución) nunca importa `internal/db/build`.** `ATTACH`, la creación de archivos y el split viven solo en `build` (ADR-007).
- **Determinismo (invariante I-3):** nada de timestamps, nada de iterar mapas al escribir y PRAGMAs fijos. Dos ejecuciones con la misma semilla deben producir archivos byte a byte idénticos.
- **No dar por hecho nada que no se ejecutó.** Antes de afirmar que algo pasa, córrelo.

## Regresión del corpus

`data/corpus.db` está versionado. Cualquier refactor de `internal/corpus`, `internal/corpus/templates`, `internal/db/build` o `internal/db/schema` debe dejar intacto su SHA-256:

```bash
d=$(mktemp -d); go run ./cmd/gencorpus -out $d/c.db -report $d/r.md
sha256sum $d/c.db   # debe ser c0343583040f51afa54cfddb5bebe608bac3a02aa1416286582ae97744adfb05
diff $d/r.md docs/reports/generator_validation_report.md   # debe salir vacío
```

Cosas que cambian el SHA sin que se note:
- **El texto de la DDL en `internal/db/schema/*.sql`.** SQLite guarda literalmente cada sentencia `CREATE` en `sqlite_master`, así que un espacio, una sangría o un comentario de más lo alteran. Los `.sql` no llevan comentarios; la explicación va en `schema.go`.
- **El orden o la longitud de cualquier slice de `templates`** y el orden de las llamadas al rng en `internal/corpus`.

Si un cambio *debe* alterar el corpus, es una decisión de diseño: regenera `data/corpus.db` y el reporte, actualiza el SHA aquí y en `OPTIMIZATIONS.md`, y regístralo en un ADR.

## Metodología: TDD estricto

El Trabajo Escrito (§5.2.1) compromete el ciclo rojo-verde-refactor y los sinodales pueden pedir evidencia, que es el historial de git:
- Commits separados: `test: … (red)` → `feat: … (green)` → `refactor: …`.
- Nombra las pruebas con el ID del plan cuando exista, por ejemplo `TestCPU01_ParserPrecedencia` o `TestCPP02_ConservacionFrecuencias`.
- La cobertura debe ser **≥ 80 % por paquete** (`make cover-gate`; se excluyen `cmd/...` e `internal/testutil`).
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
3. El plan cita qué reutiliza (`archivo:línea`) y justifica cualquier código nuevo que se parezca a algo existente.
4. Ciclo rojo → verde → refactor, con un commit por paso.
5. `make verify` en verde.
6. Pedir aprobación para el merge.

## Criterios de éxito

### Por tarea (Definition of Done)
- [ ] `make verify` en verde: `fmt-check`, `vet`, `lint-regexp`, `test-integration`, `cover-gate` (≥ 80 % por paquete), `corpus-sha` y `cross`.
- [ ] Commits rojo → verde → refactor visibles en la rama.
- [ ] Pruebas nombradas con su ID del plan cuando exista.
- [ ] Sin código duplicado respecto a lo que reportó `mapeo-reuso`.
- [ ] Toda desviación del Trabajo Escrito tiene su ADR.
- [ ] `OPTIMIZATIONS.md` y el "Mapa del código" de este archivo actualizados si cambiaron.
- [ ] Ningún `[x]` marcado sin haberlo ejecutado.

**Estado actual:** `make verify` falla en `cover-gate` por `internal/db/build` (75.9 %). La primera tarea es `fix/cover-build`; hasta entonces ningún merge cumple este criterio.

### Por fase (umbrales de `PLAN.md`)

| Fase | Criterios de aceptación |
|---|---|
| 1 — Corpus | SHA de `data/corpus.db` = `CORPUS_SHA`; conteos exactos 21,900/14,600/14,600/14,600/7,300; distribución ±1 %; 0 textos vacíos; sin fuga temporal en el split |
| 2 — Analizador | I-1 round-trip = 0 fallos; I-2 = 0 filas (consulta agrupada por `(l.id, f.corrida_id)`); I-3 mismo SHA-256 de `resultados.db` en dos corridas; I-4 R² ≥ 0.99 y T(73k) < 60 s; CPU01–CPU14, CPP01–CPP04 y CPE01–CPE04 en verde; `docs/PLAN_PRUEBAS.md` con L-1 y L-2 |
| 3 — Modelo | CPM01 Macro-F1 ≥ 0.85; CPM02 recall de Reconocimientos ≥ 0.70; CPA06 paridad Python↔Go L1 = 0 en 1,000 comentarios; métricas también sobre test sin duplicados en train |
| 4 — Servicio | CPA01–CPA05 en verde; P99 < 50 ms en 100 llamadas a `/api/clasificar`; `--selftest` < 3 s con exit 0/1; fallback 8471 → 8472 |
| 5 — CI/CD | CI verde en ubuntu, windows y macos; `--selftest` en Windows x86-64 y macOS ARM64; ADR-001 a ADR-008 presentes; checklist final ↔ `docs/PLAN_PRUEBAS.md` consistentes en ambos sentidos |

## Comandos

```bash
make verify       # criterio de terminado: todos los gates en orden
make test | test-integration | cover-gate | lint-regexp | corpus-sha | cross
make gen          # se niega a sobrescribir data/corpus.db; usa make clean-data antes
go test -count=1 -run NombrePrueba ./internal/...
```

## Mapa del código

- `cmd/<herramienta>/main.go`: solo flags; la lógica vive en `internal/`.
- `internal/corpus`: generador del corpus sintético (semilla 10433, 73,000 filas, 30/20/20/20/10).
- `internal/db`: acceso en tiempo de ejecución (`Categories`, `OpenReadOnly`, `DSN`, constantes del split).
- `internal/db/build`: solo generación (`CreateCorpus`, `CreateResultados`, `Finalize`, `ApplySplit`, `CheckCategoriaIntegrity`).
- `internal/db/schema`: DDL embebida (Anexo I + `particiones` + `uq_frecuencias_global`).
- `internal/testutil`: helpers de prueba compartidos (`MustExec`, `ExpectErr`, `Scalar`).
- Por crear según `PLAN.md`: `internal/regex`, `internal/lexer` (stoplist versionada en `internal/lexer/stoplist/v1.json`), `cmd/simulador`, `ml/` (Python, offline), `internal/ml`, `internal/web`, `internal/assets` (`go:embed` no admite `..`, así que los artefactos embebidos van ahí), `cmd/autoatend`.

## Desviaciones del diseño ya decididas

No las "corrijas" de vuelta al Trabajo Escrito:
- `journal_mode=DELETE` en ambas BD, no WAL. En ejecución, BD en memoria con `Deserialize` y `PRAGMA query_only=ON` (ADR-007, pendiente de redactar).
- Índice parcial `uq_frecuencias_global`, porque SQLite trata los NULL como distintos dentro de la PK (ADR-006).
- Stoplist: `se`, `no`, `este` y `eso` **no** se filtran. Los homógrafos `bajo`, `para`, `sobre`, `entre`, `como`, `más`, `solo`, `aun` y `hacia` se conservan. Solo `según`/`segun` y `olé`/`ole` tienen variante sin tilde; no hay plegado global de acentos (ADR-008, pendiente).
- Consulta I-2 agrupada por `(l.id, f.corrida_id)`.
- El vocabulario y el IDF del modelo salen solo de la partición train; el catálogo de `resultados.db` se construye con el corpus completo.

## Pendientes conocidos que afectan el diseño

Consulta `OPTIMIZATIONS.md` §2 antes de empezar la Fase 2:
- Las abreviaturas de ruido `q` (por `que`) y `x` (por `por`) no están en la stoplist v1.
- El 19.5 % de los comentarios de test tiene un duplicado exacto en train, lo que infla CPM01.

## Idioma

Código, identificadores y comentarios de Go en inglés; los nombres de dominio y de esquema van en español (`comentarios`, `categoria_id`, `corridas`). Los documentos para los sinodales (`docs/PLAN_PRUEBAS.md`, reportes) van en español.
