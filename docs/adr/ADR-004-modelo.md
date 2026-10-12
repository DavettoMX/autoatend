# ADR-004: elección del modelo y formato JSON del modelo

**Estado:** aceptada (Fase 3, Bloque C). Es la decisión que el Trabajo Escrito (TE) §6.4 promete documentar.
**Implementación:** `ml/autoatend_ml` (entrenamiento offline en Python), `cmd/tokenizar` e `internal/ml` (tokenizador por lotes, carga del JSON e inferencia en Go).
**Evidencia:** `docs/reports/evaluacion_modelo.md` (generado por `make train`), `data/modelo.json`, `data/modelo_alternativo.json`, `testdata/paridad_1000.json`.
**Pruebas:** CPA06 (`TestCPA06_ParidadPythonGo`), equivalencia lineal con scikit-learn, `gold5`, y CPM01/CPM02 (medidos, no bloquean).

## Contexto

El TE §5.1.3 compara Naive Bayes multinomial (NBM) con regresión logística softmax (RLS) sobre TF-IDF, con tres estrategias de desbalance y validación cruzada estratificada; §6.4.2 fija que el vocabulario y el IDF son "precalculados durante el entrenamiento" y que el servicio calcula `z_k = b_k + Σ w_{k,i}·x_i`. El binario final es un único `.exe` sin Python, sin red y sin instalación (`PLAN.md` §1), así que el modelo debe viajar como dato y la inferencia debe ser nativa en Go. El corpus tiene un desbalance 30/20/20/20/10 y el 25.5 % del test repite textos de train.

## Decisión

### 1. Familia y estrategia

El modelo ganador es **regresión logística softmax con penalización L2, C = 10, estrategia de pesos balanceados (`class_weight='balanced'`), solver `saga`**, exportado como `data/modelo.json`. El mejor de NBM (sobremuestreo de Reconocimientos, α = 0.5) se exporta como `data/modelo_alternativo.json`, para mostrar a los sinodales la otra familia con sus propios pesos.

Evidencia (validación cruzada de 5 particiones estratificadas sobre train, Macro-F1 promedio ± desviación; detalle y rejilla completa en el reporte):

| Familia | Pesos balanceados | Sobremuestreo | Línea base |
|---|---|---|---|
| NBM (mejor α) | 0.99987 ± 0.00013 (α = 0.5) | 0.99989 ± 0.00014 (α = 0.5) | 0.99964 ± 0.00023 (α = 0.1) |
| RLS (mejor C) | 1.00000 ± 0.00000 (C = 10) | 1.00000 ± 0.00000 (C = 10) | 1.00000 ± 0.00000 (C = 10) |

En test (una sola evaluación, sin ajuste): RLS Macro-F1 = 1.0000 y recall de Reconocimientos = 1.0000; sin duplicados de train (10,884 de 14,600), igual. NBM: 0.9997 y 0.9996 sin duplicados, recall de Reconocimientos 1.0000. Los umbrales CPM01 (≥ 0.85) y CPM02 (≥ 0.70) se cumplen con holgura.

**Cómo leer esta evidencia con honestidad:**

- **El corpus es fácil.** NBM supera 0.99 de Macro-F1, que es la señal que el plan manda registrar: las 21 combinaciones de la rejilla quedan entre 0.9994 y 1.0000. La holgura sobre los umbrales describe al corpus sintético generado con plantillas, no la dificultad del problema real de TelArg.
- **Empate entre estrategias.** Las tres estrategias de desbalance de RLS empatan en 1.00000 ± 0.00000 en CV. El ganador lo decide el criterio de desempate fijado antes de entrenar (mayor media, luego menor desviación, luego primera de la rejilla en el orden pesos balanceados, sobremuestreo, línea base), no una ventaja medida. Con este corpus no se puede afirmar que una estrategia de desbalance sea mejor que otra; sí que RLS ≥ NBM en CV y en test, por una diferencia menor a 0.001.
- **C = 10 es el borde de la rejilla** (`C ∈ {0.1, 1, 10}`), fijada por el diseño y no ampliada. Como el resultado ya es perfecto en CV, ampliarla no cambiaría la elección verificable.
- **Fuga no supervisada en CV.** El vocabulario y el IDF se calculan una vez con todo train, como pide el TE, y no por pliegue; las estadísticas de documento de los pliegues de validación entran en el IDF. Puede inflar ligeramente el Macro-F1 de CV; no afecta a test, que nunca participa.
- Se descartaron 846 lexemas del catálogo con df = 0 en train (aparecen solo en test); el vocabulario del modelo tiene 5,443.

### 2. Un solo escalador lineal para las dos familias

Ambas familias se exportan como `z_k = b_k + Σ_{i∈supp(x)} w_{k,i}·x_i` seguido de softmax con resta del máximo:

- NBM: `w = feature_log_prob_`, `b = class_log_prior_`. El softmax de la verosimilitud conjunta es el posterior de NBM, es decir, `predict_proba`.
- RLS: `w = coef_`, `b = intercept_`.

Go tiene así una sola ruta de inferencia (`internal/ml`) y cambiar de modelo es cambiar de JSON. La equivalencia con `predict_proba` de scikit-learn se verifica sobre todo el test (diferencia máxima ≈ 1e-15 para ambos modelos; `TestModeloLineal_EquivalenciaSklearn` repite la comprobación en Go sobre los 1,000 comentarios de paridad).

