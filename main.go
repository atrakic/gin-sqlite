package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
       _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db *gorm.DB
	err error
)

type Person struct {
	Id          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

func getPersons(c *gin.Context) {
	var people []Person
 	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.IndentedJSON(http.StatusOK, people)
	}
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
	db, err := gorm.Open("sqlite3", "./person.db")
	if err != nil {
		 fmt.Println(err)
 	}
        defer db.Close()
        db.AutoMigrate(&Person{})
        p1 := Person{Id: 1, FirstName: "Foo", LastName: "Bar", Email: "foo@bar.com"}
	db.Create(&p1)
	
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
