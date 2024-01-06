package internal

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/ningzio/geminal/tui"
)

type Renderer interface {
	RenderMessage(writer io.Writer, message *Message)
}

type Repository interface {
	// LoadHistory 负责加载历史聊天记录
	LoadHistory(ctx context.Context) ([]*Conversation, error)
	// GetConversationByChatID 负责根据 chat id 获取对应的聊天记录
	GetConversationByChatID(ctx context.Context, chatID string) (*Conversation, error)
	// SaveConversation 负责保存历史聊天记录
	SaveConversation(ctx context.Context, conversation *Conversation) error
}

type LLM interface {
	Name() string
	NewSession(ctx context.Context, chatID string, history ...*Message) error
	Talk(ctx context.Context, chatID string, history []*Message, messages ...*Message) (*Message, error)
}

// Conversation represent a conversation between user and AI
type Conversation struct {
	ChatID      string
	Title       string
	Messages    []*Message
	StartTime   time.Time
	UpdatedTime time.Time
}

func newConversation() *Conversation {
	return &Conversation{
		ChatID:      uuid.NewString(),
		Title:       "Untitled",
		StartTime:   time.Now(),
		UpdatedTime: time.Now(),
	}
}

type Message struct {
	ChatID      string
	Role        string
	ContentType string
	Content     string
	ErrMsg      string
}

var _ tui.Backend = (*Handler)(nil)

func NewHandler(llm LLM, repo Repository, render Renderer) *Handler {
	return &Handler{
		llm:    llm,
		repo:   repo,
		render: render,
	}
}

type Handler struct {
	render Renderer
	repo   Repository
	llm    LLM
}

// CreateConversation implements tui.Handler.
func (h *Handler) CreateConversation(ctx context.Context) (*tui.Conversation, error) {
	conv := newConversation()
	if err := h.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}
	return &tui.Conversation{
		ChatID:  conv.ChatID,
		Title:   conv.Title,
		Content: nil,
	}, nil
}

// GetConversation implements tui.Handler.
func (h *Handler) GetConversation(ctx context.Context, chatID string) (*tui.Conversation, error) {
	conv, err := h.repo.GetConversationByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	result := &tui.Conversation{
		ChatID: conv.ChatID,
		Title:  conv.Title,
	}

	buf := bytes.Buffer{}

	for _, msg := range conv.Messages {
		h.render.RenderMessage(&buf, msg)
	}

	result.Content = buf.Bytes()
	return result, nil
}

// ListConversation implements tui.Handler.
func (h *Handler) ListConversation(ctx context.Context) ([]*tui.Conversation, error) {
	conversations, err := h.repo.LoadHistory(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*tui.Conversation, 0, len(conversations))

	for _, conv := range conversations {
		c := &tui.Conversation{
			ChatID: conv.ChatID,
			Title:  conv.Title,
		}

		buf := bytes.Buffer{}

		for _, msg := range conv.Messages {
			h.render.RenderMessage(&buf, msg)
		}

		c.Content = buf.Bytes()
		result = append(result, c)
	}
	return result, nil
}

func (h *Handler) Talk(ctx context.Context, chatID string, writer io.Writer, prompt string) error {
	conv, err := h.repo.GetConversationByChatID(ctx, chatID)
	if err != nil {
		return err
	}

	message := &Message{
		ChatID:      chatID,
		Role:        "You",
		ContentType: "text",
		Content:     prompt,
	}
	h.render.RenderMessage(writer, message)

	result, err := h.llm.Talk(ctx, chatID, conv.Messages, message)
	if err != nil {
		return err
	}

	h.render.RenderMessage(writer, result)

	conv.Messages = append(conv.Messages, message)
	conv.Messages = append(conv.Messages, result)

	return h.repo.SaveConversation(ctx, conv)
}
