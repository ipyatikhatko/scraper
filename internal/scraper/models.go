package scraper

// Job represents a job listing with its relevant fields.
type Job struct {
	Title        string
	Company      string
	WorkFormat   string
	Location     string
	CompanyType  string
	Experience   string
	EnglishLevel string
	Technology   int
}

// Known company types for matching
var companyTypes = []string{"Product", "Outsource", "Outstaff", "Agency"}

// Known English levels for matching
var englishLevels = []string{"No English", "Beginner/Elementary", "Pre-Intermediate", "Intermediate", "Upper-Intermediate", "Advanced/Fluent"}

// Known work formats for matching
var workFormats = []string{"Remote", "Part-time", "Office"} 
