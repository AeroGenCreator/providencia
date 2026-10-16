package service

import (
	"saavedra/services/CuotasIMSS/types"
	"strconv"
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
	// Tope de sueldo
	sueldoTopado := editables["actualUMA"].Valor * editables["toparSueldo"].Valor
	// Cálculo del SBC
	sbc := payment * factorIntegracion
	if payment > sueldoTopado {
		sbc = sueldoTopado
	}
	return sbc
}

func esBisiesto(year string) (bool, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return false, err
	}
	if intYear%4 == 0 {
		return true, nil
	}
	return false, nil
}

func calcularDiasPorCotizar(mes string, faltas int, bisiesto bool) int {
	if bisiesto && mes == "Febrero" {
		return 29
	}
	diasMes := Meses[mes].Dias
	diasPorCotizar := diasMes - faltas
	return diasPorCotizar
}

func calculoCesantiaYVejezPatron(year string, sbc float64, diasPorCotizar int) float64 {
	// Calcula CesantiaVejez depende de buscar en una tabla de rangos de veces UMA del SBC.
	vecesUMA := (sbc / editables["actualUMA"].Valor)
	rangos := cesantiaVejez[year]
	// Para este proposito use BinarySearch (Ya que el orden es natural)
	// Indice 0 y el ultimo del arreglo.
	LimAlto := len(rangos) - 1
	limBajo := 0
	// Siempre que el indice menor siga por debajo del mayor (se itera)
	for limBajo <= LimAlto {
		// Se calcula la mitad del arreglo en cada iteración.
		mid := ((limBajo + LimAlto) / 2)
		record := rangos[mid]
		// Si las veces uma esta dentro del rango se regresa el cálculo
		if record.MinLim <= vecesUMA && vecesUMA <= record.MaxLim {
			return (record.Factor * sbc * float64(diasPorCotizar))
			// Si las veces UMA es mayor al limite superior, se busca en la mitad alta del arreglo.
		} else if vecesUMA > record.MaxLim {
			limBajo = mid + 1
			// Si las veces UMA es menor al limite inferior, se busca en la mitad baja del arreglo.
		} else {
			LimAlto = mid - 1
		}
	}
	return 0
}

func motorCuotasIMSS(
	year string, sbc, riesgo float64, bimestre bool, diasPorCotizar, diasPorCotizarMesAnterior int,
) (*types.Obrero, *types.Patron, error) {

	// Enf. Mat. en Dinero (Calculo especial sobre el excedente en 2026 (3 veces UMA))
	determinante := editables["enfYmatFactorDeterminante"].Valor * editables["actualUMA"].Valor
	diferencia := sbc - determinante
	enfMatObrero := 0.0
	enfMatPatron := 0.0
	if diferencia > 0 {
		enfMatObrero = (diferencia * factores["enfermedadMaternidadAdicionalObrero"] * float64(diasPorCotizar))
		enfMatPatron = (diferencia * factores["enfermedadMaternidadAdicionalPatron"] * float64(diasPorCotizar))
	}

	// Cálculos Bimestrales
	retiro := 0.0
	aportacionInfonavit := 0.0
	if bimestre {
		retiro = (factores["retiro"] * sbc * float64(diasPorCotizar+diasPorCotizarMesAnterior))
		aportacionInfonavit = (factores["aportacionInfonavit"] * sbc * float64(diasPorCotizar+diasPorCotizarMesAnterior))
	}

	// Calcular Cesantia Y Vejez (Patron)
	cesantiaYVejezPatron := calculoCesantiaYVejezPatron(year, sbc, diasPorCotizar)

	// Calculos Obrero
	obrero := types.Obrero{
		EnfermedadMaternidadAdicionalObrero: enfMatObrero,
		PrestacionesDineroObrero:            (factores["prestacionesDineroObrero"] * sbc * float64(diasPorCotizar)),
		GastosMedicosPensionadosObrero:      (factores["gastosMedicosPensionadosObrero"] * sbc * float64(diasPorCotizar)),
		InvalidezYVidaObrero:                (factores["invalidezYVidaObrero"] * sbc * float64(diasPorCotizar)),
		CesantiaYVejezObrero:                (factores["cesantiaYVejezObrero"] * sbc * float64(diasPorCotizar)),
	}

	// Calculos Patron
	patron := types.Patron{
		RiesgoTrabajoPrima:             (riesgo * sbc * float64(diasPorCotizar)),
		EnfermedadMaternidadCuotaFija:  (editables["actualUMA"].Valor * float64(diasPorCotizar) * factores["enfermedadMaternidadCuotaFija"]),
		EnfermedadMaternidadAdicional:  enfMatPatron,
		PrestacionesDineroPatron:       (factores["prestacionesDineroPatron"] * sbc * float64(diasPorCotizar)),
		GastosMedicosPensionadosPatron: (factores["gastosMedicosPensionadosPatron"] * sbc * float64(diasPorCotizar)),
		InvalidezYVida:                 (factores["invalidezYVidaPatron"] * sbc * float64(diasPorCotizar)),
		Guarderias:                     (factores["guarderias"] * sbc * float64(diasPorCotizar)),
		Retiro:                         retiro,
		AportacionInfonavit:            aportacionInfonavit,
		CesantiaYVejezPatron:           cesantiaYVejezPatron,
	}

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
	bisiesto, err := esBisiesto(peticion.Year)
	if err != nil {
		return nil, err
	}
	diasPorCotizaMesAnterior := 0
	if bimestre {
		// Incluso en año bisiesto, paso false a la función porque el unico mes cambiante es "par" o sea "Febrero"
		diasPorCotizaMesAnterior = calcularDiasPorCotizar(peticion.Month, 0, false)
	}
	diasPorCotizar := calcularDiasPorCotizar(peticion.Month, peticion.Faltas, bisiesto)
	obrero, patron, err := motorCuotasIMSS(peticion.Year, sbc, peticion.Riesgo, bimestre, diasPorCotizar, diasPorCotizaMesAnterior)
	if err != nil {
		return nil, err
	}
	incapacidadObject := types.IncapacidadObject{Razon: peticion.Razon, Dias: peticion.Incapacidad}
	cuotas := types.RespuestaCalculoCuotas{
		Employee:          employee,
		Obrero:            obrero,
		Patron:            patron,
		Antiguedad:        antiguedad.Antiguedad,
		FactorIntegracion: factorIntegracion,
		SBC:               sbc,
		Incapacidad:       incapacidadObject,
	}
	return &cuotas, nil
}
