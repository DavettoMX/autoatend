package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// slotVocab holds the domain vocabularies for the plain list-backed slots.
// Every value must read naturally in each frame the templates use it in
// ("sin {SERVICIO}", "con {AREA}", "mi plan {PLAN}", "por {CONCEPTO}", ...).
var slotVocab = map[string][]string{
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

var amountPool = []int{
	49, 79, 99, 129, 149, 169, 199, 229, 249, 299, 349, 399, 449, 499, 549, 599, 699, 799, 899,
	999, 1199, 1499, 1999, 2499,
}

var folioPrefixes = []string{"TKT-", "folio ", "ticket ", "caso ", "reporte ", "#"}

// slotFiller expands {SLOT} placeholders using one rng stream. It remembers the
// last two amounts drawn so that {MONTO_B} is guaranteed to differ from the
// amount before it (e.g. "me cobraron {MONTO} y debía ser {MONTO_B}").
type slotFiller struct {
	rng        *rand.Rand
	amount     int
	prevAmount int
}

func newSlotFiller(rng *rand.Rand) *slotFiller { return &slotFiller{rng: rng} }

func (f *slotFiller) fill(tpl string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(tpl); {
		c := tpl[i]
		if c != '{' {
			b.WriteByte(c)
			i++
			continue
		}
		end := strings.IndexByte(tpl[i:], '}')
		if end < 0 {
			return "", fmt.Errorf("gencorpus: unterminated slot in template %q", tpl)
		}
		name := tpl[i+1 : i+end]
		v, ok := f.value(name)
		if !ok {
			return "", fmt.Errorf("gencorpus: unknown slot {%s} in template %q", name, tpl)
		}
		b.WriteString(v)
		i += end + 1
	}
	return b.String(), nil
}

func (f *slotFiller) value(name string) (string, bool) {
	switch name {
	case "MONTO":
		return f.drawAmount(false), true
	case "MONTO_B":
		return f.drawAmount(true), true
	case "FOLIO":
		return f.folio(), true
	}
	vals, ok := slotVocab[name]
	if !ok {
		return "", false
	}
	return vals[f.rng.Intn(len(vals))], true
}

func (f *slotFiller) drawAmount(differ bool) string {
	n := amountPool[f.rng.Intn(len(amountPool))]
	for differ && n == f.amount {
		n = amountPool[f.rng.Intn(len(amountPool))]
	}
	f.prevAmount, f.amount = f.amount, n
	return formatAmount(n, f.rng.Intn(6))
}

// formatAmount renders n pesos in one of six surface formats ("variable format").
func formatAmount(n, style int) string {
	s := withThousands(n)
	switch style {
	case 0:
		return "$" + s
	case 1:
		return "$" + s + ".00"
	case 2:
		return s + " pesos"
	case 3:
		return "MXN " + s
	case 4:
		return "$" + s + " MXN"
	default:
		return "$" + s + ".50"
	}
}

func withThousands(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	return s[:len(s)-3] + "," + s[len(s)-3:]
}

func (f *slotFiller) folio() string {
	prefix := folioPrefixes[f.rng.Intn(len(folioPrefixes))]
	digits := 5 + f.rng.Intn(3)
	lo := 1
	for i := 1; i < digits; i++ {
		lo *= 10
	}
	return prefix + strconv.Itoa(lo+f.rng.Intn(9*lo))
}
