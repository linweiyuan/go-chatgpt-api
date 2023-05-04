package platform

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult
func ChatCompletions(c *gin.Context) {
	var chatCompletionsRequest ChatCompletionsRequest
	c.ShouldBindJSON(&chatCompletionsRequest)
	data, _ := json.Marshal(chatCompletionsRequest)
	req, _ := http.NewRequest(http.MethodPost, apiChatCompletions, bytes.NewBuffer(data))
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	api.HandleConversationResponse(c, resp)
}

func GetModels(c *gin.Context) {
	handleGet(c, apiGetModels)
}

func GetCreditGrants(c *gin.Context) {
	handleGet(c, apiGetCreditGrants)
}

//goland:noinspection GoUnhandledErrorResult
func Login(c *gin.Context) {
	var loginInfo api.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(api.ParseUserInfoErrorMessage))
		return
	}

	// hard refresh cookies
	resp, _ := api.Client.Get(auth0LogoutUrl)
	defer resp.Body.Close()

	userLogin := new(UserLogin)

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
	req.Header.Set("Authorization", api.GetAccessToken(getAccessTokenResponse.AccessToken))
	resp, err = api.Client.Do(req)
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

func GetSubscription(c *gin.Context) {
	handleGet(c, apiGetSubscription)
}

func GetApiKeys(c *gin.Context) {
	handleGet(c, apiGetApiKeys)
}

//goland:noinspection GoUnhandledErrorResult
func handleGet(c *gin.Context, url string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	resp, _ := api.Client.Do(req)
	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}
