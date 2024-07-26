package chatgpt

import (
	"context"
	"fmt"
	"github.com/assimon/ai-anti-bot/adapter"
	"github.com/assimon/ai-anti-bot/pkg/json"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

var _ adapter.IModel = (*ChatGpt)(nil)

type ChatGpt struct {
	adapter.Option
	Client *openai.Client
}

func NewChatGpt(option adapter.Option) *ChatGpt {
	cfg := openai.DefaultConfig(option.ApiKey)
	if option.Proxy != "" {
		cfg.BaseURL = option.Proxy
	}
	return &ChatGpt{
		Option: option,
		Client: openai.NewClientWithConfig(cfg),
	}
}

func (c *ChatGpt) RecognizeTextMessage(ctx context.Context, userInfo, message string) (adapter.RecognizeResult, error) {
	var result adapter.RecognizeResult
	prompt := fmt.Sprintf(viper.GetString("prompt.text"), userInfo, message)
	req := openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := c.Client.CreateChatCompletion(
		ctx,
		req,
	)
	if err != nil {
		return result, err
	}

	err = json.C.UnmarshalFromString(resp.Choices[0].Message.Content, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c *ChatGpt) RecognizeImageMessage(ctx context.Context, userInfo, file string) (adapter.RecognizeResult, error) {
	var result adapter.RecognizeResult
	prompt := fmt.Sprintf(viper.GetString("prompt.image"), userInfo)
	req := openai.ChatCompletionRequest{
		Model: c.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: prompt,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL:    file,
							Detail: openai.ImageURLDetailLow,
						},
					},
				},
			},
		},
	}
	resp, err := c.Client.CreateChatCompletion(
		ctx,
		req,
	)
	if err != nil {
		return result, err
	}
	err = json.C.UnmarshalFromString(resp.Choices[0].Message.Content, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
