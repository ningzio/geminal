package tui

import (
	"io"

	"github.com/rivo/tview"
)

var _ ChatWidget = (*Chat)(nil)

// NewChat creates a new Chat instance.
//
// It takes a function onChangeFunc as a parameter, which is a callback function
// that will be called whenever there is a change in the Chat.
//
// The function returns a pointer to the Chat instance.
func NewChat(onChangeFunc func()) *Chat {
	return &Chat{
		onChangeFunc: onChangeFunc,
		views:        make(map[string]*view),
		page:         tview.NewPages(),
	}
}

type view struct {
	textView *tview.TextView
	writer   io.Writer
}

type Chat struct {
	onChangeFunc func()
	// current view
	view *view
	// ChatUI can hold multi text view, when user change
	// chat chat history in side bar, chat ui should
	// switch to correspond text view
	views map[string]*view

	page *tview.Pages
}

// DeleteView deletes a chat view identified by the given chatID.
//
// Parameters:
// - chatID: The ID of the chat view to be deleted.
func (c *Chat) DeleteView(chatID string) {
	delete(c.views, chatID)
	c.page.RemovePage(chatID)
}

// SetTitle sets the title of the chat.
//
// title: the title to be set for the chat.
func (c *Chat) SetTitle(title string) {
	c.view.textView.SetTitle(title)
}

// NewChatView creates a new chat view for the given conversation.
//
// It takes a pointer to a Chat struct as its receiver and a pointer to a Conversation struct as its parameter.
// The function creates a new text view using the conversation's title and the onChangeFunc callback.
// If the conversation has content, it writes the content to the view's writer.
// The function stores the new view in the views map, sets it as the current view, and adds it to the page.
func (c *Chat) NewChatView(conversation *Conversation) {
	view := newTextView(conversation.Title, c.onChangeFunc)
	if conversation.Content != nil {
		_, _ = view.writer.Write(conversation.Content)
	}
	c.views[conversation.ChatID] = view
	c.view = view
	c.page.AddAndSwitchToPage(conversation.ChatID, view.textView, true)
}

// Primitive implements Primitive.
func (c *Chat) Primitive() tview.Primitive {
	return c.page
}
func (c *Chat) Writer() io.Writer {
	return c.view.writer
}

func (c *Chat) SwitchView(chatId string) bool {
	view, ok := c.views[chatId]
	if ok {
		c.view = view
		c.page.SwitchToPage(chatId)
	}
	return ok
}

func newTextView(title string, onChangeFunc func()) *view {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetChangedFunc(func() {
		onChangeFunc()
		textView.ScrollToEnd()
	})
	writer := tview.ANSIWriter(textView)

	return &view{
		textView: textView,
		writer:   writer,
	}
}
