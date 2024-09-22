package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
	"github.com/cooksey14/go-recipe-blog/middleware"
	"github.com/cooksey14/go-recipe-blog/routes"
	"github.com/cooksey14/go-recipe-blog/store"

	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pressly/goose/v3"
)

func main() {
	middleware.LoadEnvConfig()
	db_conn := middleware.GetEnv("DATABASE_URL")
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
