package sql

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateCuotasImssSchema(db *sql.DB) error {

	// Comienzo de transacción SQL
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback de respaldo
	defer tx.Rollback()

	// Validar si la tabla ya existe
	var exists bool
	queryCheckTable := `SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='cuotas_imss');`
	if err = tx.QueryRow(queryCheckTable).Scan(&exists); err != nil {
		log.Printf("Query error: (%v) for table (%v)", err.Error(), "cuotas_imss")
		return err
	}

	// Migrar datos si la tabla ya existe. De lo contrario, crearla.
	if exists {
		tableMigrate := `
		CREATE TABLE cuotas_imss_new(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			daily_payment REAL NOT NULL,
			year INTEGER NOT NULL,
			month INTEGER NOT NULL,
			factor_integracion REAL NOT NULL,
			sbc REAL NOT NULL,
			es_bimestre BOOLEAN NOT NULL,
			incapacidad INTEGER NOT NULL,
			faltas INTEGER NOT NULL,
			riesgo_trabajo_prima REAL NOT NULL,
			enfermedad_maternidad_cuota_fija REAL NOT NULL,
			enfermedad_maternidad_adicional REAL NOT NULL,
			prestaciones_dinero_patron REAL NOT NULL,
			gastos_medicos_pensionados_patron REAL NOT NULL,
			invalidez_y_vida REAL NOT NULL,
			retiro REAL NOT NULL,
			cesantia_y_vejez_patron REAL NOT NULL,
			guarderias REAL NOT NULL,
			aportacion_infonavit REAL NOT NULL,
			enfermedad_maternidad_adicional_obrero REAL NOT NULL,
			prestaciones_dinero_obrero REAL NOT NULL,
			gastos_medicos_pensionados_obrero REAL NOT NULL,
			invalidez_y_vida_obrero REAL NOT NULL,
			cesantia_y_vejez_obrero REAL NOT NULL
		);
		INSERT INTO cuotas_imss_new (
			id, name, daily_payment, year, month, factor_integracion, sbc, es_bimestre, incapacidad, faltas,
			riesgo_trabajo_prima, enfermedad_maternidad_cuota_fija, enfermedad_maternidad_adicional,
			prestaciones_dinero_patron, gastos_medicos_pensionados_patron, invalidez_y_vida,
			retiro, cesantia_y_vejez_patron, guarderias, aportacion_infonavit,
			enfermedad_maternidad_adicional_obrero, prestaciones_dinero_obrero,
			gastos_medicos_pensionados_obrero, invalidez_y_vida_obrero, cesantia_y_vejez_obrero
		)
		SELECT
			id, name, daily_payment, year, month, factor_integracion, sbc, es_bimestre, incapacidad, faltas,
			riesgo_trabajo_prima, enfermedad_maternidad_cuota_fija, enfermedad_maternidad_adicional,
			prestaciones_dinero_patron, gastos_medicos_pensionados_patron, invalidez_y_vida,
			retiro, cesantia_y_vejez_patron, guarderias, aportacion_infonavit,
			enfermedad_maternidad_adicional_obrero, prestaciones_dinero_obrero,
			gastos_medicos_pensionados_obrero, invalidez_y_vida_obrero, cesantia_y_vejez_obrero
		FROM cuotas_imss;
		DROP TABLE cuotas_imss;
		ALTER TABLE cuotas_imss_new RENAME TO cuotas_imss;`
		if _, err := tx.Exec(tableMigrate); err != nil {
			return fmt.Errorf("Migration error: %w", err)
		}
	} else {
		tableCreate := `
		CREATE TABLE IF NOT EXISTS cuotas_imss(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			daily_payment REAL NOT NULL,
			year INTEGER NOT NULL,
			month INTEGER NOT NULL,
			factor_integracion REAL NOT NULL,
			sbc REAL NOT NULL,
			es_bimestre BOOLEAN NOT NULL,
			incapacidad INTEGER NOT NULL,
			faltas INTEGER NOT NULL,
			riesgo_trabajo_prima REAL NOT NULL,
			enfermedad_maternidad_cuota_fija REAL NOT NULL,
			enfermedad_maternidad_adicional REAL NOT NULL,
			prestaciones_dinero_patron REAL NOT NULL,
			gastos_medicos_pensionados_patron REAL NOT NULL,
			invalidez_y_vida REAL NOT NULL,
			retiro REAL NOT NULL,
			cesantia_y_vejez_patron REAL NOT NULL,
			guarderias REAL NOT NULL,
			aportacion_infonavit REAL NOT NULL,
			enfermedad_maternidad_adicional_obrero REAL NOT NULL,
			prestaciones_dinero_obrero REAL NOT NULL,
			gastos_medicos_pensionados_obrero REAL NOT NULL,
			invalidez_y_vida_obrero REAL NOT NULL,
			cesantia_y_vejez_obrero REAL NOT NULL
		);`
		if _, err := tx.Exec(tableCreate); err != nil {
			return fmt.Errorf("Creation error: %w", err)
		}
	}

	return tx.Commit()
}
