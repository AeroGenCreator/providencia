package service

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

func (s service) CalcularCuotas(id int) {

}
