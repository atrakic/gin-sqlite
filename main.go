package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	count int = 10
)

func getPersons(c *gin.Context) {
	persons, err := DbGetPersons(count)
	checkErr(err)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": persons})
}

func getPersonByID(c *gin.Context) {
	id := c.Param("id")

	person, err := DbGetPersonByID(id)
	checkErr(err)

	if person.FirstName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": person})
}

func addPerson(c *gin.Context) {

	var json Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := DbAddPerson(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

}

func updatePerson(c *gin.Context) {

	personID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	var json Person
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := DbUpdatePerson(json, personID)
	if success {
		fmt.Printf("Updating id %d", personID)
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func deletePerson(c *gin.Context) {
	personID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	// Must auth
	//user := c.MustGet(gin.AuthUserKey).(string)
	//if secret, ok := secrets[user]; ok {
	
	success, err := DbDeletePerson(personID)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "id #" + strconv.Itoa(personID) + " deleted"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return r
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == "admin" && password == "secret" {
		log.WithFields(log.Fields{
			"user": user,
		}).Info("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		return
	}
}

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

		// Enable auth from here:
		// curl -i -X "DELETE" http://admin:secret@localhost:8080/api/v1/person/2
		v1.DELETE("person/:id", basicAuth, deletePerson)
	}

	/*
	v1.Use(gin.BasicAuth(gin.Accounts {
		"admin": "secret",
	}))
	*/

	_ = r.Run()
}