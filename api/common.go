package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultErrorMessageKey = "detail"

const (
	ChatGPTUrl         = "https://chat.openai.com/chat"
	ChatGPTWelcomeText = "Welcome to ChatGPT"
	ChatGPTTitleText   = "ChatGPT"

	AuthorizationHeader = "Authorization"

	RefreshEveryMinutes = 1
)

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
