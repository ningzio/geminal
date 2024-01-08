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
	pageConversations = "conversations"
	pageOptions       = "options"
	pageDeletePage    = "delete"
	pageRenameInput   = "rename"
	pageWarningModal  = "warning"
)

// NewHistoryTUI 创建一个历史聊天记录组件
func NewHistoryTUI(handler HistoryHandler) *History {
	page := tview.NewPages()
	conversation := newConversationList()
	deleteConversation := newDeleteModal()
	option := newOption()
	input := newInputField()
	warning := errorModal()

	history := &History{
		handler:            handler,
		conversations:      conversation,
		deleteConversation: deleteConversation,
		options:            option,
		renameTitle:        input,
		warning:            warning,
		page:               page,
	}

	history.addPages()
	history.setCallbackFunc()

	return history
}

type History struct {
	handler HistoryHandler

	// components
	conversations      *tview.List
	options            *tview.List
	deleteConversation *tview.Modal
	renameTitle        *tview.InputField
	warning            *tview.Modal

	// to organize components
	page *tview.Pages
}

// addPages adds pages to the history.
//
// It adds several pages to the history's page container. Each page is added with its corresponding content and settings.
func (h *History) addPages() {
	h.page.AddPage(pageConversations, h.conversations, true, true)
	h.page.AddPage(pageOptions, h.options, true, false)
	h.page.AddPage(pageDeletePage, h.deleteConversation, true, false)
	h.page.AddPage(pageRenameInput, h.renameTitle, true, false)
	h.page.AddPage(pageWarningModal, h.warning, true, false)
}

// setCallbackFunc 为 History 中的对话设置回调函数。
func (h *History) setCallbackFunc() {
	h.conversations.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		h.handler.OnConversationChanged(secondaryText)
	})
	h.conversations.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		h.ShowOptionPage(i, s2)
	})
}

// Primitive implements Primitive.
func (h *History) Primitive() tview.Primitive {
	return h.page
}

// NewHistory adds a new conversation to the history.
//
// Parameters:
// - conv: A pointer to the Conversation object to be added.
func (h *History) NewHistory(conv *Conversation) {
	h.conversations.InsertItem(0, conv.Title, conv.ChatID, 0, nil)
	h.conversations.SetCurrentItem(0)
}

// GetCurrentChatID returns the current chat ID.
func (h *History) GetCurrentChatID() string {
	if h.conversations.GetItemCount() == 0 {
		return ""
	}
	_, chatID := h.conversations.GetItemText(h.conversations.GetCurrentItem())
	return chatID
}

// ShowOptionPage displays the option page at the specified index for the given chat ID.
//
// Parameters:
// - index: The index of the option page to be displayed.
// - chatID: The ID of the chat for which the option page is being displayed.
func (h *History) ShowOptionPage(index int, chatID string) {
	h.options.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		switch s1 {
		case optionDelete:
			h.ShowDeletePage(index, chatID)
		case optionRename:
			h.ShowRenameTitlePage(index, chatID)
		case optionNothing:
			h.page.SwitchToPage(pageConversations)
		}
	})
	h.page.SwitchToPage(pageOptions)
}

func (h *History) ShowDeletePage(index int, chatID string) {
	h.deleteConversation.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		defer h.page.SwitchToPage(pageConversations)
		if buttonLabel != deleteConfirmButton {
			return
		}
		h.deleteHistory(index, chatID)
	})
	h.page.SwitchToPage(pageDeletePage)
}

func (h *History) deleteHistory(index int, chatID string) {
	if err := h.handler.DeleteConversation(chatID); err != nil {
		h.warning.SetText(err.Error())
		h.page.SwitchToPage(pageWarningModal)
	}
	h.conversations.RemoveItem(index)
}

func (h *History) ShowRenameTitlePage(index int, chatID string) {
	h.renameTitle.SetDoneFunc(func(key tcell.Key) {
		defer h.renameTitle.SetText("")
		if key == tcell.KeyEnter {
			if err := h.handler.RenameConversation(chatID, h.renameTitle.GetText()); err != nil {
				h.warning.SetText(err.Error())
				h.page.SwitchToPage(pageWarningModal)
			} else {
				h.conversations.SetItemText(index, h.renameTitle.GetText(), chatID)
			}
		}
		h.page.SwitchToPage(pageConversations)
	})
	h.page.SwitchToPage(pageRenameInput)
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

func newConversationList() *tview.List {
	list := tview.NewList()
	list.SetTitle("Conversations")
	list.SetBorder(true)
	list.ShowSecondaryText(false)

	return list
}
