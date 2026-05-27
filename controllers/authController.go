package controllers

import (
	"crypto/sha256"
	"encoding/hex"
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

	// บดรหัสผ่านด้วย bcrypt (รหัสผ่านปกติยังคงใช้ bcrypt ได้เพราะสั้น)
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 10)

	user := models.User{
		Username: data["username"],
		Password: string(password),
	}

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

	// ✅ ใช้ SHA-256 ในการ Hash Refresh Token (แก้ปัญหาความยาวเกิน)
	hash := sha256.Sum256([]byte(refreshToken))
	user.RefreshToken = hex.EncodeToString(hash[:])
	database.DB.Save(&user)

	// ส่ง HttpOnly Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	cookie.HTTPOnly = true
	cookie.Path = "/" // บังคับให้ Cookie ใช้ได้กับทุก Route
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"accessToken": accessToken})
}

func RefreshToken(c *fiber.Ctx) error {
	cookie := c.Cookies("refreshToken")
	if cookie == "" {
		return c.Status(403).JSON(fiber.Map{"message": "Refresh Token is required!"})
	}

	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(403).JSON(fiber.Map{"message": "Invalid Refresh Token!"})
	}

	var user models.User
	database.DB.Where("id = ?", claims["id"]).First(&user)
	if user.ID == 0 {
		return c.Status(403).JSON(fiber.Map{"message": "User not found!"})
	}

	// ✅ Hash Cookie ที่รับมาด้วย SHA-256 เพื่อนำไปเทียบกับฐานข้อมูล
	hash := sha256.Sum256([]byte(cookie))
	hashedCookie := hex.EncodeToString(hash[:])

	if user.RefreshToken != hashedCookie {
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
	userID := c.Locals("userID")

	var user models.User
	database.DB.Where("id = ?", userID).First(&user)
	user.RefreshToken = ""
	database.DB.Save(&user)

	// ลบ Cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.HTTPOnly = true
	cookie.Path = "/"
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}