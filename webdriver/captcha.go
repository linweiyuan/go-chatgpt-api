package webdriver

import (
	"github.com/tebeka/selenium"
	"log"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func HandleCaptcha(webDriver selenium.WebDriver) {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "mb-2")
		if err != nil {
			return false, nil
		}

		welcomeText, _ := element.Text()
		return welcomeText == "Welcome to ChatGPT", nil
	}, time.Second*10, time.Second*2)

	if err != nil {
		webDriver.SwitchFrame(0)

		err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
			element, err := driver.FindElement(selenium.ByCSSSelector, "input[type=checkbox]")
			if err != nil {
				return false, nil
			}

			element.Click()
			return true, nil
		}, time.Second*10, time.Second*2)

		if err != nil {
			log.Fatal("Failed to handle captcha.")
		}
	}
}
