package chatgpt

import (
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/xqdoo00o/OpenAIAuth/auth"

	"github.com/linweiyuan/go-chatgpt-api/api"
)

func Login(c *gin.Context) {
	var loginInfo api.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.ParseUserInfoErrorMessage))
		return
	}

	authenticator := auth.NewAuthenticator(loginInfo.Username, loginInfo.Password, api.ProxyUrl)
	if err := authenticator.Begin(); err != nil {
		c.AbortWithStatusJSON(err.StatusCode, api.ReturnMessage(err.Details))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": authenticator.GetAccessToken(),
	})
}
