package tui

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func init() {
	tview.Styles.ContrastBackgroundColor = tcell.ColorMidnightBlue

	// now is default
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorGreen
	tview.Styles.BorderColor = tcell.ColorWhite
	tview.Styles.TitleColor = tcell.ColorWhite
	tview.Styles.GraphicsColor = tcell.ColorWhite
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorYellow
	tview.Styles.TertiaryTextColor = tcell.ColorGreen
	tview.Styles.InverseTextColor = tcell.ColorBlue
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorNavy
}

// NewApplication initializes a new Application with the given backend.
//
// backend: The backend to use for the Application.
// Returns a pointer to the newly created Application and an error if there was any.
func NewApplication(backend Backend) (*Application, error) {
	app := &Application{
		backend: backend,
		app:     tview.NewApplication(),
		grid:    tview.NewGrid(),
		page:    tview.NewPages(),
	}

	app.initWidget()
	app.setMainLayout()
	app.setPages()
	app.bindKeys()

	convs, err := backend.ListConversation(context.Background())
	if err != nil {
		return nil, err
	}
	for _, c := range convs {
		app.history.NewHistory(c)
	}

	app.app.SetRoot(app.page, true).EnableMouse(true)
	return app, nil
}

type Application struct {
	backend Backend

	app     *tview.Application
	grid    *tview.Grid
	page    *tview.Pages
	input   InputWidget
	chat    ChatWidget
	history HistoryWidget
	warning *Warning
}

// initWidget initializes the widget in the Application struct.
//
// It creates a new chat, input, and history widget, and assigns them to the
// corresponding fields in the Application struct.
//
// Return:
// - error: an error object if there is an error during initialization, otherwise nil.
func (app *Application) initWidget() error {
	app.chat = NewChat(func() { app.app.Draw() })
	app.input = NewInputTUI(app.submitFunc())
	app.history = NewHistoryTUI(app)

	return nil
}

/*

 */

// setMainLayout sets the main layout of the Application.
//
// It configures the grid layout of the Application's UI, adding the input,
// chat, history, and view components to their respective positions.
// The view component displays a text with shortcut keys for the user.
//
// +--------------------+
// | History |			|
// |	     |   chat   |
// |         |          |
// |         |----------|
// |         |   input  |
// |---------+----------|
// | help               |
// +--------------------+
func (app *Application) setMainLayout() {
	app.grid.SetRows(-8, -2, 1).SetColumns(-2, -8)
	app.grid.AddItem(app.input.Primitive(), 1, 1, 1, 1, 0, 0, true)
	app.grid.AddItem(app.chat.Primitive(), 0, 1, 1, 1, 0, 0, false)
	app.grid.AddItem(app.history.Primitive(), 0, 0, 2, 1, 0, 0, false)

	view := tview.NewTextView()
	view.SetText("F1: history, F2: input, F3: chat, F4: new conversation")
	view.SetDynamicColors(true)
	view.SetTextColor(tcell.ColorDarkGrey)
	app.grid.AddItem(view, 2, 0, 1, 2, 0, 0, false)
}

// setPages initializes and sets up the pages for the Application.
func (app *Application) setPages() {
	app.page = tview.NewPages()
	app.page.AddPage("main", app.grid, true, true)
	app.warning = NewWarningTUI(func() { app.page.SwitchToPage("main") })
	app.page.AddPage("warning", app.warning.Primitive(), true, false)
}

// showWarning sets the error message to the warning label, sets the button text to "ok",
// sets the color to red, and switches to the warning page.
//
// err: The error to be displayed as a warning.
func (app *Application) showWarning(err error) {
	app.warning.SetText(err.Error())
	app.warning.SetButtons("ok")
	app.warning.SetColor(tcell.ColorRed)
	app.page.SwitchToPage("warning")
}

// bindKeys binds the key events to specific actions in the Application.
func (app *Application) bindKeys() {
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
				app.showWarning(err)
				return nil
			}
			app.history.NewHistory(conv)
			app.chat.NewChatView(conv)
		case tcell.KeyTab:
			switch app.app.GetFocus() {
			case app.history.Primitive():
				app.app.SetFocus(app.chat.Primitive())
			case app.chat.Primitive():
				app.app.SetFocus(app.input.Primitive())
			case app.input.Primitive():
				app.app.SetFocus(app.history.Primitive())
			}
		}
		return event
	})
}

// submitFunc returns an OnUserSubmit function that handles user input.
//
// The function takes a string input and performs the following steps:
// 1. Retrieves the current chat ID from the application's history.
// 2. If the chat ID is empty, creates a new conversation using the application's backend.
// 3. Updates the chat ID, creates a new chat view, and initializes a new history if necessary.
// 4. Calls the backend's Talk method asynchronously to send the user input to the chat.
func (app *Application) submitFunc() OnUserSubmit {
	return func(input string) {
		chatID := app.history.GetCurrentChatID()
		// new conversation
		if len(chatID) == 0 {
			conversation, err := app.backend.CreateConversation(context.Background())
			if err != nil {
				app.showWarning(err)
				return
			}
			chatID = conversation.ChatID
			app.chat.NewChatView(conversation)
			app.history.NewHistory(conversation)
		}
		go func() {
			// if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(chatID), input); err != nil {
			if err := app.backend.Talk(context.Background(), chatID, app.chat.Writer(), input); err != nil {
				app.showWarning(err)
				return
			}
		}()
	}
}

// OnConversationChanged is a function that handles the change in conversation for the Application.
//
// It takes a chatID string as a parameter and switches the view of the chat based on the chatID.
// If the view switch fails, it restores the conversation from the repository by calling GetConversation
// on the backend. If an error occurs during the conversation retrieval, it shows a warning and returns.
// Otherwise, it creates a new chat view for the retrieved conversation.
func (app *Application) OnConversationChanged(chatID string) {
	ok := app.chat.SwitchView(chatID)
	if !ok {
		// restore conversation from repository
		conversation, err := app.backend.GetConversation(context.Background(), chatID)
		if err != nil {
			app.showWarning(err)
			return
		}
		app.chat.NewChatView(conversation)
	}
}

// DeleteConversation deletes a conversation with the given chatID.
//
// Parameters:
// - chatID: the ID of the conversation to be deleted.
//
// Returns:
// - error: an error if the conversation deletion fails.
func (app *Application) DeleteConversation(chatID string) error {
	if err := app.backend.DeleteConversation(context.Background(), chatID); err != nil {
		return err
	}
	app.chat.DeleteView(chatID)
	return nil
}

// RenameConversation renames a conversation in the Application.
//
// Parameters:
// - chatID: the ID of the conversation to be renamed.
// - newTitle: the new title for the conversation.
//
// Return:
// - error: an error if the conversation renaming fails.
func (app *Application) RenameConversation(chatID, newTitle string) error {
	return app.backend.UpdateConversation(context.Background(), chatID, newTitle)
}

// Run runs the Application.
//
// It returns an error if there was a problem running the Application.

func (app *Application) Run() error {
	return app.app.Run()
}
