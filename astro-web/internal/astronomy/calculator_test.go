package astronomy

import (
	"testing"
	"time"
)

func TestCalcularFaseLunar(t *testing.T) {
	// elegimos fecha
	fechaPrueba := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	// Llamo a la funcion
	fase, iluminacion := CalcularFaseLunar(fechaPrueba)
	// verificamos
	if fase == "" {
		t.Errorf("La fase lunar no debería estar vacía")
	}
	// Verificamos que la iluminación esté entre 0% y 100%
	if iluminacion < 0 || iluminacion > 100 {
		t.Errorf("La iluminación %.2f está fuera del rango válido (0-100)", iluminacion)
	}
}
