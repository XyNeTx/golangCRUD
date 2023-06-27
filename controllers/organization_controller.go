package controllers

import (
	"context"
	"log"
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

func GetMyOrg(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result []models.MyOrg
	userId := c.Params("userId")
	// userID, err := primitive.ObjectIDFromHex(userId)
	pipeline := []bson.M{
		{"$match": bson.M{"members.userId": userId}},
	}
	myOrg, err := orgCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MyOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	defer myOrg.Close(ctx)

	if err := myOrg.All(ctx, &result); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MyOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	if len(result) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.MyOrgResponse{Status: http.StatusNotFound, Message: "data not found", Data: nil})
	}

	return c.Status(http.StatusOK).JSON(
		responses.MyOrgResponse{Status: http.StatusOK, Message: "success", Data: result},
	)
}

func GetAllOrgs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var result []models.AllOrgs
	defer cancel()

	allOrgs, err := orgCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AllOrgsResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	defer allOrgs.Close(ctx)
	for allOrgs.Next(ctx) {
		var singleOrg models.AllOrgs
		if err = allOrgs.Decode(&singleOrg); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AllOrgsResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
		}

		result = append(result, singleOrg)
	}

	if len(result) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.MyOrgResponse{Status: http.StatusNotFound, Message: "data not found", Data: nil})
	}

	return c.Status(http.StatusOK).JSON(
		responses.AllOrgsResponse{Status: http.StatusOK, Message: "success", Data: result},
	)
}

func GetOrgSummary(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orgId := c.Params("organizationId")

	pipeline := []bson.M{
		{"$match": bson.M{"organizationId": orgId}},
	}

	orgSummary, err := ideaCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OrgSummaryResponse{
			Status:  http.StatusInternalServerError,
			Message: "error79",
			Counts:  nil,
		})
	}
	defer orgSummary.Close(ctx)

	var result []models.OrgSummary

	if err := orgSummary.All(ctx, &result); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(responses.OrgSummaryResponse{
			Status:  http.StatusInternalServerError,
			Message: "error89",
			Counts:  nil,
		})
	}

	if len(result) == 0 {
		return c.Status(http.StatusNotFound).JSON(responses.OrgSummaryResponse{
			Status:  http.StatusNotFound,
			Message: "data not found",
			Counts:  nil,
		})
	}

	counts := make(map[string]int)
	benefits := make(map[string]int32)
	benefitSum := int32(0)

	for _, summary := range result {
		counts[summary.Status]++
		for _, idea := range summary.IdeaTemplate {
			benefits[summary.Status] += idea.Benefit
			benefitSum += idea.Benefit
		}
	}

	statusCounts := make([]models.StatusCount, 0, len(counts))
	for status, count := range counts {
		statusCounts = append(statusCounts, models.StatusCount{
			Status:  status,
			Count:   count,
			Benefit: benefits[status],
		})
	}

	return c.Status(http.StatusOK).JSON(responses.OrgSummaryResponse{
		Status:  http.StatusOK,
		Message: "success",
		Counts:  statusCounts,
		Benefit: benefitSum,
	})
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

	// userId := c.Params("userId")
	// pipeline := []bson.M{
	// 	{"$match": bson.M{"members.userId": userId}},
	// }
	// userOrg, err := userCollection.Aggregate(ctx, pipeline)

	newOrg := models.CreateOrg{
		Id:        primitive.NewObjectID(),
		Name:      org.Name,
		Member:    org.Member,
		Verifier:  org.Verifier,
		PowerUser: org.PowerUser,
	}

	result, err := orgCollection.InsertOne(ctx, newOrg)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	return c.Status(http.StatusCreated).JSON(responses.CreateOrgResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

// func EditOrg(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	organizationId := c.Params("organizationId")
// 	var editOrg models.CreateOrg
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(organizationId)

// 	//validate the request body
// 	if err := c.BodyParser(&editOrg); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{Status: http.StatusBadRequest, Message: "error", Data: nil})
// 	}

// 	//use the validator library to validate required fields
// 	if validationErr := validate.Struct(&editOrg); validationErr != nil {
// 		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{Status: http.StatusBadRequest, Message: "error", Data: nil})
// 	}

// 	update := bson.M{"name": editOrg.Name, "picture": editOrg.Picture, "members": editOrg.Member, "powerUser": editOrg.PowerUser, "verifier": editOrg.Verifier}

// 	result, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
// 	}
// 	//get updated user details
// 	var updatedOrg models.CreateOrg
// 	if result.MatchedCount == 1 {
// 		err := orgCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedOrg)

// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
// 		}
// 	}

//		return c.Status(http.StatusOK).JSON(responses.CreateOrgResponse{Status: http.StatusOK, Message: "success", Data: updatedOrg})
//	}

func EditOrg(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizationId := c.Params("organizationId")
	objId, _ := primitive.ObjectIDFromHex(organizationId)

	// Parse the request body into a new organization object
	var editOrg models.CreateOrg
	if err := c.BodyParser(&editOrg); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    nil,
		})
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&editOrg); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    nil,
		})
	}

	update := bson.M{}

	if editOrg.Name != "" {
		update["name"] = editOrg.Name
	}

	if editOrg.Picture != "" {
		update["picture"] = editOrg.Picture
	}

	update2 := bson.M{}

	if len(editOrg.Member) > 0 {
		update2["$push"] = bson.M{"members": bson.M{"$each": editOrg.Member}}
	}

	update3 := bson.M{}

	if len(editOrg.Verifier) > 0 {
		update3["$push"] = bson.M{"verifier": bson.M{"$each": editOrg.Verifier}}
	}

	update4 := bson.M{}

	if len(editOrg.PowerUser) > 0 {
		update4["$push"] = bson.M{"powerUser": bson.M{"$each": editOrg.PowerUser}}
	}

	// Update the organization document with the new data
	result, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	result2, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, update2)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}
	result3, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, update3)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}
	result4, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, update4)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    nil,
		})
	}

	// Retrieve the updated organization details
	var updatedOrg models.CreateOrg
	if result.MatchedCount == 1 && result2.MatchedCount == 1 && result3.MatchedCount == 1 && result4.MatchedCount == 1 {
		err := orgCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedOrg)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CreateOrgResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    updatedOrg,
	})
}

