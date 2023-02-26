package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api/auth"
	"github.com/linweiyuan/go-chatgpt-api/api/conversation"
	"github.com/linweiyuan/go-chatgpt-api/api/user"
	"github.com/linweiyuan/go-chatgpt-api/middleware"
)

func main() {
	f, _ := os.Create("api.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	router.POST("/user/login", user.Login)
	router.GET("/auth/session", auth.RenewAccessToken)

	authMiddleware := middleware.AuthMiddleware()

	router.GET("/conversations", authMiddleware, conversation.GetConversations)

	conversationGroup := router.Group("/conversation", authMiddleware)
	{
		conversationGroup.POST("", conversation.MakeConversation)
		conversationGroup.POST("/gen_title/:id", conversation.GenConversationTitle)
		conversationGroup.GET("/:id", conversation.GetConversation)

		// PATCH is official method, POST is added for Java
		// rename or delete conversation use a same API with different parameters
		conversationGroup.PATCH("/:id", conversation.PatchConversation)
		conversationGroup.POST("/:id", conversation.PatchConversation)

		conversationGroup.POST("/message_feedback", conversation.FeedbackMessage)
	}

	router.Run(":8080")
}
