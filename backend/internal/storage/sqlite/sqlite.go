package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sudhanshu042004/sandbox/internal/config"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg config.Config) (*Sqlite, error) { // create new instance of sqlite for ya
	db, err := sql.Open("sqlite3", cfg.Storage)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		password TEXT
		)`)

	if err != nil {
		return nil, fmt.Errorf("error while creating db instance %s", err)
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
