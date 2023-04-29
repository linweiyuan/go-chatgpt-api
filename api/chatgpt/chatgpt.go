package chatgpt

import (
	"encoding/json"
	"fmt"

	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/webdriver"
	"github.com/tebeka/selenium"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const (
	defaultApiTimeoutSeconds       = 30
	apiPrefix                      = "https://chat.openai.com/backend-api"
	defaultRole                    = "user"
	getConversationsErrorMessage   = "Failed to get conversations."
	generateTitleErrorMessage      = "Failed to generate title."
	getContentErrorMessage         = "Failed to get content."
	updateConversationErrorMessage = "Failed to update conversation."
	clearConversationsErrorMessage = "Failed to clear conversations."
	feedbackMessageErrorMessage    = "Failed to add feedback."
	getModelsErrorMessage          = "Failed to get models."
	getAccountCheckErrorMessage    = "Check failed." // Placeholder. Never encountered.
	parseJsonErrorMessage          = "Failed to parse json request body."
	doneFlag                       = "[DONE]"

	contentType                        = "application/x-www-form-urlencoded"
	userAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	csrfUrl                            = "https://chat.openai.com/api/auth/csrf"
	promptLoginUrl                     = "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	loginUsernameUrl                   = "https://auth0.openai.com/u/login/identifier?state="
	loginPasswordUrl                   = "https://auth0.openai.com/u/login/password?state="
	authSessionUrl                     = "https://chat.openai.com/api/auth/session"
	parseUserInfoErrorMessage          = "Failed to parse user login info."
	getCsrfTokenErrorMessage           = "Failed to get CSRF token."
	getAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	getStateErrorMessage               = "Failed to get state."
	emailInvalidErrorMessage           = "Email is not valid."
	emailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	getAccessTokenErrorMessage         = "Failed to get access token, please try again later."
)

var client tls_client.HttpClient

//goland:noinspection GoUnhandledErrorResult
func init() {
	go func() {
		ticker := time.NewTicker(api.RefreshEveryMinutes * time.Minute)
		for {
			select {
			case <-ticker.C:
				tryToRefreshPage()
			}
		}
	}()

	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(defaultApiTimeoutSeconds),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}...)
}

