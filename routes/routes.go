// routes.go

package routes

import (
	"database/sql"
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
)

// SetupRoutes sets up the HTTP routes for the application
func SetupRoutes(db *sql.DB) {
	http.HandleFunc("/recipes", handlers.GetRecipes(db))
	http.HandleFunc("/recipes/", handlers.GetRecipe(db))
	http.HandleFunc("/recipes/create", handlers.CreateRecipe(db))
	http.HandleFunc("/recipes/update/", handlers.UpdateRecipe(db))
	http.HandleFunc("/recipes/delete/", handlers.DeleteRecipe(db))
}
