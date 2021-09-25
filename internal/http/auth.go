package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const apiKeyHeader = "X-API-KEY"

func apiKeyRequired(key string) gin.HandlerFunc {
	return func(g *gin.Context) {
		if g.GetHeader(apiKeyHeader) != key {
			g.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
