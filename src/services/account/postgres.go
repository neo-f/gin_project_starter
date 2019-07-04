package account

import (
	"gin_project_starter/src/utils/token"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnMatch = errors.New("email or password is not correct")
)

type PostgresService struct {
	Conn *pg.DB
}

func (u *PostgresService) Login(email, password string) (obj Account, t string, err error) {
	err = u.Conn.Model(&obj).Where("email = ?", email).Select()
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(password))
	if err != nil {
		err = ErrUnMatch
		return
	}
	t, err = token.Create(obj.Email)
	if err != nil {
		return
	}
	_, _ = u.Conn.Model(&obj).Set("last_login_at = ?", time.Now()).Update()
	return
}

func (u *PostgresService) Update(obj *Account, columns ...string) error {
	_, err := u.Conn.Model(obj).Column(columns...).Returning("*").WherePK().Update()
	return err
}

func (u *PostgresService) Delete(id int64) error {
	_, err := u.Conn.Model(Account{ID: id}).WherePK().Delete()
	return err
}

func (u *PostgresService) Retrieve(id int64) (obj Account, err error) {
	err = u.Conn.Model(&u).Where("id = ?", id).Select()
	return
}

func (u *PostgresService) List(limit, offset int) (objs []Account, count int, err error) {
	objs = make([]Account, 0)
	count, err = u.Conn.Model(&objs).Limit(limit).Offset(limit).SelectAndCount()
	return
}

func (u *PostgresService) Create(username, password, email string) (obj Account, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return obj, err
	}
	obj = Account{
		Username: username,
		Password: string(hash),
		Email:    email,
	}
	_, err = u.Conn.Model(&obj).Insert()
	return obj, err
}
