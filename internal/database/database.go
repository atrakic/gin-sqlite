package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

// DB is ...
var DB *sql.DB

// Person is ...
type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// ConnectDatabase is ...
func ConnectDatabase() error {
	dataSourceName := os.Getenv("DATABASE_FILE")
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// DbGetPersons is ...
func DbGetPersons(count int) ([]Person, error) {

	rows, err := DB.Query("SELECT id, first_name, last_name, email from people LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.ID, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email)

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

// DbAddPerson is ...
func DbAddPerson(newPerson Person) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email)

	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return true, nil
}

// DbDeletePerson is ...
func DbDeletePerson(personID int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from people where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(personID)

	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return true, nil
}

// DbUpdatePerson is ...
func DbUpdatePerson(ourPerson Person, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, id)

	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return true, nil
}

// DbGetPersonByID is ...
func DbGetPersonByID(id string) (Person, error) {
	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email from people WHERE id = ?")

	if err != nil {
		return Person{}, err
	}

	person := Person{}
	sqlErr := stmt.QueryRow(id).Scan(&person.ID, &person.FirstName, &person.LastName, &person.Email)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}
	return person, nil
}
