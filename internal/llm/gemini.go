package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/ningzio/geminal/internal"
	"google.golang.org/api/option"
)

var _ internal.LLM = (*GeminiAI)(nil)

func NewGeminiAI(apiKey string) (*GeminiAI, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("new gemini ai: %w", err)
	}
	model := client.GenerativeModel("gemini-pro")
	return &GeminiAI{
		client:   client,
		model:    model,
		sessions: make(map[string]*genai.ChatSession),
	}, nil
}

// GeminiAI is a client for the Gemini AI API.
type GeminiAI struct {
	client *genai.Client
	model  *genai.GenerativeModel

	sessions map[string]*genai.ChatSession
}

// Name implements internal.LLM.
func (*GeminiAI) Name() string {
	return "Gemini Pro"
}

// NewSession implements internal.LLM.
func (*GeminiAI) NewSession(ctx context.Context, chatID string, history ...*internal.Message) error {
	return nil
}

// Talk implements internal.LLM.
func (ai *GeminiAI) Talk(ctx context.Context, chatID string, messages ...*internal.Message) (result *internal.Message) {
	session, ok := ai.sessions[chatID]
	if !ok {
		session = ai.model.StartChat()
		ai.sessions[chatID] = session
	}

	var prompts []genai.Part
	for _, msg := range messages {
		prompts = append(prompts, genai.Text(string(msg.Content)))
	}

	result = &internal.Message{
		ChatID: chatID,
		Role:   ai.Name(),
	}

	resp, err := session.SendMessage(ctx, prompts...)
	if err != nil {
		result.ErrMsg = fmt.Sprintf("gemini ai: %v", err)
		return
	}

	if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != 0 {
		switch resp.PromptFeedback.BlockReason {
		case genai.BlockReasonSafety:
			for _, rating := range resp.PromptFeedback.SafetyRatings {
				result.ErrMsg += fmt.Sprintf("%s: %s, block: %v", rating.Category, rating.Probability, rating.Blocked)
			}
		case genai.BlockReasonOther:
			result.ErrMsg += fmt.Sprintf("block: %v", resp.PromptFeedback.BlockReason)
		}
		return
	}

	for _, i := range resp.Candidates[0].Content.Parts {
		switch i := i.(type) {
		case genai.Text:
			result.Content += string(i)
		case genai.Blob:
		default:
		}
	}
	return
}
