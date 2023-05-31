package chatgpt

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
)

const (
	healthCheckUrl       = "https://chat.openai.com/backend-api/accounts/check"
	readyHint            = "Service go-chatgpt-api is ready."
	defaultCookiesApiUrl = "https://chatgpt.linweiyuan.com/cookies"
	errorHint403         = "If you still hit 403, do not raise new issue (will be closed directly without comment), change to a new clean IP or use legacy version first."
	errorHintBlock       = "It appears that you have been blocked from using the cookies API due to your IP being detected by Cloudflare WAF. Please send your IP address to root@linweiyuan.com so that it can be added to the whitelist. More details: https://github.com/linweiyuan/go-chatgpt-api/issues/134"
	cookieName           = "__cf_bm"
	sleepHours           = 8760 // 365 days
)

//goland:noinspection GoSnakeCaseUsage
var cfbm *Cookie
var firstTime = true

//goland:noinspection GoUnhandledErrorResult
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
	if !firstTime {
		injectCookies(req)
	}
	resp, err = api.Client.Do(req)
	return
}

//goland:noinspection GoUnhandledErrorResult
func checkHealthCheckStatus(resp *http.Response) {
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		logger.Info(readyHint)
		firstTime = false
	} else {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		alert := doc.Find(".message").Text()
		if alert != "" {
			logger.Error(strings.TrimSpace(alert) + " by OpenAI.")
			time.Sleep(time.Hour * sleepHours)
			os.Exit(1)
		} else {
			getCookies()
		}
	}
}

func getCookiesApiUrl() string {
	cookiesApiUrl := os.Getenv("GO_CHATGPT_API_COOKIES_API_URL")
	if cookiesApiUrl == "" {
		cookiesApiUrl = defaultCookiesApiUrl
	}
	return cookiesApiUrl
}

//goland:noinspection GoUnhandledErrorResult
func getCookies() {
	req, _ := http.NewRequest(http.MethodGet, getCookiesApiUrl(), nil)
	resp, err := api.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			logger.Error(errorHintBlock)
			time.Sleep(time.Hour * sleepHours)
			os.Exit(1)
		}

		logger.Error("Failed to get cookies, please try again later.")
		return
	}

	defer resp.Body.Close()
	var cookies []*Cookie
	err = json.NewDecoder(resp.Body).Decode(&cookies)
	if err != nil {
		logger.Error("Failed to parse cookies, please retry later.")
		return
	}

	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			cfbm = cookie
			break
		}
	}

	if firstTime {
		logger.Info(readyHint)
		logger.Error(errorHint403)
		firstTime = false

		go func() {
			keepCheckingCookies()
		}()
	}
}

func injectCookies(req *http.Request) {
	if cfbm != nil {
		req.Header.Set("Cookie", cookieName+"="+cfbm.Value)
	}
}

//goland:noinspection GoUnhandledErrorResult
func keepCheckingCookies() {
	for {
		now := time.Now()
		refreshTime := now.Add(time.Minute * 5) // // refresh cookie 5 minutes before it is expired
		if refreshTime.Minute() == time.Unix(cfbm.Expiry, 0).Minute() {
			// use old cookie to get back new cookie
			resp, err := healthCheck()
			if err == nil && resp.StatusCode == http.StatusUnauthorized {
				oldValue := cfbm.Value
				for _, cookie := range resp.Cookies() {
					if cookie.Name == cookieName {
						cfbm = &Cookie{
							Name:  cookie.Name,
							Value: cookie.Value,
						}
						break
					}
				}

				newValue := cfbm.Value
				if oldValue == newValue {
					go func() {
						for {
							time.Sleep(time.Minute * 20)
							getCookies()
						}
					}()
				} else {
					// if new cfbm is set, go-chatgpt-api itself will take over the task of refreshing cookie from external cookies api
					go func() {
						for {
							time.Sleep(time.Minute * 25)
							resp, err := healthCheck()
							if err != nil {
								logger.Error("Health check failed: " + err.Error())
								os.Exit(1)
							}

							for _, cookie := range resp.Cookies() {
								if cookie.Name == cookieName {
									cfbm = &Cookie{
										Name:  cookie.Name,
										Value: cookie.Value,
									}
								}
							}
							resp.Body.Close()
						}
					}()
				}
				break
			}
		}

		time.Sleep(time.Minute)
	}
}
