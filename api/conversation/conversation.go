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
	MessageID       string `json:"message_id"`
	ParentMessageID string `json:"parent_message_id"`
	ConversationID  string `json:"conversation_id"`
	Content         string `json:"content"`
}

func MakeConversation(c *gin.Context) {
	var request MakeConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse make conversation request."))
		return
	}

	conversation := Conversation{
		Action: "next",
		Messages: []Message{
			{
				Author: Author{
					Role: "user",
				},
				Content: Content{
					ContentType: "text",
					Parts:       []string{request.Content},
				},
				ID:   request.MessageID,
				Role: "user",
			},
		},
		Model:           "text-davinci-002-render-sha",
		ParentMessageID: request.ParentMessageID,
	}

	if request.ConversationID != "" {
		conversation.ConversationID = &request.ConversationID
	}

	jsonBytes, _ := json.Marshal(conversation)
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
		MessageID string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse gen conversation title request."))
		return
	}

	jsonBytes, _ := json.Marshal(map[string]string{
		"message_id": request.MessageID,
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

type PatchConversationRequest struct {
	Title     *string `json:"title"`
	IsVisible bool    `json:"is_visible"`
}

func PatchConversation(c *gin.Context) {
	var request PatchConversationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse update conversation request."))
		return
	}

	// bool default to false, then will hide (delete) the conversation
	if request.Title != nil {
		request.IsVisible = true
	}

	conversationID := c.Param("id")

	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "https://apps.openai.com/api/conversation/"+conversationID, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to update conversation."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write([]byte(body))
}

type FeedbackMessageRequest struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Rating         string `json:"rating"`
}

func FeedbackMessage(c *gin.Context) {
	var request FeedbackMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage("Failed to parse feedback conversation request."))
		return
	}

	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "https://apps.openai.com/api/conversation/message_feedback", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		c.JSON(resp.StatusCode, api.ReturnMessage("Your how selected another one before."))
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to make a message feedback."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write([]byte(body))
}

func ClearConversations(c *gin.Context) {
	jsonBytes, _ := json.Marshal(PatchConversationRequest{
		IsVisible: false,
	})
	req, _ := http.NewRequest("PATCH", "https://apps.openai.com/api/conversations", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+c.GetHeader("Authorization"))

	resp, err := client.Do(req)
	api.CheckError(c, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, api.ReturnMessage("Failed to clear conversations."))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write([]byte(body))
}
