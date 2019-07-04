package services

import (
	"gin_project_starter/src/services/account"
	"gin_project_starter/src/storages"
)

type AccountService interface {
	Retrieve(id int64) (account.Account, error)
	List(limit, offset int) ([]account.Account, int, error)
	Create(username, password, email string) (account.Account, error)
	Update(obj *account.Account, columns ...string) error
	Delete(id int64) error
	Login(email, password string) (user account.Account, t string, err error)
}

func NewAccountService() AccountService {
	conn := storages.NewStorage().Get("auth")
	return &account.PostgresService{Conn: conn}
}
