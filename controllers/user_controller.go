package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"uplevel-api/configs"
	"uplevel-api/models"
	"uplevel-api/responses"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var access_tokenCollection *mongo.Collection = configs.GetCollection(configs.DB, "accessToken")
var validate = validator.New()

// func GetValidToken(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	access_token := c.Query("access_token")
// 	defer cancel()

// 	pipeline := []bson.M{
// 		{"$match": bson.M{"access_token": bson.M{"$regex": access_token}}},
// 	}
// 	cursor, err := access_tokenCollection.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
// 	}
// 	defer cursor.Close(ctx)

// 	var results []models.AccessToken
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
// 	}

// 	if len(results) == 0 {
// 		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "user not found", Data: nil})
// 	}

// 	return c.Status(http.StatusOK).JSON(
// 		responses.AccessTokenResponse{Status: http.StatusOK, Message: "success", Data: &results[0]},
// 	)
// }

func GetValidToken(c *fiber.Ctx) error {
	config := configs.AppConfig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	access_token := c.Query("access_token")

	// Make HTTP request to external API
	resp, err := http.Get(config.LineToken + access_token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	// Parse the response body
	var response struct {
		// Define the structure of the response JSON here
		// Modify it according to the actual response format
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	// Check the response data in MongoDB
	pipeline := []bson.M{
		{"$match": bson.M{"access_token": bson.M{"$regex": access_token}}},
	}
	cursor, err := access_tokenCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}
	defer cursor.Close(ctx)

	var results []models.AccessToken
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	if len(results) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "user not found",
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responses.AccessTokenResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &results[0],
	})
}

func GetUserByID(c *fiber.Ctx) error {
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
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    &results[0]},
	)
}

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	// Validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    nil,
		})
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    nil,
		})
	}

	// Check if user with the same userId already exists
	existingUser := models.User{}
	err := userCollection.FindOne(ctx, bson.M{"userId": user.UserId}).Decode(&existingUser)
	if err == nil {
		// User with the same userId already exists, return an error
		return c.Status(http.StatusConflict).JSON(responses.UserResponse{
			Status:  http.StatusConflict,
			Message: "user already exists",
			Data:    nil,
		})
	} else if err != mongo.ErrNoDocuments {
		// Error occurred while querying the database
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to check user existence",
			Data:    nil,
		})
	}

	if user.Profile.Language == "" && (user.Profile.Organization == nil || *user.Profile.Organization == "") {
		// Assign default values to profile
		user.Profile.Language = "en"
		emptyString := ""
		user.Profile.Organization = &emptyString
	}

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		UserId:      user.UserId,
		DisplayName: user.DisplayName,
		Picture:     user.Picture,
		PictureURL:  user.PictureURL,
		Profile:     user.Profile,
		Industry:    user.Industry,
		Province:    user.Province,
	}

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create user",
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &newUser,
	})
}
