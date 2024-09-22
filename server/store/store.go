package store

import (
	"database/sql"
	"errors"

	"github.com/cooksey14/go-recipe-blog/models"
	"golang.org/x/crypto/bcrypt"
)

// Store struct to hold the database connection
type Store struct {
	DB *sql.DB
}

// NewStore creates a new Store instance
func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

// Recipe-related methods
func (s *Store) GetAllRecipes() ([]models.Recipe, error) {
	rows, err := s.DB.Query("SELECT id, title, ingredients, instructions FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (s *Store) CreateRecipe(recipe models.Recipe) error {
	_, err := s.DB.Exec(
		"INSERT INTO recipes (title, ingredients, instructions) VALUES ($1, $2, $3)",
		recipe.Title, recipe.Ingredients, recipe.Instructions,
	)
	return err
}

func (s *Store) GetRecipeByID(id string) (models.Recipe, error) {
	var recipe models.Recipe
	row := s.DB.QueryRow(
		"SELECT id, title, ingredients, instructions FROM recipes WHERE id = $1", id,
	)
	err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Ingredients, &recipe.Instructions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe, errors.New("recipe not found")
		}
		return recipe, err
	}
	return recipe, nil
}

func (s *Store) UpdateRecipe(id int, recipe models.Recipe) error {
	_, err := s.DB.Exec(
		"UPDATE recipes SET title = $1, ingredients = $2, instructions = $3 WHERE id = $4",
		recipe.Title, recipe.Ingredients, recipe.Instructions, id,
	)
	return err
}

func (s *Store) DeleteRecipe(id int) error {
	_, err := s.DB.Exec("DELETE FROM recipes WHERE id = $1", id)
	return err
}

// User-related methods
func (s *Store) IsEmailExists(email string) (bool, error) {
	var exists bool
	err := s.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email,
	).Scan(&exists)
	return exists, err
}

func (s *Store) CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		user.Email, string(hashedPassword),
	)
	return err
}

func (s *Store) GetUserPasswordHash(email string) (string, error) {
	var storedHash string
	err := s.DB.QueryRow(
		"SELECT password_hash FROM users WHERE email = $1", email,
	).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return storedHash, nil
}
