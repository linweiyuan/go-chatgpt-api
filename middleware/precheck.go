package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
	"github.com/sirupsen/logrus"
	"net/http"
)

//goland:noinspection GoUnhandledErrorResult
func PreCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(api.Authorization) == "" {
			logrus.Info("Missing access token.")
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
			logrus.Info("Session timeout, need to refresh.")
			refreshDoneChannel := make(chan bool)

			go func() {
				webdriver.WebDriver.Refresh()

				webdriver.HandleCaptcha(webdriver.WebDriver)

				refreshDoneChannel <- true
			}()

			<-refreshDoneChannel
			logrus.Info("Refresh is done.")
		}

		c.Next()
	}
}
