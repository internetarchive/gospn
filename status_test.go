package spn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

// rewriteTransport redirects every request to the given test server,
// so code that hardcodes https://web.archive.org URLs can be tested locally.
type rewriteTransport struct {
	target *url.URL
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.target.Scheme
	req.URL.Host = t.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func newTestConnector(t *testing.T, handler http.HandlerFunc) Connector {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %s", err)
	}

	return Connector{
		AccessKey:  "AK",
		SecretKey:  "SK",
		HTTPClient: &http.Client{Transport: &rewriteTransport{target: target}},
	}
}

func TestUserStatusUpdate(t *testing.T) {
	to := UserStatus{}
	from := UserStatus{DailyCaptures: 1, DailyCapturesLimit: 2, Available: 3, Processing: 4}
	to.Update(from)
	if to != from {
		t.Errorf("Expected %+v, got %+v", from, to)
	}
}

func TestGetCaptureStatusSuccess(t *testing.T) {
	c := newTestConnector(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/save/status/JOB123" {
			t.Errorf("Expected path /save/status/JOB123, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "LOW AK:SK" {
			t.Errorf("Unexpected Authorization header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != userAgent {
			t.Errorf("Unexpected User-Agent header: %s", r.Header.Get("User-Agent"))
		}
		json.NewEncoder(w).Encode(CaptureStatus{
			JobID:     "JOB123",
			Status:    "success",
			Outlinks:  []string{"https://example.com"},
			Resources: []string{"https://example.com/style.css"},
		})
	})

	status, err := c.GetCaptureStatus("JOB123")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if status.JobID != "JOB123" || status.Status != "success" {
		t.Errorf("Unexpected status: %+v", status)
	}
	if len(status.Outlinks) != 1 || status.Outlinks[0] != "https://example.com" {
		t.Errorf("Unexpected outlinks: %+v", status.Outlinks)
	}
}

func TestGetCaptureStatusNilSlices(t *testing.T) {
	c := newTestConnector(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(CaptureStatus{JobID: "JOB123", Status: "pending"})
	})

	status, err := c.GetCaptureStatus("JOB123")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if status.Outlinks == nil {
		t.Errorf("Expected Outlinks to be defaulted, got nil")
	}
	if status.Resources == nil {
		t.Errorf("Expected Resources to be defaulted, got nil")
	}
}

func TestGetCaptureStatusNonOK(t *testing.T) {
	c := newTestConnector(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := c.GetCaptureStatus("JOB123")
	if err == nil {
		t.Fatal("Expected an error for non-200 response, got nil")
	}
}

func TestGetUserStatusSuccess(t *testing.T) {
	c := newTestConnector(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/save/status/user" {
			t.Errorf("Expected path /save/status/user, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(UserStatus{
			DailyCaptures:      10,
			DailyCapturesLimit: 100,
			Available:          5,
			Processing:         2,
		})
	})

	status, err := c.GetUserStatus()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	want := UserStatus{DailyCaptures: 10, DailyCapturesLimit: 100, Available: 5, Processing: 2}
	if status != want {
		t.Errorf("Expected %+v, got %+v", want, status)
	}
}

func TestGetUserStatusNonOK(t *testing.T) {
	c := newTestConnector(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})

	_, err := c.GetUserStatus()
	if err == nil {
		t.Fatal("Expected an error for non-200 response, got nil")
	}
}

func TestGetAvailableCaptureSlot(t *testing.T) {
	c := Connector{
		cachedStatus:   &UserStatus{Available: 1},
		cachedStatusMu: &sync.Mutex{},
	}

	if err := c.GetAvailableCaptureSlot(); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if c.cachedStatus.Available != 0 {
		t.Errorf("Expected Available to be 0, got %d", c.cachedStatus.Available)
	}
	if c.cachedStatus.Processing != 1 {
		t.Errorf("Expected Processing to be 1, got %d", c.cachedStatus.Processing)
	}
}
