package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ilmudata/project-golang/controllers"
	"ilmudata/project-golang/database"
)

func setupRoutes(app *fiber.App) {
	authController := controllers.InitAuthController()
	// User
	app.Post("/api/users", authController.CreateUser)
	app.Get("/api/users", controllers.GetUsers)
	app.Get("/api/users/:id", controllers.GetUser)
	app.Delete("/api/users/:id", controllers.DeleteUser)
	app.Post("/api/login", authController.LoginUser)

	// Product
	app.Post("/api/products", controllers.CreateProduct)
	app.Get("/api/products", controllers.GetProducts)
	app.Get("/api/products/:id", controllers.GetProduct)
	app.Put("/api/products/:id", controllers.UpdateProduct)
	app.Delete("/api/products/:id", controllers.DeleteProductById)

	// Order
	app.Post("/api/orders", controllers.CreateOrder)
	app.Get("/api/orders", controllers.GetOrders)
	app.Get("/api/orders/:id", controllers.GetOrder)
}

func main() {
	database.ConnectDb()

	app := fiber.New()
	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
