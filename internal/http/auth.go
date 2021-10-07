package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/10Pines/tracker/v2/pkg/tracker"
)

func apiKeyRequired(key string) gin.HandlerFunc {
	return func(g *gin.Context) {
		if g.GetHeader(tracker.APIKeyHeader) != key {
			g.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
