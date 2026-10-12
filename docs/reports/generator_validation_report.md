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
| noise rate | PASS | realized 0.3003, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | pairwise Jaccard of top-100 tokens in [0.190, 0.342], required within [0.15, 0.80] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 20.20% duplicate texts; text length 40 / 111.0 / 232 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.307 | 0.299 | 0.325 | 0.258 |
| Cobros Indebidos | 0.307 | 1.000 | 0.342 | 0.299 | 0.212 |
| Problemas de Facturación | 0.299 | 0.342 | 1.000 | 0.250 | 0.190 |
| Quejas | 0.325 | 0.299 | 0.250 | 1.000 | 0.299 |
| Reconocimientos | 0.258 | 0.212 | 0.190 | 0.299 | 1.000 |

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
| **Realized selection rate** | **0.3003** (tolerance ±0.015) |

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

Test comments with an identical text in train: 3716 of 14600 (25.45%). A temporal split does not prevent exact repeats of a template draw.

corpus.db SHA-256: `b810bad233f9180484082e056b0aefd7676a63e5c4c87441c3ffc2b309760d71`

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. Desde hace dos días estoy sin servicio de llamadas en toda mi calle y nadie me da una solución. Nadie le da seguimiento a esto.
2. Saludos. El internet de mi plan Empresarial 899 no funciona desde diciembre y no hay quien lo repare.
3. Mi chip nuevo no activa, llevo diez días esperando que me den servicio de llamadas. Espero que lo atiendan pronto.
4. Me quedé sin servicio de mensajes justo cuando mas lo necesitaba, ¿hasta cuando van a arreglarlo? Ya perdi la cuenta de las veces que he insistido.
5. saludos. en ttoda la zona no hay cobertura toda la tarde, parece que la antena está caída. x favor, avísenme cualquier novedad!!!
6. A quien corresponda. No tengo cobertura en mi oficina en fines de semana, ya pedi q lo revisen (ref. TKT-97787). Saludos cordiales.
7. Desde que cambiaron la antena en mi barrio perdí servicio de mensajes y desde entonces sigo igual.
8. A quien corresponda. El módem se queda con la luz roja y no hay internet, llevo varias horas igual. Ya me está afectando en el trabajo.
9. Caray, llevo cinco días sin cobertura y sigo pagando mi plan Empresarial 899. Sigo pagando completo y no obtengo lo que corresponde.
10. Me quedé sin servicio de mensajes justo cuando más lo necesitaba, ¿hasta cuándo van a arreglarlo? Sigo pagando completo y no obtengo lo que corresponde.
11. En mi barrio no hay cobertura en las noches, parece que la antena está caída.
12. Desde hace cuatro días estoy sin servicio de mensajes en toda la zona y nadie me da una solución. Ya insistí varias veces y nadie me resuelve.
13. Les escribo x una falla en mi línea. El tecnico nunca llego a revisar la falla de señal, ya llevo una semzna de espera.
14. No tengo servicio de llamadas en mi domicilio toda la tarde, ya pedí que lo revisen (ref. reporte 295888). Ya perdí la cuenta de las veces que he insistido.
15. Hola. La portabilidad de mi número no se completó y estoy sin línea desde hace cinco días. Ya perdí la cuenta de las veces que he insistido.
16. Cuando intento llamar solo escucho q el número no está disponible, ya llevo quince dias asi...
17. Hice mi recarga de $49 hace cuatro dias y sigo sin internet. Quiero q algguien se haga cargo.
18. Saludos. Ya cambié mi teléfono y sigo sin señal, el problema es de la red.
19. Buenas tardes. Mi línea está suspendida sin motivo y no puedo usar nada. Ya me había pasado antes este mes.
20. No tengo datos moviles aunque mi saldo es suficiente (ref. reporte 85937)!!!
21. les escribo porque necesito ayuda. No puedo mandar ni recibir mensajes de texto en Veracruz desde hace cinco días. Quiero que alguien se haga cargo.
22. Solo me aparece emergencias en la pantalla, no puedo usar mi línea en horario laboral.
23. Cambié de SIM y quedé incomunicado, sin servicio de llamadas desde hace diez días. Ya perdí la cuenta de las veces que he insistido. Quedo atento a su respuesta.
24. Desde hace una semana estoy sin señal en mi colonia y nadie me da una solucion. Espero q lo atiendan pronto.
25. Hola. Cuando intento llamar solo escucho que el número no está disponible, ya llevo cuatro días así.
26. Hice mi recarga de $49 MXN hace cinco días y sigo sin internet. Es muy frustrante estar así. Quedo atento a su respuesta.
27. en el centro no hay cobertura en horario laboral parece q la antena esta caida. sigo pagando completo y no obtengo lo q corresponde. gracias de antemano.
28. Mi chip nuevo no activa, llevo cinco días esperando que me den servicio de llamadas. Ya probé todo lo que me dijeron y sigue igual.
29. A quien corresponda. La portabilidad de mi número no se completó y estoy sin línea desde hace cinco días.
30. Me prometieron restablecer el servicio en 7 días y siguen sin cumplir. Espero su pronta respuesta.
31. Mi chip nuvo no activa, llevo cinco días esperando q me den imternet. Muchas gracias por su atencion!!!
32. Les escribo porque necesito ayuda. Mi línea sigue sin servicio tras la reconexión que pagué (ref. ticket 648000).
33. Desde q cambiaron la antena en toda la zona perdí datos moviles y desde entonces sigo igual. Ya probe todo lo que me dijeron y sigue igual.
34. La portabilidad de mi número no se completó y estoy sin línea desde hace varias horas.
35. El internet de mi plan Familiar 549 no funciona desde octubre y no hay quien lo repare. Sigo pagando completo y no obtengo lo que corresponde.
36. Saludos. No tengo señal en toda mi calle x las mañanas ya pedi q lo revisen (ref. ticket 6477913). Ya probe todo lo q me dijeron y sigue igual.
37. Saludos. Mi Linea esta suspendida Sin Motivo Y No puedo sar Nada. ya insisti Varias Veces Y Nadie me Resuelve.
38. El técnico nunca llegó a revisar la falla de señal, ya llevo varias horas de espera.
39. Se me fue por completo la señal en Hermosillo y no he podido hacer ni recibir llamadas. Gracias.
40. Soy cliente de TelArg desde hace años. Mi chip nuevo no activa, llevo una semana esperando que me den cobertura. Espero su pronta respuesta.
41. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir.
42. No tengo servicio de mensajes en mi colonia en fines de semana, ya pedi q lo revisen (ref. caso 30497). Muchas gracias x su atención.
43. Solo me aparece emergencias en la pantalla, no puedo usar mi línea toda la tarde. Ya me había pasado antes este mes.
44. El internet de mi plan Navega 249 no funciona desde marzo y no hay quien lo repare. Quedo al pendiente.
45. El módem se queda con la luz roja y no hay internet, llevo una semana igual. Ya probé todo lo que me dijeron y sigue igual. Espero una respuesta pronto.
46. Buenos días. Mi chip nuevo no activa, llevo cuatro días esperando que me den servicio de llamadas.
47. Buenos dias. No me llegan los codigos de verificacion x SMS desde hace cuatro días, creo que hay una falla en la red.
48. Soy cliente de TelArg desde hace años. Avisé de la falla en la sucursal con reporte 87554 y todavía no me restablecen el servicio.
49. Ya cambié mi SIM y sigo sin señal, el problema es de la red. Esto afecta a toda mi familia.
50. buenas tardes. no tengo datos moviles aunque mi saldo es suficiente (ref. folio 549447). ya insisti varias vecs y nadie me resuelve.

