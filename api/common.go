package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bufio"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const (
	defaultErrorMessageKey = "errorMessage"
	AuthorizationHeader    = "Authorization"
)

var Client tls_client.HttpClient

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
