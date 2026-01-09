package pipeline

import (
	"time"
)

// StepMode defines execution mode
type StepMode string

const (
	StepModeNormal StepMode = "normal"
	StepModeGhost  StepMode = "ghost"
)

// StepStatus represents step execution state
type StepStatus string

const (
	StepStatusPending  StepStatus = "pending"
	StepStatusRunning  StepStatus = "running"
	StepStatusComplete StepStatus = "complete"
	StepStatusFailed   StepStatus = "failed"
)

// Step represents a single pipeline step
type Step struct {
	Name        string   `yaml:"name"`
	Mode        StepMode `yaml:"mode"`
	Model       string   `yaml:"model"`
	Provider    string   `yaml:"provider"`
	Providers   []string `yaml:"providers"`   // For Ghost mode
	Prompt      string   `yaml:"prompt"`      // Optional custom prompt
	Temperature *float64 `yaml:"temperature"` // Optional
	MaxTokens   *int     `yaml:"max_tokens"`  // Optional

	// Runtime fields
	Status       StepStatus        `yaml:"-"`
	Output       string            `yaml:"-"`
	GhostOutputs map[string]string `yaml:"-"` // Provider -> Output for Ghost mode
	Error        error             `yaml:"-"`
	StartTime    time.Time         `yaml:"-"`
	EndTime      time.Time         `yaml:"-"`
}

// IsGhostMode checks if step uses Ghost mode
func (s *Step) IsGhostMode() bool {
	return s.Mode == StepModeGhost && len(s.Providers) > 1
}

// GetProviderList returns list of providers to query
func (s *Step) GetProviderList() []string {
	if s.IsGhostMode() {
		return s.Providers
	}
	if s.Provider != "" {
		return []string{s.Provider}
	}
	return []string{"openai"} // Default
}

// GetModel returns the model name, inferring from provider if needed
func (s *Step) GetModel() string {
	if s.Model != "" {
		return s.Model
	}
	// Default models per provider
	defaults := map[string]string{
		"openai":    "gpt-4",
		"anthropic": "claude-sonnet-4",
		"google":    "gemini-1.5-pro",
	}
	return defaults[s.Provider]
}

// Duration returns step execution time
func (s *Step) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// GhostResult represents a single provider's response in Ghost mode
type GhostResult struct {
	Provider string
	Model    string
	Output   string
	Duration time.Duration
	Error    error
}

// GhostAnalysis contains consensus analysis
type GhostAnalysis struct {
	Results      []GhostResult
	Consensus    string   // Common themes
	Dissent      []string // Differing opinions
	BestResponse string   // Selected best response
	BestProvider string
}
