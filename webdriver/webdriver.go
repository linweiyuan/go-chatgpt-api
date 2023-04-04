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
		"--no-sandbox",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--disable-blink-features=AutomationControlled",
		"--headless=new",
	}

	networkProxyServer := os.Getenv("NETWORK_PROXY_SERVER")
	if networkProxyServer != "" {
		chromeArgs = append(chromeArgs, "--proxy-server="+networkProxyServer)
	}

	chatgptProxyServer := os.Getenv("CHATGPT_PROXY_SERVER")
	if chatgptProxyServer == "" {
		log.Fatal("Please set ChatGPT proxy server first")
	}

	WebDriver, _ = selenium.NewRemote(selenium.Capabilities{
		"chromeOptions": chrome.Capabilities{
			Args:            chromeArgs,
			ExcludeSwitches: []string{"enable-automation"},
		},
	}, chatgptProxyServer)

	WebDriver.Get(api.ChatGPTUrl)
	if !isAccessDenied(WebDriver) && !isAtCapacity(WebDriver) {
		HandleCaptcha(WebDriver)
		WebDriver.SetAsyncScriptTimeout(time.Second * api.ScriptExecutionTimeout)
	}
}
