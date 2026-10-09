package service

import (
	"saavedra/services/Employee/store"
	"saavedra/services/Employee/types"
	"saavedra/utils"
	"strconv"
)

type Service interface {
	SliceEmployee(page string) (*types.EmployeeSlice, error)
	CreateEmployee(employee *types.Employee) error
	ReadEmployee(id int) (*types.Employee, error)
	UpdateEmployee(employee *types.Employee) error
	DeleteEmployee(id int) error
	SearchEmployee(pattern, page string) (*types.EmployeeSlice, error)
}

type service struct {
	store store.Store
}

func New(store store.Store) Service {
	return service{store: store}
}

func (s service) SliceEmployee(page string) (*types.EmployeeSlice, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice

	records, count, err := s.store.SliceEmployee(utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}

	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage

	employeeSlice := types.EmployeeSlice{
		Records:     records,
		HasNextPage: hasNextPage,
	}

	return &employeeSlice, nil
}

func (s service) CreateEmployee(employee *types.Employee) error {
	if err := s.store.CreateEmployee(employee); err != nil {
		return err
	}
	return nil
}

func (s service) ReadEmployee(id int) (*types.Employee, error) {
	employee, err := s.store.ReadEmployee(id)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (s service) UpdateEmployee(employee *types.Employee) error {
	if err := s.store.UpdateEmployee(employee); err != nil {
		return err
	}
	return nil
}

func (s service) DeleteEmployee(id int) error {
	if err := s.store.DeleteEmployee(id); err != nil {
		return err
	}
	return nil
}

func (s service) SearchEmployee(pattern, page string) (*types.EmployeeSlice, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice

	slice, count, err := s.store.SearchEmployee(pattern, utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}

	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage

	records := types.EmployeeSlice{
		Records:     slice,
		HasNextPage: hasNextPage,
	}

	return &records, nil
}
