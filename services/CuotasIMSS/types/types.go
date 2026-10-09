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
type Meses struct {
	Dias  int `json:"dias"`
	Orden int `json:"orden"`
}
