package store

import "saavedra/services/CuotasIMSS/types"

// Listar empleados para selección rapida en frontend
func (s store) M2OEmployees() ([]*types.M2OEmployee, error) {
	q := `SELECT id, name FROM employee;`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*types.M2OEmployee
	for rows.Next() {
		var record types.M2OEmployee
		if err = rows.Scan(&record.Id, &record.Name); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return records, nil
}

// Obtener Registro de Empleado Seleccionado
func (s store) ReadEmployee(id int) (*types.Employee, error) {
	q := `SELECT
	id, name, hire_date, daily_payment, phone, email, nss, curp
	FROM employee
	WHERE id = ?;`
	var employee types.Employee
	if err := s.db.QueryRow(q, id).Scan(
		&employee.Id,
		&employee.Name,
		&employee.HireDate,
		&employee.DailyPayment,
		&employee.Phone,
		&employee.Email,
		&employee.Nss,
		&employee.Curp,
	); err != nil {
		return nil, err
	}
	return &employee, nil
}
