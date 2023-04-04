package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
)

//goland:noinspection GoUnhandledErrorResult
func Refresh() {
	refreshDoneChannel := make(chan bool)

	go func() {
		WebDriver.Refresh()

		HandleCaptcha(WebDriver)

		refreshDoneChannel <- true
	}()

	<-refreshDoneChannel
	logger.Info("Refresh is done")
}

//goland:noinspection GoUnhandledErrorResult
func NewSessionAndRefresh(errorMessage string) {
	logger.Error("selenium error: " + errorMessage + ", need to create a new session and refresh")
	if _, err := WebDriver.PageSource(); err != nil {
		if err.Error() == "invalid session id: invalid session id" {
			logger.Info("old session id: " + WebDriver.SessionID())
			WebDriver.NewSession()
			logger.Info("new session id: " + WebDriver.SessionID())

			WebDriver.Get(api.ChatGPTUrl)
			HandleCaptcha(WebDriver)
		}
	}
}
