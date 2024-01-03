package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	table := tview.NewTable()
	table.SetBorder(true)
	app := tview.NewApplication().SetRoot(table, true)

	maxRow := 10

	count := 0
	for name, color := range tcell.ColorNames {
		row := count / maxRow
		col := count - row*maxRow
		count += 1
		cell := tview.NewTableCell(name).SetBackgroundColor(color).SetAlign(tview.AlignCenter)

		table.SetCell(row, col, cell)
	}

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
