package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("#BA0D35")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#BA0D35")).
	PaddingTop(2).
	PaddingBottom(1).
	PaddingRight(4).
	PaddingLeft(4)

type ErrorModel struct {
	displayError string
}

func NewErrorModel(
	displayError string,
) *ErrorModel {
	m := &ErrorModel{
		displayError: displayError,
	}

	return m
}

func (m *ErrorModel) Init() tea.Cmd {
	return nil
}

func (m *ErrorModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m *ErrorModel) View() string {
	return style.Render(m.displayError)
}
