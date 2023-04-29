package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yicone/go-chatgpt-api/api/chatgpt"
	"github.com/yicone/go-chatgpt-api/api/official"
	_ "github.com/yicone/go-chatgpt-api/env"
	"github.com/yicone/go-chatgpt-api/middleware"
	"github.com/yicone/go-chatgpt-api/webdriver"
)

func init() {
	gin.ForceConsoleColor()
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				webdriver.NewSessionAndRefresh()
			}
		}()
		c.Next()
	}
}
func main() {
	router := gin.Default()
	router.Use(Recover())
	router.Use(middleware.HeaderCheckMiddleware())

	// chatgpt
	conversationsGroup := router.Group("/conversations")
	{
		conversationsGroup.GET("", chatgpt.GetConversations)

		// PATCH is official method, POST is added for Java support
		conversationsGroup.PATCH("", chatgpt.ClearConversations)
		conversationsGroup.POST("", chatgpt.ClearConversations)
	}

	conversationGroup := router.Group("/conversation")
	{
		conversationGroup.POST("", chatgpt.StartConversation)
		conversationGroup.POST("/gen_title/:id", chatgpt.GenerateTitle)
		conversationGroup.GET("/:id", chatgpt.GetConversation)

		// rename or delete conversation use a same API with different parameters
		conversationGroup.PATCH("/:id", chatgpt.UpdateConversation)
		conversationGroup.POST("/:id", chatgpt.UpdateConversation)

		conversationGroup.POST("/message_feedback", chatgpt.FeedbackMessage)
	}

	// misc
	router.GET("/models", chatgpt.GetModels)
	router.GET("/accounts/check", chatgpt.GetAccountCheck)

	// auth
	router.POST("/auth/login", chatgpt.UserLogin) // login will cause some downtime because of CORS limits

	// ----------------------------------------------------------------------------------------------------

	// official api
	apiGroup := router.Group("/v1")
	{
		apiGroup.POST("/chat/completions", official.ChatCompletions)
	}
	router.GET("/dashboard/billing/credit_grants", official.CheckUsage)

	port := os.Getenv("CHATGPT_API_PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:" + err.Error())
	}
}
