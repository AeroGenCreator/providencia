package store

import (
	"database/sql"
	"saavedra/services/CuotasIMSS/types"
	"saavedra/utils"
)

type Store interface {
	// Funciones fichero extra.go
	M2OEmployees() ([]*types.M2OEmployee, error)
	ReadEmployee(id int) (*types.Employee, error)
	// Funciones en este fichero
	SliceCuotasIMSS(limit, offset int) ([]*types.CuotasIMSS, int, error)
	CreateCuotasIMSS(cuota *types.CuotasIMSS) error
	ReadCuotasIMSS(id int) (*types.CuotasIMSS, error)
	UpdateCuotasIMSS(cuota *types.CuotasIMSS) error
	DeleteCuotasIMSS(id int) error
	SearchCuotasIMSS(pattern string, limit, offset int) ([]*types.CuotasIMSS, int, error)
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}

// Para la vista lista (Se envian pocos campos) (Importantes ID, Name, Year, Month)
func (s store) SliceCuotasIMSS(limit, offset int) ([]*types.CuotasIMSS, int, error) {
	q1 := "SELECT COUNT(id) AS count_id FROM cuotas_imss;"
	q2 := `SELECT id, name, daily_payment, year, month, factor_integracion, sbc FROM cuotas_imss LIMIT ? OFFSET ?;`
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
	var records []*types.CuotasIMSS
	for rows.Next() {
		var c types.CuotasIMSS
		err = rows.Scan(&c.Id, &c.Name, &c.DailyPayment, &c.Year, &c.Month, &c.FactorIntegracion, &c.Sbc)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, &c)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return records, count, nil
}

