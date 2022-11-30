package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

	_, err := DbAddPerson(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
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

	if _, err := DbUpdatePerson(json, personID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	fmt.Printf("Updating id %d", personID)
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func deletePerson(c *gin.Context) {
	personID, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	if _, err := DbDeletePerson(personID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "id #" + strconv.Itoa(personID) + " deleted"})
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
