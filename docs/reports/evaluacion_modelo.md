# Evaluación del modelo (Bloque C, consigna e)

Reporte generado por `make train` (`ml/autoatend_ml`). Es determinista: dos corridas producen el mismo archivo. Las cifras salen de `data/corpus.db` (SHA `37f285ff…`), `data/df_train.json` y `data/corpus_ambiguos.json`; ninguna se ajustó sobre test. El generador de este corpus cumple el criterio de separabilidad K1–K3, registrado en `PLAN.md` §3 antes de cambiar las plantillas; el modelo se entrenó una sola vez, después de que la batería pasó.

## 1. Datos y vocabulario

- Corte temporal (`db.SplitCutoff`, misma consulta SQL): `2026-04-06`. Train = `fecha < corte`: **58,400** comentarios; test = `fecha >= corte`: **14,600**.
- Catálogo de `resultados.db` (corpus completo): 6,196 lexemas. Vocabulario del modelo (df > 0 en train): **5,403** lexemas.
- Lexemas con **df = 0 en train, excluidos del modelo: 793** (aparecen solo en test; el TE §6.4.2 calcula el vocabulario y el IDF con train).
- Verificación cruzada: la frecuencia de documento recalculada en Python con la salida de `cmd/tokenizar` coincide con `df_train.json` (0 discrepancias, 0 lexemas de más).
- Comentarios de test con vector vacío (ningún lexema del vocabulario): 0.
- Duplicados exactos: **3,560 de 14,600** comentarios de test (24.4 %) tienen un texto idéntico en train; el subconjunto sin duplicados tiene **11,040**.
- Subconjunto ambiguo (K2): 6,533 comentarios del corpus cuya oración base aparece, como bolsa de tokens, en dos categorías (plantillas compartidas: ambigüedad controlada, sin ruido de etiquetas). En test son **1,252**; el test sin ellos tiene **13,348**.

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
| Naive Bayes multinomial | pesos balanceados | alpha = 0.1 | 0.97107 ± 0.00080 |
| Naive Bayes multinomial | sobremuestreo de Reconocimientos | alpha = 0.1 | 0.97231 ± 0.00045 (mejor de la otra familia) |
| Naive Bayes multinomial | línea base | alpha = 0.1 | 0.96990 ± 0.00093 |
| Regresión logística softmax | pesos balanceados | C = 1 | 0.97697 ± 0.00141 |
| Regresión logística softmax | sobremuestreo de Reconocimientos | C = 1 | 0.97748 ± 0.00124 |
| Regresión logística softmax | línea base | C = 1 | 0.97792 ± 0.00115 **(ganador)** |

Rejilla completa:

| Familia | Estrategia | Hiperparámetros | Macro-F1 CV (media ± desv.) | Convergió |
|---|---|---|---|---|
| nbm | pesos_balanceados | alpha = 0.01 | 0.97076 ± 0.00080 | n/a |
| nbm | pesos_balanceados | alpha = 0.1 | 0.97107 ± 0.00080 | n/a |
| nbm | pesos_balanceados | alpha = 0.5 | 0.97076 ± 0.00085 | n/a |
| nbm | pesos_balanceados | alpha = 1 | 0.97043 ± 0.00090 | n/a |
| nbm | sobremuestreo | alpha = 0.01 | 0.97203 ± 0.00059 | n/a |
| nbm | sobremuestreo | alpha = 0.1 | 0.97231 ± 0.00045 | n/a |
| nbm | sobremuestreo | alpha = 0.5 | 0.97187 ± 0.00029 | n/a |
| nbm | sobremuestreo | alpha = 1 | 0.97135 ± 0.00057 | n/a |
| nbm | linea_base | alpha = 0.01 | 0.96974 ± 0.00106 | n/a |
| nbm | linea_base | alpha = 0.1 | 0.96990 ± 0.00093 | n/a |
| nbm | linea_base | alpha = 0.5 | 0.96871 ± 0.00037 | n/a |
| nbm | linea_base | alpha = 1 | 0.96734 ± 0.00038 | n/a |
| rls | pesos_balanceados | C = 0.1 | 0.97497 ± 0.00128 | sí |
| rls | pesos_balanceados | C = 1 | 0.97697 ± 0.00141 | sí |
| rls | pesos_balanceados | C = 10 | 0.97663 ± 0.00099 | sí |
| rls | sobremuestreo | C = 0.1 | 0.97556 ± 0.00120 | sí |
| rls | sobremuestreo | C = 1 | 0.97748 ± 0.00124 | sí |
| rls | sobremuestreo | C = 10 | 0.97678 ± 0.00106 | sí |
| rls | linea_base | C = 0.1 | 0.97499 ± 0.00159 | sí |
| rls | linea_base | C = 1 | 0.97792 ± 0.00115 | sí |
| rls | linea_base | C = 10 | 0.97787 ± 0.00134 | sí |

