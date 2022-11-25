package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	count int = 10
	DB    *sql.DB
)

type Person struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// utils
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetPersonById(id string) (Person, error) {
	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email from people WHERE id = ?")

	if err != nil {
		return Person{}, err
	}

	person := Person{}
	sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}
	return person, nil
}

func utilsGetPersons(count int) ([]Person, error) {

	rows, err := DB.Query("SELECT id, first_name, last_name, email from people LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email)

		if err != nil {
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return people, err
}

func getPersons(c *gin.Context) {
	persons, err := utilsGetPersons(count)
	checkErr(err)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": persons})
	}
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")

	person, err := GetPersonById(id)
	checkErr(err)
	// if the name is blank we can assume nothing is found
	if person.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": person})
	}
}

// TODO:
func addPerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "addPerson Called"})
}

func updatePerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	c.JSON(http.StatusOK, gin.H{"id #i" + id: "deleted"})
}

func setupRouter() *gin.Engine {
	//gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonById)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)
		v1.DELETE("person/:id", deletePerson)
	}
	return r
}

func connectDatabase() error {
	db, err := sql.Open("sqlite3", "./person.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func main() {
	//p1 := Person{Id: 1, FirstName: "Foo", LastName: "Bar", Email: "foo@bar.com"}
	connectDatabase()

	r := setupRouter()
	r.Run()
}
