package templates

// SlotVocab holds the domain vocabularies for the plain list-backed slots.
// Every value must read naturally in each frame the templates use it in
// ("sin {SERVICIO}", "con {AREA}", "mi plan {PLAN}", "por {CONCEPTO}", ...).
var SlotVocab = map[string][]string{
	"PLAN": {
		"Libre 199", "Amigo Sin Límite", "Total 399", "Ahorro 99", "Max 699", "Familiar 549",
		"Prepago Básico", "Joven 149", "Empresarial 899", "Premium 999", "Navega 249",
		"Nocturno 79", "Social 299", "Hogar 449",
	},
	"SERVICIO": {
		"internet", "datos móviles", "señal", "servicio de llamadas", "servicio de mensajes", "cobertura",
	},
	"ZONA": {
		"mi colonia", "mi casa", "mi oficina", "toda mi calle", "el centro", "mi barrio",
		"toda la zona", "mi trabajo", "mi domicilio",
	},
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
	"COLOQ": {
		"la neta", "o sea", "híjole", "órale", "no manches", "qué onda", "chale", "de plano",
		"la verdad", "ya basta", "caray", "ándale", "ay no", "qué barbaridad",
	},
	"CANAL": {
		"por teléfono", "por chat", "en la app", "en la sucursal", "en el portal web", "por correo",
	},
	"ESPERA": {
		"veinte minutos", "cuarenta minutos", "una hora", "más de una hora", "casi dos horas", "media hora",
	},
	"ELOGIO": {
		"excelente", "impecable", "sobresaliente", "de primera", "de diez", "de lo mejor",
	},
	"DIA_CORTE": {"5", "10", "15", "20", "28"},
	"DIAS":      {"3", "5", "7", "10", "15", "30"},
	"ENTREGA": {
		"por correo electrónico", "en papel", "en la app", "por WhatsApp",
	},
	"HORARIO": {
		"por las mañanas", "en las noches", "toda la tarde", "en horario de oficina",
		"de madrugada", "en fines de semana",
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
