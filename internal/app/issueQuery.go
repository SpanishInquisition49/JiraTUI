package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type IssueQuery struct {
  style         lipgloss.Style
  input         textinput.Model
}

func NewIssueQuery(jql string) IssueQuery {
  input := textinput.New()
  input.Placeholder = "Search for issues..."
  if jql != "" {
    input.SetValue(jql)
  }
  return IssueQuery{
    style: lipgloss.NewStyle(),
    input: input,
  }
}

func (iq *IssueQuery) SetStyle(style lipgloss.Style) {
  iq.style = style
}

func (iq *IssueQuery) Update(msg tea.Msg) tea.Cmd {
  inputModel, cmd := iq.input.Update(msg)
  iq.input = inputModel
  return cmd
}

func (iq *IssueQuery) View() string {
  return iq.style.Render(iq.input.View())
}

func (iq *IssueQuery) Value() string {
  return iq.input.Value()
}

func (iq *IssueQuery) SetValue(value string) {
  iq.input.SetValue(value)
}

func (iq *IssueQuery) Focus() {
  iq.input.Focus()
}

func (iq *IssueQuery) Blur() {
  iq.input.Blur()
}

