# OPTIMIZATIONS — Checklist de desarrollo derivado de la auditoría

**Fecha de auditoría:** 2026-10-08; segunda auditoría y alineación con el TE: 2026-10-10 (rama `fase1/apego-te`)
**Alcance:** código en `proyect/` contra `PLAN.md`, `2907_10433.pdf` (fuente de verdad) y `TrabajoEscrito_Consolidado.pdf` (diseño enviado a los sinodales).

---

## 0. Estado actual (línea base verificada)

Solo existe la **Fase 1** (consigna a). Las Fases 2–5 están en 0 %.

| Fase | Estado |
|---|---|
| 1 — BD y corpus (a) | Completa y alineada con el TE; falta verificarla ejecutando en Windows |
| 2 — Analizador léxico (b, c, d) | 0 % — faltan `internal/regex`, `internal/lexer`, `cmd/simulador`, `docs/PLAN_PRUEBAS.md` |
| 3 — ML (Python → JSON → Go) | 0 % |
| 4 — Servicio web (e) | 0 % |
| 5 — CI/CD y verificación multiplataforma | 0 % — no hay `.github/` |

Verificado ejecutándolo (sirve como regresión):

- [x] `make verify` en verde: `fmt-check`, `vet`, `lint-regexp`, `test-integration` (incluida la Level 4a con 73,000 filas), `cover-gate`, `corpus-sha` y `cross`.
- [x] `corpus.db` regenerado es **byte a byte idéntico** entre corridas (SHA-256 `ee3f7d55420b455f914359d248facde22c14227b55ce7c838f424fd57ac8d219`), y el reporte de validación también. El contenido de `comentarios` es el mismo que el de la línea base anterior (hash de filas `ae30de6c…`); solo cambió el contenedor.
- [x] Compilación cruzada con `CGO_ENABLED=0` para `windows/amd64` y `darwin/arm64`.
- [x] Cero `import "regexp"` en el código.
- [x] DDL literal del Anexo I (`IF NOT EXISTS`, sin tablas ni índices extra), `corpus.db` en WAL y `resultados.db` en DELETE, lectura con `mode=ro&immutable=1`.
- [x] Split temporal global del TE §5.2.2: corte 2026-04-06, 58,400 / 14,600, `max(train.fecha) < min(test.fecha)` en todo el corpus, train por categoría entre 79.5 % y 80.3 %.

---

## 1. Metodología TDD

Las pruebas tienen estilo TDD: oráculos independientes, verificación directa sobre la BD, determinismo y "todo o nada". La Fase 1 se escribió antes de versionar el repositorio y llegó completa en el commit inicial, así que no tiene evidencia del ciclo rojo-verde-refactor que promete el §5.2.1 del Trabajo Escrito. Esa evidencia existe desde `fix/cover-build` y en `fase1/apego-te`.

- [x] **[ALTA]** Versionar el repositorio con git.
- [ ] **[ALTA]** Mantener el ciclo explícito `test: … (red)` → `feat: … (green)` → `refactor: …` en todas las ramas. Es la evidencia de TDD que pueden pedir en la defensa oral. No reescribir el historial para fabricarla.
- [ ] **[ALTA]** Nota de transparencia (en `docs/PLAN_PRUEBAS.md`): la Fase 1 se desarrolló con pruebas a la par antes de versionar; el ciclo estricto es verificable desde `fix/cover-build`.
- [x] **[ALTA]** Cobertura ≥ 80 % por paquete (`make cover-gate`).
- [ ] **[MEDIA]** Cubrir `ValidationError.Error()` y las ramas restantes de `Run`, `Persist` y `verifyCorpus`.
- [ ] **[MEDIA]** Rastrear las pruebas a los IDs del plan (CPU/CPP/CPE/CPA): nombres como `TestCPU01_...` o un comentario con el ID, más una matriz ID→prueba en `docs/PLAN_PRUEBAS.md`. Hace falta para la verificación bidireccional de la Fase 5 (A1).
- [x] **[MEDIA]** Criterios de éxito ejecutables: `make verify` y sección "Criterios de éxito" en `CLAUDE.md`.
- [x] **[MEDIA]** Agente `mapeo-reuso` para mapear código reutilizable antes de cada tarea.
- [ ] **[BAJA]** Reemplazar `testing.Short()` por la etiqueta de compilación `//go:build integration` en las pruebas de nivel 4.
- [x] **[BAJA]** `VACUUM` y `ON DELETE RESTRICT` solo los cubre `make corpus-sha` (sobreviven a una mutación en `go test`).

