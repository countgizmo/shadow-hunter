package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/transmitter"
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

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type uiState int

const (
	title = iota
	navigator
)

type mainModel struct {
	state          uiState
	edn            *ast.EDN
	table          table.Model
	menuCursor     int
	path           []int
	currentPathIdx int
	host           string
	port           string
	focusIndex     int
	inputs         []textinput.Model
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

func (m *mainModel) Init() tea.Cmd {
	//m.reset()
	return nil
}

func rowHasNestedData(row []string) bool {
	rowType := row[len(row)-1]

	return rowType == "Map" || rowType == "Vector"
}

func (m *mainModel) showCurrentData() {
	data := m.getCurrentDataSlice()
	switch data := data.(type) {
	case *ast.MapElement:
		m.table = mapToTable(data)
	}
	m.menuCursor = 0
}

func (m *mainModel) reset() tea.Msg {
	m.edn = transmitter.GetAppDB(m.host, m.port)

	switch element := m.edn.Elements[0].(type) {
	case *ast.MapElement:
		m.table = mapToTable(element)
	}

	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == title {
		return m.UpdateTitle(msg)
	}

	var cmd tea.Cmd
	maxHeight := len(m.table.Rows())
	previousPathIdx := m.currentPathIdx

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
		case "r":
			m.reset()
			m.showCurrentData()
		}

	}

	if previousPathIdx != m.currentPathIdx {
		m.showCurrentData()
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	var s string

	switch m.state {
	case title:
		s = TitleView(m)
	case navigator:
		s = baseStyle.Render(m.table.View()) + "\n"
		s += fmt.Sprintf("%v %v %v\n", m.currentPathIdx, m.path, m.menuCursor)
	}

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

func initialModel() mainModel {
	m := mainModel{
		state:          title,
		path:           []int{0},
		currentPathIdx: 0,
		menuCursor:     0,
		inputs:         make([]textinput.Model, 2)}

	m.inputs[0] = textinput.New()
	m.inputs[0].Cursor.Style = cursorStyle
	m.inputs[0].CharLimit = 50
	m.inputs[0].Placeholder = "host"
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle

	m.inputs[1] = textinput.New()
	m.inputs[1].Cursor.Style = cursorStyle
	m.inputs[1].CharLimit = 30
	m.inputs[1].Placeholder = "port"

	return m
}

func Start() {
	m := initialModel()

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
