package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api/chatgpt"
	"github.com/linweiyuan/go-chatgpt-api/api/platform"
	_ "github.com/linweiyuan/go-chatgpt-api/env"
	"github.com/linweiyuan/go-chatgpt-api/middleware"
)

func init() {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.CheckHeaderMiddleware())

	setupChatGPTAPIs(router)

	setupPlatformAPIs(router)

	router.NoRoute(handleFallbackRoute)

	port := os.Getenv("GO_CHATGPT_API_PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}
}

func setupChatGPTAPIs(router *gin.Engine) {
	chatgptGroup := router.Group("/chatgpt")
	{
		chatgptGroup.POST("/login", chatgpt.Login)

		conversationsGroup := chatgptGroup.Group("/conversations")
		{
			conversationsGroup.GET("", chatgpt.GetConversations)

			// PATCH is official method, POST is added for Java support
			conversationsGroup.PATCH("", chatgpt.ClearConversations)
			conversationsGroup.POST("", chatgpt.ClearConversations)
		}

		conversationGroup := chatgptGroup.Group("/conversation")
		{
			conversationGroup.POST("", chatgpt.CreateConversation)
			conversationGroup.POST("/gen_title/:id", chatgpt.GenerateTitle)
			conversationGroup.GET("/:id", chatgpt.GetConversation)

			// rename or delete conversation use a same API with different parameters
			conversationGroup.PATCH("/:id", chatgpt.UpdateConversation)
			conversationGroup.POST("/:id", chatgpt.UpdateConversation)

			conversationGroup.POST("/message_feedback", chatgpt.FeedbackMessage)
		}

		// misc
		chatgptGroup.GET("/models", chatgpt.GetModels)
		chatgptGroup.GET("/accounts/check", chatgpt.GetAccountCheck)
	}
}

func setupPlatformAPIs(router *gin.Engine) {
	platformGroup := router.Group("/platform")
	{
		platformGroup.POST("/login", platform.Login)

		apiGroup := platformGroup.Group("/v1")
		{
			apiGroup.GET("/models", platform.ListModels)
			apiGroup.GET("/models/:model", platform.RetrieveModel)
			apiGroup.POST("/completions", platform.CreateCompletions)
			apiGroup.POST("/chat/completions", platform.CreateChatCompletions)
			apiGroup.POST("/edits", platform.CreateEdit)
			apiGroup.POST("/images/generations", platform.CreateImage)
			apiGroup.POST("/embeddings", platform.CreateEmbeddings)
			apiGroup.GET("/files", platform.ListFiles)
			apiGroup.POST("/moderations", platform.CreateModeration)
		}

		dashboardGroup := platformGroup.Group("/dashboard")
		{
			billingGroup := dashboardGroup.Group("/billing")
			{
				billingGroup.GET("/credit_grants", platform.GetCreditGrants)
				billingGroup.GET("/subscription", platform.GetSubscription)
			}

			userGroup := dashboardGroup.Group("/user")
			{
				userGroup.GET("/api_keys", platform.GetApiKeys)
			}
		}
	}
}

func handleFallbackRoute(c *gin.Context) {
	path := c.Request.URL.Path

	if strings.HasPrefix(path, "/chatgpt") {
		trimmedPath := strings.TrimPrefix(path, "/chatgpt")
		c.Request.URL.Path = trimmedPath
		chatgpt.Fallback(c)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	}
}
