package astronomy

import (
	"fmt"
	"math"
	"time"
)

// Constantes astronómicas base
const (
	degreesToRadians     = math.Pi / 180.0
	radiansToDegrees     = 180.0 / math.Pi
	solarZenith          = 90.833
	moonStandardAltitude = 0.125 // h0 promedio de la Luna (+0°07'30") [Meeus Cap. 15]
)

type Coordinates struct {
	RA  float64 // Ascensión recta (α) en grados
	Dec float64 // Declinación (δ) en grados
}

type RiseSet struct {
	Rise time.Time
	Set  time.Time
	Has  bool
}

// ============================================================================
// ☀️ CÁLCULO DEL SOL (100% INTACTO)
// ============================================================================

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

// isLeapYear indica si un año es bisiesto
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// ============================================================================
// 🌙 CÁLCULO DE LA LUNA (Meeus + Bisección)
// ============================================================================

func normalize360(deg float64) float64 {
	res := math.Mod(deg, 360.0)
	if res < 0 {
		res += 360.0
	}
	return res
}

func normalize180(deg float64) float64 {
	res := math.Mod(deg, 360.0)
	if res < -180.0 {
		res += 360.0
	} else if res >= 180.0 {
		res -= 360.0
	}
	return res
}

func timeToJD(t time.Time) float64 {
	t = t.UTC()
	year := t.Year()
	month := int(t.Month())
	day := float64(t.Day()) + float64(t.Hour())/24.0 + float64(t.Minute())/1440.0 + float64(t.Second())/86400.0 + float64(t.Nanosecond())/86400e9

	if month <= 2 {
		year--
		month += 12
	}

	a := math.Floor(float64(year) / 100.0)
	b := 2.0 - a + math.Floor(a/4.0)

	return math.Floor(365.25*float64(year+4716)) + math.Floor(30.6001*float64(month+1)) + day + b - 1524.5
}

// MoonPositionMeeus calcula la posición geocéntrica exacta de la Luna (ELP-2000, Cap. 47)
func MoonPositionMeeus(jd float64) Coordinates {
	T := (jd - 2451545.0) / 36525.0

	LPrime := normalize360(218.3164477 + 481267.88123421*T - 0.0015786*T*T + 0.0000019*T*T*T)
	D := normalize360(297.8501921 + 445267.1114034*T - 0.0018819*T*T + T*T*T/545868.0)
	M := normalize360(357.5291092 + 35999.0502909*T - 0.0001536*T*T + T*T*T/24490000.0)
	MPrime := normalize360(134.9633964 + 477198.8675055*T + 0.0087414*T*T + T*T*T/69699.0)
	F := normalize360(93.2720950 + 483202.0175233*T - 0.0036539*T*T - T*T*T/3526000.0)

	DRad := D * degreesToRadians
	MRad := M * degreesToRadians
	MPrimeRad := MPrime * degreesToRadians
	FRad := F * degreesToRadians

	sumL := 6288774.0*math.Sin(MPrimeRad) +
		1274027.0*math.Sin(2*DRad-MPrimeRad) +
		658314.0*math.Sin(2*DRad) +
		213618.0*math.Sin(2*MPrimeRad) -
		185116.0*math.Sin(MRad) -
		114332.0*math.Sin(2*FRad) +
		58793.0*math.Sin(2*DRad-2*MPrimeRad) +
		57066.0*math.Sin(2*DRad-MRad-MPrimeRad) +
		53322.0*math.Sin(2*DRad+MPrimeRad) +
		45758.0*math.Sin(2*DRad-MRad) -
		40923.0*math.Sin(MRad-MPrimeRad) -
		34720.0*math.Sin(DRad) -
		30383.0*math.Sin(MPrimeRad+2*FRad) +
		15327.0*math.Sin(2*DRad-2*FRad) -
		12528.0*math.Sin(MPrimeRad-2*FRad) +
		10980.0*math.Sin(MPrimeRad+2*FRad)

	lambda := normalize360(LPrime + sumL/1000000.0)

	sumB := 5128122.0*math.Sin(FRad) +
		280602.0*math.Sin(MPrimeRad+FRad) +
		277693.0*math.Sin(MPrimeRad-FRad) +
		173237.0*math.Sin(2*DRad-FRad) +
		55413.0*math.Sin(2*DRad-MPrimeRad+FRad) +
		46271.0*math.Sin(2*DRad-MPrimeRad-FRad) +
		32573.0*math.Sin(2*DRad+FRad) +
		17198.0*math.Sin(2*MPrimeRad+FRad)

	beta := sumB / 1000000.0

	eps0 := 23.4392911 - 0.0130041667*T - 0.0000001639*T*T + 0.0000005036*T*T*T

	lambdaRad := lambda * degreesToRadians
	betaRad := beta * degreesToRadians
	epsRad := eps0 * degreesToRadians

	sinAlpha := math.Sin(lambdaRad)*math.Cos(epsRad) - math.Tan(betaRad)*math.Sin(epsRad)
	cosAlpha := math.Cos(lambdaRad)
	ra := normalize360(math.Atan2(sinAlpha, cosAlpha) * radiansToDegrees)

	sinDec := math.Sin(betaRad)*math.Cos(epsRad) + math.Cos(betaRad)*math.Sin(epsRad)*math.Sin(lambdaRad)
	dec := math.Asin(sinDec) * radiansToDegrees

	return Coordinates{RA: ra, Dec: dec}
}

