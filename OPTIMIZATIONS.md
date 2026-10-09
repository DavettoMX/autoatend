# OPTIMIZATIONS — Checklist de desarrollo derivado de la auditoría

**Fecha de auditoría:** 2026-10-08
**Alcance:** código en `proyect/` contra `PLAN.md`, `2907_10433.pdf` (fuente de verdad) y `TrabajoEscrito_Consolidado.pdf` (diseño enviado a los sinodales).

---

## 0. Estado actual (línea base verificada)

Solo existe la **Fase 1** (consigna a). Las Fases 2–5 están en 0 %.

| Fase | Estado |
|---|---|
| 1 — BD y corpus (a) | ~90 %, con brechas menores |
| 2 — Analizador léxico (b, c, d) | 0 % — faltan `internal/regex`, `internal/lexer`, `cmd/simulador`, `docs/PLAN_PRUEBAS.md` |
| 3 — ML (Python → JSON → Go) | 0 % |
| 4 — Servicio web (e) | 0 % |
| 5 — CI/CD y verificación multiplataforma | 0 % — no hay `.github/`, `Makefile`, `.gitignore` |

Verificado en la auditoría (no requiere acción, sirve como regresión):

- [x] `go vet ./...` sin hallazgos.
- [x] Suite completa en verde, incluida la integración Level 4a con 73,000 filas (~1.6 s).
- [x] `corpus.db` regenerado es **byte a byte idéntico** (SHA-256 `c0343583040f51afa54cfddb5bebe608bac3a02aa1416286582ae97744adfb05`); el reporte de validación también es idéntico.
- [x] Compilación cruzada con `CGO_ENABLED=0` para `windows/amd64` y `darwin/arm64`.
- [x] Cero `import "regexp"` en el código.
- [x] DDL coincide con el Anexo I, salvo las dos desviaciones previstas por el plan (journal `DELETE` en lugar de `WAL`; índice parcial `uq_frecuencias_global`, ADR-006).

---

## 1. Metodología TDD

Las pruebas tienen estilo TDD (oráculos independientes, verificación directa sobre la BD, prueba de mutación documentada en ADR-006, determinismo y "todo o nada"), pero **no hay evidencia del ciclo rojo-verde-refactor** que promete el §5.2.1 del Trabajo Escrito.

- [ ] **[ALTA]** Inicializar el repositorio git (`git init`) y versionar el estado actual como línea base.
- [ ] **[ALTA]** A partir de la Fase 2, hacer commits con el ciclo explícito: `test: … (red)` → `feat: … (green)` → `refactor: …`. Es la evidencia de TDD que pueden pedir en la defensa oral.
- [x] **[ALTA]** Subir la cobertura de `internal/db` a **≥ 80 %** (umbral del plan). Hecho en `fix/cover-build`: `internal/db/build` 75.9 % → 91.8 %, `internal/db` 85.0 %.
- [ ] **[MEDIA]** Cubrir `ValidationError.Error()` (0 %) y las ramas restantes de `run`, `Persist` y `verifyCorpus`.
- [ ] **[MEDIA]** Rastrear las pruebas a los IDs del plan (CPU/CPP/CPE/CPA): nombres como `TestCPU01_...` o un comentario con el ID, más una matriz ID→prueba en `docs/PLAN_PRUEBAS.md`. Hace falta para la verificación bidireccional de la Fase 5 (A1).
- [x] **[MEDIA]** Criterios de éxito ejecutables: `make verify` (fmt-check, vet, lint-regexp, test-integration, cover-gate, corpus-sha, cross) y sección "Criterios de éxito" en `CLAUDE.md`.
- [x] **[MEDIA]** Agente `mapeo-reuso` (`.claude/agents/mapeo-reuso.md`, Sonnet, solo lectura) para mapear código reutilizable antes de cada tarea.
- [ ] **[BAJA]** Reemplazar `testing.Short()` por la etiqueta de compilación `//go:build integration` en las pruebas de nivel 4.

---

## 2. Brechas — prioridad alta

- [ ] **Fuga de datos train/test por duplicados exactos.** 2,848 de 14,600 comentarios de test (**19.5 %**) tienen su texto idéntico en train (345 de ellos en Reconocimientos), y ningún texto se repite entre categorías distintas. ADR-002 lo reconoce, pero solo como cifra a reportar. Esto infla CPM01/CPM02.
  - [ ] En la Fase 3, reportar además el Macro-F1 y el recall de Reconocimientos sobre el subconjunto de test **sin duplicado en train**.
  - [ ] Evaluar si conviene deduplicar antes de evaluar, y registrar la decisión en ADR-004.
- [ ] **La validación de "no trivialmente separable" está mal medida.** El Jaccard de las 100 palabras más frecuentes se calcula *sin quitar palabras funcionales* (de, la, mi, me…), así que el rango 0.22–0.34 refleja sobre todo esas palabras.
  - [ ] Recalcular el solapamiento después de aplicar la stoplist v1 (Fase 2).
  - [ ] Agregar una prueba de cordura: si un Naive Bayes trivial supera Macro-F1 0.99, el corpus es demasiado fácil y hay que recalibrar las plantillas o el ruido.
