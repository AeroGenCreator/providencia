package types

import "errors"

// Permite la selección de (Razon Incapacidad)
type CuotasBloques struct {
	Id    string `json:"id"`
	Razon string `json:"razon"`
}

// Tipo Dias De Vacaciones
type DiasVacaciones struct {
	MinLim int `json:"minLim"`
	MaxLim int `json:"maxLim"`
	Dias   int `json:"dias"`
}

// Tipo Slice de Dias De Vacaciones
type SliceDiasVacaciones struct {
	Records []*DiasVacaciones `json:"records"`
}

// Tipo Cesantía Y Vejez
type CesantiaVejez struct {
	MinLim float64 `json:"minLim"`
	MaxLim float64 `json:"maxLim"`
	Factor float64 `json:"factor"`
}

// Los Editables Comparten la misma estructura (Fecha & Valor)
type EditContent struct {
	Fecha int     `json:"fecha"`
	Valor float64 `json:"valor"`
}

// El JSON editable parsea las llaves principales.
type Editables struct {
	DiasAguinaldo             EditContent `json:"diasAguinaldo"`
	ActualUma                 EditContent `json:"actualUMA"`
	PrimaVacacional           EditContent `json:"primaVacacional"`
	ToparSueldo               EditContent `json:"toparSueldo"`
	EnfYMatFactorDeterminante EditContent `json:"enfYMatFactorDeterminante"`
	SalarioMinimo             EditContent `json:"salarioMinimo"`
}

type Riesgo struct {
	Clase string  `json:"clase"`
	Value float64 `json:"value"`
}

// Tipo Meses Metadata
type MesesMetadata struct {
	Dias  int `json:"dias"`
	Orden int `json:"orden"`
}

// Tipo Empleado para Many2One
type M2OEmployee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Slcie Many2One Empaquetado
type Many2One struct {
	Employee        []*M2OEmployee   `json:"employee"`
	Months          []string         `json:"months"`
	UMA             float64          `json:"uma"`
	UMAFecha        int              `json:"umaFecha"`
	MinPayment      float64          `json:"minPayment"`
	MinPaymentFecha int              `json:"minPaymentFecha"`
	Riesgos         []*Riesgo        `json:"riesgos"`
	Razones         []*CuotasBloques `json:"razones"`
}

// Incapacidad Objeto
type IncapacidadObject struct {
	Razon string `json:"razon"`
	Dias  int    `json:"dias"`
}

// Petición De Cálculo Metadata
type Peticion struct {
	EmpId       int     `json:"empId"`
	Month       string  `json:"month"`
	Year        string  `json:"year"`
	Razon       string  `json:"razon"`
	Incapacidad int     `json:"incapacidad"`
	Faltas      int     `json:"faltas"`
	Riesgo      float64 `json:"riesgo"`
}

type ParsedDate struct {
	Antiguedad int `json:"antiguedad"`
}

// Tipo Empleado Record Completo
type Employee struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	HireDate     string  `json:"hireDate"`
	DailyPayment float64 `json:"dailyPayment"`
	Phone        int     `json:"phone"`
	Email        string  `json:"email"`
	Nss          string  `json:"nss"`
	Curp         string  `json:"curp"`
}

// Respuesta a cálculo de cuotas (Siguientes 3 Bloques)
type Patron struct {
	RiesgoTrabajoPrima             float64 `json:"riesgoTrabajoPrima"`
	EnfermedadMaternidadCuotaFija  float64 `json:"enfermedadMaternidadCuotaFija"`  // 20.40% de la UMA (por asegurado)
	EnfermedadMaternidadAdicional  float64 `json:"enfermedadMaternidadAdicional"`  // 1.10% sobre el excedente de 3 UMAs del SBC
	PrestacionesDineroPatron       float64 `json:"prestacionesDineroPatron"`       // 0.70% del SBC
	GastosMedicosPensionadosPatron float64 `json:"gastosMedicosPensionadosPatron"` // 1.05% del SBC
	InvalidezYVida                 float64 `json:"invalidezYVida"`                 // 1.75% del SBC
	Retiro                         float64 `json:"retiro"`                         // 2.00% del SBC
	CesantiaYVejezPatron           float64 `json:"cesantiaYVejezPatron"`           // Tasa progresiva (varía según rangos de UMA/Salarios Mínimos, ej. 3.15% en adelante)
	Guarderias                     float64 `json:"guarderias"`                     // 1.00% del SBC
	AportacionInfonavit            float64 `json:"aportacionInfonavit"`            // 5.00% del SBC
}

