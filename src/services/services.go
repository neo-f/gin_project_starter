package services

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

type Account struct {
	ID       int64  `json:"id" pg:",pk"`
	Username string `json:"username" pg:",unique"`
	Password string `json:"-"`
	Email    string `json:"email" pg:",unique"`

	CreatedAt   time.Time   `json:"-" pg:"default:now(),notnull"`
	UpdatedAt   pg.NullTime `json:"-" pg:"type:timestamptz"`
	LastLoginAt pg.NullTime `json:"last_login_at" pg:"type:timestamptz"`
}

func (a *Account) BeforeUpdate(ctx context.Context) (context.Context, error) {
	a.UpdatedAt = pg.NullTime{Time: time.Now()}
	return ctx, nil
}

type AccountService interface {
	Retrieve(id int64) (Account, error)
	List(limit, offset int) ([]Account, int, error)
	Create(username, password, email string) (Account, error)
	Update(obj *Account, columns ...string) error
	Delete(id int64) error
	Login(email, password string) (user Account, t string, err error)
}
