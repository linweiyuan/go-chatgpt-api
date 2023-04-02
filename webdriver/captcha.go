package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/tebeka/selenium"
	"time"
)

const (
	checkCaptchaTimeout  = 20
	checkCaptchaInterval = 2
)

//goland:noinspection GoUnhandledErrorResult
func HandleCaptcha(webDriver selenium.WebDriver) {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "mb-2")
		if err != nil {
			logger.Info("Retry to check captcha")
			return false, nil
		}

		welcomeText, _ := element.Text()
		logger.Info(welcomeText)
		return welcomeText == "Welcome to ChatGPT", nil
	}, time.Second*checkCaptchaTimeout, time.Second*checkCaptchaInterval)

	if err != nil {
		logger.Info("Switch frame")
		webDriver.SwitchFrame(0)

		err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
			element, err := driver.FindElement(selenium.ByCSSSelector, "input[type=checkbox]")
			if err != nil {
				logger.Info("Retry to click captcha")
				return false, nil
			}

			element.Click()
			logger.Info("Captcha is clicked!")
			return true, nil
		}, time.Second*10, time.Second*2)

		if err != nil {
			logger.Error("Failed to handle captcha")
		}
	}
}
