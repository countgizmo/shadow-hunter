package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	width = 96

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	titleStyle = lipgloss.NewStyle().
			Padding(1).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#874BFD"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

type titleModel struct {
	focusIndex int
	cursorMode cursor.Mode
}

func TitleView(m mainModel) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	var b strings.Builder

	// title
	{
		title := lipgloss.Place(width, 1,
			lipgloss.Center, lipgloss.Center,
			titleStyle.Render("Shadow Hunter"),
		)
		b.WriteString(title + "\n\n")
	}

	submitButtonStyle := buttonStyle
	if m.focusIndex == len(m.inputs) {
		submitButtonStyle = activeButtonStyle
	}

	question := lipgloss.NewStyle().Width(50).PaddingBottom(1).Align(lipgloss.Center).Render("PREPL connection:")
	inputs := lipgloss.JoinVertical(lipgloss.Left, m.inputs[0].View(), m.inputs[1].View())
	button := submitButtonStyle.Render("Submit")
	form := lipgloss.JoinVertical(lipgloss.Left, question, inputs)
	ui := lipgloss.JoinVertical(lipgloss.Center, form, button)

	dialog := lipgloss.Place(width, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	fmt.Fprintf(&b, dialog+"\n\n")

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(b.String())
}

func (m *mainModel) UpdateTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, connect to the REPL and switch to navigation screen.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.host = m.inputs[0].Value()
				m.port = m.inputs[1].Value()
				m.reset()
				m.state = navigator
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *mainModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
