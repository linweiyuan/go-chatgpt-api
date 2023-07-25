package chatgpt

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	statusCode, errorMessage, accessTokenResponse := GetAccessToken(loginInfo)
	if statusCode != http.StatusOK {
		c.AbortWithStatusJSON(statusCode, api.ReturnMessage(errorMessage))
		return
	}

	c.Writer.WriteString(accessTokenResponse)
}

//goland:noinspection GoUnhandledErrorResult
func GetAccessToken(loginInfo api.LoginInfo) (int, string, string) {
	userLogin := UserLogin{
		client: api.NewHttpClient(),
	}

	// get csrf token
	req, _ := http.NewRequest(http.MethodGet, csrfUrl, nil)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err := userLogin.client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), ""
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			doc, _ := goquery.NewDocumentFromReader(resp.Body)
			alert := doc.Find(".message").Text()
			if alert != "" {
				return resp.StatusCode, strings.TrimSpace(alert), ""
			}
		}

		return resp.StatusCode, getCsrfTokenErrorMessage, ""
	}

	// get authorized url
	responseMap := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&responseMap)
	authorizedUrl, statusCode, err := userLogin.GetAuthorizedUrl(responseMap["csrfToken"])
	if err != nil {
		return statusCode, err.Error(), ""
	}

	// get state
	state, statusCode, err := userLogin.GetState(authorizedUrl)
	if err != nil {
		return statusCode, err.Error(), ""
	}

	// check username
	statusCode, err = userLogin.CheckUsername(state, loginInfo.Username)
	if err != nil {
		return statusCode, err.Error(), ""
	}

	// check password
	_, statusCode, err = userLogin.CheckPassword(state, loginInfo.Username, loginInfo.Password)
	if err != nil {
		return statusCode, err.Error(), ""
	}

	// get access token
	accessToken, statusCode, err := userLogin.GetAccessToken("")
	if err != nil {
		return statusCode, err.Error(), ""
	}

	return http.StatusOK, "", accessToken
}
