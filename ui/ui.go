package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type uiState int

const (
	title = iota
	navigator
)

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == title {
		return m.UpdateTitle(msg)
	} else {
		return m.UpdateNavigator(msg)
	}
}

func (m mainModel) View() string {
	var s string

	switch m.state {
	case title:
		s = TitleView(m)
	case navigator:
		s = baseStyle.Render(m.table.View()) + "\n"
		s += fmt.Sprintf("%v %v %v \n", m.currentPathIdx, m.path, m.menuCursor)
	}

	return s
}

func initialModel() mainModel {
	m := mainModel{
		state: title,
	}

	m = initialNavigatorModel(m)

	return m
}

func Start() {
	m := initialModel()

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
