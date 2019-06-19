package account

type Service interface {
	UserRetrieve(id int64) (User, error)
	UserList(limit, offset int) ([]User, int, error)
	UserUpdate(obj *User, columns ...string) error
	UserDelete(id int64) error
	UserLogin(email, password string) (user User, t string, err error)
}
