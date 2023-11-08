package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/thomasjinlo/gochatter/internal/client"
)

func BubbleteaRender(client *client.Client) {
	model := initModel(client)
	p := tea.NewProgram(model)
	//go updateModelFromServerMessages(p, client)

	go client.Connect()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
	)

type model struct {
	client *client.Client
	list list.Model
	viewport    viewport.Model
	textinput    textinput.Model
	senderStyle lipgloss.Style
	messages    []string
	err         error
}

type serverMessage string

func updateModelFromServerMessages(p *tea.Program, c *client.Client) {
	for {
		message := <-c.ToTuiCh
		p.Send(serverMessage(message))
	}
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func initModel(client *client.Client) model {
	ta := textinput.New()
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	// Remove cursor line styling
	ta.PromptStyle = lipgloss.NewStyle()

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
		Type a message and press Enter to send.`)

	
	var channels []list.Item
	for _, channel := range client.GetChannels() {
		channels = append(channels, item{title: channel.Id})
	}

	list := list.New(channels, list.NewDefaultDelegate(), 0, 0)
	return model{
		client: client,
		list: list,
		textinput:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		liCmd tea.Cmd
	)

	m.textinput, tiCmd = m.textinput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.list, liCmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textinput.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.client.FromTuiCh <- m.textinput.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textinput.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textinput.Reset()
			m.viewport.GotoBottom()
		}

	case serverMessage:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()


	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, liCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s\n\n%s",
		m.list.View(),
		m.viewport.View(),
		m.textinput.View(),
	) + "\n\n"
}
