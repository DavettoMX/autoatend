# Guía de revisión manual: L-1 y muestra de 250 comentarios

El TE compromete dos revisiones que hace el sustentante:
- **L-1** (§6.3.10): etiquetado manual de homógrafos retenidos por la stoplist.
- **Muestra** (§6.1.4): detección manual sobre la muestra del corpus.

Los archivos son borradores generados de forma automática. **El juicio es tuyo:** no copies la etiqueta de nadie, y anota la duda cuando la tengas. Tiempo estimado: entre 1 y 1.5 h en total.

## Checklist

- [ ] **L-1**: llenar la columna `clase` de las 200 filas de `l1_borrador.csv`.
- [ ] **L-1**: calcular la tasa de sub-filtrado (fórmula abajo) y anotarla en `PLAN.md` §5.2.
- [ ] **L-1**: guardar el archivo terminado como `l1_muestra_200.csv` (no sobrescribas el borrador).
- [ ] **Muestra**: marcar *OK* u *Observación* en los 250 comentarios de `muestra_250_borrador.md`.
- [ ] **Muestra**: llenar la tabla "Constancia de revisión" (fecha, revisor, total revisado, observaciones generales).
- [ ] **Muestra**: guardar el archivo como `muestra_250_revision.md`.
- [ ] Hacer el commit de los dos archivos en una rama `docs/revision-manual`.

---

## 1. L-1: homógrafos de la stoplist (`l1_borrador.csv`)

**Qué mide.** La stoplist **no filtra** `para`, `como` ni `entre` porque a veces son palabras con significado (verbos). L-1 cuenta cuántas veces actúan en realidad como palabra funcional, es decir, cuántas veces **debieron filtrarse y no se filtraron**. Esa es la tasa de sub-filtrado.

**Columnas.**
- `contexto`: la palabra entre `[[ ]]` con sus vecinas.
- `clase`: es la única que llenas.

**Qué escribir en `clase`.** Una de tres etiquetas:

| Etiqueta | Cuándo | Ejemplos |
|---|---|---|
| `cerrada` | Funciona como preposición, conjunción o adverbio comparativo o relativo. Es lo que la stoplist debería filtrar | "llamé **para** reclamar", "es **como** no tener señal", "**entre** semana", "**como** les dije" |
| `lexica` | Es un verbo u otra palabra con significado propio. Es correcto retenerla | "no **para** de fallar" (verbo *parar*), "yo **como** en casa" (verbo *comer*), "cuando **entre** la llamada" (verbo *entrar*) |
| `duda` | El contexto no basta o el ruido lo impide | "x **para** q…" muy roto |

En este corpus casi todas las filas serán `cerrada`. Ese resultado es esperado y sigue siendo válido: es justo lo que la medición debe mostrar.

**Cálculo.**

```
sub-filtrado = filas "cerrada" / (filas "cerrada" + filas "lexica")
```

Las filas `duda` se reportan aparte y no entran en el cálculo.

Resultado esperado: una cifra cercana al 100 %. Se interpreta así: "*la retención de homógrafos introduce como cota superior X % de tokens funcionales no filtrados*". Es la cota del TE §6.3.10.

**Nota para la réplica.** El TE lista bajo, para, sobre, entre y como, pero **bajo y sobre no aparecen ninguna vez en el corpus**. Por eso la muestra se reparte en 70 de `para`, 70 de `como` y 60 de `entre`. Dilo tal cual: la muestra se tomó de lo que el corpus realmente contiene.

**Para contar en una terminal** (desde `proyect/`):

```bash
cut -d, -f7 docs/reports/l1_muestra_200.csv | sort | uniq -c
```

Si alguna fila tiene comas dentro de `contexto`, cuenta mejor con una hoja de cálculo (LibreOffice: Datos → Tabla dinámica sobre la columna `clase`).

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
| `categoría ambigua` | Un reconocimiento que también se queja |
| `contradicción` | "llevo 2 días… llevo varias semanas" en el mismo texto |
| `inverosímil` | Combinación de datos que un cliente no escribiría |
| `repetitivo` | Idéntico a otro de la muestra (anótalo, no es grave: el reporte mide un 16 % de duplicados) |

**Al terminar**, llena la constancia: fecha, revisor (tu nombre), `250/250` y un resumen como "*N observaciones: X ambiguas, Y repetidas; ningún hueco de plantilla*".

**Para la réplica.** Los comentarios con observación no se corrigen en el corpus, porque eso cambiaría el SHA y todo lo que depende de él. Se reportan como hallazgo de validez del corpus sintético, ligado a la limitación de validez externa del TE §7.3.
