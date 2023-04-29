package api

import (
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultApiTimeoutSeconds = 30
	defaultErrorMessageKey   = "errorMessage"
	AuthorizationHeader      = "Authorization"
)

var Client tls_client.HttpClient

func init() {
	Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(defaultApiTimeoutSeconds),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}...)
}

func ReturnMessage(msg string) gin.H {
	return gin.H{
		defaultErrorMessageKey: msg,
	}
}

func handleError(err error) gin.H {
	return gin.H{
		defaultErrorMessageKey: err.Error(),
	}
}

func returnError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, handleError(err))
}

func CheckError(c *gin.Context, err error) {
	if err != nil {
		returnError(c, err)
	}
}

func GetAccessToken(accessToken string) string {
	if !strings.HasPrefix(accessToken, "Bearer") {
		return "Bearer " + accessToken
	}
	return accessToken
}
