# F1 con plantillas no vistas (solo reporte)

Corpus `37f285ff…`. **No es una compuerta ni un umbral:** no participa en la selección del modelo ni en la evaluación de `evaluacion_modelo.md`, y no se usó para ajustar el generador (se midió después del único entrenamiento).

**Por qué.** El corte temporal del TE (§5.2.2) no excluye plantillas: cada plantilla aparece cientos de veces en train y en test. Por eso el test mide sobre todo si el modelo reconoce plantillas ya vistas. Este análisis mide cuánto generaliza a plantillas que **no** vio, como evidencia para la limitación L-1 (validez externa del corpus sintético).

**Método.**
- Las 216 plantillas distintas (226 entradas; las 10 compartidas cuentan una vez) se reparten en 5 pliegues según su posición en el orden alfabético, módulo 5. Una plantilla compartida queda fuera en sus dos categorías a la vez.
- En cada pliegue se entrena con los comentarios de las plantillas de los otros 4 y se evalúa con los de las plantillas excluidas.
- Características: lexemas de `cmd/tokenizar` y TF-IDF (IDF suavizado, L2) ajustado con el entrenamiento del pliegue.
- Modelos: regresión logística softmax con C = 1 sin ajuste de desbalance (la configuración de `modelo.json`) y Naive Bayes multinomial con α = 0.1. Al Naive Bayes no se le aplica el sobremuestreo de `modelo_alternativo.json`.
- La plantilla de cada comentario se recupera con un generador instrumentado. Es una copia de `internal/corpus` que registra el índice de plantilla y reproduce los 73,000 textos de `corpus.db` sin ninguna diferencia. Este generador y el script de evaluación son herramientas de análisis fuera del repositorio.

| Pliegue | Plantillas excluidas | n evaluado | RLS Macro-F1 | RLS recall Reconocimientos | NBM Macro-F1 | NBM recall Reconocimientos |
|---|---|---|---|---|---|---|
| 0 | 44 | 16,281 | 0.8499 | 0.9390 | 0.7659 | 0.9255 |
| 1 | 43 | 14,305 | 0.7637 | 0.5583 | 0.7680 | 0.5154 |
| 2 | 43 | 13,770 | 0.8381 | 0.7366 | 0.7741 | 0.7965 |
| 3 | 43 | 14,155 | 0.8378 | 0.6793 | 0.7297 | 0.7182 |
| 4 | 43 | 14,489 | 0.8855 | 0.7250 | 0.7585 | 0.7052 |
| **Media** | | | **0.8350** | 0.7277 | **0.7592** | 0.7322 |

**Lectura.** Con plantillas vistas (test temporal), la regresión logística obtiene 0.9795 de Macro-F1. Con plantillas no vistas, el promedio es 0.8350, con un rango de 0.7637 a 0.8855, y el recall de Reconocimientos baja hasta 0.5583. La diferencia es la parte de la cifra de test que se debe a haber visto la plantilla. Es la cota práctica de L-1: un comentario real con una redacción que no está en ninguna plantilla se parece más a este escenario que al test temporal.
