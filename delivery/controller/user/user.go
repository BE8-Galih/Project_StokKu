package user

import (
	"fmt"
	"net/http"
	"stokku/delivery/controller"
	"stokku/delivery/view"
	userV "stokku/delivery/view/user"
	"stokku/entities"
	"stokku/repository/user"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type UserController struct {
	Repo  user.UserDBControl
	Valid *validator.Validate
}

func NewUserControl(Ur user.UserDBControl, validate *validator.Validate) *UserController {
	return &UserController{
		Repo:  Ur,
		Valid: validate,
	}
}

// Method Untuk Membuat Data User Baru
func (u *UserController) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		newUser := userV.InsertUser{}
		fmt.Println(newUser)
		if err := c.Bind(&newUser); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusUnsupportedMediaType, view.BindData())
		}

		if err := u.Valid.Struct(newUser); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.Validate())
		}

		new := entities.User{Name: newUser.Name, Email: newUser.Email, Password: newUser.Password}

		res, err := u.Repo.CreateUser(new)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		LoginData := userV.LoginRespond{Data: res}
		token, _ := controller.CreateToken(res.ID)

		LoginData.Token = token
		return c.JSON(http.StatusCreated, LoginData)
	}
}

// Method Untuk Menampilkan Data User Berdasarkan ID
func (u *UserController) GetUserID() echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())
		}
		res, err2 := u.Repo.GetUserID(id)
		fmt.Println(res, err2)
		if err2 != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, userV.StatusGetIdOk(res))
	}
}

// Method Untuk Mengupdate Data User
func (u *UserController) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		newEdit := userV.InsertUser{}
		if err := c.Bind(&newEdit); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusUnsupportedMediaType, view.BindData())
		}
		idParam := c.Param("id")
		id, err2 := strconv.Atoi(idParam)
		if err2 != nil {
			log.Error(err2)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())
		}
		editData := entities.User{}
		if newEdit.Name != "" {
			editData.Name = newEdit.Name
		}
		if newEdit.Password != "" {
			editData.Password = newEdit.Password
		}
		if newEdit.Email != "" {
			editData.Email = newEdit.Email
		}
		errUpdate := u.Repo.UpdateUserID(editData, id)

		if errUpdate != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, userV.StatusUpdate(editData))
	}
}

// Method Untuk Menghapus User
func (u *UserController) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())
		}

		errDelete := u.Repo.DeleteUserID(id)
		if errDelete != nil {
			log.Warn(errDelete)
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, userV.StatusDelete())
	}
}

// Method Untuk Login User
func (u *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		identiti := userV.InsertLogin{}
		if err := c.Bind(&identiti); err != nil {
			return c.JSON(http.StatusUnsupportedMediaType, view.BindData())
		}

		if err := u.Valid.Struct(identiti); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.Validate())

		}

		res, err := u.Repo.Login(identiti.Email, identiti.Password)

		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		LoginData := userV.LoginRespond{Data: res}

		if LoginData.Token == "" {
			token, _ := controller.CreateToken(res.ID)
			LoginData.Token = token
			return c.JSON(http.StatusOK, userV.StatusLogin(LoginData))
		}
		return c.JSON(http.StatusOK, userV.StatusLogin(LoginData))
	}
}
