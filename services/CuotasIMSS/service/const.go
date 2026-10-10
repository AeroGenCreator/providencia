package service

import "saavedra/services/CuotasIMSS/types"

var Meses = map[string]types.MesesMetadata{
	"Enero":      {Dias: 31, Orden: 1},
	"Febrero":    {Dias: 28, Orden: 2},
	"Marzo":      {Dias: 31, Orden: 3},
	"Abril":      {Dias: 30, Orden: 4},
	"Mayo":       {Dias: 31, Orden: 5},
	"Junio":      {Dias: 30, Orden: 6},
	"Julio":      {Dias: 31, Orden: 7},
	"Agosto":     {Dias: 31, Orden: 8},
	"Septiembre": {Dias: 30, Orden: 9},
	"Octubre":    {Dias: 31, Orden: 10},
	"Noviembre":  {Dias: 30, Orden: 11},
	"Diciembre":  {Dias: 31, Orden: 12},
}

var MesesM2O = [12]string{
	"Enero",
	"Febrero",
	"Marzo",
	"Abril",
	"Mayo",
	"Junio",
	"Julio",
	"Agosto",
	"Septiembre",
	"Octubre",
	"Noviembre",
	"Diciembre",
}
