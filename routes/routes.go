package routes

import (
	"go-auth-api/controllers"
	"go-auth-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/refresh", controllers.RefreshToken)

	// Routes ที่ต้องมี Token คุ้มกัน
	api.Post("/logout", middlewares.VerifyToken, controllers.Logout)
	api.Get("/protected", middlewares.VerifyToken, func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{"message": "This is a protected route", "userId": userID})
	})
}