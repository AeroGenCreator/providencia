package service

import (
	"saavedra/services/Saavedra/Users/store"
	"saavedra/services/Saavedra/Users/types"
	"saavedra/utils"
	"strconv"
)

type Service interface {
	SliceUser(page string) (*types.UserSlice, error)
	CreateUser(user *types.User) error
	ReadUser(id int) (*types.User, error)
	UpdateUser(user *types.User) error
	DeleteUser(id int, rId string) error
	Many2One() *types.Many2One
	SearchUser(pattern, page string) (*types.UserSlice, error)
}

type service struct {
	store store.Store
}

func New(store store.Store) Service {
	return service{store: store}
}

func (s service) SliceUser(page string) (*types.UserSlice, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice
	records, count, err := s.store.SliceUser(utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}
	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage
	userSlice := types.UserSlice{
		Records:     records,
		HasNextPage: hasNextPage,
	}

	return &userSlice, nil
}

func (s service) CreateUser(user *types.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := s.store.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (s service) ReadUser(id int) (*types.User, error) {
	user, err := s.store.ReadUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) UpdateUser(user *types.User) error {
	// Sin contraseña, se actualizan los datos simples.
	if user.Password == "" {
		if err := s.store.UpdateUserNoPassword(user); err != nil {
			return err
		}
		return nil
	}
	// Si hay contraseña se encripta el valor antes de actualizar todo el registro.
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err = s.store.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

func (s service) DeleteUser(id int, rId string) error {
	// 'id' del del usuario que hizo el request.
	intRId, err := strconv.Atoi(rId)
	if err != nil {
		return err
	}
	// protección; cuenta admin
	if id == utils.AdminId {
		return types.ErrAdminProtection
	}
	// protección; no eliminar la misma cuenta con la que se inicio sesión.
	if id == intRId {
		return types.ErrSelfProtection
	}
	if err := s.store.DeleteUser(id); err != nil {
		return err
	}
	return nil
}

func (s service) Many2One() *types.Many2One {
	return &types.Many2One{Parties: utils.Parties}
}

func (s service) SearchUser(pattern, page string) (*types.UserSlice, error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	var offset int
	offset = (intPage - 1) * utils.RecordsPerSlice
	slice, count, err := s.store.SearchUser(pattern, utils.RecordsPerSlice, offset)
	if err != nil {
		return nil, err
	}
	totalPages := utils.CalculateTotalPages(count, utils.RecordsPerSlice)
	hasNextPage := totalPages > intPage
	records := types.UserSlice{
		Records:     slice,
		HasNextPage: hasNextPage,
	}
	return &records, err
}
