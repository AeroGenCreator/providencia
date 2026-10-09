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
	queryCheckUsers := `SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='users');`
	if err = tx.QueryRow(queryCheckUsers).Scan(&exists); err != nil {
		log.Printf("Query error: (%v) for table (%v)", err.Error(), "users")
		return err
	}

	// Migrar datos, si la tabla ya existe. De lo contrario crearla.
	if exists {
		usersMigrate := `
		CREATE TABLE users_new(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		party TEXT NOT NULL
		);
		INSERT INTO users_new (id, name, email, password, party)
		SELECT id, name, email, password, party FROM users;
		DROP TABLE users;
		ALTER TABLE users_new RENAME TO users`
		if _, err := tx.Exec(usersMigrate); err != nil {
			return fmt.Errorf("Migration error: %w", err)
		}
	} else {
		usersCreate := `
		CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		party TEXT NOT NULL
		);`
		if _, err := tx.Exec(usersCreate); err != nil {
			return fmt.Errorf("Creation error: %w", err)
		}
	}

	return tx.Commit()
}
