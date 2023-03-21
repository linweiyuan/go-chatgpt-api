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
	chatgptProxyServer := os.Getenv("CHATGPT_PROXY_SERVER")
	if chatgptProxyServer == "" {
		log.Fatal("Please set ChatGPT proxy server first.")
	}

	WebDriver, _ = selenium.NewRemote(selenium.Capabilities{
		"chromeOptions": chrome.Capabilities{
			Args:            chromeArgs,
			ExcludeSwitches: []string{"enable-automation"},
		},
	}, chatgptProxyServer)

	WebDriver.Get(api.ChatGPTUrl)

	WebDriver.SetAsyncScriptTimeout(time.Second * api.ScriptExecutionTimeout)
}
