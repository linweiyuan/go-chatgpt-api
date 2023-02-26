package conversation

import (
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
