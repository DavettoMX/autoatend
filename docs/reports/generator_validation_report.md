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
| noise rate | PASS | realized 0.2997, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | pairwise Jaccard of top-100 tokens in [0.205, 0.351], required within [0.15, 0.80] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 20.07% duplicate texts; text length 41 / 111.3 / 234 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.325 | 0.299 | 0.299 | 0.242 |
| Cobros Indebidos | 0.325 | 1.000 | 0.351 | 0.290 | 0.227 |
| Problemas de Facturación | 0.299 | 0.351 | 1.000 | 0.235 | 0.205 |
| Quejas | 0.299 | 0.290 | 0.235 | 1.000 | 0.299 |
| Reconocimientos | 0.242 | 0.227 | 0.205 | 0.299 | 1.000 |

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
| **Realized selection rate** | **0.2997** (tolerance ±0.015) |

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

Test comments with an identical text in train: 3685 of 14600 (25.24%). A temporal split does not prevent exact repeats of a template draw.

corpus.db SHA-256: `5efce99c1fbd489829cfa78339a202ca9c33178f3b505a234de1f50824cf2b9f`

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. No puedo mandar ni recibir mensajes de texto en Chihuahua desde hace tres días.
2. mi chip nuevo no activa, llevo dos dias esperando que me den cobertura. ya perdi la cuenta de las veces que he insistido.
3. No tengo datos móviles aunque mi saldo es suficiente (ref. #96750). Nadie me da una fecha concreta de solución.
4. Hola. Mi línea no registra en la red desde hace tres días, llevo todo ese tiempo sin poder comunicarme.
5. saludos. en ttoda la zona no hay cobertura toda la tarde, parece que la antena está caída. x favor, avísenme cualquier novedad!!!
6. Hola. Pagué mi plan Joven 149 puntual y aun así me suspendieron el servicio.
7. Tengo la línea a mi nombre. Ya cambié el router y sigo sin señal, el problema es de la red.
8. Me quedé sin señal justo cuando más lo necesitaba, ¿hasta cuándo van a arreglarlo? Nadie me da una fecha concreta de solución. Muchas gracias por su atención.
9. La neta, otra vez sin internet, ya es la tercera vez este mes que se cae. Cada vez que llamo me dicen que ya está en proceso.
10. Dicen que hay mantenimiento en Monterrey pero ya pasaron una semana y seguimos sin internet. Por favor, avísenme cualquier novedad.
11. Mi línea esta suspendida sin motivo y no puedo usar nada.
12. Soy cliente de TelArg desde hace años. No tengo datos moviles aunque mi saldo es suficiente (ref. ticket 83142). x favor, avisenme cualquier novedad!!!
13. La portabilidad de mi número no se completó y estoy sin línea desde hace cuatro días. Muchas gracias por su atención.
14. Solo me aparece emergencias en la pantalla, no puedo usar mi línea en horario de oficina.
15. Hola. Se me fue x completo la señal en Querétaro y no he podido hacer ni reccibir llamadas. Ya me habia pasado antes este mes. Quedo al pendiente.
16. No me llegan los códigos de verificación por SMS desde hace diez días, creo que hay una falla en la red. Espero su pronta respuesta.
17. Mi chip nuevo no activa, llevo diez dias esperando que me den servicio de mensajes. Muchas gracias x su atencion.
18. El tecnico nunca llegó a revisar la fallla de señal, ya van cuatro días de espera. Ya me habia pasado antes este mes.
19. Hice mi recarga de MXN 399 hace una semana y sigo sin cobertura. Muchas gracias por su atención.
20. En mi casa nunca hay señal, tengo que salir a la calle para poder llamar.
21. Se me fue por completo la señal en Veracruz y no he podido hacer ni recibir llamadas. Muchas gracias por su atención.
22. Me quedé sin servicio de mensajes justo cuando más lo necesitaba, ¿hasta cuándo van a arreglarlo? Ya probé todo lo que me dijeron y sigue igual.
23. Ya cambie el modem y siyo sin señal, el problema es de la red.
24. Se me fue por completo la señal en Saltillo y no he podido hacer ni recibir llamadas. Gracias de antemano.
25. No tengo cobertura en el centro en horario de oficina, ya levanté un reporte (ref. caso 776634). Lo necesito a diario y ya me está afectando en el trabajo.
26. Hola. Cuando intento llamar solo escucho que el número no está disponible, ya van cuatro días así.
27. Cambie de SIM y quede incomunicado sin servicio de llamadas desde hace dos dias. Quedo al pendiente.
28. Mi línea no registra en la red desde hace cinco días, llevo todo ese tiempo sin poder comunicarme.
29. Tengo un problema con mi servicio. Mi linea no registra en la red desde hace diez dias, llevo todo ese tiempo sin poder comunicarme. Nadie me da una fecha concreta de solución. Muchas gracias por su atencion.
30. Soy cliente de TelArg desde hace años. Ayer hubo una caída general de la red en Tampico y hoy todavía sigue fallando todo. Ya llame tres veces y nadie me resuelve.
31. Hice mi recarga de 129 pesos hace cinco días y sigo sin cobertura. Espero una respuesta pronto.
32. El técnico nunca llegó a revisar la falla de servicio de mensajes, ya van cinco días de espera.
33. Hola. Hice mi recarga de $199 MXN hace diez días y sigo sin cobertura.
34. La verdad, llevo una semana sin servicio de llamadsa y sigo pagando mi plan Empresarial 899. Muxhas gracias por su atencion...
35. Buenas noches. Mi línea no registra en la red desde hace dos días, llevo todo ese tiempo sin poder comunicarme. Espero su pronta respuesta.
36. Mi línea sigue sin servicio tras la reconexión que pagué (ref. reporte 5361551). Esto afecta a toda mi familia. Quedo al pendiente.
37. ya basta, otra vez sin datos moviles, ya es la tercera vez este mes que se cae!!!
38. Mi celular dice sin servicio desde hace quince dias, ya reinicié todo y nada. Espero una respuesta pronto.
39. El tecnico nunca llego a revisar la falla de servicio de mensajes, ya van tres dias de espera.
40. Mi linea esta suspendida sin motivo y no puedo usar nada. Saludos cordiales.
41. Saludos. Mi línea está suspendida sin motivo y no puedo usar nada.
42. El módem se queda con la luz roja y no hay internet, llevo dos días igual.
43. Mi chip nuevo no activa, llevo cinco días esperando que me den internet. Espero su pronta respuesta.
44. Ayer hubo una caida general de la red en San Luis Potosi y hoy todavía sigue ffallando todo.
45. Mi chip nuevo no activa, llevo cinco dias esperando que me den internet.
46. Tengo un problema con mi servicio. El módem se queda con la luz roja y no hay internet, llevo una semana igual. Ya probé todo lo que me dijeron y sigue igual. Espero su pronta respuesta.
47. No tengo datos móviles aunque mi saldo es suficiente (ref. #409964). Ya llamé tres veces y nadie me resuelve.
48. Caray, llevo dos días sin servicio de mensajes y sigo pagando mi plan Total 399. Espero que lo atiendan hoy mismo. Muchas gracias por su atención.
49. Mi chip nuevo no activa, llevo dos días esperando que me den servicio de mensajes.
50. El modm se queda con la luz roja y no hay internet, llevo quince días igual.

### Cobros Indebidos

1. Les escribo por un cobro que no reconozco. Me cobran recarga automática como si lo hubiera pedido, ya me cansé de pagar por lo que no uso. Quiero una explicación clara de este cargo. Gracias de antemano.
2. Mi plan Ahorro 99 cuesta $99 y me están cobrando $129.00, no es lo que acordamos. Quiero que me devuelvan mi dinero. Gracias de antemano.
3. Pedí la devolución de $699 hace cinco días y todavía no me abonan nada. Quedo atento a su respuesta.
4. Les Escribo Porque necesito Ayuda. Me Aplicaron un Cargo X reconexion De $129.50 cuando mi línea Nunca estuvo Suspendida. Ya revise mi cuenta y no Hay explicacion Alguna.
5. Quiero que me reembolsen los $549 MXN que me cobraron de más (caso 5725858). Por favor, avísenme cualquier novedad.
6. Cargos que no reconozco en noviembre: $249.00 x renta de equipo. Exijo que los eliminen. Quiero que me devuelvan mi dinero. Por favor, avisenme cualquier novdead
7. Me cobran servicio premium como si lo hubiera pedido ya me canse de pagar x lo q no uso. No pienso pagar algo q no pedi.
8. En mi tarjeta aparece un cargo de $899.00 de TelArg que desconozco.
9. Pedi la devolucion de MXN 1,999 hace cuatro días y todavia no me abonan nada. Gracias de antemano.
10. Quiero que me reembolsen los 549 pesos que me cobraron de más (TKT-8924738).
11. Terminé de pagar el equipo hace meses y aún me cargan $999.00 mensuales. Es muy frustrante que pase esto con mi dinero.
12. Desde que cambié de plan a Empresarial 899 me cobran $49 MXN de más, es un cobro indebido.
13. Ya reclamé el cobro de reconexión en la app y nadie me devuelve mi dinero (reporte 187262). Quiero que me devuelvan mi dinero.
14. Cancelé renta de equipo hace meses y me siguen descontando MXN 99 cada mes. Quiero una explicación clara de este cargo.
15. Me cobraron roaming internacional de $1,999 MXN sin salir de Toluca. Ya pedí una aclaración y sigo sin respuesta.
16. Me activaron paquete de datos extra sin avisarme y ahora me cobran 249 pesos...
17. Cargos que no reconozco en marzo: $79 MXN por renta de equipo. Exijo que los eliminen.
18. Mi plan Familiar 549 cuesta MXN 549 y me están cobrando MXN 899, no es lo que acordamos.
19. Mi banco me reporta un Cargo No reconocido de TelArg por 2,499 pesos, quiero Aclararlo. nunca me avisaron De este Cargo. Muchas gracias x su Atencion.
20. Hola. Ya reclame el cobro de recarga automática en la app y nadie me devuelve mi dinero (ticket 226778).
21. Chale, aparece un cargo de $49.00 por membresía de juegos y yo no autoricé nada. Nadie me da una respuesta concreta. Espero una respuesta pronto.
22. Soy cliente de TelArg desde hace años. Me cobraron $2,499.50 por una promoción que nunca acepté (TKT-64993). Por favor, avísenme cualquier novedad.
23. Les escribo porque necesito ayuda. Me aplicaron un cargo por reconexión de MXN 79 cuando mi línea nunca estuvo suspendida.
24. Cobro indebido: me descontaron $1,999 dos veces el mismo día x lo mismo.
25. YA RECLAME EL COBRO DE PAQUETE DE DATOS EXTRA POR CORREO Y NADIE ME DEVUELVE MI DINERO (#514346).
26. Me suscribieron a contenido de terceros sin mi consentimiento y ya van 4 cobros.
27. El debito automatico me quito $249 MXN dos veces en septiembre, necesiito el reembolso. Es dinero q me hace falta.
28. Me siguen cobrando la línea que di de baja, ya llevo 12 meses pagando de más. Es dinero que me hace falta. Muchas gracias por su atención.
29. Saludos. Me activaron seguro de equipo sin avisarme y ahora me cobran $449. Muchas gracias por su atención.
30. Se me cobro una penalizacion de 399 pesos q no correspponde porque ya cumpli mi contrato. No pienso pagar algo q no pedi.
31. descconozco el cobro de $799.50 del dia 20 de junio, quiero que lo reviertan. quedo atento a su respuesta.
32. Les escribo desde Puebla. Mi banco me reporta un cargo no reconocido de TelArg por $499 MXN, quiero aclararlo. Ya revisé mi cuenta y no hay explicación alguna.
33. Activaron seguro de equipo sin pedírmelo y me cobraron $449.50 en el estado de mi tarjeta.
34. Hola. Cada mes se me va dinero extra por mensajes premium y no sé cómo quitarlo. Quiero una explicación clara de este cargo.
35. Activaron recarga automática sin pedírmelo y me cobraron $399 en el estado de mi tarjeta. Muchas gracias por su atención.
36. En mi tarjeta aparece un cargo de $229.00 de TelArg que desconozco. Gracias de antemano.
37. Me hicieron un cobro doble de $699 en julio, exijo la devolucion.
38. Les escribo desde Veracruz. Me quitaron el saldo de mi recarga de $79.00 sin explicación alguna.
39. Desde que cambié de plan a Hogar 449 me cobran MXN 229 de más, es un cobro indebido. Por favor, avísenme cualquier novedad.
40. Ay no, aparece un cargo de $199.50 por reconexión y yo no autoricé nada.
41. Cancelé renta de equipo hace meses y me siguen descontando $99.50 cada mes. Ya pedí una aclaración y sigo sin respuesta.
42. Me cobraron $199 por una promoción que nunca acepté (ticket 77264). Ya revisé mi cuenta y no hay explicación alguna. Gracias de antemano.
43. Me suscribieron a suscripción de música sin mi consentimiento y ya van 2 cobros. Exijo que lo revisen cuanto antes.
44. Quiero que me reembolsen los $299 q me cobraron de más (TKT-10484). Nadie me da una respuesta concreta...
45. Me activaron servicio premium sin avisarme y ahora me cobran $249.00. Quiero que me devuelvan mi dinero.
46. Me cobraron $999.00 por una promoción que nunca acepté (reporte 3778724). Quedo atento a su respuesta.
47. Pague mi mensualidad completa en la sucursal y me volvieron a cobrar 249 pesos.
48. Les escribo desde San Luis Potosí. Me descontaron 549 pesos de mi saldo sin haber usado nada (TKT-2399968). Ya llamé tres veces y nadie me resuelve.
49. Pedí la devolución de $449 MXN hace una semana y todavía no me abonan nada. Nunca me avisaron de este cargo.
50. Pedí la cancelación de mi línea y aun así me cobraron $399.00 en diciembre. Quiero una explicación clara de este cargo. Saludos cordiales.

### Problemas de Facturación

1. Recibo mi factura 3 días tarde cada mes, ya casi vence cuando la tengo.
2. Mi plan Premium 999 cuesta 999 pesos, pero en mi fact el monto total es de $1,199.50. Ya me habia pasado antes con otro periodo. Quedo al pendiente.
3. Mi pago de $799.50 sí se hizo pero la factura todavía marca vencido.
4. La factura trae un saldo anterior que ya pagué, el pago de $129 no aparece reflejado. Necesito el documento correcto con urgencia.
5. Necesito el comprobante fiscal (CFDI) de marzo y el sistema no me lo genera (caso 668448). Nadie me da una explicación clara.
6. Buenas tardes. Mi factura de noviembre dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente. No me gusta quedarme con dudas de lo que pago. Saludos cordiales.
7. Quiero cambiar mi factura a entrega en papel y no encuentro la opción por ningún lado.
8. Recibo mi factura 5 días tarde cada mes, ya casi vence cuando la tengo.
9. Mi fact sale a nombre de otra persona y necesito que corrijan los datos.
10. Me emitieron dos facturas por el mismo periodo de mayo, una de ellas está duplicada!!!
11. Mi recibo viene con efcha límite de pago distinta a la q me indicaron por teléfono
12. Hola. No me llega mi factura por WhatsApp desde febrero.
13. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualice. No me gusta quedarme con dudas de lo que pago.
14. No puedo registrar mi correo para facturación, la página me regresa un error (TKT-594213).
15. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Necesito que lo corrijan cuanto antes.
16. Buenos días. Falta el folio fiscal en mi factura de marzo y no puedo deducirla.
17. Mi fact llega con los datos de mi domicilio anterior aunque ya lo actualice. Espero una respuesta pronto...
18. Cambié de plan a Nocturno 79 y la factura sigue mostrando la tarifa anterior. Ya me había pasado antes con otro periodo.
19. Les escribo por un tema de facturación. No me llega mi factura por correo electrónico desde junio. No me gusta quedarme con dudas de lo que pago.
20. La factura trae un saldo anterior que ya pagué, el pago de $499.00 no aparece reflejado. Quedo al pendiente.
21. Los consumos en mi factura de noviembre no corresponden con mi historial de llamadas.
22. Caray, mi estado de cuenta no coincide con lo q me dijeron en el área de cobranza. Me urge tenerlo en regla...
23. La fecha de vencimiento de mi recibo no es la q me dieron, por eso se me juntan los periodos. Muchas gracias x su atencion.
24. Mi recibo no muestra el descuento prometido en ventas, solo el total de 1,499 pesos.
25. Quiero cambiar mi fact a entrega x correo electronico y no encuentro la opcion x ningun lado. Saludos cordiakes.
26. Falta el folio fiscal en mi factura de ferbero y no puedo deducirla. No me gusta quedarme con dudas de lo q pago!!!
27. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos por correo. Es muy molesto andar aclarando lo mismo cada mes.
28. O sea, mi estado de cuenat no coincide con lo que me dijeron en el área de cobranza...
29. No aparece el pago que hice en la app en mi estado de cuenta, sigo con saldo pendiente.
30. El ciclo de facturación cambió y ahora tengo dos recibos casi juntos, uno de $499.
31. Cambie de plan a Max 699 y la factura sigue mostrando la tarifa anterior...
32. Hola. Mi factura llega con los datos de mi domicilio anterior ahnque ya lo actualice. Mucchas gracias x su atencion.
33. Soy cliente de TelArg desde hace años. No puedo registrar mi correo para facturación, la página me regresa un error (#74626).
34. Soy cliente de TelArg desde hace años. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos en la app. Ya llamé tres veces y nadie me resuelve. Espero su pronta respuesta.
35. Buenas tardes. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el otal no cuadra. Espero una respuesta pdonto.
36. Mi pago de $299 sí se hizo pero la factura todavía marca vencido. Quedo al pendiente.
37. Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé.
38. Solicité factura electrónica en la app y me siguen enviando el recibo impreso. No me gusta quedarme con dudas de lo que pago. Saludos cordiales.
39. mi fechha de corte cambio al dia 5 sin avisarme y ahora el estado de cuenta viene desfasado. saludos cordiales.
40. Buenos días. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. Espero una respuesta pronto.
41. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. Espero una respuesta pronto.
42. Saludos. Pedi la factura de mi equipo y me mandaron la de otro numero de linea. Gracias.
43. Mi fact sale a nombre de otra persona y necesito que corrijan los datos. Gracias.
44. Mi factura de julio dice que tengo adeudo pero el sistema de pagos dice que estoy al corriente.
45. La factura trae un saldo anterior que ya pagué, el pago de $999.50 no aparece reflejado. Ya lo reporté y sigue igual.
46. No me llega mi factura en la app desde enero.
47. Les escribo por un tema de facturación. Mi fecha de corte cambió al día 20 sin avisarme y ahora el estado de cuenta viene desfasado. Me urge tenerlo en regla. Quedo al pendiente.
48. Recibo mi fact 7 dias tarde cada mes, ya casi vence cuando la tengo.
49. Necesito el comprobante fiscal (CFDI) de agosto y el sistema no me lo genera (reporte 56184).
50. La fecha de vencimiento de mi recibo no es la que me dieron, x eso se me juntan los periodos. Ya llame tres veces y nadie me resuelve.

### Quejas

1. De planno el trato de Alejando en atencion a clientes fue grosero no pienso aguantar eso. q alguien se haga responsable por favor. Gracias.
2. Una vergüenza la atención: me dejaron esperando media hora y se me cortó la llamada. Pido que tomen en cuenta mi queja.
3. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio.
4. Me hicieron perrder toda la tarde en la sucursal y al final no resolvieron nada. Llevo años con ustedes y no merezco este trato.
5. Tengo una inconformidad que presentar. Mi queja: el call center cuelga cuando pido hablra con un supervisor. Gracias.
6. Subieron el precio de mi plan Nocturno 79 sin avisarme, esto me parece un abuso.
7. Saludos. La atención en la sucursal de Morelia es pésima, hacen fila más de una hora para nada.
8. A quien corresponda. Esperé casi dos horas en el area de cobranza y nadie me contestó. Quiero q alguien me explique x q
9. Me vendieron el plan Empresarial 899 con condiciones que no cumplen, publicidad engañosa. Es una falta de respeto con los clientes.
10. Nadie da seguimiento a mi reporte (TKT-62722), siemmpre me dicen que ya esta en proceso
11. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. Espero que alguien con autoridad me atienda.
12. Me llaman seguido en fines de semana para ofrecerme promociones aunque ya les pedí que no. Espero su pronta respuesta.
13. Me prometieron una llamada de vuelta y después de tres días sigo esperando (ticket 70145).
14. Estoy muy molesto con la forma en que Ricardo me atendió, sin ninguna empatía. Espero su pronta respuesta.
15. La neta, qué mala atención, desde febrero estoy tratando de aclarar un tema y nadie me resuelve.
16. Me vendieron el plan Empresarial 899 con condiciones que no cumplen, publicidad engañosa.
17. Cada vez que llamo me piden los mismos datos y no resuelven nada estoy harto del servicio de atencion. Quiero que alguien me explique por que.
18. El proceso para cambiar de plan a Empresarial 899 fue lentísimo y lleno de requisitos absurdos.
19. Me llaman seguido toda la tarde para ofrecerme promociones aunque ya les pedi que no. q alguien se haga responsable, x favor.
20. Me cambiaron las condiciiones de mi plan Libre 199 sin pedirme autorizacion, ya no es lo que contrate. Estoy pensando seriamente en cambiarme de compañía.
21. El chat de TelArg es inútil, solo responden con mensajes automáticos. Ya llamé tres veces y nadie me resuelve.
22. Subieron el precio de mi plan Hogar 449 sin avisarme, esto me parece un abuso.
23. Buenas tardes. Chale, no es posible que una empresa tan grande tenga una atención tan lenta. Quedo atento a su respuesta.
24. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio. Espero una respuesta pronto.
25. Me prometieron una llamada de vuelta y después de cinco días sigo esperando (#74248).
26. Subieron el precio de mi plan Hogar 449 sin avisarme, esto me parece un abuso. Espero que alguien por fin me responda.
27. El chat de TelArg es inutil, solo responden con mensajes automaticos. Que alguien se haga responsable, x favor.
28. Esperé media hora en ventas y nadie me contestó. No es justo que uno tenga que insistir tanto. Espero su pronta respuesta.
29. Me cambiaron las condiciones de mi plan Empresarial 899 sin pedirme autorización, ya no es lo que contraté. Estoy pensando seriamente en cambiarme de compañía. Espero su pronta respuesta.
30. Subieron el precio de mi plan Hogar 449 sin avisarme, esto me parece un abuso. Quiero q alguien me expilque por que!!!
31. Soy cliente de TelArg desde hace años. Subieron el precio de mi plan Total 399 sin avisarme, esto me parece un abuso.
32. Subieron el precio de mi plan Prepago Basico sin avisarme, esto me parece un abuso.
33. Me llaman seguido toda la tarde para ofrecerme promociones aunque ya les pedí que no. No es justo que uno tenga que insistir tanto.
34. La aplicacion se cierra sola y soporte no hace nada, queja formal con ticket 710568. Llevo años con ustedes y no merezco este trato!!!
35. Les escribo porque quiero presentar una queja. Nadie da seguimiento a mi reporte (TKT-3557024), siempre me dicen que ya esta en proceso...
36. Mi queja: el call center cuelga cuando pido hablar con un supervisor. Todo esto es muy frustrante.
37. Esperé una hora en soporte técnico y nadie me contestó. Estoy pensando seriamente en cambiarme de compañía.
38. Mi queja es por la falta de interés en el call center, nunca devuelven las llamadas.
39. Ya presenté mi inconformidad en la app y no tengo respuesta que falta de seriedad. No es justo que uno tenga q insistir tanto. Espero su pronta respuesta.
40. La atención en la sucursal de Puebla es pésima, hacen fila cuarenta minutos para nada. Quiero que alguien me explique por qué.
41. El chat de TelArg es inútil, solo responden con mensajes automáticos. Gracias.
42. Subieron el precio de mi plan Joven 149 sin avisarme, esto me parece un abuso.
43. Me prometieron una llamada de vuelta y después de cuatro días sigo esperando (folio 67020). Espero que alguien por fin me responda.
44. En la sucursal me dijeron que no era su responsabilidad y me mandaron a otra área.
45. Saludos. Una vergüenza la atención: me dejaron esperando media hora y se me cortó la llamada. Quiero que alguien me explique por qué.
46. Hola. En atención a clientes me dijeron que no era su responsabilidad y me mandaron a otra área. Que alguien se haga responsable, por favor.
47. Me transfirieron 3 veces entre áreas y nadie pudo resolver mi problema. Que alguien se haga responsable, por favor.
48. Quiero dejar constancia de mi inconformidad. En atención a clientes me dijeron que no era su responsabilidad y me mandaron a otra área. Quedo atento a su respuesta.
49. De plano, qué mala atención, desde diciembre estoy tratando de aclarar un tema y nadie me resuelve. Pido que tomen en cuenta mi queja. Espero que alguien por fin me responda.
50. Quiero dejar constancia de mi inconformidad. Quiero presentar mi queja porque Carlos me colgó la llamada en plena aclaración.

### Reconocimientos

1. QUIERO AGRADECER LA RAPIDEZ CON LA Q ATENDIERON MI REPORTE (REPORTE 9263662), TODO QUEDO RESUELTO. POR ESTO SIGO CON TELARG.
2. En la sucursal de Chihuahua me atendieron de maravilla, todo fue rápido y claro.
3. Quiero reconocer públicamente la labor de Héctor en soporte técnico, un trabajo verdaderamente profesional. Realmente superó mis expectativas. Sigan así.
4. MIRA NOMAS, QUÉ GUSTO Q CON SOPORTE TECNICO SE RESOLVIO MI PROBLEMA AL PRIMER INTENTO. SEGUIRÉ RECOMENDANDOLOS.
5. En la sucursal de Toluca me atendieron de maravilla, todo fue rápido y claro. Ojalá todas las empresas atendieran así.
6. Les escribo con mucho gusto. Me encantó la atención por teléfono, respondieron en minutos y con mucha amabilidad.
7. Excelente cobertura en Guadalajara, nunca he tenido problemas con mi plan Total 399. Gracias.
8. Gracias por restablecer mi servicio el mismo día, atención de primera.
9. Quiero compartir mi experiencia. Un diez para el area de cobranza: resolvieron mi duda en el portal web sin darme vueltas.
10. Les escribo con mucho gusto. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. Saludos cordiales.
11. Muy contento con el internet de mi plan Prepago Básico, la velocidad es increíble. Eso es lo que uno espera de una buena compañía.
12. Mi señal es excelente incluso en toda mi calle, felicidades a todo el equipo de redes. Todo el equipo merece un aplauso.
13. Tengo la linea a mi nombre. Agradezco que me hayan devuelto mi dinero tan rapido, excelente respuessta.
14. QUIERO COMPARTIR MI EXPERIENCIA. QUIERO FELICITAR A ALEJANDRO EN VENTAS, LA ATENCION FUE SOBRESALIENTE.
15. Les escribo desde San Luis Potosí. Me encanto la atencion en el portal web, respondieron en minutos y con mucha amabilidad. Con gusto volveria a acudir con usfedes.
16. La verdad, qué gusto que con la sucursal se resolvió mi problema al primer intento.
17. A quien corresponda. El cambio a plan Prepago Básico fue sencillo y sin trámites, gracias por facilitarlo. Un saludo a todo el equipo.
18. Les escribo desde Morelia. Reconoozco el esfuerzo q hacen en ventas soempre estan al pendiente de los clientes.
19. Soy cliente de TelArg desde hace años. Agradezco que me hayan devuelto mi dinero tan rápido, excelente respuesta. Realmente superó mis expectativas.
20. Les escribo con una buena noticia. Mi experiencia en el portal web fue de lo mejor Miriam se toko el tiempo de explicarme todo.
21. Buenas noches. Felicidades por la cobertura en Aguascalientes, hasta en carretera tengo señal.
22. Buenas tardes. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo. Eso es lo que uno espera de una buena compañía.
23. Verónica me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Sigan así, de verdad se agradece.
24. ME ENCANTO LA ATENCIÓN EN LA SUCURSAL, RESPONDIERON EN MINUTOS Y CON MUCHA AMABILIDAD. SEGUIRÉ RECOMENDANDOLOS.
25. Mi reconocimiento al tecnico que llego puntual y dejo mi internet funcionando perfecto.
26. Quiero agradecer la rapidez con la que atendieron mi reporte (reporte 290575), todo quedó resuelto. Con gusto volvería a acudir con ustedes.
27. Las mejores promociones del mercado, sigo con TelArg por el buen servicio. Un saludo a todo el equipo.
28. El servicio de TelArg ha mejorado muchisimo se nota el compromio con los clientes.
29. Gracias por La Llamada de Seguimiento, nunca me Habia pasado q una Empresa se Preocupara Asi. es Un gusto q existan Personas asi.
30. Un diez para el área de cobranza: resolvieron mi duda por chat sin darme vueltas. De verdad quedé muy satisfecho. Muchas gracias por su atención.
31. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. De verdad quedé muy satisfecho.
32. Muchas gracias por resolver mi problema en el portal web, el servicio fue de diez. De verdad quedé muy satisfecho. Seguiré recomendándolos.
33. Mira Nomas, q Buen Servicio, En menos De media hora Me Solucionaron todo. me dio mucho Gusto ser su Cliente. Seguire Recomendandolos.
34. El equipo que me entregaron llego antes de lo prometido y en perfecto estado, muchas gracias.
35. Mi experiencia en la sucursal fue de lo mejor, Miriam se tomó el tiempo de explicarme todo. Se nota que se preocupan por sus clientes. Seguiré recomendándolos.
36. A quien corresponda. Excelente cobertura en Guadalajara, nunca he tenido problemas con mi plan Hogar 449. Que tengan un excelente día.
37. Mi señal es excelente incluso en el centro, felicidades a todo el equipo de redes. Que tengan un excelente día.
38. Quiero dejar constancia de lo bien que me atendio Fernanda, un trato impecable. q tengan un excelente día.
39. Un diez para la sucursal: resolvieron mi duda en la sucursal sin darme vueltas. Saludos cordiales.
40. El Cambio A plan nocturno 79 Fue sencillo y sin Tramites, gracias X Facilitarlo. Defihitivamente Los Voy a Recomendar con mi Familia...
41. Carlos me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Realmente superó mis expectativas.
42. Buenos dias. Quiero felicitar a Jorge en soporte tecnico, la atencion fue de diez. Saludos cordiales.
43. Carlos me atendio con mucha atencion y profesionalismo, merece un reconocimiento por su trabajo. Que tengan un excelente dia.
44. La promoción del plan Nocturno 79 es muy buena, gracias por pensar en los clientes. Realmente superó mis expectativas.
45. Les escribo con mucho gusto. Muy contento con el internet de mi plan Social 299, la velocidad es increíble.
46. A quien corresponda. La promoción del plan Libre 199 es muy buena, gracias por pensar en los clientes. De verdad quedé muy satisfecho. Sigan así.
47. Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante. Sigan así, de verdad se agradece. Un saludo a todo el equipo.
48. Agradezco que me hayan devuelto mi dinero tan rápido, excelente respuesta. Todo el equipo merece un aplauso.
49. Quiero agradecer la rapidez con la que atendieron mi reporte (caso 6665465), todo quedó resuelto.
50. Les escribo con una buena noticia. Gracias grracias y gracias x fin alguien resolvio mi problema de manera agil.
