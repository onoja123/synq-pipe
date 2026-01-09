package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/onoja123/synq-pipe/pkg/pipeline"
)

// Color styles
var (
	colorMutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
	colorWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	colorErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))
	colorSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#33FF99"))
)

// Model represents the TUI state
type Model struct {
	pipeline     *pipeline.Pipeline
	width        int
	height       int
	ready        bool
	err          error
	done         bool
	currentStep  int
	ghostMode    bool
	lastAnalysis *pipeline.GhostAnalysis
}

// StepUpdateMsg wraps pipeline updates
type StepUpdateMsg pipeline.StepUpdate

// PipelineCompleteMsg signals completion
type PipelineCompleteMsg struct{}

// PipelineErrorMsg signals error
type PipelineErrorMsg struct {
	err error
}

// NewModel creates a new TUI model
func NewModel(p *pipeline.Pipeline) Model {
	return Model{
		pipeline:    p,
		currentStep: -1,
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForStepUpdate(m.pipeline),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case StepUpdateMsg:
		update := pipeline.StepUpdate(msg)
		m.currentStep = update.StepIndex

		if update.Analysis != nil {
			m.ghostMode = true
			m.lastAnalysis = update.Analysis
		} else {
			m.ghostMode = false
		}

		// Check if pipeline is complete
		if m.isComplete() {
			return m, tea.Quit
		}

		return m, waitForStepUpdate(m.pipeline)

	case PipelineCompleteMsg:
		m.done = true
		return m, tea.Quit

	case PipelineErrorMsg:
		m.err = msg.err
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.err != nil {
		return errorView(m.err)
	}

	// Header
	header := m.renderHeader()

	// Two-panel layout
	leftPanel := m.renderLeftPanel()
	rightPanel := m.renderRightPanel()

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanelStyle.Render(leftPanel),
		rightPanelStyle.Width(m.width-45).Render(rightPanel),
	)

	// Footer
	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		panels,
		footer,
	)
}

// renderHeader creates the header section
func (m Model) renderHeader() string {
	title := fmt.Sprintf("🚀 Synqly Multiverse CLI  │  %s", m.pipeline.Name)
	if m.pipeline.Description != "" {
		title += fmt.Sprintf("  │  %s", m.pipeline.Description)
	}
	return headerStyle.Width(m.width).Render(title)
}

// renderLeftPanel creates the step list panel
func (m Model) renderLeftPanel() string {
	var steps []string

	steps = append(steps, labelStyle.Render("PIPELINE STEPS"))
	steps = append(steps, "")

	for i, step := range m.pipeline.Steps {
		stepView := m.renderStep(i, step)
		steps = append(steps, stepView)
	}

	steps = append(steps, "")
	progress := m.pipeline.GetProgress()
	progressBar := RenderProgressBar(progress, 35)
	steps = append(steps, progressBar)
	steps = append(steps, valueStyle.Render(fmt.Sprintf("%.0f%% Complete", progress)))

	return strings.Join(steps, "\n")
}

// renderStep renders a single step
func (m Model) renderStep(index int, step *pipeline.Step) string {
	icon := GetStatusIcon(string(step.Status))

	var style lipgloss.Style
	switch step.Status {
	case pipeline.StepStatusPending:
		style = stepPendingStyle
	case pipeline.StepStatusRunning:
		style = stepRunningStyle
	case pipeline.StepStatusComplete:
		style = stepCompleteStyle
	case pipeline.StepStatusFailed:
		style = stepFailedStyle
	}

	stepName := step.Name
	if step.IsGhostMode() {
		stepName = "👻 " + stepName
		style = stepGhostStyle
	}

	line := fmt.Sprintf("%s %s", icon, stepName)

	// Add duration if complete
	if step.Status == pipeline.StepStatusComplete || step.Status == pipeline.StepStatusFailed {
		duration := step.Duration().Round(time.Millisecond)
		line += colorMutedStyle.Render(fmt.Sprintf(" (%s)", duration))
	}

	return style.Render(line)
}

