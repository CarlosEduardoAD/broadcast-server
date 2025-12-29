package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck_Returns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	Healthcheck(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestHealthcheck_AllMethodsReturn200(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut}
	for _, m := range methods {
		t.Run(m, func(t *testing.T) {
			req := httptest.NewRequest(m, "/health", nil)
			rr := httptest.NewRecorder()

			Healthcheck(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Fatalf("method %s: expected status %d, got %d", m, http.StatusOK, res.StatusCode)
			}
		})
	}
}
