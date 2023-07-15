package token

import (
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"

	http "github.com/bogdanfinn/fhttp"
)

func GetArkoseToken(c *gin.Context) {
	token, err := chatgpt.GetArkoseToken()
	if err != nil || token == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage(api.GetArkoseTokenErrorMessage))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
