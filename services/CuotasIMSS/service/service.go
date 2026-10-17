package service

import (
	"saavedra/services/CuotasIMSS/store"
	"saavedra/services/CuotasIMSS/types"
	"saavedra/utils"
	"strconv"
)

type Service interface {
	// Funciones fichero logics.go
	CalcularCuotas(peticion *types.Peticion) (*types.RespuestaCalculoCuotas, error)
	// Funciones en este fichero
	SliceCuotasIMSS(page string) (*types.SliceCuotasIMSS, error)
	SearchCuotasIMSS(pattern, page string) (*types.SliceCuotasIMSS, error)
	Many2OneCuotasIMSS() (*types.Many2One, error)
	CreateCuotasIMSS(peticion *types.Peticion) error
}

type service struct {
	store store.Store
}

func New(store store.Store) Service {
	return service{store: store}
}

func (s service) SliceCuotasIMSS(page string) (*types.SliceCuotasIMSS, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice
	records, count, err := s.store.SliceCuotasIMSS(utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}
	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage
	slice := types.SliceCuotasIMSS{
		Records:     records,
		HasNextPage: hasNextPage,
	}
	return &slice, nil
}

func (s service) SearchCuotasIMSS(pattern, page string) (*types.SliceCuotasIMSS, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice
	records, count, err := s.store.SearchCuotasIMSS(pattern, utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}
	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage
	data := types.SliceCuotasIMSS{
		Records:     records,
		HasNextPage: hasNextPage,
	}
	return &data, nil
}

func (s service) Many2OneCuotasIMSS() (*types.Many2One, error) {
	employees, err := s.store.M2OEmployees()
	if err != nil {
		return nil, err
	}
	data := types.Many2One{
		Employee:        employees,
		Months:          MesesM2O[:],
		UMA:             editables.ActualUma.Valor,
		UMAFecha:        editables.ActualUma.Fecha,
		MinPayment:      editables.SalarioMinimo.Valor,
		MinPaymentFecha: editables.SalarioMinimo.Fecha,
		Riesgos:         primaRiesgo,
		Razones:         Razones,
	}
	return &data, nil
}

func (s service) CreateCuotasIMSS(peticion *types.Peticion) error {
	cuotas, err := s.CalcularCuotas(peticion)
	if err != nil {
		return err
	}
	yearInt, err := strconv.Atoi(peticion.Year)
	if err != nil {
		return err
	}
	mesNum := Meses[peticion.Month].Orden
	mesDias := Meses[peticion.Month].Dias
	esBimestre := esBimestre(peticion.Month)
	diasPorCotizar := cuotas.DiasPorCotizaMesAnterior
	if esBimestre {
		diasPorCotizar += cuotas.DiasPorCotizaMesAnterior
	}
	cuotaRegistroNuevo := types.CuotasIMSS{
		Name:                                cuotas.Employee.Name,
		DailyPayment:                        cuotas.Employee.DailyPayment,
		Year:                                yearInt,
		Month:                               mesNum,
		FactorIntegracion:                   cuotas.FactorIntegracion,
		Sbc:                                 cuotas.SBC,
		EsBimestre:                          esBimestre,
		DiasDelMes:                          mesDias,
		DiasPorCotizar:                      diasPorCotizar,
		Incapacidad:                         cuotas.Incapacidad.Dias,
		Razon:                               cuotas.Incapacidad.Razon,
		Faltas:                              peticion.Faltas,
		RiesgoTrabajoPrima:                  cuotas.Patron.RiesgoTrabajoPrima,
		EnfermedadMaternidadCuotaFija:       cuotas.Patron.EnfermedadMaternidadCuotaFija,
		EnfermedadMaternidadAdicional:       cuotas.Patron.EnfermedadMaternidadAdicional,
		PrestacionesDineroPatron:            cuotas.Patron.PrestacionesDineroPatron,
		GastosMedicosPensionadosPatron:      cuotas.Patron.GastosMedicosPensionadosPatron,
		InvalidezYVida:                      cuotas.Patron.InvalidezYVida,
		Retiro:                              cuotas.Patron.Retiro,
		CesantiaYVejezPatron:                cuotas.Patron.CesantiaYVejezPatron,
		Guarderias:                          cuotas.Patron.Guarderias,
		AportacionInfonavit:                 cuotas.Patron.AportacionInfonavit,
		EnfermedadMaternidadAdicionalObrero: cuotas.Obrero.EnfermedadMaternidadAdicionalObrero,
		PrestacionesDineroObrero:            cuotas.Obrero.PrestacionesDineroObrero,
		GastosMedicosPensionadosObrero:      cuotas.Obrero.GastosMedicosPensionadosObrero,
		InvalidezYVidaObrero:                cuotas.Obrero.InvalidezYVidaObrero,
		CesantiaYVejezObrero:                cuotas.Obrero.CesantiaYVejezObrero,
	}
	if err = s.store.CreateCuotasIMSS(&cuotaRegistroNuevo); err != nil {
		return err
	}
	return nil
}
