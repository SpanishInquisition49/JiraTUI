package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/SpanishInquisition49/JiraTUI/internal/jira"
)

type IssueList struct {
	style         lipgloss.Style
  titleStyle    lipgloss.Style
	issuesList    list.Model
	issues        []jira.Issue
	selectedIssue *jira.Issue
}

func NewIssueList() IssueList {
	// Issues list
	items := []list.Item{}
	list := list.New(
		items,
		itemDelegate{},
		15,
		20,
	)

  sp := spinner.New()
  sp.Spinner = spinner.Dot
  sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	list.Title = "Issues"
	list.SetFilteringEnabled(false)
	list.SetShowStatusBar(false)
  list.SetSpinner(sp.Spinner)

	return IssueList{
    style:         lipgloss.NewStyle(),
    titleStyle:    lipgloss.NewStyle(),
		issuesList:    list,
		issues:        []jira.Issue{},
		selectedIssue: nil,
	}
}

/**
 * SetIssues sets the issues in the list
 * @param issues []jira.Issue - The issues to set
 */
func (il *IssueList) SetIssues(issues []jira.Issue) {
  if issues == nil {
    issues = []jira.Issue{}
  }
	il.issues = issues
	items := []list.Item{}
	for _, issue := range issues {
		items = append(items, item(issue.Key))
	}
	il.issuesList.SetItems(items)
  il.issuesList.StopSpinner()
}

func (il *IssueList) StartSpinner() {
  il.issuesList.StartSpinner()
}

/**
 * Handle the update of the component
 * @param msg Msg - The message to forward to the list
 */
func (il *IssueList) Update(msg tea.Msg) (list.Model, tea.Cmd) {
  issueList, issueCmd := il.issuesList.Update(msg)
  il.issuesList = issueList
	selectedIssueIndex := il.issuesList.Index()
	if selectedIssueIndex >= 0 && selectedIssueIndex < len(il.issues) {
		il.selectedIssue = &il.issues[selectedIssueIndex]
	} else {
		il.selectedIssue = nil
	}
  return issueList, issueCmd
}

func (il IssueList) View() string {
  return il.style.Render(il.issuesList.View())
}

func (il IssueList) GetSelectedIssue() *jira.Issue {
  if il.selectedIssue == nil {
    return nil
  } else if len(il.issues) == 0 {
    return nil
  }
  return il.selectedIssue
}

func (il *IssueList) SetStyle(style lipgloss.Style) {
  il.style = style
}

func (il *IssueList) SetTitleStyle(style lipgloss.Style) {
  il.titleStyle = style
  il.issuesList.Styles.Title = style
}
