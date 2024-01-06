package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type WarningTUI struct {
	modal *tview.Modal
	flex  *tview.Flex
}

func NewWarningTUI(doneFunc func()) *WarningTUI {
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

	return &WarningTUI{modal: modal, flex: flex}
}

func (w *WarningTUI) SetText(message string) {
	w.modal.SetText(message)
}

func (w *WarningTUI) SetButtons(buttons ...string) {
	w.modal.ClearButtons()
	w.modal.AddButtons(buttons)
}

func (w *WarningTUI) SetColor(color tcell.Color) {
	w.modal.SetTextColor(color)
}

func (w *WarningTUI) Primitive() tview.Primitive {
	return w.flex
}
