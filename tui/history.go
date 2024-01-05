package tui

import (
	"github.com/rivo/tview"
)

// OnHistoryListChangeFunc 当用户浏览历史记录列表时的处理方法
type OnHistoryListChangeFunc func(index int, title, chatID string, shortcut rune)

var _ HistoryWight = (*HistoryTUI)(nil)

// NewHistoryTUI 创建一个历史聊天记录组件
func NewHistoryTUI(handler OnHistoryListChangeFunc, messages ...*Conversation) *HistoryTUI {
	list := tview.NewList()
	list.SetTitle("History")
	list.SetBorder(true)
	list.ShowSecondaryText(false)

	list.SetChangedFunc(handler)

	for _, msg := range messages {
		list.AddItem(msg.Title, msg.ChatID, 0, nil)
	}

	return &HistoryTUI{
		list: list,
	}
}

type HistoryTUI struct {
	list *tview.List
}

// Primitive implements Primitive.
func (h *HistoryTUI) Primitive() tview.Primitive {
	return h.list
}

func (h *HistoryTUI) NewHistory(conv *Conversation) {
	h.list.InsertItem(0, conv.Title, conv.ChatID, 0, nil)
	h.list.SetCurrentItem(0)
}

func (h *HistoryTUI) GetCurrentChatID() string {
	if h.list.GetItemCount() == 0 {
		return ""
	}
	_, chatID := h.list.GetItemText(h.list.GetCurrentItem())
	return chatID
}
