package action

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "/helthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	var (
		rr = httptest.NewRecorder()
		handler = http.NewServeMux()
	)

	handler.HandleFunc("/healthcheck", HealthCheck)
	handler.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler return an unexpected HTTP status code: '%v' | expected: '%v'", status, http.StatusOK)
	}
}