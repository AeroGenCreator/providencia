package service

import (
	"saavedra/services/CuotasIMSS/store"
	"saavedra/services/CuotasIMSS/types"
	"saavedra/utils"
	"strconv"
)

type Service interface {
	// Funciones fichero logics.go
	CalcularCuotas(id int)
	// Funciones en este fichero
	SliceCuotasIMSS(page string) (*types.SliceCuotasIMSS, error)
	SearchCuotasIMSS(pattern, page string) (*types.SliceCuotasIMSS, error)
	Many2OneCuotasIMSS() (*types.Many2One, error)
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
		Employee: employees,
		Months:   MesesM2O[:],
	}
	return &data, nil
}
