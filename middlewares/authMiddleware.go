package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(403).JSON(fiber.Map{"message": "No token provided!"})
	}

	// ตัดคำว่า "Bearer " ออก (ถ้ามี) หรือเอาตัว Token มาตรงๆ
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Failed to authenticate token."})
	}

	claims := token.Claims.(jwt.MapClaims)
	
	// ฝาก ID ของ User ไว้ใน Ctx เพื่อให้ Controller ตัวอื่นหยิบไปใช้ได้
	c.Locals("userID", claims["id"])

	return c.Next()
}