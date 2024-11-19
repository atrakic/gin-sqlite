package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atrakic/gin-sqlite/api"
	"github.com/atrakic/gin-sqlite/database"
	"github.com/gin-gonic/gin"
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
