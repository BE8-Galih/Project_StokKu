package user

type UserControll interface {
	GetUserID() error
	CreateUser() error
	UpdateUser() error
	DeleteUser() error
	Login() error
}
