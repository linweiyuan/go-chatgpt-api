package chatgpt

type StartConversationRequest struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	Model           string    `json:"model"`
	ParentMessageID string    `json:"parent_message_id"`
	ConversationID  *string   `json:"conversation_id"`
	ContinueText    string    `json:"continue_text"`
}

type Message struct {
	Author  Author  `json:"author"`
	Content Content `json:"content"`
	ID      string  `json:"id"`
}

type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ConversationResponse struct {
	Message struct {
		ID      string `json:"id"`
		Content struct {
			Parts []string `json:"parts"`
		} `json:"content"`
		EndTurn  bool `json:"end_turn"`
		Metadata struct {
			FinishDetails struct {
				Type string `json:"type"`
			} `json:"finish_details"`
		} `json:"metadata"`
	} `json:"message"`
	ConversationID string `json:"conversation_id"`
}

type FeedbackMessageRequest struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Rating         string `json:"rating"`
}

type GenerateTitleRequest struct {
	MessageID string `json:"message_id"`
	Model     string `json:"model"`
}

type PatchConversationRequest struct {
	Title     *string `json:"title"`
	IsVisible bool    `json:"is_visible"`
}
