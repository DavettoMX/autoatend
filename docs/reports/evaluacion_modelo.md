# Evaluación del modelo (Bloque C, consigna e)

Reporte generado por `make train` (`ml/autoatend_ml`). Es determinista: dos corridas producen el mismo archivo. Las cifras salen de `data/corpus.db` (SHA `b810bad2…`) y de `data/df_train.json`; ninguna se ajustó sobre test.

## 1. Datos y vocabulario

- Corte temporal (`db.SplitCutoff`, misma consulta SQL): `2026-04-06`. Train = `fecha < corte`: **58,400** comentarios; test = `fecha >= corte`: **14,600**.
- Catálogo de `resultados.db` (corpus completo): 6,289 lexemas. Vocabulario del modelo (df > 0 en train): **5,443** lexemas.
- Lexemas con **df = 0 en train, excluidos del modelo: 846** (aparecen solo en test; el TE §6.4.2 calcula el vocabulario y el IDF con train).
- Verificación cruzada: la frecuencia de documento recalculada en Python con la salida de `cmd/tokenizar` coincide con `df_train.json` (0 discrepancias, 0 lexemas de más).
- Comentarios de test con vector vacío (ningún lexema del vocabulario): 0.
- Duplicados exactos: **3,716 de 14,600** comentarios de test (25.5 %) tienen un texto idéntico en train; el subconjunto sin duplicados tiene **10,884**.

## 2. Método

- **Tokenización**: `cmd/tokenizar` (Go) invocado una vez como subproceso (TE §6.3.6). Cada forma → `lexer.Stem` → lexema canónico del catálogo; números y stoplist fuera; raíces fuera del catálogo se descartan.
- **Características** (TE §5.1.3, §6.4.2): TF = conteo crudo; IDF suavizado `idf = ln((1 + N) / (1 + df)) + 1` con N y df de train (`df_train.json`); vector normalizado en L2. Los productos y la suma de cuadrados se hacen en Python puro, en orden ascendente de índice, para que Go reproduzca los bits.
- **Familias**: Naive Bayes multinomial (suavizado de Laplace/Lidstone, α ∈ {0.01, 0.1, 0.5, 1}) y regresión logística softmax multinomial con penalización L2 (C ∈ {0.1, 1, 10}). La regresión usa `solver='saga'` (gradiente promedio estocástico, de la familia del descenso de gradiente estocástico), `max_iter = 300`, `tol = 0.0001`.
- **Desbalance (30/20/20/20/10 %)**, tres estrategias: (a) pesos balanceados (`class_weight='balanced'`; en Naive Bayes, `sample_weight` de `compute_sample_weight('balanced')`); (b) sobremuestreo aleatorio de Reconocimientos hasta el tamaño mediano de clase (20 %), con semilla fija, aplicado solo a las particiones de entrenamiento de cada pliegue y al ajuste final, nunca a validación ni a test; (c) línea base sin ajuste.
- **Selección**: validación cruzada estratificada de 5 particiones (`StratifiedKFold(5, shuffle=True, random_state=10433)`) sobre train. Gana la celda con mayor Macro-F1 promedio; en un empate, la de menor desviación y luego la primera de la rejilla. `n_jobs = 1`, semillas fijas.
- **Limitación declarada**: el vocabulario y el IDF se calculan una vez con todo train (como pide el TE), no por pliegue. Es una fuga no supervisada de estadísticas de documento hacia los pliegues de validación; puede inflar ligeramente el Macro-F1 de CV y no afecta a test.
- **Modelo lineal único** para las dos familias: `z_k = b_k + Σ_i w_{k,i}·x_i`, softmax con resta del máximo. Naive Bayes: `w = feature_log_prob_`, `b = class_log_prior_`; regresión: `coef_`, `intercept_`.

## 3. Validación cruzada (train)

Mejor hiperparámetro de cada celda; Macro-F1 promedio ± desviación estándar muestral de las 5 particiones.

