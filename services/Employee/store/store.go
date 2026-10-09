package store

import (
	"database/sql"
	"saavedra/services/Employee/types"
	utils "saavedra/utils"
)

type Store interface {
	SliceEmployee(limit, offset int) ([]*types.Employee, int, error)
	CreateEmployee(employee *types.Employee) error
	ReadEmployee(id int) (*types.Employee, error)
	UpdateEmployee(employee *types.Employee) error
	DeleteEmployee(id int) error
	SearchEmployee(pattern string, limit, offset int) ([]*types.Employee, int, error)
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}

func (s store) SliceEmployee(limit, offset int) ([]*types.Employee, int, error) {
	q1 := "SELECT COUNT(id) AS count_id FROM employee;"
	q2 := "SELECT id, name, hire_date, daily_payment, phone, email, nss, curp FROM employee LIMIT ? OFFSET ?;"

	var count int
	err := s.db.QueryRow(q1).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.db.Query(q2, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*types.Employee
	for rows.Next() {
		var emp types.Employee
		err = rows.Scan(&emp.Id, &emp.Name, &emp.HireDate, &emp.DailyPayment, &emp.Phone, &emp.Email, &emp.Nss, &emp.Curp)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, &emp)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return records, count, nil
}

func (s store) CreateEmployee(employee *types.Employee) error {
	c := `SELECT id FROM employee WHERE email = ? OR nss = ? OR curp = ?;`
	q := `INSERT INTO employee(name, hire_date, daily_payment, phone, email, nss, curp) VALUES (?, ?, ?, ?, ?, ?, ?);`

	var id int
	err := s.db.QueryRow(c, employee.Email, employee.Nss, employee.Curp).Scan(&id)
	if err == nil {
		return utils.DuplicatedDataError
	}
	if err != sql.ErrNoRows {
		return err
	}

	_, err = s.db.Exec(
		q,
		employee.Name,
		employee.HireDate,
		employee.DailyPayment,
		employee.Phone,
		employee.Email,
		employee.Nss,
		employee.Curp,
	)
	if err != nil {
		return utils.DuplicatedDataError
	}
	return nil
}

func (s store) ReadEmployee(id int) (*types.Employee, error) {
	q := "SELECT id, name, hire_date, daily_payment, phone, email, nss, curp FROM employee WHERE id = ?;"

	var emp types.Employee
	err := s.db.QueryRow(q, id).Scan(&emp.Id, &emp.Name, &emp.HireDate, &emp.DailyPayment, &emp.Phone, &emp.Email, &emp.Nss, &emp.Curp)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (s store) UpdateEmployee(employee *types.Employee) error {
	q := `
	UPDATE employee
	SET name = ?, hire_date = ?, daily_payment = ?, phone = ?, email = ?, nss = ?, curp = ?
	WHERE id = ?;`

	_, err := s.db.Exec(q, employee.Name, employee.HireDate, employee.DailyPayment, employee.Phone, employee.Email, employee.Nss, employee.Curp, employee.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s store) DeleteEmployee(id int) error {
	q := "DELETE FROM employee WHERE id = ?;"
	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}
	return nil
}

func (s store) SearchEmployee(pattern string, limit, offset int) ([]*types.Employee, int, error) {
	formatedPattern := "%" + pattern + "%"

	q1 := `SELECT id, name, hire_date, daily_payment, phone, email, nss, curp FROM employee WHERE name LIKE ? LIMIT ? OFFSET ?;`
	q2 := `SELECT COUNT(id) AS count_id FROM employee WHERE name LIKE ?;`

	rows, err := s.db.Query(q1, formatedPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*types.Employee
	for rows.Next() {
		var record types.Employee
		if err := rows.Scan(&record.Id, &record.Name, &record.HireDate, &record.DailyPayment, &record.Phone, &record.Email, &record.Nss, &record.Curp); err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err() // Corregido el error original donde retornaba 'err' en lugar de 'rows.Err()'
	}

	var count int
	if err := s.db.QueryRow(q2, formatedPattern).Scan(&count); err != nil {
		return nil, 0, err
	}

	return records, count, nil
}
