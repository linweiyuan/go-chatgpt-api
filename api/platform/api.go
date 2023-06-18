package platform

import (
	"bufio"
	"bytes"
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

//goland:noinspection GoUnhandledErrorResult
func CreateChatCompletions(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var request struct {
		Stream bool `json:"stream"`
	}
	json.Unmarshal(body, &request)

	url := c.Request.URL.Path
	if strings.Contains(url, "/chat") {
		url = apiCreateChatCompletions
	} else {
		url = apiCreateCompletions
	}

	resp, err := handlePost(c, url, body, request.Stream)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if request.Stream {
		handleCompletionsResponse(c, resp)
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

func CreateCompletions(c *gin.Context) {
	CreateChatCompletions(c)
}

//goland:noinspection GoUnhandledErrorResult
func handleCompletionsResponse(c *gin.Context, resp *http.Response) {
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

	reader := bufio.NewReader(resp.Body)
	for {
		if c.Request.Context().Err() != nil {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "" {
			continue
		}

		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()
	}
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
