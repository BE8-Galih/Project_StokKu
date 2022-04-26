package user

import (
	"errors"
	"stokku/entities"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UsersDB struct {
	Db *gorm.DB
}

func NewDBUser(db *gorm.DB) *UsersDB {
	return &UsersDB{
		Db: db,
	}
}

//Method Untuk Menambahkan Data User di Database
func (u *UsersDB) CreateUser(NewUser entities.User) (entities.User, error) {
	if err := u.Db.Create(&NewUser).Error; err != nil {
		log.Warn(err)
		return entities.User{}, errors.New("Cannot Create Data")
	}
	log.Info()
	return NewUser, nil
}

//Method Untuk Mengambil Data User di Database
func (u *UsersDB) GetUserID(id int) (entities.User, error) {
	User := entities.User{}
	if err := u.Db.Where("id = ?", id).Find(&User).Error; err != nil {
		log.Warn(err)
		return entities.User{}, errors.New("Cannot Access Database")
	}
	return User, nil
}

// Method Untuk Mengupdate Data User di Database
func (u *UsersDB) UpdateUserID(NewData entities.User, id int) error {
	if err := u.Db.Table("users").Where("id = ?", id).Updates(entities.User{Name: NewData.Name, Email: NewData.Email, Password: NewData.Password}).Error; err != nil {
		log.Warn(err)
		return errors.New("Error Updating Data")
	}
	log.Info()
	return nil
}

// Method Untuk Menghapus Data User di Database
func (u *UsersDB) DeleteUserID(id int) error {
	if err := u.Db.Delete(&entities.User{}, id).Error; err != nil {
		log.Warn(err)
		return errors.New("Error Access Database")
	}
	log.Info()
	return nil
}

//Method Untuk Authentifikasi Ke Database Ketika User Login
func (u *UsersDB) Login(email string, password string) (entities.User, error) {
	identiti := entities.User{}
	if err := u.Db.Where("email = ? AND password = ?", email, password).First(&identiti).Error; err != nil {
		log.Warn(err)
		return entities.User{}, errors.New("Error Access Database")
	}
	return identiti, nil
}
