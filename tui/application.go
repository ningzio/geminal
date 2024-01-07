package tui

import (
	"context"
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Conversation struct {
	ChatID  string
	Title   string
	Content []byte
}

type Backend interface {
	GetConversation(ctx context.Context, chatID string) (*Conversation, error)
	CreateConversation(ctx context.Context) (*Conversation, error)
	DeleteConversation(ctx context.Context, chatID string) error
	UpdateConversation(ctx context.Context, chatID, title string) error
	ListConversation(ctx context.Context) ([]*Conversation, error)

	Talk(ctx context.Context, chatID string, writer io.Writer, prompt string) error
}

func NewApplication(handler Backend) (*Application, error) {
	grid := tview.NewGrid().SetRows(-8, -2, 1).SetColumns(-2, -8)
	tApp := tview.NewApplication()
	tApp.SetRoot(grid, true).EnableMouse(true)
	app := &Application{
		backend: handler,
		app:     tApp,
		grid:    grid,
	}

	chat := NewChatTUI(func() { app.app.Draw() })
	app.chat = chat

	input := NewInputTUI(app.submitFunc())
	app.input = input

	conversations, err := app.backend.ListConversation(context.Background())
	if err != nil {
		return nil, fmt.Errorf("init application: %w", err)
	}
	history := NewHistoryTUI(app, conversations...)
	app.history = history

	app.addInputToGrid(input)
	app.addChatToGrid(chat)
	app.addHistoryToGrid(history)
	app.addShortcut()

	page := tview.NewPages()
	page.AddAndSwitchToPage("main", grid, true)

	warning := NewWarningTUI(func() { app.page.RemovePage("warning") })
	app.warning = warning

	app.page = page

	app.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			app.app.SetFocus(app.history.Primitive())
			return nil
		case tcell.KeyF2:
			app.app.SetFocus(app.input.Primitive())
			return nil
		case tcell.KeyF3:
			app.app.SetFocus(app.chat.Primitive())
			return nil
		case tcell.KeyF4:
			conv, err := app.backend.CreateConversation(context.Background())
			if err != nil {
				app.warning.SetText(err.Error())
				app.warning.SetButtons("ok")
				app.warning.SetColor(tcell.ColorRed)
				app.page.AddPage("warning", app.warning.Primitive(), true, true)
				return nil
			}
			app.history.NewHistory(conv)
			app.chat.NewChatView(conv)
		}
		return event
	})

	app.app.SetRoot(app.page, true)
	return app, nil
}

type Application struct {
	backend Backend

	app  *tview.Application
	grid *tview.Grid
	page *tview.Pages

	input   InputWidget
	chat    ChatWidget
	history HistoryWidget
	warning WarningWidget
}

/*
********************
*   *              *
*   *              *
* 1 *      2       *
*   *              *
*   ****************
*   *      3       *
********************
*       4          *
********************
 */

func (app *Application) addInputToGrid(input InputWidget) {
	app.grid.AddItem(input.Primitive(), 1, 1, 1, 1, 0, 0, true)
}

func (app *Application) addChatToGrid(chat ChatWidget) {
	app.grid.AddItem(chat.Primitive(), 0, 1, 1, 1, 0, 0, false)
}

func (app *Application) addHistoryToGrid(history HistoryWidget) {
	app.grid.AddItem(history.Primitive(), 0, 0, 2, 1, 0, 0, false)
}

func (app *Application) addShortcut() {
	view := tview.NewTextView()
	view.SetText("F1: history, F2: input, F3: chat, F4: new conversation")
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorDarkGrey)

	app.grid.AddItem(view, 2, 0, 1, 2, 0, 0, false)
}

func (app *Application) submitFunc() OnUserSubmit {
	return func(input string) {
		chatID := app.history.GetCurrentChatID()
		// new conversation
		if len(chatID) == 0 {
			conversation, err := app.backend.CreateConversation(context.Background())
			if err != nil {
				app.warning.SetText(err.Error())
				app.warning.SetButtons("ok")
				app.warning.SetColor(tcell.ColorRed)
				app.page.AddPage("warning", app.warning.Primitive(), true, true)
				return
			}
			chatID = conversation.ChatID
			app.chat.NewChatView(conversation)
			app.history.NewHistory(conversation)
		}
		go func() {
			// if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(chatID), input); err != nil {
			if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(), input); err != nil {
				app.warning.SetText(err.Error())
				app.warning.SetButtons("ok")
				app.warning.SetColor(tcell.ColorRed)
				app.page.AddPage("warning", app.warning.Primitive(), true, true)
				return
			}
		}()
	}
}

func (app *Application) OnConversationChanged(chatID string) {
	ok := app.chat.SwitchView(chatID)
	if !ok {
		// restore conversation from repository
		conversation, err := app.backend.GetConversation(context.Background(), chatID)
		if err != nil {
			app.warning.SetText(err.Error())
			app.warning.SetButtons("ok")
			app.warning.SetColor(tcell.ColorRed)
			app.page.AddPage("warning", app.warning.Primitive(), true, true)
			return
		}
		app.chat.NewChatView(conversation)
	}
}

func (app *Application) DeleteConversation(chatID string) error {
	if err := app.backend.DeleteConversation(context.Background(), chatID); err != nil {
		return err
	}
	app.chat.DeleteView(chatID)
	return nil
}
func (app *Application) RenameConversation(chatID, newTitle string) error {
	return app.backend.UpdateConversation(context.Background(), chatID, newTitle)
}

func (app *Application) Run() error {
	return app.app.Run()
}
