package main

import (
	"time"

	"github.com/rivo/tview"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		v := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
		v.SetBorder(true)
		v.SetDynamicColors(true)
		return v
	}
	menu := newPrimitive("Menu")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		// SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, true).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	sideBar1 := newPrimitive("Side Bar1")
	// grid.RemoveItem(sideBar)
	grid.AddItem(sideBar1, 0, 0, 0, 0, 0, 0, false)
	grid.AddItem(sideBar1, 1, 2, 1, 1, 0, 100, false)

	app := tview.NewApplication()

	go func() {
		time.Sleep(time.Second * 5)
		v := main.(*tview.TextView)
		v.SetText("timeout")

		grid.RemoveItem(sideBar)
	}()

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
