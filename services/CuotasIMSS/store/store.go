package store

import (
	"database/sql"
	"saavedra/services/CuotasIMSS/types"
)

type Store interface {
	// Funciones en fichero extra.go
	M2OEmployees() ([]*types.M2OEmployee, error)
	ReadEmployee(id int) (*types.Employee, error)
	// Funciones en este fichero
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}
