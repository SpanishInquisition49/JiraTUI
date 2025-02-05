package app

import (
	"cmp"
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/SpanishInquisition49/JiraTUI/internal/jira"
)

type IssueCard struct {
	style               lipgloss.Style
	dimensions          Size
	titleStyle          lipgloss.Style
	labelStyle          lipgloss.Style
	valueStyle          lipgloss.Style
	issue               *jira.Issue
	descriptionViewport viewport.Model
	commentBox          textarea.Model
}

func NewIssueCard() IssueCard {
	dv := viewport.New(0, 0)
	cm := textarea.New()
	cm.Placeholder = "Add a comment..."
	cm.SetWidth(50)
	cm.SetHeight(10)
	cm.ShowLineNumbers = false

	return IssueCard{
		style:               lipgloss.NewStyle(),
		titleStyle:          lipgloss.NewStyle(),
		labelStyle:          lipgloss.NewStyle(),
		valueStyle:          lipgloss.NewStyle(),
		issue:               nil,
		descriptionViewport: dv,
		commentBox:          cm,
	}
}

func (ic *IssueCard) SetTitleStyle(style lipgloss.Style) {
	ic.titleStyle = style
}

func (ic *IssueCard) SetStyle(style lipgloss.Style) {
	ic.style = style
}

func (ic *IssueCard) SetLabelStyle(style lipgloss.Style) {
	ic.labelStyle = style
}

func (ic *IssueCard) SetValueStyle(style lipgloss.Style) {
	ic.valueStyle = style
}

func (ic *IssueCard) SetIssue(issue *jira.Issue) {
	ic.issue = issue
}

func (ic *IssueCard) Update(msg tea.Msg) (viewport.Model, textarea.Model, tea.Cmd) {
	descriptionModel, descriptionCmd := ic.descriptionViewport.Update(msg)
	commentModel, commentCmd := ic.commentBox.Update(msg)
	ic.commentBox = commentModel
	ic.descriptionViewport = descriptionModel
	return descriptionModel, commentModel, tea.Batch(descriptionCmd, commentCmd)
}

func (ic *IssueCard) View() string {
	summary := "No Summary"
	description := "No Description"
	status := "Unknown"
	assignee := "Unassigned"
	reporter := "Unknown"

	if ic.issue != nil {
		summary = cmp.Or(ic.issue.Summary, summary)
		description = cmp.Or(ic.issue.Description, description)
		status = cmp.Or(ic.issue.Status, status)
		assignee = cmp.Or(ic.issue.Assignee, assignee)
		reporter = cmp.Or(ic.issue.Reporter, reporter)
	}

	// TODO: adjust word wrap
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(50),
	)
	description, _ = r.Render(description)
	ic.descriptionViewport.SetContent(description)

	// Card content
	card := lipgloss.JoinVertical(
		lipgloss.Left,
		ic.titleStyle.Render(summary),
		fmt.Sprintf("%s %s", ic.labelStyle.Render("Status:"), ic.valueStyle.Render(status)),
		fmt.Sprintf("%s %s", ic.labelStyle.Render("Assignee:"), ic.valueStyle.Render(assignee)),
		fmt.Sprintf("%s %s", ic.labelStyle.Render("Reporter:"), ic.valueStyle.Render(reporter)),
		fmt.Sprintf("%s", ic.labelStyle.Render("Description:")),
		ic.descriptionViewport.View(),
	)
	return ic.style.Render(card)
}
