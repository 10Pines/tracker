package http

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiKeyRequired(t *testing.T) {
	apiKey := "123"

	tests := []struct {
		name           string
		headers        map[string]string
		expectedStatus int
	}{
		{
			name:           "without header",
			headers:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "incorrect header",
			headers:        map[string]string{"key": "value"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "incorrect value",
			headers:        map[string]string{apiKeyHeader: "test"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "correct header",
			headers:        map[string]string{apiKeyHeader: apiKey},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			context, router := gin.CreateTestContext(resp)
			router.Use(apiKeyRequired(apiKey))
			router.GET("/")

			request, _ := http.NewRequest(http.MethodGet, "/", nil)
			for k, v := range test.headers {
				request.Header.Set(k, v)
			}

			context.Request = request
			router.ServeHTTP(resp, context.Request)

			assert.Equal(t, test.expectedStatus, resp.Code)
		})
	}
}
