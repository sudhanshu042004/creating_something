package sqlite

import (
	"database/sql"

	"github.com/sudhanshu042004/sandbox/internal/config"
)

type Sqlite struct {
	Db *sql.DB
}

func NewUser(cfg config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Address)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXITS user(
	id INTEGER PRIMARY KEY AUTOINCREMENT
	name TEXT
	email TEXT
	password TEXT
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateUser(name string, email string, password string) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO user (name,email,password) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, email, password)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