// renderRightPanel creates the output panel
func (m Model) renderRightPanel() string {
	if m.currentStep < 0 || m.currentStep >= len(m.pipeline.Steps) {
		return labelStyle.Render("Waiting for execution...")
	}

	step := m.pipeline.Steps[m.currentStep]

	var content []string
	content = append(content, labelStyle.Render(fmt.Sprintf("STEP %d: %s", m.currentStep+1, step.Name)))
	content = append(content, "")

	if step.IsGhostMode() && m.lastAnalysis != nil {
		content = append(content, m.renderGhostOutput(step, m.lastAnalysis))
	} else {
		content = append(content, m.renderNormalOutput(step))
	}

	return strings.Join(content, "\n")
}

// renderNormalOutput renders single provider output
func (m Model) renderNormalOutput(step *pipeline.Step) string {
	var parts []string

	// Provider info
	provider := step.Provider
	if provider == "" && len(step.GetProviderList()) > 0 {
		provider = step.GetProviderList()[0]
	}

	providerStyle := GetProviderStyle(provider)
	parts = append(parts, providerStyle.Render(fmt.Sprintf("Provider: %s", provider)))
	parts = append(parts, valueStyle.Render(fmt.Sprintf("Model: %s", step.GetModel())))
	parts = append(parts, "")

	// Output
	if step.Status == pipeline.StepStatusRunning {
		parts = append(parts, colorWarningStyle.Render("⚡ Executing..."))
	} else if step.Output != "" {
		output := step.Output
		if len(output) > 500 {
			output = output[:500] + "..."
		}
		parts = append(parts, outputStyle.Render(output))
	}

	return strings.Join(parts, "\n")
}

// renderGhostOutput renders Ghost mode comparison
func (m Model) renderGhostOutput(step *pipeline.Step, analysis *pipeline.GhostAnalysis) string {
	var parts []string

	parts = append(parts, stepGhostStyle.Render("👻 GHOST MODE ANALYSIS"))
	parts = append(parts, "")

	// Consensus
	if analysis.Consensus != "" {
		parts = append(parts, consensusStyle.Render("Consensus: ")+analysis.Consensus)
	}

	// Dissent
	if len(analysis.Dissent) > 0 {
		parts = append(parts, dissentStyle.Render("Dissent:"))
		for _, d := range analysis.Dissent {
			parts = append(parts, "  • "+d)
		}
	}
	parts = append(parts, "")

	// Individual outputs
	parts = append(parts, labelStyle.Render("PROVIDER RESPONSES:"))
	parts = append(parts, "")

	for _, result := range analysis.Results {
		providerStyle := GetProviderStyle(result.Provider)
		header := providerStyle.Render(fmt.Sprintf("▸ %s (%s)", result.Provider, result.Model))

		if result.Error != nil {
			parts = append(parts, header+" "+colorErrorStyle.Render("ERROR"))
		} else {
			output := result.Output
			if len(output) > 200 {
				output = output[:200] + "..."
			}
			parts = append(parts, header)
			parts = append(parts, ghostOutputStyle.Render(output))
		}
		parts = append(parts, "")
	}

	// Best response
	if analysis.BestProvider != "" {
		parts = append(parts, consensusStyle.Render(fmt.Sprintf("✓ Selected: %s", analysis.BestProvider)))
	}

	return strings.Join(parts, "\n")
}

// renderFooter creates the footer
func (m Model) renderFooter() string {
	if m.done {
		return colorSuccessStyle.Render("✓ Pipeline complete! Press q to exit.")
	}
	return colorMutedStyle.Render("Press q or ctrl+c to quit")
}

// isComplete checks if all steps are done
func (m Model) isComplete() bool {
	for _, step := range m.pipeline.Steps {
		if step.Status != pipeline.StepStatusComplete && step.Status != pipeline.StepStatusFailed {
			return false
		}
	}
	return true
}

// errorView renders error state
func errorView(err error) string {
	return colorErrorStyle.Render(fmt.Sprintf("Error: %v", err))
}

// waitForStepUpdate waits for pipeline updates
func waitForStepUpdate(p *pipeline.Pipeline) tea.Cmd {
	return func() tea.Msg {
		update := <-p.Updates()
		return StepUpdateMsg(update)
	}
}
