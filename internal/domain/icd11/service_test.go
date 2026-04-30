package icd11_test

import (
	"testing"

	"github.com/dinithshenuka/icd-code-api/internal/domain/icd11"
)

func TestService_GetCodes_Delegates(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)

	codes := svc.GetCodes()
	if len(codes) != 2 {
		t.Errorf("expected 2 codes, got %d", len(codes))
	}
}

func TestService_GetCodes_EmptyDB(t *testing.T) {
	db := newTestDB(t)
	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)

	codes := svc.GetCodes()
	if codes != nil {
		t.Errorf("expected nil for empty DB, got %v", codes)
	}
}

func TestService_SearchCodes_ReturnsMatches(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)

	results := svc.SearchCodes("Cholera")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Code != "A00" {
		t.Errorf("expected code A00, got %s", results[0].Code)
	}
}

func TestService_SearchCodes_NoMatch(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)

	results := svc.SearchCodes("ZZZNOTFOUND")
	if results != nil {
		t.Errorf("expected nil for no match, got %v", results)
	}
}

func TestService_SearchCodes_EmptyQuery(t *testing.T) {
	db := newTestDB(t)
	seedCodes(t, db, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})
	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)

	// Empty query matches everything via LIKE '%%'
	results := svc.SearchCodes("")
	if len(results) == 0 {
		t.Error("expected at least one result for empty query (matches all)")
	}
}
