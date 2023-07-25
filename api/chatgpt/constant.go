package chatgpt

//goland:noinspection SpellCheckingInspection
const (
	defaultRole           = "user"
	parseJsonErrorMessage = "Failed to parse json request body."

	csrfUrl                  = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl           = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	getCsrfTokenErrorMessage = "Failed to get CSRF token."
	authSessionUrl           = "https://chat.openai.com/api/auth/session"

	gpt4Model                          = "gpt-4"
	actionContinue                     = "continue"
	responseTypeMaxTokens              = "max_tokens"
	responseStatusFinishedSuccessfully = "finished_successfully"
	noModelPermissionErrorMessage      = "You have no permission to use this model, maybe you Plus has expired, or this model is temporary disabled, or the arkoseToken is invalid."

	refreshPuidErrorMessage = "Failed to refresh PUID."
)
