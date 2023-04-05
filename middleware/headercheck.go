package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"net/http"
)

//goland:noinspection GoUnhandledErrorResult
func HeaderCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.Authorization) == "" {
			logger.Info("Missing access token")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage("Missing accessToken."))
			return
		}

		c.Next()
	}
}
