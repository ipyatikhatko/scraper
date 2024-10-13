package main

import (
	"jobs-scraper/internal/scraper"
	"log"

	"github.com/joho/godotenv"
)

func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	// Initialize the database
	if err := scraper.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	
	// Start scraping
	if err := scraper.StartScraping(); err != nil {
		log.Fatalf("Error during scraping: %v", err)
	}
	
	// Close the database connection
	defer scraper.CloseDB() // Close the database connection when main exits
}