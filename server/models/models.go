package models

// Struct for recipes
type Recipe struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

// Email holds email name and address info
type Email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email,omitempty"`
}

// Struct for User authentication
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
