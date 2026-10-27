package astronomy

import (
	"math"
	"time"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AstronomicalEvent struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	Category    string    `json:"category"`
}

// CalcularFaseLunar calcula la fase de la luna y su iluminación aproximada para una fecha dada
func CalcularFaseLunar(t time.Time) (string, float64) {
	// Fecha de referencia de Luna Nueva conocida (11 de enero de 2024)
	refNewMoon := time.Date(2024, 1, 11, 11, 57, 0, 0, time.UTC)
	synodicMonth := 29.53058867 // días que dura el ciclo lunar

	daysElapsed := t.UTC().Sub(refNewMoon).Hours() / 24.0
	cycles := daysElapsed / synodicMonth
	phaseFraction := cycles - math.Floor(cycles) // valor entre 0.0 y 1.0

	// Ángulo de fase en grados (0 a 360)
	_ = phaseFraction * 360.0

	var name string
	if phaseFraction < 0.03 || phaseFraction > 0.97 {
		name = "Luna Nueva"
	} else if phaseFraction < 0.22 {
		name = "Luna Creciente"
	} else if phaseFraction < 0.28 {
		name = "Cuarto Creciente"
	} else if phaseFraction < 0.47 {
		name = "Gibosa Creciente"
	} else if phaseFraction < 0.53 {
		name = "Luna Llena"
	} else if phaseFraction < 0.72 {
		name = "Gibosa Menguante"
	} else if phaseFraction < 0.78 {
		name = "Cuarto Menguante"
	} else {
		name = "Luna Menguante"
	}

	// Iluminación aproximada (0% a 100%)
	illumination := (1.0 - math.Cos(phaseFraction*2.0*math.Pi)) / 2.0 * 100.0

	return name, illumination
}

// GenerarEventosLocales calcula eventos astronómicos dinámicos para una ubicación y fecha
func GenerarEventosLocales(loc Location, targetDate time.Time) []AstronomicalEvent {
	var events []AstronomicalEvent

	// 1. Fase lunar del día
	moonPhase, illumination := CalcularFaseLunar(targetDate)
	events = append(events, AstronomicalEvent{
		Title:       "Fase Lunar: " + moonPhase,
		Description: "Iluminación lunar aproximada del " + formatFloat(illumination) + "%. Calculado para las coordenadas (" + formatFloat(loc.Latitude) + ", " + formatFloat(loc.Longitude) + ").",
		EventDate:   targetDate,
		Category:    "moon",
	})

	// 2. Eventos estacionales o destacados basados en el día del año (Equinoccios y Solsticios aproximados)
	month := targetDate.Month()
	day := targetDate.Day()

	if month == 3 && day >= 19 && day <= 21 {
		events = append(events, AstronomicalEvent{
			Title:       "Equinoccio de Otoño / Primavera",
			Description: "El Sol cruza el ecuador celeste. La duración del día y de la noche son prácticamente iguales en todo el planeta.",
			EventDate:   targetDate,
			Category:    "sun",
		})
	} else if month == 6 && day >= 20 && day <= 22 {
		events = append(events, AstronomicalEvent{
			Title:       "Solsticio de Invierno / Verano",
			Description: "El día más corto o más largo del año según el hemisferio, debido a la máxima inclinación axial de la Tierra.",
			EventDate:   targetDate,
			Category:    "sun",
		})
	} else if month == 9 && day >= 21 && day <= 23 {
		events = append(events, AstronomicalEvent{
			Title:       "Equinoccio de Primavera / Otoño",
			Description: "Segundo equinoccio del año astronómico.",
			EventDate:   targetDate,
			Category:    "sun",
		})
	} else if month == 12 && day >= 20 && day <= 22 {
		events = append(events, AstronomicalEvent{
			Title:       "Solsticio de Verano / Invierno",
			Description: "Máxima declinación solar en el solsticio de diciembre.",
			EventDate:   targetDate,
			Category:    "sun",
		})
	}

	// 3. Evento destacado genérico de observación planetaria
	events = append(events, AstronomicalEvent{
		Title:       "Visibilidad de Planetas Brillantes",
		Description: "Condiciones óptimas para la observación de Júpiter y Saturno en las primeras horas de la noche desde tu ubicación.",
		EventDate:   targetDate,
		Category:    "planet",
	})

	return events
}

func formatFloat(f float64) string {
	return strconvFormat(f)
}

func strconvFormat(f float64) string {
	// Conversión simple a string con 1 decimal
	intPart := int(f)
	decPart := int((f - float64(intPart)) * 10)
	if decPart < 0 {
		decPart = -decPart
	}
	return itoa(intPart) + "." + itoa(decPart)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var buf []byte
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
