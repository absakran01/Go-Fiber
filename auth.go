package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func AuthHandlers(router fiber.Router, db *gorm.DB){
	router.Post("/login", func(c *fiber.Ctx) error{
		user := &User{
			Username: c.FormValue("username"), 
			Password: c.FormValue("password"),
		}
		userFromDB := &User{}
		

		if user.Username == "" || user.Password == ""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user is empty",
			})}
		db.Where("username = ?", user.Username).First(userFromDB)


		if userFromDB.Id ==0{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		log.Printf("Login attempt: %+v\n", user)



		err :=bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Password",
			})
		}


		secret := []byte(os.Getenv("AUTH_SECRET"))
		expTime, _ := strconv.Atoi(os.Getenv("AUTH_EXP_TIME"))
		method := jwt.SigningMethodHS256
		claims := jwt.MapClaims{
			"userId": user.Id, 
			"username": user.Username, 
			"exp": expTime,
		}

		token, err := jwt.NewWithClaims(method, claims).SignedString(secret)


		if err != nil{
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: !c.IsFromLocal(),
		Secure:   !c.IsFromLocal(),
		MaxAge:   expTime,
	})


		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token" : token,
		})
	})
	router.Post("/register", func(c *fiber.Ctx) error{
		user := &User{
			Username: c.FormValue("username"), 
			Password: c.FormValue("password"),
		}
		if user.Username == "" || user.Password == ""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user is empty",
			})}


		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPass)

		db.Create(user)

		secret := []byte(os.Getenv("AUTH_SECRET"))
		expTime, _ := strconv.Atoi(os.Getenv("AUTH_EXP_TIME"))
		method := jwt.SigningMethodHS256
		claims := jwt.MapClaims{
			"userId": user.Id, 
			"username": user.Username, 
			"exp": expTime,
		}

		token, err := jwt.NewWithClaims(method, claims).SignedString(secret)


		if err != nil{
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: !c.IsFromLocal(),
		Secure:   !c.IsFromLocal(),
		MaxAge:   expTime,
	})


		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token" : token,
		})
	})
	//..
}