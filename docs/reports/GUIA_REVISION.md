# Guía de revisión manual: L-1, muestra de 250 comentarios y corpus de oro

El TE compromete tres revisiones que hace el sustentante:
- **L-1** (§6.3.10): etiquetado manual de homógrafos retenidos por la stoplist.
- **Muestra** (§6.1.4): detección manual sobre la muestra del corpus.
- **Corpus de oro** (§6.3.5c): 100 comentarios con los lexemas esperados, anotados por ti.

Los archivos son borradores generados de forma automática. **El juicio es tuyo:** no copies la etiqueta de nadie, y anota la duda cuando la tengas. Tiempo estimado: entre 1 y 1.5 h en total para L-1 y la muestra, más entre 40 y 60 min para el corpus de oro.

## Checklist

- [ ] **L-1**: llenar la columna `clase` de las filas `PENDIENTE` de `l1_contextos.csv` (16 contextos distintos que cubren las 200 filas).
- [ ] **L-1**: pedir la regeneración de `l1_muestra_200.csv` a partir de `l1_contextos.csv` (las 200 filas quedan llenas y consistentes).
- [ ] **L-1**: calcular la tasa de sub-filtrado (comando abajo) y anotarla en `PLAN.md` §5.2.
- [ ] **Muestra**: marcar *OK* u *Observación* en los 250 comentarios de `muestra_250_borrador.md`.
- [ ] **Muestra**: llenar la tabla "Constancia de revisión" (fecha, revisor, total revisado, observaciones generales).
- [ ] **Muestra**: guardar el archivo como `muestra_250_revision.md`.
- [ ] **Corpus de oro**: revisar las 100 filas de `corpus_oro.csv` (sección 3): `x` en `revisado` si la fila está bien; si no, la lista correcta en `notas` y la causa en `limitacion`.
- [ ] Hacer el commit de los archivos revisados en una rama `docs/revision-manual`.

---

## 1. L-1: homógrafos de la stoplist (`l1_contextos.csv`)

**Qué mide.** La stoplist **no filtra** `para`, `como` ni `entre` porque a veces son palabras con significado (verbos). L-1 cuenta cuántas veces actúan en realidad como palabra funcional, es decir, cuántas veces **debieron filtrarse y no se filtraron**. Esa es la tasa de sub-filtrado.

**Por qué 16 contextos y no 200 filas.** La muestra tiene 200 apariciones (`l1_borrador.csv`, partición de test), pero como el corpus es sintético salen de solo **16 frases de plantilla** distintas. Por ejemplo, las 60 de `entre` son todas "me transfirieron … veces entre áreas". Una frase idéntica siempre lleva la misma etiqueta, así que se etiqueta **una vez por contexto** y se pondera por su `frecuencia`. Da la misma tasa que etiquetar las 200 filas, sin inconsistencias.

**Columnas de `l1_contextos.csv`.**
- `ejemplo`: una aparición real, con la palabra entre `[[ ]]`.
- `frecuencia`: cuántas de las 200 filas tienen ese contexto.
- `clase`: es la única que llenas.
- `nota`: dice si la etiqueta ya viene de tu revisión o si está `PENDIENTE`.

**Qué escribir en `clase`.** Una de tres etiquetas. El criterio es si la palabra cumple una de las **4 clases que el caso manda ignorar** (artículo, preposición, conjunción, interjección):

| Etiqueta | Cuándo | Ejemplos |
|---|---|---|
| `cerrada` | Funciona como **preposición, conjunción, artículo o interjección**. Es lo que la stoplist debería filtrar | "llamé **para** reclamar", "para nada", "es **como** no tener señal" (conjunción comparativa), "**como** si lo hubiera pedido", "veces **entre** áreas" |
| `lexica` | Cualquier otra cosa: verbo, sustantivo o **adverbio**, incluidos el interrogativo y el relativo. Es correcto retenerla | "no sé **cómo** quitarlo" (interrogativo sin tilde por el ruido), "no **para** de fallar" (verbo *parar*), "yo **como** en casa" (verbo *comer*) |
| `duda` | El contexto no basta para decidir | |

Caso discutible: "atiendan **como** cliente" ("en calidad de" + sustantivo, función de preposición). Elige una etiqueta y aplícala igual a todo el contexto.

**Cálculo** (desde `proyect/`), ponderado por frecuencia:

```bash
awk -F, 'NR>1 { t[$6] += $5 } END { c = t["cerrada"]; l = t["lexica"]; printf "cerrada %d, lexica %d, duda %d -> sub-filtrado %.1f %%\n", c, l, t["duda"], 100*c/(c+l) }' docs/reports/l1_contextos.csv
```

Las filas `duda` se reportan aparte y no entran en el cálculo. Se interpreta así: "*la retención de homógrafos introduce como cota superior X % de tokens funcionales no filtrados*". Es la cota del TE §6.3.10.

