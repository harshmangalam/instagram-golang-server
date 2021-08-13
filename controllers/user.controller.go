package controllers

import (
	"context"
	"fmt"
	"instagram/database"
	"instagram/ent"
	"instagram/ent/post"
	"instagram/ent/user"

	"github.com/gofiber/fiber/v2"
)

func FetchUserById(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	userId, err := c.ParamsInt("userId")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    nil,
		})
	}

	user, err := client.User.Get(ctx, userId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Fetch users",
		"data":    user,
	})
}

func FetchUserByUsername(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	username := c.Params("username")

	userData, err := client.User.Query().WithPosts(func(pq *ent.PostQuery) {
		pq.Select(post.FieldID, post.FieldImage).WithLikes(func(uq *ent.UserQuery) {
			uq.Select(user.FieldID)
		})
	}).WithFollowers(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID)
	}).WithFollowings(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID)
	}).Where(user.Username(username)).First(ctx)

	fmt.Println(userData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Fetch users",
		"data":    userData,
	})
}

func FetchUserSuggestion(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	userId := int(c.Locals("userId").(float64))

	users, err := client.User.Query().Where(user.IDNEQ(userId)).WithFollowers(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID)
	}).WithFollowings(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID)
	}).All(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "can`t fetch user suggestions",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Fetch user suggestions",
		"data":    users,
	})
}

func FollowUnfollowUser(c *fiber.Ctx) error {
	client := database.Client
	ctx := context.Background()

	userId, err := c.ParamsInt("userId")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	userData, err := client.User.Get(ctx, userId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	currentUserId := int(c.Locals("userId").(float64))

	currentUserData, err := client.User.Get(ctx, currentUserId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	hasAlreadyFollowing, err := currentUserData.QueryFollowings().Where(user.ID(userId)).Exist(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Can`t like",
			"data":    err.Error(),
		})
	}

	var message string

	if hasAlreadyFollowing {
		_, err = currentUserData.Update().RemoveFollowingIDs(userId).Save(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Can`t unfollow",
				"data":    err.Error(),
			})
		}

		message = "You are no longer following " + userData.Name
	} else {
		_, err = currentUserData.Update().AddFollowingIDs(userId).Save(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Can`t follow",
				"data":    err.Error(),
			})
		}

		message = "You are following " + userData.Name
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    nil,
	})

}
