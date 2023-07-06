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
	builtInArkoseTokenIllegalMessage   = "arkose token is illegal, try using an external arkose token service"
	ArkoseTokenIllegalMessage          = "arkose token is illegal,the arkose service url :"
	getArkoseTokenErrorMessage         = "Failed to get arkose token."
	noModelPermissionErrorMessage      = "You have no permission to use this model, maybe you Plus has expired, or this model is temporary disabled."
)
