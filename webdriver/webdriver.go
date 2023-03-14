package webdriver

import (
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
	"time"
)

var WebDriver selenium.WebDriver

//goland:noinspection GoUnhandledErrorResult
func init() {
	chromeArgs := []string{
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-blink-features=AutomationControlled",
	}
	proxyServer := os.Getenv("PROXY_SERVER")
	if proxyServer != "" {
		chromeArgs = append(chromeArgs, "--proxy-server="+proxyServer)
	}
	webDriverUrl := os.Getenv("WEB_DRIVER_URL")
	if webDriverUrl == "" {
		log.Fatal("Please set web driver url first.")
	}

	WebDriver, _ = selenium.NewRemote(selenium.Capabilities{
		"chromeOptions": chrome.Capabilities{
			Args:            chromeArgs,
			ExcludeSwitches: []string{"enable-automation"},
		},
	}, webDriverUrl)

	WebDriver.Get(api.ChatGPTUrl)

	WebDriver.SetAsyncScriptTimeout(time.Second * api.ScriptExecutionTimeout)
}
