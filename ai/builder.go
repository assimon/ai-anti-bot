package ai

import (
	"errors"
	"fmt"
	"github.com/assimon/ai-anti-bot/adapter"
	"github.com/assimon/ai-anti-bot/ai/chatgpt"
)

const (
	TypeChatGpt = "chatgpt"
)

func New(model string, option adapter.Option) (adapter.IModel, error) {
	if model == "" {
		return nil, errors.New("missing inbound type")
	}
	switch model {
	case TypeChatGpt:
		return chatgpt.NewChatGpt(option), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown ai model: %v", model))
	}
}
