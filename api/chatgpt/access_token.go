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
func (userLogin *UserLogin) GetAuthorizedUrl(csrfToken string) (string, int, error) {
	params := fmt.Sprintf(
		"callbackUrl=/&csrfToken=%s&json=true",
		csrfToken,
	)
	req, err := http.NewRequest(http.MethodPost, promptLoginUrl, strings.NewReader(params))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := userLogin.client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New(api.GetAuthorizedUrlErrorMessage)
	}

	responseMap := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&responseMap)
	return responseMap["url"], http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (userLogin *UserLogin) GetState(authorizedUrl string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, authorizedUrl, nil)
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err := userLogin.client.Do(req)
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
func (userLogin *UserLogin) CheckUsername(state string, username string) (int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&js-available=true&webauthn-available=true&is-brave=false&webauthn-platform-available=false&action=default",
		state,
		username,
	)
	req, _ := http.NewRequest(http.MethodPost, api.LoginUsernameUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err := userLogin.client.Do(req)
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
func (userLogin *UserLogin) CheckPassword(state string, username string, password string) (string, int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&password=%s&action=default",
		state,
		username,
		password,
	)
	req, err := http.NewRequest(http.MethodPost, api.LoginPasswordUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	userLogin.client.SetFollowRedirect(false) // make sure the cookie is injected with host chat.openai.com
	resp, err := userLogin.client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		alert := doc.Find("#prompt-alert").Text()
		if alert != "" {
			return "", resp.StatusCode, errors.New(strings.TrimSpace(alert))
		}

		return "", resp.StatusCode, errors.New(api.EmailOrPasswordInvalidErrorMessage)
	}

	if resp.StatusCode == http.StatusFound {
		req, _ := http.NewRequest(http.MethodGet, api.Auth0Url+resp.Header.Get("Location"), nil)
		req.Header.Set("User-Agent", api.UserAgent)
		resp, err := userLogin.client.Do(req)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusFound {
			req, _ := http.NewRequest(http.MethodGet, resp.Header.Get("Location"), nil)
			req.Header.Set("User-Agent", api.UserAgent)
			api.InjectCookies(req) // if not set this, will get 403 in some IPs
			resp, err := userLogin.client.Do(req)
			if err != nil {
				return "", http.StatusInternalServerError, err
			}

			defer resp.Body.Close()
			if resp.StatusCode == http.StatusFound {
				return "", http.StatusOK, nil
			}

			if resp.StatusCode == http.StatusTemporaryRedirect {
				errorDescription := req.URL.Query().Get("error_description")
				if errorDescription != "" {
					return "", resp.StatusCode, errors.New(errorDescription)
				}
			}

			return "", resp.StatusCode, errors.New(api.GetAccessTokenErrorMessage)
		}

		return "", resp.StatusCode, errors.New(api.EmailOrPasswordInvalidErrorMessage)
	}

	return "", resp.StatusCode, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat,GoUnusedParameter
func (userLogin *UserLogin) GetAccessToken(code string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, authSessionUrl, nil)
	req.Header.Set("User-Agent", api.UserAgent)
	api.InjectCookies(req)
	resp, err := userLogin.client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			responseMap := make(map[string]string)
			json.NewDecoder(resp.Body).Decode(&responseMap)
			return "", resp.StatusCode, errors.New(responseMap["detail"])
		}

		return "", resp.StatusCode, errors.New(api.GetAccessTokenErrorMessage)
	}

	data, _ := io.ReadAll(resp.Body)
	return string(data), http.StatusOK, nil
}
