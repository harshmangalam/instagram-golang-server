package routes

import (
	"instagram/controllers"
	"instagram/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(route fiber.Router) {

	user := route.Group("/user")

	user.Get("/suggestions", middlewares.Protected(), controllers.FetchUserSuggestion)
	user.Get("/:userId", controllers.FetchUserById)
	user.Get("/username/:username", controllers.FetchUserByUsername)
	user.Patch("/:userId/followUnfollow", middlewares.Protected(), controllers.FollowUnfollowUser)

}
