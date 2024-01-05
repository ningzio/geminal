package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type Shortcut struct {
	textView *tview.TextView
}

func NewShortcut(shortcut string) *Shortcut {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetTextAlign(tview.AlignLeft)
	view.SetText(fmt.Sprintf("press %s to focus", shortcut))

	return &Shortcut{textView: view}
}

func (s *Shortcut) Primitive() tview.Primitive {
	return s.textView
}
