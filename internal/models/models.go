// Package models defines the API models and DTOs for Swagger documentation
package models

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