type Obrero struct {
	EnfermedadMaternidadAdicionalObrero float64 `json:"enfermedadMaternidadAdicionalObrero"` // 0.40% sobre el excedente de 3 UMAs del SBC
	PrestacionesDineroObrero            float64 `json:"prestacionesDineroObrero"`            // 0.25% del SBC
	GastosMedicosPensionadosObrero      float64 `json:"gastosMedicosPensionadosObrero"`      // 0.375% del SBC
	InvalidezYVidaObrero                float64 `json:"invalidezYVidaObrero"`                // 0.625% del SBC
	CesantiaYVejezObrero                float64 `json:"cesantiaYVejezObrero"`                // 1.125% del SBC
}

type RespuestaCalculoCuotas struct {
	Employee                 *Employee         `json:"employee"`
	Obrero                   *Obrero           `json:"obrero"`
	Patron                   *Patron           `json:"patron"`
	Antiguedad               int               `json:"antiguedad"`
	FactorIntegracion        float64           `json:"factorIntegracion"`
	SBC                      float64           `json:"sbc"`
	Incapacidad              IncapacidadObject `json:"incapacidad"`
	DiasPorCotizar           int               `json:"diasPorCotizar"`
	DiasPorCotizaMesAnterior int               `json:"diasPorCotizarMesAnterior"`
}

// Registro CuotasIMSS
type CuotasIMSS struct {
	Id                                  int     `json:"id"`
	Name                                string  `json:"name"`
	DailyPayment                        float64 `json:"daily_payment"`
	Year                                int     `json:"year"`
	Month                               int     `json:"month"`
	FactorIntegracion                   float64 `json:"factor_integracion"`
	Sbc                                 float64 `json:"sbc"`
	EsBimestre                          bool    `json:"es_bimestre"`
	DiasDelMes                          int     `json:"diasDelMes"`
	DiasPorCotizar                      int     `json:"diasPorCotizar"`
	Incapacidad                         int     `json:"incapacidad"`
	Razon                               string  `json:"razon"`
	Faltas                              int     `json:"faltas"`
	RiesgoTrabajoPrima                  float64 `json:"riesgo_trabajo_prima"`
	EnfermedadMaternidadCuotaFija       float64 `json:"enfermedad_maternidad_cuota_fija"`
	EnfermedadMaternidadAdicional       float64 `json:"enfermedad_maternidad_adicional"`
	PrestacionesDineroPatron            float64 `json:"prestaciones_dinero_patron"`
	GastosMedicosPensionadosPatron      float64 `json:"gastos_medicos_pensionados_patron"`
	InvalidezYVida                      float64 `json:"invalidez_y_vida"`
	Retiro                              float64 `json:"retiro"`
	CesantiaYVejezPatron                float64 `json:"cesantia_y_vejez_patron"`
	Guarderias                          float64 `json:"guarderias"`
	AportacionInfonavit                 float64 `json:"aportacion_infonavit"`
	EnfermedadMaternidadAdicionalObrero float64 `json:"enfermedad_maternidad_adicional_obrero"`
	PrestacionesDineroObrero            float64 `json:"prestaciones_dinero_obrero"`
	GastosMedicosPensionadosObrero      float64 `json:"gastos_medicos_pensionados_obrero"`
	InvalidezYVidaObrero                float64 `json:"invalidez_y_vida_obrero"`
	CesantiaYVejezObrero                float64 `json:"cesantia_y_vejez_obrero"`
}

type SliceCuotasIMSS struct {
	Records     []*CuotasIMSS `json:"records"`
	HasNextPage bool          `json:"hasNextPage"`
}

type Record struct {
	Id int `json:"id"`
}

// Fila de Tabla Report (SQL)
type ReportRow struct {
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	Name           string  `json:"name"`
	SubtotalObrero float64 `json:"subtotalObrero"`
	SubtotalPatron float64 `json:"subtotalPatron"`
	Total          float64 `json:"total"`
}

// Slice de Tabla Report (SQL)
type ReportSlice struct {
	Records []*ReportRow `json:"records"`
}

// Solicitud del Reporte
type ReportRequest struct {
	Year string `json:"year"`
}

var SalaryError = errors.New("Imposible cálculo de cuotas cuando el salario del empleado es menor que el salario mínimo")
