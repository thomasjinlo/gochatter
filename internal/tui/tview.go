package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/thomasjinlo/gochatter/internal/client"
)

func TviewRender(client *client.Client) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
        SetChangedFunc(func() {app.Draw()})
    inputView := tview.NewInputField().SetLabel("| ")

    go receiveFromServer(client, textView)
    go client.Connect()
    go func() {
        <-client.CloseCh
        app.Stop()
    }()

    inputView.SetDoneFunc(func(key tcell.Key) {
        switch key {
        case tcell.KeyEnter:
            text := inputView.GetText()
            inputView.SetText("")
            client.FromTuiCh <- text
            textView.Write([]byte("\n" + "You: " + text))
            textView.ScrollToEnd()
        }
    })

    flex := tview.NewFlex().SetDirection(tview.FlexRow)
    flex.AddItem(textView, 0, 10, false)
    flex.AddItem(inputView, 3, 2, false)

	textView.SetBorder(true)
	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(inputView).Run(); err != nil {
		panic(err)
	}
}

func receiveFromServer(client *client.Client, textView *tview.TextView) {
    for {
        select {
        case message := <-client.ToTuiCh:
            textView.Write([]byte("\n" + message))
            textView.ScrollToEnd()
        }
    }
}
