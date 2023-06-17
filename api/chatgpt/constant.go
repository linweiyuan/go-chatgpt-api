package chatgpt

//goland:noinspection SpellCheckingInspection
const (
	apiPrefix                       = "https://chat.openai.com/backend-api"
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

	gpt4Model                = "gpt-4"
	gpt4BrowsingModel        = "gpt-4-browsing"
	gpt4PluginsModel         = "gpt-4-plugins"
	gpt4ArkoseTokenPublicKey = "35536E1E-65B4-4D96-9D97-6ADB7EFF8147"
	arkoseTokenTemplate      = `%s.%s|r=us-east-1|meta=3|meta_width=300|metabgclr=transparent|metaiconclr=%%23555555|guitextcolor=%%23000000|pk=%s|at=40|rid=%d|ag=101|cdn_url=https%%3A%%2F%%2Ftcr9i.chat.openai.com%%2Fcdn%%2Ffc|lurl=https%%3A%%2F%%2Faudio-us-east-1.arkoselabs.com|surl=https%%3A%%2F%%2Ftcr9i.chat.openai.com|smurl=https%%3A%%2F%%2Ftcr9i.chat.openai.com%%2Fcdn%%2Ffc%%2Fassets%%2Fstyle-manager`

	actionContinue                     = "continue"
	responseTypeMaxTokens              = "max_tokens"
	responseStatusFinishedSuccessfully = "finished_successfully"
)
