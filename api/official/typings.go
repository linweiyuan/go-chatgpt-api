package official

type ChatCompletionsRequest struct {
	Model    string                   `json:"model"`
	Messages []ChatCompletionsMessage `json:"messages"`
	Stream   bool                     `json:"stream"`
}

type ChatCompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
