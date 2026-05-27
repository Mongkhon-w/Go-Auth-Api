package controllers

import (
	"go-auth-api/database"
	"go-auth-api/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid data"})
	}

	// บดรหัสผ่าน (Hash)
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 10)

	user := models.User{
		Username: data["username"],
		Password: string(password),
	}

	// บันทึกลง DB
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Username already exists"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid data"})
	}

	var user models.User
	database.DB.Where("username = ?", data["username"]).First(&user)
	if user.ID == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// เทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	// สร้าง Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	accessToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	// สร้าง Refresh Token
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	refreshToken, _ := refreshTokenString.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))

	// Hash Refresh Token ก่อนลง DB
	hashedRefresh, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	user.RefreshToken = string(hashedRefresh)
	database.DB.Save(&user)

	// ส่ง Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"accessToken": accessToken})
}

func RefreshToken(c *fiber.Ctx) error {
	// อ่าน Token จาก Cookie
	cookie := c.Cookies("refreshToken")
	if cookie == "" {
		return c.Status(403).JSON(fiber.Map{"message": "Refresh Token is required!"})
	}

	// แกะกล่อง JWT
	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(403).JSON(fiber.Map{"message": "Invalid Refresh Token!"})
	}

	var user models.User
	database.DB.Where("id = ?", claims["id"]).First(&user)

	// เทียบ Hash Refresh Token ใน DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(cookie)); err != nil {
		return c.Status(403).JSON(fiber.Map{"message": "Invalid Refresh Token!"})
	}

	// สร้าง Access Token ใบใหม่
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	newAccessToken, _ := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return c.JSON(fiber.Map{"accessToken": newAccessToken})
}

func Logout(c *fiber.Ctx) error {
	// สมมติว่าดึง userID มาจาก Context ของ Middleware (เดี๋ยวเราจะทำในสเตปต่อไป)
	userID := c.Locals("userID")

	var user models.User
	database.DB.Where("id = ?", userID).First(&user)
	user.RefreshToken = "" // เคลียร์เป็นค่าว่าง
	database.DB.Save(&user)

	// ล้าง Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}