package handler

import (
	"net/http"

	"github.com/Nekodigi/gpt-playground-backend/config"
	"github.com/Nekodigi/gpt-playground-backend/handler/chatgpt"
	"github.com/Nekodigi/gpt-playground-backend/lib/charge"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

var (
	openaiClient *openai.Client
	chrg         *charge.Charge
)

func init() {
	conf := config.Load()

	openaiClient = openai.NewClient(conf.ChatGPTToken)
	chrg = charge.InitCharge(conf)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Router(e *gin.Engine) {
	e.Use(CORSMiddleware())
	(&chatgpt.ChatGpt{OpenAI: openaiClient, Chrg: chrg}).Handle(e)
	e.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
}
