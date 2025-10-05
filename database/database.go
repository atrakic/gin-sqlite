package database

import (
	"os"
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Create table if not exists
	sqlStmt := `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA temp_store  = MEMORY;
	CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY AUTOINCREMENT unique,
		first_name TEXT not null,
		last_name TEXT not null,
		email TEXT not null unique);
	INSERT or IGNORE INTO people VALUES (1, 'Foo','Bar','foo@bar.com');
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
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
