package external

import (
	"context"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

var dailyTaskListPrompt string

var openaiClient *openai.Client

func initOpenaiClient() {
	apiKey := os.Getenv("OPENAI_KEY")
	openaiClient = openai.NewClient(apiKey)
}

func MakeChatGPTRequest(ctx context.Context, content string) (string, error) {
	if openaiClient == nil {
		initOpenaiClient()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()

	resp, err := openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func GetDailyTaskList(ctx context.Context) (string, error) {
	return MakeChatGPTRequest(ctx, dailyTaskListPrompt)
}

func GetLibraryMessage(ctx context.Context) (string, error) {
	return MakeChatGPTRequest(ctx, "Miks peaksin täna raamatukokku minema?")
}
