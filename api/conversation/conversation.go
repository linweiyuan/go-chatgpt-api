package conversation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 0,
	}
}

func GetConversations(c *gin.Context) {
	req, _ := http.NewRequest("GET", "https://apps.openai.com/api/conversations?offset=0&limit=100", nil)
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get conversations."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write([]byte(body))
}

type Conversation struct {
	Action          string    `json:"action"`
	ConversationID  *string   `json:"conversation_id"`
	Messages        []Message `json:"messages"`
	Model           string    `json:"model"`
	ParentMessageID string    `json:"parent_message_id"`
}

type Message struct {
	Author  Author  `json:"author"`
	Content Content `json:"content"`
	ID      string  `json:"id"`
	Role    string  `json:"role"`
}

type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type MakeConversationRequest struct {
	MessageId       string  `json:"message_id"`
	ParentMessageId string  `json:"parent_message_id"`
	ConversationId  *string `json:"conversation_id"`
	Content         string  `json:"content"`
}

func MakeConversation(c *gin.Context) {
	var makeConversationRequest MakeConversationRequest
	if err := c.ShouldBindJSON(&makeConversationRequest); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse make conversation request."))
		return
	}

	jsonBytes, _ := json.Marshal(Conversation{
		Action:         "next",
		ConversationID: makeConversationRequest.ConversationId,
		Messages: []Message{
			{
				Author: Author{
					Role: "user",
				},
				Content: Content{
					ContentType: "text",
					Parts:       []string{makeConversationRequest.Content},
				},
				ID:   makeConversationRequest.MessageId,
				Role: "user",
			},
		},
		Model:           "text-davinci-002-render-sha",
		ParentMessageID: makeConversationRequest.ParentMessageId,
	})
	req, _ := http.NewRequest("POST", "https://apps.openai.com/api/conversation", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))
	req.Header.Set("Accept", "text/event-stream")

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		c.JSON(resp.StatusCode, api.ReturnMessage("Too many requests in 1 hour, please try again later."))
		return
	}

	if resp.StatusCode == http.StatusInternalServerError {
		c.JSON(resp.StatusCode, api.ReturnMessage("Server error, please try again later."))
		return
	}

	io.Copy(c.Writer, resp.Body)
}

func GenConversationTitle(c *gin.Context) {
	var request struct {
		MessageId string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse gen conversation title request."))
		return
	}

	jsonBytes, _ := json.Marshal(map[string]string{
		"message_id": request.MessageId,
		"model":      "text-davinci-002-render-sha",
	})
	req, _ := http.NewRequest("POST", "https://apps.openai.com/api/conversation/gen_title/"+c.Param("id"), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to gen conversation title."))
		return
	}

	io.Copy(c.Writer, resp.Body)
}

func GetConversation(c *gin.Context) {
	req, _ := http.NewRequest("GET", "https://apps.openai.com/api/conversation/"+c.Param("id"), nil)
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to get conversation."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write([]byte(body))
}
