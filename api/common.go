package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const (
	defaultErrorMessageKey             = "errorMessage"
	AuthorizationHeader                = "Authorization"
	ContentType                        = "application/x-www-form-urlencoded"
	UserAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	Auth0Url                           = "https://auth0.openai.com"
	LoginUsernameUrl                   = Auth0Url + "/u/login/identifier?state="
	LoginPasswordUrl                   = Auth0Url + "/u/login/password?state="
	ParseUserInfoErrorMessage          = "Failed to parse user login info."
	GetAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	GetStateErrorMessage               = "Failed to get state."
	EmailInvalidErrorMessage           = "Email is not valid."
	EmailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	GetAccessTokenErrorMessage         = "Failed to get access token."

	healthCheckUrl       = "https://chat.openai.com/backend-api/accounts/check"
	welcomeHint          = "Welcome to ChatGPT"
	defaultCookiesApiUrl = "https://api.linweiyuan.com/chatgpt/cookies"
	errorHint403         = "If you still hit 403, do not raise new issue (will be closed directly without comment), change to a new clean IP or use legacy version first."
	errorHintBlock       = "You have been blocked to use cookies api because your IP is detected by Cloudflare WAF."
)

var Client tls_client.HttpClient

//goland:noinspection GoSnakeCaseUsage
var __cf_bm = ""
var firstTime = true

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLogin interface {
	GetAuthorizedUrl(csrfToken string) (string, int, error)
	GetState(authorizedUrl string) (string, int, error)
	CheckUsername(state string, username string) (int, error)
	CheckPassword(state string, username string, password string) (string, int, error)
	GetAccessToken(code string) (string, int, error)
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithTimeoutSeconds(0),
	}...)

	//goland:noinspection SpellCheckingInspection
	proxyUrl := os.Getenv("GO_CHATGPT_API_PROXY")
	if proxyUrl != "" {
		err := Client.SetProxy(proxyUrl)
		if err != nil {
			logger.Error("Failed to config proxy: " + err.Error())
			return
		}
		logger.Info("GO_CHATGPT_API_PROXY: " + proxyUrl)

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
	} else {
		resp, err := healthCheck()
		if err == nil {
			defer resp.Body.Close()

			checkHealthCheckStatus(resp)
		}
	}
}

func ReturnMessage(msg string) gin.H {
	return gin.H{
		defaultErrorMessageKey: msg,
	}
}

func GetAccessToken(accessToken string) string {
	if !strings.HasPrefix(accessToken, "Bearer") {
		return "Bearer " + accessToken
	}
	return accessToken
}

//goland:noinspection GoUnhandledErrorResult
func HandleConversationResponse(c *gin.Context, resp *http.Response) {
	reader := bufio.NewReader(resp.Body)
	for {
		if c.Request.Context().Err() != nil {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "" {
			continue
		}

		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()
	}
}

//goland:noinspection GoUnhandledErrorResult
func checkHealthCheckStatus(resp *http.Response) {
	cookiesApiUrl := os.Getenv("GO_CHATGPT_API_COOKIES_API_URL")

	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == http.StatusUnauthorized && cookiesApiUrl == "" {
		logger.Info(welcomeHint)
		firstTime = false
	} else {
		if cookiesApiUrl == "" {
			logger.Info("GO_CHATGPT_API_COOKIES_API_URL defaults to: " + defaultCookiesApiUrl)
			cookiesApiUrl = defaultCookiesApiUrl
		} else {
			logger.Info("GO_CHATGPT_API_COOKIES_API_URL: " + cookiesApiUrl)
		}

		go getCookiesSSE(cookiesApiUrl)
	}
}

func healthCheck() (resp *http.Response, err error) {
	req, _ := http.NewRequest(http.MethodGet, healthCheckUrl, nil)
	req.Header.Set("User-Agent", UserAgent)
	resp, err = Client.Do(req)
	return
}

//goland:noinspection GoUnhandledErrorResult,GoUnusedFunction
func getCookiesSSE(cookiesApiUrl string) {
	req, _ := http.NewRequest(http.MethodGet, cookiesApiUrl, nil)
	resp, err := Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			logger.Error(errorHintBlock)
			time.Sleep(time.Hour)
			os.Exit(1)
		}

		time.Sleep(time.Minute)
		getCookiesSSE(cookiesApiUrl)
		return
	}

	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "" {
			continue
		}

		if len(line) > 6 {
			responseMap := make(map[string]string)
			err = json.Unmarshal([]byte(line[6:]), &responseMap)
			if err == nil {
				__cf_bm = responseMap["__cf_bm"]

				if firstTime && __cf_bm != "" {
					logger.Info(welcomeHint)
					logger.Error(errorHint403)
					firstTime = false
				}
			}
		}
	}

	getCookiesSSE(cookiesApiUrl)
}

func InjectCookies(req *http.Request) {
	if __cf_bm != "" {
		req.Header.Set("Cookie", "__cf_bm="+__cf_bm)
	}
}

//goland:noinspection GoUnhandledErrorResult
func NewHttpClient() tls_client.HttpClient {
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}...)

	proxyUrl := os.Getenv("GO_CHATGPT_API_PROXY")
	if proxyUrl != "" {
		client.SetProxy(proxyUrl)
	}

	return client
}
