package webdriver

import (
	"log"
	"time"

	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/tebeka/selenium"
)

const (
	checkWelcomeTextTimeout  = 5
	checkCaptchaTimeout      = 15
	checkAccessDeniedTimeout = 3
	checkAvailabilityTimeout = 3
	checkCaptchaInterval     = 1
	checkNextInterval        = 5
)

var isFirstTimeRun = true

//goland:noinspection GoUnhandledErrorResult
func HandleCaptcha(webDriver selenium.WebDriver) {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "mb-2")
		if err != nil {
			return false, nil
		}

		welcomeText, _ := element.Text()
		isWelcomed := welcomeText == "Welcome to ChatGPT"
		if isFirstTimeRun && isWelcomed {
			logger.Info(welcomeText)
			isFirstTimeRun = false
		}
		return isWelcomed, nil
	}, time.Second*checkWelcomeTextTimeout, time.Second*checkCaptchaInterval)

	if err != nil {
		webDriver.SwitchFrame(0)

		err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
			element, err := driver.FindElement(selenium.ByCSSSelector, "input")
			if err != nil {
				return false, nil
			}

			element.Click()
			return true, nil
		}, time.Second*checkCaptchaTimeout, time.Second*checkCaptchaInterval)

		if err != nil {
			logger.Error("Failed to handle captcha: " + err.Error())

			webDriver.Refresh()
			HandleCaptcha(webDriver)
		} else {
			time.Sleep(time.Second * checkNextInterval)

			title, _ := webDriver.Title()
			if title == "" {
				log.Fatal("Failed to handle captcha, looks like infinite loop, please remove CHATGPT_PROXY_SERVER to use API mode first until I find a way to fix it.")
			}

			if title == "Just a moment..." {
				HandleCaptcha(webDriver)
			}
		}
	}
}

func isAccessDenied(webDriver selenium.WebDriver) bool {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "cf-error-details")
		if err != nil {
			return false, nil
		}

		accessDeniedText, _ := element.Text()
		logger.Error(accessDeniedText)
		return true, nil
	}, time.Second*checkAccessDeniedTimeout, time.Second*checkCaptchaInterval)

	if err != nil {
		return false
	}

	return true
}

func isAtCapacity(webDriver selenium.WebDriver) bool {
	err := webDriver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
		element, err := driver.FindElement(selenium.ByClassName, "text-3xl")
		if err != nil {
			return false, nil
		}

		atCapacityText, _ := element.Text()
		logger.Error(atCapacityText)
		return true, nil
	}, time.Second*checkAvailabilityTimeout, time.Second*checkCaptchaInterval)

	if err != nil {
		return false
	}

	return true
}
