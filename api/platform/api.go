package platform

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-logger/logger"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection GoUnhandledErrorResult
func CreateChatCompletions(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var request struct {
		Stream bool `json:"stream"`
	}
	json.Unmarshal(body, &request)

	url := c.Request.URL.Path
	if strings.Contains(url, "/chat") {
		url = apiCreateChatCompletions
	} else {
		url = apiCreateCompletions
	}

	resp, err := handlePost(c, url, body, request.Stream)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		logger.Error(fmt.Sprintf(api.AccountDeactivatedErrorMessage, c.GetString(api.EmailKey)))
		responseMap := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseMap)
		c.AbortWithStatusJSON(resp.StatusCode, responseMap)
		return
	}

	if request.Stream {
		handleCompletionsResponse(c, resp)
	} else {
		io.Copy(c.Writer, resp.Body)
	}
}

func CreateCompletions(c *gin.Context) {
	CreateChatCompletions(c)
}

//goland:noinspection GoUnhandledErrorResult
func handleCompletionsResponse(c *gin.Context, resp *http.Response) {
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

	reader := bufio.NewReader(resp.Body)
	for {
		if c.Request.Context().Err() != nil {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "" {
			continue
		}

		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()
	}
}

func handlePost(c *gin.Context, url string, data []byte, stream bool) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set(api.AuthorizationHeader, api.GetAccessToken(c))
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return nil, err
	}

	return resp, nil
}
