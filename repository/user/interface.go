package user

import "stokku/entities"

type UserDBControl interface {
	CreateUser(NewUser entities.User) (entities.User, error)
	GetUserID(id int) (entities.User, error)
	UpdateUserID(NewData entities.User, id int) error
	DeleteUserID(id int) error
	Login(email string, password string) (entities.User, error)
}