- [ ] **Abreviaturas del ruido que la stoplist no contempla.** El inyector produce `que→q` y `por→x`; `que` y `por` se filtran, pero `q` y `x` no están en la stoplist v1 del plan, así que entrarán al catálogo como lexemas falsos. ADR-008 solo contempla variantes sin tilde.
  - [ ] Decidir la política (agregar `q`→conjunción y `x`→preposición a `v1.json`, o documentarlo como limitación medida en L-1).
  - [ ] Extender ADR-008 y CPU14 con esa decisión antes de escribir `v1.json`.

---

## 3. Brechas — prioridad media

- [ ] **Procedencia del corpus.** §6.1.4 dice que la tasa de ruido "se documenta en el campo `parametros` de la tabla `corridas`", pero `corpus.db` no guarda la semilla, la configuración de ruido ni la versión del generador (`PRAGMA user_version = 0`).
  - [ ] Agregar una tabla `metadatos` (o usar `user_version` más una tabla clave/valor) con la semilla, el `NoiseConfig`, el total, la fecha de inicio y la versión del generador.
  - [ ] Hacer que `cmd/simulador` copie esa configuración a `corridas.parametros`.
  - [ ] Cubrir el entregable "migrations" del plan, al menos con un esquema versionado mediante `user_version`.
- [ ] **El generador no se configura por línea de comandos.** El diseño dice que el ruido es configurable; `main.go` solo acepta `-out` y `-report`.
  - [ ] Agregar flags `-seed`, `-total` y `-noise-rate` (los valores por defecto deben dar el mismo SHA del corpus de línea base).
  - [ ] Necesario para la prueba de linealidad I-4 (barridos de 1k/10k/30k/50k/73k).
- [ ] **Revisión manual de la muestra.** §6.1.4 dice "detección manual"; lo implementado es un chequeo de cordura automático más la impresión de la muestra.
  - [ ] Revisar a mano los 250 comentarios y dejar constancia (fecha, revisor, observaciones) en `docs/reports/`.
- [ ] **ADR-007 se cita pero no existe.** Lo mencionan `schema.go`, `attach.go` y ADR-001; el plan lo entrega hasta la Fase 4.
  - [ ] Redactar ADR-007 ahora (journal `DELETE`, `Deserialize`, integridad por base de datos), porque la Fase 1 ya depende de él.
- [ ] **`frecuencias.categoria_id` no se valida en tiempo de ejecución.**
  - [ ] Agregar `CHECK (categoria_id IS NULL OR categoria_id BETWEEN 1 AND 5)` como defensa local (desviación menor; registrarla en ADR-006 o ADR-007).
- [ ] **Inconsistencia del diseño en `lexemas.forma_canonica`.** Es `UNIQUE` en toda la tabla, pero el texto del diseño dice "unicidad dentro de una corrida". No se nota hoy porque cada corrida usa un archivo propio.
  - [ ] Documentar la decisión: dejarlo así (una corrida por archivo) o cambiarlo a `UNIQUE(forma_canonica, corrida_id)`.

---

## 4. Brechas — prioridad baja

- [ ] `cmd/gencorpus/report.go`: `defer f.Close()` ignora el error al cerrar, así que un reporte truncado pasaría sin aviso. Revisar el error de `Close`.
- [ ] `cmd/gencorpus/pipeline.go`: si `verifyCorpus` falla después de `Persist`, el `corpus.db` se queda en disco. Borrarlo para mantener el "todo o nada".
- [ ] `cmd/gencorpus/noise.go`: `typo()` sobre una letra mayúscula siempre termina duplicándola, porque `qwertyNeighbors` solo tiene minúsculas. Normalizar con `unicode.ToLower` antes de buscar.
- [x] `corpus.db` (16 MB) está en la raíz: moverlo a `data/` (ver §5).
- [x] Crear `.gitignore` (`dist/`, `*.out`, `*.test`; `data/corpus.db` sí se versiona).
- [x] Crear `README.md` (cómo generar, probar y compilar) y `CLAUDE.md` (reglas para agentes).

---

## 5. Reestructura de directorios

Estructura objetivo:

