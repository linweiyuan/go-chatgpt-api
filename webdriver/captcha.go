package webdriver

import (
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func HandleCaptcha(webDriver selenium.WebDriver) {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "mb-2")
		if err != nil {
			logrus.Info("Retry to check captcha.")
			return false, nil
		}

		welcomeText, _ := element.Text()
		logrus.Info(welcomeText)
		return welcomeText == "Welcome to ChatGPT", nil
	}, time.Second*10, time.Second*2)

	if err != nil {
		logrus.Info("switch frame")
		webDriver.SwitchFrame(0)

		err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
			element, err := driver.FindElement(selenium.ByCSSSelector, "input[type=checkbox]")
			if err != nil {
				logrus.Info("Retry to click captcha.")
				return false, nil
			}

			element.Click()
			logrus.Info("Captcha is clicked!")
			return true, nil
		}, time.Second*10, time.Second*2)

		if err != nil {
			logrus.Fatal("Failed to handle captcha.")
		}
	}
}
