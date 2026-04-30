package icd11_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinithshenuka/icd-code-api/internal/domain/icd11"
	_ "modernc.org/sqlite"
)

// newTestDB creates an in-memory SQLite database with the icd_codes schema.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE icd_codes (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		code        TEXT NOT NULL,
		description TEXT NOT NULL,
		version     TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// seedCodes inserts rows into the in-memory database.
func seedCodes(t *testing.T, db *sql.DB, codes []icd11.ICDCode) {
	t.Helper()
	for _, c := range codes {
		_, err := db.Exec(
			"INSERT INTO icd_codes (code, description, version) VALUES (?, ?, ?)",
			c.Code, c.Description, c.Version,
		)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
}

// --- GetAll ---

func TestRepository_GetAll_EmptyDB(t *testing.T) {
	db := newTestDB(t)
	repo := icd11.NewRepositoryWithDB(db)

	codes := repo.GetAll()
	if codes != nil {
		t.Errorf("expected nil slice for empty DB, got %v", codes)
	}
}

func TestRepository_GetAll_ReturnsCodes(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	codes := repo.GetAll()
	if len(codes) != 2 {
		t.Errorf("expected 2 codes, got %d", len(codes))
	}
}

func TestRepository_GetAll_LimitedTo100(t *testing.T) {
	db := newTestDB(t)
	// Insert 150 rows.
	batch := make([]icd11.ICDCode, 150)
	for i := range batch {
		batch[i] = icd11.ICDCode{
			Code:        fmt.Sprintf("X%03d", i),
			Description: fmt.Sprintf("Disease %d", i),
			Version:     "ICD-10",
		}
	}
	seedCodes(t, db, batch)
	repo := icd11.NewRepositoryWithDB(db)

	codes := repo.GetAll()
	if len(codes) != 100 {
		t.Errorf("expected 100 codes (LIMIT), got %d", len(codes))
	}
}

func TestRepository_GetAll_FieldsPopulated(t *testing.T) {
	db := newTestDB(t)
	want := icd11.ICDCode{Code: "A00", Description: "Cholera", Version: "ICD-10"}
	seedCodes(t, db, []icd11.ICDCode{want})
	repo := icd11.NewRepositoryWithDB(db)

	codes := repo.GetAll()
	if len(codes) == 0 {
		t.Fatal("expected at least one code")
	}
	got := codes[0]
	if got.Code != want.Code || got.Description != want.Description || got.Version != want.Version {
		t.Errorf("field mismatch: got %+v, want %+v", got, want)
	}
}

// --- Search ---

func TestRepository_Search_EmptyDB(t *testing.T) {
	db := newTestDB(t)
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("cholera")
	if results != nil {
		t.Errorf("expected nil for empty DB, got %v", results)
	}
}

func TestRepository_Search_MatchesByCode(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("A00")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Code != "A00" {
		t.Errorf("expected code A00, got %s", results[0].Code)
	}
}

func TestRepository_Search_MatchesByDescription(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("Varicella")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Code != "B01" {
		t.Errorf("expected code B01, got %s", results[0].Code)
	}
}

func TestRepository_Search_CaseInsensitive(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	// SQLite LIKE is case-insensitive for ASCII by default.
	results := repo.Search("cholera")
	if len(results) != 1 {
		t.Errorf("expected case-insensitive match, got %d results", len(results))
	}
}

func TestRepository_Search_NoMatch(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("ZZZNOTFOUND")
	if results != nil {
		t.Errorf("expected nil for no match, got %v", results)
	}
}

func TestRepository_Search_LimitedTo50(t *testing.T) {
	db := newTestDB(t)
	batch := make([]icd11.ICDCode, 80)
	for i := range batch {
		batch[i] = icd11.ICDCode{
			Code:        fmt.Sprintf("A%02d", i),
			Description: "Matching disease",
			Version:     "ICD-10",
		}
	}
	seedCodes(t, db, batch)
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("Matching")
	if len(results) != 50 {
		t.Errorf("expected 50 results (LIMIT), got %d", len(results))
	}
}

func TestRepository_Search_PartialMatch(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00.1", Description: "Cholera due to Vibrio", Version: "ICD-10"},
		{Code: "A00.9", Description: "Cholera, unspecified", Version: "ICD-10"},
		{Code: "B01",   Description: "Varicella",            Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	results := repo.Search("Cholera")
	if len(results) != 2 {
		t.Errorf("expected 2 partial matches, got %d", len(results))
	}
}

func TestRepository_Search_SQLInjectionSafe(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)

	// This should not cause a SQL error or return unexpected rows.
	results := repo.Search("' OR '1'='1")
	// Should return nil/empty, not all rows.
	if len(results) != 0 {
		t.Errorf("SQL injection attempt returned %d rows, expected 0", len(results))
	}
}

// --- Ping ---

func TestRepository_Ping(t *testing.T) {
	db := newTestDB(t)
	repo := icd11.NewRepositoryWithDB(db)

	if err := repo.Ping(); err != nil {
		t.Errorf("Ping failed on open database: %v", err)
	}
}
