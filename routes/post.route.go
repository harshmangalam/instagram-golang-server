package routes

import (
	"instagram/controllers"
	"instagram/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitPostRoute(route fiber.Router) {
	post := route.Group("/post")

	post.Get("/", controllers.FetchPosts)
	post.Get("/:postId", controllers.FetchPostById)
	post.Patch("/:postId/heart", middlewares.Protected(), controllers.AddRemoveHeart)
	post.Post("/", middlewares.Protected(), controllers.CreatePost)
	post.Delete("/:postId", middlewares.Protected(), controllers.DeletePost)
}
