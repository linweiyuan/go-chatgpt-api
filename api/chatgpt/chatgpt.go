package chatgpt

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoSnakeCaseUsage
var __cf_bm = "" // https://developers.cloudflare.com/fundamentals/get-started/reference/cloudflare-cookies/#__cf_bm-cookie-for-cloudflare-bot-products
var firstTime = true

//goland:noinspection GoUnhandledErrorResult
func init() {
	go func() {
		ticker := time.NewTicker(time.Minute * healthCheckInterval)
		for {
			select {
			case <-ticker.C:
				resp, err := healthCheck()
				if err != nil && resp.StatusCode != http.StatusOK {
					getCookies()
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Minute * getCookiesInterval)
		for {
			select {
			case <-ticker.C:
				getCookies()
			}
		}
	}()

	proxyUrl := os.Getenv("GO_CHATGPT_API_PROXY")
	if proxyUrl == "" {
		resp, _ := healthCheck()
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		if string(data) == "error code: 1020" {
			logger.Error(accessDeniedText)
			return
		}

		checkHealthCheckStatus(resp)
	} else {
		for {
			resp, err := healthCheck()
			if err != nil {
				// wait for warp-svc to be ready
				time.Sleep(time.Second)
				continue
			}

			checkHealthCheckStatus(resp)
			break
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func checkHealthCheckStatus(resp *http.Response) {
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == http.StatusOK {
		logger.Info(welcomeText)
	} else {
		logger.Warn(healthCheckFailedText)

		getCookies()
	}
}

func healthCheck() (resp *http.Response, err error) {
	req, _ := http.NewRequest(http.MethodGet, authSessionUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	injectCookies(req)
	resp, err = api.Client.Do(req)
	return
}

//goland:noinspection GoUnhandledErrorResult
func getCookies() {
	req, _ := http.NewRequest(http.MethodGet, getCookiesUrl, nil)
	resp, err := api.Client.Do(req)
	if err != nil {
		logger.Error("Failed to get cookies: " + err.Error())
		return
	}

	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	responseMap := make(map[string]string)
	json.Unmarshal(data, &responseMap)
	__cf_bm = responseMap["__cf_bm"]
	if __cf_bm == "" {
		logger.Error(getCookiesFailedText)
		return
	}

	if firstTime {
		logger.Info("Get cookies successfully: " + __cf_bm)
		logger.Info(welcomeText)
		firstTime = false
	}
}

func injectCookies(req *http.Request) {
	if __cf_bm != "" {
		req.Header.Set("Cookie", "__cf_bm="+__cf_bm)
	}
}

//goland:noinspection GoUnhandledErrorResult
func handleGet(c *gin.Context, url string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	injectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, api.ReturnMessage(errorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func handlePost(c *gin.Context, url string, requestBody string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(requestBody))
	handlePostOrPatch(c, req, errorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handlePatch(c *gin.Context, url string, requestBody string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodPatch, url, strings.NewReader(requestBody))
	handlePostOrPatch(c, req, errorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handlePostOrPatch(c *gin.Context, req *http.Request, errorMessage string) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	injectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, api.ReturnMessage(errorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}