| Familia | Estrategia | Mejores hiperparámetros | Macro-F1 CV (media ± desv.) |
|---|---|---|---|
| Naive Bayes multinomial | pesos balanceados | alpha = 0.5 | 0.99987 ± 0.00013 |
| Naive Bayes multinomial | sobremuestreo de Reconocimientos | alpha = 0.5 | 0.99989 ± 0.00014 (mejor de la otra familia) |
| Naive Bayes multinomial | línea base | alpha = 0.1 | 0.99964 ± 0.00023 |
| Regresión logística softmax | pesos balanceados | C = 10 | 1.00000 ± 0.00000 **(ganador)** |
| Regresión logística softmax | sobremuestreo de Reconocimientos | C = 10 | 1.00000 ± 0.00000 |
| Regresión logística softmax | línea base | C = 10 | 1.00000 ± 0.00000 |

Rejilla completa:

| Familia | Estrategia | Hiperparámetros | Macro-F1 CV (media ± desv.) | Convergió |
|---|---|---|---|---|
| nbm | pesos_balanceados | alpha = 0.01 | 0.99970 ± 0.00022 | n/a |
| nbm | pesos_balanceados | alpha = 0.1 | 0.99976 ± 0.00018 | n/a |
| nbm | pesos_balanceados | alpha = 0.5 | 0.99987 ± 0.00013 | n/a |
| nbm | pesos_balanceados | alpha = 1 | 0.99984 ± 0.00015 | n/a |
| nbm | sobremuestreo | alpha = 0.01 | 0.99971 ± 0.00021 | n/a |
| nbm | sobremuestreo | alpha = 0.1 | 0.99977 ± 0.00018 | n/a |
| nbm | sobremuestreo | alpha = 0.5 | 0.99989 ± 0.00014 | n/a |
| nbm | sobremuestreo | alpha = 1 | 0.99984 ± 0.00020 | n/a |
| nbm | linea_base | alpha = 0.01 | 0.99958 ± 0.00026 | n/a |
| nbm | linea_base | alpha = 0.1 | 0.99964 ± 0.00023 | n/a |
| nbm | linea_base | alpha = 0.5 | 0.99960 ± 0.00030 | n/a |
| nbm | linea_base | alpha = 1 | 0.99942 ± 0.00027 | n/a |
| rls | pesos_balanceados | C = 0.1 | 0.99986 ± 0.00010 | sí |
| rls | pesos_balanceados | C = 1 | 0.99998 ± 0.00004 | sí |
| rls | pesos_balanceados | C = 10 | 1.00000 ± 0.00000 | sí |
| rls | sobremuestreo | C = 0.1 | 0.99986 ± 0.00010 | sí |
| rls | sobremuestreo | C = 1 | 0.99998 ± 0.00004 | sí |
| rls | sobremuestreo | C = 10 | 1.00000 ± 0.00000 | sí |
| rls | linea_base | C = 0.1 | 0.99994 ± 0.00005 | sí |
| rls | linea_base | C = 1 | 0.99998 ± 0.00004 | sí |
| rls | linea_base | C = 10 | 1.00000 ± 0.00000 | sí |

Las 21 combinaciones de la rejilla quedan entre 0.99942 y 1.00000 de Macro-F1 promedio: el corpus es tan separable que la estrategia de desbalance y el hiperparámetro cambian el resultado en la cuarta o quinta cifra decimal.

**Empate o diferencia menor a una desviación estándar** entre el ganador y 2 celda(s) más: la elección entre ellas no es estadísticamente firme y la decide el criterio de desempate declarado en el método (menor desviación, luego primera de la rejilla).

El ganador queda en el borde de la rejilla (C = 10 es el extremo superior de la rejilla); la rejilla es la fijada por el diseño y no se amplió.

## 4. Ganador

**Regresión logística softmax**, estrategia **pesos balanceados**, C = 10; Macro-F1 de CV 1.00000 ± 0.00000. Reajustado con todo train y exportado como `data/modelo.json`.
El mejor de la otra familia (Naive Bayes multinomial, sobremuestreo de Reconocimientos, alpha = 0.5; CV 0.99989 ± 0.00014) se exporta como `data/modelo_alternativo.json`.

## 5. Evaluación final en test (una sola vez, sin ajustar nada)

| Criterio | Umbral | Test completo | Test sin duplicados de train |
|---|---|---|---|
| CPM01 Macro-F1 | ≥ 0.85 | 1.0000 (CUMPLE) | 1.0000 (CUMPLE) |
| CPM02 recall de Reconocimientos | ≥ 0.7 | 1.0000 (CUMPLE) | 1.0000 (CUMPLE) |
| n | | 14,600 | 10,884 |