Las 21 combinaciones de la rejilla quedan entre 0.96734 y 0.97792 de Macro-F1 promedio.

**Empate o diferencia menor a una desviación estándar** entre el ganador y 2 celda(s) más: la elección entre ellas no es estadísticamente firme y la decide el criterio de desempate declarado en el método (menor desviación, luego primera de la rejilla).

## 4. Ganador

**Regresión logística softmax**, estrategia **línea base**, C = 1; Macro-F1 de CV 0.97792 ± 0.00115. Reajustado con todo train y exportado como `data/modelo.json`.
El mejor de la otra familia (Naive Bayes multinomial, sobremuestreo de Reconocimientos, alpha = 0.1; CV 0.97231 ± 0.00045) se exporta como `data/modelo_alternativo.json`.

## 5. Evaluación final en test (una sola vez, sin ajustar nada)

| Criterio | Umbral | Test completo | Test sin duplicados de train | Test sin subconjunto ambiguo |
|---|---|---|---|---|
| CPM01 Macro-F1 | ≥ 0.85 | 0.9795 (CUMPLE) | 0.9826 (CUMPLE) | 0.9997 (CUMPLE) |
| CPM02 recall de Reconocimientos | ≥ 0.7 | 0.9819 (CUMPLE) | 0.9877 (CUMPLE) | 0.9986 (CUMPLE) |
| n | | 14,600 | 11,040 | 13,348 |

#### Modelo ganador, test completo (Regresión logística softmax)

n = 14,600; exactitud = 0.9793; **Macro-F1 = 0.9795**; **recall de Reconocimientos = 0.9819**.

| Categoría | Precisión | Recall | F1 | Soporte |
|---|---|---|---|---|
| Falta de Servicio | 0.9790 | 0.9972 | 0.9880 | 4,308 |
| Cobros Indebidos | 0.9767 | 0.9740 | 0.9754 | 2,928 |
| Problemas de Facturación | 0.9720 | 0.9807 | 0.9764 | 2,905 |
| Quejas | 0.9803 | 0.9558 | 0.9679 | 2,964 |
| Reconocimientos | 0.9980 | 0.9819 | 0.9899 | 1,495 |

Matriz de confusión normalizada por fila (cada fila suma 1):

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 0.9972 | 0.0000 | 0.0000 | 0.0028 | 0.0000 |
| Cobros Indebidos | 0.0000 | 0.9740 | 0.0215 | 0.0044 | 0.0000 |
| Problemas de Facturación | 0.0000 | 0.0169 | 0.9807 | 0.0024 | 0.0000 |
| Quejas | 0.0304 | 0.0064 | 0.0064 | 0.9558 | 0.0010 |
| Reconocimientos | 0.0013 | 0.0000 | 0.0000 | 0.0167 | 0.9819 |

Conteos:

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 4,296 | 0 | 0 | 12 | 0 |
| Cobros Indebidos | 0 | 2,852 | 63 | 13 | 0 |
| Problemas de Facturación | 0 | 49 | 2,849 | 7 | 0 |
| Quejas | 90 | 19 | 19 | 2,833 | 3 |
| Reconocimientos | 2 | 0 | 0 | 25 | 1,468 |

#### Modelo ganador, test sin duplicados exactos de train

n = 11,040; exactitud = 0.9825; **Macro-F1 = 0.9826**; **recall de Reconocimientos = 0.9877**.

| Categoría | Precisión | Recall | F1 | Soporte |
|---|---|---|---|---|
| Falta de Servicio | 0.9843 | 0.9976 | 0.9909 | 3,278 |
| Cobros Indebidos | 0.9806 | 0.9792 | 0.9799 | 2,687 |
| Problemas de Facturación | 0.9763 | 0.9806 | 0.9785 | 2,062 |
| Quejas | 0.9796 | 0.9611 | 0.9703 | 1,952 |
| Reconocimientos | 0.9990 | 0.9877 | 0.9934 | 1,061 |

