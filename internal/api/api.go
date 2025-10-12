package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/atrakic/gin-sqlite/internal/database"
	"github.com/gin-gonic/gin"
)

var (
	count int = 10
)

func GetPersons(c *gin.Context) {
	persons, err := database.DbGetPersons(count)
	checkErr(err)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": persons})
}

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

func AddPerson(c *gin.Context) {
	var json database.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DbAddPerson(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

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
		log.Fatal(err)
	}
}
