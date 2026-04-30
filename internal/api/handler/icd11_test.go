package handler_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/dinithshenuka/icd-code-api/internal/api/handler"
	"github.com/dinithshenuka/icd-code-api/internal/domain/icd11"
	_ "modernc.org/sqlite"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestHandler builds a real ICD11Handler backed by an in-memory SQLite db.
func newTestHandler(t *testing.T, codes []icd11.ICDCode) *handler.ICD11Handler {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
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
	for _, c := range codes {
		_, err = db.Exec(
			"INSERT INTO icd_codes (code, description, version) VALUES (?, ?, ?)",
			c.Code, c.Description, c.Version,
		)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	t.Cleanup(func() { db.Close() })

	repo := icd11.NewRepositoryWithDB(db)
	svc := icd11.NewService(repo)
	return handler.NewICD11Handler(svc)
}

// --- GetIcd10 ---

func TestGetIcd10_ReturnsOK(t *testing.T) {
	h := newTestHandler(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)

	h.GetIcd10(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetIcd10_ReturnsJSON(t *testing.T) {
	want := icd11.ICDCode{Code: "A00", Description: "Cholera", Version: "ICD-10"}
	h := newTestHandler(t, []icd11.ICDCode{want})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)

	h.GetIcd10(c)

	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 1 || codes[0].Code != want.Code {
		t.Errorf("got %+v, want %+v", codes, []icd11.ICDCode{want})
	}
}

func TestGetIcd10_EmptyDB_ReturnsEmptyJSON(t *testing.T) {
	h := newTestHandler(t, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)

	h.GetIcd10(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	// Body should be valid JSON (null or [])
	var v interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}
}

// --- SearchCodes ---

func TestSearchCodes_NilQuery_Returns400(t *testing.T) {
	h := newTestHandler(t, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/search", nil)

	h.SearchCodes(c, handler.SearchCodesParams{Q: nil})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSearchCodes_EmptyQuery_Returns400(t *testing.T) {
	h := newTestHandler(t, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/search?q=", nil)

	empty := ""
	h.SearchCodes(c, handler.SearchCodesParams{Q: &empty})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSearchCodes_NilQuery_ErrorBody(t *testing.T) {
	h := newTestHandler(t, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/search", nil)

	h.SearchCodes(c, handler.SearchCodesParams{Q: nil})

	var resp handler.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON in error response: %v", err)
	}
	if resp.Code == "" || resp.Message == "" {
		t.Errorf("error response missing code/message: %+v", resp)
	}
}

func TestSearchCodes_ValidQuery_Returns200(t *testing.T) {
	h := newTestHandler(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/search?q=Cholera", nil)

	q := "Cholera"
	h.SearchCodes(c, handler.SearchCodesParams{Q: &q})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestSearchCodes_ValidQuery_ReturnsMatchingCodes(t *testing.T) {
	h := newTestHandler(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := "Cholera"
	h.SearchCodes(c, handler.SearchCodesParams{Q: &q})

	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 1 || codes[0].Code != "A00" {
		t.Errorf("unexpected results: %+v", codes)
	}
}

func TestSearchCodes_ValidQuery_NoMatch_ReturnsEmptyJSON(t *testing.T) {
	h := newTestHandler(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := "ZZZNOTFOUND"
	h.SearchCodes(c, handler.SearchCodesParams{Q: &q})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var v interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Errorf("response body is not valid JSON: %v", err)
	}
}

func TestSearchCodes_SQLInjection_Returns200(t *testing.T) {
	h := newTestHandler(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	q := "' OR '1'='1"
	h.SearchCodes(c, handler.SearchCodesParams{Q: &q})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("SQL injection returned rows, expected 0 got %d", len(codes))
	}
}

func TestSearchCodes_LargeQuery_Returns200(t *testing.T) {
	h := newTestHandler(t, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 10 000-character query string
	q := fmt.Sprintf("%010000d", 0)
	h.SearchCodes(c, handler.SearchCodesParams{Q: &q})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