### Cobros Indebidos

1. Pagué mi mensualidad completa por teléfono y me volvieron a cobrar $399 MXN. Nunca me avisaron de este cargo.
2. Me quitaron el saldo de mi recarga de MXN 149 sin explicación alguna. Me urge que lo corrijan cuanto antes. Muchas gracias por su atención.
3. Desconozco el cobro de 49 pesos del día 15 de julio, quiero que lo reviertan.
4. Después de cancelar mi plan Social 299 me cargaron $299 MXN a mi tarjeta.
5. Buenas noches. Me siguen cobrando la línea que di de baja, ya llevo 3 meses pagando de más. Espero una respuesta pronto.
6. Me cobraron $399.00 por seguro de equipo que nunca contraté. Exijo que lo revisen cuanto antes.
7. me cobran servicio premium como si lo hubiera pedido, ya me canse de pagar x lo q no uso. no pienso pagar un cargo q no corresponde.
8. Les escribo x un cobro q no reconozco. Activaron recarga automatica sin pedirmelo y me cobraron $149.00 en el estado de mi tarjeta. Gracias.
9. En mi tarjeta aparece un cargo de $999.00 de TelArg que desconozco. Me urge que lo corrijan cuanto antes. Por favor, avísenme cualquier novedad.
10. Cobro indebido: me descontaron $1,999 MXN dos veces el mismo día por lo mismo. Quiero que me devuelvan mi dinero. Por favor, avísenme cualquier novedad.
11. Activaron recarga automatica sin pedirmelo y me cobraron 79 pesos en el estado de mi tarjeta. No pienso pagar un cargo que no corresponde. Espero su pronta respuesta.
12. Pedí la cancelación de mi línea y aun así me cobraron $549 MXN en abril. Quiero una explicación clara de este cargo.
13. Buenos dias. Me quitaron el saldo de mi recarga de $149 MXN sin explicacion alguna...
14. Desconozco el cobro de $79.50 del dia 28 de septiembre, quiero q lo reviertan...
15. Me activaron servicio premium sin avisarme y ahora me cobran 169 pesos.
16. Me activaron seguro de equipo sin avisarme y ahora me cobran $249.50. Quiero una explicacion clara de este cargo. Quedo al pendiente.
17. Se me cobro una penalizacion de $49.50 q no corresponde porque ya cumpli mi contrato. Quiero una explicacion clara de este cwrgo.
18. Me cobraron MXN 199 por una promoción que nunca acepté (reporte 984866). Ya pedí una aclaración y sigo sin respuesta.
19. desconozco el Cobro de $79.00 del Dia 5 De agosto, Quiero Que lo reviertan.
20. Desconozco el cobro de MXN 899 del dia 5 de enero, quiero q lo reviertan. Quedo al pendiente...
21. Mi banco me reporta un cargo no reconocido de TelArg por $299.00, quiero aclararlo.
22. Pedí la devolución de $399 MXN hace cuatro días y todavía no me abonan nada.
23. Pagué mi mensualidad completa en la app y me volvieron a cobrar $299 MXN.
24. Ya reclame el cbro de reconexion por correo y nadie me devuelve mi dinero (caso 68166). Es muy frustrante q pase esto con mi dinero.
25. ACTIVARON RECARGA AUTOMATICA SIN PEDÍRMELO Y ME COBRARON 129 PESOS EN EL ESTADO DE MI TARJETA. GRAACIAS...
26. Sin mi consentimiento me dieron de alta en mensajes premium y ya van 4 cobros.
27. Me siguen cobrando la línea que di de baja, ya llevo 2 meses pagando de más.
28. Buenas tardes. En mi tarjeta aparece un cargo de $249 MXN de TelArg que desconozco. No pienso pagar un cargo que no corresponde.
29. Buenas noches. Me cobraron MXN 349 por mensajes premium que nunca contraté. Ya insistí varias veces y nadie me resuelve.
30. Desde que cambié de plan a Navega 249 me cobran $99.00 de más, es un cobro indebido. Nunca me avisaron de este cargo. Espero una respuesta pronto.
31. soy cliente de telarg desde hace años. q onda, me robaron 349 pesos de saldo en un solo dia sin hacer llamadas ni usar datos.
32. desconozco el cobro de $169.50 del día 10 de mayo, quiero que lo reviertan.
33. En mi tarjeta aparece un cargo de $699 de TelArg que desconozco. Saludos cordiales.
34. Me cobran renta de equipo como si lo hubiera pedido, ya me cansé de pagar por lo que no uso. No pienso pagar un cargo que no corresponde.
35. Activaron servicio premium sin pedírmelo y me cobraron $199 MXN en el estado de mi tarjeta. No pienso pagar un cargo que no corresponde.
36. Me cobraron $499.50 por seguro de equipo que nunca contraté. Espero una respuesta pronto.
37. SALUDOS. NO MANCHES, APARECE UN CARGO DE $229 POR RECARGA AUTOMATICA Y YO NO AUTORICÉ NAADA. POR FAVOR, AVISENME CUALQUIER NOVEDAD.
38. Me cobraron 299 pesos por una promoción que nunca acepté (folio 376072). Ya pedí una aclaración y sigo sin respuesta.
39. Me descontaron $149.50 de mi saldo sin haber usado nada (#3328019). Por favor, avísenme cualquier novedad.
40. En mi tarjeta aparece un cargo de $549.50 de TelArg que desconozco. Ya revisé mi cuenta y no hay explicación alguna.
41. Buenas tardes. Mi banco me reporta un cargo no reconocido de TelArg por $499, quiero aclararlo.
42. Me activaron contenido de terceros sin avisarme y ahora me cobran 79 pesos. x favor, avisenme cualquier novedad.
43. Me cobraron roaming internacional de $1,199 MXN sin salir de Veracruz. Me urge q lo corrijan cuanto antes.
44. Saludos. Activaron contenido de terceros sin pedírmelo y me cobraron $49 en el estado de mi tarjeta. Ya pedí una aclaración y sigo sin respuesta.
45. En mi tarjeta aparece un cargo de $229 MXN de TelArg que desconozco. Nunca me avisaron de este cargo.
46. Les escribo por un cobro que no reconozco. Híjole, me robaron $149 de saldo en un solo día sin hacer llamadas ni usar datos. Quedo atento a su respuesta.
47. El debito automatico me quitó $899.00 dos veces en enero, necesito el reembolso!!!
48. Desconozco el cobro de $129 MXN del día 5 de mayo, quiero que lo reviertan.
49. Me aplicaron un cargo x reconexion de $349 cuando mi linea nunca estvuo suspendida. Espero una respuesta pronto.
50. Me descontaron $899 MXN de mi saldo sin haber usado nada (ticket 56830). Gracias.

### Problemas de Facturación

1. El portal muestra un monto de factura distinto al del recibo que me llegó, $999.00 contra $699.00. Espero su pronta respuesta.
2. Buenos dias. Pedí la fact de mi equipo y me mandaron la de otro número de línea. No me gusta quedarme con dudas de lo q pago!!!
3. Cambié de plan a Joven 149 y la factura sigue mostrando la tarifa anterior. Ya insistí varias veces y nadie me resuelve. Espero su pronta respuesta.
4. Solicité factura electrónica por correo electrónico y me siguen enviando el recibo impreso. Quiero que me confirmen cuándo estará corregido.
5. Mi fecha de corte cambió al día 10 sin avisarme y ahora el estado de cuenta viene desfasado. Ya intenté arreglarlo por mi cuenta y no pude.
6. Tengo la línea a mi nombre. Mi plan Nocturno 79 cuesta $79 MXN, pero en mi factura el monto total es de $129.
7. No puedo descargar mi factura en el portal web, marca error cada vez.
8. Les escribo porque necesito ayuda. Solicite factura electronica por corrreo electronico y me siguen enviando el recibo impreso. Necsito que esto quede aclarado con urgencia. Espero una respuesta pronto...
9. La factura de abril viene con el plan equivocado, aparece Ahorro 99 y yo contraté otro. Es muy molesto andar ackarando esto. Espero una respuesta pronto.
10. Mi recibo no muestra el descuento prometido en el call center, solo el total de $499.00. Ya intenté arreglarlo por mi cuenta y no pude.
11. Mi factura de noviembre dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente. Gracias de antemano.
12. El ciclo de facturación cambió y ahora tengo dos recibos casi juntos de $2,499 MXN cada uno.
13. necesito el comprobante fiscal (CFDI) de agosto y el sistema no me lo genera (TKT-5310885). Es muy molesto andar aclarando esto.
14. La fecha de vencimiento de mi recibo no es la que me dieron, por eso se me juntan los periodos. Ya intenté arreglarlo por mi cuenta y no pude.
15. La factura trae un saldo anterior que ya pagué, el pago de $229.50 no aparece reflejado. Espero su pronta respuesta.
16. MI FECHA DE CORTE CAMBIO AL DIA 15 SIN AVISARME Y AHORA EL ESTADO DE CUENTA VIENE DESFASADO. QUEDO AL PENDIENTE.
17. Les escribo por un tema de facturación. Cambié de plan a Total 399 y la factura sigue mostrando la tarifa anterior. Nadie me da una explicación clara.
18. Hola. El estado de cuenta de agosto no detalla los conceptos, solo trae un total de $599.50. Muchas gracias por su atención.
19. Les escribo porque necesito ayuda. O sea, mi estado de cuenta no coincide con lo que me dijeron en la sucursal. Nadie me da una explicación clara.
20. La factura de octubre viene con el plan equivocado, aparece Premium 999 y yo contraté otro. Necesito que lo resuelvan cuanto antes. Saludos cordiales.
21. Mi fact llega con los datos de mi domicilio anterior aunque ya lo actualicé.
22. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos.
23. Hola. La factura de septiembre viene con el plan equivocado, aparece Amigo Sin Límite y yo contraté otro.
24. Los consumos en mi factura de diciembre no corresponden con mi historial de llamadas.
25. Quiero cambiar mi factura a entrega por correo electrónico y no encuentro la opción por ningún lado. Ya me había pasado antes y nadie lo arregló de fondo.
26. la Fecha de Vencimiento de mi recibo No Es La Que Me dieron, x eso se Me juntan los Periodos...
27. Mi factura sale a nombre de otra persona y necesito que corrijan los datos.
28. me emitieron dos facturas por el mismo periodo de marzo, una de ellas está duplicada. Nadie me da una explicación clara.
29. Buenas tardes. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. Ya intenté arreglarlo por mi cuenta y no pude.
30. No puedo registrar mi correo para facturación, la página me regresa un error (reporte 442051). Ya insistí varias veces y nadie me resuelve.
31. No aparece el pago que hice en la app en mi estado de cuenta, sigo con saldo pendiente.
32. El estado de cuenta de julio no detalla los conceptos, solo trae un total de $299.50. Necesito que esto quede aclarado con urgencia. Gracias de antemano.
33. Les escribo porque necesito ayuda. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. No me gusta quedarme con dudas de lo que pago.
34. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Es muy molesto andar aclarando esto. Quedo al pendiente.
35. El ciclo de fact cambio y ahora ttengo dos recibos casi juntos de $1,499 cada uno. Quedo atento a su respuesta.
36. Pedí la factura de mi equipo y me mandaron la de otro número de línea. Es muy molesto andar aclarando esto. Muchas gracias por su atención.
37. Hola. Mi recibo viene con fecha límite de pago distinta a la que me indicaron por teléfono. Quiero que me confirmen cuándo estará corregido.
38. Mi plan Empresarial 899 cuesta $899 MXN, pero en mi factura el monto total es de $1,199.00. Gracias.
39. No puedo registrar mi correo para facturación, la página me regresa un error (TKT-43672). Ya intenté arreglarlo por mi cuenta y no pude.
40. Les escribo por un tema de facturación. Pedí la factura de mi equipo y me mandaron la de otro número de línea.
41. Mi plan Empresarial 899 cuesta $899 MXN, pero en mi factura el monto total es de $1,499.50. Nadie me da una explicación clara. Quedo atento a su respuesta.
42. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Necesito que lo resuelvan cuanto antes.
43. Solicité fact electronica en la app y me siguen enviando el recibo impreso
44. Quiero cambiar mi factura a entrega por correo electrónico y no encuentro la opción por ningún lado. Quiero que me confirmen cuándo estará corregido. Por favor, avísenme cualquier novedad.
45. El portal muestra un monto de factura distinto al del recibo que me llegó, 449 pesos contra $549 MXN.
46. La factura de septiembre viene con el plan equivocado, aparece Hogar 449 y yo contrate otro.
47. No puedo registrar mi correo para facturación, la página me regresa un error (folio 27586).
48. El portal muestra un monto de factura distinto al del recibo que me llegó, $899 MXN contra $799.50. Espero una respuesta pronto.
49. El estado de cuenta de enero no detalla los conceptos, solo trae un total de MXN 49. Gracias de antemano.
50. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos en la sucursal.

### Quejas

1. Mala experiencia por correo: tuve que repetir mi problema a cuatro personas distintas!!!
2. Ay no, el trato de Jorge en la sucursal fue grosero, no pienso aguantar eso...
3. Mi queja es por la falta de interés en el área de cobranza, nunca devuelven las llamadas.
4. La neta, no es posible que una empresa tan grande tenga una atención tan lenta. Ya insistí varias veces y nadie me resuelve.
5. Les escribo desde Hermosillo. No hay quien resuelva dudas en el call center, todo son excusas y tiempos de espera. Espero su pronta respuesta.
6. Hola. Para cancelar mi linea me ponen mil trabas, es un proceso interminable.
7. A quien corresponda. El proceso para cambiar de plan a Prsmium 999 fue lentisimo y lleno de requisitos absurdos.
8. La sucursal cierra antes del horario que anuncian y ya van 5 veces que voy en balde. Ya insistí varias veces y nadie me resuelve.
9. Los asesores no saben nada, me dieron 4 respuestas distintas en la misma semana. Me parece inaceptable q nadie de la carra.
10. me transfirieron 3 veces entre áreas y nadie pudo resolver mi problema.
11. Quiero presentar mi queja porque Fernanda me colgó la llamada en plena aclaración. Me parece inaceptable que nadie dé la cara.
12. Buenas noces. Mala experiencia x teléfono: tuve q repetir mi problema a cuatro personas distintas. Pido que tomen en cuenta mi queja
13. Nadie da seguimiento a lo que pedí (reporte 2456444), siempre me dicen que ya está en proceso. Me parece inaceptable que nadie dé la cara.
14. El proceso para cambiar de plan a Hogar 449 fue lentísimo y lleno de requisitos absurdos. Que alguien se haga responsable, por favor. Espero que alguien por fin me responda.
15. Mi queja es por la falta de interes en soporte técnico, nunca devuelven las llamadas.
16. Tengo la línea a mi nombre. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Saludos cordiales.
17. No hay quien resuelva dudas en soporte técnico, todo son excusas y tiempos de espera.
18. En ventas me dijeron que no era su responsabilidad y me mandaron a otra area. Ya perdi demasiado tiempo en esto. Gracias.
19. Me vendieron el plan Joven 149 con condiciones q no cumplen, publicidad engañosa. Espero q alguien con autoridad me atienda.
20. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atencion. Estoy pensando seriamente en cambiarme de compañia.
21. soy Cliente de TelArg desde hace Años. Exijo Que Me Atiendan como cliente, no es La primera Vez Que me ignoran X chat. Todo Esto Es muy frustrante.
22. Buenos días. Subieron el precio de mi plan Social 299 sin avisarme, esto me parece un abuso. Estoy pensando seriamente en cambiarme de compañía.
23. Tengo una inconformidad que presentar. Híjole, no es posible que una empresa tan grande tenga una atención tan lenta. Espero una respuesta pronto.
24. Ya presenté mi inconformidad en la app y no tengo respuesta, qué falta de seriedad.
25. esperé Acsi dos horas En El Area de Cobranza Y Nadie Me Contesto. Saludos cordiales.
26. Hola. El proceso para cambiar de plan a Max 699 fue lentísimo y lleno de requisitos absurdos. Es una falta de respeto con los clientes.
27. tengo la línea a mi nombre. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención.
28. Quiero presentar mi queja porque Verónica me colgó la llamada en plena aclaración. Ya perdí demasiado tiempo en esto. Espero que alguien por fin me responda.
29. UNA VERGUENZA LA ATENICON: ME DEJARON ESPERANDO CUARENTA MINUTOS Y SE ME CORTÓ LA LLAMADA.
30. La aplicación se cierra sola y soporte no hace nada, queja formal con ticket 2624965.
31. La atención en la sucursal de León es pésima, hacen fila veinte minutos para nada.
32. Nadie da seguimiento A lo Que pedi (csao 120038), siempre Me dicen q ya Esta En Proceso. estoy pensando seriamente En Ccambiarme De compañia
33. Los asesores no saben nada, me dieron 2 respuestas distintas en la misma semana. Ya perdi demasiado tiempo en esto!!!
34. Buenas noches. La aplicación se cierra sola y soporte no hace nada, queja formal con caso 5353945. Muchas gracias por su atención.
35. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad.
36. Les escribo porque quiero presentar una queja. Espere una hora en el call center y nadie me contesto.
37. No manches, no es posible que una empresa tan grande tenga una atención tan lenta. Es una falta de respeto con los clientes. Gracias.
38. Buenos días. Ya basta, qué mala atención, desde abril estoy tratando de aclarar un tema y nadie me resuelve. Ya perdí demasiado tiempo en esto.
39. Quiero dejar constancia de mi inconformidad. Los asesores no saben nada, me dieron 5 respuestas distintas en la misma semana. Todo esto es muy frustrante. Espero su pronta respuesta.
40. Les escrbo desde Tijuana. Me cambiaron las condiciones de mi plan Amigo Sin Limite sin pedirme autorización ya no es lo que contrate. Ya perdí demasiado tiempo en esto.
41. Chale, qué mala atención, desde julio estoy tratando de aclarar un tema y nadie me resuelve.
42. A QUIEN CORRESPONDA. Q ONDA, QUE MALA ATENCIÓN, DESDE FEBRERO ESTOY TARTANDO DE ACLARAR UN TEMA Y NADIE ME RESUELVE.
43. Subieron el precio de mi plan Amigo Sin Límite sin avisarme, esto me parece un abuso. Muchas gracias por su atención.
44. Nadie da seguimiento a lo que pedí (folio 9327387), siempre me dicen que ya está en proceso.
45. Estoy muy molesto con la forma en que Fernanda me atendio, sin ninguna empatia.
46. Me prometieron una llamada de vuelta y después de una semana sigo esperando (#425687). Ya insistí varias veces y nadie me resuelve.
47. Nadie da seguimiento a lo que pedí (ticket 5609371), siempre me dicen que ya está en proceso. Ya insistí varias veces y nadie me resuelve.
48. Espere cuarenta minutos en el área de cobranza y nadie me contestó. Espero que alguien con autoridad me atienda.
49. Mi queja: el call center cuelga cuando pido hablar con un supervisor. Espero su pronta respuesta.
50. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio. Ya insistí varias veces y nadie me resuelve.

### Reconocimientos

1. Muchas gracias por resolver mi problema por correo, el servicio fue de primera.
2. Hola. Les agradezco la atención durante julio, me ayudaron mucho con mi cambio de línea. Gracias.
3. El cambio a plan Total 399 fue sencillo y sin trámites, gracias por facilitarlo. Se los quería hacer saber. Seguiré recomendándolos.
4. Les escribo con mucho gusto. Quiero agradecer la rapidez con la q atendieron mi solicitud (caso 444134), todo quedo resuelto. Todo el equipo merece un aplauso. Gracias.
5. La promoción del plan Premium 999 es muy buena, gracias por pensar en los clientes.
6. rGacias, gracias y gracias, x fin alguien resolvió mi roblema de manera agil. Se nota la diferencia.
7. Excelente cobertura en Guadalajara, nunca he tenido problemas con mi plan Total 399. Gracias.
8. Gracias, gracias y gracias, por fin alguien resolvio mi problema de manera agil.
9. Mi experiencia en la sucursal fue de lo mejor, Héctor se tomo el tiempo de explicarme todo. Por esto sigo con TelArg.
10. En la sucursal de Toluca me atendieron de maravilla, todo fue rápido y claro. Que tengan un excelente día.
11. buenas tardes. Mi señal es excelente incluso en mi oficina, felicidades a todo el equipo de redes.
12. Mi experiencia por correo fue de lo mejor, Alejandro se tomó el tiempo de explicarme todo.
13. Soy cliente de TelArg desde hace años. Mi experiencia en la app fue de lo mejor, Carlos se tomo el tiempo de explicarme todo.
14. Quiero compartir mi experiencia. Gracias por la llamada de seguimiento, nunca me había pasado que una empresa se preocupara así.
15. El cambio a plan Joven 149 fue sencillo y sin tramites, gracias x facilitarlo...
16. Hola. Mi experiencia en el portal web fue de lo mejor, Gerardo se tomó el tiempo de explicarme todo. Saludos cordiales.
17. QUEIRO AGRADECER LA RAPIDEZ CON LA Q ATENDIERON MI SOLICITUD (TICKET 34626), TODO QUEDÓ RESUELTO. SE NOTA LA DIFERENCIA. SEGUIRE RECOMENDANDOLOS.
18. TENGO LA LINEA A MI NOMBRE. LAS MEJORES PROMOCIONES DEL MERCADO, SIIGO CON TELARG X EL BUEN SERVICIO.
19. Quiero dejar un comentario positivo. Quiero dejar constancia de lo bien que me atendió Verónica, un trato impecable.
20. El cambio a plan Premium 999 fue sencillo y sin trámites, gracias por facilitarlo. Todo el equipo merece un aplauso. Saludos cordiales.
21. Quiero felicitar a Ricardo en atención a clientes, el trato fue de diez.
22. Quiero dejar un comentario positivo. Un diez para el call center: resolvieron mi duda sin darme vueltas.
23. Verónica me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Sigan así, de verdad se agradece.
24. Felicito a TelArg por la claridad de mi factura, ahora sí entiendo cada concepto. Voy a seguir con ustedes muchos años más. Gracias.
25. quiero Reconocer Publicamente la Labor de gerardo en la sucursal un trabajo verdaderamente Profesional.
26. Felicito a TelArg por la claridad de mi factura, ahora sí entiendo cada concepto.
27. Hola. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. Todo el equipo merece un aplauso.
28. Buenos días. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad. Me dio mucho gusto ser su cliente. Saludos cordiales.
29. El cambio a plan Nocturno 79 fue sencillo y sin trámites, gracias por facilitarlo.
30. Agradezco que me hayan devuelto mi dinero tan rápido, excelente respuesta.
31. Les escribo desde Cancún. La promoción del plan Social 299 es muy buena, gracias por pensar en los clientes.
32. En la sucursal de Culiacán me atendieron de maravilla, todo fue rápido y claro. Sigan así.
33. Gracias a Hector x la paciencia y la amabilidad al explicarme mi plan Familiar 549. De verdad quede muy satisfecho. Un saludo a todo el equipo.
34. A quien corresponda. Ricardo me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Seguiré recomendándolos.
35. Quiero dejar constancia de lo bien que me atendió Carlos, un trato impecable. Seguiré recomendándolos.
36. Excelente cobertura en Toluca, nunca he tenido problemas con mi plan Nocturno 79.
37. Gracias por la llamada de seguimiento, nunca me había pasado que una empresa se preocupara así. Definitivamente los voy a recomendar con mi familia.
38. Muy contento con el internet de mi plan Premium 999, la velocidad es increíble. Sigan así.
39. Les escribo con una buena noticia. Quiero felicitar a Gerardo en ventas, el trato fue de diez.
40. Quiero dejar constancia de lo bien que me atendió Verónica, un trato impecable.
41. En la sucursal de Hermosillo me atendieron de maravilla, todo fue rápido y claro.
42. Les escribo desde Guadalajara. Mi experiencia x chat fue de lo mejor, Laura se tomo el tiempo de explicarme todo.
43. Las mejores promociones del mercado, sigo con TelArg x el buen servicio.
44. Quiero dejar constancia de lo bien que me atendió Héctor, un trato impecable.
45. Quiero agradecer la rapidez con la que atendieron mi solicitud (caso 828580), todo quedó resuelto.
46. Muchas gracias por resolver mi problema en el portal web, el servicio fue sobresaliente. Seguiré recomendándolos.
47. Mi reconocimiento al tecnico q llegó puhtual y dejo mi internet funcionahdo perfecto. Sigan así.
48. QUIERO AGRADECER LA RAPIDEZ CON LA Q ATENDIERON MI SOLICITUD (CASO 2383031), TODO QUESO RESUELTO.
49. Hola. Mi experiencia en la app fue de lo mejor, Daniela se tomó el tiempo de explicarme todo. Todo el equipo merece un aplauso.
50. En la sucursal de Culiacán me atendieron de maravilla, todo fue rápido y claro.
