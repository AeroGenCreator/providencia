package service

import (
	"saavedra/services/CuotasIMSS/types"
	"time"
)

// Parse (String Fecha) (Objeto Fecha) retorna: Diferencia entre fecha de contratación y actualidad.
func obtenerAntiguedad(fechaContratacion string) (*types.ParsedDate, error) {
	fechaObj, _ := time.Parse("2006-01-02", fechaContratacion)
	today := time.Now()
	diff := today.Year() - fechaObj.Year()
	// Si mes de contratación menor o dia de contratacion menor a la fecha de hoy restar en 1 a la differencia.
	if today.Month() < fechaObj.Month() || (today.Month() == fechaObj.Month() && today.Day() < fechaObj.Day()) {
		diff--
	}
	// Si la diferencia es menor a 0, retornar una antiguedad de 0.
	if diff < 0 {
		diff = 0
	}
	return &types.ParsedDate{Antiguedad: diff}, nil
}

func esBimestre(mes string) bool {
	metadata := Meses[mes]
	if modulo := metadata.Orden % 2; modulo == 0 {
		return true
	}
	return false
}

func calcularDiasVacaciones(antiguedad int) int {
	var dias int
	ultimo := diasVacaciones[len(diasVacaciones)-1]
	if ultimo.MaxLim <= antiguedad {
		entero := antiguedad / 5
		dias = 20 + (2 * entero)
		return dias
	}
	for _, dict := range diasVacaciones {
		if dict.MinLim <= antiguedad && antiguedad < dict.MaxLim {
			dias = dict.Dias
		}
	}
	return dias
}

func calcularFactorIntegracion(diasVacaciones float64) float64 {
	var diasAnho float64 = 365
	var salarioNeutro float64 = 1
	diasAguinaldo := editables["diasAguinaldo"].Valor
	primaVacacional := editables["primaVacacional"].Valor
	factorIntegracion := salarioNeutro + ((diasAguinaldo + (diasVacaciones * primaVacacional)) / diasAnho)
	return factorIntegracion
}

func calcularSBC(payment, factorIntegracion float64) float64 {
	return payment * factorIntegracion
}

func motorCuotasIMSS(sbc, riesgo float64, bimestre bool) (*types.Obrero, *types.Patron, error) {
	// Calculos Obrero
	obrero := types.Obrero{}

	// Calculos Patron
	patron := types.Patron{}

	// Calculos Bimestre

	return &obrero, &patron, nil
}

func (s service) CalcularCuotas(peticion *types.Peticion) (*types.RespuestaCalculoCuotas, error) {
	employee, err := s.store.ReadEmployee(peticion.EmpId)
	if err != nil {
		return nil, err
	}
	antiguedad, err := obtenerAntiguedad(employee.HireDate)
	diasVacaciones := calcularDiasVacaciones(antiguedad.Antiguedad)
	factorIntegracion := calcularFactorIntegracion(float64(diasVacaciones))
	sbc := calcularSBC(employee.DailyPayment, factorIntegracion)
	bimestre := esBimestre(peticion.Month)
	obrero, patron, err := motorCuotasIMSS(sbc, peticion.Riesgo, bimestre)
	if err != nil {
		return nil, err
	}
	cuotas := types.RespuestaCalculoCuotas{
		Employee: employee,
		Obrero:   obrero,
		Patron:   patron,
	}
	return &cuotas, nil
}