#### Modelo ganador, test completo (Regresión logística softmax)

n = 14,600; exactitud = 1.0000; **Macro-F1 = 1.0000**; **recall de Reconocimientos = 1.0000**.

| Categoría | Precisión | Recall | F1 | Soporte |
|---|---|---|---|---|
| Falta de Servicio | 1.0000 | 1.0000 | 1.0000 | 4,308 |
| Cobros Indebidos | 1.0000 | 1.0000 | 1.0000 | 2,928 |
| Problemas de Facturación | 1.0000 | 1.0000 | 1.0000 | 2,905 |
| Quejas | 1.0000 | 1.0000 | 1.0000 | 2,964 |
| Reconocimientos | 1.0000 | 1.0000 | 1.0000 | 1,495 |

Matriz de confusión normalizada por fila (cada fila suma 1):

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 1.0000 | 0.0000 | 0.0000 | 0.0000 | 0.0000 |
| Cobros Indebidos | 0.0000 | 1.0000 | 0.0000 | 0.0000 | 0.0000 |
| Problemas de Facturación | 0.0000 | 0.0000 | 1.0000 | 0.0000 | 0.0000 |
| Quejas | 0.0000 | 0.0000 | 0.0000 | 1.0000 | 0.0000 |
| Reconocimientos | 0.0000 | 0.0000 | 0.0000 | 0.0000 | 1.0000 |

Conteos:

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 4,308 | 0 | 0 | 0 | 0 |
| Cobros Indebidos | 0 | 2,928 | 0 | 0 | 0 |
| Problemas de Facturación | 0 | 0 | 2,905 | 0 | 0 |
| Quejas | 0 | 0 | 0 | 2,964 | 0 |
| Reconocimientos | 0 | 0 | 0 | 0 | 1,495 |

#### Modelo ganador, test sin duplicados exactos de train

n = 10,884; exactitud = 1.0000; **Macro-F1 = 1.0000**; **recall de Reconocimientos = 1.0000**.

| Categoría | Precisión | Recall | F1 | Soporte |
|---|---|---|---|---|
| Falta de Servicio | 1.0000 | 1.0000 | 1.0000 | 3,307 |
| Cobros Indebidos | 1.0000 | 1.0000 | 1.0000 | 2,710 |
| Problemas de Facturación | 1.0000 | 1.0000 | 1.0000 | 1,988 |
| Quejas | 1.0000 | 1.0000 | 1.0000 | 1,880 |
| Reconocimientos | 1.0000 | 1.0000 | 1.0000 | 999 |

Matriz de confusión normalizada por fila (cada fila suma 1):

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 1.0000 | 0.0000 | 0.0000 | 0.0000 | 0.0000 |
| Cobros Indebidos | 0.0000 | 1.0000 | 0.0000 | 0.0000 | 0.0000 |
| Problemas de Facturación | 0.0000 | 0.0000 | 1.0000 | 0.0000 | 0.0000 |
| Quejas | 0.0000 | 0.0000 | 0.0000 | 1.0000 | 0.0000 |
| Reconocimientos | 0.0000 | 0.0000 | 0.0000 | 0.0000 | 1.0000 |

Conteos:

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 3,307 | 0 | 0 | 0 | 0 |
| Cobros Indebidos | 0 | 2,710 | 0 | 0 | 0 |
| Problemas de Facturación | 0 | 0 | 1,988 | 0 | 0 |
| Quejas | 0 | 0 | 0 | 1,880 | 0 |
| Reconocimientos | 0 | 0 | 0 | 0 | 999 |

### Modelo alternativo (mejor de la otra familia)

| Modelo | Test: Macro-F1 | Test: recall Reconocimientos | Sin duplicados: Macro-F1 | Sin duplicados: recall Reconocimientos |
|---|---|---|---|---|
| `modelo.json` (rls, pesos_balanceados) | 1.0000 | 1.0000 | 1.0000 | 1.0000 |
| `modelo_alternativo.json` (nbm, sobremuestreo) | 0.9997 | 1.0000 | 0.9996 | 1.0000 |

### Nota de cordura (corpus sintético)

