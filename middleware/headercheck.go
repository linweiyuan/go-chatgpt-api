package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
)

//goland:noinspection GoUnhandledErrorResult
func HeaderCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.AuthorizationHeader) == "" && c.Request.URL.Path != "/auth/login" {
			logger.Info("Missing access token")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage("Missing accessToken."))
			return
		}

		c.Next()
	}
}
