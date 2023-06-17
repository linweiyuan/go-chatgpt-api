package chatgpt

//goland:noinspection GoSnakeCaseUsage
import tls_client "github.com/bogdanfinn/tls-client"

type UserLogin struct {
	client tls_client.HttpClient
}

type CreateConversationRequest struct {
	Action                     string     `json:"action"`
	Messages                   *[]Message `json:"messages"`
	Model                      string     `json:"model"`
	ParentMessageID            string     `json:"parent_message_id"`
	ConversationID             *string    `json:"conversation_id"`
	PluginIDs                  []string   `json:"plugin_ids"`
	TimezoneOffsetMin          int        `json:"timezone_offset_min"`
	ArkoseToken                string     `json:"arkose_token"`
	HistoryAndTrainingDisabled bool       `json:"history_and_training_disabled"`
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

type CreateConversationResponse struct {
	Message struct {
		ID     string `json:"id"`
		Author struct {
			Role     string      `json:"role"`
			Name     interface{} `json:"name"`
			Metadata struct {
			} `json:"metadata"`
		} `json:"author"`
		CreateTime float64     `json:"create_time"`
		UpdateTime interface{} `json:"update_time"`
		Content    struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
		Status   string  `json:"status"`
		EndTurn  bool    `json:"end_turn"`
		Weight   float64 `json:"weight"`
		Metadata struct {
			MessageType   string `json:"message_type"`
			ModelSlug     string `json:"model_slug"`
			FinishDetails struct {
				Type string `json:"type"`
			} `json:"finish_details"`
		} `json:"metadata"`
		Recipient string `json:"recipient"`
	} `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}

type FeedbackMessageRequest struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Rating         string `json:"rating"`
}

type GenerateTitleRequest struct {
	MessageID string `json:"message_id"`
}

type PatchConversationRequest struct {
	Title     *string `json:"title"`
	IsVisible bool    `json:"is_visible"`
}

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Expiry int64  `json:"expiry"`
}
