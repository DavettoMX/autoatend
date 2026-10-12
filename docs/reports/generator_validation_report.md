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
| noise rate | PASS | realized 0.2984, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | pairwise Jaccard of top-100 tokens in [0.205, 0.333], required within [0.15, 0.80] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 20.02% duplicate texts; text length 41 / 111.5 / 244 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.307 | 0.290 | 0.299 | 0.250 |
| Cobros Indebidos | 0.307 | 1.000 | 0.333 | 0.325 | 0.227 |
| Problemas de Facturación | 0.290 | 0.333 | 1.000 | 0.242 | 0.205 |
| Quejas | 0.299 | 0.325 | 0.242 | 1.000 | 0.316 |
| Reconocimientos | 0.250 | 0.227 | 0.205 | 0.316 | 1.000 |

## Generator parameters

| Parameter | Value |
|---|---:|
| Templates per category | 36 / 34 / 34 / 34 / 32 |
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
| **Realized selection rate** | **0.2984** (tolerance ±0.015) |

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

Test comments with an identical text in train: 3617 of 14600 (24.77%). A temporal split does not prevent exact repeats of a template draw.

corpus.db SHA-256: `0ac4ac13bb7bd81b6c254244614f150d9e1226bc3483b169b6e8a0bb6b8515b6`

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. Mi teléfono dice sin servicio desde hace más de 48 horas, ya reinicié todo y nada. Necesito que me den una solución definitiva. Por favor, avísenme cualquier novedad.
2. mi chip nuevo no activa, llevo dos dias esperando que me den cobertura. ya perdi la cuenta de las veces que he insistido.
3. Necesito servicio de llamadas para trabajar y llevo cinco días sin nada, esto ya es urgente. Gracias de antemano.
4. De plano, otra vez sin cobertura, ya es la tercera vez este mes que se cae.
5. Necesito internet para trabajar y llevo tres días sin nada, esto ya es urgente.
6. A quien corresponda. La verdad, otra vez sin datos moviles, ya es la tercera vez este mes q se cae. Quedo atento a su resspuesta.
7. soy cliente de telarg desde hace años. en mi trabajo nunca hay señal, tengo q salur a la calle para poder llamar. quedo atento a su respuesta.
8. Mi línea sigue sin servicio tras la reconexión que pagué (ref. TKT-347757).
9. La neta, otra vez sin internet, ya es la tercera vez este mes que se cae. Cada vez que llamo me dicen que ya está en proceso.
10. Dicen que hay mantenimiento en Monterrey pero ya llevamos una semana y seguimos sin internet. Por favor, avísenme cualquier novedad.
11. Desde hace cuatro días estoy sin señal en mi colonia y nadie me da una solución.
12. No tengo internet en mi casa toda la tarde, ya pedí que lo revisen (ref. ticket 523851). Espero una respuesta pronto.
13. Hice mi recarga de MXN 49 hace más de 48 horas y sigo sin servicio de llamadas. Ya probé todo lo que me dijeron y sigue igual.
14. Activé el roaming para viajar y al llegar a Tampico no tuve internet. Espero que lo atiendan hoy mismo.
15. No tngo servicio de ensajes en toda la zona toda la tarde, ya pedi que lo revisen (ref. ticket 4778244). Quedo al pendiente.
16. Activé el roaming para viajar y al llegar a Cancún no tuve internet.
17. Saludos. El módem se queda con la luz roja y no hay internet, llevo diez días igual.
18. La portabilidad de mi número no se completó y estoy sin línea desde hace diez días. Cada vez que llamo me dicen que ya está en proceso. Gracias.
19. No tengo datos móviles en mi trabajo en fines de semana, ya pedí que lo revisen (ref. caso 1711095).
20. El internet de mi plan Familiar 549 no funciona desde marzo y no hay quien lo repare.
21. No puedo mandar ni recibir mensajes de texto en Culiacán desde hace cuatro días.
22. Activé el roaming para viajar y al llegar a Morelia no tuve internet.
23. A quien corresponda. Mi línea no registra en la red desde hace cinco días, llevo todo ese tiempo sin poder comunicarme.
24. A quien corresponda. O sea, otra vez sin servicio de llamadas, ya es la tercera vez este mes que se cae. Ya llamé tres veces y nadie me resuelve. Espero una respuesta pronto.
25. Mi línea no registra en la red desde hace tres días, llevo todo ese tiempo sin poder comunicarme. Espero su pronta respuesta.
26. En mi colonia no hay cobertura en horario laboral, parece que la antena está caída. Ya probé todo lo que me dijeron y sigue igual.
27. ME QUEDÉ SIN SERVICIO DE MENSAJES JUSTO CUANDO MAS LO NECESITABA ¿HASTA CUANDO VAN A ARREGLARLO? NECESITO Q ME DEN UNA SOLUCION DEFINITIVA.
28. Activé el roaming para viajar y al llegar a Villahermosa no tuve cobertura. Gracias de antemano.
29. Hay pérdida de servicio de mensajes por las mañanas en toda mi calle, y ya son varios días con el mismo problema. Espero su pronta respuesta.
30. Qué onda, llevo diez días sin servicio de llamadas y sigo pagando mi plan Ahorro 99. Muchas gracias por su atención.
31. Me prometieron restablecer el servicio en 5 días y siguen sin cumplir. Ya perdí la cuenta de las veces que he insistido. Quedo al pendiente.
32. En mi barrio no hay cobertura en horario laboral, parece que la antena esta caida. Esto afecta a toda mi familia.
33. Solo me aparece emergenias en la pantalla, no puedo usar mi línea por las mañanas.
34. Mi teléfono dice sin servicio desde hace varias horas, ya reinicié todo y nada. Quedo al pendiente.
35. No tengo servicio de llamadas en mi domicilio de madrugada, ya pedí que lo revisen (ref. #6840106). Ya llamé tres veces y nadie me resuelve.
36. Desde hace quince días estoy sin datos móviles en mi trabajo y nadie me da una solución.
37. En mi colonia no hay cobertura por las mañanas, parece que la antena esta caída
38. no tengo datos móviles en toda mi calle de madrugada, ya pedí que lo revisen (ref. reporte 833912). ya me había pasado antes este mes.
39. Me prometieron restablecer el servicio en 30 días y siguen sin cumplir.
40. No me llegan los códigos de verificación por SMS desde hace varias horas, creo que hay una falla en la red.
41. Buenas tardes. Se me fue por completo la señal en Cancún y no he podido hacer ni recibir llamadas. Quiero que alguien me explique qué está pasando. Muchas gracias por su atención.
42. A quien corresponda. Mi línea no registra en la red desde hace dos días, llevo todo ese tiempo sin poder comunicarme. Saludos cordiales.
43. Buenas tardes. No tengo cobertura en el centro toda la tarde, ya pedí que lo revisen (ref. #183865). Ya me había pasado antes este mes. Gracias.
44. Hola. solo me Aparece Emergencias En La Pantalla No puedo Usar mi linea de madrugada.
45. Saludos. Mi línea está suspendida sin motivo y no puedo usar nada. Quedo atento a su respuesta.
46. Mi chip nuevo no activa, llevo mas de 48 horas esperando q me den datos móvilles. Nadie me da una fecha concreta de solución.
47. Buenas noches. Solo me aparece emergencias en la pantalla, no puedo usar mi línea en horario laboral.
48. Activé el roaming para viajar y al llegar a Culiacán no tuve servicio de mensajes.
49. La verdad, otra vez sin internet, ya es la tercera vez este mes q se cae. Esto afecta a toda mi familia...
50. En mi trabajo nunca hay señal, tengo que salir a la calle para poder llamar.

### Cobros Indebidos

1. Pagué mi mensualidad completa por teléfono y me volvieron a cobrar $399 MXN. Nunca me avisaron de este cargo.
2. Me quitaron el saldo de mi recarga de MXN 149 sin explicación alguna. Me urge que lo corrijan cuanto antes. Muchas gracias por su atención.
3. Les escribo Porque necesito ayuda. Cobro indebido: me descontaron mXN 249 dos Veces El Mismo dia por Lo mismo. es Muy Frustrante Q Pase esto con mi Dinero!!!
4. Soy cliente de TelArg desde hace años. Me cobraron roaming internacional de MXN 1,199 sin salir de Veracruz. Por favor, avísenme cualquier novedad.
5. Buenas noches. Me siguen cobrando la línea que di de baja, ya llevo 3 meses pagando de más. Espero una respuesta pronto.
6. Tengo un problema con un cargo en mi cuenta. Pedí la devolución de $1,999 hace diez días y todavía no me abonan nada. Ya revisé mi cuenta y no hay explicación alguna.
7. Me cobran servicio premium como si lo hubiera pedido ya me canse de pagar x lo q no uso. No pienso pagar algo q no pedi.
8. A quien corresponda. Despues de cancelar mi plan Navega 249 me cargaron $249.00 a mi tarjeta.
9. Se me cobró una penalización de $149.50 que no corresponde porque ya cumplí mi contrato.
10. Cancelé renta de equipo hace meses y me siguen descontando $149 cada mes.
11. El débito automático me quitó $299.00 dos veces en agosto, necesito el reembolso. Exijo que lo revisen cuanto antes. Espero su pronta respuesta.
12. me aplicaron un cargo por reconexion de $699 mxn cuando mi linea nunca estuvo suspendida. ya revise mi cuenta y no hay explicación alguna. x favor, avisenme cualquier novedad!!!
13. Soy cliente de TelArg desde hace años. Ya reclamé el cobro de reconexión en el portal web y nadie me devuelve mi dinero (caso 4592462).
14. Mi banco me reporta un cargo no reconocido de TelArg por $129 MXN, quiero aclararlo.
15. Mi banco me reporta un cargo no reconocido de TelArg por $799, quiero aclararlo. Saludos cordiales.
16. Sin mi consentimiento me dieron de alta en seguro de equipo y ya van 12 cobros. Nadie me da una respuesta concreta.
17. mi banco me reporta un cargo no reconocido de TelArg por $899, quiero aclararlo.
18. Sin mi consentimiento me dieron de alta en contenido de terceros y ya van 4 cobros. Gracias de antemano.
19. Soy cliente de TelArg desde hace años. Sin mi consentimiento me dieron de alta en recarga automatica y ya van 2 cobros. Nadie me da una respuesta concreta. Por favor, avisenme ucalquier novedad.
20. Sin mi consentimiento me dieron de alta en mensajes premium y ya van 8 cobros. Gracias de antemano.
21. En mi tarjeta aparece un cargo de 399 pesos de TelArg que desconozco.
22. Hice la recarga de $49.00 y me cobraron una comisión que no me explicaron.
23. Se me cobró una penalización de $169.50 que no corresponde porque ya cumplí mi contrato.
24. Hice la recarga de 99 pesos y me cobraron una comision q no me explicaron. Muxhas gracias x su atencion...
25. Mi plan Libre 199 cuesta $199 MXN y me estan cobrando $349.00, no es lo q acordamos.
26. Pedí la cancelación de mi línea y aun así me cobraron MXN 899 en abril.
27. Soy cliente de TelArg desde hace años. Quiero que me reembolsen los $349 que me cobraron de más (reporte 6973504).
28. Termine de pagar el equipo hace meses y aún me cargan $199 MXN memsuales.
29. Pague mi mensualidad completa en la app y me volvieron a cobrar $1,199.00. Ya pedi una aclaracion y sigo sin respuesta. Quedo atento a su respuesta.
30. Buenas noches. Pedí la cancelación de mi línea y aun así me cobraron MXN 449 en marzo. Esto ya me pasó antes y nunca lo arreglaron.
31. Desde que cambié de plan a Joven 149 me cobran $49 MXN de mas, es un cobro indebido
32. Pedí la cancelación de mi línea y aun así me cobraron MXN 1,999 en marzo. Me urge que lo corrijan cuanto antes. Espero su pronta respuesta.
33. Me cobraron $129 MXN por servicio premium que nunca contraté.
34. Desde que cambié de plan a Premium 999 me cobran $599 MXN de más, es un cobro indebido. Nunca me avisaron de este cargo.
35. mi banco me reporta un cargo no reconocido de TelArg por 49 pesos, quiero aclararlo.
36. Les escribo porque necesito ayuda. El débito automático me quitó $899.00 dos veces en enero, necesito el reembolso.
37. Terminé de pagar el equipo hace meses y aún me cargan 99 pesos mensuales. Quedo al pendiente.
38. Les escribo por un cobro que no reconozco. Me siguen cobrando la línea que di de baja, ya llevo 6 meses pagando de más.
39. Buenas noches. Cargos que no reconozco en enero: $79.50 por reconexión. Exijo que los eliminen.
40. Desconozco el cobro de $49 del día 20 de febrero, quiero que lo reviertan.
41. A quien corresponda. Pedí la devolución de $349.00 hace cinco días y todavía no me abonan nada.
42. Buenos días. Desde que cambié de plan a Nocturno 79 me cobran MXN 49 de más, es un cobro indebido. Quiero que me devuelvan mi dinero. Espero una respuesta pronto.
43. Cancele suscripción de música hace meses y me siguen descontando $169.00 cada mes. Exijo que lo revisen cuanto antes.
44. Me cobraron MXN 199 por seguro de equipo que nunca contraté. Saludos cordiales.
45. Me cobraron $349 MXN por paquete de datos extra que nunca contraté.
46. Activaron contenido de terceros sin pedírmelo y me cobraron $99.00 en el estado de mi tarjeta. Ya pedí una aclaración y sigo sin respuesta.
47. Cancele paquete de datos extra hace meses y me siguen descontando $349 cada mes.
48. Me activaron plan adicional sin avisarme y ahora me cobran $199 MXN. Ya revisé mi cuenta y no hay explicación alguna.
49. Soy cliente de TelArg desde hace años. Me cobraron $799 MXN por una promoción que nunca acepté (TKT-34887).
50. Les escribo porque necesito ayuda. Activaron recarga automática sin pedírmelo y me cobraron $229.00 en el estado de mi tarjeta. Es muy frustrante que pase esto con mi dinero.

### Problemas de Facturación

1. Buenas tardes. No puedo registrar mi correo para facturacion, la pagina me regresa un error (TKT-57642). Gracias de antemano.
2. Buenos dias. Pedí la fact de mi equipo y me mandaron la de otro número de línea. No me gusta quedarme con dudas de lo q pago!!!
3. Tengo la línea a mi nombre. Mi factura de abril dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente.
4. Buenas noches. Me llego la factura en blanco, sin importes, desde que active la entrega en papel. Nadie me da una explicación clara. Quedo atento a su respuesta.
5. Necesito el comprobante fiscal (CFDI) de marzo y el sistema no me lo genera (caso 668448). Nadie me da una explicación clara.
6. Buenas tardes. Mi factura de noviembre dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente. No me gusta quedarme con dudas de lo que pago. Saludos cordiales.
7. Mi plan Joven 149 cuesta $149 MXN, pero en mi factura el monto total es de $169 MXN. Ya me había pasado antes con otro periodo. Quedo atento a su respuesta.
8. Quiero cambiar mi factura a entrega en la app y no encuentro la opcion x inngun lado.
9. Quiero cambiar mi factura a entrega en papel y no encuentro la opción por ningún lado. Necesito el documento correcto con urgencia. Gracias.
10. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra.
11. Tengo la línea a mi nombre. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos.
12. Quiero cambiar mi fact a entrega por WhatsApp y no encuentro la opcion por ningun lado.
13. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Gracias de antemano.
14. Mi direccion de fact es incorrecta y x eso no llegan mis recibos. Espero su pronta respuesta.
15. Tengo un problema con mi fact. Necesito el comprobante fiscal (CFDI) de enero y el sistema no me lo genera (folio 16125). Espero una respuesta pronto.
16. Ay no mi estado de cuenta no coincide con lo que me dijeron en el call xenter.
17. Hola. La factura trae un saldo anterior que ya pagué, el pago de $899 no aparece reflejado.
18. Hola. No puedo descargar mi factura en el portal web, marca error cada vez. Espero una respuesta pronto.
19. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Ya me había pasado antes con otro periodo. Quedo atento a su respuesta.
20. Cambié de plan a Social 299 y la factura sigue mostrando la tarifa anterior.
21. Me llegó la factura en blanco, sin importes, desde que activé la entrega en la app.
22. El estado de cunta de enero no detalla los conceptos, solo trae un toatl de MXN 129.
23. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Necesito que lo corrijan cuanto antes. Quedo atento a su respuesta.
24. Falta el folio fiscal en mi factura de marzo y no puedo deducirla.
25. Buenas noches. Mi fecha de corte cambió al dia 5 sin avisarmme y ahora el estado de cuenta vviene desfasado.
26. Me llego la factura en blanco, sin importes, desde q active la entrega por correo electronico.
27. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Necesito que lo corrijan cuanto antes.
28. Caray, mi estado de cuenta no coincide con lo que me dijeron en el área de cobranza.
29. Les escribo porque necesito ayuda. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. Quedo al pendiente.
30. El portal muestra un monto de factura distinto al del recibo que me llegó, $99.50 contra $2,499.
31. Quiero cambiar mi factura a entrega por correo electrónico y no encuentro la opción por ningún lado.
32. Falta el folio fiscal en mi factura de diciembre y no puedo deducirla. Necesito q lo corrijan cuanto antes. Quedo atento a su respuesat
33. Tengo la línea a mi nombre. Mi fact llega con los datos de mi domicilio anterior aunque ya lo actualice.
34. El ciclo de facturación cambió y ahora tengo dos recibos casi juntos de 399 pesos cada uno.
35. Mi factura de noviembre dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente.
36. Tengo la línea a mi nombre. Mi pago de $249.50 sí se hizo pero la factura todavía marca vencido.
37. Pedí la factura de mi equipo y me mandaron la de otro número de línea. Ya envié mis datos otra vez y no pasa nada.
38. Tengo un problema con mi factura. q barbaridad, mi estado de cuenta no coincide con lo que me dijeron en el call center. Por favor, avísenme cualquier novedad.
39. recibbo mi factura 7 dias tarde cada mes, ya casi vence cuando la tengo...
40. La fact de mi linea empresarial llega a otra razon social, ya mande los documentos en el portal web. Nadie me da una explicacion clara. Quedo al pendiente.
41. Los consumos en mi factrua de enero no corresponden con mi historial de llamadas. Ya me habia pasado antes con otro periodo. Gracias
42. Buenas noches. Solicité factura electrónica por correo electrónico y me siguen enviando el recibo impreso. Necesito que lo corrijan cuanto antes.
43. No aparece el pago que hice en el portal web en mi estado de cuenta, sigo con saldo pendiente. Ya lo reporté y sigue igual. Gracias.
44. Tengo la línea a mi nombre. Mi recibo viene con fecha límite de pago distinta a la que me indicaron por teléfono. Ya llamé tres veces y nadie me resuelve.
45. Mi recibo viene con fecha limite de pago distinta a la q me indicarln en el portal web. Necesito el documento correcto con urgecnia.
46. Recibo mi factura 7 días tarde cada mes, ya casi vence cuando la tengo. Nadie me da una explicación clara.
47. Mi recibo no muestra el descuento prometido en atención a clientes, solo el total de $79.50. Necesito que lo corrijan cuanto antes.
48. Los consumos en mi factura de marzo no corresponden con mi historial de llamadas.
49. Mi direccion de facturación es incorrecta y x eso no llegan mis recibos. Quiero q me confirmen cuando estara corregido.
50. Les escribo por un tema de facturación. La fecha de vencimiento de mi recibo no es la que me dieron, por eso se me juntan los periodos. Es muy molesto andar aclarando lo mismo cada mes.

### Quejas

1. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Quedo atento a su respuesta.
2. En el call center me dijeron que no era su responsabilidad y me mandaron a otra área. Pido que tomen en cuenta mi queja.
3. EL PROCESO PARA CAMBIAR DE PLAN A FAMILIAR 549 FUE LENTISIMO Y LLENO DE REQUISITOS ABSURDOS...
4. La neta, no es posible que una empresa tan grande tenga una atención tan lenta. Ya llamé tres veces y nadie me resuelve.
5. El proceso para cambiar de plan a Premium 999 fue lentísimo y lleno de requisitos absurdos. Que alguien se haga responsable, por favor.
6. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por chat.
7. espere una hora en el call center y nadie me contesto. me parece inaceptable q nadie de la cara.
8. El equipo q me entregaron llego tarde y con la caja abierta, muy mal sercicio. Es una falta de respeto con los clientes.
9. De plano, qué mala atención, desde julio estoy tratando de aclarar un tema y nadie me resuelve.
10. Espere mas de una hora en ventas y nadie me contesto. Pido q tomen en cuenta mi queja.
11. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. Llevo años con ustedes y no merezco este trato. Quedo atento a su respuesta.
12. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por teléfono.
13. Me prometieron una llamada de vuelta y después de tres días sigo esperando (ticket 70145).
14. O sea, no es posible que una empresa tan grande tenga una atención tan lenta.
15. Me llaman seguido de madrugada para ofrecerme promociones aunque ya les pedí que no. Ya llamé tres veces y nadie me resuelve.
16. Me vendieron el plan Empresarial 899 con condiciones que no cumplen, publicidad engañosa.
17. mi queja es x la falta de interes en soporte técnico, nunca devuelven las llamadas. ya perdi demasiado tiempo en esto. espero q alguin por fin me responda.
18. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas.
19. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. Muchas gracias por su atención.
20. Estoy muy molesto con la forma en que Paola me atendió, sin ninguna empatía. Quedo atento a su respuesta.
21. La sucursal cierra antes del horario que anuncian y ya van 3 veces que voy en balde.
22. Me prometieron una llamada de vuelta y después de más de 48 horas sigo esperando (caso 5403946). Muchas gracias por su atención.
23. Ay no, no es posible que una empresa tan grande tenga una atención tan lenta. Quiero que alguien me explique por qué.
24. En atención a clientes me dijeron que no era su responsabilidad y me mandaron a otra área. Espero una respuesta pronto.
25. El proceso para cambiar de plan a Max 699 fue lentísimo y lleno de requisitos absurdos. Pido que tomen en cuenta mi queja!!!
26. Una vergüenza la atención: me dejaron esperando una hora y se me cortó la llamada. Quiero que alguien me explique por qué.
27. Quiero presentar mi queja porque Héctor me colgó la llamada en plena aclaración. Quedo atento a su respuesta.
28. Me prometieron una llamada de vuelta y después de diez días sigo esperando (caso 6787831). Espero que alguien con autoridad me atienda. Gracias.
29. Esperé cuarenta minutos en el call center y nadie me contestó.
30. Exijo q me atiendan como cliente, no es la primera vez q me ignoran en la sucursal...
31. Les escribo desde Culiacán. Subieron el precio de mi plan Max 699 sin avisarme, esto me parece un abuso. No es justo que uno tenga que insistir tanto. Muchas gracias por su atención.
32. Subieron el precio de mi plan Empresarial 899 sin avisarme, esto me parece un abuso. Pido que tomen en cuenta mi queja.
33. Tengo la linea a mi nombre. En la sucursal me dijeron que no era su responsabilidad y me mandaron a otra área...
34. Espere mas de una hora en spoorte tecnico y nadie me contesto. Llevo años con ustedes y no merezco este trato. Espero q alguien x fin me responda.
35. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas.
36. Me camiaron las condiciones de mi plan Max 699 sin pedirme autorizacion, ya no es lo q contrate. Muchas gracias x su atencion...
37. Espere mas de una hora en el area de cobranza y nadie me contesto. Ya llame tres veces y nadie me resuelve. Quedo atento a su respuesta...
38. Buenas tardes. Qué barbaridad, el trato de Gerardo en el call center fue grosero, no pienso aguantar eso.
39. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio. Que alguien se haga responsable, por favor. Quedo atento a su respuesta.
40. Saludos. La atencion en la sucursal de Oaxaca es ppesima, acen fila cuarenta minutos para nada. Espero que alguien con autoridad me atienda.
41. La atención en la sucursal de Villahermosa es pésima, hacen fila cuarenta minutos para nada.
42. Los asesores no saben nada, me dieron 2 respuestas distintas en la misma semana.
43. Soy cliente de TelArg desde hace años. La atención en la sucursal de Oaxaca es pésima, hacen fila veinte minutos para nada.
44. Tengo una inconformidad que presentar. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. Ya perdí demasiado tiempo en esto.
45. Buenos días. Me vendieron el plan Premium 999 con condiciones que no cumplen, publicidad engañosa. Llevo años con ustedes y no merezco este trato.
46. No hay quien resuelva dudas en atencion a clientes, todo son excusas y tiempos de espera. Ya perdí demasiado tiempo en esto.
47. Tengo una inconformidad que presentar. En el call center me dijeron que no era su responsabilidad y me mandaron a otra área.
48. Les escribo desde León. Una vergüenza la atención: me dejaron esperando casi dos horas y se me cortó la llamada. Quedo atento a su respuesta.
49. Qué barbaridad, el trato de Gerardo en atención a clientes fue grosero, no pienso aguantar eso.
50. Buenos días. En la sucursal me dijeron que no era su responsabilidad y me mandaron a otra área. Pido que tomen en cuenta mi queja.

### Reconocimientos

1. Quiero dejar constancia de lo bien q me atendio Verónica, un trato impecable. Todo el equipo merece un aplauso.
2. Hola. Les agradezco la atención durante julio, me ayudaron mucho con mi cambio de línea. Gracias.
3. El cambio a plan Total 399 fue sencillo y sin trámites, gracias por facilitarlo. Se los quería hacer saber. Seguiré recomendándolos.
4. Mi señal es excelente incluso en mi trabajo, felicidades a todo el equipo de redes.
5. En la sucursal de Morelia me atendieron de maravilla, todo fue rápido y claro.
6. Gracias a Jorge por la paciencia y la amabilidad al explicarme mi plan Amigo Sin Límite. Sigan así, de verdad se agradece.
7. Hector me atendio con mucha atencion y profesionalismo, merece un reconocimiento x su trabajo. Me dio mucho gusto ser su cliente.
8. Felicito a TelArg por la claridad de mi factura, ahora sí entiendo cada concepto. Eso es lo que uno espera de una buena compañía.
9. Soy cliente de TelArg desde hace años. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad.
10. Quiero compartir mi experiencia. De plano, qué gusto que con atención a clientes se resolvió mi problema al primer intento. Es un gusto que existan personas así.
11. Mi expereincia x tel fue de lo mejor, Jorge se tomo el tiempo de explicarme todo.
12. La verdad, qué gusto que con atención a clientes se resolvió mi problema al primer intento. De verdad quedé muy satisfecho. Gracias.
13. Les escribo desde Queretaro. El equipo que me entregaron llego antes de lo prometido y en perfecto estado, muchas grcias.
14. Buenos dias. Mi reconocimiento al tecnico que llegó puntual y dejo mi internet funcionando perfecto. Se nota que se preocupan por sus clientes. Gracias!!!
15. Quiero compartir mi experiencia. Reconozco el esfuerzo que hacen en el call center, siempre están al pendiente de nosotros. Muchas gracias por su atención.
16. Soy cliente de TelArg desde hace años. Reconozco el esfuerzo que hacen en el call center, siempre están al pendiente de nosotros. Con gusto volvería a acudir con ustedes.
17. Me encanto la atencion en la sucursal, respondieron en minutos y con mucga amabilidad.
18. Les escribo desde Morelia. Me encantó la atención en la sucursal, respondieron en minutos y con mucha amabilidad. Eso es lo que uno espera de una buena compañía. Que tengan un excelente día.
19. Excelente cobertura en Villahermosa, nunca he tenido problemas con mi plan Amigo Sin Límite. Es un gusto que existan personas así.
20. Felicidades por la cobertura en Veracruz, hasta en carretera tengo señal. Sigan así.
21. Mi experiencia en el portal web fue de lo mejor, Ricardo se tomó el tiempo de explicarme todo. Se nota que se preocupan por sus clientes.
22. Me encanto la atencion por chat, respondieron en minutos y con mucha amabilidad. De verdad quede muy satisfecho. Saludos cordiales.
23. Gracias gracias y gracias x fin alguien resolvio mi problema de manera agil. Muchas gracias por su atencion.
24. Gracias, gracias y gracias, por fin alguien resolvió mi problema de manera ágil. Con gusto volvería a acudir con ustedes.
25. Quiero reconocer públicamente la labor de Ricardo en el área de cobranza, un trabajo verdaderamente profesional. Sigan así, de verdad se agradece.
26. La neta, q gusto que con el area de cobranza se resolvió mi problema al primer intento.
27. Quiero dejar un comentario positivo. Mi reconocimiento al tecnico que llegó puntual y dejó mi internet funcionando perfecto.
28. Felicito a TelArg por la claridad de mi factura, ahora sí entiendo cada concepto.
29. Un diez para soporte técnico: resolvieron mi duda sin darme vueltas.
30. Gracias, gracias y gracias, por fin alguien resolvió mi problema de manera ágil.
31. Tengo la linea a mi nombre. Excelente cobertura en Cancun, nunca he tenido problemas con mi plan Social 299. Muchas gracias por su atencion.
32. Quiero agradecer la rapidez con la que atendieron mi solicitud (folio 937249), todo quedó resuelto. Ojalá siempre sea así de fácil. Saludos cordiales.
33. Felicidades por la cobertura en San Luis Potosí, hasta en carretera tengo señal.
34. Felicito a TelArg por la claridad de mi factura, ahora sí entiendo cada concepto.
35. Tengo la línea a mi nombre. Quiero dejar constancia de lo bien que me atendió Miriam, un trato impecable. Es un gusto que existan personas así.
36. Me encantó la atención en la sucursal, respondieron en minutos y con mucha amabilidad.
37. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad. Se los quería hacer saber. Saludos cordiales.
38. Buenas tardes. Gracias, gracias y gracias, por fin alguien resolvió mi problema de manera ágil. De verdad quedé muy satisfecho.
39. Quiero felicitar a Miriam en soporte técnico, el trato fue impecable. Sigan así, de verdad se agradece.
40. El cambio a plan Amigo Sin Límite fue sencillo y sin trámites, gracias por facilitarlo.
41. Las mejores promociones del mercado, sigo con TelArg por el buen servicio.
42. Quiero felicitar a Fernanda en soporte técnico, el trato fue de diez.
43. Gracias por restablecer mi servicio el mismo día, atención de primera.
44. Muchas gracias por resolver mi problema por chat, el servicio fue de primera. Por esto sigo con TelArg.
45. La promoción del plan Nocturno 79 es muy buena, gracias por pensar en los clientes. De verdad quedé muy satisfecho. Gracias.
46. Gracias x la llamada de sguimiento, nunca me habia pasado que una empresa se preocupara asi...
47. Muy contento con el internet de mi plan Premium 999, la velocidad es increíble.
48. Gracias a Itzel por la paciencia y la amabilidad al explicarme mi plan Prepago Básico. Me dio mucho gusto ser su cliente.
49. Les agradezco la atención durante marzo, me ayudaron mucho con mi cambio de línea.
50. Verónica me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Me dio mucho gusto ser su cliente. Un saludo a todo el equipo.
