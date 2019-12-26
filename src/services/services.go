package services

import (
	"time"
)

type Account struct {
	ID          int64     `json:"id" pg:",pk"`
	Username    string    `json:"username" pg:",unique"`
	Password    string    `json:"-"`
	Email       string    `json:"email" pg:",unique"`
	CreatedAt   time.Time `json:"-" pg:"default:now()"`
	UpdatedAt   time.Time `json:"-"`
	LastLoginAt time.Time `json:"last_login_at"`
}

type AccountService interface {
	Retrieve(id int64) (Account, error)
	List(limit, offset int) ([]Account, int, error)
	Create(username, password, email string) (Account, error)
	Update(obj *Account, columns ...string) error
	Delete(id int64) error
	Login(email, password string) (user Account, t string, err error)
}
