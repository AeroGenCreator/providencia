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
	queryCheckSession := `SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='session');`
	if err = tx.QueryRow(queryCheckSession).Scan(&exists); err != nil {
		log.Printf("Query error: (%v) for table (%v)", err.Error(), "session")
		return err
	}

	// Migrar datos, si la tabla ya existe. De lo contrario crearla.
	if exists {
		sessionMigrate := `
		CREATE TABLE session_new(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL UNIQUE,
		token TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		INSERT INTO session_new (id, user_id, token)
		SELECT id, user_id, token FROM session;
		DROP TABLE session;
		ALTER TABLE session_new RENAME TO session`
		if _, err := tx.Exec(sessionMigrate); err != nil {
			return fmt.Errorf("Migration error: %w", err)
		}
	} else {
		sessionCreate := `
		CREATE TABLE session(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL UNIQUE,
		token TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`
		if _, err := tx.Exec(sessionCreate); err != nil {
			return fmt.Errorf("Creation error: %w", err)
		}
	}
	return tx.Commit()
}
