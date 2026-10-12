# ADR-004: elección del modelo y formato JSON del modelo

**Estado:** aceptada (Fase 3, Bloque C). Es la decisión que el Trabajo Escrito (TE) §6.4 promete documentar.
**Implementación:** `ml/autoatend_ml` (entrenamiento offline en Python), `cmd/tokenizar` e `internal/ml` (tokenizador por lotes, carga del JSON e inferencia en Go).
**Evidencia:** `docs/reports/evaluacion_modelo.md` (generado por `make train`), `docs/reports/plantillas_no_vistas.md`, `data/modelo.json`, `data/modelo_alternativo.json`, `testdata/paridad_1000.json`.
**Pruebas:** CPA06 (`TestCPA06_ParidadPythonGo`), equivalencia lineal con scikit-learn, `gold5`, y CPM01/CPM02 (medidos, no bloquean).

## Contexto

El TE §5.1.3 compara Naive Bayes multinomial (NBM) con regresión logística softmax (RLS) sobre TF-IDF, con tres estrategias de desbalance y validación cruzada estratificada; §6.4.2 fija que el vocabulario y el IDF son "precalculados durante el entrenamiento" y que el servicio calcula `z_k = b_k + Σ w_{k,i}·x_i`. El binario final es un único `.exe` sin Python, sin red y sin instalación (`PLAN.md` §1), así que el modelo debe viajar como dato y la inferencia debe ser nativa en Go. El corpus tiene un desbalance 30/20/20/20/10 y el 24.4 % del test repite textos de train.

## Decisión

### 1. Familia y estrategia

El modelo ganador es **regresión logística softmax con penalización L2, C = 1, sin ajuste de desbalance (línea base), solver `saga`**. Se exporta como `data/modelo.json`. El mejor de NBM (sobremuestreo de Reconocimientos, α = 0.1) se exporta como `data/modelo_alternativo.json`, para mostrar a los sinodales la otra familia con sus propios pesos.

Las cifras son del corpus `37f285ff…`. Ese corpus se regeneró con el criterio de separabilidad K1–K3, registrado en `PLAN.md` §3 antes de cambiar las plantillas, y el modelo se entrenó **una sola vez**, después de que la batería pasó.

Evidencia: validación cruzada de 5 particiones estratificadas sobre train, Macro-F1 promedio ± desviación. El detalle y la rejilla completa están en el reporte.

| Familia | Pesos balanceados | Sobremuestreo | Línea base |
|---|---|---|---|
| NBM (mejor α) | 0.97107 ± 0.00080 (α = 0.1) | 0.97231 ± 0.00045 (α = 0.1) | 0.96990 ± 0.00093 (α = 0.1) |
| RLS (mejor C) | 0.97697 ± 0.00141 (C = 1) | 0.97748 ± 0.00124 (C = 1) | 0.97792 ± 0.00115 (C = 1) |

En test (una sola evaluación, sin ajuste):

| Medida | RLS | NBM |
|---|---|---|
| Macro-F1, test completo | 0.9795 | 0.9708 |
| Macro-F1 sin duplicados de train (11,040) | 0.9826 | 0.9726 |
| Macro-F1 sin el subconjunto ambiguo (13,348) | 0.9997 | 0.9952 |
| Recall de Reconocimientos (RLS) | 0.9819; sin ambiguos, 0.9986 | |

Los umbrales CPM01 (≥ 0.85) y CPM02 (≥ 0.70) se cumplen.

**Cómo leer esta evidencia con honestidad:**