**Notas para la réplica.**
- El TE lista bajo, para, sobre, entre y como, pero **bajo y sobre no aparecen ninguna vez en el corpus**. La muestra se reparte en 70 de `para`, 70 de `como` y 60 de `entre`, que es lo que el corpus realmente contiene.
- Que 200 apariciones se reduzcan a 16 contextos es un hallazgo de la poca diversidad del corpus sintético (validez externa, TE §7.3).
- El ruido que quita tildes convierte "cómo" en "como", el mismo efecto que ADR-003 documenta para "qué" → "que".

---

## 2. Muestra de 250 comentarios (`muestra_250_borrador.md`)

**Qué mide.** Si el corpus sintético es **plausible** (parece escrito por clientes reales de TelArg) y si cada comentario **pertenece a su categoría**. Es la "detección manual" del TE §6.1.4.

**Cómo revisar cada comentario.** Marca `x` en *OK* si se cumplen las tres:
1. **Plausible**: un cliente podría escribirlo. Las erratas, la falta de tildes, las abreviaturas `q` y `x` y las mayúsculas irregulares **son ruido intencional** del generador, así que no cuentan como defecto.
2. **Categoría correcta**: el tema principal coincide con la sección donde aparece (Falta de Servicio, Cobros Indebidos, Problemas de Facturación, Quejas, Reconocimientos).
3. **Sin defectos del generador**: no hay huecos de plantilla sin llenar, frases contradictorias ni texto cortado.

Si algo falla, deja *OK* vacío y escribe un motivo corto en *Observación*:

| Motivo | Ejemplo |
|---|---|
| `slot incoherente` | Dentro de **una misma oración** el hueco no concuerda: "el internet está tan lento… que es como no tener **servicio de llamadas**" |
| `temas mezclados` | Dos oraciones sobre problemas distintos sin conexión: "hay mantenimiento en Oaxaca…" + "pagué puntual y me suspendieron" |
| `contradicción` | Las oraciones se contradicen: "desde hace varias horas…" + "llevo dos semanas…" |
| `categoría ambigua` | Un reconocimiento que también se queja |
| `repetitivo` | Idéntico a otro de la muestra (anótalo, no es grave: el reporte mide un 20 % de duplicados) |

**Qué cambió respecto a la primera versión del corpus.** La primera revisión (corpus con SHA `ee3f7d55…`) encontró defectos reales: en el 40 % de los comentarios se pegaban dos plantillas de la misma categoría sin relación entre sí, en 298 de 588 comentarios con "mi plan cuesta X y me están cobrando Y" el cobro era menor que el precio, el precio no coincidía con el número del plan, algunos huecos chocaban dentro de la oración ("tan lento… como no tener servicio de llamadas") y las aperturas y los cierres no distinguían la categoría. El generador se corrigió y el corpus se regeneró (SHA en `PLAN.md` §3): ahora `compose` (`internal/corpus/generate.go`) arma cada comentario con un saludo opcional (compartido o de la categoría), **una sola plantilla base**, a veces un **complemento** de la misma categoría que amplía el mismo problema sin añadir otro, y una despedida opcional compatible con la categoría. Los planes tienen precio fijo y los montos de "cuesta/cobrando" salen de él. Esta muestra se saca del corpus nuevo.

**La columna *Auto*.** Un script revisó los 250 comentarios con verificaciones automáticas (coherencia de montos y precios, una sola plantilla base, apertura/cierre/complemento compatible con la categoría, huecos que concuerdan). `ok` no significa que el comentario sea plausible: solo que esas verificaciones no encontraron nada. Tu revisión sigue siendo la que cuenta, y el ruido intencional puede ocultar un caso al script.

**Cómo decidir rápido.**
- `slot incoherente` y `contradicción` → *OK* vacío.
- `temas mezclados` → *OK* con `x` si un cliente real podría quejarse de las dos cosas a la vez (es común listar varios problemas). Vacío si no tiene sentido juntarlas. Anota el motivo en los dos casos.
- No corrijas nada: el objetivo de esta revisión es **medir** cuántos hay. Ese número es el resultado.

**Al terminar**, llena la constancia: fecha, revisor (tu nombre), `250/250` y un resumen como "*N de 250 con observación: X slot incoherente, Y temas mezclados, Z contradicción; categoría correcta en W de 250*".

**Para la réplica.** El TE (§6.1.4) pide que la revisión manual detecte inconsistencias y que el corpus solo persista si pasan las validaciones. La primera revisión las detectó y se corrigió el generador (cambió el SHA y todo lo que depende de él, y quedó registrado en `PLAN.md` §3). Los comentarios con observación que aparezcan en esta segunda muestra tampoco se corrigen a mano: se reportan como hallazgo de validez del corpus sintético, ligado a la limitación de validez externa del TE §7.3 (plantillas finitas: el 20 % de los textos es duplicado exacto).

---

## 3. Corpus de oro (TE §6.3.5c)

