package astronomy

import (
	"fmt"
	"math"
	"time"
)

const (
	degreesToRadians = math.Pi / 180
	radiansToDegrees = 180.0 / math.Pi
	solarZenith      = 90.833
)

type RiseSet struct {
	Rise time.Time
	Set  time.Time
	// Has va a indicar si ese astro sale y se pone ese día
	// Útil para zonas polares (día sin noche o noche sin día)
	Has bool
}

// SunRiseSet calcula la salida y puesta del Sol para una fecha y ubicación dadas
// Basado en una aproximación del algoritmo NOAA
func SunRiseSet(lat, lon float64, date time.Time) (RiseSet, error) {
	// Validar coordenadas
	if math.IsNaN(lat) || math.IsInf(lat, 0) || math.IsNaN(lon) || math.IsInf(lon, 0) {
		return RiseSet{}, fmt.Errorf("las coordenadas deben ser números válidos")
	}
	if lat < -90 || lat > 90 {
		return RiseSet{}, fmt.Errorf("latitud fuera de rango")
	}
	if lon < -180 || lon > 180 {
		return RiseSet{}, fmt.Errorf("longitud fuera de rango")
	}

	// Obtener día del año
	day := date.YearDay()

	// Determinar cantidad de días del año
	daysInYear := 365.0
	if isLeapYear(date.Year()) {
		daysInYear = 366.0
	}

	// Fracción anual en radianes
	gamma := 2 * math.Pi / daysInYear * float64(day-1)

	sin := math.Sin
	cos := math.Cos

	// Ecuación del tiempo en minutos
	eqTime := 229.18 * (0.000075 +
		0.001868*cos(gamma) -
		0.032077*sin(gamma) -
		0.014615*cos(2*gamma) -
		0.040849*sin(2*gamma))

	// Declinación solar en radianes
	decl := 0.006918 -
		0.399912*cos(gamma) +
		0.070257*sin(gamma) -
		0.006758*cos(2*gamma) +
		0.000907*sin(2*gamma) -
		0.002697*cos(3*gamma) +
		0.001480*sin(3*gamma)

	// Latitud en radianes
	latitude := lat * degreesToRadians
	zenith := solarZenith * degreesToRadians

	denominator := cos(latitude) * cos(decl)
	// Evitar división problemática cerca de los polos
	if math.Abs(denominator) < 1e-12 {
		return RiseSet{Has: false}, nil
	}

	// Coseno del ángulo horario
	cosHourAngle := (cos(zenith) - sin(latitude)*sin(decl)) / denominator

	// Comprobar si hay salida y puesta ese día
	if cosHourAngle < -1 || cosHourAngle > 1 {
		return RiseSet{Has: false}, nil
	}

	// Ángulo horario en grados
	hourAngle := math.Acos(cosHourAngle) * radiansToDegrees

	// Mediodía solar en minutos UTC
	solarNoon := 720 - 4*lon - eqTime

	// Salida y puesta en minutos desde medianoche UTC
	riseMinutes := solarNoon - 4*hourAngle
	setMinutes := solarNoon + 4*hourAngle

	// Construir instantes a partir de medianoche UTC del día solicitado
	utcMidnight := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	rise := utcMidnight.Add(time.Duration(riseMinutes * float64(time.Minute)))
	set := utcMidnight.Add(time.Duration(setMinutes * float64(time.Minute)))

	// Devolver en la zona horaria de la fecha recibida
	return RiseSet{
		Rise: rise.In(date.Location()),
		Set:  set.In(date.Location()),
		Has:  true,
	}, nil
}

// MoonRiseSet calcula la salida y puesta de la Luna para una fecha y ubicación dadas
func MoonRiseSet(lat, lon float64, date time.Time) (RiseSet, error) {
	// Validar coordenadas
	if math.IsNaN(lat) || math.IsInf(lat, 0) || math.IsNaN(lon) || math.IsInf(lon, 0) {
		return RiseSet{}, fmt.Errorf("las coordenadas deben ser números válidos")
	}
	if lat < -90 || lat > 90 {
		return RiseSet{}, fmt.Errorf("latitud fuera de rango")
	}
	if lon < -180 || lon > 180 {
		return RiseSet{}, fmt.Errorf("longitud fuera de rango")
	}

	// TODO: Implementar cálculo de salida/puesta lunar
	h0 := -0.125 * math.Pi / 180.0

	return RiseSet{Has: false}, nil
}

// isLeapYear indica si un año es bisiesto
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// julianDay calcula el dia Juliano para una fecha UTC a las 00:00
func julianDay(t time.Time) float64 {
	year := t.Year()
	month := int(t.Month())
	day := float64(t.Day()) + (float64(t.Hour())+float64(t.Minute())/60+float64(t.Second())/3600)/24

	if month <= 2 {
		year--
		month += 12
	}

	A := year / 100
	B := 2 - A + A/4

	jd := float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		day + float64(B) - 1524.5

	return jd
}

// julianCenturies pasa de JD a siglos julianos desde J2000.0
func julianCenturies(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}
