# Muestra de 250 comentarios: borrador de revisión (TE §6.1.4)

Borrador generado a partir de la sección *Sample review* de `generator_validation_report.md` (corpus `37f285ff…`; 50 comentarios por categoría, flujo aleatorio propio con semilla+1). El usuario revisa cada comentario y deja constancia abajo; no se marca nada por él. La revisión anterior (250/250 sin observaciones) correspondía al corpus `b810bad2…` y no se traslada.

## Constancia de revisión

| Campo | Valor |
|---|---|
| Fecha de revisión | |
| Revisor | |
| Comentarios revisados | /250 |
| Observaciones | |

Criterio de la revisión: cada texto es plausible como comentario de un cliente de TelArg y pertenece a la categoría bajo la que aparece. Marca `x` en la columna *OK* si cumple; si no, anota el motivo en *Observación*. Para una plantilla **compartida** basta con que la categoría sea una lectura razonable del texto: el mismo texto existe a propósito en otra categoría (ambigüedad controlada, K2, `PLAN.md` §3); su lectura alternativa está en `plantillas_nuevas.md`.

La columna *Auto* la llena un script con verificaciones automáticas, **no sustituye tu juicio**: coherencia de montos (cobro mayor que el precio), precio igual al del plan, una sola plantilla base por comentario, apertura, cierre y complemento compatibles con la categoría y huecos que concuerdan dentro de la oración. Dice `ok` si no detectó nada o el problema encontrado. La columna *Tipo* dice de qué clase de plantilla base viene el comentario: `base`, `contexto cruzado` (vocabulario de otra categoría usado en el sentido de esta) o `compartida`. El ruido intencional (erratas, abreviaturas) puede ocultar la plantilla a la coincidencia exacta del script (`?`).

## Falta de Servicio

