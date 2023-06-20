package platform

import "github.com/linweiyuan/go-chatgpt-api/api"

//goland:noinspection SpellCheckingInspection
const (
	apiCreateChatCompletions = api.PlatformApiUrlPrefix + "/v1/chat/completions"
	apiCreateCompletions     = api.PlatformApiUrlPrefix + "/v1/completions"

	platformAuthClientID      = "DRivsnm2Mu42T3KOpqdtwB3NYviHYzwD"
	platformAuthAudience      = "https://api.openai.com/v1"
	platformAuthRedirectURL   = "https://platform.openai.com/auth/callback"
	platformAuthScope         = "openid profile email offline_access"
	platformAuthResponseType  = "code"
	platformAuthGrantType     = "authorization_code"
	platformAuth0Url          = api.Auth0Url + "/authorize?"
	getTokenUrl               = api.Auth0Url + "/oauth/token"
	auth0Client               = "eyJuYW1lIjoiYXV0aDAtc3BhLWpzIiwidmVyc2lvbiI6IjEuMjEuMCJ9" // '{"name":"auth0-spa-js","version":"1.21.0"}'
	auth0LogoutUrl            = api.Auth0Url + "/v2/logout?returnTo=https%3A%2F%2Fplatform.openai.com%2Floggedout&client_id=" + platformAuthClientID + "&auth0Client=" + auth0Client
	dashboardLoginUrl         = "https://api.openai.com/dashboard/onboarding/login"
	getSessionKeyErrorMessage = "Failed to get session key."
)
