package repo

import (
	"context"

	"github.com/ningzio/geminal/internal"
)

var _ internal.Repository = (*Repository)(nil)

type Repository struct{}

// GetConversationByChatID implements internal.Repository.
func (*Repository) GetConversationByChatID(ctx context.Context, chatID string) (*internal.Conversation, error) {
	if chatID == "test" {
		return &internal.Conversation{
			ChatID: chatID,
			Title:  "test",
			Messages: []*internal.Message{
				{
					ChatID:      chatID,
					Role:        "You",
					ContentType: "text",
					Content:     "hello",
				},
			},
		}, nil
	} else if chatID == "test2" {
		return &internal.Conversation{
			ChatID: chatID,
			Title:  "test2",
			Messages: []*internal.Message{
				{
					ChatID:      chatID,
					Role:        "You",
					ContentType: "text",
					Content:     "hello111",
				},
				{
					ChatID:      chatID,
					Role:        "Gemini Pro",
					ContentType: "text",
					Content:     "hello there",
				},
			},
		}, nil
	}
	return &internal.Conversation{
		ChatID: chatID,
		Title:  "",
		Messages: []*internal.Message{
			{
				ChatID:      chatID,
				Role:        "You",
				ContentType: "text",
				Content:     "",
			},
		},
	}, nil
}

// LoadHistory implements internal.Repository.
func (*Repository) LoadHistory(ctx context.Context) ([]*internal.Conversation, error) {
	return nil, nil
	// return []*internal.Conversation{
	// 	{
	// 		ChatID: "test",
	// 		Title:  "test",
	// 		Messages: []*internal.Message{
	// 			{
	// 				ChatID:      "test",
	// 				Role:        "You",
	// 				ContentType: "text",
	// 				Content:     []byte("hello"),
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ChatID: "test2",
	// 		Title:  "test2",
	// 		Messages: []*internal.Message{
	// 			{
	// 				ChatID:      "test2",
	// 				Role:        "You",
	// 				ContentType: "text",
	// 				Content:     []byte("hello111"),
	// 			},
	// 			{
	// 				ChatID:      "test2",
	// 				Role:        "Gemini Pro",
	// 				ContentType: "text",
	// 				Content:     []byte("hello there"),
	// 			},
	// 		},
	// 	},
	// }, nil
}

// SaveConversation implements internal.Repository.
func (*Repository) SaveConversation(ctx context.Context, conversation *internal.Conversation) error {
	return nil
}
