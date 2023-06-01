package routes

import (
	"go-auth/database"
	"go-auth/models"
	"go-auth/sessions"
	"go-auth/utils"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func AttachAuthRoutesV1(app *fiber.App) {
	router := app.Group("/api/v1")

	router.Get("/me", meHandlerV1)
	router.Post("/login", loginHandlerV1)
	router.Post("/signup", signupHandlerV1)
	router.Post("/logout", logoutHandlerV1)
}

func meHandlerV1(c *fiber.Ctx) error {
	if sessions.RSS == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"uid": nil,
		})
	}

	sess, err := sessions.RSS.Get(c)
	if err != nil || sess == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"uid": nil,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"uid": sess.Get("uid"),
	})
}

func loginHandlerV1(c *fiber.Ctx) error {
	var body struct {
		Email    string
		Password string
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "bad request body: " + err.Error(),
		})
	}

	var user models.User
	database.DB.Model(&models.User{}).First(&user, "email = ?", body.Email)

	if reflect.DeepEqual(user, models.User{}) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err": "invalid credentials",
		})
	}
	if !utils.NewHasher("bcrypt").Compare([]byte(user.Password), []byte(body.Password)) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err": "invalid credentials",
		})
	}

	sess, err := sessions.RSS.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}
	sess.Set("uid", user.ID)

	if err := sess.Save(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": "unable to save user session: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "logged in successfully",
	})
}

func signupHandlerV1(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": "bad request body: " + err.Error(),
		})
	}

	var user models.User
	database.DB.Model(&models.User{}).Where("email = ? OR username = ?", body.Email, body.Username).First(&user)

	if !reflect.DeepEqual(user, models.User{}) {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": "email or username already exists",
		})
	}

	hasher := utils.NewHasher("bcrypt")
	raw, err := hasher.Hash([]byte(body.Password))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": "unable to hash user password: " + err.Error(),
		})
	}
	hashedPassword := string(raw)
	body.Password = hashedPassword

	if err := database.DB.Model(&models.User{}).Create(&body).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": "couldn't create new user: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": body.ID,
	})
}

func logoutHandlerV1(c *fiber.Ctx) error {
	if sessions.RSS == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"err": "already logged out",
		})
	}

	sess, err := sessions.RSS.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": "couldn't log out: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "logged out successfully",
	})
}
