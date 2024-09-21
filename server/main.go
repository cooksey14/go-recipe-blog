package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cooksey14/go-recipe-blog/models"
	"github.com/cooksey14/go-recipe-blog/routes"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pressly/goose/v3"
)

func main() {
	config := models.Config{
		DBHost:     getEnv("DBHost", "localhost"),
		DBPort:     getEnvAsInt("DBPort", 5432),
		DBUser:     getEnv("DBUser", "postgres"),
		DBPassword: getEnv("DBPassword", "postgres"),
		DBName:     getEnv("DBName", "recipes_db"),
	}
	log.Printf("Loaded configuration: %+v", config)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	} else {
		log.Println("Connected to the database successfully")
	}

	// Run migrations
	runMigrations(db)

	// Set up routes
	routes.SetupRoutes(db)

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// runMigrations applies database migrations on startup
func runMigrations(db *sql.DB) {
	migrationsDir := "./migrations"

	goose.SetDialect("postgres")

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations applied successfully")
	}
}

// Helper function to get environment variables with a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variables as an int
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
