package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/thomasjinlo/gochatter/internal/client"
)

func TviewRender(c *client.Client) {
    app := tview.NewApplication()
    var currChannel *client.Channel
    var currChatView *tview.TextView
    var channelSet bool

    chatContainer := tview.NewFlex()
    channelListView := tview.NewList()
    channelSocket := c.Connect()

    for _, channel := range c.GetChannels() {
        channelSocket.JoinChannel(channel.Id, c.DisplayName)
        chatView := tview.NewTextView()

        if !channelSet {
            currChannel = channel
            currChatView = chatView
            chatContainer.AddItem(currChatView, 0, 5, false)
            channelSet = true
        }

        channelSocket.HandleMessageReceive(channel.Id, func(sender, message string) {
            messageBytes := []byte(fmt.Sprintf("%s: %s\n", sender, message))
            chatView.Write(messageBytes)
            app.Draw()
        })

        channelSocket.HandleUserJoin(channel.Id, func(message string) {
            chatView.Write([]byte(message))
            app.Draw()
        })

        channelListView.AddItem(channel.Name, string(len(channel.Users)), 0, func() {
            currChannel = channel
            chatContainer.RemoveItem(currChatView)
            chatContainer.AddItem(chatView, 0, 1, false)
            currChatView = chatView
        })
    }

    chatFlex := tview.NewFlex().SetDirection(tview.FlexRow)
    chatFlex.AddItem(chatContainer, 0, 5, false)
    inputView := tview.NewInputField()
    inputView.SetDoneFunc(func(key tcell.Key){
        switch key {
        case tcell.KeyEnter:
            message := inputView.GetText()
            currChatView.Write([]byte("You: " + message + "\n"))
            inputView.SetText("")
            channelSocket.Broadcast(currChannel.Id, c.DisplayName, message)
        }
    })

    chatFlex.AddItem(inputView, 0, 5, false)

    clientFlex := tview.NewFlex()

    clientFlex.AddItem(channelListView, 0, 1, false)
    clientFlex.AddItem(chatFlex, 0, 5, false)

    if err := app.SetRoot(clientFlex, true).EnableMouse(true).SetFocus(inputView).Run(); err != nil {
        panic(err)
    }
}
