package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

func ListModels(c *gin.Context) {
	handleGet(c, apiListModels)
}

func RetrieveModel(c *gin.Context) {
	model := c.Param("model")
	handleGet(c, fmt.Sprintf(apiRetrieveModel, model))
}

//goland:noinspection GoUnhandledErrorResult
func CreateCompletions(c *gin.Context) {
	var request CreateCompletionsRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreateCompletions, data, request.Stream)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if request.Stream {
		api.HandleConversationResponse(c, resp)
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

//goland:noinspection GoUnhandledErrorResult
func CreateChatCompletions(c *gin.Context) {
	var request ChatCompletionsRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreataeChatCompletions, data, request.Stream)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if request.Stream {
		api.HandleConversationResponse(c, resp)
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

//goland:noinspection GoUnhandledErrorResult
func CreateEdit(c *gin.Context) {
	var request CreateEditRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreateEdit, data, false)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func CreateImage(c *gin.Context) {
	var request CreateImageRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreateImage, data, false)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func CreateEmbeddings(c *gin.Context) {
	var request CreateEmbeddingsRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreateEmbeddings, data, false)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}

func CreateModeration(c *gin.Context) {
	var request CreateModerationRequest
	c.ShouldBindJSON(&request)
	data, _ := json.Marshal(request)
	resp, err := handlePost(c, apiCreateModeration, data, false)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}

func ListFiles(c *gin.Context) {
	handleGet(c, apiListFiles)
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
	req.Header.Set("Authorization", api.GetAccessToken(getAccessTokenResponse.AccessToken))
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

func handlePost(c *gin.Context, url string, data []byte, stream bool) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return nil, err
	}

	return resp, nil
}
