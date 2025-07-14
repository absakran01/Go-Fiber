package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(router fiber.Router, db *gorm.DB){
	router.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("post")
	})
	router.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("get by id")
	})
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("get all")
	})
	router.Put("/:id", func(c *fiber.Ctx) error {
		return c.SendString("put by id")
	})
	router.Delete("/:id", func(c *fiber.Ctx) error {
		return c.SendString("dlt by id")
	})
}