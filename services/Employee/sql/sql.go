package sql

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateSchema(db *sql.DB) error {

	// Comienzo de transacción SQL
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback de respaldo
	defer tx.Rollback()

	// Validar si la tabla ya existe
	var exists bool
	queryCheckEmployee := `SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='employee');`
	if err = tx.QueryRow(queryCheckEmployee).Scan(&exists); err != nil {
		log.Printf("Query error: (%v) for table (%v)", err.Error(), "employee")
		return err
	}

	// Migrar datos si la tabla ya existe. De lo contrario, crearla.
	if exists {
		employeeMigrate := `
		CREATE TABLE employee_new(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			hire_date TEXT NOT NULL,
			daily_payment REAL NOT NULL,
			phone TEXT,
			email TEXT UNIQUE,
			nss TEXT UNIQUE,
			curp TEXT UNIQUE
		);
		INSERT INTO employee_new (id, name, hire_date, daily_payment, phone, email, nss, curp)
		SELECT id, name, hire_date, daily_payment, phone, email, nss, curp FROM employee;
		DROP TABLE employee;
		ALTER TABLE employee_new RENAME TO employee;`
		if _, err := tx.Exec(employeeMigrate); err != nil {
			return fmt.Errorf("Migration error: %w", err)
		}
	} else {
		employeeCreate := `
		CREATE TABLE IF NOT EXISTS employee(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			hire_date TEXT NOT NULL,
			daily_payment REAL NOT NULL,
			phone TEXT,
			email TEXT UNIQUE,
			nss TEXT UNIQUE,
			curp TEXT UNIQUE
		);`
		if _, err := tx.Exec(employeeCreate); err != nil {
			return fmt.Errorf("Creation error: %w", err)
		}
	}

	return tx.Commit()
}
