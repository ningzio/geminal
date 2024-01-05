package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// OnUserSubmit 当用户提交输入时的处理方法
type OnUserSubmit = func(input string)

// InputTUI 负责控制用户的输入
type InputTUI struct {
	textArea  *tview.TextArea
	primitive tview.Primitive
}

// Primitive implements Primitive.
func (i *InputTUI) Primitive() tview.Primitive {
	return i.primitive
}

func NewInputTUI(submitFunc OnUserSubmit) *InputTUI {
	textArea := tview.NewTextArea()
	textArea.SetPlaceholder("Type something here...")
	textArea.SetWordWrap(true)

	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)

	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfos)
	updateInfos()

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			input := textArea.GetText()
			textArea.SetText("", true)
			submitFunc(input)
			return nil
		}
		return event
	})

	grid := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textArea, 0, 0, 1, 1, 0, 0, true).
		AddItem(position, 1, 0, 1, 1, 0, 0, false)
	grid.SetBorder(true)

	return &InputTUI{
		textArea:  textArea,
		primitive: grid,
	}
}
