// Package models defines the API models and DTOs for Swagger documentation
package models

import "github.com/golang-jwt/jwt/v5"

// Person represents a person in the database
// @Description Person information
// @ID Person
type Person struct {
	ID        uint64 `json:"id" example:"1" format:"uint64"`                      // Person ID
	FirstName string `json:"first_name" example:"John" maxLength:"50"`            // First name
	LastName  string `json:"last_name" example:"Doe" maxLength:"50"`              // Last name
	Email     string `json:"email" example:"john.doe@example.com" format:"email"` // Email address
} // @name Person

// CreatePersonRequest represents the request body for creating a person
// @Description Request body for creating a new person
type CreatePersonRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"John" maxLength:"50"`                  // First name (required)
	LastName  string `json:"last_name" binding:"required" example:"Doe" maxLength:"50"`                    // Last name (required)
	Email     string `json:"email" binding:"required,email" example:"john.doe@example.com" format:"email"` // Email address (required)
} // @name CreatePersonRequest

// UpdatePersonRequest represents the request body for updating a person
// @Description Request body for updating an existing person
type UpdatePersonRequest struct {
	FirstName string `json:"first_name,omitempty" example:"Jane" maxLength:"50"`              // First name (optional)
	LastName  string `json:"last_name,omitempty" example:"Smith" maxLength:"50"`              // Last name (optional)
	Email     string `json:"email,omitempty" example:"jane.smith@example.com" format:"email"` // Email address (optional)
} // @name UpdatePersonRequest

// APIResponse represents a generic API response
// @Description Generic API response
type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`    // Response data
	Message string      `json:"message,omitempty"` // Response message
	Error   string      `json:"error,omitempty"`   // Error message
} // @name APIResponse

// HealthCheckResponse represents the health check response
// @Description Health check response
type HealthCheckResponse struct {
	Message string `json:"message" example:"pong 1697123456"` // Health check message with timestamp
} // @name HealthCheckResponse

// LoginRequest represents the login request body
// @Description Login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`  // Username (required)
	Password string `json:"password" binding:"required" example:"secret"` // Password (required)
} // @name LoginRequest

// LoginResponse represents the login response body
// @Description Login response body
type LoginResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT token
	ExpiresAt int64  `json:"expires_at" example:"1697209856"`                         // Token expiration timestamp
} // @name LoginResponse

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
} // @name JWTClaims

// PaginationRequest represents pagination query parameters
// @Description Pagination request parameters
type PaginationRequest struct {
	Page     int `form:"page" json:"page" example:"1" minimum:"1"`                          // Page number (starting from 1)
	PageSize int `form:"page_size" json:"page_size" example:"10" minimum:"1" maximum:"100"` // Number of items per page
} // @name PaginationRequest

// PaginationMeta represents pagination metadata
// @Description Pagination metadata information
type PaginationMeta struct {
	CurrentPage int   `json:"current_page" example:"1"`      // Current page number
	PageSize    int   `json:"page_size" example:"10"`        // Number of items per page
	TotalPages  int   `json:"total_pages" example:"5"`       // Total number of pages
	TotalItems  int64 `json:"total_items" example:"50"`      // Total number of items
	HasNextPage bool  `json:"has_next_page" example:"true"`  // Whether there is a next page
	HasPrevPage bool  `json:"has_prev_page" example:"false"` // Whether there is a previous page
} // @name PaginationMeta

// PaginatedResponse represents a paginated API response
// @Description Paginated API response
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`       // Response data
	Pagination PaginationMeta `json:"pagination"` // Pagination metadata
} // @name PaginatedResponse
