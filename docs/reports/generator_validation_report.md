# Generator validation report — corpus.db

Produced by `cmd/gencorpus`. Seed **10433**, 73000 comments, 100 per day from 2024-08-30 (2024-08-30 to 2026-08-29 in the file).
The corpus is persisted only if every check of the battery passes (verdict: **PASS**).

## Class distribution

| Category | Target | Generated | Persisted | Share |
|---|---:|---:|---:|---:|
| Falta de Servicio | 21900 | 21900 | 21900 | 30.00% |
| Cobros Indebidos | 14600 | 14600 | 14600 | 20.00% |
| Problemas de Facturación | 14600 | 14600 | 14600 | 20.00% |
| Quejas | 14600 | 14600 | 14600 | 20.00% |
| Reconocimientos | 7300 | 7300 | 7300 | 10.00% |
| **Total** | 73000 | 73000 | 73000 | 100.00% |

## Validation battery

| Check | Result | Detail |
|---|---|---|
| category counts | PASS | got [21900 14600 14600 14600 7300], target [21900 14600 14600 14600 7300], rows 73000/73000 |
| no empty texts | PASS | 0 empty texts |
| no unfilled slots | PASS | 0 texts with unfilled slots |
| dates | PASS | 730 days, 2024-08-30 to 2026-08-29, 100 per day |
| noise rate | PASS | realized 0.2979, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | K1: pairwise Jaccard of top-100 tokens after stoplist v1 in [0.156, 0.299], required within [0.15, 0.80] |
| base collision | PASS | K2: global 0.0895 (required [0.05, 0.10]), per category 0.0661 / 0.0928 / 0.0858 / 0.1507 / 0.0381 (each >= 0.03, Reconocimientos <= 0.05) |
| probe classifier | PASS | K3: Naive Bayes probe (alpha 1, 5 stratified folds) Macro-F1 0.9726, required within [0.80, 0.99] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 19.35% duplicate texts; text length 42 / 110.8 / 228 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.258 | 0.250 | 0.299 | 0.205 |
| Cobros Indebidos | 0.258 | 1.000 | 0.250 | 0.227 | 0.198 |
| Problemas de Facturación | 0.250 | 0.250 | 1.000 | 0.190 | 0.156 |
| Quejas | 0.299 | 0.227 | 0.190 | 1.000 | 0.212 |
| Reconocimientos | 0.205 | 0.198 | 0.156 | 0.212 | 1.000 |

## Generator parameters

| Parameter | Value |
|---|---:|
| Templates per category | 45 / 44 / 44 / 47 / 46 |
| Base / cross-context / shared templates | 36·6·3 / 34·6·4 / 34·6·4 / 34·6·7 / 38·6·2 |
| Complements per category | 13 / 13 / 12 / 12 / 14 |
| Shared openers / closers | 9 / 3 |
| Own openers per category | 3 / 3 / 3 / 3 / 4 |
| Own closers per category | 6 / 6 / 6 / 4 / 4 |
| Opener / complement / closer probability | 0.25 / 0.40 / 0.30 |
| Amount pool / surface formats | 24 / 6 |
| Folio prefixes / digits | 6 / 5–7 |

## Noise calibration

| Parameter | Value |
|---|---:|
| Comment selection rate | 0.30 |
| Tilde drop (per accented letter) | 0.80 |
| Abbreviation q/x/tel/fact (per word) | 0.50 |
| Keyboard typo (per word of 4+ letters) | 0.04 |
| Case jitter (per comment) | 0.25 |
| Irregular punctuation (per comment) | 0.30 |
| **Realized selection rate** | **0.2979** (tolerance ±0.015) |

## Persisted corpus verification

Measured by re-opening the finished file read-only.

| Category | Rows | Train | Test | max(train.fecha) | min(test.fecha) |
|---|---:|---:|---:|---|---|
| Falta de Servicio | 21900 | 17592 | 4308 | 2026-04-05 | 2026-04-06 |
| Cobros Indebidos | 14600 | 11672 | 2928 | 2026-04-05 | 2026-04-06 |
| Problemas de Facturación | 14600 | 11695 | 2905 | 2026-04-05 | 2026-04-06 |
| Quejas | 14600 | 11636 | 2964 | 2026-04-05 | 2026-04-06 |
| Reconocimientos | 7300 | 5805 | 1495 | 2026-04-05 | 2026-04-06 |

Empty texts: 0. Split rule (global, temporal): the cut-off is the fecha at position 80% of the corpus in chronological order; train is fecha < 2026-04-06, test is fecha >= 2026-04-06, so every training comment precedes every test comment.

