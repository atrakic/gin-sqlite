package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/atrakic/gin-sqlite/internal/auth"
	"github.com/atrakic/gin-sqlite/internal/database"
	"github.com/atrakic/gin-sqlite/internal/models"
	"github.com/gin-gonic/gin"
)

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate user credentials and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 401 {object} models.APIResponse "Invalid credentials"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: "Invalid request format",
		})
		return
	}

	// Validate credentials
	if !auth.ValidateCredentials(loginRequest.Username, loginRequest.Password) {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Error: "Invalid username or password",
		})
		return
	}

	// Generate JWT token
	token, expiresAt, err := auth.GenerateJWT(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}

// GetPersons retrieves persons from the database with pagination
// @Summary Get all persons with pagination
// @Description Get a paginated list of all persons in the database
// @Tags persons
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1) example(1)
// @Param page_size query int false "Number of items per page (default: 10, max: 100)" minimum(1) maximum(100) example(10)
// @Success 200 {object} models.PaginatedResponse "Paginated list of persons"
// @Failure 400 {object} models.APIResponse "Invalid pagination parameters"
// @Failure 500 {object} models.APIResponse "Internal server error"
// @Router /api/v1/person [get]
func GetPersons(c *gin.Context) {
	// Parse pagination parameters
	var pagination models.PaginationRequest

	// Set defaults
	pagination.Page = 1
	pagination.PageSize = 10

	// Bind query parameters
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Error: "Invalid pagination parameters: " + err.Error(),
		})
		return
	}

	// Validate pagination parameters
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 {
		pagination.PageSize = 10
	}
	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.PageSize

	// Get total count
	totalCount, err := database.DbGetPersonsCount()
	if err != nil {
		log.Printf("Database error getting count: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: "Failed to get total count",
		})
		return
	}

	// Get persons with pagination
	persons, err := database.DbGetPersons(pagination.PageSize, offset)
	if err != nil {
		log.Printf("Database error getting persons: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Error: "Failed to retrieve persons",
		})
		return
	}

	// Calculate pagination metadata
	totalPages := int((totalCount + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	paginationMeta := models.PaginationMeta{
		CurrentPage: pagination.Page,
		PageSize:    pagination.PageSize,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
		HasNextPage: pagination.Page < totalPages,
		HasPrevPage: pagination.Page > 1,
	}

	// Return paginated response
	response := models.PaginatedResponse{
		Data:       persons,
		Pagination: paginationMeta,
	}

	c.JSON(http.StatusOK, response)
}

// GetPersonByID retrieves a person by their ID
// @Summary Get person by ID
// @Description Get a single person by their ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} database.Person "Person details"
// @Failure 404 {object} models.APIResponse "Person not found"
// @Router /api/v1/person/{id} [get]
func GetPersonByID(c *gin.Context) {
	id := c.Param("id")

	person, err := database.DbGetPersonByID(id)
	checkErr(err)

	if person.FirstName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": person})
}

// AddPerson creates a new person
// @Summary Create a new person
// @Description Create a new person with the provided information
// @Tags persons
// @Accept json
// @Produce json
// @Param person body models.CreatePersonRequest true "Person to create"
// @Success 200 {object} models.APIResponse "Success message"
// @Failure 400 {object} models.APIResponse "Invalid input"
// @Security BearerAuth
// @Router /api/v1/person [post]
func AddPerson(c *gin.Context) {
	var json database.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DbAddPerson(json)
	checkErr(err)

	c.JSON(http.StatusOK, gin.H{"message": "Person added successfully"})
}

// UpdatePerson updates an existing person
// @Summary Update a person
// @Description Update an existing person by ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body models.UpdatePersonRequest true "Person to update"
// @Success 200 {object} models.APIResponse "Success message"
// @Failure 400 {object} models.APIResponse "Invalid input or ID"
// @Security BearerAuth
// @Router /api/v1/person/{id} [put]
func UpdatePerson(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	var json database.Person
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := database.DbUpdatePerson(json, personID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	fmt.Printf("Updating id %d", personID)
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// DeletePerson deletes a person by ID
// @Summary Delete a person
// @Description Delete a person by ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} models.APIResponse "Success message"
// @Failure 400 {object} models.APIResponse "Invalid ID"
// @Failure 500 {object} models.APIResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/person/{id} [delete]
func DeletePerson(c *gin.Context) {
	personID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	if _, err := database.DbDeletePerson(personID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "id #" + strconv.Itoa(personID) + " deleted"})
}

// checkErr is ...
func checkErr(err error) {
	if err != nil {
		log.Printf("Database error: %v", err)
	}
}
