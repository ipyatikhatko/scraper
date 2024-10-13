package scraper

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var technologiesMap map[string]int // To hold technology names and their corresponding IDs

// getRandomDelay returns a random duration between 5 and 15 seconds
func getRandomDelay() time.Duration {
	return time.Duration(rand.Intn(11)+5) * time.Second // Generates a number between 5 and 15
}

// FetchAllTechnologies retrieves all technologies from the database and stores them in a map
func FetchAllTechnologies() error {
	technologiesMap = make(map[string]int) // Initialize the map

	query := `SELECT id, djinni_keyword FROM technologies`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error fetching technologies: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var djinni_keyword string
		if err := rows.Scan(&id, &djinni_keyword); err != nil {
			return fmt.Errorf("error scanning technology row: %v", err)
		}
		technologiesMap[djinni_keyword] = id // Store in the map
	}
	return nil
}

func scrapeKeyword(c *colly.Collector, keyword string, techID int) {
	fmt.Printf("Scraping for keyword: %s\n", keyword)

	// Set up job scraping handlers, passing techID to correctly associate the technology
	setupJobScrapingHandlers(c, techID)

	// Visit the first page to find the total number of jobs and scrape jobs at the same time
	url := fmt.Sprintf("https://djinni.co/jobs/?primary_keyword=%s&region=UKR", keyword)
	totalPages := 0 // Initialize totalPages to 0

	// On visiting the first page, calculate total pages and process jobs
	c.OnHTML("header", func(e *colly.HTMLElement) {
		jobsCountText := e.ChildText("h1 > span")
		jobsCount, err := strconv.Atoi(strings.ReplaceAll(jobsCountText, ",", "")) // Remove commas and convert to int
		if err == nil {
			// Calculate total pages (15 jobs per page)
			totalPages = (jobsCount + 14) / 15 // Ceiling division
			fmt.Printf("Total jobs: %d, Pages to scrape: %d\n", jobsCount, totalPages)
		}
	})

	// Scrape the first page and process jobs
	c.Visit(url)

	// Wait for the first page to finish processing and calculate total jobs
	c.Wait()

	// If there are more pages, scrape the remaining pages
	if totalPages > 1 {
		scrapePages(c, keyword, techID, totalPages)
	}
}

// scrapePages scrapes all the remaining pages for a given keyword and technology ID
func scrapePages(c *colly.Collector, keyword string, techID int, totalPages int) {
	// Loop through pages 2 to totalPages for the current keyword
	for page := 2; page <= totalPages; page++ {
		// Formulate the URL for the current page
		url := fmt.Sprintf("https://djinni.co/jobs/?primary_keyword=%s&region=UKR&page=%d", keyword, page)

		// Visit the URL
		c.Visit(url)

		// Introduce a random delay before visiting the next URL
		if page < totalPages {
			delay := getRandomDelay()
			fmt.Printf("Waiting for %v seconds before visiting page %d for keyword: %s...\n", delay.Seconds(), page, keyword)
			time.Sleep(delay)
		}
	}

	// Wait for the scraping to complete before moving to the next keyword
	c.Wait()
}

func setupJobScrapingHandlers(c *colly.Collector, techID int) {
	// Set up event listeners only once
	c.OnHTML("main > ul", func(e *colly.HTMLElement) {
		e.ForEach("li[id^='job-item-']", func(i int, h *colly.HTMLElement) {
			// Extract job information here
			job := Job{
				Title:        h.ChildText("h3 > a.job-item__title-link"),
				Company:      h.ChildText("div.d-inline-flex > a.text-body"),
				WorkFormat:   "Not specified",
				Location:     "Not specified",
				CompanyType:  "Not specified",
				Experience:   "0",
				EnglishLevel: "Not specified",
				Technology:   techID,
			}

			// Loop through each .text-nowrap element to extract specific information
			h.ForEach(".fw-medium .text-nowrap", func(i int, elem *colly.HTMLElement) {
				text := strings.TrimSpace(elem.Text)

				// Use string matching to determine the job attributes
				switch {
				case containsAny(text, workFormats):
					job.WorkFormat = text
				case strings.Contains(text, "Worldwide"), strings.Contains(text, "Ukraine"):
					job.Location = text
				case containsAny(text, companyTypes):
					job.CompanyType = text
				case strings.Contains(text, "year"):
					experience := extractExperience(text)
					if experience != "" {
						job.Experience = experience
					}
				case containsAny(text, englishLevels):
					job.EnglishLevel = text
				}
			})

			err := insertJob(job)
			if err != nil {
				fmt.Printf("Error inserting job into database: %v\n", err)
			} else {
				fmt.Printf("Inserted job into database: %s at %s\n", job.Title, job.Company)
			}
		})
	})
}


func StartScraping() error {
    // Fetch all technologies before scraping
    err := FetchAllTechnologies()
    if err != nil {
        return fmt.Errorf("failed to fetch technologies: %v", err)
    }

    // Loop through each programming keyword sequentially
    for keyword, techID := range technologiesMap {
        // Create a new Colly collector for each keyword to avoid shared state issues
        c := colly.NewCollector()

        // Set up session ID and headers
        djiniSessionId := os.Getenv("DJINI_SESSION_ID")

        c.OnRequest(func(r *colly.Request) {
            r.Headers.Set("Cookie", fmt.Sprintf("sessionid=%s;", djiniSessionId))
            r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
            r.Headers.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
        })

        // Scrape the keyword and wait for it to finish before moving to the next
        scrapeKeyword(c, keyword, techID)

        // Wait for the scraping to complete before moving to the next keyword
        c.Wait()
    }

    return nil
}

