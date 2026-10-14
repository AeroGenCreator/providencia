package types

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

// Tipo Editables (UMA y Salario Mínimo)
type Editables struct {
	Fecha int     `json:"fecha"`
	Valor float64 `json:"valor"`
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
	Employee        []*M2OEmployee `json:"employee"`
	Months          []string       `json:"months"`
	UMA             float64        `json:"uma"`
	UMAFecha        int            `json:"umaFecha"`
	MinPayment      float64        `json:"minPayment"`
	MinPaymentFecha int            `json:"minPaymentFecha"`
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

// Petición De Cálculo Metadata
type Peticion struct {
	EmpId           int    `json:"empId"`
	Month           string `json:"month"`
	Year            string `json:"year"`
	PrevIncapacidad int    `json:"prevIncapacidad"`
	Incapacidad     int    `json:"incapacidad"`
	Faltas          int    `json:"faltas"`
}

type ParsedDate struct {
	Antiguedad int `json:"antiguedad"`
}

// Respuesta a cálculo de cuotas (Siguientes 3 Bloques)
type Patron struct {
	RiesgoTrabajoPrima             float64
	EnfermedadMaternidadCuotaFija  float64 // 20.40% de la UMA (por asegurado)
	EnfermedadMaternidadAdicional  float64 // 1.10% sobre el excedente de 3 UMAs del SBC
	PrestacionesDineroPatron       float64 // 0.70% del SBC
	GastosMedicosPensionadosPatron float64 // 1.05% del SBC
	InvalidezYVida                 float64 // 1.75% del SBC
	Retiro                         float64 // 2.00% del SBC
	CesantiaYVejezPatron           float64 // Tasa progresiva (varía según rangos de UMA/Salarios Mínimos, ej. 3.15% en adelante)
	Guarderias                     float64 // 1.00% del SBC
	AportacionInfonavit            float64 // 5.00% del SBC
}

type Obrero struct {
	EnfermedadMaternidadAdicionalObrero float64 // 0.40% sobre el excedente de 3 UMAs del SBC
	PrestacionesDineroObrero            float64 // 0.25% del SBC
	GastosMedicosPensionadosObrero      float64 // 0.375% del SBC
	InvalidezYVidaObrero                float64 // 0.625% del SBC
	CesantiaYVejezObrero                float64 // 1.125% del SBC
}

type RespuestaCalculoCuotas struct {
	Employee *Employee `json:"employee"`
	Obrero   *Obrero   `json:"obrero"`
	Patron   *Patron   `json:"patron"`
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
	Incapacidad                         int     `json:"incapacidad"`
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
