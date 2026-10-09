package store

import (
	"database/sql"
	"saavedra/services/Saavedra/Session/types"
)

type Store interface {
	ReadSession(id int) (*types.Session, error)
	UpdateSession(id int, token string) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}

func (s store) ReadSession(id int) (*types.Session, error) {
	q := "SELECT id, user_id, token FROM session WHERE user_id = ?;"
	var dbSession types.Session
	err := s.db.QueryRow(q, id).Scan(&dbSession.Id, &dbSession.UserId, &dbSession.Token)
	if err == sql.ErrNoRows {
		dbSession.Token = ""
		return &dbSession, nil
	} else if err != nil {
		return nil, err
	}
	return &dbSession, nil
}

func (s store) UpdateSession(id int, token string) error {
	q := "UPDATE session SET token = ? WHERE user_id = ?"
	_, err := s.db.Exec(q, token, id)
	if err != nil {
		return err
	}
	return nil
}
