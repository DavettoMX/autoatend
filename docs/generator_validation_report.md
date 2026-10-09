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
| noise rate | PASS | realized 0.2994, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | pairwise Jaccard of top-100 tokens in [0.220, 0.342], required within [0.15, 0.80] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 16.25% duplicate texts; text length 40 / 127.9 / 278 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.258 | 0.266 | 0.307 | 0.266 |
| Cobros Indebidos | 0.258 | 1.000 | 0.266 | 0.274 | 0.220 |
| Problemas de Facturación | 0.266 | 0.266 | 1.000 | 0.235 | 0.220 |
| Quejas | 0.307 | 0.274 | 0.235 | 1.000 | 0.342 |
| Reconocimientos | 0.266 | 0.220 | 0.220 | 0.342 | 1.000 |

## Noise calibration

| Parameter | Value |
|---|---:|
| Comment selection rate | 0.30 |
| Tilde drop (per accented letter) | 0.80 |
| Abbreviation q/x/tel/fact (per word) | 0.50 |
| Keyboard typo (per word of 4+ letters) | 0.04 |
| Case jitter (per comment) | 0.25 |
| Irregular punctuation (per comment) | 0.30 |
| **Realized selection rate** | **0.2994** (tolerance ±0.015) |

## Persisted corpus verification

Measured by re-opening the finished file read-only.

| Category | Rows | Train | Test | max(train.fecha) | min(test.fecha) |
|---|---:|---:|---:|---|---|
| Falta de Servicio | 21900 | 17520 | 4380 | 2026-04-03 | 2026-04-03 |
| Cobros Indebidos | 14600 | 11680 | 2920 | 2026-04-06 | 2026-04-06 |
| Problemas de Facturación | 14600 | 11680 | 2920 | 2026-04-05 | 2026-04-05 |
| Quejas | 14600 | 11680 | 2920 | 2026-04-08 | 2026-04-08 |
| Reconocimientos | 7300 | 5840 | 1460 | 2026-04-10 | 2026-04-10 |

Empty texts: 0. Split rule: per category, ordered by (fecha, id), oldest 80% train, newest remainder test.

