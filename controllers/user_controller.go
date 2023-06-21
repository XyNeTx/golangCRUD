package controllers

import (
	"context"
	"net/http"
	"time"
	"uplevel-api/configs"
	"uplevel-api/models"
	"uplevel-api/responses"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var access_tokenCollection *mongo.Collection = configs.GetCollection(configs.DB, "accessToken")
var validate = validator.New()

func GetValidToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	access_token := c.Query("access_token")
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"access_token": bson.M{"$regex": access_token}}},
	}
	cursor, err := access_tokenCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}
	defer cursor.Close(ctx)

	var results []models.AccessToken
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	if len(results) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "user not found", Data: nil})
	}

	return c.Status(http.StatusOK).JSON(
		responses.AccessTokenResponse{Status: http.StatusOK, Message: "success", Data: &results[0]},
	)
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"userId": userId}},
	}
	cursor, err := userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}
	defer cursor.Close(ctx)

	var results []models.User
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	if len(results) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "user not found", Data: nil})
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &results[0]},
	)
}
