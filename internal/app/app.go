package app

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/SpanishInquisition49/JiraTUI/internal/jira"
)

type (
	status    uint8
	issuesMsg struct{ issues []jira.Issue }
)

const (
	StatusDefault status = iota
	StatusSearch
	StatusIssueDetail
	StatusComment
)

type Styles struct {
	BorderColor        lipgloss.Color
	FocusedBorderColor lipgloss.Color
}

type model struct {
	state       status
	width       int
	height      int
	style       AppStyles
	jiraClient  *jira.Client
	searchInput IssueQuery
	issuesList  IssueList
	detailCard  IssueCard
	isStacked   bool
}

func NewModel(jiraClient *jira.Client) *model {
	var s AppStyles = DefaultStyles()

	il := NewIssueList()
	il.SetStyle(s.FocusedStyle)
	il.SetTitleStyle(s.ListTitleStyle)
	ic := NewIssueCard()
	ic.SetStyle(s.DefaultStyle)
	ic.SetTitleStyle(s.CardTitleStyle)
	ic.SetLabelStyle(s.CardLabelStyle)
	ic.SetValueStyle(s.CardValueStyle)

	jql := os.Getenv("JIRA_DEFAULT_JQL")
	si := NewIssueQuery(jql)
	si.SetStyle(s.DefaultStyle)

	return &model{
		state:       StatusDefault,
		style:       s,
		jiraClient:  jiraClient,
		searchInput: si,
		issuesList:  il,
		detailCard:  ic,
	}
}

func (m model) Init() tea.Cmd {
	if m.searchInput.Value() != "" {
		// TODO: add the start spinner command
		return tea.Batch(
			searchIssues(&m),
		)
	}
	return nil
}

func searchIssues(m *model) tea.Cmd {
  m.issuesList.StartSpinner()
	return func() tea.Msg {
		issues := m.jiraClient.SearchIssues(m.searchInput.Value())
		return issuesMsg{issues}
	}
}

func (m *model) handleEnter() tea.Cmd {
	switch m.state {
	case StatusIssueDetail:
		m.ChangeStatus(StatusDefault)
		return nil
	case StatusDefault:
		m.ChangeStatus(StatusIssueDetail)
	case StatusSearch:
		m.ChangeStatus(StatusDefault)
		return searchIssues(m)
	default:
		return nil
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	commands := []tea.Cmd{}
	// Update the search input
	cmd = m.searchInput.Update(msg)
	commands = append(commands, cmd)
	if m.state != StatusIssueDetail {
		_, cmd = m.issuesList.Update(msg)
		commands = append(commands, cmd)
	}
	m.detailCard.SetIssue(m.issuesList.GetSelectedIssue())
	_, _, cmd = m.detailCard.Update(msg)
	commands = append(commands, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		resize(&m)
	case issuesMsg:
		m.issuesList.SetIssues(msg.issues)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.state != StatusSearch {
				return m, tea.Quit
			}
		case "ctrl-c":
			return m, tea.Quit
		case "enter":
      return m, m.handleEnter()
		case "esc":
			if m.state == StatusComment {
				m.ChangeStatus(StatusIssueDetail)
				return m, nil
			}
			if m.state == StatusSearch {
				m.ChangeStatus(StatusDefault)
				return m, nil
			}
		case "/":
			if m.state != StatusSearch {
				m.ChangeStatus(StatusSearch)
				return m, nil
			}
		}
	}

	return m, tea.Batch(commands...)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string

	if m.isStacked {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			m.issuesList.View(),
			m.detailCard.View(),
		)
	} else {
		content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.issuesList.View(),
			m.detailCard.View(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.searchInput.View(),
		content,
	)
}

func (m *model) ChangeStatus(newStatus status) {
	// Reset
	m.searchInput.SetStyle(m.style.DefaultStyle)
	m.issuesList.SetStyle(m.style.DefaultStyle)
	m.detailCard.SetStyle(m.style.DefaultStyle)
	m.searchInput.Blur()
	switch newStatus {
	case StatusSearch:
		m.searchInput.SetStyle(m.style.FocusedStyle)
		m.searchInput.Focus()
	case StatusIssueDetail:
		m.detailCard.SetStyle(m.style.FocusedStyle)
	case StatusDefault:
		m.issuesList.SetStyle(m.style.FocusedStyle)
	}
  m.state = newStatus
}

func resize(m *model) {
	// Decide layout: side-by-side or stacked
	m.isStacked = m.width <= 80
	// Set responsive dimensions for the input field
	// m.searchInput.Width = m.width - 7 // Leave padding for borders

	// Adjust the issue list and issue card dynamically
	//availableHeight := m.height - 4 // Account for input field and padding
	//listHeight := availableHeight / 2
	//if availableHeight > 10 {
	//listHeight = availableHeight * 2 / 3 // Allocate more space to the list
	//}

	// Allocate proportions for wide layout
	// listWidth := m.width / 3 // Give one-third of the width to the issue list
	// cardWidth := m.width - listWidth  // Remaining space for the issue card

	// TODO: change the size of the components

	// m.issuesList.SetHeight(listHeight)
	// m.issuesList.SetWidth(listWidth)
}
