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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orgCollection *mongo.Collection = configs.GetCollection(configs.DB, "organization")
var ideaCollection *mongo.Collection = configs.GetCollection(configs.DB, "ideas")
var validateOrg = validator.New()

// func GetOrgSummary(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	var orgSummary []models.OrgSummary
// 	organizationId := c.Params("organizationId")
// 	// userID, err := primitive.ObjectIDFromHex(userId)
// 	pipeline := []bson.M{
// 		{"$match": bson.M{"organizationId": organizationId}},
// 	}
// 	results, err := ideaCollection.Aggregate(ctx, pipeline)

// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.OrgSummaryResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
// 	}

// 	defer results.Close(ctx)

// 	if err := results.All(ctx, &orgSummary); err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.OrgSummaryResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error",
// 			Data:    nil,
// 		})
// 	}

// 	return c.Status(http.StatusOK).JSON(
// 		responses.MyOrgResponse{Status: http.StatusOK, Message: "success", Data: orgSummary},
// 	)
// }

func GetMyOrg(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var myOrg []models.MyOrg
	userId := c.Params("userId")
	// userID, err := primitive.ObjectIDFromHex(userId)
	pipeline := []bson.M{
		{"$match": bson.M{"members.userId": userId}},
	}
	results, err := orgCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MyOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	defer results.Close(ctx)

	if err := results.All(ctx, &myOrg); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MyOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	return c.Status(http.StatusOK).JSON(
		responses.MyOrgResponse{Status: http.StatusOK, Message: "success", Data: myOrg},
	)
}

func GetAllOrgs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var allOrgs []models.AllOrgs
	defer cancel()

	results, err := orgCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AllOrgsResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleOrg models.AllOrgs
		if err = results.Decode(&singleOrg); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AllOrgsResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}

		allOrgs = append(allOrgs, singleOrg)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AllOrgsResponse{Status: http.StatusOK, Message: "success", Data: allOrgs},
	)
}

func CreateOrg(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var org models.CreateOrg
	//validate the request body
	if err := c.BodyParser(&org); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{Status: http.StatusBadRequest, Message: "error", Data: nil})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&org); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{Status: http.StatusBadRequest, Message: "error", Data: nil})
	}

	newOrg := models.CreateOrg{
		Id:   primitive.NewObjectID(),
		Name: org.Name,
	}

	result, err := orgCollection.InsertOne(ctx, newOrg)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	return c.Status(http.StatusCreated).JSON(responses.CreateOrgResponse{Status: http.StatusCreated, Message: "success", Data: result})
}
