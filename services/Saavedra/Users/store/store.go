package store

import (
	"database/sql"
	"saavedra/services/Saavedra/Users/types"
)

type Store interface {
	SliceUser(limit, offset int) ([]*types.User, int, error)
	CreateUser(user *types.User) error
	ReadUser(id int) (*types.User, error)
	UpdateUser(user *types.User) error
	UpdateUserNoPassword(user *types.User) error
	DeleteUser(id int) error
	SearchUser(pattern string, limit, offset int) ([]*types.User, int, error)
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return store{db: db}
}

func (s store) SliceUser(limit, offset int) ([]*types.User, int, error) {
	q1 := "SELECT COUNT(id) AS count_id FROM users;"
	q2 := "SELECT id, name, email, party FROM users LIMIT ? OFFSET ?;"
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
	var records []*types.User
	for rows.Next() {
		var user types.User
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Party)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, &user)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return records, count, nil
}

func (s store) CreateUser(user *types.User) error {
	q := `
	INSERT INTO users(
	name, email, password, party)
	VALUES (?, ?, ?, ?)
	ON CONFLICT DO NOTHING;`
	res, err := s.db.Exec(q, user.Name, user.Email, user.Password, user.Party)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = int(id)
	return nil
}

func (s store) ReadUser(id int) (*types.User, error) {
	q := "SELECT id, name, email, party FROM users WHERE id = ?;"
	var user types.User
	err := s.db.QueryRow(q, id).Scan(&user.Id, &user.Name, &user.Email, &user.Party)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s store) UpdateUser(user *types.User) error {
	q := "UPDATE users SET name = ?, email = ?, password = ?, party = ? WHERE id = ?;"
	_, err := s.db.Exec(q, user.Name, user.Email, user.Password, user.Party, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s store) UpdateUserNoPassword(user *types.User) error {
	q := "UPDATE users SET name = ?, email = ?, party = ? WHERE id = ?;"
	_, err := s.db.Exec(q, user.Name, user.Email, user.Party, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s store) DeleteUser(id int) error {
	q := "DELETE FROM users WHERE id = ?;"
	if _, err := s.db.Exec(q, id); err != nil {
		return err
	}
	return nil
}

func (s store) SearchUser(pattern string, limit, offset int) ([]*types.User, int, error) {
	formatedPattern := "%" + pattern + "%"
	q1 := `SELECT id, name, email, party FROM users WHERE name LIKE ? LIMIT ? OFFSET ?;`
	q2 := `SELECT COUNT(id) AS count_id FROM users WHERE name LIKE ?;`
	rows, err := s.db.Query(q1, formatedPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var records []*types.User
	for rows.Next() {
		var record types.User
		if err := rows.Scan(&record.Id, &record.Name, &record.Email, &record.Party); err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	var count int
	if err := s.db.QueryRow(q2, formatedPattern).Scan(&count); err != nil {
		return nil, 0, err
	}
	return records, count, nil
}
