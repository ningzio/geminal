package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Warning struct {
	modal *tview.Modal
	flex  *tview.Flex
}

func NewWarningTUI(doneFunc func()) *Warning {
	modal := tview.NewModal()
	modal.AddButtons([]string{"OK"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		doneFunc()
	})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	return &Warning{modal: modal, flex: flex}
}

func (w *Warning) SetText(message string) {
	w.modal.SetText(message)
}

func (w *Warning) SetButtons(buttons ...string) {
	w.modal.ClearButtons()
	w.modal.AddButtons(buttons)
}

func (w *Warning) SetColor(color tcell.Color) {
	w.modal.SetTextColor(color)
}

func (w *Warning) Primitive() tview.Primitive {
	return w.flex
}
