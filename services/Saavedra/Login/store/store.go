package store

import (
	"database/sql"
	sessionTypes "saavedra/services/Saavedra/Session/types"
	usersTypes "saavedra/services/Saavedra/Users/types"
)

type Store interface {
	ReadUser(user *usersTypes.UserStr) (*usersTypes.User, error)
	CreateSession(session *sessionTypes.Session) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}

func (s store) ReadUser(user *usersTypes.UserStr) (*usersTypes.User, error) {
	q := `
	SELECT id, name, email, password, party
	FROM users
	WHERE email = ?;`
	var dbUser usersTypes.User
	err := s.db.QueryRow(
		q, user.Email).Scan(
		&dbUser.Id,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.Party)
	if err == sql.ErrNoRows {
		dbUser.Name = ""
		dbUser.Email = ""
		dbUser.Password = ""
		dbUser.Party = ""
		return &dbUser, nil
	} else if err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func (s store) CreateSession(session *sessionTypes.Session) error {
	query := `
	INSERT INTO session (user_id, token)
	VALUES (?, ?)
	ON CONFLICT(user_id) DO UPDATE SET token = excluded.token;
	`
	_, err := s.db.Exec(query, session.UserId, session.Token)
	if err != nil {
		return err
	}
	return nil
}