Matriz de confusión normalizada por fila (cada fila suma 1):

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 0.9976 | 0.0000 | 0.0000 | 0.0024 | 0.0000 |
| Cobros Indebidos | 0.0000 | 0.9792 | 0.0160 | 0.0048 | 0.0000 |
| Problemas de Facturación | 0.0000 | 0.0160 | 0.9806 | 0.0034 | 0.0000 |
| Quejas | 0.0256 | 0.0097 | 0.0031 | 0.9611 | 0.0005 |
| Reconocimientos | 0.0019 | 0.0000 | 0.0000 | 0.0104 | 0.9877 |

Conteos:

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 3,270 | 0 | 0 | 8 | 0 |
| Cobros Indebidos | 0 | 2,631 | 43 | 13 | 0 |
| Problemas de Facturación | 0 | 33 | 2,022 | 7 | 0 |
| Quejas | 50 | 19 | 6 | 1,876 | 1 |
| Reconocimientos | 2 | 0 | 0 | 11 | 1,048 |

#### Modelo ganador, test sin el subconjunto ambiguo

n = 13,348; exactitud = 0.9998; **Macro-F1 = 0.9997**; **recall de Reconocimientos = 0.9986**.

| Categoría | Precisión | Recall | F1 | Soporte |
|---|---|---|---|---|
| Falta de Servicio | 0.9993 | 1.0000 | 0.9996 | 4,031 |
| Cobros Indebidos | 1.0000 | 1.0000 | 1.0000 | 2,663 |
| Problemas de Facturación | 1.0000 | 1.0000 | 1.0000 | 2,686 |
| Quejas | 1.0000 | 0.9996 | 0.9998 | 2,524 |
| Reconocimientos | 1.0000 | 0.9986 | 0.9993 | 1,444 |

Matriz de confusión normalizada por fila (cada fila suma 1):

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 1.0000 | 0.0000 | 0.0000 | 0.0000 | 0.0000 |
| Cobros Indebidos | 0.0000 | 1.0000 | 0.0000 | 0.0000 | 0.0000 |
| Problemas de Facturación | 0.0000 | 0.0000 | 1.0000 | 0.0000 | 0.0000 |
| Quejas | 0.0004 | 0.0000 | 0.0000 | 0.9996 | 0.0000 |
| Reconocimientos | 0.0014 | 0.0000 | 0.0000 | 0.0000 | 0.9986 |

Conteos:

| Real \ Predicho | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---|---|---|---|---|
| Falta de Servicio | 4,031 | 0 | 0 | 0 | 0 |
| Cobros Indebidos | 0 | 2,663 | 0 | 0 | 0 |
| Problemas de Facturación | 0 | 0 | 2,686 | 0 | 0 |
| Quejas | 1 | 0 | 0 | 2,523 | 0 |
| Reconocimientos | 2 | 0 | 0 | 0 | 1,442 |

### Modelo alternativo (mejor de la otra familia)

| Modelo | Test: Macro-F1 | Test: recall Reconocimientos | Sin duplicados: Macro-F1 | Sin duplicados: recall Reconocimientos | Sin ambiguos: Macro-F1 | Sin ambiguos: recall Reconocimientos |
|---|---|---|---|---|---|---|
| `modelo.json` (rls, linea_base) | 0.9795 | 0.9819 | 0.9826 | 0.9877 | 0.9997 | 0.9986 |
| `modelo_alternativo.json` (nbm, sobremuestreo) | 0.9708 | 0.9746 | 0.9726 | 0.9793 | 0.9952 | 0.9910 |

### Nota de cordura (corpus sintético)

El Naive Bayes multinomial alcanza Macro-F1 = 0.9708 en test, por debajo de 0.99.

Además, el 24.4 % del test repite textos de train; por eso las métricas se reportan también sin ellos (sección 5).

**F1 con plantillas no vistas (solo reporte, no es compuerta ni umbral).** El corte temporal no excluye plantillas: todas aparecen en train y en test. Para medir cuánto depende el modelo de haber visto cada plantilla, `docs/reports/plantillas_no_vistas.md` reentrena con el 80 % de las plantillas de cada categoría y evalúa con el 20 % restante, en 5 pliegues. Es un análisis aparte, que no participa en la selección ni en la evaluación de esta sección.

