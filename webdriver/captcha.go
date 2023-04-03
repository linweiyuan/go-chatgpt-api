package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/tebeka/selenium"
	"time"
)

const (
	checkWelcomeTextTimeout = 5
	checkCaptchaTimeout     = 15
	checkCaptchaInterval    = 1
	checkNextInterval       = 5
)

//goland:noinspection GoUnhandledErrorResult
func HandleCaptcha(webDriver selenium.WebDriver) {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "mb-2")
		if err != nil {
			return false, nil
		}

		welcomeText, _ := element.Text()
		logger.Info(welcomeText)
		return welcomeText == "Welcome to ChatGPT", nil
	}, time.Second*checkWelcomeTextTimeout, time.Second*checkCaptchaInterval)

	if err != nil {
		webDriver.SwitchFrame(0)

		logger.Info("Checking captcha")
		err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
			element, err := driver.FindElement(selenium.ByCSSSelector, "input[type=checkbox]")
			if err != nil {
				return false, nil
			}

			element.Click()
			logger.Info("Captcha is clicked!")
			return true, nil
		}, time.Second*checkCaptchaTimeout, time.Second*checkCaptchaInterval)

		if err != nil {
			logger.Error("Failed to handle captcha: " + err.Error())
			if pageSource, err := webDriver.PageSource(); err == nil {
				title, _ := webDriver.Title()
				logger.Error(title)
				logger.Error(pageSource)
			}
			webDriver.Refresh()
		} else {
			time.Sleep(time.Second * checkNextInterval)

			title, _ := webDriver.Title()
			logger.Info(title)
			if title == "Just a moment..." {
				logger.Info("Still get a captcha")

				HandleCaptcha(webDriver)
			}
		}
	}
}