---

## 2. Brechas — prioridad alta

- [x] **El split violaba el TE §5.2.2** ("todos los comentarios del conjunto de entrenamiento son temporalmente anteriores a todos los del conjunto de prueba"): era temporal por categoría y globalmente el train llegaba al 2026-04-10 con el test empezando el 2026-04-03. Corregido: corte global (`db.SplitCutoff`).
- [x] **URI de archivo inválida en Windows.** `file://C:%5C…` que SQLite rechaza ("invalid uri authority"); afectaba a `gencorpus.exe`, `OpenReadOnly` y el `ATTACH`. Corregido: `db.FileURI` produce `file:///C:/…`.
  - [ ] Ejecutar `dist/gencorpus.exe` en Windows real o con Wine (no está instalado) y comprobar que el SHA del corpus es el mismo.
- [ ] **Duplicados exactos entre train y test.** 2,860 de 14,600 comentarios de test (**19.6 %**) tienen su texto idéntico en train (355 de ellos en Reconocimientos); ningún texto se repite entre categorías distintas. El reporte de validación lo mide. Esto infla CPM01/CPM02.
  - [ ] En la Fase 3, reportar además el Macro-F1 y el recall de Reconocimientos sobre el subconjunto de test **sin duplicado en train**.
  - [ ] Evaluar si conviene deduplicar antes de evaluar, y registrar la decisión en ADR-004.
- [ ] **La validación de "no trivialmente separable" está mal medida.** El Jaccard de las 100 palabras más frecuentes se calcula *sin quitar palabras funcionales* (de, la, mi, me…), así que el rango 0.22–0.34 refleja sobre todo esas palabras.
  - [ ] Recalcular el solapamiento después de aplicar la stoplist v1 (Fase 2).
  - [ ] Agregar una prueba de cordura: si un Naive Bayes trivial supera Macro-F1 0.99, el corpus es demasiado fácil y hay que recalibrar las plantillas o el ruido.
- [x] **Abreviaturas del ruido que la stoplist no contemplaba** (`q` 5,570 y `x` 3,609 comentarios). Decidido en `PLAN.md` (ADR-003): `q` → conjunción, `x` → preposición, junto con `segun` y `ole`.
  - [ ] Implementarlo en `v1.json` y CPU14 (Fase 2).

---

## 3. Brechas — prioridad media

- [ ] **Procedencia del corpus.** El TE §6.1.4 dice que la tasa de ruido "se documenta en el campo `parametros` de la tabla `corridas`". No se agrega nada a `corpus.db` (sería una desviación del Anexo I).
  - [ ] En la Fase 2, `cmd/simulador` escribe en `corridas.parametros` la semilla, el `NoiseConfig` y el total del corpus, además de la versión de la stoplist.
- [x] **El generador no se configura por línea de comandos.** El diseño dice que el ruido es configurable; `main.go` solo acepta `-out` y `-report`.
  - [x] Agregar flags `-seed`, `-total` y `-noise-rate` (los valores por defecto deben dar el mismo SHA del corpus).
  - [ ] Necesario para la prueba de linealidad I-4 (barridos de 1k/10k/30k/50k/73k).
- [ ] **Revisión manual de la muestra.** §6.1.4 dice "detección manual"; lo implementado es un chequeo de cordura automático más la impresión de la muestra.
  - [ ] Revisar a mano los 250 comentarios y dejar constancia (fecha, revisor, observaciones) en `docs/reports/`.