corpus.db SHA-256: `c0343583040f51afa54cfddb5bebe608bac3a02aa1416286582ae97744adfb05`

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. Hice mi recarga de $129.50 hace quince días y sigo sin datos móviles. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir.
2. mi línae sigue sin servicio tras la reconexion q pague (ref. reporte 3569964). quedo atento a su respuesta.
3. Hice mi recarga de $1,499.00 hace cuatro días y sigo sin cobertura. Espero su pronta respuesta.
4. el modem se queda con la luz roja y no hay internet, llevo quince dias igual. no tengo datos moviles aunque mi saldo es suficiente (ref. ticket 3282927).
5. Desde que cambiaron la antena en mi casa perdí señal y no lo he recuperado. Me quedé sin cobertura justo cuando más lo necesitaba, ¿hasta cuándo van a arreglarlo?
6. Les escribo desde San Luis Potosí. El router dice sin servicio desde hace una semana, ya reinicié todo y nada. Activé el roaming para viajar y al llegar a Hermosillo no tuve señal.
7. El modem se queda con la luz roja y no hay internet, llevo varias horas igual.
8. Hice mi recarga de $49.50 hace diez días y sigo sin señal.
9. Activé el roaming para viajar y al llegar a Tijuana no tuve internet. Mi celular dice sin servicio desde hace cuatro días, ya reinicié todo y nada.
10. Cuando intento llamar solo escucho q el numero no esta disponible, ya van mas de 48 horas asi.
11. Hola. Hice mi recarga de $599.50 hace dos dias y sigo sin datos moviles. Se pierde servicio de mensajes en las noches en mi barrio, y ya son varios dias con el mismo problema.
12. Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace mas de 48 horas. x favor avísenme cualquier novedad.
13. Hola. Me prometieron restablecer el servicio en 7 dias y siguen sin cumplir.
14. Les escribo desde Puebla. En toda la zona no hay cobertura toda la tarde, parece que la antena está caída.
15. el internet de mi plan Libre 199 no funciona desde agosto y no hay quien lo repare.
16. Ayer hubo una caída general de la red en Saltillo y hoy todavía sigue fallando todo.
17. El internet de mi plan Amigo Sin Límite no funciona desde septiembre y no hay quien lo repare. Necesito señal para trabajar y llevo cinco días sin nada, esto ya es urgente. Espero una respuesta pronto.
18. Buenas noches. Ya cambié mi teléfono y sigo sin señal, el problema es de la red.
19. Ya cambié mi teléfono y sigo sin señal, el problema es de la red. Necesito internet para trabajar y llevo dos días sin nada, esto ya es urgente.
20. Buenos diias. Se pierde internet en las noches en toda la zona, y ya son varios dias con el mismo problema. Quedo atento a su rspuesta...
21. Saludos. Dicen q hay mantenimiento en León pero ya pasaron tres dias y seguimos sin datos moviles. x favor, avísenme cualquier novedad...
22. Les escribo porque necesito ayuda. El técnico nunca llegó a revisar la falla de señal, ya van diez días de espera.
23. La portabilidad de mi numero no se completo y estoy sin linea desde hace diez dias. Dicen que hay mantenimiento en Hermosillo perro ya pasaron una semana y seguimos sin cobertura. Es mi número de linea el q aparece en el sistema.
24. Mi chip nuevo no activa, llevo quince días esperando que me den servicio de llamadas. Espero su pronta respuesta.
25. No tengo servicio de mensajes en toda la zona en las noches, ya levanté un reporte (ref. ticket 801567).
26. Dicen que hay mantenimiento en Monterrey pero ya pasaron dos días y seguimos sin servicio de llamadas.
27. Buenas tardes. Necesito internet para trabajar y llevo tres dias sin nada, esto ya es urgente.
28. Mi línea está suspendida sin motivo y no puedo usar nada. Es mi número de línea el que aparece en el sistema.
29. Les escribo desde Villahermosa. Reporté la falla en la sucursal con reporte 4157883 y todavía no me restablecen el servicio. Gracias.
30. Cuando intento llamar solo escucho que el número no está disponible, ya van cinco días así.
31. Hice mi recarga de MXN 699 hace varias horas y sigo sin internet.
32. Dicen que hay mantenimiento en Oaxaca pero ya pasaron dos días y seguimos sin servicio de llamadas. Pagué mi plan Social 299 puntual y aun así me suspendieron el servicio.
33. Ayer hubo una caída general de la red en Querétaro y hoy todavía sigue fallando todo.
34. En mi casa no hay cobertura por las mañanas, parece que la antena está caída. Espero una respuesta pronto.
35. Desde hace una semana estoy sin servicio de llamadas en mi colonia y nadie me da una solución.
36. No tengo internet en mi trabajo toda la tarde, ya levanté un reporte (ref. caso 9576037).
37. LES ESCRIBO DESDE PUEBLA. ME QUEDE SIN SERVICIO DE LLAMADAS JUSTO UCANDO MAS LO NECESITABA, ¿HASTA CUÁNDO VAN A ARREGLARLO? GRACIAS.
38. Mi línea sigue sin servicio tras la reconexión que pagué (ref. ticket 5189648). Es mi número de línea el que aparece en el sistema.
39. me prometireon restablecer el servicio en 3 dias y siguen sin cumplir...
40. El módem se queda con la luz roja y no hay internet, llevo diez días igual.
41. Buenos días. No me llegan los códigos de verificación por SMS desde hace varias horas, creo que hay una falla en la red. Desde que cambiaron la antena en toda mi calle perdí servicio de llamadas y no lo he recuperado.
42. No me llegan los códigos de verificación por SMS desde hace tres días, creo que hay una falla en la red. Gracias de antemano.
43. La verdad, el internet está tan lento en horario de oficina que es como no tener servicio de llamadas. Me quedé sin servicio de llamadas justo cuando más lo necesitaba, ¿hasta cuándo van a arreglarlo?
44. El técnico nunca llegó a revisar la falla de señal, ya van más de 48 horas de espera. Caray, el internet está tan lento por las mañanas que es como no tener internet.
45. No tengo señal en mi colonia toda la tarde, ya levanté un reporte (ref. #12948). Solo me aparece emergencias en la pantalla, no puedo usar mi línea en fines de semana.
46. Les escribo desde Tijuana. Cambié de SIM y quedé incomunicado, sin servicio de llamadas desde hace quince días. Por favor, avísenme cualquier novedad.
47. No puedo mandar ni recibir mensajes de texto en Saltillo desde hace más de 48 horas.
48. No tengo datos móviles aunque mi saldo es suficiente (ref. reporte 2870089).
49. Les escribo desde Monterrey. Se me fue por completo la señal en Chihuahua y no he podido hacer ni recibir llamadas.
50. La verdad, llevo varias horas sin señal y sigo pagando mi plan Nocturno 79. Dicen que hay mantenimiento en Tijuana pero ya pasaron varias horas y seguimos sin servicio de mensajes.

### Cobros Indebidos

1. Quiero que me reembolsen los $399 que me cobraron de más (ticket 70362). Activaron renta de equipo sin pedírmelo y me cobraron MXN 149 en el estado de mi tarjeta.
2. El débito automático me quitó $2,499 MXN dos veces en octubre, necesito el reembolso. Gracias.
3. Quiero que me reembolsen los $129 que me cobraron de más (reporte 2613641).
4. Mi banco me reporta un cargo no reconocido de TelArg por $149.50, quiero aclararlo.
5. el debito automatico me quitó $199 dos veces en junio, necesito el reembolso.
6. Saludos. Pedí la cancelación de mi línea y aun así me cobraron $1,499.50 en septiembre.
7. Cobran 79 pesos por plan adicional y nunca me avisaron q se activaba.
8. Mi plan Libre 199 cuesta $149 MXN y me estan cobrando $599.50, no es lo que acordamos...
9. me descontaron 449 pesos de mi saldo sin haber usado nada (ticket 68518). por favor, avisenme cualquier novedad...
10. hice la recarga de $349.50 y me cobraron una comision q no me explicaron. me hicieron un cobro doble de $79 mxn en abril, exijo la devolución
11. Me cobran 1,199 pesos de más desde que cambié de plan a Navega 249, es un cobro indebido. Gracias.
12. Me suscribieron a servicio premium sin mi consentimiento y ya van 3 cobros. Saludos cordiales.
13. Cobro indebido: me descontaron $169 y 999 pesos el mismo dia x lo mismo...
14. Activaron renta de equipo sin pedírmelo y me cobraron 249 pesos en el estado de mi tarjeta. Me cobraron 349 pesos por plan adicional que nunca contraté.
15. Pagué mi mensualidad completa y me volvieron a cobrar $49.00 por chat. Por favor, avísenme cualquier novedad.
16. Buenos días. Cobran $199 MXN por roaming internacional y nunca me avisaron que se activaba. Mi banco me reporta un cargo no reconocido de TelArg por $149.00, quiero aclararlo.
17. Ya reclamé el cobro de renta de equipo por correo y nadie me devuelve mi dinero (TKT-88848). Me cobraron MXN 169 por plan adicional que nunca contraté.
18. Activaron membresía de juegos sin pedírmelo y me cobraron 899 pesos en el estado de mi tarjeta. Me hicieron un cobro doble de MXN 99 en abril, exijo la devolución.
19. Me aplicaron un cargo por reconexión de $229.50 cuando mi línea nunca estuvo suspendida.
20. Me cobraron $79 por una promoción que nunca acepté (#181138). Me descontaron $1,999 de mi saldo sin haber usado nada (folio 6108588).
21. Cancelé paquete de datos extra hace meses y me siguen descontando $129 cada mes. Me cobran $169 MXN de más desde que cambié de plan a Ahorro 99, es un cobro indebido.
22. Saludos. Chale, me robaron $449 de saldo en un solo día sin hacer llamadas ni usar datos.
23. Me cobran plan adicional como si lo hubiera pedido, ya me cansé de pagar por lo que no uso.
24. Buenas noches. Desconozco el cobro de $549.00 del día 15 de marzo, quiero que lo reviertan. Cobran $1,999.00 por servicio premium y nunca me avisaron que se activaba.
25. Hola. Pedí la cancelación de mi línea y aun así me cobraron $1,499 MXN en mayo. Muchas gracias por su atención.
26. Mi plan Joven 149 cuesta 699 pesos y me estan cobrando $399.50 no es lo q acordamos. Cargos que no reconozco en enero: MXN 699 por suscripcion de musica. Exijo q los eliminen.
27. Les escribo desde Cancún. Pagué mi mensualidad completa y me volvieron a cobrar $169 en la app. Gracias de antemano.
28. Desconozco el cobro de MXN 249 del día 10 de febrero, quiero que lo reviertan. Es mi número de línea el que aparece en el sistema.
29. Tengo la líjea a mi nombre. Me hicieron un cobro doble de $129.50 en febrero, exijo la devolución. Me quitaron el saldo de mi recarga de $79.50 sin explicacion alguna. Quedo atento a su respuesta...
30. Ya reclamé el cobro de contenido de terceros por chat y nadie me devuelve mi dinero (TKT-224585). Cargos que no reconozco en mayo: $169.00 por reconexión. Exijo que los eliminen. Saludos cordiales.
31. pedi la devolucion de $299.00 hace varias horas y todavia no me abonan nada. me cobraron $549 mxn por una promocion que nunca acepé (tkt-1577520)...
32. A quien corresponda. En mi tarjeta aparece un cargo de $499.50 de TelArg que desconozco. Espero una respuesta pronto.
33. Pagué mi mensualidad completa y me volvieron a cobrar $499.00 por correo. Me cobraron 449 pesos por una promoción que nunca acepté (TKT-7049339).
34. Pedí la devolución de $999.50 hace una semana y todavia no me abonan nada. Cargos que no reconozco en julio: $2,499.50 x plan adicional. Exijo q los eliminen. Saludos cordiales.
35. Chale, me robaron 399 pesos de saldo en un solo día sin hacer llamadas ni usar datos. Me descontaron $1,199 MXN de mi saldo sin haber usado nada (folio 399218).
36. Mi plan Premium 999 cuesta $229 y me están cobrando $699.50, no es lo que acordamos. Gracias.
37. Me cargaron MXN 1,199 a mi tarjeta despues de cancelar mi plan Total 399.
38. Pagué mi mensualidad completa y me volvieron a cobrar MXN 1,999 por chat.
39. Se me cobró una penalización de $399.00 que no corresponde porque ya cumplí mi contrato. Espero una respuesta pronto.
40. Me aplicaron un cargo por reconexión de $599 MXN cuando mi línea nunca estuvo suspendida. Por favor, avísenme cualquier novedad.
41. Me cargaron MXN 349 a mi tarjeta despies de cancelar mi plan Joven 149.
42. Me aplicaron un cargo por reconexión de $449 cuando mi línea nunca estuvo suspendida.
43. Mi plan Familiar 549 cuesta 129 pesos y me están cobrando MXN 1,999, no es lo que acordamos. Cargos que no reconozco en mayo: $1,499 por roaming internacional. Exijo que los eliminen.
44. Buenos días. Me descontaron MXN 699 de mi saldo sin haber usado nada (caso 4155474). Me cargaron MXN 129 a mi tarjeta después de cancelar mi plan Social 299. Gracias de antemano.
45. Buenas tardes. Activaron plan adicional sin pedírmelo y me cobraron $699.00 en el estado de mi tarjeta.
46. Me cobraron 799 pesos por paquete de datos extra que nunca contraté. Mi banco me reporta un cargo no reconocido de TelArg por $199, quiero aclararlo.
47. Saludos. Me cobran servicio premium como si lo hubiera pedido, ya me cansé de pagar por lo que no uso. Mi banco me reporta un cargo no reconocido de TelArg por $699.50, quiero aclararlo.
48. Ya reclame el cobro de servicio premium por tel y nadie me devuelve mi dinero (reporte 74921). Quiero que me reembolsen los $999 que me cobraron de más (caso 8531462).
49. Me siguen cobrando la línea que di de baja, ya llevo 15 meses pagando de más. Me cobraron roaming internacional de 299 pesos sin salir de Saltillo.
50. hola. me Aplicaron Un Cargo X Reconexion De $1,999.00 cuando mi lijea nunca estuvo suspendida. Mi plan hogar 449 Cuesta $999 MXN y me Estan cobrando $169.50, No Es lo que acordamos. saludos Cordiales.

### Problemas de Facturación

1. Falta el Folio fiscal en mi Factura De abril Y No Puedo deducirla. Gracias.
2. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Los consumos en mi factura de septiembre no corresponden con mi historial de llamadas.
3. El RFC de mi factura está mal capturado, necesito una factura corregida. Recibo mi factura 15 días tarde cada mes, ya casi vence cuando la tengo. Espero una respuesta pronto.
4. No me llega mi factura por ccorreo electronico desde marzo. Muchas gracias x su atencion.
5. Quiero cambiar mi factura a entrega por correo electrónico y no encuentro la opción por ningún lado. Mi factura sale a nombre de otra persona y necesito que corrijan los datos.
6. Buenas tardes. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a 499 pesos. El portal muestra un monto de factura distinto al del recibo que me llegó, $129.00 contra MXN 229.
7. Les escribo porque necesito ayuda. El estado de cuenta de marzo no detalla los conceptos, solo trae un total de 129 pesos. Es mi numero de linea el q aparece en el sistema.
8. Pedí la factura de mi equipo y me mandaron la de otro número de línea. Mi pago de 79 pesos sí se hizo pero la factura todavía marca vencido.
9. Pedí la factura de mi equipo y me mandaron la de otro número de línea. Cambié de plan a Amigo Sin Límite y la factura sigue mostrando la tarifa anterior.
10. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualice
11. Mi pago de $399 sí se hizo pero la factura todavía marca vencido. Me llegó la factura en blanco, sin importes, desde que activé la entrega en la app. Gracias de antemano.
12. Falta el folio fiscal en mi factura de febrero y no puedo deducirla.
13. saludos. el portal muestra un monto de factura distinto al del recibo que me lllego, 79 pesos contra $249 mxn. la fecha de vencimiento de mi recibo no es la q me dieron, por eso se me juntan los periodos.
14. Buenas tardes. La factura de agoosto viene con el plan equivocado, aparece Social 299 y yo contrate otro.
15. Mi factura de diciembre dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente. Mi recibo viene con fecha límite de pago distinta a la que me indicaron por chat. Espero una respuesta pronto.
16. Solicite fact electrónica por correo electrónico y me siguen enviando el recibo impreso. Quedo atento a su respuesta
17. No puedo registrar mi correo para facturación, la página me regresa un error (TKT-341466).
18. La fecha de vencimiento de mi recibo no es la que me dieron, por eso se me juntan los periodos. Espero una respuesta pronto.
19. La factura trae un saldo anterior que ya pagué, el pago de $499.00 no aparece reflejado. Quiero cambiar mi factura a entrega por WhatsApp y no encuentro la opción por ningún lado. Quedo al pendiente.
20. mi factura llega con los datos de mi domicilio anterior aunque ya lo actualice. no puedo registrar mi correo para fact, la pagina me regresa un error (caso 9210069)...
21. MI DIRECCION DE FACT ES INCORRECTA Y X ESO NO LLEGAN MIS RECIBOS. QUEDO ATENTO A SU RESPUESTA
22. Buenos días. La fecha de vencimiento de mi recibo no es la que me dieron, por eso se me juntan los periodos.
23. Hola. Los consumos en mi factura de noviembre no corresponden con mi historial de llamadas. El RFC de mi factura está mal capturado, necesito una factura corregida.
24. La factura trae un saldo anterior que ya pagué, el pago de $549 MXN no aparece reflejado. Saludos cordiales.
25. Soy cliente de TelArg desde hace años. Necesito el comprobante fiscal (CFDI) de febrero y el sistema no me lo genera (ticket 147951).
26. Solicité factura electrónica en papel y me siguen enviando el recibo impreso.
27. necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. Quedo atento a su respuesta.
28. La fact de mi linea empresarial llega a otra razon social, ya mabde los documentos por correo. No puedo descargar mi fctura en el portal web, marac error cada vez. x favor, avísenme cualquier novedad.
29. Hola. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a 99 pesos.
30. No puedo registrar mi correo para facturación, la página me regresa un error (ticket 801409). Saludos cordiales.
31. Mi factura sale a nombre de otra persona y necesito que corrijan los datos. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a $199.
32. la Fact Trae Un Saldo Anterior Q ya Pague, el Pago de $399 no aparece reflejado. eEspero su Pronta Respuesta.
33. Mi recibo viene con fecha límite de pago distinta a la que me indicaron en la sucursal.
34. les escribo desde morelia. la fact de septiembre viene con el plan equivocado, aparece joven 149 y yo contrate otro. mi recibo viene con fecha limite de pago distinta a la que me indicaron x chat. espero una respuesta pronto.
35. BUENAS TARDES. EL RFC DE MI FACTURA ESTA MAL CAPTURADO, NECESITO UNA FACT CORREGIDA. NO PUEDO REGISTRAR MI ORREO PARA FACT, LA PAGINA ME REGRESA UN ERROR (TKT-424701).
36. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. Recibo mi factura 30 días tarde cada mes, ya casi vence cuando la tengo.
37. Quiero cambiar mi factura a entrega por WhatsApp y no encuentro la opción por ningún lado. Muchas gracias por su atención.
38. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. En mi factura el monto total es de $449.50 pero mi plan Empresarial 899 debería ser distinto.
39. EL CICLO DE FACTURACION CAMBIO Y AHORA TENGO DOS RECIBOS CASI JUNTOS, UNO DE $1,999
40. El estado de cuenta de octubre no detalla los conceptos, solo trae un total de 2,499 pesos. Solicite factura electronica x correo electronico y me siguen enviando el recibo impreso.
41. Mi pago de $1,499.50 sí se hizo pero la factura todavía marca vencido.
42. El RFC de mi factura está mal capturado, necesito una factura corregida. Mi recibo no muestra el descuento prometido en el área de cobranza, solo el total de $399.00.
43. El portal muestra un monto de fact distinto al del recibo que me llego, MXN 149 contra $49.00.
44. El RFC de mi factura está mal capturado, necesito una factura corregida. Espero una respuesta pronto.
45. Mi direccion de ffacturación es incorrecta y x eso no llegan mis recibos...
46. Buenas tardes. Recibo mi factura 5 días tarde cada mes, ya casi vence cuando la tengo. El ciclo de facturación cambió y ahora tengo dos recibos casi juntos, uno de $449 MXN.
47. Cambié de plan a Premium 999 y la factura sigue mostrando la tarifa anterior. No puedo registrar mi correo para facturación, la página me regresa un error (TKT-40588).
48. Falta el folio fiscal en mi factura de noviembre y no puedo deducirla. Quiero cambiar mi factura a entrega en papel y no encuentro la opción por ningún lado.
49. La factura trae un saldo anterior que ya pagué, el pago de $249 no aparece reflejado.
50. hola. solicite fact electronica en la app y me siguen enviando el recibo impreso.

### Quejas

1. Ya presenté mi inconformidad en la app y no tengo respuesta, qué falta de seriedad. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por chat.
2. Les escribo porque necesito ayuda. Me prometieron una llamada de vuelta y después de diez días sigo esperando (#7096680). Me vendieron el plan Social 299 con condiciones que no cumplen, publicidad engañosa.
3. Mi queja: el call center cuelga cuando pido hablar con un supervisor. Qué barbaridad, el trato de Alejandro en atención a clientes fue grosero, no pienso aguantar eso.
4. En atención a clientes me dijeron que no era su responsabilidad y me mandaron a otra área. El chat de TelArg es inútil, solo responden con mensajes automáticos.
5. Exijo que me atiendan como cliente, no es la primera vez que me ignoran en la app.
6. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. Espero su pronta respuesta.
7. Mi queja es por la falta de interés en atención a clientes, nunca devuelven las llamadas.
8. Esperé cuarenta minutos en el área de cobranza y nadie me contestó. Quedo al pendiente.
9. El chat de TelArg es inútil, solo responden con mensajes automáticos. Qué onda, qué mala atención, desde agosto estoy tratando de aclarar un tema y nadie me resuelve.
10. Mala experiencia x correo: tuve q repetir mi problema a cuatro personas distintas. Las promociones q anuncian no aplican nunca para clientes que ya tenemos años con ustedes!!!
11. Esperé una hora en soporte técnico y nadie me contestó.
12. No manches, el trato de Laura en el área de cobranza fue grosero, no peinso aguantar eso. Me vendieron el plan Hogar 449 con condiciones q no cumplen, publicidad engañosa.
13. Buenas tardes. Me vendieron el plan Hogar 449 con condiciones que no cumplen, publicidad engañosa.
14. Para cancelar mi linea me pnen mil trabas, es un proceso interminable. No hay quien resuelva dudas en ventas, todo son excusas y tiempos de espera.
15. ME PROMETIERON UNA LLAMADA DE VUELTA Y DESPUES DE DOS DIAS SIGO ESPERANDO (REPORTE 960624)...
16. esperé media hora en soporte técnico y nadie me contesto.
17. Me cambiaron las condiciones de mi plan Empresarial 899 sin pedirme autorizacion ya no es lo q contraté. Saludos cordiales.
18. Buenos días. El chat de TelArg es inútil, solo responden con mensajes automáticos.
19. El chat de TelArg es inútil, solo responden con mensajes automáticos. Me vendieron el plan Max 699 con condiciones que no cumplen, publicidad engañosa.
20. Mala experiencia por correo: tuve que repetir mi problema a cuatro personas distintas. Quedo atento a su respuesta.
21. Caray, no es posible que una empresa tan grande tenga una atención tan lenta.
22. Ándale, no es posible que una empresa tan grande tenga una atención tan lenta. Muchas gracias por su atención.
23. Ya presenté mi inconformidad en la sucursal y no tengo respuesta, qué falta de seriedad.
24. Me prometieron una llamada de vuelta y después de más de 48 horas sigo esperando (reporte 122991).
25. Saludos. La atención en la sucursal de Querétaro es pésima, hacen fila más de una hora para nada. Exijo que me atiendan como cliente, no es la primera vez que me ignoran en el portal web. Gracias.
26. El chat de TelArg es inútil, solo responden con mensajes automáticos.
27. Los asesores no saben nada, me dieron 5 respuestas distintas en la misma semana.
28. La atención en la sucursal de Querétaro es pésima, hacen fila una hora para nada. Espero su pronta respuesta.
29. Les escribo porque necesito ayuda. La atención en la sucursal de Mérida es pésima, hacen fila cuarenta minutos para nada. Una vergüenza la atención: me dejaron esperando veinte minutos y se me cortó la llamada.
30. Les escribo porque necesito ayuda. Mala experiencia en la sucursal: tuve que repetir mi problema a cuatro personas distintas.
31. Cada vez Que Llamo Me Piden los mismos datos y No resuelven Nada, estoy Harto Del Servicio De Atención. Las condiciones del Contrato cambian cada Mes Y nadie Me las explica con Claridad.
32. Tengo la línea a mi nombre. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes.
33. Saludos. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio. El proceso para cambiar de plan a Total 399 fue lentísimo y lleno de requisitos absurdos.
34. me prometieron una llamada de vuelta y despues de una semaan sigo esperando (reporte 87712).
35. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio.
36. Ya basta, que mala atencion, desde octubre estoy tratando de aclarar un tema y nadie me resuelve. Espero una respuesta pronto.
37. El chat de TelArg es inútil, solo responden con mensajes automáticos. Me prometieron una llamada de vuelta y después de dos días sigo esperando (folio 2817080).
38. Me llaman seguido toda la tarde para ofrecerme promociones aunque ya les pedí que no. Quedo atento a su respuesta.
39. La sucursal cierra antes del horario que anuncian y ya van 15 veces que voy en balde. Me prometieron una llamada de vuelta y después de varias horas sigo esperando (ticket 89719).
40. La sucursal cierra qntes del horario que anunciaan y ya van 15 veces q voy en balde. Quiero presentar mi queja porque Ferannda me colgo la llamada en plena aclaracion. Quedo atento a su respuesta...
41. Me vendieron el plan Total 399 con condiciones que no cuumplen, publicidad engañosa. x favor, avísenme cualquier novedad
42. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. Muchas gracias por su atención.
43. A quien corresponda. Caray, el trato de Ricardo en la sucursal fue grosero, no pienso aguantar eso. Una vergüenza la atención: me dejaron esperando veinte minutos y se me cortó la llamada.
44. A quien corresponda. La sucursal cierra antes del horario q anuncian y ya van 30 veces q voy en balde.
45. Exijo que me atiendan como cliente, no es la primera vez que me ignoran por chat. Es mi número de línea el que aparece en el sistema.
46. La sucursal cierra antes del horario que anuncian y ya van 7 veces que voy en balde. Es mi número de linea el q aparece en el sistema.
47. Buenas tardes. Nadie da seguimiento a mi reporte (reporte 814191), siempre me dicen que ya está en proceso.
48. Saludos. La atención en la sucursal de Cancún es pésima, hacen fila casi dos horas para nada.
49. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atencion. Gracias.
50. Subieron el precio de mi plan Total 399 sin avisarme, esto me parece un abuso.

### Reconocimientos

1. A quien corresponda. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Agradezco que me hayan devuelto mi dinero tan rápido, excelente respuesta.
2. Reconozco el esfuerzo que hacen en atención a clientes, siempre están al pendiente de los clientes. Gracias.
3. Mi experiencia x tel fue de lo mejor, Gerardo se tomo el tiempo de explicarme todo. El cambio a plan Prepago Basico fue sencillo y sin tramites, gracias x facilitarlo.
4. En la sucursal de Cancún me atendieron de maravilla, todo fue rápido y claro.
5. Gracias, gracias y gracias, por fin alguien resolvió mi problema de manera ágil.
6. El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias. La neta, qué gusto que con el área de cobranza se resolvió mi problema al primer intento.
7. Lllevo años con TelArg y siempre he tenido buen servicio, los recomiendo
8. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo. Quiero reconocer públicamente la labor de Miriam en ventas, un verdadero profesional. Espero una respuesta pronto.
9. No manches, qué buen servicio, en menos de más de una hora me solucionaron todo. Mi señal es excelente incluso en el centro, felicidades a todo el equipo de redes. Gracias de antemano.
10. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. Gracias por la llamada de seguimiento, nunca me había pasado que una empresa se preocupara así. Quedo al pendiente.
11. En la sucursal de Morelia me atendieron de maravilla, todo fue rápido y claro. Quiero felicitar a Daniela en atención a clientes, la atención fue de lo mejor.
12. Un diez para atención a clientes: resolvieron mi duda por chat sin darme vueltas.
13. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad. Quiero felicitar a Alejandro en el call center, la atención fue impecable.
14. Soy cliente de TelArg desde hace años. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Excelente cobertura en Chihuahua, nunca he tenido problemas con mi plan Premium 999.
15. En la sucursal de Hermosillo me atendieron de maravilla, todo fue rápido y claro. Un diez para el área de cobranza: resolvieron mi duda por correo sin darme vueltas.
16. Alejandro fue muy atento y profesional, merece un reconocimiento por su trabajo. Espero una respuesta pronto.
17. Soy cliente de TelArg desde hace años. El cambio a plan Max 699 fue sencillo y sin trámites, gracias por facilitarlo. Excelente cobertura en Saltillo, nunca he tenido problemas con mi plan Libre 199.
18. gracias A alejandro x la paciencia y la amabilidad Al explicarme Mi Plan navega 249. Espero Su Pronta respueesta.
19. Felicidades por la cobertura en Querétaro, hasta en carretera tengo señal.
20. El cambio a plan Joven 149 fue sencillo y sin trámites garcias x facilitarlo. Quiiero agradecer la rapidez con la q atendieron mi reporte (caso 7952123) todo quedo resuelto.
21. Mi experiencia por teléfono fue de lo mejor, Alejandro se tomó el tiempo de explicarme todo. Quedo al pendiente.
22. Saludos. Excelente cobertura en Monterrey, nunca he tenido problemas con mi plan Familiar 549. Un diez para ventas: resolvieron mi duda en la app sin darme vueltas.
23. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes.
24. Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante. Quiero reconocer públicamente la labor de Jorge en soporte técnico, un verdadero profesional. Espero una respuesta pronto.
25. MUCHAS GACIAS X RESOLVER MI PROBLEMA POR CHAT, EL SERVICIO FUE SOBRESALIENTE.
26. Saludos. Quiero dejar constancia de lo bien que me atendio Jorge, un trato impecable. Gracias.
27. El cambio a plan Max 699 fue sencillo y sin tramites, gracias x facilitarlo. Gracias, gracias y gracias, x fin alguien resolvio mi problema de manera ágil.
28. muy contento con el internet de mi plan Hogar 449, la velocidad es increíble.
29. Reconozco el esfuerzo que hacen en la sucursal, siempre están al pendiente de los clientes. Les agradezco la atención durante diciembre, me ayudaron mucho con mi cambio de línea.
30. Tengo la linea a mi nombre. Gracias x restablecer mi servicio el mismo dia, atención de primera...
31. Les escribo porque necesito ayuda. Gracias por la llamada de seguimiento, nunca me había pasado que una empresa se preocupara así.
32. A quien corresponda. El cambio a plan Prepago Basico fue sencillo y sin tramites, gracias x facilitarlo. Es mi numero de línea el q aparece en el sistema.
33. quiero felicitar a alejandro en la sucursal, la atención fue excelente.
34. Soy cliente de TelArg desde hace años. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo. Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante. Gracias de antemano.
35. Les felicito por la aplicacion, es muy fácil de usar y puedo ver mis consumos al instante. Quiero dejr constancia de lo bien que me atendio Daniela, un trato impecable.
36. Buenos días. Felicidades por la cobertura en Guadalajara, hasta en carretera tengo señal.
37. Mi experiencia por chat fue de lo mejor, Laura se tomo el tiempo de explicarme todo. Les felicito por la aplicacion, es muy facil de usar y puedo ver mis consumos al instante. x favor, avisenme cualquier novedad.
38. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad.
39. Buenos días. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. Quiero agradecer la rapidez con la que atendieron mi reporte (TKT-5836679), todo quedó resuelto.
40. Mi señal es excepente incluso en mi casa, felicidades a todo el equipo de redes. Me encanto la atencion en el portal web, respondieron en minutos y con mucha amabilidad...
41. Me encantó la atención por correo, respondieron en minutos y con mucha amabilidad. Mi señal es excelente incluso en mi trabajo, felicidades a todo el equipo de redes.
42. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. Saludos cordiales.
43. Les escribo porque necesito ayuda. Quiero reconocer públicamente la labor de Fernanda en el área de cobranza, un verdadero profesional.
44. Agradezco q me hayan devuelto mi dinero tan rápido, excelente respuesta.
45. Mi señal es excelente incluso en mi oficina, felicidades a todo el equipo de redes. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo.
46. Muchsa gracias x resolver mi problema en el portal web, el servicio fue sobresaliente.
47. Muy contento con el internet de mi plan Navega 249, la velocidad es increible.
48. Buenos dias. O sea, que gusto que con el area de cobranza se resolvio mi problema al primer intento. Quiero agradecer la rapidez con la que atendieron mi reporte (reporte 79901), todo quedó resuelto.
49. Quiero reconocer públicamente la labor de Alejandro en el área de cobranza, un verdadero profesional.
50. Quiero reconocer públicamente la labor de Daniela en la sucursal, un verdadero profesional. Un diez para la sucursal: resolvieron mi duda por correo sin darme vueltas. Gracias de antemano.
