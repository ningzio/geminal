package main

import (
	"time"

	"github.com/rivo/tview"
)

func loading() {
	app := tview.NewApplication()
	view := tview.NewTextView()
	view.SetChangedFunc(func() { app.Draw() })

	view.SetBorder(false)

	go func() {
		text := []string{"🌕", "🌔", "🌓", "🌒", "🌑", "🌘", "🌗", "🌖"}
		count := 0
		tick := time.NewTicker(time.Millisecond * 150)
		for range tick.C {
			view.SetText(text[count])
			count = (count + 1) % len(text)
		}
	}()

	app.SetRoot(view, true).SetFocus(view).Run()
}

func main() {
	loading()
}
