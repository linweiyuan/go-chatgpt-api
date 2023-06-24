package chatgpt

import (
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection SpellCheckingInspection
const (
	healthCheckUrl = "https://chat.openai.com/backend-api/accounts/check"
	errorHintBlock = "Looks like you have bean blocked -> curl https://chat.openai.com | grep '<p>' | awk '{$1=$1;print}'"
	errorHint403   = "Failed to handle 403, have a look at https://github.com/linweiyuan/java-chatgpt-api or use other more powerful alternatives (do not raise new issue about 403)."
	sleepHours     = 8760 // 365 days
)

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func init() {
	proxyUrl := os.Getenv("GO_CHATGPT_API_PROXY")
	if proxyUrl != "" {
		logger.Info("GO_CHATGPT_API_PROXY: " + proxyUrl)
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
			logger.Error(errorHint403)
		}
		time.Sleep(time.Hour * sleepHours)
		os.Exit(1)
	}
}