Test comments with an identical text in train: 3560 of 14600 (24.38%). A temporal split does not prevent exact repeats of a template draw.

corpus.db SHA-256: `37f285ff1da70f4bd5205c51e255a612d1a72ad45b70fc9fda35dbac1b3e6e63`

Ambiguous subset (K2 collision): 6533 comments (8.95%), ids in `data/corpus_ambiguos.json`.

## Iteraciones del generador

El criterio K1–K3 quedó registrado (`PLAN.md` §3, commit `test: pre-registered separability criterion K1-K3 … (red)`) antes de cambiar cualquier plantilla. Cada iteración se ajustó mirando **solo** la salida de la batería; ningún modelo se entrenó ni se evaluó hasta que la batería pasó.

| # | Cambio | K1 (Jaccard mín.–máx.) | K2 (global; Reconocimientos) | K3 (Macro-F1 sonda) | Veredicto |
|---|---|---|---|---|---|
| 0 | Plantillas originales (corpus `b810bad2…`) | 0.111–0.282 FAIL | 0.0000; 0.0000 FAIL | 1.0000 FAIL | FAIL |
| 1 | +6 plantillas de contexto cruzado por categoría (a partir de `docs/data/grok_examples.txt` y de los ejemplos del usuario), +10 plantillas compartidas en los pares del plan, marco "(ref. {FOLIO})" en las 5 categorías | 0.099–0.316 FAIL (Facturación–Reconocimientos) | 0.0925; 0.0500 PASS | 0.9710 PASS | FAIL |
| 2 | Vocabulario de facturación y cobros en 3 plantillas de contexto cruzado de Reconocimientos y "atención"/"todo" en 2 de Facturación | 0.111–0.316 FAIL | 0.0926; 0.0518 FAIL (Reconocimientos > 0.05) | 0.9702 PASS | FAIL |
| 3 | +6 plantillas base de Reconocimientos sobre pagos y facturación (diluyen las 2 compartidas de Reconocimientos y comparten vocabulario con Cobros y Facturación) | 0.156–0.307 PASS | 0.0895; 0.0381 PASS | 0.9726 PASS | PASS |
| 4 | Coherencia: la plantilla compartida Quejas↔Reconocimientos decía "tuve que llamar" tras un canal escrito; ahora "tuve que insistir" | 0.156–0.299 PASS | 0.0895; 0.0381 PASS | 0.9726 PASS | PASS (corpus persistido) |

