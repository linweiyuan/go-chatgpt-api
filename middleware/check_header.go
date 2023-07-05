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
		if c.GetHeader(api.AuthorizationHeader) == "" &&
			c.Request.URL.Path != "/chatgpt/login" &&
			c.Request.URL.Path != "/platform/login" &&
			c.Request.URL.Path != "/" &&
			!strings.HasPrefix(c.Request.URL.Path, "/chatgpt/public-api") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage("Please provide a valid access token or api key in 'Authorization' header."))
			return
		}

		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
