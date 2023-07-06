package chatgpt

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type (
	ChatGpt struct {
		OpenAI *openai.Client
	}
	ChatGptReq struct {
		Prompt []openai.ChatCompletionMessage `json:"prompt"`
	}
)

func (c *ChatGpt) Handle(e *gin.Engine) {
	e.POST("/chatgpt", func(ctx *gin.Context) {
		var chatgptReq ChatGptReq
		ctx.Bind(&chatgptReq)
		fmt.Println(chatgptReq.Prompt)
		resp, err := c.OpenAI.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, Messages: chatgptReq.Prompt,
		})
		//TODO CHAT GPT request
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
		}
		fmt.Println(resp.Usage.TotalTokens)
		ctx.JSON(http.StatusAccepted, resp.Choices[0].Message.Content)
	})
}
