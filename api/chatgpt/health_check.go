package chatgpt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"

	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-logger/logger"
)

const (
	healthCheckUrl         = "https://chat.openai.com/backend-api/accounts/check"
	errorHintBlock         = "looks like you have bean blocked by OpenAI, please change to a new IP or have a try with WARP"
	errorHintFailedToStart = "failed to start, please try again later: %s"
	sleepHours             = 8760 // 365 days
)

func init() {
	proxyUrl := os.Getenv("PROXY")
	if proxyUrl != "" {
		logger.Info("PROXY: " + proxyUrl)
		err1 := api.Client.SetProxy(proxyUrl)
		fmt.Println(err1)
		if err := api.Client.SetProxy(proxyUrl); err != nil {
			log.Fatal(err)
		}

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
			logger.Error("failed to health check: " + err.Error())
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

func checkHealthCheckStatus(resp *http.Response) {
	if resp != nil {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			logger.Info(api.ReadyHint)
		} else {
			doc, _ := goquery.NewDocumentFromReader(resp.Body)
			alert := doc.Find(".message").Text()
			if alert != "" {
				logger.Error(errorHintBlock)
			} else {
				logger.Error(fmt.Sprintf(errorHintFailedToStart, resp.Status))
			}
			time.Sleep(time.Hour * sleepHours)
			os.Exit(1)
		}
	}
}
