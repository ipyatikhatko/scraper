package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the connection to PostgreSQL and executes the setup SQL script
func InitDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
		return err // Return the error
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
		return err // Return the error
	}

	fmt.Println("Connected to the database")

	// Execute the SQL setup file
	if err := executeSQLFile("cmd/scraper/setup.sql"); err != nil {
		return err // Return the error
	}

	fmt.Println("Database setup is complete")
	return nil // Return nil if everything went well
}

// executeSQLFile reads an SQL file and executes its contents
func executeSQLFile(filename string) error {
	sqlFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Execute the SQL commands
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error executing SQL file: %v", err)
	}

	return nil
}

func insertJob(job Job) error {
	// Check if job already exists in the database
	if jobExists(job.Title, job.Company) {
		fmt.Printf("Job already exists in the database: %s at %s\n", job.Title, job.Company)
	} else {
		// Proceed with inserting the job if no duplicates are found
		insertQuery := `INSERT INTO jobs (title, company, work_format, location, company_type, experience_years, english_level, technology) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := db.Exec(insertQuery, job.Title, job.Company, job.WorkFormat, job.Location, job.CompanyType, job.Experience, job.EnglishLevel, job.Technology)
		if err != nil {
			return fmt.Errorf("error inserting job: %v", err)
		}
	}
	return nil
}

func jobExists(title, company string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM jobs WHERE title = $1 AND company = $2)`
	err := db.QueryRow(query, title, company).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking if job exists: %v\n", err)
		return false
	}
	return exists
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}
