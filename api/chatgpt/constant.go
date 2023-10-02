package chatgpt

const (
	defaultRole           = "user"
	parseJsonErrorMessage = "failed to parse json request body"

	csrfUrl                  = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl           = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	getCsrfTokenErrorMessage = "failed to get CSRF token"
	authSessionUrl           = "https://chat.openai.com/api/auth/session"

	gpt4Model                          = "gpt-4"
	actionContinue                     = "continue"
	responseTypeMaxTokens              = "max_tokens"
	responseStatusFinishedSuccessfully = "finished_successfully"
	noModelPermissionErrorMessage      = "you have no permission to use this model"
)
