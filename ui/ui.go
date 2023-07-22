package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ziggytwister.com/shadow-hunter/ast"
)

var tableHeaderStyle = table.DefaultStyles().Header.
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	BorderBottom(true).
	Bold(false)

var tableSelectedStyle = table.DefaultStyles().Selected.
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Bold(false)

var tableStyle = table.Styles{
	Header:   tableHeaderStyle,
	Selected: tableSelectedStyle,
	Cell:     table.DefaultStyles().Cell,
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type mainModel struct {
	edn   *ast.EDN
	table table.Model
}

func (m mainModel) Init() tea.Cmd { return nil }

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func vectorToTable(v *ast.VectorElement) table.Model {
	columns := []table.Column{
		{Title: "Idx", Width: 10},
		{Title: "Value", Width: 10},
	}

	rows := []table.Row{}
	var row table.Row

	for i, element := range v.Elements {
		switch element := element.(type) {
		case *ast.VectorElement:
			row = []string{strconv.Itoa(i), vectorToTable(element).View()}
		default:
			row = []string{strconv.Itoa(i), element.String()}
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	t.SetStyles(tableStyle)

	return t
}

func mapToTable(m *ast.MapElement) table.Model {
	columns := []table.Column{
		{Title: "Key", Width: 10},
		{Title: "Value", Width: 10},
	}

	rows := []table.Row{}
	var row table.Row

	for i, key := range m.Keys {
		switch value := m.Values[i].(type) {
		case *ast.MapElement:
			row = []string{key.String(), "..."}
		case *ast.VectorElement:
			row = []string{key.String(), vectorToTable(value).View()}
		default:
			row = []string{key.String(), value.String()}
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	t.SetStyles(tableStyle)

	return t
}

func Render(edn *ast.EDN) {
	m := mainModel{edn: edn}

	for _, element := range edn.Elements {
		switch element := element.(type) {
		case *ast.MapElement:
			m.table = mapToTable(element)
		}
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
