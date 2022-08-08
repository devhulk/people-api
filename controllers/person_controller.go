package controllers

import (
	"context"
	"net/http"
	"people-api/configs"
	"people-api/models"
	"people-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var peopleCollection *mongo.Collection = configs.GetCollection(configs.DB, "people")
var validate = validator.New()

func CreatePerson(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var person models.Person
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&person); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PeopleResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&person); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PeopleResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newPerson := models.Person{
		Id:          primitive.NewObjectID(),
		FirstName:   person.FirstName,
		LastName:    person.LastName,
		Address:     person.Address,
		PhoneNumber: person.PhoneNumber,
	}

	result, err := peopleCollection.InsertOne(ctx, newPerson)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.PeopleResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAPerson(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	personId := c.Params("personId")
	var person models.Person
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(personId)

	err := peopleCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&person)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PeopleResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": person}})
}

func EditAPerson(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	personId := c.Params("personId")
	var person models.Person
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(personId)

	//validate the request body
	if err := c.BodyParser(&person); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PeopleResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&person); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PeopleResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	opts := options.Update().SetUpsert(true)

	update := bson.M{"firstname": person.FirstName, "lastname": person.LastName, "address": person.Address, "phonenumber": person.PhoneNumber}

	result, err := peopleCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update}, opts)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated person details
	var updatedPerson models.Person
	if result.MatchedCount == 1 {
		err := peopleCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPerson)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.PeopleResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedPerson}})
}

func DeleteAPerson(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	personId := c.Params("personId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(personId)

	result, err := peopleCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.PeopleResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Person with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PeopleResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Person successfully deleted!"}},
	)
}

func GetAllPeople(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var people []models.Person
	defer cancel()

	results, err := peopleCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePerson models.Person
		if err = results.Decode(&singlePerson); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PeopleResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		people = append(people, singlePerson)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PeopleResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": people}},
	)
}
