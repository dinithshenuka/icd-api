package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	_ "modernc.org/sqlite"
)

func main() {
	// 1. Open the Excel file
	f, err := excelize.OpenFile("data/LinearizationMiniOutput-MMS-en.xlsx")
	if err != nil {
		log.Fatalf("Failed to open excel: %v", err)
	}
	defer f.Close()

	// 2. Setup SQLite Database in the database folder
	db, err := sql.Open("sqlite", "database/icd11.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 3. Clear existing data
	_, _ = db.Exec("DELETE FROM icd_codes")

	// 4. Read and Import
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO icd_codes(code, description, version) VALUES(?, ?, ?)")

	count := 0
	for i, row := range rows {
		if i < 1 { // Skip header
			continue
		}

		if len(row) < 5 {
			continue
		}

		code := strings.TrimSpace(row[2])
		description := strings.TrimSpace(row[4])
		
		if code == "" || description == "" {
			continue
		}

		// Clean description
		description = strings.TrimLeft(description, "- ")
		description = strings.TrimSpace(description)

		_, err = stmt.Exec(code, description, "ICD-11")
		if err != nil {
			log.Printf("Error inserting row %d: %v", i, err)
		}
		count++
	}

	tx.Commit()
	fmt.Printf("Data Import completed! Successfully imported %d codes into database/icd11.db.\n", count)
}
