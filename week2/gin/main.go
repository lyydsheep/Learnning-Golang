package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, GO")
	})

	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "name is "+name+"\n")
		ctx.String(http.StatusOK, "id is "+id)
	})

	server.GET("/views/*abcd", func(ctx *gin.Context) {
		page := ctx.Param("abcd")
		ctx.String(http.StatusOK, "page is "+page)
	})

	server.Run(":8080")
}
