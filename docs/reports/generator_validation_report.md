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
| noise rate | PASS | realized 0.3025, configured 0.30, tolerance ±0.015 |
| lexical overlap | PASS | pairwise Jaccard of top-100 tokens in [0.198, 0.361], required within [0.15, 0.80] |
| sample review | PASS | 0 of 250 sampled comments rejected (empty, unfilled slot, <10 or >800 runes, <2 words) |

Corpus statistics: 20.20% duplicate texts; text length 41 / 111.1 / 237 runes (min / mean / max).

## Lexical overlap

Jaccard similarity of the 100 most frequent tokens per category (accepted range for distinct pairs: [0.15, 0.80]).

| | Falta de Servicio | Cobros Indebidos | Problemas de Facturación | Quejas | Reconocimientos |
|---|---:|---:|---:|---:|---:|
| Falta de Servicio | 1.000 | 0.325 | 0.282 | 0.316 | 0.250 |
| Cobros Indebidos | 0.325 | 1.000 | 0.361 | 0.307 | 0.220 |
| Problemas de Facturación | 0.282 | 0.361 | 1.000 | 0.266 | 0.198 |
| Quejas | 0.316 | 0.307 | 0.266 | 1.000 | 0.316 |
| Reconocimientos | 0.250 | 0.220 | 0.198 | 0.316 | 1.000 |

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
| **Realized selection rate** | **0.3025** (tolerance ±0.015) |

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

Test comments with an identical text in train: 3690 of 14600 (25.27%). A temporal split does not prevent exact repeats of a template draw.

corpus.db SHA-256: `86eb43035be94edb7a6e5c653187696aa4fd7eddc7837cff1adc05ad201f1876`

## Sample review

50 comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.

### Falta de Servicio

