package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
)

//goland:noinspection SpellCheckingInspection
func CheckHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.AuthorizationHeader) == "" &&
			c.Request.URL.Path != "/chatgpt/login" &&
			c.Request.URL.Path != "/platform/login" &&
			c.Request.URL.Path != "/healthCheck" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage("Hello World."))
			return
		}

		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
