package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const defaultErrorMessageKey = "errorMessage"

func ReturnMessage(msg string) gin.H {
	return gin.H{
		defaultErrorMessageKey: msg,
	}
}

func CheckError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			defaultErrorMessageKey: "ChatGPT is at capacity right now, please try again later.",
		})
	}
}
