package api

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"

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