func (s store) CreateCuotasIMSS(c *types.CuotasIMSS) error {
	checkQuery := `SELECT id FROM cuotas_imss WHERE name = ? AND year = ? AND month = ?;`
	insertQuery := `
	INSERT INTO cuotas_imss (
		name, daily_payment, year, month, factor_integracion, sbc, es_bimestre, incapacidad, faltas,
		riesgo_trabajo_prima, enfermedad_maternidad_cuota_fija, enfermedad_maternidad_adicional,
		prestaciones_dinero_patron, gastos_medicos_pensionados_patron, invalidez_y_vida,
		retiro, cesantia_y_vejez_patron, guarderias, aportacion_infonavit,
		enfermedad_maternidad_adicional_obrero, prestaciones_dinero_obrero,
		gastos_medicos_pensionados_obrero, invalidez_y_vida_obrero, cesantia_y_vejez_obrero
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	var id int
	err := s.db.QueryRow(checkQuery, c.Name, c.Year, c.Month).Scan(&id)
	if err == nil {
		return utils.DuplicatedDataError
	}
	if err != sql.ErrNoRows {
		return err
	}
	_, err = s.db.Exec(
		insertQuery,
		c.Name, c.DailyPayment, c.Year, c.Month, c.FactorIntegracion, c.Sbc, c.EsBimestre, c.Incapacidad, c.Faltas,
		c.RiesgoTrabajoPrima, c.EnfermedadMaternidadCuotaFija, c.EnfermedadMaternidadAdicional,
		c.PrestacionesDineroPatron, c.GastosMedicosPensionadosPatron, c.InvalidezYVida,
		c.Retiro, c.CesantiaYVejezPatron, c.Guarderias, c.AportacionInfonavit,
		c.EnfermedadMaternidadAdicionalObrero, c.PrestacionesDineroObrero,
		c.GastosMedicosPensionadosObrero, c.InvalidezYVidaObrero, c.CesantiaYVejezObrero,
	)
	if err != nil {
		return utils.DuplicatedDataError
	}
	return nil
}

func (s store) ReadCuotasIMSS(id int) (*types.CuotasIMSS, error) {
	q := `SELECT
		id, name, daily_payment, year, month, factor_integracion, sbc, es_bimestre, incapacidad, faltas,
		riesgo_trabajo_prima, enfermedad_maternidad_cuota_fija, enfermedad_maternidad_adicional,
		prestaciones_dinero_patron, gastos_medicos_pensionados_patron, invalidez_y_vida,
		retiro, cesantia_y_vejez_patron, guarderias, aportacion_infonavit,
		enfermedad_maternidad_adicional_obrero, prestaciones_dinero_obrero,
		gastos_medicos_pensionados_obrero, invalidez_y_vida_obrero, cesantia_y_vejez_obrero
		FROM cuotas_imss WHERE id = ?;`
	var c types.CuotasIMSS
	err := s.db.QueryRow(q, id).Scan(
		&c.Id, &c.Name, &c.DailyPayment, &c.Year, &c.Month, &c.FactorIntegracion, &c.Sbc, &c.EsBimestre, &c.Incapacidad, &c.Faltas,
		&c.RiesgoTrabajoPrima, &c.EnfermedadMaternidadCuotaFija, &c.EnfermedadMaternidadAdicional,
		&c.PrestacionesDineroPatron, &c.GastosMedicosPensionadosPatron, &c.InvalidezYVida,
		&c.Retiro, &c.CesantiaYVejezPatron, &c.Guarderias, &c.AportacionInfonavit,
		&c.EnfermedadMaternidadAdicionalObrero, &c.PrestacionesDineroObrero,
		&c.GastosMedicosPensionadosObrero, &c.InvalidezYVidaObrero, &c.CesantiaYVejezObrero,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s store) UpdateCuotasIMSS(c *types.CuotasIMSS) error {
	q := `
	UPDATE cuotas_imss
	SET name = ?, daily_payment = ?, year = ?, month = ?, factor_integracion = ?, sbc = ?, es_bimestre = ?, incapacidad = ?, faltas = ?,
	    riesgo_trabajo_prima = ?, enfermedad_maternidad_cuota_fija = ?, enfermedad_maternidad_adicional = ?,
	    prestaciones_dinero_patron = ?, gastos_medicos_pensionados_patron = ?, invalidez_y_vida = ?,
	    retiro = ?, cesantia_y_vejez_patron = ?, guarderias = ?, aportacion_infonavit = ?,
	    enfermedad_maternidad_adicional_obrero = ?, prestaciones_dinero_obrero = ?,
	    gastos_medicos_pensionados_obrero = ?, invalidez_y_vida_obrero = ?, cesantia_y_vejez_obrero = ?
	WHERE id = ?;`
	_, err := s.db.Exec(
		q,
		c.Name, c.DailyPayment, c.Year, c.Month, c.FactorIntegracion, c.Sbc, c.EsBimestre, c.Incapacidad, c.Faltas,
		c.RiesgoTrabajoPrima, c.EnfermedadMaternidadCuotaFija, c.EnfermedadMaternidadAdicional,
		c.PrestacionesDineroPatron, c.GastosMedicosPensionadosPatron, c.InvalidezYVida,
		c.Retiro, c.CesantiaYVejezPatron, c.Guarderias, c.AportacionInfonavit,
		c.EnfermedadMaternidadAdicionalObrero, c.PrestacionesDineroObrero,
		c.GastosMedicosPensionadosObrero, c.InvalidezYVidaObrero, c.CesantiaYVejezObrero,
		c.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s store) DeleteCuotasIMSS(id int) error {
	q := "DELETE FROM cuotas_imss WHERE id = ?;"
	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}
	return nil
}

func (s store) SearchCuotasIMSS(pattern string, limit, offset int) ([]*types.CuotasIMSS, int, error) {
	formatedPattern := "%" + pattern + "%"
	q1 := `SELECT id, name, daily_payment, year, month, factor_integracion, sbc FROM cuotas_imss WHERE name LIKE ? LIMIT ? OFFSET ?;`
	q2 := `SELECT COUNT(id) AS count_id FROM cuotas_imss WHERE name LIKE ?;`
	rows, err := s.db.Query(q1, formatedPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var records []*types.CuotasIMSS
	for rows.Next() {
		var r types.CuotasIMSS
		if err := rows.Scan(
			&r.Id, &r.Name, &r.DailyPayment, &r.Year, &r.Month, &r.FactorIntegracion, &r.Sbc,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, &r)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	var count int
	if err := s.db.QueryRow(q2, formatedPattern).Scan(&count); err != nil {
		return nil, 0, err
	}
	return records, count, nil
}
