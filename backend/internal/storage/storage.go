package storage

type Storage interface {
	CreateUser(name string, email string, password string) (int64, error)
}
