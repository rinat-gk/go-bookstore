package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rinatgk/go-fiber-postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-fiber-postgres/models"

	"github.com/gofiber/fiber/v2"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	DB = db
}

func main() {
	// initialize the Database
	InitDB()

	// Migrate the User model
	DB.AutoMigrate(&models.User{})

	// Start the Fiber app
	app := fiber.New()

	//Define routes
	app.Get("/users", GetUsers)
	app.Get("/users/:id", GetUser)
	app.Post("/users", CreateUser)
	app.Put("/users/:id", UpdateUser)
	app.Delete("/users/:id", DeleteUser)

	app.Listen(":3000")
}

// Create Handlers

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	DB.Find(&users)
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if result := DB.First(&user, id); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.JSON(user)
}
