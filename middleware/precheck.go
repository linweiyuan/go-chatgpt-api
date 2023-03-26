package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
	"net/http"
)

//goland:noinspection GoUnhandledErrorResult
func PreCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.Authorization) == "" {
			logger.Info("Missing access token")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage("Missing accessToken."))
			return
		}

		url := api.PreCheckUrl
		xhrStatus, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
			const xhr = new XMLHttpRequest();
			xhr.open('GET', '%s', false);
			xhr.send();
			return xhr.status;`, url), nil)

		if xhrStatus == float64(http.StatusForbidden) {
			logger.Warn("Session timeout, need to refresh")
			refreshDoneChannel := make(chan bool)

			go func() {
				webdriver.WebDriver.Refresh()

				webdriver.HandleCaptcha(webdriver.WebDriver)

				refreshDoneChannel <- true
			}()

			<-refreshDoneChannel
			logger.Info("Refresh is done")
		}

		c.Next()
	}
}
