package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Omkardalvi01/redis_go.git/internal/engine"
)

type Handler struct {
	Engine *engine.Engine
}

func New(eng *engine.Engine) http.Handler {
	return &Handler{Engine: eng}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmd := strings.TrimSpace(r.URL.Query().Get("cmd"))
	if cmd == "" {
		http.Error(w, "(error) missing cmd query parameter", http.StatusBadRequest)
		return
	}

	result := h.Engine.ExecuteLine(cmd, true)
	if result.Err != nil {
		status := http.StatusBadRequest
		if strings.Contains(strings.ToLower(result.Err.Error()), "replay") {
			status = http.StatusInternalServerError
		}
		http.Error(w, result.Err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, result.Response)
}
