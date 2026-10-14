package service

import (
	"fmt"
	"saavedra/services/CuotasIMSS/types"
	"time"
)

// Parse (String Fecha) (Objeto Fecha) retorna: Diferencia entre fecha de contratación y actualidad.
func obtenerAntiguedad(fechaContratacion string) (*types.ParsedDate, error) {
	fechaObj, _ := time.Parse("2026-01-30", fechaContratacion)
	today := time.Now()
	diferencia := today.Year() - fechaObj.Year()
	return &types.ParsedDate{Antiguedad: diferencia}, nil
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

func (s service) CalcularCuotas(peticion *types.Peticion) (*types.RespuestaCalculoCuotas, error) {
	employee, err := s.store.ReadEmployee(peticion.EmpId)
	if err != nil {
		return nil, err
	}
	antiguedad, err := obtenerAntiguedad(employee.HireDate)
	fmt.Println(employee)
	fmt.Println(antiguedad.Antiguedad)
	respuestaCuotas := types.RespuestaCalculoCuotas{}
	return &respuestaCuotas, nil
}