## 6. Equivalencia del escalador lineal con scikit-learn

| Modelo | Máx. \|p_lineal − predict_proba\| sobre todo test | Acuerdo de clase |
|---|---|---|
| `modelo.json` (rls) | 5.00e-16 | 1.000000 |
| `modelo_alternativo.json` (nbm) | 3.11e-15 | 1.000000 |

`softmax(b + W·x)` con los pesos exportados reproduce `predict_proba` de scikit-learn (Naive Bayes: `softmax` de la verosimilitud conjunta = posterior).

## 7. Paridad Python ↔ Go (CPA06)

Medida por la prueba Go `TestCPA06_ParidadPythonGo` (`go test -run TestCPA06 ./internal/ml`) sobre `testdata/paridad_1000.json`: 1,000 comentarios de test elegidos con la semilla 10433. Python vectoriza con `cmd/tokenizar`; Go vuelve a tokenizar y vectorizar el texto con `internal/ml`.

| Medida | Resultado |
|---|---|
| Comentarios comparados | 1,000 |
| Soportes distintos (índices) | 0 |
| L1 entre vectores TF-IDF (Σ\|Δ\|, igualdad de bits) | 0 |
| Clase predicha idéntica | 1,000 de 1,000 |
| Máx. \|Δprob\| (Go vs Python) | 1.110e-16 |
| Resultado de la prueba | PASS |

Las probabilidades pasan por `exp`: Python usa la libm del sistema y Go `math.Exp`, que pueden diferir en el último bit, por eso se comparan con tolerancia 1e-12 y se reporta el máximo. Los vectores, en cambio, se comparan bit a bit.

## 8. Comentarios de la demostración

`data/gold5.json`: un comentario de test por categoría, sin duplicado en train, fuera del subconjunto ambiguo (nunca de una plantilla compartida), longitud legible, clasificado correctamente con confianza ≥ 0.80 (elegido al azar con la semilla 10433 entre los elegibles).

| Categoría | id | Confianza | Texto |
|---|---|---|---|
| Falta de Servicio | 72686 | 0.957 | No tengo señal en el centro por las mañanas, ya pedí que lo revisen (ref. #33172). |
| Cobros Indebidos | 71403 | 0.974 | Me hicieron un cobro doble de $349.00 despues de reportar que no tenia internet. |
| Problemas de Facturación | 61169 | 0.993 | Mi paago de $449.00 si se hizo pero la fact todavía marca vencido. Espero que lo resuelvan pronto. |
| Quejas | 69477 | 0.965 | Me cambiaron las condiciones de mi plan Nocturno 79 sin pedirme autorizacion, ya no es lo que contrate. |
| Reconocimientos | 70117 | 0.975 | Hola. Gracias a Daniela por la paciencia y la amabilidad al explicarme mi plan Navega 249. |

Candidatos para la demostración oral (no están en `gold5.json`):

- **Caso claro, Falta de Servicio** (id 68425): «Avisé de la falla por teléfono con folio 55024 y todavía no me restablecen el servicio. Quedo al pendiente.». Real: Falta de Servicio; predicción: Falta de Servicio con 0.992; segunda: Problemas de Facturación con 0.004. Lexemas: avisé, falla, folio, me, no, pendiente, quedo, restablecen, servicio, teléfono, todavía.
- **Caso frontera, Reconocimientos** (id 60639): «Me atendieron x correo y quedo resuelto, aunque tuve que insistir 5 vexes. Gracias.». Real: Reconocimientos; predicción: Reconocimientos con 0.550; segunda: Quejas con 0.391. Lexemas: atendieron, correo, gracias, insistir, me, quedo, resuelto, tuve, vexes.

El caso frontera es un Reconocimiento del subconjunto ambiguo (una plantilla compartida con Quejas, el «reconocimiento con tono ambiguo» del guion del TE §6.4.3): entre los del test, el de menor margen entre las dos clases más probables, acertado si alguno lo está. El caso claro sale del test sin duplicados ni ambiguos.

## 9. Reproducción

```bash
make train        # compila cmd/tokenizar, entrena, escribe modelo.json, modelo_alternativo.json, gold5.json, testdata/paridad_1000.json y este reporte
make train-sha    # dos entrenamientos independientes producen el mismo SHA-256
```
