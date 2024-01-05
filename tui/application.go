package tui

import (
	"context"
	"fmt"
	"io"

	"github.com/ningzio/geminal/internal"
	"github.com/rivo/tview"
)

type Handler interface {
	LoadHistory(ctx context.Context) ([]*internal.Conversation, error)
	Talk(ctx context.Context, chatID string, writer io.Writer, prompt string)
	NewConversation(ctx context.Context) *internal.Conversation
}

type Primitive interface {
	Primitive() tview.Primitive
}

// HistoryWight 历史记录组件
type HistoryWight interface {
	Primitive

	// NewHistory 插入一个新的历史记录, 并且放在第一个位置
	NewHistory(conv *internal.Conversation)
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
	Writer() io.Writer

	// NewChatView 新建一个聊天窗口, 并切换到该窗口
	NewChatView(chatId string, title string, content []byte)

	// SwitchView 切换 chat view, 如果 chatId 不存在, 则新建一个, 并返回 false
	SwitchView(chatId string) bool

	SetTitle(title string)
}

func NewApplication(handler Handler) (*Application, error) {
	grid := tview.NewGrid().SetRows(-8, -2).SetColumns(-2, -8)
	tApp := tview.NewApplication()
	tApp.SetRoot(grid, true).EnableMouse(true)
	app := &Application{
		h:    handler,
		app:  tApp,
		grid: grid,
	}

	chat := NewChatTUI(func() { app.app.Draw() })
	app.chat = chat

	input := NewInputTUI(app.submitFunc())
	app.input = input

	conversations, err := app.h.LoadHistory(context.Background())
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
	h Handler

	app  *tview.Application
	grid *tview.Grid

	input   InputWight
	chat    ChatWight
	history HistoryWight

	chatView tview.Primitive
}

func (app *Application) replaceChatView() {
	app.grid.RemoveItem(app.chatView)
	app.addChatToGrid(app.chat)
}

func (app *Application) addInputToGrid(input InputWight) {
	app.grid.AddItem(input.Primitive(), 1, 1, 1, 1, 0, 0, true)
}

func (app *Application) addChatToGrid(chat ChatWight) {
	app.chatView = chat.Primitive()
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
			conversation := app.h.NewConversation(context.Background())
			app.chat.NewChatView(conversation.ChatID, conversation.Title, nil)
			app.history.NewHistory(conversation)
		}
		go app.h.Talk(context.Background(), chatID, app.chat.Writer(), input)
	}
}

func (app *Application) onHistoryChange(index int, title, chatID string, shortcut rune) {
	ok := app.chat.SwitchView(chatID)
	if !ok {
		// title, content, err := app.h.GetConversationByChatID(context.Background(), chatID)
		// if err != nil {
		// 	// TODO: error handling
		// 	return
		// }
		app.chat.NewChatView(chatID, "Untitled", nil)
	}
	app.replaceChatView()
}

func (app *Application) Run() error {
	return app.app.Run()
}
