package platform

type UserLogin struct{}

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

//goland:noinspection SpellCheckingInspection
type CreateCompletionsRequest struct {
	Model            string                 `json:"model"`
	Prompt           string                 `json:"prompt,omitempty"`
	Suffix           string                 `json:"suffix,omitempty"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	Temperature      int                    `json:"temperature,omitempty"`
	TopP             int                    `json:"top_p,omitempty"`
	N                int                    `json:"n,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	Logprobs         int                    `json:"logprobs,omitempty"`
	Echo             bool                   `json:"echo,omitempty"`
	Stop             string                 `json:"stop,omitempty"`
	PresencePenalty  int                    `json:"presence_penalty,omitempty"`
	FrequencyPenalty int                    `json:"frequency_penalty,omitempty"`
	BestOf           int                    `json:"best_of,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
}

type ChatCompletionsRequest struct {
	Model            string                   `json:"model"`
	Messages         []ChatCompletionsMessage `json:"messages"`
	Temperature      int                      `json:"temperature,omitempty"`
	TopP             int                      `json:"top_p,omitempty"`
	N                int                      `json:"n,omitempty"`
	Stream           bool                     `json:"stream,omitempty"`
	Stop             string                   `json:"stop,omitempty"`
	MaxTokens        int                      `json:"max_tokens,omitempty"`
	PresencePenalty  int                      `json:"presence_penalty,omitempty"`
	FrequencyPenalty int                      `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{}   `json:"logit_bias,omitempty"`
	User             string                   `json:"user,omitempty"`
}

type ChatCompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type CreateEditRequest struct {
	Model       string `json:"model"`
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	N           int    `json:"n,omitempty"`
	Temperature int    `json:"temperature,omitempty"`
	TopP        int    `json:"top_p,omitempty"`
}

type CreateImageRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type CreateEmbeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	User  string `json:"user,omitempty"`
}
