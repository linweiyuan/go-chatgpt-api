package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/user"
)

func RenewAccessToken(c *gin.Context) {
	client := &http.Client{
		Timeout: 0,
	}

	req, _ := http.NewRequest("GET", "https://explorer.api.openai.com/api/auth/session", nil)
	req.Header.Set("Cookie", c.GetHeader("Cookie"))
	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to refresh access token, please try again later."))
		return
	}

	body, _ := io.ReadAll(resp.Body)

	var authSession user.AuthSession
	json.Unmarshal(body, &authSession)
	authSession.Cookies = user.GetCookiesString(req.Cookies())

	c.JSON(http.StatusOK, authSession)
}
