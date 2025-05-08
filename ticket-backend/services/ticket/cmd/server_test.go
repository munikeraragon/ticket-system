package ticket

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestListTicketsIntegration(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Clear and seed
	db.Exec("DELETE FROM tickets")
	db.Exec(`INSERT INTO tickets (id, customer_name, email, created_at, status, notes)
	         VALUES (1, 'Test User', 'test@example.com', '2024-01-01T00:00:00Z', 'open', 'note')`)

	server := NewTicketServer(db)

	go func() {
		StartHTTPServer(server, ":8082")
	}()
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}

	// GET /tickets
	resp, err := http.Get("http://localhost:8082/tickets")
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	t.Logf("GET /tickets Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// GET /tickets/1
	resp, err = http.Get("http://localhost:8082/tickets/1")
	if err != nil {
		t.Fatalf("GET by ID failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	t.Logf("GET /tickets/1 Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for GET /tickets/1, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// PATCH /tickets/1 (valid)
	reqBody := `{"status": "done", "notes": "updated notes"}`
	req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8082/tickets/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("PATCH request failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	t.Logf("PATCH /tickets/1 Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for PATCH, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// GET /tickets/9999 (not found)
	resp, err = http.Get("http://localhost:8082/tickets/9999")
	if err != nil {
		t.Fatalf("GET by ID (not found) failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	t.Logf("GET /tickets/9999 Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for GET /tickets/9999, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// PATCH /tickets/9999 (not found)
	reqBody = `{"status": "done", "notes": "not found update"}`
	req, _ = http.NewRequest(http.MethodPatch, "http://localhost:8082/tickets/9999", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("PATCH (not found) request failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	t.Logf("PATCH /tickets/9999 Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for PATCH /tickets/9999, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// PATCH /tickets/1 (invalid status)
	reqBody = `{"status": "INVALID", "notes": "bad status"}`
	req, _ = http.NewRequest(http.MethodPatch, "http://localhost:8082/tickets/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("PATCH (invalid status) failed: %v", err)
	}
	body, _ = io.ReadAll(resp.Body)
	t.Logf("PATCH /tickets/1 (invalid status) Status: %d", resp.StatusCode)
	t.Logf("Body: %s", string(body))
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for PATCH /tickets/1 with invalid status, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
