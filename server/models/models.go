package models

// The Recipe struct holds metadata for a recipe
type Recipe struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

// The Email struct holds name and email address
type Email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email,omitempty"`
}

// The User Struct for User authentication and JWT authorization
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
