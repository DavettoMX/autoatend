# Plantillas nuevas del generador: revisión del usuario

Corpus `37f285ff…`, rama `fix/corpus-separabilidad`. Estas plantillas se escribieron para cumplir el criterio de separabilidad K1–K3, que se registró en `PLAN.md` §3 antes de cambiar el generador. Todas están en `internal/corpus/templates/templates.go`.

- **Contexto cruzado:** palabras de otra categoría usadas en el sentido de esta (sarcasmo, problema resuelto, mención lateral). Se escribieron a partir de `docs/data/grok_examples.txt` (número de línea = número de ejemplo) y de los 3 ejemplos del usuario. Se descartó el #8 por incoherente y se reescribió el #93 ("me ajustaron la factura mal").
- **Compartidas:** el mismo texto es plantilla de dos categorías (ambigüedad controlada). La etiqueta es una lectura que el texto admite, no la única; nunca contradice el texto. Son el 8.95 % del corpus (6,533 comentarios, ids en `data/corpus_ambiguos.json`).

Marca `x` en *OK* si la plantilla es plausible y pertenece a su categoría. En las compartidas, marca `x` si las **dos** lecturas son razonables. Si no, anota el motivo.

## Constancia

| Campo | Valor |
|---|---|
| Fecha | |
| Revisor | |
| Plantillas revisadas | /47 |
| Observaciones | |

## 1. Contexto cruzado (6 por categoría)

| # | Categoría | Plantilla | Origen | OK | Observación |
|---|---|---|---|---|---|
| 1 | Falta de Servicio | Llevo {DURACION} sin {SERVICIO} y me siguen cobrando el recibo completo, puro cuento. | Grok #1 | [ ] | |
| 2 | Falta de Servicio | Es increíble la falta de atención, mi negocio sigue sin internet desde hace {DURACION}. | usuario, Grok #2 | [ ] | |
| 3 | Falta de Servicio | Gracias, TelArg: la factura llega puntual, pero el internet no funciona desde {MES}. | Grok #3 | [ ] | |
| 4 | Falta de Servicio | La factura dice servicio activo, pero el módem no da internet desde hace {DURACION}. | Grok #12 | [ ] | |
| 5 | Falta de Servicio | El asesor fue amable, pero sigo sin {SERVICIO} después de {VECES} visitas del técnico. | Grok #9 | [ ] | |
| 6 | Falta de Servicio | Qué buena atención, ya van {DIAS} días sin {SERVICIO_DATOS} y nadie me da fecha de solución. | Grok #11 | [ ] | |
| 7 | Cobros Indebidos | Cancelé mi servicio porque nunca tuve señal y me siguen haciendo cobros de {MONTO} a mi tarjeta. | usuario, Grok #31 | [ ] | |
| 8 | Cobros Indebidos | Me hicieron un cobro doble de {MONTO} después de reportar que no tenía internet. | Grok #32 | [ ] | |
| 9 | Cobros Indebidos | Gracias por el cargo no reconocido de {MONTO}, ni idea de qué es. | Grok #33 | [ ] | |
| 10 | Cobros Indebidos | Ya di de baja la línea y me siguen cobrando el plan {PLAN}, la atención no sirve. | Grok #34 | [ ] | |
| 11 | Cobros Indebidos | El asesor me dijo que no me cobrarían más y en {MES} me volvió a salir el cargo de {MONTO}. | Grok #37 | [ ] | |
| 12 | Cobros Indebidos | Increíble, cancelé por mala atención y me siguen cobrando {MONTO_CHICO} cada mes en la tarjeta. | Grok #38 | [ ] | |
| 13 | Problemas de Facturación | Gracias por nada, la factura de {MES} llegó otra vez con el plan equivocado. | Grok #51 | [ ] | |
| 14 | Problemas de Facturación | Me cobran el monto completo, pero la factura dice un descuento que nunca aplicaron. | Grok #53 | [ ] | |
| 15 | Problemas de Facturación | El CFDI de {MES} no llega y el cobro sí se hizo. | Grok #54 | [ ] | |
| 16 | Problemas de Facturación | Increíble, la factura dice que no he pagado todo y no aparece mi abono de {MONTO}. | Grok #58 | [ ] | |
| 17 | Problemas de Facturación | La factura dice servicio activo, pero en {ZONA} no tengo internet desde hace días. | Grok #66 | [ ] | |
| 18 | Problemas de Facturación | El recibo no trae el descuento que me prometieron en atención a clientes por la falla de {SERVICIO}. | Grok #70 | [ ] | |
| 19 | Quejas | Me quitan el servicio por un peso de adeudo, cuando el problema es la atención. | usuario, Grok #71 | [ ] | |
| 20 | Quejas | El asesor de la sucursal de {CIUDAD} fue grosero y me cortó la línea por pagar un día tarde. | Grok #74 | [ ] | |
| 21 | Quejas | Publicidad engañosa: prometen cobertura en {CIUDAD} y no hay señal. | Grok #76 | [ ] | |
| 22 | Quejas | Mala atención, me reagendan al técnico {VECES} veces y nadie da la cara. | Grok #77 | [ ] | |
| 23 | Quejas | Para cancelar te piden mil cosas y mientras tanto te siguen cobrando el plan {PLAN}. | Grok #79 | [ ] | |
| 24 | Quejas | Me cortan por un adeudo mínimo cuando ellos no mandan bien la factura. | Grok #82 | [ ] | |
| 25 | Reconocimientos | Por fin me devolvieron el cobro doble de {MONTO}, el asesor fue muy amable. | Grok #91 | [ ] | |
| 26 | Reconocimientos | Gracias, TelArg: después de {DURACION} sin internet lo arreglaron en cuanto lo reporté. | Grok #92 | [ ] | |
| 27 | Reconocimientos | Excelente atención, me corrigieron la factura de {MES} y el recibo ya cuadra con mi pago. | Grok #93 (reescrito) | [ ] | |
| 28 | Reconocimientos | Por fin resolvieron la falta de señal en {ZONA}, el técnico llegó puntual. | Grok #94 | [ ] | |
| 29 | Reconocimientos | Gracias, me quitaron de la cuenta el cargo de {MONTO_CHICO} que no reconocía y me devolvieron mi dinero. | Grok #95 | [ ] | |
| 30 | Reconocimientos | Buena atención {CANAL}, me aclararon el monto de mi estado de cuenta en un momento. | Grok #96 | [ ] | |

