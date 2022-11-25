package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string   `json:"email"`
}

// test data
var persons = []Person{
	{ID: 1, FirstName: "Foo", LastName: "Bar", Email: "foo@bar.com"},
}

func getPersons(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, persons)
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "getPersonById " + id +" Called"})
}

func addPerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "addPerson Called"})
}

func updatePerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "deletePerson " + id + " Called"})
}

func main() {
	r := gin.Default()	
	v1 := r.Group("/api/v1")
	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonById)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)
		v1.DELETE("person/:id", deletePerson)
	}
	r.Run()
}
