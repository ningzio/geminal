package tui

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/rivo/tview"
)

type Conversation struct {
	ChatID  string
	Title   string
	Content []byte
}

type Handler interface {
	CreateConversation(ctx context.Context) (*Conversation, error)
	ListConversation(ctx context.Context) ([]*Conversation, error)
	GetConversation(ctx context.Context, chatID string) (*Conversation, error)

	Talk(ctx context.Context, chatID string, writer io.Writer, prompt string) error
}

type Primitive interface {
	Primitive() tview.Primitive
}

// HistoryWight 历史记录组件
type HistoryWight interface {
	Primitive

	// NewHistory 插入一个新的历史记录, 并且放在第一个位置
	NewHistory(conv *Conversation)
	// GetCurrentChatID 获取当前聊天窗口的 chat id
	GetCurrentChatID() string
}

// InputWight 用户输入组件
type InputWight interface {
	Primitive
}

// ChatWight 聊天窗口组件
type ChatWight interface {
	Primitive
	// Writer 返回当前的 chat view 的 writer, 用于写入聊天内容
	// Writer(chatID string) io.Writer
	Writer() io.Writer

	// NewChatView 新建一个聊天窗口, 并切换到该窗口
	NewChatView(conversation *Conversation)

	// SwitchView 切换 chat view
	SwitchView(chatId string) bool

	SetTitle(title string)
}

func NewApplication(handler Handler) (*Application, error) {
	grid := tview.NewGrid().SetRows(-8, -2).SetColumns(-2, -8)
	tApp := tview.NewApplication()
	tApp.SetRoot(grid, true).EnableMouse(true)
	app := &Application{
		backend: handler,
		app:     tApp,
		grid:    grid,
	}

	chat := NewChatTUI(func() { app.app.Draw() })
	app.chat = chat

	input := NewInputTUI(app.submitFunc())
	app.input = input

	conversations, err := app.backend.ListConversation(context.Background())
	if err != nil {
		return nil, fmt.Errorf("init application: %w", err)
	}
	history := NewHistoryTUI(app.onHistoryChange, conversations...)
	app.history = history

	app.addInputToGrid(input)
	app.addChatToGrid(chat)
	app.addHistoryToGrid(history)

	app.app.SetRoot(app.grid, true).EnableMouse(true)
	return app, nil
}

type Application struct {
	backend Handler

	app  *tview.Application
	grid *tview.Grid

	input   InputWight
	chat    ChatWight
	history HistoryWight
}

func (app *Application) addInputToGrid(input InputWight) {
	app.grid.AddItem(input.Primitive(), 1, 1, 1, 1, 0, 0, true)
}

func (app *Application) addChatToGrid(chat ChatWight) {
	app.grid.AddItem(chat.Primitive(), 0, 1, 1, 1, 0, 0, false)
}

func (app *Application) addHistoryToGrid(history HistoryWight) {
	app.grid.AddItem(history.Primitive(), 0, 0, 2, 1, 0, 0, false)
}

func (app *Application) submitFunc() OnUserSubmit {
	return func(input string) {
		chatID := app.history.GetCurrentChatID()
		// new conversation
		if len(chatID) == 0 {
			conversation, err := app.backend.CreateConversation(context.Background())
			if err != nil {
				// TODO: error handling
				log.Println(err)
				return
			}
			chatID = conversation.ChatID
			app.chat.NewChatView(conversation)
			app.history.NewHistory(conversation)
		}
		go func() {
			// if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(chatID), input); err != nil {
			if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(), input); err != nil {
				// TODO: error handling
				log.Println(err)
				return
			}
		}()
	}
}

func (app *Application) onHistoryChange(index int, title, chatID string, shortcut rune) {
	ok := app.chat.SwitchView(chatID)
	if !ok {
		// restore conversation from repository
		conversation, err := app.backend.GetConversation(context.Background(), chatID)
		if err != nil {
			// TODO: error handling
			log.Println(err)
			return
		}
		app.chat.NewChatView(conversation)
	}
}

func (app *Application) Run() error {
	return app.app.Run()
}
