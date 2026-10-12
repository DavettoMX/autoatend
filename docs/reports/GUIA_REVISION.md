# Guía de revisión manual: L-1 y muestra de 250 comentarios

El TE compromete dos revisiones que hace el sustentante:
- **L-1** (§6.3.10): etiquetado manual de homógrafos retenidos por la stoplist.
- **Muestra** (§6.1.4): detección manual sobre la muestra del corpus.

Los archivos son borradores generados de forma automática. **El juicio es tuyo:** no copies la etiqueta de nadie, y anota la duda cuando la tengas. Tiempo estimado: entre 1 y 1.5 h en total.

## Checklist

- [ ] **L-1**: llenar la columna `clase` de las filas `PENDIENTE` de `l1_contextos.csv` (16 contextos distintos que cubren las 200 filas).
- [ ] **L-1**: pedir la regeneración de `l1_muestra_200.csv` a partir de `l1_contextos.csv` (las 200 filas quedan llenas y consistentes).
- [ ] **L-1**: calcular la tasa de sub-filtrado (comando abajo) y anotarla en `PLAN.md` §5.2.
- [ ] **Muestra**: marcar *OK* u *Observación* en los 250 comentarios de `muestra_250_borrador.md`.
- [ ] **Muestra**: llenar la tabla "Constancia de revisión" (fecha, revisor, total revisado, observaciones generales).
- [ ] **Muestra**: guardar el archivo como `muestra_250_revision.md`.
- [ ] Hacer el commit de los dos archivos en una rama `docs/revision-manual`.

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
