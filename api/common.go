package api

import (
	"github.com/gin-gonic/gin"
)

const defaultErrorMessageKey = "errorMessage"

const (
	ChatGPTUrl             = "https://chat.openai.com/chat"
	ScriptExecutionTimeout = 10

	Authorization = "Authorization"

	DoneFlag = "[DONE]"
)

func ReturnMessage(msg string) gin.H {
	return gin.H{
		defaultErrorMessageKey: msg,
	}
}
