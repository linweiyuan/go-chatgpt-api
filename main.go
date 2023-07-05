package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-chatgpt-api/api/platform"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-chatgpt-api/middleware"

	http "github.com/bogdanfinn/fhttp"
)

func init() {
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
