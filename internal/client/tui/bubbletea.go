package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thomasjinlo/gochatter/internal/client"
)

// init - create channel items
// init - websocket connection
// init - receive from web server

// update - channel click
// update - msg send

// view -

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type channelItem string

func (i channelItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func BubbleTeaRender(c *client.Client) {
	var channelItems []list.Item

	channelSocket := c.Connect()
	for _, channel := range c.GetChannels() {
		channelSocket.JoinChannel(channel.Id, c.DisplayName)
		channelItems = append(channelItems, channelItem(channel.Name))
	}


	channels := list.New([]list.Item{channelItem("hello")}, list.DefaultDelegate{}, 50, 25)
	chatview := viewport.New(30, 5)
	textinput := textinput.New()

	m := Model{channels: channels, textinput: textinput, chatview: chatview}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type Model struct {
	channels list.Model
	chatview viewport.Model
	textinput textinput.Model
}


func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.channels.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.channels, cmd = m.channels.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.channels.View()
}
