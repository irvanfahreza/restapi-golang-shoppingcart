package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"ilmudata/project-golang/database"
	"ilmudata/project-golang/models"
)

var store = session.New()

type User struct {
	// This is not the model, more like a serializer
	ID       int    `json:"id"`
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type AuthController struct {
	// declare variables
	Db    *gorm.DB
	store *session.Store
}

// func CreateResponseLogin(input LoginForm) LoginForm {
// 	return LoginForm{Username: input.Username, Password: input.Password}
// }

func InitAuthController() *AuthController {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})

	return &AuthController{Db: db}
}

func CreateResponseUser(user models.User) User {
	return User{ID: user.ID, Username: user.Username, Email: user.Email, Password: user.Password}
}

func createAkun(db *gorm.DB, newUser *models.User) (err error) {
	err = db.Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (controller *AuthController) CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	sHash := string(bytes)

	user.Password = sHash

	err := createAkun(controller.Db, &user)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	errs := readAkunByUsername(controller.Db, &user, user.Username, user.Password)
	if errs != nil {
		return c.Status(500).JSON(errs.Error())
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func findUsername(username string, user *models.User) error {
	database.Database.Db.Where("username=?", username).First(user)
	// if user.ID == nil {
	// 	return errors.New("user does not exist")
	// }
	return nil
}

func readAkunByUsername(db *gorm.DB, user *models.User, username string, password string) (err error) {
	err = database.Database.Db.Where(&User{
		Username: username}, &User{Password: password}).First(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (controller *AuthController) LoginUser(c *fiber.Ctx) error {
	// Username := c.Params("username")
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	var input User
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// bytes, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	// sHash := string(bytes)

	// input.Password = sHash

	err = readAkunByUsername(controller.Db, &user, user.Username, user.Password)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// if err = database.Database.Db.First(&user).Error; err != nil {
	// 	return c.Status(400).JSON(err.Error())
	// }

	plainPass := input.Password
	dataPass := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(dataPass), []byte(plainPass))
	if err != nil {
		sess.Set("username", user.Username)
		// sess.Set("password", user.Password)
		sess.Save()

		responseUser := CreateResponseUser(user)
		return c.Status(200).JSON(responseUser)

	}

	return c.Status(400).JSON(err.Error())

	// database.Database.Db.First(&user)

}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("error")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("error")
	}

	err = findUser(id, &user)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		Username string `json:"username"`
		Email    string `json:Email`
		Password string `json:"password"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.Username = updateData.Username
	user.Email = updateData.Email
	user.Password = updateData.Password

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)

}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("error")
	}

	err = findUser(id, &user)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted User")
}
