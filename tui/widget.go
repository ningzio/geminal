package tui

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Primitive interface {
	Primitive() tview.Primitive
}

// HistoryWidget 历史记录组件
type HistoryWidget interface {
	Primitive

	// NewHistory 插入一个新的历史记录, 并且放在第一个位置
	NewHistory(conv *Conversation)
	// GetCurrentChatID 获取当前聊天窗口的 chat id
	GetCurrentChatID() string
}

// InputWidget 用户输入组件
type InputWidget interface {
	Primitive
}

// ChatWidget 聊天窗口组件
type ChatWidget interface {
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

type WarningWidget interface {
	Primitive
	SetText(message string)
	SetButtons(buttons ...string)
	SetColor(color tcell.Color)
}