| # | Comentario | Auto | Tipo | OK | Observación |
|---|---|---|---|---|---|
| 1 | El asesor fue amable, pero sigo sin señal después de 5 visitas del técnico. Espero una respuesta pronto. | ok | contexto cruzado | | |
| 2 | Cambié de SIM y quedé incomunicado, sin servicio de llamadas desde hace quince días. Por favor, avísenme cualquier novedad. | ok | base | | |
| 3 | hola. El asesor fue amable, pero sigo sin servicio de mensajes después de 3 visitas del técnico. Esto afecta a toda mi familia. Espero una respuesta pronto. | ok | contexto cruzado | | |
| 4 | La neta, llevo más de 48 horas sin internet y sigo pagando mi plan Nocturno 79. Ya probé todo lo que me dijeron y sigue igual. | ok | base | | |
| 5 | Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace cinco días. | ok | base | | |
| 6 | Mi chip nuevo no activa, llevo varias horas esperando que me den servicio de llamadas. | ok | base | | |
| 7 | Buenas tardes. La verdad, llevo cuatro días sin datos móviles y sigo pagando mi plan Familiar 549. Ya me había pasado antes este mes. | ok | base | | |
| 8 | El internet de mi plan Premium 999 no funciona desde enero y no hay quien lo repare. Por favor, avísenme cualquier novedad. | ok | base | | |
| 9 | Mi línea está suspendida sin motivo y no puedo usar nada. | ok | base | | |
| 10 | Ya reporté la falla 5 veces y nadie de soporte me resuelve. Necesito que me den una solución definitiva. | ok | compartida | | |
| 11 | La neta, el internet está tan lento en fines de semana que es como no tener internet. Quedo atento a su respuesta. | ok | base | | |
| 12 | Buenos días. Gracias, TelArg: la factura llega puntual, pero el internet no funciona desde agosto. | ok | contexto cruzado | | |
| 13 | Qué barbaridad, otra vez sin datos móviles, ya es la tercera vez este mes que se cae. | ok | base | | |
| 14 | Hola. La portabilidad de mi número no se completó y estoy sin línea desde hace cuatro días. Ya me había pasado antes este mes. | ok | base | | |
| 15 | En mi casa nunca hay señal, tengo que salir a la calle para poder llamar. | ok | base | | |
| 16 | El asesor fue amable, pero sigo sin servicio de mensajes después de 4 visitas del técnico. Ya me está afectando en el trabajo. Gracias. | ok | contexto cruzado | | |
| 17 | Ayer hubo una caída general de la red en Saltillo y hoy todavía sigue fallando todo. Gracias. | ok | base | | |
| 18 | Las llamadas se cortan cada minuto y el servicio no funciona bien desde hace más de 48 horas. Ya me había pasado antes este mes. | ok | base | | |
| 19 | Es increible la falta de atencion, mi negocio sigue sin internet desde hace mas de 48 horas. Espero una respuesta pronto. | ok | contexto cruzado | | |
| 20 | LES ESCRIBO POR UNA FALLA EN MI LÍNEA. NO PUEDO MANDAR NI RECIBIR MENSAJES DE TEXTO EN PUEBLA DESDE HACE UNA SEMANA. YA ME HABIA PASADO ANTES ESTE MES... | ok | base | | |
| 21 | Es increíble la falta de atención, mi negocio sigue sin internet desde hace tres días. Ya perdí la cuenta de las veces que he insistido. Muchas gracias por su atención. | ok | contexto cruzado | | |
| 22 | Mi chip nuevo no activa, llevo diez días esperando que me den datos móviles. Quiero que alguien me explique qué está pasando. Quedo atento a su respuesta. | ok | base | | |
| 23 | Buenos días. Pagué mi plan Empresarial 899 puntual y aun así me suspendieron el servicio. Ya probé todo lo que me dijeron y sigue igual. Gracias. | ok | base | | |
| 24 | En mi trabajo nunca hay señal, tengo que salir a la calle para poder llamar. Muchas gracias por su atención. | ok | base | | |
| 25 | Les escribo desde Culiacán. Gracias, TelArg: la factura llega puntual, pero el internet no funciona desde octubre. Espero que lo atiendan pronto. | ok | contexto cruzado | | |
| 26 | El internet de mi plan Ahorro 99 no funciona desde julio y no hay quien lo repare. Espero su pronta respuetsa. | ok | base | | |
| 27 | O sea, el internet está tan lento toda la tarde que es como no tener internet. Ya me había pasado antes este mes. | ok | base | | |
| 28 | Buenas noches. Mi tel dice sin servicio desde hace diez dias, ya reinicié todo y nada. Quiero que alguien me explique qué está pasando | ok | base | | |
| 29 | Hice mi recarga de 499 pewos hace una semana y sigo sin internet... | ok | base | | |
| 30 | Desde hace quince días estoy sin servicio de mensajes en mi trabajo y nadie me da una solución. Ya me había pasado antes este mes. | ok | base | | |
| 31 | Tengo la línea a mi nombre. Qué buena atención, ya van 5 días sin datos móviles y nadie me da fecha de solución. Quedo atento a su respuesta. | ok | contexto cruzado | | |
| 32 | Híjole, llevo tres días sin datos móviles y sigo pagando mi plan Empresarial 899. | ok | base | | |
| 33 | No manches, otra vez sin datos móviles, ya es la tercera vez este mes que se cae. Sigo pagando completo y no obtengo lo que corresponde. | ok | base | | |
| 34 | No me llegan los códigos de verificación por SMS desde hace varias horas, creo que hay una falla en la red. Gracias. | ok | base | | |
| 35 | No tengo señal en mi casa en horario laboral, ya pedí que lo revisen (ref. caso 8205023). Ya perdí la cuenta de las veces que he insistido. Gracias de antemano. | ok | base | | |
| 36 | Buenas noches. Me prometieron restablecer el servicio en 10 días y siguen sin cumplir. | ok | base | | |
| 37 | Soy cliente de TelArg desde hace años. En mi barrio no hay cobertura de madrugada, parece que la antena está caída. Nadie le da seguimiento a esto. Por favor, avísenme cualquier novedad. | ok | base | | |
| 38 | Ya Cambie Mi Tel Y Sigo Sin señal, El Problema Es De La red. espero Que lo Atiendan prono. espero Su Pronta Respuesta. | ok | base | | |
| 39 | dicen que hay mantenimiento en Aguascalientes pero ya llevamos una semana y seguimos sin servicio de llamadas. | ok | base | | |
| 40 | No puedo mandar ni recibir mensajes de texto en Villahermosa desde hace quince días. | ok | base | | |
| 41 | Buenas tardes. No puedo mandar ni recibir mensajes de texto en Puebla desde hace tres días. | ok | base | | |
| 42 | Me prometieron restablecer el servicio en 10 días y siguen sin cumplir. Por favor, avísenme cualquier novedad. | ok | base | | |
| 43 | Cuando intento llamar solo escucho que el número no esta disponible, ya llevo diez dias así. Quiero que alguien se haga cargo!!! | ok | base | | |
| 44 | Hola. Dicen que hay mantenimiento en Monterrey pero ya llevamos dos días y seguimos sin cobertura. Quedo atento a su respuesta. | ok | base | | |
| 45 | Ya reporte la falla 2 veces y nadie de soporte me resuelve. | ok | compartida | | |
| 46 | La neta, llevo dos dias sin señal y sigo pagando mi plan Social 299. Ya probe todo lo que me dijeron y sigue igual. | ok | base | | |
| 47 | Qué buena atención, ya van 30 días sin datos móviles y nadie me da fecha de solución. Esto afecta a toda mi familia. Gracias. | ok | contexto cruzado | | |
| 48 | El técnico nunca llegó a revisar la falla de internet, ya llevo más de 48 horas de espera. | ok | base | | |
| 49 | Buenas tardes. No tengo datos móviles en mi oficina en las noches, ya pedí que lo revisen (ref. reporte 548391). Nadie le da seguimiento a esto. Saludos cordiales. | ok | base | | |
| 50 | Hay pérdida de cobertura toda la tarde en toda la zona, y ya son varios días con el mismo problema. | ok | base | | |

## Cobros Indebidos

