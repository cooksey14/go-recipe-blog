package routes

import (
	"database/sql"
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
	"github.com/cooksey14/go-recipe-blog/middleware"
)

// SetupRoutes sets up the HTTP routes for the application
func SetupRoutes(db *sql.DB) {
	// Public endpoints
	http.Handle("/signup", handlers.HandleCORS(http.HandlerFunc(handlers.SignUpUser(db))))
	http.Handle("/login", handlers.HandleCORS(http.HandlerFunc(handlers.LoginUser(db))))
	http.Handle("/recipes", handlers.HandleCORS(http.HandlerFunc(handlers.ListRecipes(db))))
	http.Handle("/recipes/", handlers.HandleCORS(http.HandlerFunc(handlers.GetRecipe(db))))

	// Protected endpoints
	http.Handle("/recipes/create", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handlers.CreateRecipe(db)))))
	http.Handle("/recipes/update/", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handlers.UpdateRecipe(db)))))
	http.Handle("/recipes/delete/", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handlers.DeleteRecipe(db)))))
}
