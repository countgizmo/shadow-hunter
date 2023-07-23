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
	edn            *ast.EDN
	table          table.Model
	menuCursor     int
	path           []int
	currentPathIdx int
}

func (m mainModel) getCurrentDataSlice() ast.Element {
	result := m.edn.Elements[0]

	for i := 1; i <= m.currentPathIdx; i++ {
		switch element := result.(type) {
		case *ast.MapElement:
			result = element.Values[m.path[i]]
		case *ast.VectorElement:
			result = element.Elements[m.path[i]]
		}
	}

	return result
}

func (m mainModel) Init() tea.Cmd { return nil }

func rowHasNestedData(row []string) bool {
	rowType := row[len(row)-1]

	return rowType == "Map" || rowType == "Vector"
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	maxHeight := len(m.table.Rows())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "down", "j":
			if m.menuCursor < maxHeight {
				m.menuCursor++
			}
		case "enter":
			row := m.table.Rows()[m.menuCursor]
			if rowHasNestedData(row) {
				m.path = append(m.path, m.menuCursor)
				m.currentPathIdx++
			}
		case "backspace":
			if m.currentPathIdx > 0 {
				m.path = m.path[:len(m.path)-1]
				m.currentPathIdx--
			}
		}
	}

	data := m.getCurrentDataSlice()
	switch data := data.(type) {
	case *ast.MapElement:
		m.table = mapToTable(data)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	s := baseStyle.Render(m.table.View()) + "\n"
	s += fmt.Sprintf("%v %v\n", m.currentPathIdx, m.path)

	return s
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
		{Title: "Value", Width: 30},
		{Title: "Type", Width: 10},
	}

	rows := []table.Row{}
	var row table.Row

	for i, key := range m.Keys {
		switch value := m.Values[i].(type) {
		case *ast.MapElement:
			row = []string{key.String(), value.String(), "Map"}
		case *ast.VectorElement:
			row = []string{key.String(), vectorToTable(value).View(), "Vector"}
		default:
			row = []string{key.String(), value.String(), ""}
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
	m := mainModel{edn: edn, path: []int{0}, currentPathIdx: 0}

	// NOTE(Evgheni): I assume the data starts with a single root element

	switch element := edn.Elements[0].(type) {
	case *ast.MapElement:
		m.table = mapToTable(element)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
