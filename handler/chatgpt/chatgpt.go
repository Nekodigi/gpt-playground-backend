package chatgpt

import (
	"fmt"
	"net/http"

	"github.com/Nekodigi/gpt-playground-backend/lib/charge"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type (
	ChatGpt struct {
		OpenAI *openai.Client
		Chrg   *charge.Charge
	}
	ChatGptReq struct {
		UserId string                         `json:"userId"`
		Prompt []openai.ChatCompletionMessage `json:"prompt"`
	}
)

func (c *ChatGpt) Handle(e *gin.Engine) {
	e.POST("/chatgpt", func(ctx *gin.Context) {
		var chatgptReq ChatGptReq
		ctx.Bind(&chatgptReq)
		fmt.Println(chatgptReq.Prompt)

		if !c.Chrg.EnsureQuota(ctx, chatgptReq.UserId, 1) {
			return
		}

		resp, err := c.OpenAI.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, Messages: chatgptReq.Prompt,
		})

		if !c.Chrg.UseQuota(ctx, chatgptReq.UserId, float64(resp.Usage.TotalTokens)*0.4/1000) { //0.004 / 1k
			return
		}

		//TODO CHAT GPT request
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
		}
		fmt.Println(resp.Usage.TotalTokens)
		ctx.JSON(http.StatusAccepted, resp.Choices[0].Message.Content)
	})
}
