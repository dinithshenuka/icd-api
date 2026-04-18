package icdcodes

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// Repository handles database logic for ICD codes
type Repository struct {
	db *sql.DB
}

// NewRepository opens the SQLite database and returns a repository
func NewRepository(dbPath string) *Repository {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return &Repository{db: db}
}

// GetAll returns a limited set of codes (for sanity)
func (r *Repository) GetAll() []ICDCode {
	rows, err := r.db.Query("SELECT code, description, version FROM icd_codes LIMIT 100")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var codes []ICDCode
	for rows.Next() {
		var c ICDCode
		rows.Scan(&c.Code, &c.Description, &c.Version)
		codes = append(codes, c)
	}
	return codes
}

// Search filters codes using a SQL LIKE query for high-performance lookup
func (r *Repository) Search(query string) []ICDCode {
	searchTerm := "%" + query + "%"
	rows, err := r.db.Query("SELECT code, description, version FROM icd_codes WHERE code LIKE ? OR description LIKE ? LIMIT 50", searchTerm, searchTerm)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []ICDCode
	for rows.Next() {
		var c ICDCode
		rows.Scan(&c.Code, &c.Description, &c.Version)
		results = append(results, c)
	}
	return results
}
