package platform

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult
func Login(c *gin.Context) {
	var loginInfo api.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.ParseUserInfoErrorMessage))
		return
	}

	userLogin := UserLogin{
		client: api.NewHttpClient(),
	}

	// hard refresh cookies
	resp, _ := userLogin.client.Get(auth0LogoutUrl)
	defer resp.Body.Close()

	// get authorized url
	authorizedUrl, statusCode, err := userLogin.GetAuthorizedUrl("")
	if err != nil {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(err.Error()))
		return
	}

	// get state
	state, _, _ := userLogin.GetState(authorizedUrl)

	// check username
	statusCode, err = userLogin.CheckUsername(state, loginInfo.Username)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(err.Error()))
		return
	}

	// check password
	code, statusCode, err := userLogin.CheckPassword(state, loginInfo.Username, loginInfo.Password)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(err.Error()))
		return
	}

	// get access token
	accessToken, statusCode, err := userLogin.GetAccessToken(code)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(err.Error()))
		return
	}

	// get session key
	var getAccessTokenResponse GetAccessTokenResponse
	json.Unmarshal([]byte(accessToken), &getAccessTokenResponse)
	req, _ := http.NewRequest(http.MethodPost, dashboardLoginUrl, strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set(api.AuthorizationHeader, "Bearer "+getAccessTokenResponse.AccessToken)
	resp, err = userLogin.client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(resp.StatusCode, api.ReturnMessage(getSessionKeyErrorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}