1. No tengo datos móviles aunque mi saldo es suficiente (ref. folio 718269). Saludos cordiales.
2. Saludos. El internet de mi plan Empresarial 899 no funciona desde diciembre y no hay quien lo repare.
3. Mi chip nuevo no activa, llevo diez días esperando que me den servicio de llamadas. Espero que lo atiendan pronto.
4. Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace cinco días. Ya insistí tres veces y nadie me resuelve.
5. saludos. en ttoda la zona no hay cobertura toda la tarde, parece que la antena está caída. x favor, avísenme cualquier novedad!!!
6. A quien corresponda. No tengo cobertura en mi oficina en fines de semana, ya pedi q lo revisen (ref. TKT-97787). Saludos cordiales.
7. Desde que cambiaron la antena en mi barrio perdí servicio de mensajes y desde entonces sigo igual.
8. A quien corresponda. El módem se queda con la luz roja y no hay internet, llevo varias horas igual. Ya me está afectando en el trabajo.
9. Caray, llevo cinco días sin cobertura y sigo pagando mi plan Empresarial 899. Sigo pagando completo y no obtengo lo que corresponde.
10. Buenos días. El internet de mi plan Amigo Sin Límite no funciona desde febrero y no hay quien lo repare. Quedo al pendiente.
11. hola. hay perdida de servicio de mensajes en las noches en toda mi calle, y ya son varios dias con el mismo problema. espero su pronta respuesta.
12. Mi celular dice sin servicio desde hace cuatro días, ya reinicié todo y nada. Quiero que alguien se haga cargo.
13. La portabilidad de mi número no se completó y estoy sin línea desde hace varias horas. Quiero que alguien me explique qué está pasando. Espero su pronta respuesta.
14. El tecnico nunca llego a devisar la falla de internet, ya llevo varrias horas de espera. Ya me habia pasado antes este mes.
15. Me quede sin servicio de llamadas justo cuando más lo necesitaba, ¿hasta cuáhdo van a arreglarlo? Ya me esta afectando en el trabajo.
16. Necesito internet para trabajar y llevo dos días sin nada, esto ya es urgente. Sigo pagando completo y no obtengo lo que corresponde. Muchas gracias por su atención.
17. Mi línea no registra en la red desde hace dos días, llevo todo ese tiempo sin poder comunicarme.
18. mi linea sigue sin servicio tras la reconexión q pague (ref. caso 38483). esto afecta a toda mi familia. x favor, avisenme cualquier novedad.
19. Hice mi recarga de MXN 399 hace una semana y sigo sin cobertura. Muchas gracias por su atención.
20. No tengo datos moviles aunque mi saldo es suficiente (ref. reporte 85937)!!!
21. les escribo porque necesito ayuda. No puedo mandar ni recibir mensajes de texto en Veracruz desde hace cinco días. Quiero que alguien se haga cargo.
22. Tengo un problema con mi servicio. Pagué mi plan Amigo Sin Límite puntual y aun así me suspendieron el servicio. Necesito que me den una solución definitiva. Espero su pronta respuesta.
23. Cambié de SIM y quedé incomunicado, sin servicio de llamadas desde hace diez días. Ya perdí la cuenta de las veces que he insistido. Quedo atento a su respuesta.
24. Buenas noches. No tengo internet en mi trabajo en las noches, ya pedí que lo revisen (ref. TKT-73936).
25. Hola. Cuando intento llamar solo escucho que el número no está disponible, ya llevo cuatro días así.
26. Hice mi recarga de $49 MXN hace cinco días y sigo sin internet. Es muy frustrante estar así. Quedo atento a su respuesta.
27. Buenos días. Activé el roaming para viajar y al llegar a Aguascalientes no tuve señal.
28. Avisé de la falla por teléfono con ticket 308956 y todavía no me restablecen el servicio.
29. Mi celular dice sin servicio desde hace cuatro dias, ya reinicie todo y nada. Ya me esta afectando en el trabajo.
30. Desde hace una semana estoy sin señal en mi trabajo y nadie me da una solución. Ya me había pasado antes este mes. Gracias de antemano.
31. Les escribo por una falla en mi línea. Mi chip nuevo no activa, llevo más de 48 horas esperando que me den datos móviles. Ya probé todo lo que me dijeron y sigue igual.
32. Hola. Avisé de la falla por chat con caso 700567 y todavía no me restablecen el servicio. Quedo atento a su respuesta.
33. No tengo datos móviles aunque mi saldo es suficiente (ref. reporte 4502542). Espero su pronta respuesta.
34. La portabilidad de mi número no se completó y estoy sin línea desde hace diez días. Esto afecta a toda mi familia.
35. Hola. Cuando intento llamar solo escucho que el número no está disponible, ya llevo tres días así.
36. El técnico nunca llegó a revisar la falla de datos móviles, ya llevo varias horas de espera. Espero su pronta respuesta.
37. Hice mi recarga de 99 pesos hace varias horas y sigo sin datos móviles. Nadie le da seguimiento a esto. Quedo al pendiente.
38. Buenas tardes. No tengo servicio de llamadas en mi casa de madrugada, ya pedi que lo revisen (ref. #547424).
39. Se me fue por completo la señal en León y no he podido hacer ni recibir llamadas.
40. Les escribo por una falla en mi línea. Avisé de la falla en la app con folio 87438 y todavía no me restablecen el servicio. Es muy frustrante estar así.
41. La portabilidad de mi número no se completó y estoy sin línea desde hace más de 48 horas. Sigo pagando completo y no obtengo lo que corresponde. Espero una respuesta pronto.
42. El internet de mi plan Premium 999 no funciona desde marzo y no hay quien lo repare.
43. Cuando intento llamar solo escucho que el número no está disponible, ya llevo diez días así.
44. Buenas tardes. No puedo mandar ni recibir mensajes de texto en Guadalajara desde hace una semana.
45. Desde que cambiaron la antena en el centro perdí cobertura y desde entonces sigo igual. Gracias de antemano.
46. En mi trabajo nunca hay señal, tengo que salir a la calle para poder llamar. Es muy frustrante estar así. Quedo al pendiente.
47. Soy cliente de TelArg desde hace años. Mi chip nuevo no activa, llevo varias horas esperando que me den servicio de llamadas. Ya insisti tres veces y nadie me resuelve. Gracias de antemano.
48. Active el roaming para viajar y al llegar a Guadalajara no tuve servicio de llamadas. Muchas gracias x su atención!!!
49. MI LINEA ESTA SUSPENDIDA SIN MOTIVO Y NO PUEDO USAR NADA. MUCHAS GRACAIS POR SU ATENCIÓN.
50. Buenas tardes. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir.

### Cobros Indebidos

1. Pagué mi mensualidad completa por teléfono y me volvieron a cobrar $399 MXN. Nunca me avisaron de este cargo.
2. Me quitaron el saldo de mi recarga de MXN 149 sin explicación alguna. Me urge que lo corrijan cuanto antes. Muchas gracias por su atención.
3. Desconozco el cobro de 49 pesos del día 15 de julio, quiero que lo reviertan.
4. Cargos q no reconozco en agosto: $229 MXN x servicio premium. Exijo que los eliminen...
5. Mi banco me reporta un cargo no reconocido de TelArg por $449.00, quiero aclararlo. Esto ya me pasó antes y nunca lo arreglaron.
6. Cancelé paquete de datos extra hace meses y me siguen descontando $79 MXN cada mes. Espero una respuesta pronto.
7. Buenos días. Cada mes se me va dinero por contenido de terceros y no sé cómo quitarlo.
8. Me descontaron 149 pesos de mi saldo sin haber usado nada (folio 941897). Quiero que me devuelvan mi dinero.
9. Cancelé mensajes premium hace meses y me siguen descontando $79 cada mes.
10. A quien corresponda. Mi banco me reporta un cargo no reconocido de TelArg por $699, quiero aclararlo. Muchas gracias por su atención.
11. Me siguen cobrando la línea que di de baja, ya llevo 3 meses pagando de más. Me urge que lo corrijan cuanto antes. Gracias.
12. Híjole, me robaron $99 MXN de saldo en un solo día sin hacer llamadas ni usar datos. Gracias.
13. Soy cliente de TelArg desde hace años. Ya reclamé el cobro de reconexión en el portal web y nadie me devuelve mi dinero (caso 4592462).
14. Buenso dias. Me aplicaron un cargo x reconexion de $229.50 cuando mi linea nunca estuvo suspendida. Ya pedi una aclaracion y sigo sin respuesta
15. Hola. Se me cobró una penalización de $1,199 que no corresponde porque ya cumplí mi contrato.
16. Me hicieron un cobto doble de $999 en octubre, exijo la devolucion.
17. Me cobraron roaming internacional de $229 sin salir de Culiacán. Espero su pronta respuesta.
18. Quiero que me reembolsen los MXN 499 que me cobraron de más (caso 365044). Espero una respuesta pronto.
19. me cobraron roaming internacional de $549 MXN sin salir de Oaxaca.
20. Desde que cambié de plan a Total 399 me cobran $49.50 de más, es un cobro indebido.
21. Cada mes se me va dinero por paquete de datos extra y no sé cómo quitarlo. Gracias.
22. Me quitaron $249 x mensajes premium que jamas pedí. x favor, avisenme cualquier novedad.
23. Les escribo porque necesito ayuda. Me aplicaron un cargo por reconexión de MXN 79 cuando mi línea nunca estuvo suspendida.
24. Cancelé suscripción de música hace meses y me siguen descontando 299 pesos cada mes.
25. ACTIVARON RECARGA AUTOMATICA SIN PEDÍRMELO Y ME COBRARON 129 PESOS EN EL ESTADO DE MI TARJETA. GRAACIAS...
26. Sin mi consentimiento me dieron de alta en mensajes premium y ya van 4 cobros.
27. Me siguen cobrando la línea que di de baja, ya llevo 2 meses pagando de más.
28. Buenas tardes. En mi tarjeta aparece un cargo de $249 MXN de TelArg que desconozco. No pienso pagar un cargo que no corresponde.
29. Me quitaron el saldo de mi recarga de $199 sin explicación alguna.
30. me Quitaron mXN 299 Por mensajes premium Que jamás pedi!!!
31. Se me cobró una penalización de MXN 899 que no corresponde porque ya cumplí mi contrato.
32. Me quitaron el saldo de mi recarga de MXN 449 sin explicación alguna. Quiero una explicación clara de este cargo.
33. Ya reclamé el cobro de plan adicional por correo y nadie me devuelve mi dinero (#7340995).
34. Cargos que no reconozco en junio: $129.00 por renta de equipo. Exijo que los eliminen. Quiero una explicación clara de este cargo. Quedo al pendiente.
35. Tengo un problema con un cargo en mi cuenta. Me descontaron $249 de mi saldo sin haber usado nada (TKT-520611).
36. Después de cancelar mi plan Max 699 me cargaron $699.00 a mi tarjeta. Espero su pronta respuesta.
37. Terminé de pagar el equipo hace meses y aún me cargan MXN 49 mensuales.
38. Sin mi consentimiento me dieron de alta en suscripción de música y ya van 8 cobros.
39. Quiero que me reembolsen los 799 pesos que me cobraron de más (ticket 3152382). Ya pedí una aclaración y sigo sin respuesta. Espero su pronta respuesta.
40. Activaron membresía de juegos sin pedírmelo y me cobraron $79.00 en el estado de mi tarjeta. Gracias.
41. Desde que cambié de plan a Ahorro 99 me cobran $79 de más, es un cobro indebido. Quiero que me devuelvan mi dinero. Espero su pronta respuesta.
42. Activaron renta de equipo sin pedírmelo y me cobraron MXN 229 en el estado de mi tarjeta. Espero una respuesta pronto.
43. A quien corresponda. Cobro indebido: me descontaron $249 dos veces el mismo día por lo mismo. Quiero una explicación clara de este cargo.
44. Se me cobró una penalización de 1,499 pesos que no corresponde porque ya cumplí mi contrato. Ya insistí tres veces y nadie me resuelve. Quedo al pendiente.
45. Después de cancelar mi plan Max 699 me cargaron MXN 699 a mi tarjeta. Nunca me avisaron de este cargo.
46. Me siguen cobrando la línea que di de baja, ya llevo 2 meses pagando de más. Esto ya me pasó antes y nunca lo arreglaron.
47. Me descontaron $1,999 de mi saldo sin haber usado nada (reporte 602750). Nadie me da una respuesta concreta.
48. Pedí la devolución de $229 MXN hace diez días y todavía no me abonan nada. Exijo que lo revisen cuanto antes.
49. activaron servicio premium sin pedirmelo y me cobraron $229 mxn en el estado de mi tarjeta!!!
50. Me cobraron $1,499.50 por una promoción que nunca acepté (TKT-33636).

### Problemas de Facturación

1. El portal muestra un monto de factura distinto al del recibo que me llegó, $999.00 contra $699.00. Espero su pronta respuesta.
2. Buenos dias. Pedí la fact de mi equipo y me mandaron la de otro número de línea. No me gusta quedarme con dudas de lo q pago!!!
3. Cambié de plan a Joven 149 y la factura sigue mostrando la tarifa anterior. Ya insistí tres veces y nadie me resuelve. Espero su pronta respuesta.
4. Solicité factura electrónica por correo electrónico y me siguen enviando el recibo impreso. Quiero que me confirmen cuándo estará corregido.
5. Mi fecha de corte cambió al día 10 sin avisarme y ahora el estado de cuenta viene desfasado. Ya intenté arreglarlo por mi cuenta y no pude.
6. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos por correo. No me gusta quedarme con dudas de lo que pago.
7. Qué barbaridad, mi estado de cuenta no coincide con lo que me dijeron en la sucursal. Ya insistí tres veces y nadie me resuelve. Saludos cordiales.
8. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Espero su pronta respuesta.
9. Mi fecha de corte cambió al día 20 sin avisarme y ahora el estado de cuenta viene desfasado.
10. Mi dirección de facturacion es incorrecta y por eso no llegan mis recibos. Gracias de antemano.
11. Saludos. No me llega mi factura por WhatsApp desde marzo. Necesito que esto quede aclarado con urgencia. Quedo al pendiente.
12. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a $999.50. Quiero que me confirmen cuándo estará corregido.
13. Quiero cambiar mi factura a entrega en papel y no encuentro la opción por ningún lado.
14. Buenas noches. Me emitieron dos facturas por el mismo periodo de marzo, una de ellas está duplicada. Gracias de antemano.
15. Los consumos en mi factura de abril no corresponden con mi historial de llamadas. Saludos cordiales.
16. Me emitieron dos facturas por el mismo periodo de agosto, una de ellas está duplicada. Muchas gracias por su atención.
17. la factura de diciembre viene con el plan equivocado aparece amigo sin límite y yo contrate otro.
18. Me llegó la factura en blanco, sin importes, desde que activé la entrega en papel. Gracias de antemano.
19. mi factura de febrero dice que tengo adeudo pero el sistema de pagos dice q estoy al corriente. gracias de antemano...
20. El ciclo de facturación cambió y ahora tengo dos recibos casi juntos de MXN 199 cada uno. Nadie me da una explicación clara.
21. Mi fact llega con los datos de mi domicilio anterior aunque ya lo actualicé.
22. Mi dirección de facturación es incorrecta y por eso no llegan mis recibos.
23. Hola. La factura de septiembre viene con el plan equivocado, aparece Amigo Sin Límite y yo contraté otro.
24. Los consumos en mi factura de diciembre no corresponden con mi historial de llamadas.
25. Quiero cambiar mi factura a entrega por correo electrónico y no encuentro la opción por ningún lado. Ya me había pasado antes y nadie lo arregló de fondo.
26. la Fecha de Vencimiento de mi recibo No Es La Que Me dieron, x eso se Me juntan los Periodos...
27. Mi factura sale a nombre de otra persona y necesito que corrijan los datos.
28. me emitieron dos facturas por el mismo periodo de marzo, una de ellas está duplicada. Nadie me da una explicación clara.
29. Buenas tardes. Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. Ya intenté arreglarlo por mi cuenta y no pude.
30. les escribo por un tema de facturacion. o sea, mi estado de cuenta no coincide con lo q me dijeron en el area de cobranza. ya lo reporte y sigue igual.
31. No aparece el pago que huce en la app en mi estado de cuenta, sigo con saldo pendiente. Saludos cordiales.
32. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos en el portal web. Ya intenté arreglarlo por mi cuenta y no pude.
33. El esstado de cuenta de enero no detalla los conceptos, solo trae un total de $1,999.50.
34. mi factura sale a nombre de otra persona y necesito que corrijan los datos.
35. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a 549 pesos. Es muy molesto andar aclarando esto.
36. Tengo un problema con mi factura. Cambié de plan a Prepago Básico y la factura sigue mostrando la tarifa anterior. Nadie me da una explicación clara. Saludos cordiales.
37. Tengo la línea a mi nombre. Me emitieron dos facturas x el mismo periodo de abril, una de ellas esta duplicada
38. la fact trae un saldo anterior que ya pague el pago de $229.50 no aparece reflejado. necesito q esto quee aclarado con urgencia.
39. Les escribo porque necesito ayuda. Mi factura sale a nombre de otra persona y necesito que corrijan los datos.
40. Saludos. Mi recibo viene con fecha límite de pago distinta a la que me indicaron por chat. Quedo al pendiente.
41. Buenos dias. La fecha de vencimiento de mi recibo no es la que me dieron, x eso se me juntan los periodos...
42. NNecesito aclarar mi estado de cuenta, hay conceptos q no entiendo y el total no cuadra. Por favor, avisenme cualquier novedad.
43. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra.
44. El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a $129.50. Quedo al pendiente.
45. Mi recibo viene con fecha limite de pago distinta a la que me indicaron por tel. Ya lo reporte y sigue igual. Espero su pronta respuesta.
46. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. Quedo atento a su respuesta.
47. El portal muestra un monto de factura distinto al del recibo que me llego, $149 MXN contra 499 pesos. Espero su pronta respuesta!!!
48. no aparece el pago que hice por teléfono en mi estado de cuenta, sigo con saldo pendiente. Es muy molesto andar aclarando esto.
49. No puedo descargar mi factura en el portal web, marca error cada vez. Por favor, avísenme cualquier novedad.
50. Buenos días. Recibo mi factura 5 días tarde cada mes, ya casi vence cuando la tengo.

### Quejas

1. Mi queja: el call center cuelga cuando pido hablar con un supervisor.
2. Ay no, el trato de Jorge en la sucursal fue grosero, no pienso aguantar eso...
3. Mi queja es por la falta de interés en el área de cobranza, nunca devuelven las llamadas.
4. La neta, no es posible que una empresa tan grande tenga una atención tan lenta. Ya insistí tres veces y nadie me resuelve.
5. Les escribo desde Hermosillo. No hay quien resuelva dudas en el call center, todo son excusas y tiempos de espera. Espero su pronta respuesta.
6. Hola. Para cancelar mi linea me ponen mil trabas, es un proceso interminable.
7. Para cancelar mi línea me ponen mil trabas, es un proceso interminable.
8. La sucursal cierra antes del horario que anuncian y ya van 5 veces que voy en balde. Ya insistí tres veces y nadie me resuelve.
9. Los asesores no saben nada, me dieron 4 respuestas distintas en la misma semana. Me parece inaceptable q nadie de la carra.
10. me transfirieron 3 veces entre áreas y nadie pudo resolver mi problema.
11. Quiero presentar mi queja porque Fernanda me colgó la llamada en plena aclaración. Me parece inaceptable que nadie dé la cara.
12. Buenos dias. El chat de TelArg es inutil, solo responden con mensajes automaticos
13. Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Me parece inaceptable que nadie dé la cara.
14. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas. Quiero que alguien me explique por qué.
15. Quiero dejar constancia de mi inconformidad. Caray, el trato de Itzel en ventas fue grosero, no pienso aguantar eso. Gracias.
16. Me cambiaron las condiciones de mi plan Empresarial 899 sin pedirme autorizacion, ya no es lo q contrate.
17. Los asesores no saben nada, me dieron 3 respuestas distintas en la misma semana. Quedo atento a su respuesta.
18. Mi queja: el call center cuelga cuando pido hablar con un supervisor. Espero que alguien por fin me responda.
19. El chat de TelArg es inútil, solo responden con mensajes automáticos.
20. En el área de cobranza me dijeron que no era su responsabilidad y me mandaron a otra área. Ya perdí demasiado tiempo en esto. Espero que alguien por fin me responda.
21. Una verguenza la atencion: me dejaron esperando veinte minutos y se me corto la llamada. Quiero q alguidn me explique por que. Saludos cordiales.
22. El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio.
23. Tengo una inconformidad que presentar. Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Ya perdí demasiado tiempo en esto.
24. estoy muy molesto Con La Forma en Q Daniela me atendio, Sin Ninguna empatia. todo esto Es muy Frustrante.
25. esperé Acsi dos horas En El Area de Cobranza Y Nadie Me Contesto. Saludos cordiales.
26. Hola. El proceso para cambiar de plan a Max 699 fue lentísimo y lleno de requisitos absurdos. Es una falta de respeto con los clientes.
27. tengo la línea a mi nombre. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención.
28. Buenas tardes. El chat de TelArg es inútil, solo responden con mensajes automáticos. Ya perdí demasiado tiempo en esto. Espero que alguien por fin me responda.
29. UNA VERGUENZA LA ATENICON: ME DEJARON ESPERANDO CUARENTA MINUTOS Y SE ME CORTÓ LA LLAMADA.
30. Buenas noches. Mala experiencia en la sucursal: tuve que repetir mi problema a cuatro personas distintas. Es una falta de respeto con los clientes. Muchas gracias por su atención.
31. La atención en la sucursal de León es pésima, hacen fila veinte minutos para nada.
32. Nadie da seguimiento A lo Que pedi (csao 120038), siempre Me dicen q ya Esta En Proceso. estoy pensando seriamente En Ccambiarme De compañia
33. Los asesores no saben nada, me dieron 2 respuestas distintas en la misma semana. Ya perdi demasiado tiempo en esto!!!
34. Buenas noches. La aplicación se cierra sola y soporte no hace nada, queja formal con caso 5353945. Muchas gracias por su atención.
35. Quiero presentar mi queja porque Miriam me colgó la llamada en plena aclaración. Saludos cordiales.
36. Tengo una inconformidad que presentar. La atención en la sucursal de Villahermosa es pésima, hacen fila cuarenta minutos para nada.
37. Mi queja es x la falta de interes en la sucursal, nunca devuelven las llamadas.
38. Ya presenté mi inconformidad en la sucursal y no tengo respuesta, qué falta de seriedad.
39. Me llaman seguido toda la tarde para ofrecerme promociones aunque ya les pedí que no. Espero que alguien con autoridad me atienda.
40. A quien corresponda. Exijo que me atiendan como cliente, no es la primera vez que me ignoran en el portal web. Ya insistí tres veces y nadie me resuelve. Gracias.
41. Quiero dejar constancia de mi inconformidad. Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas.
42. Me transfirieron 4 veces entre áreas y nadie pudo resolver mi problema.
43. Mi queja es por la falta de interés en ventas, nunca devuelven las llamadas.
44. Les escribo porque quiero presentar una queja. Exijo que me atiendan como cliente, no es la primera vez que me ignoran en el portal web.
45. Mala experiencia por correo: tuve que repetir mi problema a cuatro personas distintas. Espero su pronta respuesta.
46. Qué onda, qué mala atención, desde abril estoy tratando de aclarar un tema y nadie me resuelve.
47. Buenas noches. Me prometieron una llamada de vuelta y después de varias horas sigo esperando (caso 70386). Quiero que alguien me explique por qué.
48. Subieron el precio de mi plan Total 399 sin avisarme esto me parece un abuso. Ya peri demasiado tiempo en wsto.
49. LAS CONDICIONES DEL CONTRATO CAMBIAN CADA MES Y NADIE ME LAS EXPLICA CON CLARIDAD. PIDO Q TOMEN EN CUENTA MI QUEJA. SALUDOS CORDIALES.
50. Esperé cuarenta minutos en el área de cobranza y nadie me contestó.

### Reconocimientos

1. Muchas gracias por resolver mi problema por correo, el servicio fue de primera.
2. Quiero agradecer la rapidez con la que atendieron mi solicitud (#732275), todo quedó resuelto.
3. El cambio a plan Total 399 fue sencillo y sin trámites, gracias por facilitarlo. Se los quería hacer saber. Seguiré recomendándolos.
4. Les escribo con mucho gusto. Quiero agradecer la rapidez con la q atendieron mi solicitud (caso 444134), todo quedo resuelto. Todo el equipo merece un aplauso. Gracias.
5. La promoción del plan Premium 999 es muy buena, gracias por pensar en los clientes.
6. Gracias, gracias y gracias, x fin alguien resolvio mi roblema de manera agil. Se nota q se preocupan x sus clientes.
7. Excelente cobertura en Guadalajara, nunca he tenido problemas con mi plan Total 399. Gracias.
8. Gracias, gracias y gracias, por fin alguien resolvio mi problema de manera agil.
9. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Se nota que se preocupan por sus clientes.
10. Mi reconocimiento al técnico q llego puntual y dejo mi internet funcionando perfecto. Saludos cordiales.
11. Jorge me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Todo el equipo merece un aplauso.
12. Quiero dejar constancia de lo bien que me atendió Alejandro, un trato impecable. Ojalá todas las empresas fueran así.
13. Mi reconocimiento al tecnico que llego puntual y dejo mi internet funcionando perfecto. Voy a seguir con ustedes muchos años mas. Que tengan un excelente dia.
14. Mi experiencia en el portal web fue de lo mejor, Itzel se tomo el tiempo de explicarme todo.
15. El cambio a plan Joven 149 fue sencillo y sin tramites, gracias x facilitarlo...
16. Les escribo con mucho gusto. Las mejores promociones del mercado, sigo con TelArg por el buen servicio.
17. Quiero agradecer la rapidez con la q atendieron mi solicitud (ticket 34626), todo quedo resuelto. Se nota q se preocupan por sus clientes. Seguire recomendandolos...
18. Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo.
19. Les escribo con una buena noticia. Quiero agradecer la rapidez con la q atendieron mi solicitud (folio 35755), todo quedo resuelto. Sigan así.
20. El cambio a plan Premium 999 fue sencillo y sin trámites, gracias por facilitarlo. Todo el equipo merece un aplauso. Saludos cordiales.
21. Quiero felicitar a Ricardo en atención a clientes, el trato fue de diez.
22. En la sucursal de Oaxaca me atendieron de maravilla, todo fue rápido y claro.
23. Verónica me atendió con mucha atención y profesionalismo, merece un reconocimiento por su trabajo. Sigan así, de verdad se agradece.
24. Saludos. Muy contento con el internet de mi plan Libre 199, la velocidad es increíble. Se los quería hacer saber. Saludos cordiales.
25. Muchas gracias por resolver mi problema por correo, el servicio fue de lo mejor. Me dio mucho gusto ser su cliente.
26. Buenos dias. En la sucursal de Toluca me atendieron de maravilla, todo fue rapido y claro.
27. quiero reconocer publicamente la labor de laura en soporte tecnico, un trabajo verdaderamente profesional. seguiré recomendandolos
28. Les escribo con mucho gusto. Mi reconocimiento al tecnico que llego puntual y dejo mi internet funcionando perfecto!!!
29. Mi señal es excelente incluso en toda mi calle, felicidades a todo el equipo de redes. Gracias!!!
30. Les escribo con una buena noticia. Me encantó la atención en el portal web, respondieron en minutos y con mucha amabilidad. Un saludo a todo el equipo.
31. Quiero dejar constancia de lo bien q me atendio Paola, un trato impecable!!!
32. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes.
33. Gracias, gracias y gracias, por fin alguien resolvió mi problrma de manera agil. Me dio mucho gusto ser su cliente. Saludos cordiales.
34. Mi experiencia en la sucursal fue de lo mejor, Ricardo se tomó el tiempo de explicarme todo. Sigan así.
35. Quiero dejar constancia de lo bien q me atendio Daniela, un trato impecable.
36. Qué padre, qué gusto que con la sucursal se resolvió mi problema al primer intento.
37. Muchas gracias por resolver mi problema por correo, el servicio fue de primera.
38. El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Todo el equipo merece un aplauso.
39. Mi experiencia en la app fue de lo mejor, Alejandro se tomó el tiempo de explicarme todo.
40. Les escribo con una buena noticia. El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias. De verdad quedé muy satisfecho.
41. De plano, qué buen servicio, en menos de diez minutos me solucionaron todo. Voy a seguir con ustedes muchos años más.
42. Saludos. Muy contento con el internet de mi plan Amigo Sin Límite, la velocidad es increíble. Ojalá todas las empresas fueran así.
43. Hector me atendio con mucha atención y profesionalismo, merece un reconocimiento x su trabajo. Me dio mucho gusto ser su clienet. Sigan así.
44. Buenos días. Las mejores promociones del mercado, sigo con TelArg por el buen servicio.
45. Que padre que buen servicio en menos de media hora me solucionaron todo.
46. Buenos días. Les agradezco la atencion durante abril, me ayudaron mucho con mi cambio de linea. Se nita que se preocupan x sus clientes.
47. Mi señal es excelente incluso en toda mi calle, felicidades a todo el equipo de redes.
48. Mi señal es excelente incluso en el centro, felicidades a todo el equipo de redes. Voy a seguir con ustedes muchos años más.
49. Les escribo con mucho gusto. Felicidades por la cobertura en Mérida, hasta en carretera tengo señal.
50. Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante.
