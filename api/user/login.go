package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthSession struct {
	User struct {
		ID      string        `json:"id"`
		Name    string        `json:"name"`
		Email   string        `json:"email"`
		Image   string        `json:"image"`
		Picture string        `json:"picture"`
		Groups  []interface{} `json:"groups"`
	} `json:"user"`
	Expires     time.Time `json:"expires"`
	AccessToken string    `json:"accessToken"`
	Cookies     string    `json:"cookies"`
}

func Login(c *gin.Context) {
	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse user login info."))
		return
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 0,
		Jar:     jar,
	}

	resp, err := client.Get("https://explorer.api.openai.com/api/auth/csrf")
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get CSRF token."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	responseMap := make(map[string]string)
	json.Unmarshal(body, &responseMap)
	resp, err = client.PostForm("https://explorer.api.openai.com/api/auth/signin/auth0?", url.Values{
		"callbackUrl": {"/"},
		"csrfToken":   {responseMap["csrfToken"]},
		"json":        {"true"},
	})
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get authorized url."))
		return
	}

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &responseMap)
	resp, err = client.Get(responseMap["url"])
	api.CheckError(c, err)
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get state."))
		return
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	state, _ := doc.Find("input[name=state]").Attr("value")
	resp, err = client.PostForm(fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state), url.Values{
		"username": {loginInfo.Username},
	})
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		c.JSON(resp.StatusCode, api.ReturnMessage("Email is not valid."))
		return
	}

	resp, err = client.PostForm(fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state), url.Values{
		"username": {loginInfo.Username},
		"password": {loginInfo.Password},
	})
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		c.JSON(resp.StatusCode, api.ReturnMessage("Email or password is not correct."))
		return
	}

	resp, err = client.Get("https://explorer.api.openai.com/api/auth/session")
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get access token, please try again later."))
		return
	}

	body, _ = io.ReadAll(resp.Body)
	if string(body) == "{}" {
		c.JSON(http.StatusForbidden, api.ReturnMessage("OpenAI's services are not available in your country."))
		return
	}

	var authSession AuthSession
	json.Unmarshal(body, &authSession)
	authSession.Cookies = GetCookiesString(resp.Cookies())

	c.JSON(http.StatusOK, authSession)
}

func GetCookiesString(cookies []*http.Cookie) string {
	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookieStrings, ":")
}
