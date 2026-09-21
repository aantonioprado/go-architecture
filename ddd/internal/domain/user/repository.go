package user

type Repository interface {
	Save(u *User) error
	FindAll() ([]*User, error)
	FindByID(id string) (*User, error)
	FindByEmail(email Email) (*User, error)
	Delete(id string) error
}
