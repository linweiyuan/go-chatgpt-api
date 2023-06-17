package chatgpt

//goland:noinspection GoSnakeCaseUsage
import tls_client "github.com/bogdanfinn/tls-client"

type UserLogin struct {
	client tls_client.HttpClient
}

type CreateConversationRequest struct {
	Action            string    `json:"action"`
	Messages          []Message `json:"messages"`
	Model             string    `json:"model"`
	ParentMessageID   string    `json:"parent_message_id"`
	ConversationID    *string   `json:"conversation_id"`
	PluginIDs         []string  `json:"plugin_ids"`
	TimezoneOffsetMin int       `json:"timezone_offset_min"`
	ArkoseToken       string    `json:"arkose_token"`
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
