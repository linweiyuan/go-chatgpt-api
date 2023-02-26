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
