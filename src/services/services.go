package services

import (
	"time"
)

type Account struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
