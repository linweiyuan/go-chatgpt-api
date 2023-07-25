package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-chatgpt-api/api/imitate"
	"github.com/linweiyuan/go-chatgpt-api/api/platform"
	"github.com/linweiyuan/go-chatgpt-api/api/token"
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

	router.Use(middleware.CORS())
	router.Use(middleware.Authorization())

	setupChatGPTAPIs(router)
	setupPlatformAPIs(router)
	setupPandoraAPIs(router)
	setupImitateAPIs(router)
	setupTokenAPIs(router)
	router.NoRoute(api.Proxy)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, api.ReadyHint)
	})

	port := os.Getenv("PORT")
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
		chatgptGroup.POST("/backend-api/login", chatgpt.Login) // add support for other projects

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
		platformGroup.POST("/v1/login", platform.Login)

		apiGroup := platformGroup.Group("/v1")
		{
			apiGroup.POST("/chat/completions", platform.CreateChatCompletions)
			apiGroup.POST("/completions", platform.CreateCompletions)
		}
	}
}

//goland:noinspection SpellCheckingInspection
func setupPandoraAPIs(router *gin.Engine) {
	router.Any("/api/*path", func(c *gin.Context) {
		c.Request.URL.Path = strings.ReplaceAll(c.Request.URL.Path, "/api", "/chatgpt/backend-api")
		router.HandleContext(c)
	})
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

func setupTokenAPIs(router *gin.Engine) {
	tokenGroup := router.Group("/token")
	{
		tokenGroup.GET("/arkose", token.GetArkoseToken)
	}
}