- **El corpus anterior era trivialmente separable** (`b810bad2…`). Con él, RLS obtuvo 1.0000 y una matriz diagonal, y NBM 0.9997, lo que contradice el TE §6.1.4. El generador se corrigió con plantillas de contexto cruzado y con plantillas compartidas entre dos categorías, y la batería lo verifica con K1–K3. La regla fue fijar el criterio primero, sin umbral sobre el modelo final.
- **El error vive en el subconjunto ambiguo.** Sin los 1,252 comentarios ambiguos de test, RLS queda en 0.9997. La matriz de confusión del test completo muestra justo las confusiones que el TE §5.2.2 anticipa (Quejas → Falta 3.0 %; Cobros ↔ Facturación 2.2 % y 1.7 %; Reconocimientos → Quejas 1.7 %). Fuera de ese subconjunto, el corpus sigue siendo casi separable, porque el test temporal contiene las mismas plantillas que train.
- **Plantillas no vistas** (`docs/reports/plantillas_no_vistas.md`; solo reporte, no participa en la selección): con el 20 % de las plantillas fuera del entrenamiento, RLS promedia 0.8350 de Macro-F1 (rango 0.7637–0.8855) y NBM 0.7592. Es la cota práctica de la limitación L-1.
- **RLS > NBM** en CV (0.978 contra 0.972, con desviaciones ≈ 0.001) y en test (0.9795 contra 0.9708). Ahora la diferencia se puede medir.
- **Estrategias de desbalance casi empatadas.** Las tres de RLS quedan a menos de una desviación estándar entre sí (0.97697 a 0.97792). La línea base gana por la media, y la elección no es estadísticamente firme. En NBM, el sobremuestreo supera a las otras dos por más de una desviación.
- **Fuga no supervisada en CV.** El vocabulario y el IDF se calculan una vez con todo train, como pide el TE, y no por pliegue; las estadísticas de documento de los pliegues de validación entran en el IDF. Puede inflar ligeramente el Macro-F1 de CV; no afecta a test, que nunca participa.
- Se descartaron 793 lexemas del catálogo con df = 0 en train (aparecen solo en test); el vocabulario del modelo tiene 5,403.

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

### 6. Ambigüedad controlada del corpus

**Problema.** Con el corpus `b810bad2…` el modelo obtuvo Macro-F1 = 1.0000 y una matriz de confusión diagonal, incluso sin duplicados. El TE §6.1.4 exige que el generador "evite que el corpus sea trivialmente separable, lo que inflaría artificialmente las métricas", y que la batería verifique un "solapamiento léxico inter-categoría acotado". Con un 1.0 no se puede demostrar la comparación NBM contra RLS ni las tres estrategias de desbalance (§5.1.3), porque empataban. El diagnóstico (`docs/reports/generator_validation_report.md`, "Iteraciones del generador") fue que cada plantilla base tenía al menos un lexema exclusivo de su categoría y que ningún texto base aparecía con dos etiquetas.

**Criterio registrado antes de cambiar el generador** (commit en rojo `4a36e53`, en el que el corpus anterior falla las tres compuertas; umbrales fijos desde entonces, `PLAN.md` §3):

| Compuerta | Qué mide | Rango | Corpus `b810bad2…` | Corpus `37f285ff…` |
|---|---|---|---|---|
| K1 | Jaccard de los 100 términos más frecuentes por par de categorías, **después de la stoplist v1** | [0.15, 0.80] | 0.111 | 0.156–0.299 |
| K2 | Colisión de la oración base: fracción de comentarios cuya bolsa base (texto limpio, con stoplist) aparece en ≥ 2 categorías | global [0.05, 0.10]; ≥ 0.03 por categoría; ≤ 0.05 en Reconocimientos | 0.0000 | 0.0895 (Rec. 0.0381) |
| K3 | Sonda Naive Bayes (α = 1, 5 pliegues estratificados, semilla 10433) | Macro-F1 en [0.80, 0.99] | 1.0000 | 0.9726 |

Las plantillas se ajustaron mirando solo la batería, en 4 iteraciones registradas. El modelo se entrenó **una sola vez**, después de que la batería pasó, y no hay ningún umbral sobre el modelo final.

