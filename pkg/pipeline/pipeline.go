package pipeline

import (
	"fmt"
	"time"

	"github.com/onoja123/synqly-go/pkg/synqly"
)

// Pipeline represents a complete workflow
type Pipeline struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Steps       []*Step `yaml:"steps"`

	// Runtime
	client      *synqly.Client
	ghostExec   *GhostExecutor
	currentStep int
	updateChan  chan StepUpdate
}

// StepUpdate represents a step state change
type StepUpdate struct {
	StepIndex int
	Step      *Step
	Analysis  *GhostAnalysis // Only for Ghost mode
}

// NewPipeline creates a new pipeline instance
func NewPipeline(name, description string, steps []*Step, client *synqly.Client) *Pipeline {
	return &Pipeline{
		Name:        name,
		Description: description,
		Steps:       steps,
		client:      client,
		ghostExec:   NewGhostExecutor(client),
		updateChan:  make(chan StepUpdate, 10),
	}
}

// Execute runs the entire pipeline
func (p *Pipeline) Execute(initialInput string) error {
	if len(p.Steps) == 0 {
		return fmt.Errorf("pipeline has no steps")
	}

	currentOutput := initialInput

	for i, step := range p.Steps {
		p.currentStep = i

		// Execute step
		output, analysis, err := p.executeStep(step, currentOutput)
		if err != nil {
			step.Status = StepStatusFailed
			step.Error = err
			p.sendUpdate(i, step, analysis)
			return fmt.Errorf("step %d (%s) failed: %w", i, step.Name, err)
		}

		// Update output for next step
		currentOutput = output
		step.Output = output
		step.Status = StepStatusComplete

		// Send update to TUI
		p.sendUpdate(i, step, analysis)
	}

	return nil
}

// executeStep executes a single step
func (p *Pipeline) executeStep(step *Step, input string) (string, *GhostAnalysis, error) {
	step.Status = StepStatusRunning
	step.StartTime = time.Now()
	p.sendUpdate(p.currentStep, step, nil)

	var output string
	var analysis *GhostAnalysis
	var err error

	if step.IsGhostMode() {
		// Ghost mode execution
		analysis, err = p.ghostExec.Execute(step, input)
		if err != nil {
			step.EndTime = time.Now()
			return "", nil, err
		}

		// Store all outputs
		step.GhostOutputs = make(map[string]string)
		for _, result := range analysis.Results {
			if result.Error == nil {
				step.GhostOutputs[result.Provider] = result.Output
			}
		}

		output = analysis.BestResponse
	} else {
		// Normal mode execution
		output, err = p.executeNormalStep(step, input)
		if err != nil {
			step.EndTime = time.Now()
			return "", nil, err
		}
	}

	step.EndTime = time.Now()
	return output, analysis, nil
}

// executeNormalStep executes a single provider step
func (p *Pipeline) executeNormalStep(step *Step, input string) (string, error) {
	providers := step.GetProviderList()
	if len(providers) == 0 {
		return "", fmt.Errorf("no provider specified")
	}

	provider := providers[0]
	model := step.GetModel()

	messages := []synqly.Message{
		{Role: "user", Content: input},
	}

	if step.Prompt != "" {
		messages = []synqly.Message{
			{Role: "system", Content: step.Prompt},
			{Role: "user", Content: input},
		}
	}

	params := synqly.ChatCreateParams{
		Provider:    provider,
		Model:       model,
		Messages:    messages,
		Temperature: step.Temperature,
		MaxTokens:   step.MaxTokens,
	}

	response, err := p.client.Chat.Create(params)
	if err != nil {
		return "", err
	}

	return response.GetContent(), nil
}

// sendUpdate sends step update to TUI
func (p *Pipeline) sendUpdate(index int, step *Step, analysis *GhostAnalysis) {
	select {
	case p.updateChan <- StepUpdate{
		StepIndex: index,
		Step:      step,
		Analysis:  analysis,
	}:
	default:
		// Non-blocking send
	}
}

// Updates returns the update channel for TUI
func (p *Pipeline) Updates() <-chan StepUpdate {
	return p.updateChan
}

// CurrentStepIndex returns current executing step
func (p *Pipeline) CurrentStepIndex() int {
	return p.currentStep
}

// GetProgress returns completion percentage
func (p *Pipeline) GetProgress() float64 {
	if len(p.Steps) == 0 {
		return 0
	}

	completed := 0
	for _, step := range p.Steps {
		if step.Status == StepStatusComplete {
			completed++
		}
	}

	return float64(completed) / float64(len(p.Steps)) * 100
}
