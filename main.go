package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api/auth"
	"github.com/linweiyuan/go-chatgpt-api/api/user"
)

func main() {
	f, _ := os.Create("api.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	router.POST("/user/login", user.Login)
	router.GET("/auth/session", auth.RenewAccessToken)

	router.Run(":8080")
}
