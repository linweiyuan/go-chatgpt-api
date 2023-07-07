package imitate

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"

	http "github.com/bogdanfinn/fhttp"
)

type ContinueInfo struct {
	ConversationID string `json:"conversation_id"`
	ParentID       string `json:"parent_id"`
}

type APIRequest struct {
	Messages  []ApiMessage `json:"messages"`
	Stream    bool         `json:"stream"`
	Model     string       `json:"model"`
	PluginIDs []string     `json:"plugin_ids"`
}

type ApiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func HandleRequestError(c *gin.Context, response *http.Response) bool {
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			// Read response body
			body, _ := io.ReadAll(response.Body)
			c.JSON(500, gin.H{"error": gin.H{
				"message": "Unknown error",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "500",
				"details": string(body),
			}})
			return true
		}
		c.JSON(response.StatusCode, gin.H{"error": gin.H{
			"message": errorResponse["detail"],
			"type":    response.Status,
			"param":   nil,
			"code":    "error",
		}})
		return true
	}
	return false
}
