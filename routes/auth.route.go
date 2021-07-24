package routes

import (
	"instagram/controllers"
	"instagram/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoute(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/login", controllers.Login)
	auth.Get("/check_username/:username", controllers.CheckUsername)

	auth.Get("/me", middlewares.Protected(), controllers.CurrentUser)
	auth.Get("/logout", middlewares.Protected(), controllers.Logout)

}
