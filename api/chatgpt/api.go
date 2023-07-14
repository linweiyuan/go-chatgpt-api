package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/acheong08/OpenAIAuth/auth"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/funcaptcha"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-logger/logger"

	http "github.com/bogdanfinn/fhttp"
)

//goland:noinspection SpellCheckingInspection
var (
	arkoseTokenUrl string
	puid           string
	bx             string
)

//goland:noinspection SpellCheckingInspection
func init() {
	arkoseTokenUrl = os.Getenv("GO_CHATGPT_API_ARKOSE_TOKEN_URL")
	bx = os.Getenv("GO_CHATGPT_API_BX")

	setupPUID()
}

//goland:noinspection GoUnhandledErrorResult
func CreateConversation(c *gin.Context) {
	var request CreateConversationRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(parseJsonErrorMessage))
		return
	}

	if request.ConversationID == nil || *request.ConversationID == "" {
		request.ConversationID = nil
	}

	if len(request.Messages) != 0 {
		if request.Messages[0].Author.Role == "" {
			request.Messages[0].Author.Role = defaultRole
		}
	}

	logger.Info(request.Model)

	if strings.HasPrefix(request.Model, gpt4Model) {
		arkoseToken, err := GetArkoseToken()
		if err != nil || arkoseToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, api.ReturnMessage(err.Error()))
			return
		}

		request.ArkoseToken = arkoseToken
	}

	resp, done := sendConversationRequest(c, request)
	if done {
		return
	}

	handleConversationResponse(c, resp, request)
}

//goland:noinspection GoUnhandledErrorResult
func sendConversationRequest(c *gin.Context, request CreateConversationRequest) (*http.Response, bool) {
	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, api.ChatGPTApiUrlPrefix+"/backend-api/conversation", bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Authorization", api.GetAccessTokenFromHeader(c.Request.Header))
	req.Header.Set("Accept", "text/event-stream")
	if puid != "" {
		//goland:noinspection SpellCheckingInspection
		req.Header.Set("Cookie", "_puid="+puid)
	}
	resp, err := api.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return nil, true
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		req, _ := http.NewRequest(http.MethodGet, api.ChatGPTApiUrlPrefix+"/backend-api/models?history_and_training_disabled=false", nil)
		req.Header.Set("User-Agent", api.UserAgent)
		req.Header.Set("Authorization", api.GetAccessTokenFromHeader(c.Request.Header))
		response, err := api.Client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
			return nil, true
		}

		defer response.Body.Close()
		modelAvailable := false
		var getModelsResponse GetModelsResponse
		json.NewDecoder(response.Body).Decode(&getModelsResponse)
		for _, model := range getModelsResponse.Models {
			if model.Slug == request.Model {
				modelAvailable = true
				break
			}
		}
		if !modelAvailable {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"detail": noModelPermissionErrorMessage,
			})
			return nil, true
		}

		responseMap := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseMap)
		c.AbortWithStatusJSON(resp.StatusCode, responseMap)
		return nil, true
	}

	return resp, false
}

//goland:noinspection GoUnhandledErrorResult
func handleConversationResponse(c *gin.Context, resp *http.Response, request CreateConversationRequest) {
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

	isMaxTokens := false
	continueParentMessageID := ""
	continueConversationID := ""

	defer resp.Body.Close()
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

		responseJson := line[6:]
		if strings.HasPrefix(responseJson, "[DONE]") && isMaxTokens && request.AutoContinue {
			continue
		}

		// no need to unmarshal every time, but if response content has this "max_tokens", need to further check
		if strings.TrimSpace(responseJson) != "" && strings.Contains(responseJson, responseTypeMaxTokens) {
			var createConversationResponse CreateConversationResponse
			json.Unmarshal([]byte(responseJson), &createConversationResponse)
			message := createConversationResponse.Message
			if message.Metadata.FinishDetails.Type == responseTypeMaxTokens && createConversationResponse.Message.Status == responseStatusFinishedSuccessfully {
				isMaxTokens = true
				continueParentMessageID = message.ID
				continueConversationID = createConversationResponse.ConversationID
			}
		}

		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()
	}

	if isMaxTokens && request.AutoContinue {
		continueConversationRequest := CreateConversationRequest{
			ArkoseToken:                request.ArkoseToken,
			HistoryAndTrainingDisabled: request.HistoryAndTrainingDisabled,
			Model:                      request.Model,
			TimezoneOffsetMin:          request.TimezoneOffsetMin,

			Action:          actionContinue,
			ParentMessageID: continueParentMessageID,
			ConversationID:  &continueConversationID,
		}
		resp, done := sendConversationRequest(c, continueConversationRequest)
		if done {
			return
		}

		handleConversationResponse(c, resp, continueConversationRequest)
	}
}

//goland:noinspection SpellCheckingInspection
func setupPUID() {
	username := os.Getenv("GO_CHATGPT_API_OPENAI_EMAIL")
	password := os.Getenv("GO_CHATGPT_API_OPENAI_PASSWORD")
	if username != "" && password != "" {
		go func() {
			for {
				authenticator := auth.NewAuthenticator(username, password, os.Getenv("GO_CHATGPT_API_PROXY"))
				err := authenticator.Begin()
				if err != nil {
					logger.Error("Failed to login to get PUID.")
					break
				}

				puid, err = authenticator.GetPUID()
				if err != nil {
					logger.Error("Failed to get PUID.")
					break
				}
				time.Sleep(24 * time.Hour * 7)
			}
		}()
	}
}

//goland:noinspection GoUnhandledErrorResult
func GetArkoseToken() (string, error) {
	if arkoseTokenUrl == "" {
		return funcaptcha.GetOpenAITokenWithBx(bx)
	}

	req, _ := http.NewRequest(http.MethodGet, arkoseTokenUrl, nil)
	resp, err := api.Client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	defer resp.Body.Close()
	responseMap := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseMap)
	arkoseToken, ok := responseMap["token"]
	if !ok || arkoseToken == "" {
		return "", nil
	}

	return arkoseToken.(string), nil
}
