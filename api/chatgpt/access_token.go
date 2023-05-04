package chatgpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"

	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *ChatGPTUserLogin) GetAuthorizedUrl(csrfToken string) (string, int, error) {
	params := fmt.Sprintf(
		"callbackUrl=/&csrfToken=%s&json=true",
		csrfToken,
	)
	req, err := http.NewRequest(http.MethodPost, promptLoginUrl, strings.NewReader(params))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New(api.GetAuthorizedUrlErrorMessage)
	}

	responseMap := make(map[string]string)
	data, _ := io.ReadAll(resp.Body)
	json.Unmarshal(data, &responseMap)
	return responseMap["url"], http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *ChatGPTUserLogin) GetState(authorizedUrl string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, authorizedUrl, nil)
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New(api.GetStateErrorMessage)
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	state, _ := doc.Find("input[name=state]").Attr("value")
	return state, http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *ChatGPTUserLogin) CheckUsername(state string, username string) (int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&js-available=true&webauthn-available=true&is-brave=false&webauthn-platform-available=false&action=default",
		state,
		username,
	)
	req, err := http.NewRequest(http.MethodPost, api.LoginUsernameUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(api.EmailInvalidErrorMessage)
	}

	return http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *ChatGPTUserLogin) CheckPassword(state string, username string, password string) (string, int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&password=%s&action=default",
		state,
		username,
		password,
	)
	req, err := http.NewRequest(http.MethodPost, api.LoginPasswordUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { // ChatGPT is 200
		return "", resp.StatusCode, errors.New(api.EmailOrPasswordInvalidErrorMessage)
	}

	return "", http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat,GoUnusedParameter
func (user *ChatGPTUserLogin) GetAccessToken(code string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, api.AuthSessionUrl, nil)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New(api.GetAccessTokenErrorMessage)
	}

	data, _ := io.ReadAll(resp.Body)
	return string(data), http.StatusOK, nil
}
