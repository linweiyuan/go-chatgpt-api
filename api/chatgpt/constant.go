package chatgpt

//goland:noinspection SpellCheckingInspection
const (
	defaultRole           = "user"
	parseJsonErrorMessage = "Failed to parse json request body."

	csrfUrl                  = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl           = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	getCsrfTokenErrorMessage = "Failed to get CSRF token."
	authSessionUrl           = "https://chat.openai.com/api/auth/session"

	gpt4Model                = "gpt-4"
	gpt4BrowsingModel        = "gpt-4-browsing"
	gpt4PluginsModel         = "gpt-4-plugins"
	gpt4ArkoseTokenPublicKey = "35536E1E-65B4-4D96-9D97-6ADB7EFF8147"
	gpt4TokenUrl             = "https://tcr9i.chat.openai.com/fc/gt2/public_key/" + gpt4ArkoseTokenPublicKey

	actionContinue                     = "continue"
	responseTypeMaxTokens              = "max_tokens"
	responseStatusFinishedSuccessfully = "finished_successfully"
)
