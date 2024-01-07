package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	dropdown()
}

func dropdown() {
	app := tview.NewApplication()

	list := tview.NewList()

	for i := range []string{"a", "b", "c"} {
		list.AddItem(fmt.Sprintf("Item %d", i+1), "", 0, nil)
	}

	page := tview.NewPages()
	page.AddPage("conversation", list, true, true)

	dropdown := tview.NewList()
	dropdown.SetDoneFunc(func() {
		page.SwitchToPage("conversation")
	})

	page.AddPage("dropdown", dropdown, true, false)

	list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		dropdown.Clear()

		dropdown.AddItem("Delete this conversation?", s2, 0, func() {
			modal := tview.NewModal()
			modal.SetText("delete this conversation?")
			modal.AddButtons([]string{"cancel", "delete"})
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonLabel {
				case "cancel":
					page.RemovePage("modal")
					page.SwitchToPage("conversation")
				case "delete":
					page.RemovePage("modal")
					list.RemoveItem(i)
					page.SwitchToPage("conversation")
				}
			})
			page.AddAndSwitchToPage("modal", modal, true)
		})
		dropdown.AddItem("Rename this conversation?", s2, 0, func() {
			input := tview.NewInputField()
			input.SetLabel("new title: ")
			input.SetDoneFunc(func(key tcell.Key) {
				switch key {
				case tcell.KeyEnter:
					list.SetItemText(i, input.GetText(), s2)
					page.RemovePage("input")
					page.SwitchToPage("conversation")
				case tcell.KeyEsc:
					page.RemovePage("input")
					page.SwitchToPage("dropdown")
				}
			})
			page.AddAndSwitchToPage("input", input, true)
		})
		page.SwitchToPage("dropdown")
	})

	app.SetRoot(page, true).EnableMouse(true).Run()
}
