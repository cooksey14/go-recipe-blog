package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cooksey14/go-recipe-blog/middleware"
	"github.com/cooksey14/go-recipe-blog/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(recipes); err != nil {
			http.Error(w, "Failed to encode recipes", http.StatusInternalServerError)
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

// SignUpUser handles user registration
func SignUpUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", user.Email).Scan(&exists)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User signed up successfully")
	}
}

// LoginUser authenticates the user and returns a JWT token
func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.User
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var storedHash string
		err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", creds.Email).Scan(&storedHash)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"email": creds.Email,
			"exp":   time.Now().Add(72 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(middleware.SignKey)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}

// Handler to get a recipe by ID
func GetRecipe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe
		id := strings.TrimPrefix(r.URL.Path, "/recipes/")
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
		fmt.Fprintf(w, "Recipe updated successfully")
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
		fmt.Fprintf(w, "Recipe deleted successfully")
	}
}
