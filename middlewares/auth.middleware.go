package middlewares

import (
	"fmt"
	"instagram/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt"
)

func Protected() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(config.Get("JWT_SECRET")),
		ErrorHandler:   jwtError,
		SuccessHandler: successHandler,
		ContextKey:     "userId",
		AuthScheme:     "Bearer ",
		TokenLookup:    "cookie:token",
	})

}

func successHandler(c *fiber.Ctx) error {

	token := c.Locals("userId").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	userId := claims["userId"].(float64)

	c.Locals("userId", userId)

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {

	fmt.Println(c.Cookies("token"))
	fmt.Println(err.Error())

	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
