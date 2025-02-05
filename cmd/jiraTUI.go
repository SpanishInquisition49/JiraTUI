package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/SpanishInquisition49/JiraTUI/internal/app"
	"github.com/SpanishInquisition49/JiraTUI/internal/jira"
)

func main() {
  fmt.Println("Starting Jira TUI...")
  godotenv.Load()
  
	client := jira.CreateClient(
		os.Getenv("JIRA_EMAIL"),
		os.Getenv("JIRA_TOKEN"),
		os.Getenv("JIRA_URL"),
	)

  app := app.NewModel(client)
  p := tea.NewProgram(app, tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v", err)
    os.Exit(1)
  }


}
