package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cooksey14/go-recipe-blog/models"
)

// CORS middleware function to handle preflight requests
func HandleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, HX-Request, HX-Current-URL, HX-Target, HX-Trigger, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handler to get all recipes
func ListRecipes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, ingredients, instructions FROM recipes")
		if err != nil {
			http.Error(w, "Failed to fetch recipes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var recipes []models.Recipe
		for rows.Next() {
			var recipe models.Recipe
			if err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions); err != nil {
				log.Println(err)
				continue
			}
			recipes = append(recipes, recipe)
		}

		w.Header().Set("Content-Type", "text/html")
		for _, recipe := range recipes {
			fmt.Fprintf(w, "<div><h2>%s</h2><p>%s</p><p>%s</p></div>", recipe.Title, recipe.Ingredients, recipe.Instructions)
		}
	}
}

// CreateRecipe handles the creation of a new recipe
func CreateRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO recipes (title, ingredients, instructions) VALUES ($1, $2, $3)", recipe.Title, recipe.Ingredients, recipe.Instructions)
		if err != nil {
			log.Println("Failed to create recipe:", err)
			http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Recipe created successfully")
	}
}

// Handler to get a recipe by ID
func GetRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe
		id := r.URL.Path[len("/recipes/"):]
		row := db.QueryRow("SELECT id, title, ingredients, instructions FROM recipes WHERE id = $1", id)
		err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions)
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)
	}
}

// Handler to update a recipe by ID
func UpdateRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/recipes/update/")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}

		var recipe models.Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE recipes SET title = $1, ingredients = $2, instructions = $3 WHERE id = $4", recipe.Title, recipe.Ingredients, recipe.Instructions, idInt)
		if err != nil {
			http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Handler to delete a recipe by ID
func DeleteRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/recipes/delete/")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM recipes WHERE id = $1", idInt)
		if err != nil {
			http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
