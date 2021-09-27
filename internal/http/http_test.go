package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/10Pines/tracker/internal/logic"
)

func TestHealthEndpoints(t *testing.T) {
	r := NewRouter(logic.Logic{}, "test")
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz/ready", nil)
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
