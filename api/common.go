package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-logger/logger"
	"github.com/xqdoo00o/funcaptcha"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

//goland:noinspection SpellCheckingInspection
const (
	ChatGPTApiPrefix    = "/chatgpt"
	ImitateApiPrefix    = "/imitate/v1"
	ChatGPTApiUrlPrefix = "https://chat.openai.com"

	PlatformApiPrefix    = "/platform"
	PlatformApiUrlPrefix = "https://api.openai.com"

	defaultErrorMessageKey             = "errorMessage"
	AuthorizationHeader                = "Authorization"
	XAuthorizationHeader               = "X-Authorization"
	ContentType                        = "application/x-www-form-urlencoded"
	UserAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	Auth0Url                           = "https://auth0.openai.com"
	LoginUsernameUrl                   = Auth0Url + "/u/login/identifier?state="
	LoginPasswordUrl                   = Auth0Url + "/u/login/password?state="
	ParseUserInfoErrorMessage          = "Failed to parse user login info."
	GetAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	GetStateErrorMessage               = "Failed to get state."
	EmailInvalidErrorMessage           = "Email is not valid."
	EmailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	GetAccessTokenErrorMessage         = "Failed to get access token."
	GetArkoseTokenErrorMessage         = "Failed to get arkose token."
	defaultTimeoutSeconds              = 600 // 10 minutes

	EmailKey                       = "email"
	AccountDeactivatedErrorMessage = "Account %s is deactivated."

	ReadyHint = "Service go-chatgpt-api is ready."
)

var (
	Client       tls_client.HttpClient
	ArkoseClient tls_client.HttpClient
)

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
	GetAccessTokenFromHeader(c *gin.Context) (string, int, error)
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithTimeoutSeconds(defaultTimeoutSeconds),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
	}...)
	ArkoseClient = getHttpClient()
	funcaptcha.SetTLSClient(&Client)
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func NewHttpClient() tls_client.HttpClient {
	client := getHttpClient()

	proxyUrl := os.Getenv("PROXY")
	if proxyUrl != "" {
		client.SetProxy(proxyUrl)
	}

	return client
}

func getHttpClient() tls_client.HttpClient {
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
	}...)
	return client
}

//goland:noinspection GoUnhandledErrorResult
func Proxy(c *gin.Context) {
	url := c.Request.URL.Path
	if strings.Contains(url, ChatGPTApiPrefix) {
		url = strings.ReplaceAll(url, ChatGPTApiPrefix, ChatGPTApiUrlPrefix)
	} else if strings.Contains(url, ImitateApiPrefix) {
		url = strings.ReplaceAll(url, ImitateApiPrefix, ChatGPTApiUrlPrefix+"/backend-api")
	} else {
		url = strings.ReplaceAll(url, PlatformApiPrefix, PlatformApiUrlPrefix)
	}

	method := c.Request.Method
	queryParams := c.Request.URL.Query().Encode()
	if queryParams != "" {
		url += "?" + queryParams
	}

	// if not set, will return 404
	c.Status(http.StatusOK)

	var req *http.Request
	if method == http.MethodGet {
		req, _ = http.NewRequest(http.MethodGet, url, nil)
	} else {
		body, _ := io.ReadAll(c.Request.Body)
		req, _ = http.NewRequest(method, url, bytes.NewReader(body))
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set(AuthorizationHeader, GetAccessToken(c))
	resp, err := Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			logger.Error(fmt.Sprintf(AccountDeactivatedErrorMessage, c.GetString(EmailKey)))
		}

		responseMap := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseMap)
		c.AbortWithStatusJSON(resp.StatusCode, responseMap)
		return
	}

	io.Copy(c.Writer, resp.Body)
}

func ReturnMessage(msg string) gin.H {
	logger.Warn(msg)

	return gin.H{
		defaultErrorMessageKey: msg,
	}
}

func GetAccessToken(c *gin.Context) string {
	accessToken := c.GetString(AuthorizationHeader)
	if !strings.HasPrefix(accessToken, "Bearer") {
		return "Bearer " + accessToken
	}

	return accessToken
}
