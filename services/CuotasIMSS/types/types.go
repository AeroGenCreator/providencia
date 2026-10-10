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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Slcie M2O Empaquetado
type M2OSlice struct {
	Records []*M2OEmployee `json:"records"`
	Meses   []string
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
	Id   int    `json:""`
	Mes  string `json:""`
	Year string `json:""`
}

// Respuesta a cálculo de cuotas
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
