package main

import (
	"fmt"
	"time"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//p1 := Person{Id: 1, FirstName: "Foo", LastName: "Bar", Email: "foo@bar.com"}
	if err := ConnectDatabase(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server...")
	r := setupRouter()

	v1 := r.Group("/api/v1")
	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonByID)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)

		// Basic Auth from here:
		// curl -i -X "DELETE" http://admin:secret@localhost:8080/api/v1/person/2
		v1.DELETE("person/:id", basicAuth, deletePerson)
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
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == "admin" && password == "secret" {
		log.Println("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
}