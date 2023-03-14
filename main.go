package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api/conversation"
	"github.com/linweiyuan/go-chatgpt-api/middleware"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.PreCheckMiddleware())

	conversationsGroup := router.Group("/conversations")
	{
		conversationsGroup.GET("", conversation.GetConversations)

		// PATCH is official method, POST is added for Java support
		conversationsGroup.PATCH("", conversation.ClearConversations)
		conversationsGroup.POST("", conversation.ClearConversations)
	}

	conversationGroup := router.Group("/conversation")
	{
		conversationGroup.POST("", conversation.StartConversation)
		conversationGroup.POST("/gen_title/:id", conversation.GenerateTitle)
		conversationGroup.GET("/:id", conversation.GetConversation)

		// rename or delete conversation use a same API with different parameters
		conversationGroup.PATCH("/:id", conversation.UpdateConversation)
		conversationGroup.POST("/:id", conversation.UpdateConversation)

		conversationGroup.POST("/message_feedback", conversation.FeedbackMessage)
	}

	//goland:noinspection GoUnhandledErrorResult
	router.Run(":8080")
}
