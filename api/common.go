package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const (
	defaultErrorMessageKey             = "errorMessage"
	AuthorizationHeader                = "Authorization"
	ContentType                        = "application/x-www-form-urlencoded"
	UserAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	LoginUsernameUrl                   = "https://auth0.openai.com/u/login/identifier?state="
	LoginPasswordUrl                   = "https://auth0.openai.com/u/login/password?state="
	ParseUserInfoErrorMessage          = "Failed to parse user login info."
	GetAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	GetStateErrorMessage               = "Failed to get state."
	EmailInvalidErrorMessage           = "Email is not valid."
	EmailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	GetAccessTokenErrorMessage         = "Failed to get access token, please try again later."

	AuthSessionUrl   = "https://chat.openai.com/api/auth/session"
	accessDeniedText = "Access denied, please set environment variable GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535 or something like this."
	welcomeText      = "Welcome to ChatGPT"
	getCookiesUrl    = "https://get-chatgpt-cookies.linweiyuan.com"

	healthCheckInterval = 15
	getCookiesInterval  = 25
)

var Client tls_client.HttpClient

//goland:noinspection GoSnakeCaseUsage
var __cf_bm = "" // https://developers.cloudflare.com/fundamentals/get-started/reference/cloudflare-cookies/#__cf_bm-cookie-for-cloudflare-bot-products
var firstTime = true

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin interface {
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
		resp, _ := healthCheck()
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		if string(data) == "error code: 1020" {
			logger.Error(accessDeniedText)
			return
		}

		checkHealthCheckStatus(resp)
	}

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
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == http.StatusOK {
		logger.Info(welcomeText)
		firstTime = false
	} else {
		getCookies()
	}
}

func healthCheck() (resp *http.Response, err error) {
	req, _ := http.NewRequest(http.MethodGet, AuthSessionUrl, nil)
	req.Header.Set("User-Agent", UserAgent)
	InjectCookies(req)
	resp, err = Client.Do(req)
	return
}

//goland:noinspection GoUnhandledErrorResult
func getCookies() {
	req, _ := http.NewRequest(http.MethodGet, getCookiesUrl, nil)
	resp, err := Client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	responseMap := make(map[string]string)
	json.Unmarshal(data, &responseMap)
	__cf_bm = responseMap["__cf_bm"]
	if __cf_bm == "" {
		return
	}

	if firstTime {
		logger.Info(welcomeText)
		firstTime = false
	}
}

func InjectCookies(req *http.Request) {
	if __cf_bm != "" {
		req.Header.Set("Cookie", "__cf_bm="+__cf_bm)
	}
}