//goland:noinspection GoUnhandledErrorResult
func tryToRefreshPage() {
	tabs, _ := webdriver.WebDriver.WindowHandles()
	if len(tabs) < 2 {
		webdriver.OpenNewTabAndChangeBackToOldTab()
		tabs, _ = webdriver.WebDriver.WindowHandles()
	}
	webdriver.WebDriver.SwitchWindow(tabs[1]) // new tab to refresh cookies
	webdriver.Refresh()
	webdriver.WebDriver.SwitchWindow(tabs[0]) // old tab for API handling
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
	handleGet(c, url, getConversationsErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func StartConversation(c *gin.Context) {
	xhrMap, _ := webdriver.WebDriver.ExecuteScript("return window.xhrMap;", nil)
	if xhrMap == nil {
		webdriver.InitXhrMap()
	}
	conversationMap, _ := webdriver.WebDriver.ExecuteScript("return window.conversationMap;", nil)
	if conversationMap == nil {
		webdriver.InitConversationMap()
	}

	var callbackChannel = make(chan string)

	var request StartConversationRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(parseJsonErrorMessage))
		return
	}

	if request.ConversationID == nil || *request.ConversationID == "" {
		request.ConversationID = nil
	}
	if request.Messages[0].Author.Role == "" {
		request.Messages[0].Author.Role = defaultRole
	}
	if request.VariantPurpose == "" {
		request.VariantPurpose = "none"
	}

	oldContentToResponse := ""
	messageID := request.Messages[0].ID
	if !sendConversationRequest(c, callbackChannel, request, oldContentToResponse, messageID) {
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for eventDataString := range callbackChannel {
		c.Writer.Write([]byte("data: " + eventDataString + "\n\n"))
		c.Writer.Flush()
	}
}

//goland:noinspection GoUnhandledErrorResult
func sendConversationRequest(c *gin.Context, callbackChannel chan string, request StartConversationRequest, oldContent string, messageID string) bool {
	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation"
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScriptForStartConversation(url, accessToken, string(jsonBytes), messageID)
	_, err := webdriver.WebDriver.ExecuteScript(script, nil)
	if handleSeleniumError(err, script, c) {
		return false
	}

	go func() {
		defer func() {
			webdriver.WebDriver.ExecuteScript(fmt.Sprintf("conversationMap.delete('%s');xhrMap.delete('%s');", messageID, messageID), nil)
		}()

		// temp for performance optimisation
		temp := ""
		var conversationResponse ConversationResponse
		maxTokens := false
		for {
			if c.Request.Context().Err() != nil {
				stopGenerate(messageID)
				close(callbackChannel)
				break
			}

			conversationResponseData, err := webdriver.WebDriver.ExecuteScript(fmt.Sprintf("return conversationMap.get('%s');", messageID), nil)
			if err != nil {
				if strings.Contains(err.Error(), "conversationMap is not defined") {
					webdriver.InitConversationMap()
					continue
				}
			}

			if conversationResponseData == nil || conversationResponseData == "" {
				time.Sleep(api.IdleTimeMilliseconds * time.Millisecond)
				continue
			}

			conversationResponseDataString := conversationResponseData.(string)
			if conversationResponseDataString[0:1] == strconv.Itoa(4) || conversationResponseDataString[0:1] == strconv.Itoa(5) {
				statusCode, _ := strconv.Atoi(conversationResponseDataString[0:3])
				if statusCode == http.StatusForbidden {
					webdriver.Refresh()
				}
				c.AbortWithStatusJSON(statusCode, api.ReturnMessage(conversationResponseDataString[3:]))
				close(callbackChannel)
				break
			}

			if conversationResponseDataString[0:1] == "!" {
				callbackChannel <- conversationResponseDataString[1:]
				callbackChannel <- doneFlag
				close(callbackChannel)
				break
			}

			if temp != "" {
				if temp == conversationResponseDataString {
					time.Sleep(api.IdleTimeMilliseconds * time.Millisecond)
					continue
				}
			}
			temp = conversationResponseDataString

			err = json.Unmarshal([]byte(conversationResponseDataString), &conversationResponse)
			if err != nil {
				time.Sleep(api.IdleTimeMilliseconds * time.Millisecond)
				continue
			}

			message := conversationResponse.Message
			if oldContent == "" {
				callbackChannel <- conversationResponseDataString
			} else {
				message.Content.Parts[0] = oldContent + (message.Content.Parts[0])
				withOldContentJsonString, _ := json.Marshal(conversationResponse)
				callbackChannel <- string(withOldContentJsonString)
			}

			maxTokens = message.Metadata.FinishDetails.Type == "max_tokens"
			if maxTokens {
				if request.ContinueText == "" {
					callbackChannel <- doneFlag
					close(callbackChannel)
				} else {
					oldContent = message.Content.Parts[0]
				}
				break
			}

			endTurn := message.EndTurn
			if endTurn {
				callbackChannel <- doneFlag
				close(callbackChannel)
				break
			}
		}
		if maxTokens && request.ContinueText != "" {
			time.Sleep(time.Second)

			continueMessageID := uuid.NewString()
			parentMessageID := conversationResponse.Message.ID
			conversationID := conversationResponse.ConversationID
			requestBodyJson := fmt.Sprintf(`
			{
				"action": "next",
				"messages": [{
					"id": "%s",
					"author": {
						"role": "%s"
					},
					"role": "%s",
					"content": {
						"content_type": "text",
						"parts": ["%s"]
					}
				}],
				"parent_message_id": "%s",
				"model": "%s",
				"conversation_id": "%s",
				"timezone_offset_min": %d,
				"variant_purpose": "%s",
				"continue_text": "%s"
			}`, continueMessageID,
				defaultRole,
				defaultRole,
				request.ContinueText,
				parentMessageID,
				request.Model,
				conversationID,
				request.TimezoneOffsetMin,
				request.VariantPurpose,
				request.ContinueText)
			var request StartConversationRequest
			json.Unmarshal([]byte(requestBodyJson), &request)
			sendConversationRequest(c, callbackChannel, request, oldContent, continueMessageID)
		}
	}()
	return true
}

//goland:noinspection GoUnhandledErrorResult
func stopGenerate(id string) {
	webdriver.WebDriver.ExecuteScript(fmt.Sprintf("xhrMap.get('%s').abort();", id), nil)
}

//goland:noinspection GoUnhandledErrorResult
func GenerateTitle(c *gin.Context) {
	var request GenerateTitleRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(parseJsonErrorMessage))
		return
	}

	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation/gen_title/" + c.Param("id")
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScript(url, accessToken, string(jsonBytes), generateTitleErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == generateTitleErrorMessage {
		tryToRefreshPage()
		GenerateTitle(c)
	} else {
		c.Writer.Write([]byte(responseText.(string)))
	}
}

