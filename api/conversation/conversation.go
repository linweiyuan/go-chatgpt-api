package conversation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
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
	accessToken := c.GetHeader(api.Authorization)
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getGetScript(url, accessToken, getConversationsErrorMessage), nil)
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
	accessToken := c.GetHeader(api.Authorization)
	webdriver.WebDriver.ExecuteScript(getPostScriptForStartConversation(url, accessToken, string(jsonBytes)), nil)

	var callbackChannel = make(chan string)

	go func() {
		for {
			eventData, _ := webdriver.WebDriver.ExecuteScriptAsync(getCallbackScriptForStartConversation(), nil)

			// sometimes callback will not return the final data
			if eventData == nil {
				callbackChannel <- api.DoneFlag
				close(callbackChannel)
				break
			}

			eventDataString := eventData.(string)
			if eventDataString == "429" || eventDataString == api.DoneFlag {
				callbackChannel <- eventDataString
				close(callbackChannel)
				break
			}

			callbackChannel <- eventDataString
		}
	}()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for eventDataString := range callbackChannel {
		c.Writer.Write([]byte("data:" + eventDataString + "\n\n"))
		c.Writer.Flush()
	}

	c.Writer.Write([]byte("event: close\ndata: close\n\n"))
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
	accessToken := c.GetHeader(api.Authorization)
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getPostScript(url, accessToken, string(jsonBytes), generateTitleErrorMessage), nil)
	if responseText == generateTitleErrorMessage {
		c.JSON(http.StatusInternalServerError, api.ReturnMessage(generateTitleErrorMessage))
		return
	}
	c.Writer.Write([]byte(responseText.(string)))
}

//goland:noinspection GoUnhandledErrorResult
func GetConversation(c *gin.Context) {
	url := apiPrefix + "/conversation/" + c.Param("id")
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getGetScript(url, accessToken, getContentErrorMessage), nil)
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
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getPatchScript(url, accessToken, string(jsonBytes), updateConversationErrorMessage), nil)
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
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getPostScript(url, accessToken, string(jsonBytes), feedbackMessageErrorMessage), nil)
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
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(getPatchScript(url, accessToken, string(jsonBytes), clearConversationsErrorMessage), nil)
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
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', true);
		xhr.setRequestHeader('Accept', 'text/event-stream');
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.onreadystatechange = function() {
			if (xhr.readyState === xhr.LOADING && xhr.status === 200) {
				window.postMessage(xhr.responseText);
			} else if (xhr.status === 429) {
				window.postMessage("429");
			} if (xhr.readyState === xhr.DONE) {

			}
		};
		xhr.send('%s');
		return xhr.responseText;
	`, url, accessToken, jsonString)
}

func getCallbackScriptForStartConversation() string {
	return `
		const callback = arguments[0];
		const handleFunction = function(event) {
			const list = event.data.split('\n\n');
			list.pop();
			const eventData = list.pop();
			if (eventData.startsWith('event')) {
				callback(eventData.substring(55));
			} else {
				callback(eventData.substring(6));
			}
		};
		window.removeEventListener('message', handleFunction);
		window.addEventListener('message', handleFunction);
	`
}

func getPostScript(url string, accessToken string, jsonString string, errorMessage string) string {
	return fmt.Sprintf(`
		fetch('%s', {
			method: 'POST',
			headers: {
				'Authorization': '%s',
				'Content-Type': 'application/json'
			},
			body: '%s'
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
			body: '%s'
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
