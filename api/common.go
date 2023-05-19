package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bufio"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/linweiyuan/go-chatgpt-api/env"

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
	defaultTimeoutSeconds              = 300 // 5 minutes
)

var Client tls_client.HttpClient

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
		tls_client.WithTimeoutSeconds(defaultTimeoutSeconds),
	}...)
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
