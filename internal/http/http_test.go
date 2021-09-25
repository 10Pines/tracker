package http

import (
	"github.com/10Pines/tracker/internal/logic"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoints(t *testing.T) {
	r := NewRouter(logic.Logic{}, "test")
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz/ready", nil)
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
