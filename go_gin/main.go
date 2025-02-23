package main

import (
	"github.com/gin-gonic/gin"
)

func plain_get_handler(ctx *gin.Context) {
	ctx.String(200, "Hello world!")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	http_engine := gin.Default()

	http_engine.GET("/test_plain", plain_get_handler)

	http_engine.Run()
}
