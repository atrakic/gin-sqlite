// Package main provides the entry point for the Gin SQLite Demo API server.
//
//	@title			Gin SQLite Demo API
//	@version		1.0
//	@description	A simple REST API for managing persons using Gin and SQLite
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	https://github.com/atrakic/gin-sqlite-demo
//
//	@license.name	MIT
//	@license.url	https://github.com/atrakic/gin-sqlite-demo/blob/main/LICENSE
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/atrakic/gin-sqlite/internal/api"
	"github.com/atrakic/gin-sqlite/internal/auth"
	"github.com/atrakic/gin-sqlite/internal/database"
	"github.com/atrakic/gin-sqlite/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	if err := database.ConnectDatabase(); err != nil {
		log.Fatal(err)
	}

	// Initialize database tables
	if err := database.InitializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Println("Starting server...")
	r := setupRouter()

	// Auth endpoints
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", api.Login)
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("person", api.GetPersons)
		v1.GET("person/:id", api.GetPersonByID)

		// Needs JWT authentication
		v1.POST("person", jwtAuth, api.AddPerson)
		v1.PUT("person/:id", jwtAuth, api.UpdatePerson)
		v1.DELETE("person/:id", jwtAuth, api.DeletePerson)
	}

	_ = r.Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Serve only swagger.json file
	r.GET("/docs/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	// Swagger endpoint
	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/docs/swagger.json")))

	// PingHandler handles the ping endpoint
	// @Summary Health check endpoint
	// @Description Returns a pong message with timestamp
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} models.HealthCheckResponse "Pong response with timestamp"
	// @Router /ping [get]
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong " + fmt.Sprint(time.Now().Unix())})
	})
	return r
}

// jwtAuth validates JWT tokens from Authorization header
func jwtAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Error: "Authorization header required",
		})
		c.Abort()
		return
	}

	// Check for Bearer token format
	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Error: "Authorization header must be Bearer token",
		})
		c.Abort()
		return
	}

	// Validate JWT token
	claims, err := auth.ValidateJWT(tokenParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Error: "Invalid or expired token",
		})
		c.Abort()
		return
	}

	// Set user context for use in handlers
	c.Set("username", claims.Username)
	log.Printf("User authenticated: %s", claims.Username)
	c.Next()
}
