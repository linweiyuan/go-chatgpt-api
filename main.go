package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-chatgpt-api/api/imitate"
	"github.com/linweiyuan/go-chatgpt-api/api/platform"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-chatgpt-api/middleware"

	http "github.com/bogdanfinn/fhttp"
	"github.com/acheong08/OpenAIAuth/auth"
)

type auth_struct struct {
	OpenAI_Email    string `json:"openai_email"`
	OpenAI_Password string `json:"openai_password"`
}

var (
	authorizations auth_struct
)

func init() {
	setupPUID()
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

//goland:noinspection SpellCheckingInspection
func main() {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.CheckHeaderMiddleware())

	setupChatGPTAPIs(router)
	setupPlatformAPIs(router)
	setupPandoraAPIs(router)
	setupImitateAPIs(router)
	router.NoRoute(api.Proxy)

	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, api.ReadyHint)
	})

	port := os.Getenv("GO_CHATGPT_API_PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}
}

//goland:noinspection SpellCheckingInspection
func setupChatGPTAPIs(router *gin.Engine) {
	chatgptGroup := router.Group("/chatgpt")
	{
		chatgptGroup.POST("/login", chatgpt.Login)

		conversationGroup := chatgptGroup.Group("/backend-api/conversation")
		{
			conversationGroup.POST("", chatgpt.CreateConversation)
		}
	}
}

func setupPlatformAPIs(router *gin.Engine) {
	platformGroup := router.Group("/platform")
	{
		platformGroup.POST("/login", platform.Login)

		apiGroup := platformGroup.Group("/v1")
		{
			apiGroup.POST("/chat/completions", platform.CreateChatCompletions)
			apiGroup.POST("/completions", platform.CreateCompletions)
		}
	}
}

//goland:noinspection SpellCheckingInspection
func setupPandoraAPIs(router *gin.Engine) {
	pandoraEnabled := os.Getenv("GO_CHATGPT_API_PANDORA") != ""
	if pandoraEnabled {
		router.Any("/api/*path", func(c *gin.Context) {
			c.Request.URL.Path = strings.ReplaceAll(c.Request.URL.Path, "/api", "/chatgpt/backend-api")
			router.HandleContext(c)
		})
	}
}

func setupImitateAPIs(router *gin.Engine) {
	imitateGroup := router.Group("/imitate")
	{
		imitateGroup.POST("/login", chatgpt.Login)

		apiGroup := imitateGroup.Group("/v1")
		{
			apiGroup.POST("/chat/completions", imitate.CreateChatCompletions)
		}
	}
}

func setupPUID() {
	authorizations.OpenAI_Email = os.Getenv("GO_CHATGPT_OPENAI_EMAIL")
	authorizations.OpenAI_Password = os.Getenv("GO_CHATGPT_OPENAI_PASSWORD")
	if authorizations.OpenAI_Email != "" && authorizations.OpenAI_Password != "" {
		go func() {
			for {
				authenticator := auth.NewAuthenticator(authorizations.OpenAI_Email, authorizations.OpenAI_Password, "")
				err := authenticator.Begin()
				if err != nil {
					log.Println(err)
					break
				}
				puid, err := authenticator.GetPUID()
				if err != nil {
					break
				}
				os.Setenv("GO_CHATGPT_API_PUID", puid)
				println(puid)
				time.Sleep(24 * time.Hour * 7)
			}
		}()
	}
}