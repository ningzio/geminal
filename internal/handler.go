package internal

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
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
	Talk(ctx context.Context, chatID string, messages ...*Message) *Message
}

// Conversation represent a conversation between user and AI
type Conversation struct {
	ChatID      string
	Title       string
	Messages    []*Message
	StartTime   time.Time
	UpdatedTime time.Time
}

type Message struct {
	ChatID      string
	Role        string
	ContentType string
	Content     string
	ErrMsg      string
}

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

func (h *Handler) NewConversation(ctx context.Context) *Conversation {
	return &Conversation{
		ChatID: uuid.NewString(),
		Title:  "Untitled",
	}
}

func (h *Handler) LoadHistory(ctx context.Context) ([]*Conversation, error) {
	return h.repo.LoadHistory(ctx)
}

func (h *Handler) Talk(ctx context.Context, chatID string, writer io.Writer, prompt string) {
	message := &Message{
		ChatID:      chatID,
		Role:        "You",
		ContentType: "text",
		Content:     prompt,
	}

	h.render.RenderMessage(writer, message)
	h.render.RenderMessage(writer, h.llm.Talk(ctx, chatID, message))
}
