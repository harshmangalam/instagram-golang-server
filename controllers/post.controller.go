package controllers

import (
	"context"

	"instagram/database"
	"instagram/ent"
	"instagram/ent/user"

	"github.com/gofiber/fiber/v2"
)

func FetchPosts(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	posts, err := client.Post.Query().WithCreator(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID, user.FieldName, user.FieldUsername, user.FieldProfilePic)
	}).WithLikes(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID)
	}).All(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "can`t fetch posts",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Fetch posts",
		"data":    posts,
	})
}

func FetchPostById(c *fiber.Ctx) error {
	postId, err := c.ParamsInt("postId")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    nil,
		})
	}

	client := database.Client
	ctx := context.Background()

	post, err := client.Post.Get(ctx, postId)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Fetch post by ID",
		"data":    post,
	})
}

func CreatePost(c *fiber.Ctx) error {

	type PostBody struct {
		Image string
	}

	postBody := new(PostBody)

	if err := c.BodyParser(postBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid post field input",
			"data":    err.Error(),
		})
	}

	client := database.Client
	ctx := context.Background()

	userId := int(c.Locals("userId").(float64))
	post, err := client.Post.Create().SetImage(postBody.Image).SetCreatorID(userId).Save(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn`t create post",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Post created successfully",
		"data":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	userId := int(c.Locals("userId").(float64))
	postId, err := c.ParamsInt("postId")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    err.Error(),
		})
	}

	post, err := client.Post.Get(ctx, postId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    err.Error(),
		})
	}

	postCreator, err := post.QueryCreator().First(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    err.Error(),
		})
	}

	if postCreator.ID != userId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not the owner of this post",
			"data":    nil,
		})
	}

	err = client.Post.DeleteOneID(postId).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn`t delete post",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
		"data":    nil,
	})
}

func AddRemoveHeart(c *fiber.Ctx) error {

	client := database.Client
	ctx := context.Background()

	postId, err := c.ParamsInt("postId")

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    err.Error(),
		})
	}

	postData, err := client.Post.Get(ctx, postId)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Post not found",
			"data":    err.Error(),
		})
	}

	userId := int(c.Locals("userId").(float64))

	hasAlreadyLiked, err := postData.QueryLikes().Where(user.ID(userId)).Exist(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Can`t like",
			"data":    err.Error(),
		})
	}

	var message string

	if hasAlreadyLiked {
		_, err = postData.Update().RemoveLikeIDs(userId).Save(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Can`t remove heart",
				"data":    err.Error(),
			})
		}

		message = "Heart removed from  post"

	} else {
		_, err = postData.Update().AddLikeIDs(userId).Save(ctx)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Can`t add heart",
				"data":    err.Error(),
			})
		}

		message = "Heart added to post"

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    postData,
	})
}
