# AutoAtend (CENEVAL caso 10433): plan de ejecución hacia la demo

**Demo:** 13 de octubre de 2026. **Reescrito:** 11 de octubre de 2026.
**Objetivo:** llegar a la evaluación oral con un sistema funcionando que siga el guion del TE §6.4.3 y el orden del documento, con resultados medidos para la sección "Resultados" de la presentación.
**Plataformas:** Windows x86-64 (equipo del CENEVAL: Core i5, 8 GB) y macOS ARM64 (Mac propia, que el CENEVAL autorizó).

## 0. Fuentes de verdad

En orden de precedencia:
1. `../documentos/2907_10433.pdf`: enunciado del caso.
2. `../documentos/TrabajoEscrito_Consolidado.pdf` (TE): diseño entregado a los sinodales, congelado.
3. Este plan.

La guía del Acuerdo 286 ("Exposición y réplica oral") fija lo que se evalúa:
- La presentación dura 20 minutos y lleva portada, diagnóstico, resolución de las consignas y resultados.
- Lo que se diseñó se presenta funcionando.
- Los sinodales evalúan el fundamento teórico-metodológico, la pertinencia y la consistencia entre el trabajo escrito, la exposición y la réplica.

Por eso **cada diferencia con el TE es una pregunta potencial**: se evita, y cuando no se puede evitar se declara en §1.3 con su respuesta preparada.

---

## 1. Alineación caso ↔ TE

### 1.1 Consignas

| Consigna | TE | Artefacto | Estado |
|---|---|---|---|
| a) Base de datos del corpus | §6.1, Anexo I | `data/corpus.db` (`cmd/gencorpus`) | ✅ Fase 1 |
| b) Programa analizador | §6.2 | `internal/lexer`, `internal/regex`, `cmd/simulador` → `data/resultados.db` | Bloques A y B |
| c) Plan de pruebas | §6.3, Anexo II | El TE mismo, más las pruebas en el repo y `docs/CASOS_PRUEBA.md` (matriz ID → prueba, porque el Anexo II dice "la tabla completa… se encuentra en el repositorio") | Bloques A–E |
| d) Justificación del plan | §6.3.1, §6.3.8 (ISO/IEC 25010) | El TE mismo; las cifras medidas van en §5.2 | Bloque E |
| e) Servicio web | §6.4 | `cmd/autoatend` (`internal/ml`, `internal/web`, `internal/assets`) | Bloques C y D |

Solo la consigna e se exige "presentar el día de su evaluación oral". La b la alimenta, y el paso 4 del guion conecta las dos.

### 1.2 Política de desviaciones

Se implementa lo que dice el TE. Si no se puede cumplir tal cual, o si parece haber una opción mejor, se valida primero con el usuario. La consulta incluye la sección afectada, el motivo, las alternativas fieles al TE que se evaluaron, la evidencia y una recomendación. Solo después se implementa y se registra en el ADR que el TE promete para ese tema (ADR-001 SQLite, ADR-002 simulación del AFN-ε, ADR-003 stoplist, ADR-004 modelo, ADR-005 despliegue).

### 1.3 Desviaciones aprobadas (11-oct)

| Desviación | TE | Respuesta para la réplica | Registro |
|---|---|---|---|
| **Sin CI con GitHub Actions** (Capa 1) **ni Wine** (Capa 2) | §6.3.9 | Se sustituyen por ejecución nativa en Windows x86-64 y en macOS ARM64, que es evidencia más fuerte que un emulador o un runner remoto. La Capa 3 (`--selftest` en la máquina del examen) se conserva. Sin CGo no hay dependencias de DLL que detectar. | ADR-005 |
| **Los umbrales de §6.3.8 se miden y no bloquean** (I-4, P99, paridad, cobertura, CPM01/CPM02) | §6.3.8 | Cada umbral se mide con una prueba o un script reproducible y se reporta el valor observado. Durante el sprint solo bloquean I-1, I-2 y `--selftest`, que son los que protegen la corrección de lo que se muestra. | §5.2 de este plan |
| **No hay documento `docs/PLAN_PRUEBAS.md` aparte** | — (solo del plan anterior) | El plan de pruebas y su justificación son el TE §6.3. En el repo vive la matriz de casos (`docs/CASOS_PRUEBA.md`) que pide el Anexo II. | — |

