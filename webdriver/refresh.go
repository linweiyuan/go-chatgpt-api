package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/api"
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
}

//goland:noinspection GoUnhandledErrorResult
func NewSessionAndRefresh() {
	if _, err := WebDriver.PageSource(); err != nil {
		if err.Error() == "invalid session id: invalid session id" {
			WebDriver.NewSession()
			WebDriver.Get(api.ChatGPTUrl)
			HandleCaptcha(WebDriver)
		}
	}
}
