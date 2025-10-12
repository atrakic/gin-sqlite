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
//	@securityDefinitions.basic	BasicAuth
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atrakic/gin-sqlite/internal/api"
	"github.com/atrakic/gin-sqlite/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	if err := database.ConnectDatabase(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server...")
	r := setupRouter()

	v1 := r.Group("/api/v1")
	{
		v1.GET("person", api.GetPersons)
		v1.GET("person/:id", api.GetPersonByID)

		// Needs authentication
		v1.POST("person", basicAuth, api.AddPerson)
		v1.PUT("person/:id", basicAuth, api.UpdatePerson)
		v1.DELETE("person/:id", basicAuth, api.DeletePerson)
	}

	_ = r.Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Serve swagger.json file
	r.Static("/docs", "./docs")

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
	// @Success 200 {object} map[string]interface{} "Pong response with timestamp"
	// @Router /ping [get]
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong " + fmt.Sprint(time.Now().Unix())})
	})
	return r
}

func basicAuth(c *gin.Context) {
	_admin := os.Getenv("ADMIN_USER")
	if _admin == "" {
		_admin = "admin"
	}
	_password := os.Getenv("ADMIN_PASSWORD")
	if _password == "" {
		_password = "secret"
	}

	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == _admin && password == _password {
		log.Println("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
}