| # | Comentario | Auto | Tipo | OK | Observación |
|---|---|---|---|---|---|
| 1 | Me quitaron el saldo de mi recarga de $499 MXN sin explicación alguna. | ok | base | | |
| 2 | Pedí la devolución de $169 MXN hace una semana y todavía no me abonan nada. Espero su pronta respuesta. | ok | base | | |
| 3 | Gracias x el cargo no reconocido de $299.00, ni idea de q es. Ya revise mi cuenta y no hay explicacion alguna. Saludos cordiales. | ok | contexto cruzado | | |
| 4 | Me cobraron $99 MXN x servicio premium q nunca contraté. Es dinero que me hace falta | ok | base | | |
| 5 | Desde que cambié de plan a Empresarial 899 me cobran $799 de más, es un cobro indebido. Ya revisé mi cuenta y no hay explicación alguna. | ok | base | | |
| 6 | Me activaron renta de equipo sin avisarme y ahora me cobran $229.00. Nunca me avisaron de este cargo. | ok | base | | |
| 7 | Después de cancelar mi plan Premium 999 me cargaron $999.00 a mi tarjeta. Quiero una explicación clara de este cargo. | ok | base | | |
| 8 | Activaron suscripción de musica sin pedirmelo y me cobraron $129.00 en el estado de mi tarjeta. No pienso pagar un cargo q no corresponde. | ok | base | | |
| 9 | Cargos que no reconozco en mayo: $129 por recarga automática. Exijo que los eliminen. | ok | base | | |
| 10 | YA RECLAMÉ EL COBRO DE MEMBRESIA DE JUEGOS EN LA SUCURSAL Y NADIE ME DEVUELVE MI DINERO (#37895). ES DUNERO Q ME HACE FALTA. QUEDO ATENTO A SU RESPUESTA. | ok | base | | |
| 11 | Les escribo desde Hermosillo. Me siguen cobrando la línea que di de baja, ya llevo 8 meses pagando de más. Me urge que lo corrijan cuanto antes. Por favor, avísenme cualquier novedad. | ok | base | | |
| 12 | El monto de este mes no es el q me dijeron x telefono, quiero una aclaración. | ok | compartida | | |
| 13 | Chale, aparece un cargo de $399.00 por suscripción de música y yo no autoricé nada. Quedo al pendiente. | ok | base | | |
| 14 | Tengo la línea a mi nombre. Cargos que no reconozco en junio: 149 pesos por paquete de datos extra. Exijo que los eliminen. Me urge que lo corrijan cuanto antes. Muchas gracias por su atención. | ok | base | | |
| 15 | Fui A aclarar un Cobro de 99 pesos Y en la Sucursal me Dejaron Esperando Sin darme respuesta. quiero Q Me devuelvan mi Dinero. quedo al pendiente. | ok | compartida | | |
| 16 | Después de cancelar mi plan Ahorro 99 me cargaron 99 pesos a mi tarjeta. Ya pedí una aclaración y sigo sin respuesta. Gracias. | ok | base | | |
| 17 | Mi plan Joven 149 cuesta $149 MXN y me están cobrando $169, no es lo q acordamos. Quiero una explicación clara de este cargo. | ok | base | | |
| 18 | Les escribo desde Monterrey. El monto de este mes no es el que me dijeron en la app quiero una aclaracion. Quiero que me devuelvan mi dinero. | ok | compartida | | |
| 19 | increible, cancele por mala atencion y me siguen cobrando $349.50 cada mes en la tarjeta. ya insisti varias veces y nadie me resuelve. | ok | contexto cruzado | | |
| 20 | Ya reclame el cobro de contenido de terceros x telefono y nadie me devuelve mi dinero (#2692734). | ok | base | | |
| 21 | El asesor me dijo que no me cobrarían más y en noviembre me volvió a salir el cargo de MXN 599. Me urge que lo corrijan cuanto antes. Gracias de antemano. | ok | contexto cruzado | | |
| 22 | Mi plan Familiar 549 cuesta MXN 549 y me están cobrando $999 MXN no es lo que acordamos. Quiero una explicacion clara de este cargo. | ok | base | | |
| 23 | Fui a aclarar un cobro de $99 y en el área de cobranza me dejaron esperando sin darme respuesta. | ok | compartida | | |
| 24 | Buenas noches. Desconozco el cobro de $129.00 del día 5 de mayo, quiero que lo reviertan. Es muy frustrante que pase esto con mi dinero. | ok | base | | |
| 25 | Desconozco el cobro de $349 del dia 10 de agosto, quiero que lo reviertan. Muchas gracias por su atencion. | ok | base | | |
| 26 | Quiero q me reembolsen los $1,499 MXN que me cobraron de más (folio 3408217). Exijo q lo revisen cuanto antes. | ok | base | | |
| 27 | Pedí la devolución de $49.00 hace tres días y todavía no me abonan nada. | ok | base | | |
| 28 | Quiero que me reembolsen los $699.50 que me cobraron de más (#90347). Gracias de antemano. | ok | base | | |
| 29 | Les escribo porque necesito ayuda. Pagué mi mensualidad completa en el portal web y me volvieron a cobrar $169.50. | ok | base | | |
| 30 | En mi tarjeta aparece un cargo de 199 pesos de TelArg que desconozco. | ok | base | | |
| 31 | La neta, me robaron 449 pesos de saldo en un solo día sin hacer llamadas ni usar datos. | ok | base | | |
| 32 | Buenas noches. Desde que camibe de plan a Aigo Sin Limite me cobran $99.00 de más, es un cobro indebido. | ok | base | | |
| 33 | Buenos días. Cargos que no reconozco en agosto: MXN 399 por roaming internacional. Exijo que los eliminen. Me urge que lo corrijan cuanto antes. Gracias de antemano. | ok | base | | |
| 34 | Me quitaron el saldo de mi recarga de $449.50 sin explicacion alguna. | ok | base | | |
| 35 | Buenas tardes. Me activaron plan adicional sin avisarme y ahora me cobran 79 pesos. Esto ya me pasó antes y nunca lo arreglaron. Quedo al pendiente. | ok | base | | |
| 36 | Me cobraron $599.50 por una promoción que nunca acepté (ref. ticket 9195662). | ok | base | | |
| 37 | Ya basta, aparece un cargo de MXN 499 por renta de equipo y yo no autoricé nada. Esto ya me pasó antes y nunca lo arreglaron. | ok | base | | |
| 38 | En mi tarjeta aparece un cargo de 169 pesos de TelArg que desconozco. Es dinero que me hace falta. Espero una respuesta pronto. | ok | base | | |
| 39 | En mi tarjeta aparece un cargo de $79.00 de TelArg que desconozco. | ok | base | | |
| 40 | Ya di de baja la línea y me siguen cobrando el plan Libre 199, la atención no sirve. Nadie me da una respuesta concreta. | ok | contexto cruzado | | |
| 41 | Les escribo porque necesito ayuda. No manches, aparece un cargo de MXN 129 por recarga automática y yo no autoricé nada. Ya revisé mi cuenta y no hay explicación alguna. | ok | base | | |
| 42 | Me activaron contenido de terceros sin avisarme y ahora me cobran $99 MXN. Gracias. | ok | base | | |
| 43 | Después de cancelar mi plan Joven 149 me cargaron MXN 149 a mi tarjeta. | ok | base | | |
| 44 | Se me cobró una penalización de $129.00 que no corresponde porque ya cumplí mi contrato. | ok | base | | |
| 45 | Tengo un problema con un cargo en mi cuenta. Me aplicaron un cargo por reconexión de $2,499 cuando mi línea nunca estuvo suspendida. Quiero una explicación clara de este cargo. | ok | base | | |
| 46 | Cancelé mi servicio porque nunca tuve señal y me siguen haciendo cobros de $199.00 a mi tarjeta. Es muy frustrante que pase esto con mi dinero. Por favor, avísenme cualquier novedad. | ok | contexto cruzado | | |
| 47 | Sin mi consentimiento me dieron de alta en renta de equipo y ya van 12 cobros. | ok | base | | |
| 48 | Me siguen cobrando la línea que di de baja, ya llevo 6 meses pagando de más. | ok | base | | |
| 49 | Me siguen cobrando la línea que di de baja, ya llevo 3 meses pagando de más. Nunca me avisaron de este cargo. | ok | base | | |
| 50 | Pedi la devolucion de $249 hace cuatro dias y todavia no me abonan nada. Ya pedi una aclaracion y sigo sin respuesta. | ok | base | | |

## Problemas de Facturación

| # | Comentario | Auto | Tipo | OK | Observación |
|---|---|---|---|---|---|
| 1 | Tengo un problema con mi factura. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. Necesito que esto quede aclarado con urgencia. Gracias. | ok | base | | |
| 2 | Quiero cambiar mi factura a entrega en la app y no encuentro la opción por ningún lado. | ok | base | | |
| 3 | a quien corresponda. Necesito aclarar mi estado de cuenta, hay conceptos que no entiendo y el total no cuadra. | ok | base | | |
| 4 | Tengo la línea a mi nombre. Mi pago de $79 sí se hizo pero la factura todavía marca vencido. | ok | base | | |
| 5 | Buenas tardes. Gracias por nada, la factura de septiembre llegó otra vez con el plan equivocado. Necesito que lo resuelvan cuanto antes. | ok | contexto cruzado | | |
| 6 | Saludos. Mi fecha de corte cambió al día 15 sin avisarme y ahora el estado de cuenta viene desfasado. | ok | base | | |
| 7 | Mi dirección de facturación es incorrecta y por eso no llegan mis recibos. Ya lo reporté y sigue igual. Espero una respuesta pronto. | ok | base | | |
| 8 | falta el folio fiscal en mi factura de enero y no puedo deducirla. | ok | base | | |
| 9 | Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. Quedo al pendiente. | ok | base | | |
| 10 | Los consumos en mi factura de mayo no corresponden con mi historial de llamadas. Muchas gracias por su atención. | ok | base | | |
| 11 | La factura de septiembre viene con el plan equivocado, aparece Empresarial 899 y yo contraté otro. Ya me había pasado antes y nadie lo arregló de fondo. Espero su pronta respuesta. | ok | base | | |
| 12 | Mi fact de julio dice que tego adeudo pero el sistema de pagos dice q estoy al corriente. Espero una respuesta pronto... | ok | base | | |
| 13 | Me cobran el monto completo, pero la factura dice un descuento que nunca aplicaron | ok | contexto cruzado | | |
| 14 | El monto de este mes no es el que me dijeron en la sucursal, quiero una aclaración. Quedo atento a su respuesta. | ok | compartida | | |
| 15 | Tengo un problema con mi factura. El CFDI de enero no llega y el cobro sí se hizo. Ya intenté arreglarlo por mi cuenta y no pude. | ok | contexto cruzado | | |
| 16 | Mi fecha de corte cambió al día 28 sin avisarme y ahora el estado de cuenta viene desfasado. Por favor, avísenme cualquier novedad. | ok | base | | |
| 17 | El estado de cuenta de marzo no detalla los conceptos, solo trae un total de 199 pesos. | ok | base | | |
| 18 | Les escribo porque necesito ayuda. El recibo no trae el descuento que me prometieron en atención a clientes por la falla de servicio de llamadas. Es muy molesto andar aclarando esto. | ok | contexto cruzado | | |
| 19 | Recibo mi factura 10 días tarde cada mes, ya casi vence cuando la tengo. | ok | base | | |
| 20 | La factura de octubre viene con el plan equivocado, aparece Ahorro 99 y yo contraté otro. | ok | base | | |
| 21 | Mi fecha de corte cambió al día 15 sin avisarme y ahora el estado de cuenta viene desfasado. | ok | base | | |
| 22 | Mi factura llega con los datos de mi domicilio anterior aunque ya lo actualicé. | ok | base | | |
| 23 | Solicité factura electrónica por WhatsApp y me siguen enviando el recibo impreso. Quedo atento a su respuesta. | ok | base | | |
| 24 | Me cobran el monto completo, pero la factura dice un descuento que nunca aplicaron. Espero su pronta respuesta. | ok | contexto cruzado | | |
| 25 | tengo un problema con mi fact. el cfdi de octubre no llega y el cobro si se hizo... | ok | contexto cruzado | | |
| 26 | Los consumos en mi factura de octubre no corresponden con mi historial de llamadas. | ok | base | | |
| 27 | Me llego la factura en blanco, sin importes, desde que active la entrega en papel. | ok | base | | |
| 28 | El recibo no trae el descuento Que Me Prometieron En atencion a clientes Por La falla De Internet. muchas gracias por su atencion. | ok | contexto cruzado | | |
| 29 | Tengo la línea a mi nombre. La factura trae un saldo anterior que ya pagué, el pago de $449.00 no aparece reflejado. Espero que lo resuelvan pronto. | ok | base | | |
| 30 | Les escribo desde Chihuahua. La factura de mi línea empresarial llega a otra razón social, ya mandé los documentos en el portal web. | ok | base | | |
| 31 | No puedo registrar mi correo para facturacion, la pagina me regresa un error (ref. caso 9138778). Necesito que lo resuelvan cuanto antes. Saludos cordiales. | ok | base | | |
| 32 | El monto de este mes no es el que me dijeron x correo quiero una aclaracion. No me gusta quedarme con dudas de lo q pago. | ok | compartida | | |
| 33 | Mi fecha de corte cambió al día 28 sin avisarme y ahora el estado de cuenta viene desfasado. Espero que lo resuelvan pronto. | ok | base | | |
| 34 | el estado de cuenta de junio no detalla los conceptos, solo trae un total de mxn 1,199. | ok | base | | |
| 35 | Buenos días. Mi plan Max 699 cuesta $699, pero en mi factura el monto total es de $999.50. | ok | base | | |
| 36 | El desglose de mi factura no incluye los impuestos y no entiendo cómo llegan a 399 pesos. | ok | base | | |
| 37 | Mi fact sale a nombre de otra persona y necesito que corrijan los datos. Me urge tenerlo en regla. | ok | base | | |
| 38 | Falta el folio fiscal en mi factura de febrero y no puedo deducirla. Espero que lo resuelvan pronto. | ok | base | | |
| 39 | Buenas noches. El RFC de mi factura está mal capturado, necesito una factura corregida. Muchas gracias por su atención. | ok | base | | |
| 40 | Necesito el comprobante fiscal (CFDI) de noviembre y el sistema no me lo genera (ref. reporte 760569). Necesito que lo resuelvan cuanto antes. | ok | base | | |
| 41 | Buenas noches. Me cobran el monto completo, pero la fact dice un descuento que nunca aplicaron. Espero que lo reshelvan pronto. | ok | contexto cruzado | | |
| 42 | Gracias por nada, la factura de abril llegó otra vez con el plan equivocado. Ya intenté arreglarlo por mi cuenta y no pude. Quedo al pendiente. | ok | contexto cruzado | | |
| 43 | Mi recibo viene con fecha limite de pago distinta a la q me indicaron en el poratl web. Necesito q lo resuelvan cuanto antes. Quedo atento a su respuesta. | ok | base | | |
| 44 | mi recibo no muestra el descuento prometido en soporte técnico, solo el total de MXN 129. Espero que lo resuelvan pronto. Quedo al pendiente. | ok | base | | |
| 45 | Mi pago de MXN 99 sí se hizo pero la factura todavía marca vencido. | ok | base | | |
| 46 | A quien corresponda. Necesito el comprobante fiscal (CFDI) de enero y el sistema no me lo genera (ref. folio 4979469). | ok | base | | |
| 47 | La factura de septiembre viene con el plan equivocado, aparece Amigo Sin Límite y yo contraté otro. Espero una respuesta pronto. | ok | base | | |
| 48 | EL ESTADO DE CUENTA DE FEBRERO NO DETALLA LOS CONCEPTOS, SOLO TRAE UN TOTAL DE 99 PESOS. | ok | base | | |
| 49 | Buenas tardes. Falta el folio fiscal en mi factura de marzo y no puedo deducirla. Me urge tenerlo en regla. | ok | base | | |
| 50 | La factura Dice Servicio Activo, pero En mi Ofiina No Tengo Internet Desde hace Dias. | ok | contexto cruzado | | |

## Quejas

| # | Comentario | Auto | Tipo | OK | Observación |
|---|---|---|---|---|---|
| 1 | El equipo que me entregaron llegó tarde y con la caja abierta, muy mal servicio. | ok | base | | |
| 2 | Esperé media hora en ventas y nadie me contestó. | ok | base | | |
| 3 | a quien corresponda. Para cancelar te piden mil cosas y mientras tanto te siguen cobrando el plan Social 299. Espero su pronta respuesta. | ok | contexto cruzado | | |
| 4 | Mi queja: el call center cuelga cuando pido hablar con un supervisor. | ok | base | | |
| 5 | Llevo varias horas esperando una solución y nadie me da respuesta. Quedo atento a su respuesta. | ok | compartida | | |
| 6 | Me llaman seguido en horario laboral para ofrecerme promociones aunque ya les pedí que no. Espero que alguien por fin me responda. | ok | base | | |
| 7 | Publicidad engañosa: prometen cobertura en Culiacan y no hay señal. Que alguien se haga responsable, x favor. | ok | contexto cruzado | | |
| 8 | Publicidad engañosa: prometen cobertura en Chihuahua y no hay señal. Pido que tomen en cuenta mi queja. | ok | contexto cruzado | | |
| 9 | Me cortan por un adeudo mínimo cuando ellos no mandan bien la factura. | ok | contexto cruzado | | |
| 10 | Llevo más de 48 horas esperando una solución y nadie me da respuesta. | ok | compartida | | |
| 11 | De plano, no es posible que una empresa tan grande tenga una atencion tan lenta. Ya perdí demasiado tiempo en esto. | ok | base | | |
| 12 | Pedí 4 veces que corrijan mi factura y en ventas nadie me hace caso. | ok | compartida | | |
| 13 | Una verguenza la atención: me dejaron esperando una hora y se me corto la llamada. Espero su pronta respuesta. | ok | base | | |
| 14 | Buenas tardes. Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Pido que tomen en cuenta mi queja. | ok | base | | |
| 15 | Me vendieron el plan Amigo Sin Límite con condiciones que no cumplen, publicidad engañosa. Espero que alguien por fin me responda. | ok | base | | |
| 16 | Les escribo desde Cancún. Chale, no es posible que una empresa tan grande tenga una atención tan lenta. Pido que tomen en cuenta mi queja. | ok | base | | |
| 17 | Buenas noches. Fui a aclarar un cobro de 199 pesos y en atencion a clientes me dejaron esperando sin darme respuesta. Estoy pensando seriamente en cambiarme de compañia. Muchas gracias por su atencion. | ok | compartida | | |
| 18 | Buenos días. Qué onda, qué mala atención, desde noviembre estoy tratando de aclarar un tema y nadie me resuelve. | ok | base | | |
| 19 | A quien corresponda. No hay quien resuelva dudas en atención a clientes, todo son excusas y tiempos de espera. | ok | base | | |
| 20 | EL ASESOR DE LA SUCURSAL DE SAN LUIS POTOSI FUE GROSERO Y ME CORTO LA LINEA X PAGAR UN DIA TARDE. SALUDOS CORDIALES. | ok | contexto cruzado | | |
| 21 | Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. | ok | base | | |
| 22 | Los asesores no saben nada, me dieron 2 respuestas distintas en la misma semana. | ok | base | | |
| 23 | El proceso para cambiar de plan a Joven 149 fue lentísimo y lleno de requisitos absurdos. Todo esto es muy frustrante. Gracias. | ok | base | | |
| 24 | Exijo que me atiendan como cliente, no es la primera vez que me ignoran por chat. Muchas gracias por su atención. | ok | base | | |
| 25 | Me vendieron el plan Prepago Básico con condiciones que no cumplen, publicidad engañosa. Que alguien se haga responsable, por favor. | ok | base | | |
| 26 | Quiero dejar constancia de mi inconformidad. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. Que alguien se haga responsable, por favor. | ok | base | | |
| 27 | Tengo una inconformidad q presentar. Las promociones que anuncian no aplican nunca para clientes que ya tenemos años con ustedes. | ok | base | | |
| 28 | Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Todo esto es muy frustrante. | ok | base | | |
| 29 | Tengo una inconformidad q presentar. En atencion a clientes me dijeron q no era su responsabilidad y me mandaron a otra area. | ok | base | | |
| 30 | Exijo que me atiendan como cliente, no es la primera vez que me ignoran por correo. | ok | base | | |
| 31 | Las condiciones del contrato cambian cada mes y nadie me las explica con claridad. Gracias. | ok | base | | |
| 32 | El proceso para cambiar de plan a Familiar 549 fue lentísimo y lleno de requisitos absurdos. Estoy pensando seriamente en cambiarme de compañía. | ok | base | | |
| 33 | Quiero presentar mi queja porque Héctor me colgó la llamada en plena aclaración. Quiero que alguien me explique por qué. Saludos cordiales. | ok | base | | |
| 34 | Pedí una aclaración por escrito y solo me mandan respuestas copiadas y pegadas. | ok | base | | |
| 35 | estoy muy molesto con la forma en que fernanda me atendio, sin ninguna empatia. | ok | base | | |
| 36 | Mi queja: el call center cuelga cuando pido hablar con un supervisor. | ok | base | | |
| 37 | Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. | ok | base | | |
| 38 | Los asesores no saben nada, me dieron 3 respuestas distintas en la misma semana. Espero q alguien con autoridad me atienda. | ok | base | | |
| 39 | Subieron el precio de mi plan Premium 999 sin avisarme, esto me parece un abuso. | ok | base | | |
| 40 | No manches, el trato de Miriam en el área de cobranza fue grosero, no pienso aguantar eso. Espero una respuesta pronto. | ok | base | | |
| 41 | Una vergüenza la atención: me dejaron esperando media hora y se me cortó la llamada. Todo esto es muy frustrante. | ok | base | | |
| 42 | Quiero presentar mi queja porque Fernanda me colgó la llamada en plena aclaración. | ok | base | | |
| 43 | mala experiencia en la sucursal: tuve que repetir mi problema a cuatrro personas distintas. no es justo q uno tenga que insistir tanto | ok | base | | |
| 44 | Otra vez me dejaron sin servicio y en atención a clientes solo me dan largas (ref. reporte 61940). Llevo años con ustedes y no merezco esto. | ok | compartida | | |
| 45 | Para cancelar mi línea me ponen mil trabas, es un proceso interminable. Llevo años con ustedes y no merezco esto. | ok | base | | |
| 46 | Buenas tardes. Me cortan por un adeudo mínimo cuando ellos no mandan bien la factura. | ok | contexto cruzado | | |
| 47 | Nadie da seguimiento a lo que pedí (ref. TKT-19656), siempre me dicen que ya está en proceso. Ya insistí varias veces y nadie me resuelve. | ok | base | | |
| 48 | Una vergüenza la atención: me dejaron esperando veinte minutos y se me cortó la llamada. Llevo años con ustedes y no merezco esto. Saludos cordiales. | ok | base | | |
| 49 | Buenas noches. Cada vez que llamo me piden los mismos datos y no resuelven nada, estoy harto del servicio de atención. Saludos cordiales. | ok | base | | |
| 50 | LA ATENCION EN LA SUCURSAL DE SAN LUIS POTOSÍ ES PÉSIMA HACEN FILA UNA HORA PARA NADA. TODO ESTO ES MUY FRUSTRANTE. QQUEDO ATENTO A SU RESPUESTA. | ok | base | | |

## Reconocimientos

| # | Comentario | Auto | Tipo | OK | Observación |
|---|---|---|---|---|---|
| 1 | Les escribo con mucho gusto. La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. Realmente superó mis expectativas. | ok | base | | |
| 2 | LAS MEJORES PROMOCIONES DEL MERCADO, SIGO CON TELARG X EL BUEN SERVICIO. SIGAN ASI. | ok | base | | |
| 3 | Las mejores promociones del mercado, sigo con TelArg por el buen servicio. Realmente superó mis expectativas. Que tengan un excelente día. | ok | base | | |
| 4 | Excelente cobertura en Veracruz, nunca he tenido problemas con mi plan Navega 249. Seguiré recomendándolos. | ok | base | | |
| 5 | Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. caso 5692106), todo quedó resuelto. Realmente superó mis expectativas. | ok | base | | |
| 6 | Les felicito por la aplicación, es muy fácil de usar y puedo ver mis consumos al instante. | ok | base | | |
| 7 | Gracias a Jorge por la paciencia y la amabilidad al explicarme mi plan Max 699. Por esto sigo con TelArg. | ok | base | | |
| 8 | Pagar mi recibo en la app es muy fácil, todo queda registrado al momento. | ok | base | | |
| 9 | Llevo años con TelArg y siempre he tenido buen servicio, los recomiendo. | ok | base | | |
| 10 | Mi señal es excelente incluso en mi colonia, felicidades a todo el equipo de redes. Realmente supero mis exprctativas. | ok | base | | |
| 11 | El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. | ok | base | | |
| 12 | Mi señal es excelente incluso en mi casa, felicidades a todo el equipo de redes. | ok | base | | |
| 13 | buenas Tardes. andale, Qué Buen Servicio, en Menos De Media Hora Me solucionaron todo. da Gusto ver Q hay Empresas Q Cumplen. saludos cordiales... | ok | base | | |
| 14 | El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias. Voy a seguir con ustedes muchos años más. | ok | base | | |
| 15 | buenos dias. mi señal es excelente incluso en mi barrio, felicidades a todo el equipo de redes. | ok | base | | |
| 16 | Gracias x la llamada de seguimiento, nunca me habia pasado q una empresa se preocupara asi. Eso es lo q uno espera de una buena compañía!!! | ok | base | | |
| 17 | Quiero compartir mi experiencia. Muchas gracias por resolver mi problema por teléfono, el servicio fue impecable. Ojalá todas las empresas fueran así. | ok | base | | |
| 18 | mi señal es excelente incluso en mi domicilio, felicidades a todo el equipo de redes. | ok | base | | |
| 19 | Buenas tardes. Gracias por restablecer mi servicio el mismo día, atención de primera. De verdad quedé muy satisfecho. | ok | base | | |
| 20 | La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. | ok | base | | |
| 21 | Las mejores promociones del mercado, sigo con TelArg x el buen servicio. | ok | base | | |
| 22 | Buenas tardes. Mi señal es excelente incluso en toda la zona, felicidades a todo el equipo de redes. | ok | base | | |
| 23 | Quiero dejar constancia de lo bien que me atendió Fernanda, un trato impecable. Un saludo a todo el equipo. | ok | base | | |
| 24 | Gracias x la llamada de seguimiento, nunca me habia pasado q una empresa se preocupara asi. | ok | base | | |
| 25 | Me encanta que la factura llegue cada mes por WhatsApp sin falta. Que tengan un excelente día. | ok | base | | |
| 26 | Quiero dejar constancia de lo bien que me atendió Itzel, un trato impecable. Sigan así, de verdad se agradece. | ok | base | | |
| 27 | Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. | ok | base | | |
| 28 | Las mejores promociones del mercado, sigo con TelArg por el buen servicio. | ok | base | | |
| 29 | El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Voy a seguir con ustedes muchos años más. | ok | base | | |
| 30 | Quiero dejar un comentario positivo. La promoción del plan Empresarial 899 es muy buena, gracias por pensar en los clientes. Ojalá todas las empresas fueran así. | ok | base | | |
| 31 | Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. TKT-467455), todo quedó resuelto. | ok | base | | |
| 32 | Les escribo con una buena noticia. Quiero felicitar a Daniela en el call center, el trato fue de lo mejor. Ojalá todo se mantenga así. Muchas gracias por su atención. | ok | base | | |
| 33 | Mi señal es excelente incluso en mi colonia, felicidades a todo el equipo de redes. | ok | base | | |
| 34 | Les felicito por la aplicacion, es muy facil de usar y pusdo ver mis consumos al instante. Realmente supero mis expectativas. Seguire reocmendándolos | ok | base | | |
| 35 | Al final si me resolvieron en el call center, aunque tardaron bastante!!! | ok | compartida | | |
| 36 | El servicio de TelArg ha mejorado muchísimo, se nota el compromiso con los clientes. Ojalá todas las empresas fueran así. | ok | base | | |
| 37 | Por fin me devolvieron el cobro doble de $999.50, el asesor fue muy amable. Saludos cordiales. | ok | contexto cruzado | | |
| 38 | Las mejores promociones del mercado, sigo con TelArg por el buen servicio. Se los quería hacer saber. | ok | base | | |
| 39 | Excelente cobertura en Saltillo, nunca he tenido problemas con mi plan Social 299. Ojalá todas las empresas fueran así. Sigan así. | ok | base | | |
| 40 | El equipo que me entregaron llegó antes de lo prometido y en perfecto estado, muchas gracias. | ok | base | | |
| 41 | Al final sí me resolvieron en la sucursal, aunque tardaron bastante. Gracias. | ok | compartida | | |
| 42 | Gracias por avisarme antes de la fecha de corte, así no se me pasa el pago. Se los quería hacer saber. Muchas gracias por su atención. | ok | base | | |
| 43 | La domiciliación de mi pago funciona perfecto, nunca he tenido un cargo de más. Eso es lo que uno espera de una buena compañía. | ok | base | | |
| 44 | Un diez para atencion a clientes: resolvieron mi duda sin darme vueltas. Sigan asi, de verdad se agradece | ok | base | | |
| 45 | A quien corresponda. Quiero agradecer la rapidez con la que atendieron mi solicitud (ref. #8106181), todo quedó resuelto. | ok | base | | |
| 46 | En soporte técnico me explicaron con paciencia el cobro de mi plan Libre 199, quedé muy conforme. Eso es lo que uno espera de una buena compañía. | ok | base | | |
| 47 | Buenos días. Mi reconocimiento al técnico que llegó puntual y dejó mi internet funcionando perfecto. | ok | base | | |
| 48 | Gracias a Carlos por la paciencia y la amabilidad al explicarme mi plan Social 299. Ojalá todas las empresas fueran así. Muchas gracias por su atención. | ok | base | | |
| 49 | gracias, telarg: despues de dos dias sin internet lo arreglaron en cuanto lo reporte. da gusto ver q hay empresas q cumplen. | ok | contexto cruzado | | |
| 50 | A quien corresponda. Gracias, TelArg: después de diez días sin internet lo arreglaron en cuanto lo reporté. Todo el equipo merece un aplauso. Gracias. | ok | contexto cruzado | | |