**Qué es y para qué sirve.** Son 100 comentarios reales del corpus (de la partición de **test**) con los lexemas que el sistema produce para cada uno. El TE pide que **tú** los revises a mano, para que sirvan de **oráculo independiente**: una vez que tú confirmes que la lista de cada fila es la correcta, una prueba de integración compara lo que el analizador produce hoy contra tu lista. Si el analizador cambia y se rompe algo, esa prueba lo detecta. Es la evidencia de la consigna (b): que el pipeline da el resultado correcto sobre entradas conocidas.

**Qué hay en el archivo `corpus_oro.csv`.** Un borrador generado de forma automática con la semilla 10433: 30 de Falta de Servicio, 20 de Cobros Indebidos, 20 de Problemas de Facturación, 20 de Quejas y 10 de Reconocimientos, sin textos repetidos. Incluye frases compartidas entre categorías, frases de contexto cruzado, comentarios con ruido (erratas, sin tildes, `q`, `x`, mayúsculas), homógrafos (`para`, `como`, `entre`) y números o folios. **El juicio es tuyo:** la lista viene del propio sistema, así que puede traer sus mismos errores y tu trabajo es encontrarlos.

**Columnas.**
- `n`, `id`, `categoria`, `texto`: la fila y el comentario.
- `lexemas_esperados`: pares `forma→canónico` en el orden del texto, separados por `; `. La forma es como aparece escrita (ya en minúsculas) y el canónico es el lexema del catálogo al que el sistema la asigna. Ejemplo: `cobraron→cobro; doble→doble; factura→factura`.
- `revisado`, `limitacion`, `notas`: las llenas tú.

**Qué revisar en cada fila (unos 25 segundos por fila).**

1. **Qué palabras aparecen.** Debe estar **toda palabra con significado**. No debe estar ningún artículo, preposición, conjunción ni interjección (stoplist v1, ADR-003).
   - Se filtran: `el, la, de, en, con, por, a, sin, desde, hasta, y, o, ni, que, pero, porque, aunque, ah, ay`, y las abreviaturas de ruido `q` y `x`.
   - Se **conservan** a propósito los homógrafos `para`, `como`, `entre`, `bajo` y `sobre` (L-1), y también `se`, `no`, `este`, `eso`, `mi`, `su`, `me`.
   - Los **números no se catalogan**: `$299`, `3`, `48` y la parte numérica de un folio (`ref. folio 4985234`) no deben aparecer; sí aparece la palabra (`folio`, `ticket`, `ref`).
2. **Que cada forma vaya al lexeme correcto.** `cobraron→cobro` es correcto (el mismo verbo cobrar). Es **error** si el agrupador (stemmer) juntó dos palabras distintas, por ejemplo `estado→esto` (el sustantivo *estado* quedó junto con *este/esta/esto*) o `plano→plan`. Que el lexema canónico tenga otra terminación no es error por sí solo (`esperando→espero`): solo importa que sea la misma palabra.
3. **Marcar.** Si la fila está bien, escribe `x` en `revisado`. Si el sistema se equivoca, **no cambies la columna `lexemas_esperados`**: escribe en `notas` la lista correcta (o solo lo que cambia, por ejemplo `estado→estado`) y en `limitacion` una de estas causas:

| Valor en `limitacion` | Cuándo |
|---|---|
| `L-1` | Un homógrafo (`para`, `como`, `entre`) que aquí cumple función de preposición o conjunción y debió filtrarse (o al revés) |
| `L-2` | Morfología irregular: formas del mismo verbo o palabra que quedan en lexemas distintos (por ejemplo `tuve` y `tengo`) |
| `stemmer` | El stemmer juntó dos palabras distintas (`estado→esto`) |
| `sin plegado de acentos` | Dos escrituras de la misma palabra con y sin tilde quedan como lexemas distintos (`facturacion` y `facturación`) |
| `ruido/erratas` | La errata o la abreviatura impide reconocer la palabra (`ecsribo`, `diinero`) |
| `error` | Ninguna de las anteriores: un defecto no documentado |

Si una fila falla por varias causas, anota la principal y menciona las otras en `notas`. Una fila **con limitación documentada no es un fallo del sistema**: es una limitación conocida del diseño. Lo que importa es que quede registrada.

**Al terminar.** Cuenta cuántas filas quedaron con `x`, cuántas con cada limitación y cuántas con `error`, y avísame para escribir la prueba de integración del corpus de oro (4c) con tu lista ya validada. Si hay filas con limitación, la prueba las separa en un grupo aparte con lo que tú indicaste.

**Dónde mirar con más cuidado.** Las agrupaciones del stemmer (por ejemplo `este/esta/estado/estar→esto`, `seguido/seguir→seguimiento`, `plano→plan`, `calle→call`), las palabras de saludo o conectores que no están en la stoplist (`hola`, `si`) y los homógrafos.

**Para la réplica.** El corpus de oro es independiente en el sentido que importa: la lista final la decide una persona leyendo el texto, no el código. Los defectos que encuentres no se corrigen a mano: se registran con su limitación, porque la stoplist y el stemmer son una versión congelada (v1).
