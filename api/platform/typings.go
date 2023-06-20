package platform

//goland:noinspection GoSnakeCaseUsage
import tls_client "github.com/bogdanfinn/tls-client"

type UserLogin struct {
	client tls_client.HttpClient
}

type GetAccessTokenRequest struct {
	ClientID    string `json:"client_id"`
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
