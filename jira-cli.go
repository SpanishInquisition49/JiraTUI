package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type Styles struct {
	BorderColor       lipgloss.Color
	FocuedBorderColor lipgloss.Color
	inputField        lipgloss.Style
	issueList         lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("12")
	s.FocuedBorderColor = lipgloss.Color("11")
	s.inputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.RoundedBorder()).Padding(1).Width(80)
	s.issueList = s.inputField
	return s
}

type model struct {
	width         int
	height        int
	jqlField      textinput.Model
	jiraClient    *jira.Client
	styles        *Styles
	status        status
	issuesList    list.Model
	issues        []jira.Issue
	selectedIssue *jira.Issue
}

type item string

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("10")).
				Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i item) FilterValue() string {
	return string(i)
}

type status int

type (
	issuesMsg struct{ issues []jira.Issue }
	errorMsg  struct{ err error }
)

const (
	StatusDefault status = iota
	StatusSearch
)

func New(jiraClient *jira.Client) *model {
	styles := DefaultStyles()
	jqlField := textinput.New()
	jqlField.Placeholder = "Enter JQL"

	// Initialize the issue list component
	// The list will be initially empty
	items := []list.Item{}
	issuesList := list.New(
		items,
		itemDelegate{},
		15,
		20,
	)
	issuesList.Title = "Issues"
	issuesList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	issuesList.SetFilteringEnabled(false)
	issuesList.SetShowStatusBar(false)

	return &model{
		jiraClient:    jiraClient,
		jqlField:      jqlField,
		styles:        styles,
		status:        StatusDefault,
		issuesList:    issuesList,
		issues:        []jira.Issue{},
		selectedIssue: nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func ToggleSearch(model *model) {
	if model.status == StatusSearch {
		model.status = StatusDefault
		model.styles.inputField = model.styles.inputField.BorderForeground(model.styles.BorderColor)
		model.jqlField.Blur()
	} else {
		model.status = StatusSearch
		model.styles.inputField = model.styles.inputField.BorderForeground(model.styles.FocuedBorderColor)
		model.jqlField.Focus()
	}
}

func Resize(m *model) {
	// Set responsive dimensions for the input field
	m.jqlField.Width = m.width - 4 // Leave padding for borders

	// Adjust the issue list and issue card dynamically
	availableHeight := m.height - 4 // Account for input field and padding
	listHeight := availableHeight / 2
	if availableHeight > 10 {
		listHeight = availableHeight * 2 / 3 // Allocate more space to the list
	}

	// Allocate proportions for wide layout
	listWidth := m.width / 3 // Give one-third of the width to the issue list
	// cardWidth := m.width - listWidth  // Remaining space for the issue card

	m.issuesList.SetHeight(listHeight)
	m.issuesList.SetWidth(listWidth)
	m.styles.issueList.Width(listWidth) // Adjust issue list style width
}

func IssueCard(issue *jira.Issue) string {
	// Styles for different parts of the card
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true).
		Padding(0, 1)

	fieldLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Bold(true)

	fieldValueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(50)

	// Extract issue details
	summary := "No Summary"
	status := "Unknown"
	assignee := "Unassigned"
	reporter := "Unknown Reporter"

	if issue != nil {
		if issue.Fields.Summary != "" {
			summary = issue.Fields.Summary
		}
		if issue.Fields.Status != nil && issue.Fields.Status.Name != "" {
			status = issue.Fields.Status.Name
		}
		if issue.Fields.Assignee != nil && issue.Fields.Assignee.DisplayName != "" {
			assignee = issue.Fields.Assignee.DisplayName
		}
		if issue.Fields.Reporter != nil && issue.Fields.Reporter.DisplayName != "" {
			reporter = issue.Fields.Reporter.DisplayName
		}
	}

	// Create the card content
	cardContent := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(summary),
		fmt.Sprintf("%s %s", fieldLabelStyle.Render("Status:"), fieldValueStyle.Render(status)),
		fmt.Sprintf("%s %s", fieldLabelStyle.Render("Assignee:"), fieldValueStyle.Render(assignee)),
		fmt.Sprintf("%s %s", fieldLabelStyle.Render("Reporter:"), fieldValueStyle.Render(reporter)),
	)

	// Apply card styling
	return cardStyle.Render(cardContent)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		log.Printf("Width: %d\tHeight: %d\n", m.width, m.height)
		Resize(&m)
	case issuesMsg:
		// Update the issue list with the new issues
		m.issues = msg.issues
		items := []list.Item{}
		for _, issue := range msg.issues {
			items = append(items, item(issue.Key))
		}
		m.issuesList.SetItems(items)

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.status != StatusSearch {
				return m, tea.Quit
			}
		case "ctrl-c":
			return m, tea.Quit
		case "enter":
			if m.status == StatusSearch {
				cmd = searchIssues(m, m.jqlField.Value())
				ToggleSearch(&m)
				return m, cmd
			}
		case "/":
			ToggleSearch(&m)
			return m, nil
		}
	}

	// Update the models
	if m.status == StatusSearch {
		m.jqlField, cmd = m.jqlField.Update(msg)
	} else {
		m.issuesList, cmd = m.issuesList.Update(msg)
	}
	selectedIssueIndex := m.issuesList.Index()
	log.Printf("Selected Issue Index: %d\n", selectedIssueIndex)
	if selectedIssueIndex >= 0 && selectedIssueIndex < len(m.issues) {
		m.selectedIssue = &m.issues[selectedIssueIndex]
		log.Printf("Selected Issue: %s\n", m.selectedIssue.Key)
	} else {
		m.selectedIssue = nil
	}
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Decide layout: side-by-side or stacked
	isWide := m.width > 80
	var content string

	if isWide {
		// Side-by-side layout with more space for issue card
		content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.styles.issueList.Render(m.issuesList.View()),
			lipgloss.NewStyle().Width(m.width-m.issuesList.Width()).
				Render(IssueCard(m.selectedIssue)), // Adjust the width dynamically
		)
	} else {
		// Stacked layout
		content = lipgloss.JoinVertical(
			lipgloss.Top,
			m.styles.issueList.Render(m.issuesList.View()),
			IssueCard(m.selectedIssue),
		)
	}

	// Combine the input field and content
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.styles.inputField.Render(m.jqlField.View()),
		content,
	)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_EMAIL"),
		Password: os.Getenv("JIRA_TOKEN"),
	}
	jiraClient, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		log.Fatal(err)
	}

	m := New(jiraClient)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func searchIssues(m model, jql string) tea.Cmd {
	return func() tea.Msg {
		issues, _, err := m.jiraClient.Issue.Search(jql, nil)
		if err != nil {
			return err
		}
		return issuesMsg{issues: issues}
	}
}