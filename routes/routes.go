// routes.go

package routes

import (
	"database/sql"
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
)

// SetupRoutes sets up the HTTP routes for the application
func SetupRoutes(db *sql.DB) {
	//List all Recipes
	http.HandleFunc("/recipes", handlers.ListRecipes(db))
	// Get a Recipe
	http.HandleFunc("/recipes/", handlers.GetRecipe(db))
	// Create a new Recipe
	http.HandleFunc("/recipes/create", handlers.CreateRecipe(db))
	// Update a Recipe
	http.HandleFunc("/recipes/update/", handlers.UpdateRecipe(db))
	// Delete a Recipe
	http.HandleFunc("/recipes/delete/", handlers.DeleteRecipe(db))
}
