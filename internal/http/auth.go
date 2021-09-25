package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const ApiKeyHeader = "X-API-KEY"

func apiKeyRequired(key string) gin.HandlerFunc {
	return func(g *gin.Context) {
		if g.GetHeader(ApiKeyHeader) != key {
			g.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
