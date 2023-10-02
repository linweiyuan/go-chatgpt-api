package imitate

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-logger/logger"
)

var (
	reg *regexp.Regexp
)

func init() {
	reg, _ = regexp.Compile("[^a-zA-Z0-9]+")
}

func CreateChatCompletions(c *gin.Context) {
	var originalRequest APIRequest
	err := c.BindJSON(&originalRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Request must be proper JSON",
			"type":    "invalid_request_error",
			"param":   nil,
			"code":    err.Error(),
		}})
		return
	}

	authHeader := c.GetHeader(api.AuthorizationHeader)
	token := os.Getenv("IMITATE_ACCESS_TOKEN")
	if authHeader != "" {
		customAccessToken := strings.Replace(authHeader, "Bearer ", "", 1)
		// Check if customAccessToken starts with sk-
		if strings.HasPrefix(customAccessToken, "eyJhbGciOiJSUzI1NiI") {
			token = customAccessToken
		}
	}

	// 将聊天请求转换为ChatGPT请求。
	translatedRequest, model := convertAPIRequest(originalRequest)

	response, done := sendConversationRequest(c, translatedRequest, token)
	if done {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if HandleRequestError(c, response) {
		return
	}

	var fullResponse string

	id := generateId()

	for i := 3; i > 0; i-- {
		var continueInfo *ContinueInfo
		var responsePart string
		var continueSignal string
		responsePart, continueInfo = Handler(c, response, originalRequest.Stream, id, model)
		fullResponse += responsePart
		continueSignal = os.Getenv("CONTINUE_SIGNAL")
		if continueInfo == nil || continueSignal == "" {
			break
		}
		println("Continuing conversation")
		translatedRequest.Messages = nil
		translatedRequest.Action = "continue"
		translatedRequest.ConversationID = &continueInfo.ConversationID
		translatedRequest.ParentMessageID = continueInfo.ParentID
		response, done = sendConversationRequest(c, translatedRequest, token)

		if done {
			return
		}

		// 以下修复代码来自ChatGPT
		// 在循环内部创建一个局部作用域，并将资源的引用传递给匿名函数，保证资源将在每次迭代结束时被正确释放
		func() {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(response.Body)
		}()

		if HandleRequestError(c, response) {
			return
		}
	}

	if !originalRequest.Stream {
		c.JSON(200, newChatCompletion(fullResponse, model, id))
	} else {
		c.String(200, "data: [DONE]\n\n")
	}
}

func generateId() string {
	id := uuid.NewString()
	id = strings.ReplaceAll(id, "-", "")
	id = base64.StdEncoding.EncodeToString([]byte(id))
	id = reg.ReplaceAllString(id, "")
	return "chatcmpl-" + id
}

func convertAPIRequest(apiRequest APIRequest) (chatgpt.CreateConversationRequest, string) {
	chatgptRequest := NewChatGPTRequest()

	var model = "gpt-3.5-turbo-0613"

	if strings.HasPrefix(apiRequest.Model, "gpt-3.5") {
		chatgptRequest.Model = "text-davinci-002-render-sha"
	}

	if strings.HasPrefix(apiRequest.Model, "gpt-4") {
		arkoseToken, err := api.GetArkoseToken()
		if err == nil {
			chatgptRequest.ArkoseToken = arkoseToken
		} else {
			fmt.Println("Error getting Arkose token: ", err)
		}
		chatgptRequest.Model = apiRequest.Model
		model = "gpt-4-0613"
	}

	if apiRequest.PluginIDs != nil {
		chatgptRequest.PluginIDs = apiRequest.PluginIDs
		chatgptRequest.Model = "gpt-4-plugins"
	}

	for _, apiMessage := range apiRequest.Messages {
		if apiMessage.Role == "system" {
			apiMessage.Role = "critic"
		}
		chatgptRequest.AddMessage(apiMessage.Role, apiMessage.Content)
	}

	return chatgptRequest, model
}

func NewChatGPTRequest() chatgpt.CreateConversationRequest {
	enableHistory := os.Getenv("ENABLE_HISTORY") == ""
	return chatgpt.CreateConversationRequest{
		Action:                     "next",
		ParentMessageID:            uuid.NewString(),
		Model:                      "text-davinci-002-render-sha",
		HistoryAndTrainingDisabled: !enableHistory,
	}
}

func sendConversationRequest(c *gin.Context, request chatgpt.CreateConversationRequest, accessToken string) (*http.Response, bool) {
	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, api.ChatGPTApiUrlPrefix+"/backend-api/conversation", bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set(api.AuthorizationHeader, accessToken)
	req.Header.Set("Accept", "text/event-stream")
	if api.PUID != "" {
		req.Header.Set("Cookie", "_puid="+api.PUID)
	}
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return nil, true
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			logger.Error(fmt.Sprintf(api.AccountDeactivatedErrorMessage, c.GetString(api.EmailKey)))
		}

		responseMap := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseMap)
		c.AbortWithStatusJSON(resp.StatusCode, responseMap)
		return nil, true
	}

	return resp, false
}

func Handler(c *gin.Context, response *http.Response, stream bool, id string, model string) (string, *ContinueInfo) {
	maxTokens := false

	// Create a bufio.Reader from the response body
	reader := bufio.NewReader(response.Body)

	// Read the response byte by byte until a newline character is encountered
	if stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	var finishReason string
	var previousText StringStruct
	var originalResponse ChatGPTResponse
	var isRole = true
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil
		}
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		// Check if line starts with [DONE]
		if !strings.HasPrefix(line, "[DONE]") {
			// Parse the line as JSON

			err = json.Unmarshal([]byte(line), &originalResponse)
			if err != nil {
				continue
			}
			if originalResponse.Error != nil {
				c.JSON(500, gin.H{"error": originalResponse.Error})
				return "", nil
			}
			if originalResponse.Message.Author.Role != "assistant" || originalResponse.Message.Content.Parts == nil {
				continue
			}
			if originalResponse.Message.Metadata.MessageType != "next" && originalResponse.Message.Metadata.MessageType != "continue" || originalResponse.Message.EndTurn != nil {
				continue
			}
			if (len(originalResponse.Message.Content.Parts) == 0 || originalResponse.Message.Content.Parts[0] == "") && !isRole {
				continue
			}
			responseString := ConvertToString(&originalResponse, &previousText, isRole, id, model)
			isRole = false
			if stream {
				_, err = c.Writer.WriteString(responseString)
				if err != nil {
					return "", nil
				}
			}
			// Flush the response writer buffer to ensure that the client receives each line as it's written
			c.Writer.Flush()

			if originalResponse.Message.Metadata.FinishDetails != nil {
				if originalResponse.Message.Metadata.FinishDetails.Type == "max_tokens" {
					maxTokens = true
				}
				finishReason = originalResponse.Message.Metadata.FinishDetails.Type
			}

		} else {
			if stream {
				if finishReason == "" {
					finishReason = "stop"
				}
				finalLine := StopChunk(finishReason, id, model)
				_, err := c.Writer.WriteString("data: " + finalLine.String() + "\n\n")
				if err != nil {
					return "", nil
				}
			}
		}
	}
	if !maxTokens {
		return previousText.Text, nil
	}
	return previousText.Text, &ContinueInfo{
		ConversationID: originalResponse.ConversationID,
		ParentID:       originalResponse.Message.ID,
	}
}
