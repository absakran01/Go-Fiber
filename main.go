package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main(){
	//init db
	db, err := InitDb()
	if err != nil {
		log.Panic(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})
	
	AuthHandlers(app.Group("/auth"), db)

	protected := app.Use(TokenAuthMiddleware(db))

	BookHandlers(protected.Group("/book"), db)

	app.Listen(":3000")
}