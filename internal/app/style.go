package app

import "github.com/charmbracelet/lipgloss"

type AppStyles struct {
	DefaultStyle lipgloss.Style
	FocusedStyle lipgloss.Style
	ListTitleStyle   lipgloss.Style
  CardTitleStyle lipgloss.Style
  CardLabelStyle lipgloss.Style
  CardValueStyle lipgloss.Style
}

func DefaultStyles() AppStyles {
	baseStyle := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12"))

	focuedStyle := baseStyle
	focuedStyle = focuedStyle.BorderForeground(lipgloss.Color("11"))

  titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

	return AppStyles{
		DefaultStyle: baseStyle,
		FocusedStyle: focuedStyle,
    ListTitleStyle:   titleStyle,
    CardTitleStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
    CardLabelStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true),
    CardValueStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
}
