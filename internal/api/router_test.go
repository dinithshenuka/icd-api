package api_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/dinithshenuka/icd-code-api/internal/api"
	"github.com/dinithshenuka/icd-code-api/internal/api/handler"
	"github.com/dinithshenuka/icd-code-api/internal/domain/icd11"
	_ "modernc.org/sqlite"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestRouter builds a full *gin.Engine backed by an in-memory SQLite db.
func newTestRouter(t *testing.T, codes []icd11.ICDCode) *gin.Engine {
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
	h := handler.NewICD11Handler(svc)
	return api.NewRouter(h)
}

// --- /health ---

func TestHealth_Returns200(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHealth_ReturnsStatusUp(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "up" {
		t.Errorf("expected status=up, got %q", body["status"])
	}
}

// --- GET /v1/icd10 ---

func TestGetIcd10_Endpoint_Returns200(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetIcd10_Endpoint_ReturnsValidJSON(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)
	r.ServeHTTP(w, req)

	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 2 {
		t.Errorf("expected 2 codes, got %d", len(codes))
	}
}

func TestGetIcd10_Endpoint_EmptyDB_Returns200(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/icd10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- GET /v1/search ---

func TestSearch_Endpoint_MissingQ_Returns400(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSearch_Endpoint_EmptyQ_Returns400(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSearch_Endpoint_WithQ_Returns200(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=Cholera", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestSearch_Endpoint_WithQ_ReturnsMatches(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
		{Code: "B01", Description: "Varicella", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=A00", nil)
	r.ServeHTTP(w, req)

	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 1 || codes[0].Code != "A00" {
		t.Errorf("unexpected results: %+v", codes)
	}
}

func TestSearch_Endpoint_NoMatch_Returns200WithEmptyJSON(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=ZZZNOTFOUND", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var v interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Errorf("body is not valid JSON: %v", err)
	}
}

func TestSearch_Endpoint_SQLInjection_Returns200(t *testing.T) {
	r := newTestRouter(t, []icd11.ICDCode{
		{Code: "A00", Description: "Cholera", Version: "ICD-10"},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=%27+OR+%271%27%3D%271", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var codes []icd11.ICDCode
	if err := json.Unmarshal(w.Body.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("SQL injection returned %d rows, expected 0", len(codes))
	}
}

func TestSearch_Endpoint_ContentTypeIsJSON(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search?q=test", nil)
	r.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct == "" {
		t.Error("expected Content-Type header to be set")
	}
}

func TestSearch_Endpoint_ErrorBody_HasCodeAndMessage(t *testing.T) {
	r := newTestRouter(t, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/search", nil)
	r.ServeHTTP(w, req)

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["code"] == "" || resp["message"] == "" {
		t.Errorf("error response missing code/message: %+v", resp)
	}
}