Antes del registro, el análisis exploratorio de `/tmp/separabilidad/` (variantes F y FA) sí entrenó modelos de prueba sobre corpus alternativos; esos resultados motivaron el criterio, pero no se usaron para fijar sus umbrales ni para ajustar estas iteraciones.

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. El asesor fue amable, pero sigo sin señal después de 5 visitas del técnico. Espero una respuesta pronto.
2. Cambié de SIM y quedé incomunicado, sin servicio de llamadas desde hace quince días. Por favor, avísenme cualquier novedad.
3. hola. El asesor fue amable, pero sigo sin servicio de mensajes después de 3 visitas del técnico. Esto afecta a toda mi familia. Espero una respuesta pronto.
4. La neta, llevo más de 48 horas sin internet y sigo pagando mi plan Nocturno 79. Ya probé todo lo que me dijeron y sigue igual.
5. Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace cinco días.
6. Mi chip nuevo no activa, llevo varias horas esperando que me den servicio de llamadas.
7. Buenas tardes. La verdad, llevo cuatro días sin datos móviles y sigo pagando mi plan Familiar 549. Ya me había pasado antes este mes.
8. El internet de mi plan Premium 999 no funciona desde enero y no hay quien lo repare. Por favor, avísenme cualquier novedad.
9. Mi línea está suspendida sin motivo y no puedo usar nada.
10. Ya reporté la falla 5 veces y nadie de soporte me resuelve. Necesito que me den una solución definitiva.
11. La neta, el internet está tan lento en fines de semana que es como no tener internet. Quedo atento a su respuesta.
12. Buenos días. Gracias, TelArg: la factura llega puntual, pero el internet no funciona desde agosto.
13. Qué barbaridad, otra vez sin datos móviles, ya es la tercera vez este mes que se cae.
14. Hola. La portabilidad de mi número no se completó y estoy sin línea desde hace cuatro días. Ya me había pasado antes este mes.
15. En mi casa nunca hay señal, tengo que salir a la calle para poder llamar.
16. El asesor fue amable, pero sigo sin servicio de mensajes después de 4 visitas del técnico. Ya me está afectando en el trabajo. Gracias.
17. Ayer hubo una caída general de la red en Saltillo y hoy todavía sigue fallando todo. Gracias.
18. Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace más de 48 horas. Ya me había pasado antes este mes.
19. Es increible la falta de atencion, mi negocio sigue sin internet desde hace mas de 48 horas. Espero una respuesta pronto.
20. LES ESCRIBO POR UNA FALLA EN MI LÍNEA. NO PUEDO MANDAR NI RECIBIR MENSAJES DE TEXTO EN PUEBLA DESDE HACE UNA SEMANA. YA ME HABIA PASADO ANTES ESTE MES...
21. Es increíble la falta de atención, mi negocio sigue sin internet desde hace tres días. Ya perdí la cuenta de las veces que he insistido. Muchas gracias por su atención.
22. Mi chip nuevo no activa, llevo diez días esperando que me den datos móviles. Quiero que alguien me explique qué está pasando. Quedo atento a su respuesta.
23. Buenos días. Pagué mi plan Empresarial 899 puntual y aun así me suspendieron el servicio. Ya probé todo lo que me dijeron y sigue igual. Gracias.
24. En mi trabajo nunca hay señal, tengo que salir a la calle para poder llamar. Muchas gracias por su atención.
25. Les escribo desde Culiacán. Gracias, TelArg: la factura llega puntual, pero el internet no funciona desde octubre. Espero que lo atiendan pronto.
26. El internet de mi plan Ahorro 99 no funciona desde julio y no hay quien lo repare. Espero su pronta respuetsa.
27. O sea, el internet está tan lento toda la tarde que es como no tener internet. Ya me había pasado antes este mes.
28. Buenas noches. Mi tel dice sin servicio desde hace diez dias, ya reinicié todo y nada. Quiero que alguien me explique qué está pasando
29. Hice mi recarga de 499 pewos hace una semana y sigo sin internet...
30. Desde hace quince días estoy sin servicio de mensajes en mi trabajo y nadie me da una solución. Ya me había pasado antes este mes.
31. Tengo la línea a mi nombre. Qué buena atención, ya van 5 días sin datos móviles y nadie me da fecha de solución. Quedo atento a su respuesta.
32. Híjole, llevo tres días sin datos móviles y sigo pagando mi plan Empresarial 899.
33. No manches, otra vez sin datos móviles, ya es la tercera vez este mes que se cae. Sigo pagando completo y no obtengo lo que corresponde.
34. No me llegan los códigos de verificación por SMS desde hace varias horas, creo que hay una falla en la red. Gracias.
35. No tengo señal en mi casa en horario laboral, ya pedí que lo revisen (ref. caso 8205023). Ya perdí la cuenta de las veces que he insistido. Gracias de antemano.
36. Buenas noches. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir.
37. Soy cliente de TelArg desde hace años. En mi barrio no hay cobertura de madrugada, parece que la antena está caída. Nadie le da seguimiento a esto. Por favor, avísenme cualquier novedad.
38. Ya Cambie Mi Tel Y Sigo Sin señal, El Problema Es De La red. espero Que lo Atiendan prono. espero Su Pronta Respuesta.
39. dicen que hay mantenimiento en Aguascalientes pero ya llevamos una semana y seguimos sin servicio de llamadas.
40. No puedo mandar ni recibir mensajes de texto en Villahermosa desde hace quince días.
41. Buenas tardes. No puedo mandar ni recibir mensajes de texto en Puebla desde hace tres días.
42. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir. Por favor, avísenme cualquier novedad.
43. Cuando intento llamar solo escucho que el número no esta disponible, ya llevo diez dias así. Quiero que alguien se haga cargo!!!
44. Hola. Dicen que hay mantenimiento en Monterrey pero ya llevamos dos días y seguimos sin cobertura. Quedo atento a su respuesta.
45. Ya reporte la falla 2 veces y nadie de soporte me resuelve.
46. La neta, llevo dos dias sin señal y sigo pagando mi plan Social 299. Ya probe todo lo que me dijeron y sigue igual.
47. Qué buena atención, ya van 30 días sin datos móviles y nadie me da fecha de solución. Esto afecta a toda mi familia. Gracias.
48. El técnico nunca llegó a revisar la falla de internet, ya llevo más de 48 horas de espera.
49. Buenas tardes. No tengo datos móviles en mi oficina en las noches, ya pedí que lo revisen (ref. reporte 548391). Nadie le da seguimiento a esto. Saludos cordiales.
50. Hay pérdida de cobertura toda la tarde en toda la zona, y ya son varios días con el mismo problema.

### Cobros Indebidos

