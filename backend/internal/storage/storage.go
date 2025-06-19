package storage

import "github.com/sudhanshu042004/sandbox/internal/types"

type Storage interface {
	CreateUser(name string, email string, password string) (int64, error)
	GetUser(email string) (*types.User, error)
}