## 2. Compartidas (10 textos, cada uno en dos categorías)

| # | Plantilla | Categorías | Lectura A | Lectura B | OK | Observación |
|---|---|---|---|---|---|---|
| 31 | Llevo {DURACION} esperando una solución y nadie me da respuesta. | Falta ↔ Quejas | Falta: sigue sin servicio y no lo reparan | Quejas: nadie de atención le responde | [ ] | |
| 32 | Ya reporté la falla {VECES} veces y nadie de soporte me resuelve. | Falta ↔ Quejas | Falta: la falla técnica sigue | Quejas: soporte no da seguimiento | [ ] | |
| 33 | Otra vez me dejaron sin servicio y en atención a clientes solo me dan largas (ref. {FOLIO}). | Falta ↔ Quejas | Falta: corte del servicio | Quejas: le dan largas en atención | [ ] | |
| 34 | Hay un cargo de {MONTO} en mi recibo de {MES} que no me cuadra. | Cobros ↔ Facturación | Cobros: un cargo que no reconoce | Facturación: el recibo no cuadra | [ ] | |
| 35 | Necesito que revisen lo que me llegó en {MES}, no coincide con lo que pago. | Cobros ↔ Facturación | Cobros: le cobraron de más | Facturación: el documento no coincide | [ ] | |
| 36 | El monto de este mes no es el que me dijeron {CANAL}, quiero una aclaración. | Cobros ↔ Facturación | Cobros: monto cobrado indebido | Facturación: monto facturado distinto | [ ] | |
| 37 | Fui a aclarar un cobro de {MONTO} y en {AREA} me dejaron esperando sin darme respuesta. | Cobros ↔ Quejas | Cobros: el cobro sigue sin aclararse | Quejas: lo dejaron esperando | [ ] | |
| 38 | Pedí {VECES} veces que corrijan mi factura y en {AREA} nadie me hace caso. | Facturación ↔ Quejas | Facturación: la factura sigue mal | Quejas: nadie le hace caso | [ ] | |
| 39 | Al final sí me resolvieron en {AREA}, aunque tardaron bastante. | Quejas ↔ Reconocimientos | Quejas: reclama la tardanza | Reconocimientos: reconoce que se resolvió (tono ambiguo, TE §6.4.3) | [ ] | |
| 40 | Me atendieron {CANAL} y quedó resuelto, aunque tuve que insistir {VECES} veces. | Quejas ↔ Reconocimientos | Quejas: tuvo que insistir | Reconocimientos: quedó resuelto | [ ] | |

## 3. Plantillas base nuevas de Reconocimientos (6)

Hacían falta para que la colisión de Reconocimientos quedara en ≤ 0.05 con 2 plantillas compartidas, y para que su vocabulario se cruzara con Cobros y Facturación (K1 ≥ 0.15). Las razones están en "Iteraciones del generador", en el reporte del generador.

| # | Plantilla | OK | Observación |
|---|---|---|---|
| 41 | Pagar mi recibo en la app es muy fácil, todo queda registrado al momento. | [ ] | |
| 42 | Me encanta que la factura llegue cada mes {ENTREGA_ELEC} sin falta. | [ ] | |
| 43 | El estado de cuenta ahora viene muy claro, gracias por detallar cada cargo. | [ ] | |
| 44 | Gracias por avisarme antes de la fecha de corte, así no se me pasa el pago. | [ ] | |
| 45 | En {AREA} me explicaron con paciencia el cobro de mi plan {PLAN}, quedé muy conforme. | [ ] | |
| 46 | La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. | [ ] | |

## 4. Plantillas existentes modificadas: marco "(ref. {FOLIO})"

Antes el marco "(ref. {FOLIO})" solo existía en Falta de Servicio, y el rasgo `ref` pesaba 7.1 en el modelo anterior. Ahora estas 7 plantillas usan "(ref. {FOLIO})" en lugar de "({FOLIO})"; no cambia nada más.

| # | Categoría | Plantilla | OK |
|---|---|---|---|
| 47 | Cobros (2), Facturación (2), Quejas (2), Reconocimientos (1) | p. ej. "Me descontaron {MONTO} de mi saldo sin haber usado nada (ref. {FOLIO})." | [ ] |
