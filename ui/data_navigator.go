package ui

import (
	"net"
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
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("040"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Underline(true)

	focusedButton = activeButtonStyle.Render("Submit")
	blurredButton = buttonStyle.Render("Submit")
)

type mainModel struct {
	state          uiState
	conn           net.Conn
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

func rowHasNestedData(row []string) bool {
	rowType := row[len(row)-1]

	return rowType == "Map" || rowType == "Vector"
}

func (m *mainModel) showCurrentData() {
	data := m.getCurrentDataSlice()
	switch data := data.(type) {
	case *ast.MapElement:
		m.table = m.mapToTable(data)
	}
	m.menuCursor = 0
}

func (m *mainModel) reset() tea.Msg {
	if m.conn == nil {
		m.conn = transmitter.GetConnection(m.host, m.port)
	}
	m.edn = transmitter.GetAppDB(m.conn)

	switch element := m.edn.Elements[0].(type) {
	case *ast.MapElement:
		m.table = m.mapToTable(element)
	}

	return nil
}

func (m *mainModel) UpdateNavigator(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *mainModel) vectorToTable(v *ast.VectorElement) table.Model {
	columns := []table.Column{
		{Title: "Idx", Width: 10},
		{Title: "Value", Width: 10},
	}

	rows := []table.Row{}
	var row table.Row

	for i, element := range v.Elements {
		switch element := element.(type) {
		case *ast.VectorElement:
			row = []string{strconv.Itoa(i), m.vectorToTable(element).View()}
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

func (m *mainModel) mapToTable(mapElement *ast.MapElement) table.Model {
	columns := []table.Column{
		{Title: "Key", Width: 10},
		{Title: "Value", Width: 80},
		{Title: "Type", Width: 10},
	}

	rows := []table.Row{}
	var row table.Row

	for i, key := range mapElement.Keys {
		switch value := mapElement.Values[i].(type) {
		case *ast.MapElement:
			row = []string{key.String(), value.String(), "Map"}
		case *ast.VectorElement:
			row = []string{key.String(), m.vectorToTable(value).View(), "Vector"}
		default:
			row = []string{key.String(), value.String(), ""}
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	t.SetStyles(tableStyle)

	return t
}

func initialNavigatorModel(m mainModel) mainModel {
	m.path = []int{0}
	m.currentPathIdx = 0
	m.menuCursor = 0
	m.inputs = make([]textinput.Model, 2)

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
