package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func TokenAuthMiddleware(db *gorm.DB) fiber.Handler{
	return func(c *fiber.Ctx) error{
		cToken := c.Cookies("jwt")
		var tokenString string
		if cToken != ""{
			log.Warn("using token from cookies...")
			tokenString = cToken
			log.Info(`
token from cookies:
----------------------------------------------------
`, tokenString, `
----------------------------------------------------`)
		}else {
			log.Warn("using token from header...")

			authHeader := c.Get("Authorization")
			if authHeader == ""{
				log.Error("Authorization header is empty")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unable to authorize",
				})
			}
			tokenSplitParts := strings.Split(authHeader, " ")

			if len(tokenSplitParts) != 2 || tokenSplitParts[0] != "Bearer"{
				log.Error("Authorization header is not in Bearer format")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unable to authorize",
				})
			}
			tokenString = tokenSplitParts[1]
			log.Info("token from header: ", tokenString)
		}

		secret := []byte(os.Getenv("AUTH_SECRET"))
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg(){
				log.Error("Unexpected signing method: ", t.Method.Alg())
				return nil, fmt.Errorf("unecpected signing method: %s", t.Method.Alg())
			}
			return secret, nil
		})

		if err != nil || !token.Valid{
			log.Error("Failed to parse token, ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{

					"error": "unable to authorize",
				})
		}

		userId := token.Claims.(jwt.MapClaims)["userId"]
		
		if err := db.Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound){
			if err != nil || !token.Valid{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unable to authorize, user assocciated with id not found",
				})
			}
		}
		c.Locals("userId", userId)
		return c.Next()
	}
}