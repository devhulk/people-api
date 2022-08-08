package routes

import (
	"people-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func PeopleRoute(app *fiber.App) {
	//All routes related to people comes here
	app.Post("/people", controllers.CreatePerson)
	app.Get("/people/:personId", controllers.GetAPerson)
	app.Put("/people/:personId", controllers.EditAPerson)
	app.Delete("/people/:personId", controllers.DeleteAPerson)
	app.Get("/people", controllers.GetAllPeople)
}
