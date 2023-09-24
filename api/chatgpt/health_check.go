package chatgpt

import (
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-logger/logger"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection SpellCheckingInspection
const (
	healthCheckUrl         = "https://chat.openai.com/backend-api/accounts/check"
	errorHintBlock         = "Looks like you have bean blocked by OpenAI, please change to a new IP or have a try with WARP."
	errorHintFailedToStart = "Failed to start, please try again later."
	sleepHours             = 8760 // 365 days
)

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func init() {
	proxyUrl := os.Getenv("PROXY")
	if proxyUrl != "" {
		logger.Info("PROXY: " + proxyUrl)
		api.Client.SetProxy(proxyUrl)

		for {
			resp, err := healthCheck()
			if err != nil {
				// wait for proxy to be ready
				time.Sleep(time.Second)
				continue
			}

			checkHealthCheckStatus(resp)
			break
		}
	} else {
		resp, err := healthCheck()
		if err != nil {
			logger.Error("Health check failed: " + err.Error())
			os.Exit(1)
		}

		checkHealthCheckStatus(resp)
	}
}

func healthCheck() (resp *http.Response, err error) {
	req, _ := http.NewRequest(http.MethodGet, healthCheckUrl, nil)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err = api.Client.Do(req)
	return
}

//goland:noinspection GoUnhandledErrorResult
func checkHealthCheckStatus(resp *http.Response) {
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		logger.Info(api.ReadyHint)
	} else {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		alert := doc.Find(".message").Text()
		if alert != "" {
			logger.Error(errorHintBlock)
		} else {
			logger.Error(errorHintFailedToStart)
			logger.Warn(doc.Text())
		}
		time.Sleep(time.Hour * sleepHours)
		os.Exit(1)
	}
}
