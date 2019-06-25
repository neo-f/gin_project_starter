package account

type Service interface {
	Retrieve(id int64) (Account, error)
	List(limit, offset int) ([]Account, int, error)
	Create(username, password, email string) (Account, error)
	Update(obj *Account, columns ...string) error
	Delete(id int64) error
	Login(email, password string) (user Account, t string, err error)
}
