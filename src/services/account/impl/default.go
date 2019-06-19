package impl

import (
	"gin_project_starter/src/services/account"
	"gin_project_starter/src/utils/token"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnMatch = errors.New("email or password is not correct")
)

func NewAccountService(conn *pg.DB) *AccountService {
	return &AccountService{conn: conn}
}

type AccountService struct {
	conn *pg.DB
}

func (u *AccountService) UserLogin(email, password string) (user account.User, t string, err error) {
	err = u.conn.Model(&user).Where("email = ?", email).Select()
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		err = ErrUnMatch
		return
	}
	t, err = token.Create(user.Email)
	if err != nil {
		return
	}
	_, _ = u.conn.Model(&user).Set("last_login_at = ?", time.Now()).Update()
	return
}

func (u *AccountService) UserUpdate(obj *account.User, columns ...string) error {
	_, err := u.conn.Model(obj).Column(columns...).Returning("*").WherePK().Update()
	return err
}

func (u *AccountService) UserDelete(id int64) error {
	_, err := u.conn.Model(account.User{ID: id}).WherePK().Delete()
	return err
}

func (u *AccountService) UserRetrieve(id int64) (user account.User, err error) {
	err = u.conn.Model(&u).Where("id = ?", id).Select()
	return
}

func (u *AccountService) UserList(limit, offset int) (users []account.User, count int, err error) {
	users = make([]account.User, 0)
	count, err = u.conn.Model(&users).Limit(limit).Offset(limit).SelectAndCount()
	return
}