//goland:noinspection GoUnhandledErrorResult
func GetConversation(c *gin.Context) {
	url := apiPrefix + "/conversation/" + c.Param("id")
	handleGet(c, url, getContentErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func UpdateConversation(c *gin.Context) {
	var request PatchConversationRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(parseJsonErrorMessage))
		return
	}

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
		tryToRefreshPage()
		UpdateConversation(c)
	} else {
		c.Writer.Write([]byte(responseText.(string)))
	}
}

//goland:noinspection GoUnhandledErrorResult
func FeedbackMessage(c *gin.Context) {
	var request FeedbackMessageRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ReturnMessage(parseJsonErrorMessage))
		return
	}

	jsonBytes, _ := json.Marshal(request)
	url := apiPrefix + "/conversation/message_feedback"
	accessToken := api.GetAccessToken(c.GetHeader(api.AuthorizationHeader))
	script := getPostScript(url, accessToken, string(jsonBytes), feedbackMessageErrorMessage)
	responseText, err := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
	if handleSeleniumError(err, script, c) {
		return
	}

	if responseText == feedbackMessageErrorMessage {
		tryToRefreshPage()
		FeedbackMessage(c)
	} else {
		c.Writer.Write([]byte(responseText.(string)))
	}
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
		tryToRefreshPage()
		ClearConversations(c)
	} else {
		c.Writer.Write([]byte(responseText.(string)))
	}
}

//goland:noinspection GoUnhandledErrorResult
func handleSeleniumError(err error, script string, c *gin.Context) bool {
	if err != nil {
		if _, ok := err.(*selenium.Error); ok {
			webdriver.NewSessionAndRefresh()
			responseText, _ := webdriver.WebDriver.ExecuteScriptAsync(script, nil)
			c.Writer.Write([]byte(responseText.(string)))
			return true
		}
	}

	return false
}

//goland:noinspection GoUnhandledErrorResult
func GetModels(c *gin.Context) {
	url := apiPrefix + "/models"
	handleGet(c, url, getModelsErrorMessage)
}

func GetAccountCheck(c *gin.Context) {
	url := apiPrefix + "/accounts/check"
	handleGet(c, url, getAccountCheckErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handleGet(c *gin.Context, url string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", api.GetAccessToken(c.GetHeader(api.AuthorizationHeader)))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(errorMessage))
		return
	}

	data, _ := io.ReadAll(resp.Body)
	c.Writer.Write(data)
}

//goland:noinspection GoUnhandledErrorResult
func UserLogin(c *gin.Context) {
	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, api.ReturnMessage(parseUserInfoErrorMessage))
		return
	}

	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}...)

	// get csrf token
	req, _ := http.NewRequest(http.MethodGet, csrfUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(getCsrfTokenErrorMessage))
		return
	}

	data, _ := io.ReadAll(resp.Body)
	responseMap := make(map[string]string)
	json.Unmarshal(data, &responseMap)

	// get authorized url
	params := fmt.Sprintf(
		"callbackUrl=/&csrfToken=%s&json=true",
		responseMap["csrfToken"],
	)
	req, err = http.NewRequest(http.MethodPost, promptLoginUrl, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(getAuthorizedUrlErrorMessage))
		return
	}

	// get state
	data, _ = io.ReadAll(resp.Body)
	json.Unmarshal(data, &responseMap)
	req, err = http.NewRequest(http.MethodGet, responseMap["url"], nil)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(getStateErrorMessage))
		return
	}

	// check username
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	state, _ := doc.Find("input[name=state]").Attr("value")
	params = fmt.Sprintf(
		"state=%s&username=%s&js-available=true&webauthn-available=true&is-brave=false&webauthn-platform-available=false&action=default",
		state,
		loginInfo.Username,
	)
	req, err = http.NewRequest(http.MethodPost, loginUsernameUrl+state, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(emailInvalidErrorMessage))
		return
	}

	// check username and password
	params = fmt.Sprintf(
		"state=%s&username=%s&password=%s&action=default",
		state,
		loginInfo.Username,
		loginInfo.Password,
	)
	req, err = http.NewRequest(http.MethodPost, loginPasswordUrl+state, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(emailOrPasswordInvalidErrorMessage))
		return
	}

	// get access token
	req, err = http.NewRequest(http.MethodGet, authSessionUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.JSON(resp.StatusCode, api.ReturnMessage(getAccessTokenErrorMessage))
		return
	}

	data, _ = io.ReadAll(resp.Body)
	c.Writer.Write(data)
}
