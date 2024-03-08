package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB // Declare at the package level
type Config struct {
	DBHost     string `json:"host"`
	DBPort     int    `json:"port"`
	DBUser     string `json:"user"`
	DBPassword string `json:"password"`
	DBName     string `json:"dbname"`
}

func loadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
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
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Confirm database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Connected to the database successfully")

	// Initialize the Gin router
	router := gin.Default()

	// Define routes
	router.GET("/recipes", getRecipes)
	router.POST("/recipes", createRecipe)
	router.GET("/recipes/:id", getRecipe)
	router.PUT("/recipes/:id", updateRecipe)
	router.DELETE("/recipes/:id", deleteRecipe)

	// Start the server
	log.Fatal(router.Run(":8080"))
}

// Define struct for recipe
type Recipe struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

// Handler to get all recipes
func getRecipes(c *gin.Context) {
	// Fetch all recipes from the database
	rows, err := db.Query("SELECT id, title, ingredients, instructions FROM recipes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions); err != nil {
			log.Println(err)
			continue
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

// Handler to create a new recipe
func createRecipe(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Insert the new recipe into the database
	_, err := db.Exec("INSERT INTO recipes (title, ingredients, instructions) VALUES ($1, $2, $3)", recipe.Title, recipe.Ingredients, recipe.Instructions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Recipe created successfully"})
}

// Handler to get a recipe by ID
func getRecipe(c *gin.Context) {
	// Fetch the recipe from the database by ID
	var recipe Recipe
	id := c.Param("id")
	row := db.QueryRow("SELECT id, title, ingredients, instructions FROM recipes WHERE id = $1", id)
	err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// Handler to update a recipe by ID
func updateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the recipe in the database
	_, err := db.Exec("UPDATE recipes SET title = $1, ingredients = $2, instructions = $3 WHERE id = $4", recipe.Title, recipe.Ingredients, recipe.Instructions, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
}

// Handler to delete a recipe by ID
func deleteRecipe(c *gin.Context) {
	id := c.Param("id")

	// Delete the recipe from the database by ID
	_, err := db.Exec("DELETE FROM recipes WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
