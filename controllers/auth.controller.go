package controllers

import (
	"context"
	"instagram/database"
	"instagram/ent/user"
	"instagram/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {

	type FormData struct {
		Name     string
		Username string
		Password string
	}

	formData := new(FormData)

	if err := c.BodyParser(formData); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid signup input field", "data": err})

	}

	client := database.Client

	ctx := context.Background()

	isUsernameExists, err := client.User.Query().Where(user.Username(formData.Username)).Exist(ctx)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn`t signup", "data": err})
	}

	if isUsernameExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "User already exists", "data": nil})
	}

	hashPassword, err := utils.HashPassword(formData.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't signup", "data": err})
	}

	newUser, err := client.User.Create().SetUsername(formData.Username).SetPassword(hashPassword).SetName(formData.Name).Save(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't signup", "data": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Account created successfully", "data": newUser})
}

func Login(c *fiber.Ctx) error {

	type FormData struct {
		Username string
		Password string
	}

	formData := new(FormData)

	if err := c.BodyParser(formData); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid login field", "data": err})

	}

	client := database.Client

	ctx := context.Background()
	dbUser, err := client.User.Query().Where(user.Username(formData.Username)).First(ctx)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Incorrect username/password credentials", "data": err})
	}

	if !utils.MatchHashPassword(dbUser.Password, formData.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Incorrect username/password credentials", "data": err})
	}

	dbUser, err = dbUser.Update().SetIsActive(true).Save(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err})
	}

	token, err := utils.CreateJWTToken(dbUser.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn`t Login", "data": err})
	}

	cookie := new(fiber.Cookie)

	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(12 * time.Hour)

	c.Cookie(cookie)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Login successfully", "data": dbUser})

}

func CheckUsername(c *fiber.Ctx) error {
	client := database.Client
	ctx := context.Background()

	isUsernameExists, err := client.User.Query().Where(user.Username(c.Params("username"))).Exist(ctx)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if isUsernameExists {
		return c.SendStatus(fiber.StatusNotAcceptable)
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func CurrentUser(c *fiber.Ctx) error {

	userId := int(c.Locals("userId").(float64))
	client := database.Client

	ctx := context.Background()
	user, err := client.User.Get(ctx, userId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "fetch current user", "data": user})
}

func Logout(c *fiber.Ctx) error {

	client := database.Client

	ctx := context.Background()

	userId := int(c.Locals("userId").(float64))

	user, err := client.User.Get(ctx, userId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	_, err = user.Update().SetIsActive(false).Save(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Can`t logout",
			"data":    err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name: "token",
		// Set expiry date to the past
		Expires: time.Now().Add(-(time.Hour * 2)),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logout successfully",
		"data":    nil,
	})
}
