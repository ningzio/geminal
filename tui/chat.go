package tui

import (
	"io"

	"github.com/rivo/tview"
)

var _ ChatWight = (*ChatTUI)(nil)

func NewChatTUI(onChangeFunc func()) *ChatTUI {
	return &ChatTUI{
		onChangeFunc: onChangeFunc,
		views:        make(map[string]*view),
	}
}

type view struct {
	textView *tview.TextView
	writer   io.Writer
}

type ChatTUI struct {
	onChangeFunc func()
	// current view
	view *view
	// ChatUI can hold multi text view, when user change
	// chat chat history in side bar, chat ui should
	// switch to correspond text view
	views map[string]*view
}

// SetTitle implements ChatWight.
func (c *ChatTUI) SetTitle(title string) {
	c.view.textView.SetTitle(title)
}

func (c *ChatTUI) NewChatView(chatId string, title string, content []byte) {
	view := newTextView(title, c.onChangeFunc)
	_, _ = view.writer.Write(content)
	c.views[chatId] = view
	c.view = view
}

// Primitive implements Primitive.
func (c *ChatTUI) Primitive() tview.Primitive {
	if c.view == nil {
		return c.emptyView().textView
	}
	return c.view.textView
}

func (c *ChatTUI) Writer() io.Writer {
	return c.view.writer
}

func (c *ChatTUI) SwitchView(chatId string) bool {
	view, ok := c.views[chatId]
	if ok {
		c.view = view
	}
	return ok
}

func (c *ChatTUI) emptyView() *view {
	return newTextView("", c.onChangeFunc)
}

func newTextView(title string, onChangeFunc func()) *view {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetChangedFunc(func() {
		onChangeFunc()
	})
	writer := tview.ANSIWriter(textView)

	return &view{
		textView: textView,
		writer:   writer,
	}
}