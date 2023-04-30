package chatgpt

const (
	apiPrefix                      = "https://chat.openai.com/backend-api"
	defaultRole                    = "user"
	getConversationsErrorMessage   = "Failed to get conversations."
	generateTitleErrorMessage      = "Failed to generate title."
	getContentErrorMessage         = "Failed to get content."
	updateConversationErrorMessage = "Failed to update conversation."
	clearConversationsErrorMessage = "Failed to clear conversations."
	feedbackMessageErrorMessage    = "Failed to add feedback."
	getModelsErrorMessage          = "Failed to get models."
	getAccountCheckErrorMessage    = "Check failed." // Placeholder. Never encountered.
	parseJsonErrorMessage          = "Failed to parse json request body."

	contentType                        = "application/x-www-form-urlencoded"
	userAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	heartBeatUrl                       = "https://chat.openai.com"
	csrfUrl                            = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl                     = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	loginUsernameUrl                   = "https://auth0.openai.com/u/login/identifier?state="
	loginPasswordUrl                   = "https://auth0.openai.com/u/login/password?state="
	authSessionUrl                     = "https://chat.openai.com/api/auth/session"
	parseUserInfoErrorMessage          = "Failed to parse user login info."
	getCsrfTokenErrorMessage           = "Failed to get CSRF token."
	getAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	getStateErrorMessage               = "Failed to get state."
	emailInvalidErrorMessage           = "Email is not valid."
	emailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	getAccessTokenErrorMessage         = "Failed to get access token, please try again later."

	conversationErrorMessage401 = "Access token has expired."
	conversationErrorMessage403 = "Something went wrong, please try again."
	conversationErrorMessage404 = "The requested conversation cannot be found."
	conversationErrorMessage413 = "The question is too large to handle."
	conversationErrorMessage422 = "The request body is invalid."
	conversationErrorMessage429 = "Too many requests, please try again later."
	conversationErrorMessage500 = "Server error, please try again."
)
