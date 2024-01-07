package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HistoryHandler interface {
	OnConversationChanged(chatID string)
	DeleteConversation(chatID string) error
	RenameConversation(chatID, newTitle string) error
}

const (
	conversations = "conversations"
	options       = "options"
	delete        = "delete"
	renameInput   = "rename"
	warningModal  = "warning"
)

var _ HistoryWidget = (*HistoryTUI)(nil)

// NewHistoryTUI 创建一个历史聊天记录组件
func NewHistoryTUI(handler HistoryHandler, messages ...*Conversation) *HistoryTUI {
	page := tview.NewPages()

	list := tview.NewList()
	deleteConfirm := newDeleteModal()
	option := newOption()
	input := newInputField()
	warning := errorModal()

	list.ShowSecondaryText(false)
	list.SetBorder(true)

	page.AddPage(conversations, list, true, true)
	page.AddPage(options, option, true, false)
	page.AddPage(delete, deleteConfirm, true, false)
	page.AddPage(renameInput, input, true, false)
	page.AddPage(warningModal, warning, true, false)

	history := &HistoryTUI{
		handler:       handler,
		list:          list,
		deleteConfirm: deleteConfirm,
		options:       option,
		renameInput:   input,
		page:          page,
	}

	warning.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		page.SwitchToPage(conversations)
	})

	list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		handler.OnConversationChanged(secondaryText)
	})
	list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		history.ShowOptionPage(i, s2)
	})
	for _, msg := range messages {
		list.AddItem(msg.Title, msg.ChatID, 0, nil)
	}

	option.SetDoneFunc(func() { page.SwitchToPage(conversations) })
	option.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		switch s1 {
		case optionDelete:
			history.ShowDeletePage()
		case optionRename:
			history.ShowRenameTitlePage()
		case optionNothing:
			page.SwitchToPage(conversations)
		}
	})

	return history
}

type optionContext struct {
	index  int
	chatID string
}

type HistoryTUI struct {
	handler HistoryHandler

	// components
	list          *tview.List
	options       *tview.List
	deleteConfirm *tview.Modal
	renameInput   *tview.InputField
	warning       *tview.Modal

	// to organize components
	page *tview.Pages

	optionContext *optionContext
}

// Primitive implements Primitive.
func (h *HistoryTUI) Primitive() tview.Primitive {
	return h.page
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

func (h *HistoryTUI) ShowOptionPage(index int, chatID string) {
	h.optionContext = &optionContext{
		index:  index,
		chatID: chatID,
	}
	h.page.SwitchToPage(options)
}

func (h *HistoryTUI) ShowDeletePage() {
	h.deleteConfirm.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		defer func() {
			h.page.SwitchToPage(conversations)
			h.optionContext = nil
		}()
		if buttonLabel != deleteConfirmButton {
			return
		}
		h.deleteHistory()
	})
	h.page.SwitchToPage(delete)
}

func (h *HistoryTUI) deleteHistory() {
	if err := h.handler.DeleteConversation(h.optionContext.chatID); err != nil {
		h.warning.SetText(err.Error())
		h.page.SwitchToPage(warningModal)
	}
	h.list.RemoveItem(h.optionContext.index)
}

func (h *HistoryTUI) ShowRenameTitlePage() {
	h.renameInput.SetDoneFunc(func(key tcell.Key) {
		defer func() {
			h.renameInput.SetText("")
			h.optionContext = nil
		}()
		if key == tcell.KeyEnter {
			if err := h.handler.RenameConversation(h.optionContext.chatID, h.renameInput.GetText()); err != nil {
				h.warning.SetText(err.Error())
				h.page.SwitchToPage(warningModal)
			} else {
				h.list.SetItemText(h.optionContext.index, h.renameInput.GetText(), h.optionContext.chatID)
			}
		}
		h.page.SwitchToPage(conversations)
	})
	h.page.SwitchToPage(renameInput)
}

func newInputField() *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel("New Title: ")
	input.SetBorder(true)
	return input
}

var deleteConfirmButton = "ok"

func newDeleteModal() *tview.Modal {
	modal := tview.NewModal()
	modal.SetText("delete this conversation?(press ESC to cancel)")
	modal.AddButtons([]string{deleteConfirmButton})
	return modal
}

// options on conversations
var (
	optionDelete  = "Delete this conversation?"
	optionRename  = "Rename this conversation?"
	optionNothing = "Do Nothing(you can just press ESC)"
)

func newOption() *tview.List {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)

	list.AddItem(optionDelete, "", 0, nil)
	list.AddItem(optionRename, "", 0, nil)
	list.AddItem(optionNothing, "", 0, nil)

	return list
}

func errorModal() *tview.Modal {
	modal := tview.NewModal()
	modal.AddButtons([]string{"OK"})
	return modal
}
