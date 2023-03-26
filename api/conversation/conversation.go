package conversation

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
	"github.com/sirupsen/logrus"
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
	url := "https://chat.openai.com/backend-api/conversations?offset=" + offset + "&limit=" + limit
	accessToken := c.GetHeader(api.Authorization)
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('GET', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.send();
		return xhr.responseText;`, url, accessToken), nil)
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
	Role    string  `json:"role"`
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
	logrus.Info("Start conversations...")

	var request StartConversationRequest
	c.BindJSON(&request)
	if request.ConversationID == nil || *request.ConversationID == "" {
		request.ConversationID = nil
	}
	jsonBytes, _ := json.Marshal(request)
	url := "https://chat.openai.com/backend-api/conversation"
	accessToken := c.GetHeader(api.Authorization)
	webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', true);
		xhr.setRequestHeader('Accept', 'text/event-stream');
		xhr.setRequestHeader('Authorization', 'Bearer %s');
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
		return xhr.responseText;`, url, accessToken, string(jsonBytes)), nil)

	var callbackChannel = make(chan string)

	go func() {
		for {
			eventData, _ := webdriver.WebDriver.ExecuteScriptAsync(`
				const callback = arguments[arguments.length - 1];
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
				window.addEventListener('message', handleFunction);`, nil)

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

	// TODO: have no idea how to handle SSE in Go
	for eventDataString := range callbackChannel {
		c.Writer.Write([]byte("data:" + eventDataString + "\n"))
		c.Writer.Flush()
	}
}

type GenerateTitleRequest struct {
	MessageID string `json:"message_id"`
	Model     string `json:"model"`
}

//goland:noinspection GoUnhandledErrorResult
func GenerateTitle(c *gin.Context) {
	logrus.Info("Generate title...")

	var request GenerateTitleRequest
	c.BindJSON(&request)
	jsonBytes, _ := json.Marshal(request)
	url := "https://chat.openai.com/backend-api/conversation/gen_title/" + c.Param("id")
	accessToken := c.GetHeader(api.Authorization)
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.send('%s');
		return xhr.responseText;`, url, accessToken, string(jsonBytes)), nil)
	c.Writer.Write([]byte(responseText.(string)))
}

//goland:noinspection GoUnhandledErrorResult
func GetConversation(c *gin.Context) {
	logrus.Info("Get conversation...")

	url := "https://chat.openai.com/backend-api/conversation/" + c.Param("id")
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('GET', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.send();
		return xhr.responseText;`, url, accessToken), nil)
	c.Writer.Write([]byte(responseText.(string)))
}

type PatchConversationRequest struct {
	Title     *string `json:"title"`
	IsVisible bool    `json:"is_visible"`
}

//goland:noinspection GoUnhandledErrorResult
func UpdateConversation(c *gin.Context) {
	logrus.Info("Update conversation...")

	var request PatchConversationRequest
	c.BindJSON(&request)
	// bool default to false, then will hide (delete) the conversation
	if request.Title != nil {
		request.IsVisible = true
	}
	jsonBytes, _ := json.Marshal(request)
	url := "https://chat.openai.com/backend-api/conversation/" + c.Param("id")
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('PATCH', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.send('%s');
		return xhr.responseText;`, url, accessToken, string(jsonBytes)), nil)
	c.Writer.Write([]byte(responseText.(string)))
}

type FeedbackMessageRequest struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Rating         string `json:"rating"`
}

//goland:noinspection GoUnhandledErrorResult
func FeedbackMessage(c *gin.Context) {
	logrus.Info("Feedback message...")

	var request FeedbackMessageRequest
	c.BindJSON(&request)
	jsonBytes, _ := json.Marshal(request)
	url := "https://chat.openai.com/backend-api/conversation/message_feedback"
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.send('%s');
		return xhr.responseText;`, url, accessToken, string(jsonBytes)), nil)
	c.Writer.Write([]byte(responseText.(string)))
}

//goland:noinspection GoUnhandledErrorResult
func ClearConversations(c *gin.Context) {
	logrus.Info("Clear conversations...")

	jsonBytes, _ := json.Marshal(PatchConversationRequest{
		IsVisible: false,
	})
	url := "https://chat.openai.com/backend-api/conversations"
	accessToken := c.GetHeader("Authorization")
	responseText, _ := webdriver.WebDriver.ExecuteScript(fmt.Sprintf(`
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '%s', false);
		xhr.setRequestHeader('Authorization', '%s');
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.send('%s');
		return xhr.responseText;`, url, accessToken, string(jsonBytes)), nil)
	c.Writer.Write([]byte(responseText.(string)))
}