1. Me quitaron el saldo de mi recarga de $499 MXN sin explicación alguna.
2. Pedí la devolución de $169 MXN hace una semana y todavía no me abonan nada. Espero su pronta respuesta.
3. Gracias x el cargo no reconocido de $299.00, ni idea de q es. Ya revise mi cuenta y no hay explicacion alguna. Saludos cordiales.
4. Me cobraron $99 MXN x servicio premium q nunca contraté. Es dinero que me hace falta
5. Desde que cambié de plan a Empresarial 899 me cobran $799 de más, es un cobro indebido. Ya revisé mi cuenta y no hay explicación alguna.
6. Me activaron renta de equipo sin avisarme y ahora me cobran $229.00. Nunca me avisaron de este cargo.
7. Después de cancelar mi plan Premium 999 me cargaron $999.00 a mi tarjeta. Quiero una explicación clara de este cargo.
8. Activaron suscripción de musica sin pedirmelo y me cobraron $129.00 en el estado de mi tarjeta. No pienso pagar un cargo q no corresponde.
9. Cargos que no reconozco en mayo: $129 por recarga automática. Exijo que los eliminen.
10. YA RECLAMÉ EL COBRO DE MEMBRESIA DE JUEGOS EN LA SUCURSAL Y NADIE ME DEVUELVE MI DINERO (#37895). ES DUNERO Q ME HACE FALTA. QUEDO ATENTO A SU RESPUESTA.
11. Les escribo desde Hermosillo. Me siguen cobrando la línea que di de baja, ya llevo 8 meses pagando de más. Me urge que lo corrijan cuanto antes. Por favor, avísenme cualquier novedad.
12. El monto de este mes no es el q me dijeron x telefono, quiero una aclaración.
13. Chale, aparece un cargo de $399.00 por suscripción de música y yo no autoricé nada. Quedo al pendiente.
14. Tengo la línea a mi nombre. Cargos que no reconozco en junio: 149 pesos por paquete de datos extra. Exijo que los eliminen. Me urge que lo corrijan cuanto antes. Muchas gracias por su atención.
15. Fui A aclarar un Cobro de 99 pesos Y en la Sucursal me Dejaron Esperando Sin darme respuesta. quiero Q Me devuelvan mi Dinero. quedo al pendiente.
16. Después de cancelar mi plan Ahorro 99 me cargaron 99 pesos a mi tarjeta. Ya pedí una aclaración y sigo sin respuesta. Gracias.
17. Mi plan Joven 149 cuesta $149 MXN y me están cobrando $169, no es lo q acordamos. Quiero una explicación clara de este cargo.
18. Les escribo desde Monterrey. El monto de este mes no es el que me dijeron en la app quiero una aclaracion. Quiero que me devuelvan mi dinero.
19. increible, cancele por mala atencion y me siguen cobrando $349.50 cada mes en la tarjeta. ya insisti varias veces y nadie me resuelve.
20. Ya reclame el cobro de contenido de terceros x telefono y nadie me devuelve mi dinero (#2692734).
21. El asesor me dijo que no me cobrarían más y en noviembre me volvió a salir el cargo de MXN 599. Me urge que lo corrijan cuanto antes. Gracias de antemano.
22. Mi plan Familiar 549 cuesta MXN 549 y me están cobrando $999 MXN no es lo que acordamos. Quiero una explicacion clara de este cargo.
23. Fui a aclarar un cobro de $99 y en el área de cobranza me dejaron esperando sin darme respuesta.
24. Buenas noches. Desconozco el cobro de $129.00 del día 5 de mayo, quiero que lo reviertan. Es muy frustrante que pase esto con mi dinero.
25. Desconozco el cobro de $349 del dia 10 de agosto, quiero que lo reviertan. Muchas gracias por su atencion.
26. Quiero q me reembolsen los $1,499 MXN que me cobraron de más (folio 3408217). Exijo q lo revisen cuanto antes.
27. Pedí la devolución de $49.00 hace tres días y todavía no me abonan nada.
28. Quiero que me reembolsen los $699.50 que me cobraron de más (#90347). Gracias de antemano.
29. Les escribo porque necesito ayuda. Pagué mi mensualidad completa en el portal web y me volvieron a cobrar $169.50.
30. En mi tarjeta aparece un cargo de 199 pesos de TelArg que desconozco.
31. La neta, me robaron 449 pesos de saldo en un solo día sin hacer llamadas ni usar datos.
32. Buenas noches. Desde que camibe de plan a Aigo Sin Limite me cobran $99.00 de más, es un cobro indebido.
33. Buenos días. Cargos que no reconozco en agosto: MXN 399 por roaming internacional. Exijo que los eliminen. Me urge que lo corrijan cuanto antes. Gracias de antemano.
34. Me quitaron el saldo de mi recarga de $449.50 sin explicacion alguna.
35. Buenas tardes. Me activaron plan adicional sin avisarme y ahora me cobran 79 pesos. Esto ya me pasó antes y nunca lo arreglaron. Quedo al pendiente.
36. Me cobraron $599.50 por una promoción que nunca acepté (ref. ticket 9195662).
37. Ya basta, aparece un cargo de MXN 499 por renta de equipo y yo no autoricé nada. Esto ya me pasó antes y nunca lo arreglaron.
38. En mi tarjeta aparece un cargo de 169 pesos de TelArg que desconozco. Es dinero que me hace falta. Espero una respuesta pronto.
39. En mi tarjeta aparece un cargo de $79.00 de TelArg que desconozco.
40. Ya di de baja la línea y me siguen cobrando el plan Libre 199, la atención no sirve. Nadie me da una respuesta concreta.
41. Les escribo porque necesito ayuda. No manches, aparece un cargo de MXN 129 por recarga automática y yo no autoricé nada. Ya revisé mi cuenta y no hay explicación alguna.
42. Me activaron contenido de terceros sin avisarme y ahora me cobran $99 MXN. Gracias.
43. Después de cancelar mi plan Joven 149 me cargaron MXN 149 a mi tarjeta.
44. Se me cobró una penalización de $129.00 que no corresponde porque ya cumplí mi contrato.
45. Tengo un problema con un cargo en mi cuenta. Me aplicaron un cargo por reconexión de $2,499 cuando mi línea nunca estuvo suspendida. Quiero una explicación clara de este cargo.
46. Cancelé mi servicio porque nunca tuve señal y me siguen haciendo cobros de $199.00 a mi tarjeta. Es muy frustrante que pase esto con mi dinero. Por favor, avísenme cualquier novedad.
47. Sin mi consentimiento me dieron de alta en renta de equipo y ya van 12 cobros.
48. Me siguen cobrando la línea que di de baja, ya llevo 6 meses pagando de más.
49. Me siguen cobrando la línea que di de baja, ya llevo 3 meses pagando de más. Nunca me avisaron de este cargo.
50. Pedi la devolucion de $249 hace cuatro dias y todavia no me abonan nada. Ya pedi una aclaracion y sigo sin respuesta.

### Problemas de Facturación

1. Tengo un problema con mi factura. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. Necesito que esto quede aclarado con urgencia. Gracias.
2. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado.
3. a quien corresponda. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra.
4. Tengo la línea a mi nombre. Mi pago de $79 sí se hizo pero la factura todavía marca vencido.
5. Buenas tardes. Gracias por nada, la factura de septiembre llegó otra vez con el plan equivocado. Necesito que lo resuelvan cuanto antes.
6. Saludos. Mi fecha de corte cambió al día 15 sin avisarme y ahora el estado de cuenta viene desfasado.
7. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Ya lo reporté y sigue igual. Espero una respuesta pronto.
8. falta el folio fiscal en mi factura de enero y no puedo deducirla.
9. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Quedo al pendiente.
10. Los consumos en mi factura de mayo no corresponden con mi historial de llamadas. Muchas gracias por su atención.
11. La factura de septiembre viene con el plan equivocado, aparece Empresarial 899 y yo contraté otro. Ya me había pasado antes y nadie lo arregló de fondo. Espero su pronta respuesta.
12. Mi fact de julio dice que tego adeudo pero el sistema de pagos dice q estoy al corriente. Espero una respuesta pronto...
13. Me cobran el monto completo, pero la factura dice un descuento que nunca aplicaron
14. El monto de este mes no es el que me dijeron en la sucursal, quiero una aclaración. Quedo atento a su respuesta.
15. Tengo un problema con mi factura. El CFDI de enero no llega y el cobro sí se hizo. Ya intenté arreglarlo por mi cuenta y no pude.
16. Mi fecha de corte cambió al día 28 sin avisarme y ahora el estado de cuenta viene desfasado. Por favor, avísenme cualquier novedad.
17. El estado de cuenta de marzo no detalla los conceptos, solo trae un total de 199 pesos.
18. Les escribo porque necesito ayuda. El recibo no trae el descuento que me prometieron en atención a clientes por la falla de servicio de llamadas. Es muy molesto andar aclarando esto.
19. Recibo mi factura 10 días tarde cada mes, ya casi vence cuando la tengo.
20. La factura de octubre viene con el plan equivocado, aparece Ahorro 99 y yo contraté otro.
21. Mi fecha de corte cambió al día 15 sin avisarme y ahora el estado de cuenta viene desfasado.
22. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé.
23. Solicité factura electrónica por WhatsApp y me siguen enviando el recibo impreso. Quedo atento a su respuesta.
24. Me cobran el monto completo, pero la factura dice un descuento que nunca aplicaron. Espero su pronta respuesta.
25. tengo un problema con mi fact. el cfdi de octubre no llega y el cobro si se hizo...
26. Los consumos en mi factura de octubre no corresponden con mi historial de llamadas.
27. Me llego la factura en blanco, sin importes, desde que active la entrega en papel.
28. El recibo no trae el descuento Que Me Prometieron En atencion a clientes Por La falla De Internet. muchas gracias por su atencion.
29. Tengo la línea a mi nombre. La factura trae un saldo anterior que ya pagué, el pago de $449.00 no aparece reflejado. Espero que lo resuelvan pronto.
30. Les escribo desde Chihuahua. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos en el portal web.
31. No puedo registrar mi correo para facturacion, la pagina me regresa un error (ref. caso 9138778). Necesito que lo resuelvan cuanto antes. Saludos cordiales.
32. El monto de este mes no es el que me dijeron x correo quiero una aclaracion. No me gusta quedarme con dudas de lo q pago.
33. Mi fecha de corte cambió al día 28 sin avisarme y ahora el estado de cuenta viene desfasado. Espero que lo resuelvan pronto.
34. el estado de cuenta de junio no detalla los conceptos, solo trae un total de mxn 1,199.
35. Buenos días. Mi plan Max 699 cuesta $699, pero en mi factura el monto total es de $999.50.
36. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a 399 pesos.
37. Mi fact sale a nombre de otra persona y necesito que corrijan los datos. Me urge tenerlo en regla.
38. Falta el folio fiscal en mi factura de febrero y no puedo deducirla. Espero que lo resuelvan pronto.
39. Buenas noches. El RFC de mi factura está mal capturado, necesito una factura corregida. Muchas gracias por su atención.
40. Necesito el comprobante fiscal (CFDI) de noviembre y el sistema no me lo genera (ref. reporte 760569). Necesito que lo resuelvan cuanto antes.
41. Buenas noches. Me cobran el monto completo, pero la fact dice un descuento que nunca aplicaron. Espero que lo reshelvan pronto.
42. Gracias por nada, la factura de abril llegó otra vez con el plan equivocado. Ya intenté arreglarlo por mi cuenta y no pude. Quedo al pendiente.
43. Mi recibo viene con fecha limite de pago distinta a la q me indicaron en el poratl web. Necesito q lo resuelvan cuanto antes. Quedo atento a su respuesta.
44. mi recibo no muestra el descuento prometido en soporte técnico, solo el total de MXN 129. Espero que lo resuelvan pronto. Quedo al pendiente.
45. Mi pago de MXN 99 sí se hizo pero la factura todavía marca vencido.
46. A quien corresponda. Necesito el comprobante fiscal (CFDI) de enero y el sistema no me lo genera (ref. folio 4979469).
47. La factura de septiembre viene con el plan equivocado, aparece Amigo Sin Límite y yo contraté otro. Espero una respuesta pronto.
48. EL ESTADO DE CUENTA DE FEBRERO NO DETALLA LOS CONCEPTOS, SOLO TRAE UN TOTAL DE 99 PESOS.
49. Buenas tardes. Falta el folio fiscal en mi factura de marzo y no puedo deducirla. Me urge tenerlo en regla.
50. La factura Dice Servicio Activo, pero En mi Ofiina No Tengo Internet Desde hace Dias.

### Quejas

1. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio.
2. Esperé media hora en ventas y nadie me contestó.
3. a quien corresponda. Para cancelar te piden mil cosas y mientras tanto te siguen cobrando el plan Social 299. Espero su pronta respuesta.
4. Mi queja: el call center cuelga cuando pido hablar con un supervisor.
5. Llevo varias horas esperando una solución y nadie me da respuesta. Quedo atento a su respuesta.
6. Me llaman seguido en horario laboral para ofrecerme promociones aunque ya les pedí que no. Espero que alguien por fin me responda.
7. Publicidad engañosa: prometen cobertura en Culiacan y no hay señal. Que alguien se haga responsable, x favor.
8. Publicidad engañosa: prometen cobertura en Chihuahua y no hay señal. Pido que tomen en cuenta mi queja.
9. Me cortan por un adeudo mínimo cuando ellos no mandan bien la factura.
10. Llevo más de 48 horas esperando una solución y nadie me da respuesta.
11. De plano, no es posible que una empresa tan grande tenga una atencion tan lenta. Ya perdí demasiado tiempo en esto.
12. Pedí 4 veces que corrijan mi factura y en ventas nadie me hace caso.
13. Una verguenza la atención: me dejaron esperando una hora y se me corto la llamada. Espero su pronta respuesta.
14. Buenas tardes. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Pido que tomen en cuenta mi queja.
15. Me vendieron el plan Amigo Sin Límite con condiciones que no cumplen, publicidad engañosa. Espero que alguien por fin me responda.
16. Les escribo desde Cancún. Chale, no es posible que una empresa tan grande tenga una atención tan lenta. Pido que tomen en cuenta mi queja.
17. Buenas noches. Fui a aclarar un cobro de 199 pesos y en atencion a clientes me dejaron esperando sin darme respuesta. Estoy pensando seriamente en cambiarme de compañia. Muchas gracias por su atencion.
18. Buenos días. Qué onda, qué mala atención, desde noviembre estoy tratando de aclarar un tema y nadie me resuelve.
19. A quien corresponda. No hay quien resuelva dudas en atención a clientes, todo son excusas y tiempos de espera.
20. EL ASESOR DE LA SUCURSAL DE SAN LUIS POTOSI FUE GROSERO Y ME CORTO LA LINEA X PAGAR UN DIA TARDE. SALUDOS CORDIALES.
21. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad.
22. Los asesores no saben nada, me dieron 2 respuestas distintas en la misma semana.
23. El proceso para cambiar de plan a Joven 149 fue lentísimo y lleno de requisitos absurdos. Todo esto es muy frustrante. Gracias.
24. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por chat. Muchas gracias por su atención.
25. Me vendieron el plan Prepago Básico con condiciones que no cumplen, publicidad engañosa. Que alguien se haga responsable, por favor.
26. Quiero dejar constancia de mi inconformidad. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. Que alguien se haga responsable, por favor.
27. Tengo una inconformidad q presentar. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes.
28. Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Todo esto es muy frustrante.
29. Tengo una inconformidad q presentar. En atencion a clientes me dijeron q no era su responsabilidad y me mandaron a otra area.
30. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por correo.
31. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Gracias.
32. El proceso para cambiar de plan a Familiar 549 fue lentísimo y lleno de requisitos absurdos. Estoy pensando seriamente en cambiarme de compañía.
33. Quiero presentar mi queja porque Héctor me colgó la llamada en plena aclaración. Quiero que alguien me explique por qué. Saludos cordiales.
34. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas.
35. estoy muy molesto con la forma en que fernanda me atendio, sin ninguna empatia.
36. Mi queja: el call center cuelga cuando pido hablar con un supervisor.
37. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención.
38. Los asesores no saben nada, me dieron 3 respuestas distintas en la misma semana. Espero q alguien con autoridad me atienda.
39. Subieron el precio de mi plan Premium 999 sin avisarme, esto me parece un abuso.
40. No manches, el trato de Miriam en el área de cobranza fue grosero, no pienso aguantar eso. Espero una respuesta pronto.
41. Una vergüenza la atención: me dejaron esperando media hora y se me cortó la llamada. Todo esto es muy frustrante.
42. Quiero presentar mi queja porque Fernanda me colgó la llamada en plena aclaración.
43. mala experiencia en la sucursal: tuve que repetir mi problema a cuatrro personas distintas. no es justo q uno tenga que insistir tanto
44. Otra vez me dejaron sin servicio y en atención a clientes solo me dan largas (ref. reporte 61940). Llevo años con ustedes y no merezco esto.
45. Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Llevo años con ustedes y no merezco esto.
46. Buenas tardes. Me cortan por un adeudo mínimo cuando ellos no mandan bien la factura.
47. Nadie da seguimiento a lo que pedí (ref. TKT-19656), siempre me dicen que ya está en proceso. Ya insistí varias veces y nadie me resuelve.
48. Una vergüenza la atención: me dejaron esperando veinte minutos y se me cortó la llamada. Llevo años con ustedes y no merezco esto. Saludos cordiales.
49. Buenas noches. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. Saludos cordiales.
50. LA ATENCION EN LA SUCURSAL DE SAN LUIS POTOSÍ ES PÉSIMA HACEN FILA UNA HORA PARA NADA. TODO ESTO ES MUY FRUSTRANTE. QQUEDO ATENTO A SU RESPUESTA.

### Reconocimientos

1. Les escribo con mucho gusto. La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. Realmente superó mis expectativas.
2. LAS MEJORES PROMOCIONES DEL MERCADO, SIGO CON TELARG X EL BUEN SERVICIO. SIGAN ASI.
3. Las mejores promociones del mercado, sigo con TelArg por el buen servicio. Realmente superó mis expectativas. Que tengan un excelente día.
4. Excelente cobertura en Veracruz, nunca he tenido problemas con mi plan Navega 249. Seguiré recomendándolos.
5. Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. caso 5692106), todo quedó resuelto. Realmente superó mis expectativas.
6. Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante.
7. Gracias a Jorge por la paciencia y la amabilidad al explicarme mi plan Max 699. Por esto sigo con TelArg.
8. Pagar mi recibo en la app es muy fácil, todo queda registrado al momento.
9. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo.
10. Mi señal es excelente incluso en mi colonia, felicidades a todo el equipo de redes. Realmente supero mis exprctativas.
11. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes.
12. Mi señal es excelente incluso en mi casa, felicidades a todo el equipo de redes.
13. buenas Tardes. andale, Qué Buen Servicio, en Menos De Media Hora Me solucionaron todo. da Gusto ver Q hay Empresas Q Cumplen. saludos cordiales...
14. El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias. Voy a seguir con ustedes muchos años más.
15. buenos dias. mi señal es excelente incluso en mi barrio, felicidades a todo el equipo de redes.
16. Gracias x la llamada de seguimiento, nunca me habia pasado q una empresa se preocupara asi. Eso es lo q uno espera de una buena compañía!!!
17. Quiero compartir mi experiencia. Muchas gracias por resolver mi problema por teléfono, el servicio fue impecable. Ojalá todas las empresas fueran así.
18. mi señal es excelente incluso en mi domicilio, felicidades a todo el equipo de redes.
19. Buenas tardes. Gracias por restablecer mi servicio el mismo día, atención de primera. De verdad quedé muy satisfecho.
20. La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más.
21. Las mejores promociones del mercado, sigo con TelArg x el buen servicio.
22. Buenas tardes. Mi señal es excelente incluso en toda la zona, felicidades a todo el equipo de redes.
23. Quiero dejar constancia de lo bien que me atendió Fernanda, un trato impecable. Un saludo a todo el equipo.
24. Gracias x la llamada de seguimiento, nunca me habia pasado q una empresa se preocupara asi.
25. Me encanta que la factura llegue cada mes por WhatsApp sin falta. Que tengan un excelente día.
26. Quiero dejar constancia de lo bien que me atendió Itzel, un trato impecable. Sigan así, de verdad se agradece.
27. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto.
28. Las mejores promociones del mercado, sigo con TelArg por el buen servicio.
29. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Voy a seguir con ustedes muchos años más.
30. Quiero dejar un comentario positivo. La promoción del plan Empresarial 899 es muy buena, gracias por pensar en los clientes. Ojalá todas las empresas fueran así.
31. Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. TKT-467455), todo quedó resuelto.
32. Les escribo con una buena noticia. Quiero felicitar a Daniela en el call center, el trato fue de lo mejor. Ojalá todo se mantenga así. Muchas gracias por su atención.
33. Mi señal es excelente incluso en mi colonia, felicidades a todo el equipo de redes.
34. Les felicito por la aplicacion, es muy facil de usar y pusdo ver mis consumos al instante. Realmente supero mis expectativas. Seguire reocmendándolos
35. Al final si me resolvieron en el call center, aunque tardaron bastante!!!
36. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Ojalá todas las empresas fueran así.
37. Por fin me devolvieron el cobro doble de $999.50, el asesor fue muy amable. Saludos cordiales.
38. Las mejores promociones del mercado, sigo con TelArg por el buen servicio. Se los quería hacer saber.
39. Excelente cobertura en Saltillo, nunca he tenido problemas con mi plan Social 299. Ojalá todas las empresas fueran así. Sigan así.
40. El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias.
41. Al final sí me resolvieron en la sucursal, aunque tardaron bastante. Gracias.
42. Gracias por avisarme antes de la fecha de corte, así no se me pasa el pago. Se los quería hacer saber. Muchas gracias por su atención.
43. La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. Eso es lo que uno espera de una buena compañía.
44. Un diez para atencion a clientes: resolvieron mi duda sin darme vueltas. Sigan asi, de verdad se agradece
45. A quien corresponda. Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. #8106181), todo quedó resuelto.
46. En soporte técnico me explicaron con paciencia el cobro de mi plan Libre 199, quedé muy conforme. Eso es lo que uno espera de una buena compañía.
47. Buenos días. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto.
48. Gracias a Carlos por la paciencia y la amabilidad al explicarme mi plan Social 299. Ojalá todas las empresas fueran así. Muchas gracias por su atención.
49. gracias, telarg: despues de dos dias sin internet lo arreglaron en cuanto lo reporte. da gusto ver q hay empresas q cumplen.
50. A quien corresponda. Gracias, TelArg: después de diez días sin internet lo arreglaron en cuanto lo reporté. Todo el equipo merece un aplauso. Gracias.
