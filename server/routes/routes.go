package routes

import (
	"database/sql"
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
)

// SetupRoutes sets up the HTTP routes for the application
func SetupRoutes(db *sql.DB) {
	// List all Recipes
	http.Handle("/recipes", handlers.HandleCORS(http.HandlerFunc(handlers.ListRecipes(db))))
	// Get a Recipe
	http.Handle("/recipes/", handlers.HandleCORS(http.HandlerFunc(handlers.GetRecipe(db))))
	// Create a new Recipe
	http.Handle("/recipes/create", handlers.HandleCORS(http.HandlerFunc(handlers.CreateRecipe(db))))
	// Update a Recipe
	http.Handle("/recipes/update/", handlers.HandleCORS(http.HandlerFunc(handlers.UpdateRecipe(db))))
	// Delete a Recipe
	http.Handle("/recipes/delete/", handlers.HandleCORS(http.HandlerFunc(handlers.DeleteRecipe(db))))
}
