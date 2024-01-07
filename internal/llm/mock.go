package llm

import (
	"context"
	"errors"
	"os"

	"github.com/ningzio/geminal/internal"
)

var _ internal.LLM = (*Mock)(nil)

type Mock struct{}

// Name implements internal.LLM.
func (*Mock) Name() string {
	return "Mock"
}

// NewSession implements internal.LLM.
func (*Mock) NewSession(ctx context.Context, chatID string, history ...*internal.Message) error {
	panic("unimplemented")
}

// Talk implements internal.LLM.
//
//nolint:govet
func (c *Mock) Talk(ctx context.Context, chatID string, history []*internal.Message, messages ...*internal.Message) (*internal.Message, error) {
	return nil, errors.New("mock error")
	f, err := os.ReadFile("/Users/ningzi/workspace/personal/geminal/internal/llm/mark.log")
	if err != nil {
		return &internal.Message{}, nil
	}

	return &internal.Message{
		ChatID:      chatID,
		Role:        "Gemini Pro",
		ContentType: "text",
		Content:     string(f),
	}, nil
}
