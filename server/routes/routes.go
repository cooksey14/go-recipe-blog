package routes

import (
	"net/http"

	"github.com/cooksey14/go-recipe-blog/handlers"
	"github.com/cooksey14/go-recipe-blog/middleware"
)

// SetupRoutes sets up the HTTP routes for the application
func SetupRoutes(handler *handlers.Handler) {
	// Public endpoints
	http.Handle("/signup", handlers.HandleCORS(http.HandlerFunc(handler.SignUpUser)))
	http.Handle("/login", handlers.HandleCORS(http.HandlerFunc(handler.LoginUser)))
	http.Handle("/recipes", handlers.HandleCORS(http.HandlerFunc(handler.ListRecipes)))
	http.Handle("/recipes/", handlers.HandleCORS(http.HandlerFunc(handler.GetRecipe)))

	// Protected endpoints
	http.Handle("/recipes/create", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handler.CreateRecipe))))
	http.Handle("/recipes/update/", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handler.UpdateRecipe))))
	http.Handle("/recipes/delete/", handlers.HandleCORS(middleware.JwtVerify(http.HandlerFunc(handler.DeleteRecipe))))
}
