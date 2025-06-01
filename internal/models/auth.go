package models

// CreateUserRequest represents the request body for user registration
type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name,omitempty"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error string `json:"error"`
}