### 3. JSON plano en lugar de pickle, joblib u ONNX (TE §5.1.3)

El modelo se serializa como JSON determinista: versión, familia, estrategia, hiperparámetros, categorías, vocabulario, IDF, pesos dispersos por clase (`[[índice, valor], ...]`, solo valores distintos de cero), sesgos y las métricas. Razones:

- **Pickle y joblib** ejecutan código arbitrario al cargarse, solo los lee Python y dependen de la versión exacta de scikit-learn. Son inaceptables para un binario Go sin Python y para un archivo que se incrusta y se muestra.
- **ONNX** necesita un runtime (en Go, CGo o una biblioteca fuera de la regla "solo Go puro"), y el modelo es lineal: no hay grafo que justifique un formato de intercambio.
- **JSON** se lee con `encoding/json` de la biblioteca estándar, se versiona en git y es legible por los sinodales. Con `repr` de `float64` (la representación más corta que reconstruye el mismo valor) el viaje Python → JSON → Go conserva cada bit; por eso el IDF se exporta tal como se calculó y **Go no recalcula ningún logaritmo**.
- El orden de las claves es fijo, no hay marcas de tiempo y dos entrenamientos producen archivos idénticos byte a byte (`make train-sha`).

### 4. Paridad Python ↔ Go (CPA06)

Python tokeniza invocando `cmd/tokenizar` (la misma `lexer.Analyze` que el barrido y el servicio) una sola vez por lote, y el tokenizador y `Model.Vectorize` comparten `ml.Lexicon`, así que no existen dos definiciones de "token". Para que los vectores coincidan bit a bit: TF = conteo crudo; `x_i = (tf_i·idf_i) / sqrt(Σ (tf_j·idf_j)²)` con la suma en orden ascendente de índice y operaciones `float64` simples; en Python se hace con bucles puros (numpy suma por pares) y en Go con conversiones explícitas `float64(a*b)`, que prohíben la multiplicación-suma fusionada que algunas arquitecturas aplicarían. Las probabilidades pasan por `exp` (libm en Python, `math.Exp` en Go), que puede diferir en el último bit; por eso los vectores se comparan con igualdad exacta (L1 = 0) y las probabilidades con tolerancia 1e-12, reportando el máximo observado.

### 5. Decisiones de implementación que el TE no fija

- **IDF suavizado:** `idf = ln((1 + N) / (1 + df)) + 1` con N = 58,400 y df de `df_train.json`. Siempre ≥ 1, de modo que ningún término del vocabulario tiene peso cero.
- **Sobremuestreo:** muestreo aleatorio con reemplazo de Reconocimientos hasta el tamaño mediano de clase (el 20 %), con semilla fija, solo sobre la parte de entrenamiento de cada pliegue y del ajuste final.
- **`saga`:** solver de gradiente promedio estocástico (de la familia del descenso de gradiente estocástico que menciona el TE), `max_iter = 300`, `tol = 1e-4`; todas las rejillas convergieron.
- **Texto sin lexemas del vocabulario** (vacío, solo palabras vacías o solo palabras desconocidas): el vector es cero, los puntajes son los sesgos y la predicción es la del *prior* del modelo; `Prediction.Vacio` lo señala para que la capa web lo trate (por ejemplo, avisando que no se reconoció ninguna palabra) en vez de presentarlo como evidencia.
- **`lexemas_clave`:** los cinco términos del soporte con mayor `|w_{k,i}·x_i|` para la clase ganadora, con empates por índice menor; no se listan términos con contribución cero.

## Alternativas descartadas

- **`TfidfVectorizer` de scikit-learn con su propio tokenizador:** habría dos definiciones de token (la de Python y la del analizador léxico propio de la consigna b) y la paridad CPA06 sería inalcanzable. Su fórmula de IDF suavizado es la misma que se usa aquí, pero el cálculo se hace sobre la salida del tokenizador Go.
- **Servir con Python (`pickle` + `predict_proba`):** contradice el sobre de despliegue (`.exe` único, sin runtime) y la regla de que Python es solo de desarrollo.
- **Exportar solo el ganador:** se conserva también el mejor de la otra familia como `modelo_alternativo.json`, porque el plan B ante una pregunta o un umbral incumplido es mostrarlo (`PLAN.md` §6).
- **Escalar a otros clasificadores (SVM, árboles, redes):** fuera del TE §5.1.3, que compara NBM y RLS; con Macro-F1 de CV de 1.0000 no habría qué ganar.
- **Calcular el IDF por pliegue:** más estricto en CV, pero el TE pide IDF "precalculados durante el entrenamiento" con un solo vocabulario; se declara la limitación en lugar de cambiar el diseño.

## Consecuencias

- `make train` reproduce los artefactos versionados y `make train-sha` comprueba que dos entrenamientos independientes producen los mismos bytes. `make verify` no requiere Python: las pruebas Go leen los JSON versionados.
- Si cambian el corpus, la stoplist o el stemmer, hay que regenerar `resultados.db` y `df_train.json` (`make clean-results sweep`) y volver a entrenar (`make train`); la prueba de paridad y la de `gold5` detectan un desajuste.
- Como el corpus es fácil, las cifras de CPM01/CPM02 deben presentarse a los sinodales con esa salvedad: miden el corpus sintético, no la operación real.