// getMoonAltitude calcula la altitud topocéntrica/geométrica de la Luna en grados
func getMoonAltitude(tUTC time.Time, lat, lon float64) float64 {
	jd := timeToJD(tUTC)
	pos := MoonPositionMeeus(jd)
	T := (jd - 2451545.0) / 36525.0
	theta := normalize360(280.46061837 + 360.98564736629*(jd-2451545.0) + 0.000387933*T*T - T*T*T/3871000.0)

	// Ángulo horario local H (lon: Este +, Oeste -)
	H := normalize180(theta + lon - pos.RA)

	phiRad := lat * degreesToRadians
	decRad := pos.Dec * degreesToRadians
	HRad := H * degreesToRadians

	sinH := math.Sin(phiRad)*math.Sin(decRad) + math.Cos(phiRad)*math.Cos(decRad)*math.Cos(HRad)
	return math.Asin(sinH) * radiansToDegrees
}

// MoonRiseSet calcula la salida y puesta exacta de la Luna por búsqueda directa y bisección.
func MoonRiseSet(lat, lon float64, date time.Time) (RiseSet, error) {
	loc := date.Location()
	localStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	utcStart := localStart.UTC()

	// Muestrear 24 intervalos de 1 hora
	const hours = 24
	alts := make([]float64, hours+1)
	times := make([]time.Time, hours+1)

	for i := 0; i <= hours; i++ {
		t := utcStart.Add(time.Duration(i) * time.Hour)
		times[i] = t
		alts[i] = getMoonAltitude(t, lat, lon)
	}

	var riseTime, setTime time.Time
	var hasRise, hasSet bool
	h0 := moonStandardAltitude

	for i := 0; i < hours; i++ {
		h1 := alts[i]
		h2 := alts[i+1]
		t1 := times[i]
		t2 := times[i+1]

		// Verificar si cruza el horizonte h0
		if (h1-h0)*(h2-h0) <= 0 && h1 != h2 {
			// Refinamiento por bisección
			ta := t1
			tb := t2
			ha := h1 - h0
			for step := 0; step < 15; step++ {
				tc := ta.Add(tb.Sub(ta) / 2)
				hc := getMoonAltitude(tc, lat, lon) - h0
				if ha*hc <= 0 {
					tb = tc
				} else {
					ta = tc
					ha = hc
				}
			}
			tCross := ta.Add(tb.Sub(ta) / 2).In(loc)
			if h1 < h2 {
				riseTime = tCross
				hasRise = true
			} else {
				setTime = tCross
				hasSet = true
			}
		}
	}

	return RiseSet{
		Rise: riseTime,
		Set:  setTime,
		Has:  hasRise || hasSet,
	}, nil
}
