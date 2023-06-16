package api

//goland:noinspection GoSnakeCaseUsage
import (
	"bufio"
	"os"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/linweiyuan/go-chatgpt-api/env"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const (
	defaultErrorMessageKey             = "errorMessage"
	AuthorizationHeader                = "Authorization"
	ContentType                        = "application/x-www-form-urlencoded"
	UserAgent                          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	Auth0Url                           = "https://auth0.openai.com"
	LoginUsernameUrl                   = Auth0Url + "/u/login/identifier?state="
	LoginPasswordUrl                   = Auth0Url + "/u/login/password?state="
	ParseUserInfoErrorMessage          = "Failed to parse user login info."
	GetAuthorizedUrlErrorMessage       = "Failed to get authorized url."
	GetStateErrorMessage               = "Failed to get state."
	EmailInvalidErrorMessage           = "Email is not valid."
	EmailOrPasswordInvalidErrorMessage = "Email or password is not correct."
	GetAccessTokenErrorMessage         = "Failed to get access token."
	defaultTimeoutSeconds              = 300 // 5 minutes
)

var Client tls_client.HttpClient

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLogin interface {
	GetAuthorizedUrl(csrfToken string) (string, int, error)
	GetState(authorizedUrl string) (string, int, error)
	CheckUsername(state string, username string) (int, error)
	CheckPassword(state string, username string, password string) (string, int, error)
	GetAccessToken(code string) (string, int, error)
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	Client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithTimeoutSeconds(defaultTimeoutSeconds),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
	}...)
}

func ReturnMessage(msg string) gin.H {
	return gin.H{
		defaultErrorMessageKey: msg,
	}
}

func GetAccessToken(accessToken string) string {
	if !strings.HasPrefix(accessToken, "Bearer") {
		return "Bearer " + accessToken
	}
	return accessToken
}

//goland:noinspection GoUnhandledErrorResult
func HandleConversationResponse(c *gin.Context, resp *http.Response) (bool, string, string) {
	Status := false
	ParentMessageID := ""
	oldpart := c.GetString("oldpart")
	part := ""

	if len(oldpart) == 0 {
		c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	}

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

		// if len(oldpart) > 0 {
		// 	logger.Info(fmt.Sprintf("HandleConversationResponseContinue: %s", line))
		// }

		if Status {
			continue
		}

		if strings.HasPrefix(line, "[DONE]") {

		} else {
			if len(oldpart) > 0 {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(line[6:]), &result)
				if err == nil {
					message := result["message"].(map[string]interface{})
					content := message["content"].(map[string]interface{})
					parts := content["parts"].([]interface{})
					part = parts[0].(string)

					parts[0] = oldpart + part
					// message["id"] = c.GetString("msg_id")
					// metadata := message["metadata"].(map[string]interface{})
					// metadata["message_type"] = "next"

					resultJSON, err2 := json.Marshal(result)
					if err2 == nil {
						line = "data: " + string(resultJSON)
						// logger.Info(fmt.Sprintf("HandleConversationResponseAddon: %s", line))
					}
				}
			} else {

			}
		}

		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()

		if strings.HasPrefix(line, "[DONE]") {
			continue
		}

		data := line[6:]
		var result map[string]interface{}
		err2 := json.Unmarshal([]byte(data), &result)
		if err2 != nil {
			continue
		}
		message := result["message"].(map[string]interface{})
		status := message["status"].(string)

		if status == "finished_successfully" {
			if message["metadata"] != nil {
				metadata := message["metadata"].(map[string]interface{})
				if metadata["finish_details"] != nil {
					finishDetails := metadata["finish_details"].(map[string]interface{})
					finishType := finishDetails["type"].(string)
					if finishType == "max_tokens" {
						// logger.Info(fmt.Sprintf("finish_details中type的值: %s", finishType))
						// logger.Info(fmt.Sprintf("HandleConversationResponse: %s", line))
						content := message["content"].(map[string]interface{})
						parts := content["parts"].([]interface{})
						part = parts[0].(string)
						Status = true
						ParentMessageID = message["id"].(string)

						if len(oldpart) == 0 {
							c.Set("msg_id", ParentMessageID)
						}
					}
				}
			}
		}
	}
	return Status, ParentMessageID, part
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func NewHttpClient() tls_client.HttpClient {
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), []tls_client.HttpClientOption{
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
	}...)

	proxyUrl := os.Getenv("GO_CHATGPT_API_PROXY")
	if proxyUrl != "" {
		client.SetProxy(proxyUrl)
	}

	return client
}
