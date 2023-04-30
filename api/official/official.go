package official

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult
func ChatCompletions(c *gin.Context) {
	var chatCompletionsRequest ChatCompletionsRequest
	c.ShouldBindJSON(&chatCompletionsRequest)
	data, _ := json.Marshal(chatCompletionsRequest)
	req, _ := http.NewRequest("POST", apiChatCompletions, bytes.NewBuffer(data))
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := api.Client.Do(req)
	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func CheckUsage(c *gin.Context) {
	req, _ := http.NewRequest("GET", apiCheckUsage, nil)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	resp, _ := api.Client.Do(req)
	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}
