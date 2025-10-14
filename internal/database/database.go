package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/atrakic/gin-sqlite/internal/models"
	_ "modernc.org/sqlite"
)

// DB is ...
var DB *sql.DB

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

// InitializeDatabase creates the people table if it doesn't exist
func InitializeDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	// Insert some sample data if table is empty
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM people").Scan(&count)
	if err != nil {
		log.Printf("Warning: Could not check table count: %v", err)
		return nil // Don't fail initialization if we can't check count
	}

	if count == 0 {
		log.Println("Initializing database with sample data...")
		// Use INSERT OR IGNORE to avoid duplicate email constraint errors
		sampleQueries := []string{
			"INSERT OR IGNORE INTO people (first_name, last_name, email) VALUES ('John', 'Doe', 'john.doe@example.com')",
			"INSERT OR IGNORE INTO people (first_name, last_name, email) VALUES ('Jane', 'Smith', 'jane.smith@example.com')",
			"INSERT OR IGNORE INTO people (first_name, last_name, email) VALUES ('Bob', 'Johnson', 'bob.johnson@example.com')",
		}

		for _, query := range sampleQueries {
			_, err := DB.Exec(query)
			if err != nil {
				log.Printf("Warning: Error adding sample data: %v", err)
			}
		}
		log.Println("Sample data initialization completed")
	} else {
		log.Printf("Database already contains %d records, skipping sample data initialization", count)
	}

	return nil
}

// DbGetPersonsCount returns the total count of persons in the database
func DbGetPersonsCount() (int64, error) {
	var count int64
	err := DB.QueryRow("SELECT COUNT(*) FROM people").Scan(&count)
	return count, err
}

// DbGetPersons retrieves persons with pagination support
func DbGetPersons(limit, offset int) ([]models.Person, error) {
	query := "SELECT id, first_name, last_name, email FROM people ORDER BY id LIMIT ? OFFSET ?"
	rows, err := DB.Query(query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]models.Person, 0)

	for rows.Next() {
		singlePerson := models.Person{}
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
func DbAddPerson(newPerson models.Person) (bool, error) {

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
func DbUpdatePerson(ourPerson models.Person, id int) (bool, error) {

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
func DbGetPersonByID(id string) (models.Person, error) {
	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email from people WHERE id = ?")

	if err != nil {
		return models.Person{}, err
	}

	person := models.Person{}
	sqlErr := stmt.QueryRow(id).Scan(&person.ID, &person.FirstName, &person.LastName, &person.Email)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Person{}, nil
		}
		return models.Person{}, sqlErr
	}
	return person, nil
}