**Qué es y qué no es.** La ambigüedad controlada son 10 oraciones base compartidas entre los pares de categorías que el TE identifica como próximos: Falta↔Quejas (3), Cobros↔Facturación (3), Quejas↔Reconocimientos (2), Quejas↔Cobros (1) y Quejas↔Facturación (1). Cada una es plausible con las dos lecturas (`docs/reports/plantillas_nuevas.md`). Se suman 30 plantillas de contexto cruzado, en las que aparecen palabras de otra categoría (como el sarcasmo o el problema ya resuelto), inspiradas en patrones reales de quejas de telecomunicaciones en México (`docs/data/grok_examples.txt`), pero redactadas como plantillas parametrizadas.

**No es ruido de etiqueta**, porque:
- la etiqueta nunca contradice el texto;
- la determina la plantilla, de forma determinista;
- la tasa está declarada (8.95 %, ids en `data/corpus_ambiguos.json`, fuera de la BD para respetar el Anexo I).

El TE la anticipa en tres lugares: en las causas primarias del diagnóstico ("inconsistencias inter-evaluador", "ambigüedad de los textos"), en §5.2.2 ("confusiones sistemáticas entre categorías semánticamente próximas (por ejemplo, Quejas y Falta de Servicio)") y en el guion de §6.4.3 ("un reconocimiento con tono ambiguo"). El ruido de etiqueta (invertir etiquetas al azar) se descartó porque el TE solo contempla ruido léxico y ortográfico.

**Efecto medido:**
- Macro-F1 de 0.9795 en test, 0.9997 sin el subconjunto ambiguo y 0.8350 con plantillas no vistas (solo reporte, `docs/reports/plantillas_no_vistas.md`).
- RLS supera a NBM en validación cruzada (0.978 contra 0.972), así que la comparación del TE §6.4.2 produce evidencia.
- Las tres estrategias de desbalance siguen a menos de una desviación entre sí; la elección la decide el desempate y se declara así.

**Limitación.** Fuera del subconjunto ambiguo el corpus sigue siendo casi separable, porque el corte temporal no excluye plantillas: el test contiene las mismas plantillas que train. La cifra realista de generalización es la de plantillas no vistas; la cifra del TE (test temporal) se presenta junto a ella, nunca sola.

## Alternativas descartadas

- **`TfidfVectorizer` de scikit-learn con su propio tokenizador:** habría dos definiciones de token (la de Python y la del analizador léxico propio de la consigna b) y la paridad CPA06 sería inalcanzable. Su fórmula de IDF suavizado es la misma que se usa aquí, pero el cálculo se hace sobre la salida del tokenizador Go.
- **Servir con Python (`pickle` + `predict_proba`):** contradice el sobre de despliegue (`.exe` único, sin runtime) y la regla de que Python es solo de desarrollo.
- **Exportar solo el ganador:** se conserva también el mejor de la otra familia como `modelo_alternativo.json`, porque el plan B ante una pregunta o un umbral incumplido es mostrarlo (`PLAN.md` §6).
- **Escalar a otros clasificadores (SVM, árboles, redes):** fuera del TE §5.1.3, que compara NBM y RLS.
- **Calcular el IDF por pliegue:** más estricto en CV, pero el TE pide IDF "precalculados durante el entrenamiento" con un solo vocabulario; se declara la limitación en lugar de cambiar el diseño.

## Consecuencias

- `make train` reproduce los artefactos versionados y `make train-sha` comprueba que dos entrenamientos independientes producen los mismos bytes. `make verify` no requiere Python: las pruebas Go leen los JSON versionados.
- Si cambian el corpus, la stoplist o el stemmer, hay que regenerar `resultados.db` y `df_train.json` (`make clean-results sweep`) y volver a entrenar (`make train`); la prueba de paridad y la de `gold5` detectan un desajuste.
- Las cifras de CPM01/CPM02 deben presentarse a los sinodales con su salvedad: miden el corpus sintético con plantillas vistas. Hay que mostrar las tres vistas: test, sin ambiguos y plantillas no vistas.