// func EditOrg(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	organizationId := c.Params("organizationId")
// 	objId, _ := primitive.ObjectIDFromHex(organizationId)

// 	// Parse the request body into a new organization object
// 	var editOrg models.CreateOrg
// 	if err := c.BodyParser(&editOrg); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "error",
// 			Data:    nil,
// 		})
// 	}

// 	// Use the validator library to validate required fields
// 	if validationErr := validate.Struct(&editOrg); validationErr != nil {
// 		return c.Status(http.StatusBadRequest).JSON(responses.CreateOrgResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "error",
// 			Data:    nil,
// 		})
// 	}

// 	// Create the update fields based on the new organization data
// 	update := bson.M{}

// 	if editOrg.Name != "" {
// 		update["name"] = editOrg.Name
// 	}

// 	if editOrg.Picture != "" {
// 		update["picture"] = editOrg.Picture
// 	}

// 	if len(editOrg.Member) > 0 {
// 		update["$push"] = bson.M{"members": bson.M{"$each": editOrg.Member}}
// 	}

// 	if len(editOrg.PowerUser) > 0 {
// 		update["$push"] = bson.M{"powerUser": bson.M{"$each": editOrg.PowerUser}}
// 	}

// 	if len(editOrg.Verifier) > 0 {
// 		update["$push"] = bson.M{"verifier": bson.M{"$each": editOrg.Verifier}}
// 	}

// 	// if editOrg.Member != nil {
// 	// 	update["members"] = editOrg.Member
// 	// }

// 	// if editOrg.PowerUser != nil {
// 	// 	update["powerUser"] = editOrg.PowerUser
// 	// }

// 	// if editOrg.Verifier != nil {
// 	// 	update["verifier"] = editOrg.Verifier
// 	// }

// 	// Update the organization document with the new data
// 	result, err := orgCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error3",
// 			Data:    nil,
// 		})
// 	}

// 	// 	// Retrieve the updated organization details
// 	// 	var updatedOrg models.CreateOrg
// 	// 	if result.MatchedCount == 1 {
// 	// 		err := orgCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedOrg)
// 	// 		if err != nil {
// 	// 			return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
// 	// 				Status:  http.StatusInternalServerError,
// 	// 				Message: "error4",
// 	// 				Data:    nil,
// 	// 			})
// 	// 		}
// 	// 	}

// 	// 	return c.Status(http.StatusOK).JSON(responses.CreateOrgResponse{
// 	// 		Status:  http.StatusOK,
// 	// 		Message: "success",
// 	// 		Data:    updatedOrg,
// 	// 	})
// 	// }

// 	var updatedOrg models.CreateOrg
// 	if result.MatchedCount == 1 {
// 		err := orgCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedOrg)
// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{
// 				Status:  http.StatusInternalServerError,
// 				Message: "error",
// 				Data:    nil,
// 			})
// 		}
// 	}

// 	return c.Status(http.StatusOK).JSON(responses.CreateOrgResponse{
// 		Status:  http.StatusOK,
// 		Message: "success",
// 		Data:    updatedOrg,
// 	})
// }

func DeleteOrg(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	organizationId := c.Params("organizationId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(organizationId)

	result, err := orgCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CreateOrgResponse{Status: http.StatusInternalServerError, Message: "error", Data: nil})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CreateOrgResponse{Status: http.StatusNotFound, Message: "error", Data: nil},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CreateOrgResponse{Status: http.StatusOK, Message: "Organization successfully deleted!", Data: nil},
	)
}
