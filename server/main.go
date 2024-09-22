package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cooksey14/go-recipe-blog/handlers"
	"github.com/cooksey14/go-recipe-blog/routes"
	"github.com/cooksey14/go-recipe-blog/store"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or unable to load it")
	}

	db_conn := getEnv("DATABASE_URL")
	db, err := sql.Open("postgres", db_conn)
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

	// Initialize the store and handlers
	store := store.NewStore(db)
	handler := handlers.NewHandler(store)

	// Set up routes
	routes.SetupRoutes(handler)

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	} else {
		log.Println("Server started on :8080")
	}
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

// Helper function to get required environment variables
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s is not set", key)
	return ""
}

// Helper function to get required integer environment variables
func getEnvAsInt(key string) int {
	valueStr := getEnv(key)
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer", key)
	}
	return valueInt
}
