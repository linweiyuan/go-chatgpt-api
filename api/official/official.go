package official

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	apiUrl = "https://api.openai.com"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 0,
	}
}

type ChatCompletionsRequest struct {
	Model    string                   `json:"model"`
	Messages []ChatCompletionsMessage `json:"messages"`
	Stream   bool                     `json:"stream"`
}

type ChatCompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

//goland:noinspection GoUnhandledErrorResult
func ChatCompletions(c *gin.Context) {
	var chatCompletionsRequest ChatCompletionsRequest
	c.ShouldBindJSON(&chatCompletionsRequest)
	data, _ := json.Marshal(chatCompletionsRequest)
	req, _ := http.NewRequest("POST", apiUrl+"/v1/chat/completions", bytes.NewBuffer(data))
	req.Header.Set("Authorization", getHeader(c.GetHeader("Authorization")))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		} else {
			c.Writer.Write([]byte(line))
			c.Writer.Flush()
		}
	}
}

func getHeader(header string) string {
	if !strings.HasPrefix("Bearer", header) {
		return "Bearer " + header
	}
	return header
}