Se conserva, ya aprobado: la stoplist de ADR-003, el IDF solo de train, el embed escrito a temporal con `mode=ro&immutable=1` y la I-2 con la consulta literal (ver "Apego al Trabajo Escrito" en `CLAUDE.md`).

---

## 2. Arquitectura y despliegue

- **Software base** (`internal/regex`, `internal/lexer`, `cmd/simulador`), escrito a mano en Go:
  - parser de ER por descenso recursivo (operadores `| * + ? ( ) [ ]` y escape `\`, gramática del TE §6.2.3) → AST;
  - construcción de Thompson (≤ 2·|r| estados);
  - simulación del AFN-ε por ε-cerradura y movida, O(n·m);
  - normalizador NFC a minúsculas y tokenizador AFD de 4 estados;
  - stoplist curada (homógrafos retenidos), tabla de símbolos hash O(1);
  - generador de ER por prefijo común con validación round-trip;
  - índice de prefijo de 3 caracteres (§6.2.5).

  **Sin `import "regexp"`, sin generadores de parsers.**
- **Base de datos** (`internal/db`, `cmd/gencorpus`): dos archivos SQLite en 3FN con `modernc.org/sqlite` (Go puro, sin CGo). La DDL es literal del Anexo I. `corpus.db` usa WAL y `resultados.db` usa DELETE. La integridad cruzada de `categoria_id` se verifica con `ATTACH`.
- **Servicio web** (`cmd/autoatend`, `internal/web`, `internal/ml`, `internal/assets`):
  - `net/http` en `127.0.0.1:8471`, con fallback 8472–8480 y banner con la URL efectiva;
  - endpoints `GET /` (HTML embebido, sin CDN), `POST /api/clasificar`, `GET /api/lexemas` (paginado), `GET /api/categorias` y `GET /api/health`;
  - inferencia nativa sobre `modelo.json` disperso;
  - `--selftest` sin interfaz (< 3 s, código de salida 0/1).
- **Compilación cruzada** desde cualquier host, sin toolchain de C:

  ```bash
  # Principal: Windows x86-64
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.version=1.0.0" \
    -o dist/autoatend.exe ./cmd/autoatend

  # Respaldo: macOS ARM64 (Apple Silicon)
  CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
    go build -trimpath -ldflags="-s -w -X main.version=1.0.0" \
    -o dist/autoatend-macos-arm64 ./cmd/autoatend
  ```

- **Arranque:**
  - Las BD embebidas (`go:embed`, TE §7) se escriben en un directorio propio del proceso bajo el temporal del usuario (`os.MkdirTemp`, sin permisos de administrador). Se abren con `db.OpenReadOnly` (`mode=ro&immutable=1`, TE §5.1.4) y con URI `file:///C:/…` en Windows (`db.FileURI`).
  - Si escribir o abrir falla, el programa termina al arrancar con un diagnóstico. El directorio se borra al salir.
  - El único `runtime.GOOS` permitido es la apertura del navegador.
  - En macOS, un binario copiado desde otra máquina puede requerir `xattr -d com.apple.quarantine` (nivel usuario).

**Flujo de datos (camino crítico):** `corpus.db` → lexer → `resultados.db` y DF de train → modelo → inferencia → servicio. Por eso los bloques van en orden. El entrenamiento en Python (C) puede correr mientras se programa B.

---

## 3. Fase 1 (consigna a): cerrada

Entregado y verificado con `make verify` en verde:
- `data/corpus.db`: semilla 10433, 73,000 comentarios, distribución 30/20/20/20/10, del 2024-08-30 al 2026-08-29.
- SHA-256 `37f285ff1da70f4bd5205c51e255a612d1a72ad45b70fc9fda35dbac1b3e6e63`.
- Reporte en `docs/reports/generator_validation_report.md`.
- Corpus regenerado el 12-oct por defectos de coherencia detectados en la revisión de la muestra (TE §6.1.4): un solo problema por comentario (base + complemento de la misma categoría), precios de plan coherentes con los cobros, huecos que concuerdan dentro de la oración y aperturas/cierres por categoría. Pruebas `TestCoherence_*` en `internal/corpus`.
- Split temporal global `db.SplitCutoff`: corte 2026-04-06, 58,400 / 14,600.
- DDL literal del Anexo I, `mode=ro&immutable=1`, `db.FileURI` para Windows, ADR-001.
- `gencorpus` acepta los flags `-seed`, `-total` y `-noise-rate` (los usa la medición I-4).

**Hechos que afectan las fases siguientes:**
- El 24.4 % del test (3,560 de 14,600) tiene un texto idéntico en train, lo que infla CPM01/CPM02. Las métricas también se reportan sin duplicados (Bloque C).
- El ruido escribe `q` por "que" y `x` por "por" (6,007 y 3,152 comentarios). Ya están cubiertos en ADR-003.
- El Jaccard del reporte se calcula sin stoplist, así que el rango 0.19–0.34 refleja sobre todo palabras funcionales. Se recalcula con la stoplist v1 (Bloque A, P2).

**Criterio de separabilidad registrado antes de tocar las plantillas (12-oct-2026, aprobado por el usuario).** Con el corpus `b810bad2…` el modelo obtuvo Macro-F1 = 1.0000 y una matriz diagonal, lo que contradice el TE §6.1.4 ("evita que el corpus sea trivialmente separable"). Antes de cambiar el generador se fijan tres compuertas de la batería, que no se modifican después. Si alguna falla, el corpus no se persiste:

| Id | Medida | Rango aceptado |
|---|---|---|
| K1 | Jaccard de los 100 tokens más frecuentes por par de categorías, **después de la stoplist v1** | [0.15, 0.80] (el rango que ya existía) |
| K2 | Colisión de la oración base: % de comentarios cuya oración base limpia, como bolsa de tokens sin las palabras de la stoplist v1, aparece también en otra categoría | global [0.05, 0.10]; ≥ 0.03 por categoría; ≤ 0.05 en Reconocimientos |
| K3 | Sonda Naive Bayes multinomial (α = 1, 5 pliegues estratificados, semilla 10433, tokens de la batería sin las palabras de la stoplist v1): Macro-F1 promedio | [0.80, 0.99] |

Otras reglas:
- Son solo de reporte, no compuertas: la regla de una palabra y el F1 con plantillas no vistas.
- No se fija ningún umbral sobre el modelo final: se entrena una sola vez, al final, y se reporta tal cual.
- No hay ruido de etiquetas. La ambigüedad es controlada: textos compatibles con las dos etiquetas, de forma determinista y con una tasa declarada.
- El commit de prueba en rojo demuestra que el corpus `b810bad2…` falla K1, K2 y K3.

**Corpus regenerado el 12-oct con el criterio aplicado (`fix/corpus-separabilidad`).**
- Nuevo SHA `37f285ff…`. Las plantillas por categoría quedan en 45 / 44 / 44 / 47 / 46, con 6 de contexto cruzado por categoría, 10 compartidas en los pares del plan y el marco "(ref. {FOLIO})" en las cinco categorías.
- Batería: K1 [0.156, 0.299], K2 0.0895 global (0.0661 / 0.0928 / 0.0858 / 0.1507 / 0.0381) y K3 0.9726. Pasa en la iteración 3 de 4; el registro está en "Iteraciones del generador" del reporte.
- El subconjunto ambiguo son 6,533 comentarios (8.95 %), con sus ids en `data/corpus_ambiguos.json`.
- Semilla, filas, distribución e inyector de ruido no cambian.

**Deuda conocida, no bloqueante (después de la demo):**
- `internal/corpus/pipeline.go`: `defer f.Close()` ignora el error al cerrar el reporte.
- `internal/corpus/pipeline.go`: si `verifyCorpus` falla, `corpus.db` se queda en disco.
- `internal/corpus/noise.go`: `typo()` con mayúscula siempre la duplica. **No se corrige antes de la demo**: cambia 38,582 textos y el SHA del corpus.
- `ValidationError.Error()` sin cobertura (`Run` 78.6 %, `verifyCorpus` 72.2 %).
- Integración de nivel 4 con `testing.Short()` en lugar de la etiqueta `//go:build integration`.

---

## 4. Plan hasta la demo

**Prioridades:** **P0** se necesita para la demo · **P1** es compromiso del TE y barato · **P2** solo si sobra tiempo. Si un bloque se atrasa, se recorta P2 y luego P1, nunca P0.

**Método:**
- Cada pieza en su rama (`fase2/…`, `fase3/…`, `fase4/…`) con al menos el par de commits `test: … (red)` → `feat: … (green)` (TE §5.2.1). El `refactor:` solo cuando haga falta.
- Las pruebas llevan el ID del Anexo II (`TestCPU01_…`).
- Merge `--no-ff` con aprobación del usuario.

### Bloque A: lexer y barrido (11-oct, Fase 2, consigna b)

| Prioridad | Tarea | TE | Pruebas |
|---|---|---|---|
| P0 | Normalizador NFC a minúsculas con el alfabeto `[a-záéíóúüñ]`, dígitos y separadores; tokenizador AFD de 4 estados | §6.2.2 | CPE01–CPE04 (CE-1…CE-7) |
| P0 | `internal/lexer/stoplist/v1.json` (forma → clase cerrada) con la política de ADR-003 (abajo) | §6.2.3 | CPU11–CPU13, CPU14 (extensión) |
| P0 | Tabla de símbolos: contadores global y por categoría, conjunto de variantes | §6.2.2 | CPU09, CPU10 |
| P0 | Generador de ER por prefijo común con alternancia de sufijos (`factur(a\|as\|ación)`) | §6.2.2 | CPP01 (con el round-trip del Bloque B) |
| P0 | `cmd/simulador`: barre `corpus.db` y persiste `corridas`, `lexemas`, `frecuencias` (sin filas con `valor = 0`) y `variantes` en `data/resultados.db`. Escribe en `corridas.parametros` la semilla, el `NoiseConfig`, el total del corpus y la versión de la stoplist (TE §6.1.4). Escritura determinista | §6.2, §6.1.3, §6.1.4 | Integración 4b |
| P0 | Export de la frecuencia de documento (DF) de train para Python. El catálogo de `resultados.db` usa el corpus completo; el vocabulario y el IDF del modelo salen solo de train. Los términos con df = 0 en train se excluyen del modelo y se reporta cuántos son | §6.4.2 | — |
| P1 | Índice de prefijo de 3 caracteres | §6.2.5 | Unitaria |
| P1 | Borradores para el usuario: muestra L-1 de 200 tokens homógrafos (bajo, para, sobre, entre, como) con la etiqueta propuesta, y los 250 comentarios de la muestra del generador | §6.3.10, §6.1.4 | — |
| P2 | Jaccard recalculado con la stoplist v1 | §6.1.4 | — |

**Política de la stoplist (ADR-003, ya aprobada):**
- **Se filtran:**
  - artículos: el, la, los, las, un, una, unos, unas, lo; contracciones: al, del;
  - preposiciones: de, en, con, por, a, ante, tras, durante, mediante, según, sin, hasta, desde, contra, hacia;
  - conjunciones: y, e, o, u, ni, que, pero, sino, aunque, porque, pues;
  - interjecciones: ah, eh, ay, uf, olé;
  - variantes de superficie: `segun`, `ole`, `q` (conjunción) y `x` (preposición).
- **Se retienen los homógrafos del TE:** bajo, para, sobre, entre, como, más, solo, aun.
- **No se filtran** `se`, `no`, `este` y `eso`: no pertenecen a las 4 clases.
- **No hay plegado global de acentos.**
- **CPU14** verifica:
  - `segun` y `ole` quedan en la misma clase que `según` y `olé`; `q` → conjunción; `x` → preposición; `hacia` → preposición;
  - no se filtran `él`, `qué`, `sí`, `té`, `se`, `no`, `este` ni `eso`.

**Criterio de salida:** `data/resultados.db` generado, I-2 = 0 filas y `make verify` en verde.

### Bloque B: motor de ER propio e invariantes (12-oct, mañana, Fase 2, consignas b y c)

| Prioridad | Tarea | TE | Pruebas |
|---|---|---|---|
| P0 | `internal/regex`: parser (gramática `expresion→alternancia→concatenacion→factor→primario`, AST tipado, errores descriptivos, límite de profundidad) | §6.2.3 | CPU01–CPU03 |
| P0 | Construcción de Thompson (reglas de composición de fragmentos) | §6.2.4 | CPU04, CPU05 |
| P0 | Simulación por ε-cerradura y movida; sondas empíricas de O(n·m) con 10, 100, 1k y 10k caracteres | §6.2.4 | CPU06–CPU08 |
| P0 | **I-1**: round-trip de cada ER del catálogo (re-parsear → Thompson → aceptar todas las variantes observadas); 0 fallos | §6.3.3 | CPP01 |
| P0 | **I-2**: la consulta literal del TE §6.3.3 devuelve 0 filas; prueba de que `corridas` tiene una sola fila | §6.3.3 | CPP02 |
| P1 | **I-3**: dos barridos producen el mismo SHA-256 de `resultados.db` | §6.3.3 | CPP03 |
| P1 | Integración 4b: sin lexemas huérfanos, cada lexema con al menos una variante, cada ER parsea con el parser propio | §6.3.5 | Integración |
| P1 | Borrador del corpus de oro: 100 comentarios de test con sus lexemas y variantes esperados, para que el usuario lo valide (oráculo de §6.3.5c) | §6.3.5 | Integración 4c |
| P1 | L-2: proporción de lexemas del catálogo con una sola variante observada | §6.3.10 | Medición |

**Criterio de salida:** I-1 = 0 fallos, I-2 = 0 filas y `make verify` en verde.

### Bloque C: modelo (12-oct, tarde, Fase 3, consigna e)

| Prioridad | Tarea | TE | Pruebas |
|---|---|---|---|
| P0 | `ml/` (Python 3.11 y scikit-learn, venv aislado) con el split de `db.SplitCutoff` (misma consulta SQL, sin volver a partir). TF-IDF con vocabulario e IDF de train y normalización L2. El vectorizador invoca el tokenizador Go como subproceso | §6.4.2, §5.2.2 | — |
| P0 | NBM y RLS, cada uno con las 3 estrategias de desbalance (`class_weight='balanced'`, sobremuestreo, línea base) y CV estratificada de 5 particiones; el ganador se elige por Macro-F1 promedio | §6.4.2, §5.1.3 | — |
| P0 | Evaluación en test: Macro-F1, recall por clase y matriz de confusión normalizada. **También sobre el test sin duplicados en train.** Si NBM supera Macro-F1 0.99, se registra como señal de que el corpus es demasiado fácil | §5.2.2 | CPM01, CPM02 (medidos) |
| P0 | Export de `modelo.json` (pesos dispersos, bias, IDF, mapa de vocabulario, métricas de CV) y del perdedor como `modelo_alternativo.json` | §6.4.2 | — |
| P0 | `internal/ml`: carga del JSON e inferencia (tokenizar → vocabulario → TF-IDF disperso → `z_k = b_k + Σ w_{k,i}·x_i` sobre el soporte → softmax con resta del máximo → `lexemas_clave` top-5 por `\|w·x\|`) | §6.4.2 | Unitarias contra puntajes calculados en Python |
| P0 | `gold5.json`: 5 comentarios de test bien clasificados (al menos un caso claro y un caso frontera para el guion) | §6.3.9, §6.4.3 | — |
| P1 | Paridad Python ↔ Go sobre 1,000 comentarios de test (vectores con L1 = 0) | §6.3.6 | CPA06 (medido) |
| P1 | ADR-004: elección del modelo y JSON en vez de pickle | §6.4 | — |

**Criterio de salida:** `modelo.json` cargado en Go y los 5 comentarios de `gold5` correctos en Go.

### Bloque D: servicio web (12-oct, noche, Fase 4, consigna e)

| Prioridad | Tarea | TE | Pruebas |
|---|---|---|---|
| P0 | `internal/assets` con `go:embed` de `corpus.db`, `resultados.db`, `modelo.json`, `gold5.json` y `ui/`; el objetivo `make assets` los copia (nunca se editan a mano) | §7 | — |
| P0 | `cmd/autoatend` e `internal/web`: los 5 endpoints. Validación: `texto` vacío → 400 con mensaje; cuerpo demasiado grande → 413; solo consultas parametrizadas. Respuesta `{"categoria", "confianzas"{5}, "lexemas_clave"[5]}` con confianzas que suman 1.0. `/api/health` → `{"ok", "version", "lexemas", "puerto"}` | §6.4.1, §6.3.7 | CPA01–CPA04 |
| P0 | Arranque: BD a temporal + `OpenReadOnly`, `ATTACH` para la integridad de `categoria_id`, puerto 8471 con fallback 8472–8480, banner y apertura del navegador | §6.4.1, §6.1.3 | CPA05 |
| P0 | Interfaz HTML única: caja de texto, categoría, las 5 confianzas, `lexemas_clave` y enlace a `/api/lexemas` | §6.4.1 | Manual |
| P0 | `--selftest` (< 3 s): integridad y conteos de `corpus.db`, validez sintáctica de todas las ER del catálogo con el parser propio, `gold5`, puerto o fallback → salida 0/1 con diagnóstico por verificación | §6.3.9, §6.4.3 | Pruebas de código de salida |
| P0 | `make cross` con las líneas de `autoatend` activadas | §2 | — |
| P1 | Latencia: 100 llamadas a `/api/clasificar`, P99 | §6.3.7 | Medido |
| P1 | ADR-002 (simulación directa del AFN-ε) y ADR-005 (solo loopback, sin TLS, autenticación ni CORS; desviación de §6.3.9) | §6.2.3, §6.4 | — |
| P2 | Chequeo `go list -deps`: `cmd/autoatend` no depende de `internal/db/build` | — | — |

**Criterio de salida:** `dist/autoatend.exe` y `dist/autoatend-macos-arm64` generados y `--selftest` = 0 en Linux.

### Bloque E: verificación nativa y ensayo (13-oct, mañana, Fase 5)

| Prioridad | Tarea |
|---|---|
| P0 | **Windows (partición propia):** `gencorpus.exe` reproduce el SHA `37f285ff…`; `autoatend.exe --selftest` = 0; el guion completo de 5 pasos |
| P0 | **Mac:** lo mismo con `gencorpus-macos-arm64` y `autoatend-macos-arm64` |
| P1 | I-4: barrido con `gencorpus -total` de 1k, 10k, 30k, 50k y 73k → R² y T(73k) (CPP04, medido) |
| P1 | Cobertura por paquete (`make cover`, se informa) |
| P1 | `docs/CASOS_PRUEBA.md`: matriz ID del Anexo II → prueba (archivo y función) → resultado |
| P0 | Llenar §5.2 con las cifras y pasarlas a las diapositivas de "Resultados"; ensayar el guion |

### Tareas del usuario (en paralelo, entre 1 y 2 h)
- Validar y corregir el borrador L-1 de 200 tokens y calcular la tasa de sub-filtrado (TE §6.3.10).
- Validar el corpus de oro de 100 comentarios (TE §6.3.5c).
- Revisar la muestra de 250 comentarios y dejar constancia (fecha, revisor, observaciones) en `docs/reports/` (TE §6.1.4).

---

## 5. Criterios de aceptación de la demo

Esta es la checklist viva: marca `[x]` solo lo que se verificó ejecutándolo.

### 5.1 Guion del TE §6.4.3

Decisión del usuario (12-oct): la demo se presenta en la **Mac propia** (autorizada por el CENEVAL). La columna Windows queda como respaldo opcional.

| Paso | Windows | Mac |
|---|---|---|
| 1. Doble clic → banner con URL → se abre el navegador | [ ] | [ ] |
| 2. Clasificar en vivo un comentario claro y uno frontera | [ ] | [ ] |
| 3. `lexemas_clave` explica la decisión con lexemas del catálogo | [ ] | [ ] |
| 4. `/api/lexemas` muestra 3 ER del catálogo (conecta la consigna e con la b) | [ ] | [ ] |
| 5. `--selftest` → salida 0 en < 3 s | [ ] | [ ] |
| `gencorpus` reproduce el SHA del corpus | [ ] | [ ] |

### 5.2 Mediciones (TE §6.3.8 y §6.3.10)

| Criterio | Umbral (TE) | Medido | Fecha |
|---|---|---|---|
| I-1 round-trip | 0 fallos sobre el catálogo completo | 0 fallos (6,196 lexemas, 6,891 variantes; `TestCPP01_I1_RoundTripCatalogue`; corpus 37f285ff…) | 2026-10-12 |
| I-2 conservación de frecuencias | 0 lexemas con discrepancia | 0 filas (consulta literal del TE §6.3.3 sobre `data/resultados.db`; corpus 37f285ff…) | 2026-10-12 |
| I-3 idempotencia | SHA-256 idéntico en 2 corridas | Idéntico: `9620196f…79f3` (`make sweep-sha`; corpus 37f285ff…) | 2026-10-12 |
| I-4 linealidad | R² ≥ 0.99; T(73k) < 60 s | | |
| Cobertura | ≥ 80 % (líneas) | Fases 1–2: corpus 91.8 %, db 87.5 %, db/build 90.8 %, db/schema 100 %, lexer 99.4 %, regex 98.7 %, sweep 82.3 % (`make cover`) | 2026-10-12 |
| Paridad Python ↔ Go | 0 discrepancias en 1,000 | 0 discrepancias: mismo soporte y L1 = 0 (igualdad de bits) en los 1,000 vectores, clase idéntica en 1,000/1,000, máx. \|Δprob\| = 1.1e-16 (`TestCPA06_ParidadPythonGo`, `testdata/paridad_1000.json`; corpus 37f285ff…) | 2026-10-12 |
| Macro-F1 (CPM01) | ≥ 0.85 en test | 0.9795 en test (14,600) con `modelo.json` (RLS, línea base, C = 1); NBM 0.9708. Un solo entrenamiento, después de que la batería K1–K3 pasó (corpus 37f285ff…). El corpus anterior b810bad2… daba 1.0000 (trivialmente separable) | 2026-10-12 |
| Recall de Reconocimientos (CPM02) | ≥ 0.70 | 0.9819 (1,495 Reconocimientos de test); NBM 0.9746 | 2026-10-12 |
| Macro-F1 y recall sin duplicados en train | se reporta | Test sin duplicados (11,040): RLS 0.9826 y recall 0.9877; NBM 0.9726 y 0.9793. Test sin el subconjunto ambiguo (13,348): RLS 0.9997 y 0.9986; NBM 0.9952 y 0.9910. Plantillas no vistas (solo reporte, `docs/reports/plantillas_no_vistas.md`): RLS 0.8350 de media (0.7637–0.8855), NBM 0.7592 | 2026-10-12 |
| Latencia P99 de `/api/clasificar` | < 50 ms | 0.11–0.17 ms en Linux (5 corridas de `TestLatencia_P99DeClasificar`: 100 llamadas HTTP consecutivas por loopback con un comentario de longitud media, 111 caracteres; mediana 0.04 ms). Pendiente medir en el hardware objetivo (Windows y Mac) en el Bloque E | 2026-10-12 |
| `--selftest` | < 3 s, salida 0 | Salida 0 en 0.05 s en Linux (`dist/autoatend --selftest`: 73,000 comentarios, 6,289 ER compiladas y 6,967 variantes aceptadas, `gold5` 5/5, puerto 8471 libre; `TestSelftestExit0` y las pruebas de salida 1). Pendiente medir en Windows y Mac (Bloque E) | 2026-10-12 |
| Revisión manual de la muestra (TE §6.1.4) | detección manual de inconsistencias | Pendiente para el corpus 37f285ff…: borrador `docs/reports/muestra_250_borrador.md` (Auto 250/250 ok; Tipo base 202, contexto cruzado 33, compartida 15) y `docs/reports/plantillas_nuevas.md` (47 plantillas para revisar). La revisión anterior (250/250 sin observaciones) era del corpus b810bad2… | |
| L-1 sub-filtrado de homógrafos | se reporta | 86.0 % (172 cerrada / 28 léxica, 0 duda; 200 tokens de test en los mismos 16 contextos, sin contextos nuevos que etiquetar; para 70, como 70, entre 60; bajo y sobre no aparecen en el corpus; las etiquetas del usuario se reaplicaron por contexto; corpus 37f285ff…) | 2026-10-12 |
| L-2 lexemas con una sola variante | se reporta | 91.4 % (5,660 de 6,196; la mayoría son erratas del ruido o formas irregulares; corpus 37f285ff…) | 2026-10-12 |

---

## 6. Riesgos y plan B

| Riesgo | Mitigación |
|---|---|
| El Bloque A (lexer) se atrasa y arrastra todo | Se recorta P2 y luego P1. El índice de prefijo y los borradores manuales son lo primero en caer |
| La interfaz web no se ve en la pantalla del sínodo | `--selftest` en la terminal (plan B del TE §6.4.3) |
| El equipo Windows del CENEVAL falla | La Mac propia con `autoatend-macos-arm64` |
| CPM01 queda por debajo de 0.85 | Se reporta el valor observado y se muestra `modelo_alternativo.json`; nunca se ajusta el umbral ni el test |
| Pregunta por Wine o el CI (§6.3.9) | Respuesta de §1.3 |
