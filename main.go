package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cooksey14/go-recipe-blog/models"
	"github.com/cooksey14/go-recipe-blog/routes"
	_ "github.com/lib/pq"
)

func loadConfig(filename string) (models.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return models.Config{}, err
	}
	defer file.Close()

	var config models.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return models.Config{}, err
	}

	return config, nil
}

func main() {
	// Load configuration from config.json
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded configuration: %+v", config)

	// Connect to the database using configuration values
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Confirm database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Connected to the database successfully")

	// Setup HTTP routes with the database connection
	routes.SetupRoutes(db)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
