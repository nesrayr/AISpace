package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nesrayr/handlers"
	"os"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("middleware")
		fmt.Printf("Request URL: %s\n", c.OriginalURL())
		tokenString := c.Cookies("jwt")
		//if tokenString == "" {
		//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Missing JWT token"})
		//}
		if tokenString == "" {
			return c.Next()
		}
		fmt.Println("tokenString:", tokenString)
		token, _ := jwt.ParseWithClaims(tokenString, &handlers.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("CLIENT_SECRET")), nil
		})
		fmt.Println(tokenString)
		//if err != nil {
		//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		//}
		claims := token.Claims.(*handlers.Claims)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}
