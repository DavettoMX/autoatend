package templates

// Plan is a TelArg plan with its monthly price in pesos. The price of the plans
// with a number in their name is that number; the others have a fixed price.
type Plan struct {
	Name  string
	Price int
}

// Plans is the {PLAN} vocabulary. {PRECIO}, {COBRO} and {EXTRA} are derived from the
// plan drawn last in the same template, so every amount agrees with the plan.
var Plans = []Plan{
	{"Libre 199", 199}, {"Amigo Sin Límite", 299}, {"Total 399", 399}, {"Ahorro 99", 99},
	{"Max 699", 699}, {"Familiar 549", 549}, {"Prepago Básico", 99}, {"Joven 149", 149},
	{"Empresarial 899", 899}, {"Premium 999", 999}, {"Navega 249", 249}, {"Nocturno 79", 79},
	{"Social 299", 299}, {"Hogar 449", 449},
}

// SlotVocab holds the domain vocabularies for the plain list-backed slots.
// Every value must read naturally in each frame the templates use it in
// ("sin {SERVICIO}", "con {AREA}", "por {CONCEPTO}", ...). Where a frame accepts only
// part of a vocabulary, the template uses a constrained slot (SERVICIO_DATOS,
// ZONA_AREA, CONCEPTO_SERV, ...) instead of the general one.
var SlotVocab = map[string][]string{
	"SERVICIO": {
		"internet", "datos móviles", "señal", "servicio de llamadas", "servicio de mensajes", "cobertura",
	},
	"SERVICIO_DATOS": {"internet", "datos móviles"},
	"ZONA": {
		"mi colonia", "mi casa", "mi oficina", "toda mi calle", "el centro", "mi barrio",
		"toda la zona", "mi trabajo", "mi domicilio",
	},
	"ZONA_AREA":  {"mi colonia", "toda mi calle", "el centro", "mi barrio", "toda la zona"},
	"ZONA_LUGAR": {"mi casa", "mi oficina", "mi trabajo", "mi domicilio"},
	"CIUDAD": {
		"Monterrey", "Guadalajara", "Puebla", "Mérida", "León", "Querétaro", "Tijuana", "Toluca",
		"Oaxaca", "Morelia", "Chihuahua", "Culiacán", "Veracruz", "Cancún", "Saltillo",
		"Aguascalientes", "Hermosillo", "San Luis Potosí", "Villahermosa", "Tampico",
	},
	"CONCEPTO": {
		"servicio premium", "suscripción de música", "seguro de equipo", "roaming internacional",
		"mensajes premium", "recarga automática", "reconexión", "renta de equipo",
		"membresía de juegos", "contenido de terceros", "plan adicional", "paquete de datos extra",
	},
	"CONCEPTO_SERV": {
		"servicio premium", "suscripción de música", "seguro de equipo", "mensajes premium",
		"recarga automática", "renta de equipo", "membresía de juegos", "contenido de terceros",
		"plan adicional", "paquete de datos extra",
	},
	"MES": {
		"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
	},
	"AREA": {
		"atención a clientes", "soporte técnico", "el área de cobranza", "la sucursal", "el call center", "ventas",
	},
	"AGENTE": {
		"Laura", "Carlos", "Miriam", "Jorge", "Fernanda", "Ricardo", "Daniela", "Héctor", "Paola",
		"Alejandro", "Verónica", "Raúl", "Itzel", "Gerardo",
	},
	"EQUIPO": {
		"mi celular", "mi teléfono", "el módem", "el router", "mi chip", "mi SIM",
	},
	"EQUIPO_PANTALLA": {"mi celular", "mi teléfono"},
	"COLOQ": {
		"la neta", "o sea", "híjole", "no manches", "qué onda", "chale", "de plano",
		"la verdad", "ya basta", "caray", "ay no", "qué barbaridad",
	},
	"COLOQ_POS": {
		"la neta", "órale", "no manches", "la verdad", "ándale", "de plano", "qué padre", "mira nomás",
	},
	"CANAL": {
		"por teléfono", "por chat", "en la app", "en la sucursal", "en el portal web", "por correo",
	},
	"CANAL_DOC":  {"por correo", "en el portal web", "en la app", "en la sucursal", "por chat"},
	"CANAL_PAGO": {"por teléfono", "en la app", "en el portal web", "en la sucursal"},
	"ESPERA": {
		"veinte minutos", "cuarenta minutos", "una hora", "más de una hora", "casi dos horas", "media hora",
	},
	"ESPERA_CORTA": {"cinco minutos", "diez minutos", "quince minutos", "veinte minutos", "media hora"},
	"ELOGIO": {
		"excelente", "impecable", "sobresaliente", "de primera", "de diez", "de lo mejor",
	},
	"DIA_CORTE":    {"5", "10", "15", "20", "28"},
	"DIAS":         {"3", "5", "7", "10", "15", "30"},
	"DIAS_RETRASO": {"3", "5", "7", "10"},
	"MESES_N":      {"2", "3", "4", "6", "8", "12"},
	"VECES":        {"2", "3", "4", "5"},
	"ENTREGA": {
		"por correo electrónico", "en papel", "en la app", "por WhatsApp",
	},
	"ENTREGA_ELEC": {"por correo electrónico", "en la app", "por WhatsApp"},
	"HORARIO": {
		"por las mañanas", "en las noches", "toda la tarde", "en horario laboral",
		"de madrugada", "en fines de semana",
	},
	"DURACION_DIAS": {
		"dos días", "tres días", "cuatro días", "una semana", "más de 48 horas", "cinco días",
		"diez días", "quince días",
	},
	"DURACION": {
		"dos días", "tres días", "cuatro días", "una semana", "más de 48 horas", "cinco días",
		"varias horas", "diez días", "quince días",
	},
}

var AmountPool = []int{
	49, 79, 99, 129, 149, 169, 199, 229, 249, 299, 349, 399, 449, 499, 549, 599, 699, 799, 899,
	999, 1199, 1499, 1999, 2499,
}

var FolioPrefixes = []string{"TKT-", "folio ", "ticket ", "caso ", "reporte ", "#"}
