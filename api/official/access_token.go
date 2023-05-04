package official

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat,GoUnusedParameter
func (user *PlatformUserLogin) GetAuthorizedUrl(csrfToken string) (string, int, error) {
	urlParams := url.Values{
		"client_id":     {platformAuthClientID},
		"redirect_uri":  {platformAuthRedirectURL},
		"scope":         {platformAuthScope},
		"response_type": {platformAuthResponseType},
	}
	req, _ := http.NewRequest(http.MethodGet, "https://auth0.openai.com/authorize?"+urlParams.Encode(), nil)
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, errors.New(api.GetAuthorizedUrlErrorMessage)
	}

	return resp.Request.URL.String(), http.StatusOK, nil
}

func (user *PlatformUserLogin) GetState(authorizedUrl string) (string, int, error) {
	split := strings.Split(authorizedUrl, "=")
	return split[1], http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *PlatformUserLogin) CheckUsername(state string, username string) (int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&js-available=true&webauthn-available=true&is-brave=false&webauthn-platform-available=false&action=default",
		state,
		username,
	)
	req, err := http.NewRequest(http.MethodPost, api.LoginUsernameUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
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
func (user *PlatformUserLogin) CheckPassword(state string, username string, password string) (string, int, error) {
	formParams := fmt.Sprintf(
		"state=%s&username=%s&password=%s&action=default",
		state,
		username,
		password,
	)
	req, err := http.NewRequest(http.MethodPost, api.LoginPasswordUrl+state, strings.NewReader(formParams))
	req.Header.Set("Content-Type", api.ContentType)
	req.Header.Set("User-Agent", api.UserAgent)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden { // platform is 403
		return "", resp.StatusCode, errors.New(api.EmailOrPasswordInvalidErrorMessage)
	}

	return resp.Request.URL.Query().Get("code"), http.StatusOK, nil
}

//goland:noinspection GoUnhandledErrorResult,GoErrorStringFormat
func (user *PlatformUserLogin) GetAccessToken(code string) (string, int, error) {
	jsonBytes, _ := json.Marshal(GetAccessTokenRequest{
		ClientID:    platformAuthClientID,
		Code:        code,
		GrantType:   platformAuthGrantType,
		RedirectURI: platformAuthRedirectURL,
	})
	req, err := http.NewRequest(http.MethodPost, getTokenUrl, strings.NewReader(string(jsonBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", api.UserAgent)
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
