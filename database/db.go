package database

import (
	"go-auth-api/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database! \n", err)
	}

	log.Println("Database connected successfully!")
	DB = db

	// สร้างตารางอัตโนมัติ
	db.AutoMigrate(&models.User{})
}