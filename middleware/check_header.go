package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yicone/go-chatgpt-api/api"
)

//goland:noinspection GoUnhandledErrorResult
func CheckHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.AuthorizationHeader) == "" && c.Request.URL.Path != "/auth/login" {
			c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage("Missing accessToken."))
			return
		}

		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