```
proyect/
├── cmd/
│   ├── gencorpus/main.go          # SOLO flags → corpus.Run(...)
│   ├── simulador/main.go          # Fase 2
│   └── autoatend/main.go          # Fase 4
├── internal/
│   ├── corpus/                    # ← sale de cmd/gencorpus (package main)
│   │   ├── generate.go  noise.go  validate.go  report.go  pipeline.go  persist.go
│   │   └── templates/             # templates.go + vocab.go (datos separados de lógica)
│   ├── db/                        # runtime: OpenReadOnly, OpenEmbedded (Deserialize), Categories
│   │   ├── schema/corpus.sql      # DDL como .sql + go:embed = artefacto literal del Anexo I
│   │   ├── schema/resultados.sql
│   │   └── build/                 # create, Finalize, ApplySplit, CheckCategoriaIntegrity (ATTACH)
│   ├── regex/  lexer/  ml/  web/  # según PLAN
│   ├── assets/                    # go:embed: corpus.db, resultados.db, modelo.json, gold5.json, ui/
│   └── testutil/                  # newCorpus, mustExec, expectErr compartidos
├── ml/                            # Python (Fase 3)
├── data/                          # artefactos generados (corpus.db aquí, no en la raíz)
├── testdata/                      # fixtures de paridad Python↔Go
├── docs/
│   ├── adr/
│   ├── reports/                   # generator_validation_report.md, revisión manual, evaluación
│   └── PLAN_PRUEBAS.md
├── Makefile                       # gen, sweep, train, test, cover-gate, cross, lint-regexp
├── .gitignore
├── README.md
├── PLAN.md
└── OPTIMIZATIONS.md
```

Checklist (cada paso es un refactor **sin cambio de comportamiento**: el SHA-256 de `corpus.db` debe seguir siendo `c0343583…`):

- [x] **Mover el generador de `cmd/gencorpus` a `internal/corpus`** y dejar `cmd/gencorpus/main.go` solo con los flags.
  - Por qué: `package main` no se puede importar, y el generador (~1,300 líneas) se necesita desde la prueba I-4, desde las pruebas del barrido de la Fase 2 y desde `--selftest`.
- [x] **Separar `templates.go` y `vocab.go`** en `internal/corpus/templates/` (datos separados de la lógica).
- [x] **Dividir `internal/db`:**
  - [x] `internal/db` (ejecución): `Categories`, `OpenReadOnly` y, en la Fase 4, `OpenEmbedded` con `Deserialize`.
  - [x] `internal/db/build` (solo generación): `create`, `CreateCorpus`, `CreateResultados`, `InsertCategorias`, `Finalize`, `ApplySplit`, `CheckCategoriaIntegrity`.
  - Por qué: así la regla de ADR-007 ("ATTACH solo en el pipeline de generación") se cumple por la estructura misma, porque `cmd/autoatend` nunca importa `build`. También evita que `open.go` mezcle creación, solo lectura y `Deserialize`.
  - [ ] Agregar una prueba o un chequeo de CI que confirme que `cmd/autoatend` no depende de `internal/db/build` (`go list -deps`).
- [x] **Mover la DDL a `internal/db/schema/*.sql` con `go:embed`**, para que los sinodales puedan comparar el archivo directamente con el Anexo I.
- [ ] **Crear `internal/assets/`** para los archivos embebidos (`go:embed` no acepta rutas con `..`, así que tienen que estar dentro del paquete que los embebe). Los copia ahí un paso del `Makefile`; nunca se editan a mano.
- [x] **Crear `internal/testutil/`** con los helpers compartidos (`newCorpus`, `newResultados`, `mustExec`, `expectErr`, `scalar`) para los paquetes de las fases 2 a 4.
- [x] **Mover `corpus.db` a `data/`** y `generator_validation_report.md` a `docs/reports/` (actualizar los valores por defecto de los flags de `gencorpus`).
- [x] **Crear `Makefile`** con los objetivos:
  - [x] `gen` — genera `data/corpus.db` y su reporte.
  - [ ] `sweep` — `cmd/simulador` → `data/resultados.db` (Fase 2).
  - [ ] `train` — `ml/` → `modelo.json` (Fase 3).
  - [ ] `assets` — copia los artefactos a `internal/assets/`.
  - [x] `test` / `test-integration` / `cover-gate` (falla si algún paquete queda por debajo del 80 %).
  - [x] `lint-regexp` — falla si aparece `import "regexp"` o CGo.
  - [ ] `cross` — `dist/autoatend.exe` y `dist/autoatend-macos-arm64` con los comandos del §1 del plan. *(Parcial: hoy compila `gencorpus` para ambos objetivos; las líneas de `autoatend` están comentadas hasta la Fase 4.)*
  - Nota: `sweep`, `train` y `assets` existen como marcadores que fallan con "pendiente: Fase N".

---

## 6. Orden de ejecución recomendado

1. [ ] `git init` y commit de la línea base (§1).
2. [x] Subir la cobertura de `internal/db` a ≥ 80 % (§1).
3. [x] Reestructura de §5 (generador → `internal/corpus`, `internal/db/build`, DDL en `.sql`, `testutil`, `data/`, `Makefile`). Hecha el 2026-10-09: SHA del corpus sin cambios y reporte idéntico. Quedan para la Fase 4 `internal/assets/` y el chequeo `go list -deps`.
4. [ ] Procedencia del corpus y flags del generador (§3).
5. [ ] Decidir la política de `q`/`x` y redactar ADR-007, más la extensión de ADR-008 (§2, §3), **antes** de empezar la Fase 2.
6. [ ] Arrancar la Fase 2 con commits rojo-verde-refactor y pruebas nombradas con su ID del plan.
7. [ ] En la Fase 3: métrica sin duplicados y prueba de cordura con un modelo trivial (§2).