**El Naive Bayes multinomial alcanza Macro-F1 = 0.9997 (> 0.99) en test.** Es una señal de que el corpus sintético, generado a partir de plantillas por categoría, es fácil de separar. Los umbrales del TE se cumplen con holgura, pero esa holgura describe al corpus y no la dificultad del problema real de TelArg.

Además, el 25.5 % del test repite textos de train; por eso las métricas se reportan también sin ellos (sección 5).

## 6. Equivalencia del escalador lineal con scikit-learn

| Modelo | Máx. \|p_lineal − predict_proba\| sobre todo test | Acuerdo de clase |
|---|---|---|
| `modelo.json` (rls) | 4.44e-16 | 1.000000 |
| `modelo_alternativo.json` (nbm) | 2.44e-15 | 1.000000 |

`softmax(b + W·x)` con los pesos exportados reproduce `predict_proba` de scikit-learn (Naive Bayes: `softmax` de la verosimilitud conjunta = posterior).

## 7. Paridad Python ↔ Go (CPA06)

Medida por la prueba Go `TestCPA06_ParidadPythonGo` (`go test -run TestCPA06 ./internal/ml`) sobre `testdata/paridad_1000.json`: 1,000 comentarios de test elegidos con la semilla 10433. Python vectoriza con `cmd/tokenizar`; Go vuelve a tokenizar y vectorizar el texto con `internal/ml`.

| Medida | Resultado |
|---|---|
| Comentarios comparados | 1,000 |
| Soportes distintos (índices) | 0 |
| L1 entre vectores TF-IDF (Σ\|Δ\|, igualdad de bits) | 0 |
| Clase predicha idéntica | 1,000 de 1,000 |
| Máx. \|Δprob\| (Go vs Python) | 1.388e-17 |
| Resultado de la prueba | PASS |

Las probabilidades pasan por `exp`: Python usa la libm del sistema y Go `math.Exp`, que pueden diferir en el último bit, por eso se comparan con tolerancia 1e-12 y se reporta el máximo. Los vectores, en cambio, se comparan bit a bit.

## 8. Comentarios de la demostración

`data/gold5.json`: un comentario de test por categoría, sin duplicado en train, longitud legible, clasificado correctamente con confianza ≥ 0.80 (elegido al azar con la semilla 10433 entre los elegibles).

| Categoría | id | Confianza | Texto |
|---|---|---|---|
| Falta de Servicio | 72630 | 1.000 | Tengo un problema con mi servicio. Ayer hubo una caída general de la red en Puebla y hoy todavía sigue fallando todo. |
| Cobros Indebidos | 71479 | 1.000 | Cobro indebido: me descontaron $149.00 dos veces el mismo día por lo mismo. Es dinero que me hace falta. |
| Problemas de Facturación | 60892 | 0.999 | La fecha de vencimiento de mi recibo no es la q me dieron, por eso se me juntan los periodos |
| Quejas | 69682 | 0.999 | Esperé cuarenta minutos en soporte técnico y nadie me contestó. No es justo que uno tenga que insistir tanto. |
| Reconocimientos | 70347 | 0.994 | Mira nomás, qué buen servicio, en menos de media hora me solucionaron todo. |

Candidatos para la demostración oral (no están en `gold5.json`):

- **Caso claro, Falta de Servicio** (id 68369): «Hice mi recarga de $349.00 hace quince días y sigo sin internet.». Real: Falta de Servicio; predicción: Falta de Servicio con 1.000; segunda: Cobros Indebidos con 0.000. Lexemas: días, hace, hice, internet, mi, quince, recarga, sigo.
- **Caso frontera, Reconocimientos** (id 68620): «Buenos dias. Felciidades por la cobertura en Chihuahua, hasta en carretera tengo señal.». Real: Reconocimientos; predicción: Reconocimientos con 0.851; segunda: Falta de Servicio con 0.147. Lexemas: buenas, carretera, chihuahua, cobertura, dias, señal, tengo.

El caso frontera es el Reconocimientos acertado de menor confianza del test sin duplicados; en un corpus fácil puede seguir teniendo una confianza alta.

## 9. Reproducción

```bash
make train        # compila cmd/tokenizar, entrena, escribe modelo.json, modelo_alternativo.json, gold5.json, testdata/paridad_1000.json y este reporte
make train-sha    # dos entrenamientos independientes producen el mismo SHA-256
```
