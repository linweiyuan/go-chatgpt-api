package chatgpt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/util/logger"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
	"github.com/tebeka/selenium"
)

const (
	apiPrefix                      = "https://chat.openai.com/backend-api"
	defaultRole                    = "user"
	getConversationsErrorMessage   = "Failed to get conversations."
	generateTitleErrorMessage      = "Failed to generate title."
	getContentErrorMessage         = "Failed to get content."
	updateConversationErrorMessage = "Failed to update conversation."
	clearConversationsErrorMessage = "Failed to clear conversations."
	feedbackMessageErrorMessage    = "Failed to add feedback."
)

var mutex sync.Mutex

//goland:noinspection GoUnhandledErrorResult
func init() {
	go func() {
		ticker := time.NewTicker(api.RefreshEveryMinutes * time.Minute)

		for {
			select {
			case <-ticker.C:
				func() {
					defer func() {
						if err := recover(); err != nil {
							logger.Error("Failed to refresh page")
							mutex.Unlock()
						}
					}()

					if mutex.TryLock() {
						webdriver.Refresh()
						mutex.Unlock()
					}
				}()
			}
		}
	}()
}

//goland:noinspection GoUnhandledErrorResult
func GetConversations(c *gin.Context) {
	offset, ok := c.GetQuery("offset")
	if !ok {
		offset = "0"
	}
	limit, ok := c.GetQuery("limit")
	if !ok {
		limit = "20"
	}
	url := apiPrefix + "/conversations?offset=" + offset + "&limit=" + limit
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getGetScript(url, accessToken, getConversationsErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == getConversationsErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(getConversationsErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

type StartConversationRequest struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	Model           string    `json:"model"`
	ParentMessageID string    `json:"parent_message_id"`
	ConversationID  *string   `json:"conversation_id"`
}

type Message struct {
	Author  Author  `json:"author"`
	Content Content `json:"content"`
	ID      string  `json:"id"`
}

type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

//goland:noinspection GoUnhandledErrorResult
func StartConversation(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	var request StartConversationRequest
	c.BindJSON(&request)
	if request.ConversationID == nil || *request.ConversationID == "" {
		request.ConversationID = nil
	}
	if request.Messages[0].Author.Role == "" {
		request.Messages[0].Author.Role = defaultRole
	}
	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation"
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScriptForStartConversation(url, accessToken, string(jsonBytes))
	_, err := webdriver.WebDriver.ExecuteScript(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	var callbackChannel = make(chan string)

	go func() {
		webdriver.WebDriver.ExecuteScript("delete window.conversationResponseData;", nil)
		temp := ""
		for {
			conversationResponseData, _ := webdriver.WebDriver.ExecuteScript("return window.conversationResponseData;", nil)
			if conversationResponseData == nil || conversationResponseData == "" {
				continue
			}

			conversationResponseDataString := conversationResponseData.(string)
			conversationResponseDataStrings := strings.Split(conversationResponseDataString, "\n")
			result := api.DoneFlag
			for _, s := range conversationResponseDataStrings {
				s = strings.TrimSpace(s)
				if s != "" && !strings.HasPrefix(s, "event") && !strings.HasPrefix(s, "data: 2023") {
					result = s
				}
			}

			if temp != "" {
				if temp == result {
					continue
				}
			}
			temp = result

			if result == "429" || result == api.DoneFlag {
				callbackChannel <- result
				close(callbackChannel)
				break
			}

			result = result[5:]
			callbackChannel <- result
		}
	}()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for eventDataString := range callbackChannel {
		c.Writer.Write([]byte("data:" + eventDataString + "\n\n"))
		c.Writer.Flush()
	}
}

type GenerateTitleRequest struct {
	MessageID string `json:"message_id"`
	Model     string `json:"model"`
}

//goland:noinspection GoUnhandledErrorResult
func GenerateTitle(c *gin.Context) {
	var request GenerateTitleRequest
	c.BindJSON(&request)
	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation/gen_title/" + c.Param("id")
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScript(url, accessToken, string(jsonBytes), generateTitleErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == generateTitleErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(generateTitleErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

//goland:noinspection GoUnhandledErrorResult
func GetConversation(c *gin.Context) {
	url := apiPrefix + "/conversation/" + c.Param("id")
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getGetScript(url, accessToken, getContentErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == getContentErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(getContentErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

type PatchConversationRequest struct {
	Title     *string `json:"title"`
	IsVisible bool    `json:"is_visible"`
}

//goland:noinspection GoUnhandledErrorResult
func UpdateConversation(c *gin.Context) {
	var request PatchConversationRequest
	c.BindJSON(&request)
	// bool default to false, then will hide (delete) the conversation
	if request.Title != nil {
		request.IsVisible = true
	}
	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation/" + c.Param("id")
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPatchScript(url, accessToken, string(jsonBytes), updateConversationErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == updateConversationErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(updateConversationErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

type FeedbackMessageRequest struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Rating         string `json:"rating"`
}

//goland:noinspection GoUnhandledErrorResult
func FeedbackMessage(c *gin.Context) {
	var request FeedbackMessageRequest
	c.BindJSON(&request)
	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation/message_feedback"
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScript(url, accessToken, string(jsonBytes), feedbackMessageErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == feedbackMessageErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(feedbackMessageErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

//goland:noinspection GoUnhandledErrorResult
func ClearConversations(c *gin.Context) {
	jsonBytes, _ := json.Marshal(PatchConversationRequest{
		IsVisible: false,
	})
	url := apiPrefix + "/conversations"
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPatchScript(url, accessToken, string(jsonBytes), clearConversationsErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == clearConversationsErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(clearConversationsErrorMessage))
		return
	}

	c.Writer.Write([]byte(responseText.(string)))
}

func getGetScript(url string, accessToken string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			headers: {
				'Authorization': '%s'
			}
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, errorMessage)
}

func getPostScriptForStartConversation(url string, accessToken string, jsonString string) string {
	return fmt.Sprintf(`
		let conversationResponseData;
	
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', true);
		xhr.setRequestHeader('Accept', 'text/event-stream');
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onreadystatechange = function() {
			if (xhr.readyState === xhr.LOADING && xhr.status === 200) {
				window.conversationResponseData = xhr.responseText;
			} else if (xhr.status === 429) {
				window.conversationResponseData = '429';
			} if (xhr.readyState === xhr.DONE) {
				window.conversationResponseData = '[DONE]';
			}
		};
		xhr.send(JSON.stringify(%s));
	`, url, accessToken, jsonString)
}

func getPostScript(url string, accessToken string, jsonString string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			method: 'POST',
			headers: {
				'Authorization': '%s',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(%s)
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, jsonString, errorMessage)
}
func getPatchScript(url string, accessToken string, jsonString string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			method: 'PATCH',
			headers: {
				'Authorization': '%s',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(%s)
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('%s');
			}
			return response.text();
		})
		.then(text => {
			arguments[0](text);
		})
		.catch(err => {
			arguments[0](err.message);
		});
	`, url, accessToken, jsonString, errorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handleSeleniumError(err error, script string, c *gin.Context) bool {
	if err != nil {
		if seleniumError, ok := err.(*selenium.Error); ok {
			webdriver.NewSessionAndRefresh(seleniumError.Message)
			responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
			c.Writer.Write([]byte(responseText.(string)))
			return true
		}
	}

	return false
}
