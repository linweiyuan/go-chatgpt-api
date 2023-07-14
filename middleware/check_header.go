package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
)

//goland:noinspection SpellCheckingInspection
func CheckHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Authorization") == "" &&
		    c.GetHeader(api.AuthorizationHeader) == "" &&
			!strings.HasSuffix(c.Request.URL.Path, "/login") &&
			c.Request.URL.Path != "/" &&
			!strings.HasPrefix(c.Request.URL.Path, "/chatgpt/public-api") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage("Please provide a valid access token or api key in 'Authorization' header."))
			return
		}

		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
