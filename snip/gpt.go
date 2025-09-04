package main

import (
	"context"
	fmt "fmt"
	os "os"
	strings "strings"

	goopenai "github.com/sashabaranov/go-openai"
)

func GenerateTagsFromContent(content string) ([]string, error) {
	client := goopenai.NewClient(os.Getenv("YOUR_OPENAI_API_KEY"))
	ctx := context.Background()
	prompt := fmt.Sprintf(
		"Given the following code snippet, generate 3-5 concise tags separated by commas:\n\n%s", content,
	)
	resp, err := client.CreateChatCompletion(ctx, goopenai.ChatCompletionRequest{
		Model: goopenai.GPT3Dot5Turbo,
		Messages: []goopenai.ChatCompletionMessage{{
			Role:    "user",
			Content: prompt,
		}},
		MaxTokens: 60,
	})
	if err != nil {
		return nil, err
	}
	tags := strings.Split(resp.Choices[0].Message.Content, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	return tags, nil
}
