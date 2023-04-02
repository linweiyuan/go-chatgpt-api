package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/tebeka/selenium"
	"time"
)

const (
	checkWelcomeTextTimeout = 3
	checkCaptchaTimeout     = 8
	checkCaptchaInterval    = 1
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
			element, err := driver.FindElement(selenium.ByCSSSelector, "input")
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
				logger.Error(pageSource)
			}
		}
	}
}
