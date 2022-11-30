package main

import (
	"log"
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
