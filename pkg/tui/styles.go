package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	colorPrimary = lipgloss.Color("#00D9FF")
	colorSuccess = lipgloss.Color("#00FF87")
	colorWarning = lipgloss.Color("#FFD700")
	colorError   = lipgloss.Color("#FF5F87")
	colorGhost   = lipgloss.Color("#AF87FF")
	colorMuted   = lipgloss.Color("#666666")
	colorBorder  = lipgloss.Color("#333333")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 2).
			MarginBottom(1)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1).
			MarginRight(1)

	leftPanelStyle = panelStyle.Copy().
			Width(40)

	rightPanelStyle = panelStyle.Copy()

	// Step list styles
	stepStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(1)

	stepPendingStyle = stepStyle.Copy().
				Foreground(colorMuted)

	stepRunningStyle = stepStyle.Copy().
				Foreground(colorWarning).
				Bold(true)

	stepCompleteStyle = stepStyle.Copy().
				Foreground(colorSuccess)

	stepFailedStyle = stepStyle.Copy().
			Foreground(colorError)

	stepGhostStyle = stepStyle.Copy().
			Foreground(colorGhost).
			Bold(true)

	// Output styles
	outputStyle = lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorMuted).
			MarginBottom(1)

	ghostOutputStyle = outputStyle.Copy().
				BorderForeground(colorGhost)

	// Progress bar
	progressBarStyle = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Background(lipgloss.Color("#2a2a2a")).
				Padding(0, 1)

	// Labels
	labelStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	// Consensus/Dissent
	consensusStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	dissentStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)

	// Provider badges
	providerOpenAI = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#74aa9c")).
			Bold(true)

	providerAnthropic = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#d4a574")).
				Bold(true)

	providerGoogle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4285f4")).
			Bold(true)
)

// GetProviderStyle returns styled provider name
func GetProviderStyle(provider string) lipgloss.Style {
	switch provider {
	case "openai":
		return providerOpenAI
	case "anthropic":
		return providerAnthropic
	case "google":
		return providerGoogle
	default:
		return valueStyle
	}
}

// GetStatusIcon returns icon for step status
func GetStatusIcon(status string) string {
	switch status {
	case "pending":
		return "⏳"
	case "running":
		return "⚡"
	case "complete":
		return "✓"
	case "failed":
		return "✗"
	default:
		return "○"
	}
}

// RenderProgressBar creates a visual progress bar
func RenderProgressBar(progress float64, width int) string {
	filled := int(progress / 100 * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return progressBarStyle.Render(bar)
}