- [x] **ADR-007 no hace falta.** El diseño del TE (WAL, `go:embed`, `mode=ro&immutable=1`, `ATTACH`) funciona con modernc; se verificó ejecutándolo. ADR-002 y ADR-006 se retiraron por el mismo criterio (ver "Apego al Trabajo Escrito" en `CLAUDE.md`).
- [x] **`CHECK` sobre `frecuencias.categoria_id`**: descartado, sería una desviación del Anexo I; la integridad se verifica con `ATTACH`.
- [x] **`lexemas.forma_canonica UNIQUE` frente a "unicidad dentro de una corrida"**: se queda como en el Anexo I; cada `resultados.db` contiene una sola corrida, lo que también hace correcta la consulta I-2 literal.

---

## 4. Brechas — prioridad baja

- [ ] `internal/corpus/pipeline.go`: `defer f.Close()` ignora el error al cerrar el reporte, así que un reporte truncado pasaría sin aviso. Revisar el error de `Close`.
- [ ] `internal/corpus/pipeline.go`: si `verifyCorpus` falla después de `Persist`, el `corpus.db` se queda en disco. Borrarlo para mantener el "todo o nada".
- [ ] `internal/corpus/noise.go`: `typo()` sobre una letra mayúscula siempre termina duplicándola, porque `qwertyNeighbors` solo tiene minúsculas. Normalizar con `unicode.ToLower` antes de buscar.
- [x] `CHECK (fecha GLOB '????-??-??')` (Anexo I) acepta fechas imposibles como `2024-13-99`; solo la batería de validación revisa fechas reales. Limitación conocida, sin cambio de DDL.
- [x] `corpus.db` movido a `data/`.
- [x] `.gitignore` (`dist/`, `*.out`, `*.test`, `*.db-wal`, `*.db-shm`, `*.db-journal`; `data/corpus.db` sí se versiona).
- [x] `README.md` y `CLAUDE.md`.

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
│   ├── db/                        # runtime, solo lectura: OpenReadOnly, FileURI, SplitCutoff, Categories
│   │   ├── schema/corpus.sql      # DDL como .sql + go:embed = artefacto literal del Anexo I
│   │   ├── schema/resultados.sql
│   │   └── build/                 # create, Finalize, CheckCategoriaIntegrity (ATTACH)
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

Checklist (cada paso es un refactor **sin cambio de comportamiento**: el SHA-256 de `corpus.db` debía seguir siendo el de entonces, `c0343583…`, y así fue):

- [x] **Mover el generador de `cmd/gencorpus` a `internal/corpus`** y dejar `cmd/gencorpus/main.go` solo con los flags.
  - Por qué: `package main` no se puede importar, y el generador (~1,300 líneas) se necesita desde la prueba I-4, desde las pruebas del barrido de la Fase 2 y desde `--selftest`.
- [x] **Separar `templates.go` y `vocab.go`** en `internal/corpus/templates/` (datos separados de la lógica).
- [x] **Dividir `internal/db`:**
  - [x] `internal/db` (ejecución): `Categories`, `OpenReadOnly`, `FileURI`, `SplitCutoff`.
  - [x] `internal/db/build` (solo generación): `create`, `CreateCorpus`, `CreateResultados`, `InsertCategorias`, `Finalize`, `CheckCategoriaIntegrity`.
  - Por qué: la creación de archivos queda fuera del binario por la estructura misma, porque `cmd/autoatend` nunca importa `build`.
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

1. [x] Versionar el repositorio (§1).
2. [x] Subir la cobertura de `internal/db` a ≥ 80 % (§1).
3. [x] Reestructura de §5. Quedan para la Fase 4 `internal/assets/` y el chequeo `go list -deps`.
4. [x] Alinear la Fase 1 con el TE (`fase1/apego-te`): WAL, DDL literal, `mode=ro&immutable=1`, split global, URI de Windows, retiro de ADR-002 y ADR-006.
5. [ ] Verificar `gencorpus.exe` ejecutando en Windows (§2).
6. [x] Flags del generador (§3).
7. [ ] Arrancar la Fase 2 con commits rojo-verde-refactor, pruebas nombradas con su ID del plan y la stoplist según ADR-003.
8. [ ] En la Fase 3: métrica sin duplicados y prueba de cordura con un modelo trivial (§2).
