package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cooksey14/go-recipe-blog/middleware"
	"github.com/cooksey14/go-recipe-blog/models"
	"github.com/cooksey14/go-recipe-blog/sendgrid"
	"github.com/cooksey14/go-recipe-blog/store"
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

func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

type Handler struct {
	Store *store.Store
}

// ListRecipes handler
func (h *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.Store.GetAllRecipes()
	if err != nil {
		http.Error(w, "Failed to fetch recipes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		http.Error(w, "Failed to encode recipes", http.StatusInternalServerError)
	}
}

// CreateRecipe handler
func (h *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.Store.CreateRecipe(recipe); err != nil {
		log.Println("Failed to create recipe:", err)
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Recipe created successfully")
}

// SignUpUser handler
func (h *Handler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	exists, err := h.Store.IsEmailExists(user.Email)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusBadRequest)
		return
	}
	if err := h.Store.CreateUser(user); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User signed up successfully")
}

// LoginUser handler
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	storedHash, err := h.Store.GetUserPasswordHash(creds.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
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

// SendEmail handler
func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var emailReq models.Email
	err = json.Unmarshal(body, &emailReq)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if emailReq.Address == "" || emailReq.Name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err = sendgrid.SendgridSendEmail(emailReq.Address, emailReq.Name)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email sent successfully")
}

// GetRecipe handler
func (h *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/recipes/")
	recipe, err := h.Store.GetRecipeByID(id)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// UpdateRecipe handler
func (h *Handler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.Store.UpdateRecipe(id, recipe); err != nil {
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recipe updated successfully")
}

// DeleteRecipe handler
func (h *Handler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	if err := h.Store.DeleteRecipe(id); err != nil {
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recipe deleted successfully")
}
