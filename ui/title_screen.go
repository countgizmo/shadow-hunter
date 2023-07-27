package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingBottom(1).
	Width(30).
	Align(lipgloss.Center)

type titleModel struct {
	focusIndex int
	cursorMode cursor.Mode
}

func TitleView(m mainModel) string {
	var b strings.Builder
	fmt.Fprintln(&b, style.Render("Shadow Hunter"))

	b.WriteString(m.hostInput.View())
	b.WriteRune('\n')
	b.WriteString(m.portInput.View())

	button := &blurredButton
	if m.hostInput.Value() != "" && m.portInput.Value() != "" {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
