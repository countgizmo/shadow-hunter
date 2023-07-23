package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"ziggytwister.com/shadow-hunter/lexer"
	"ziggytwister.com/shadow-hunter/parser"
	"ziggytwister.com/shadow-hunter/ui"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func main() {
	input := `
	{1 {:name "yo"}
   2 {:name "ho"}} `

	l := lexer.New(input)
	p := parser.New(l)

	edn := p.ParseEDN()

	ui.Render(edn)
}
