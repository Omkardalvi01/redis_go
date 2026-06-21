package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Omkardalvi01/redis_go.git/internal/engine"
)

func TestHandlerRejectsMissingCmd(t *testing.T) {
	eng, err := engine.New(engine.Config{EnableAOF: false})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	New(eng).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rr.Code)
	}
}
