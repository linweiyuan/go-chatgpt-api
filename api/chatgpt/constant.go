package chatgpt

const (
	apiPrefix                       = "https://chat.openai.com/backend-api"
	defaultRole                     = "user"
	getConversationsErrorMessage    = "Failed to get conversations."
	generateTitleErrorMessage       = "Failed to generate title."
	getContentErrorMessage          = "Failed to get content."
	updateConversationErrorMessage  = "Failed to update conversation."
	clearConversationsErrorMessage  = "Failed to clear conversations."
	feedbackMessageErrorMessage     = "Failed to add feedback."
	getModelsErrorMessage           = "Failed to get models."
	getAccountCheckErrorMessage     = "Check failed." // Placeholder. Never encountered.
	parseJsonErrorMessage           = "Failed to parse json request body."
	fallbackErrorMessage            = "Fallback failed."
	fallbackMethodNotAllowedMessage = "Fallback method not allowed."

	csrfUrl                  = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl           = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	getCsrfTokenErrorMessage = "Failed to get CSRF token."
	authSessionUrl           = "https://chat.openai.com/api/auth/session"
)
