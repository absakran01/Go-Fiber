package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(router fiber.Router, db *gorm.DB){
	router.Post("/", func(c *fiber.Ctx) error {
		book := new(Book)

		book.UserId = int(c.Locals("userId").(float64))


		if err := c.BodyParser(book); err != nil {
			//handle book json parsing error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if book.Title == "" || book.Status == ""{
				// handle incomplete data 
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "book title and status are required",
			})
		}

		if err := db.Create(book).Error; err != nil{
			//handle db error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}



		log.Println("Book created successfully, with id:", book.Id)
		log.Println("Book name:", book.Title)
		log.Println("Book status:", book.Status)

		return c.Status(fiber.StatusCreated).SendString("book created successfuly")
	})


	router.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id") 

		if id == 0 || err != nil{
			// handle id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "book id is invalid",
			})
		}

		userId := int(c.Locals("userId").(float64))
		log.Println("User ID from context:", userId)

		if userId == 0{
			// handle user id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "user id is invalid",
			})
		}


		book := new(Book)
		db.Where("id = ? AND user_id = ?", id, userId).First(book)

		if book.Id == 0{
			// handle book not found 
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "could not Find book",
			})
		}



		return c.Status(fiber.StatusOK).JSON(book)
	})


	router.Get("/", func(c *fiber.Ctx) error {


		userId := int(c.Locals("userId").(float64))
		log.Println("User ID from context:", userId)

		if userId == 0{
			// handle user id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "user id is invalid",
			})
		}


		books := new([]Book)
		db.Where("user_id = ?", userId).Find(books)

		if len(*books) == 0{
			// handle book not found 
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "could not Find any books under user id",
			})
		}



		return c.Status(fiber.StatusOK).JSON(books)
	})


	router.Put("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id") 
		bookFromClient := new(Book)


		if err := c.BodyParser(bookFromClient); err != nil {
			//handle book json parsing error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if id == 0 || err != nil{
			// handle id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "book id is invalid",
			})
		}

		userId := int(c.Locals("userId").(float64))
		log.Println("User ID from context:", userId)

		if userId == 0{
			// handle user id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "user id is invalid",
			})
		}
		bookFromClient.Id = id



		if bookFromClient.Id == 0{
			// handle book not found 
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "could not Find book",
			})
		}

		if err := db.Save(&bookFromClient).Error; err != nil {
			// handle db error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(bookFromClient)
	})


	router.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id") 

		if id == 0 || err != nil{
			// handle id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "book id is invalid",
			})
		}

		userId := int(c.Locals("userId").(float64))
		log.Println("User ID from context:", userId)

		if userId == 0{
			// handle user id error
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "user id is invalid",
			})
		}


		book := new(Book)
		db.Where("id = ? AND user_id = ?", id, userId).First(book)

		if book.Id == 0{
			// handle book not found 
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "could not Find book",
			})
		}

		if err := db.Delete(&book).Error; err != nil{
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not Delete book",
			})
		}



		return c.Status(fiber.StatusNoContent).JSON(book)
	})
} 