package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()

	server.POST("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	server.Run(":8080")
}
