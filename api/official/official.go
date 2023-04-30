package official

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"

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
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "\r\n" {
			continue
		}
		if err != nil {
			break
		} else {
			c.Writer.Write([]byte(line))
			c.Writer.Flush()
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func CheckUsage(c *gin.Context) {
	req, _ := http.NewRequest("GET", apiCheckUsage, nil)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	resp, _ := api.Client.Do(req)
	defer resp.Body.Close()
	io.Copy(c.Writer, resp.Body)
}
