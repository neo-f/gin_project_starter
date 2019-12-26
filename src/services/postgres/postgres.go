package postgres

import (
	"gin_project_starter/src/services"
	"gin_project_starter/src/storages"
	"gin_project_starter/src/utils/token"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnMatch = errors.New("email or password is not correct")
)

type AccountService struct {
	Conn *pg.DB
}

func NewAccountService() services.AccountService {
	conn := storages.Get("auth")
	return &AccountService{Conn: conn}
}

func (u *AccountService) Login(email, password string) (obj services.Account, t string, err error) {
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

func (u *AccountService) Update(obj *services.Account, columns ...string) error {
	_, err := u.Conn.Model(obj).Column(columns...).Returning("*").WherePK().Update()
	return err
}

func (u *AccountService) Delete(id int64) error {
	_, err := u.Conn.Model(services.Account{ID: id}).WherePK().Delete()
	return err
}

func (u *AccountService) Retrieve(id int64) (obj services.Account, err error) {
	err = u.Conn.Model(&u).Where("id = ?", id).Select()
	return
}

func (u *AccountService) List(limit, offset int) (objs []services.Account, count int, err error) {
	objs = make([]services.Account, 0)
	count, err = u.Conn.Model(&objs).Limit(limit).Offset(limit).SelectAndCount()
	return
}

func (u *AccountService) Create(username, password, email string) (obj services.Account, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return obj, err
	}
	obj = services.Account{
		Username: username,
		Password: string(hash),
		Email:    email,
	}
	_, err = u.Conn.Model(&obj).Insert()
	return obj, err
}
