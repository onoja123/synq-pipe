package parser

import (
	"fmt"
	"os"

	"github.com/onoja123/synq-pipe/pkg/pipeline"
	"gopkg.in/yaml.v3"
)

// WorkflowFile represents the YAML structure
type WorkflowFile struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Input       string          `yaml:"input"`
	Steps       []pipeline.Step `yaml:"steps"`
}

// ParseWorkflow reads and parses a YAML workflow file
func ParseWorkflow(filepath string) (*WorkflowFile, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var workflow WorkflowFile
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate workflow
	if err := validateWorkflow(&workflow); err != nil {
		return nil, err
	}

	// Initialize step defaults
	for i := range workflow.Steps {
		step := &workflow.Steps[i]

		// Set default mode
		if step.Mode == "" {
			step.Mode = pipeline.StepModeNormal
		}

		// Initialize status
		step.Status = pipeline.StepStatusPending
	}

	return &workflow, nil
}

// validateWorkflow checks workflow structure
func validateWorkflow(wf *WorkflowFile) error {
	if wf.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(wf.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	for i, step := range wf.Steps {
		if step.Name == "" {
			return fmt.Errorf("step %d: name is required", i)
		}

		// Validate Ghost mode
		if step.Mode == pipeline.StepModeGhost {
			if len(step.Providers) < 2 {
				return fmt.Errorf("step %d (%s): Ghost mode requires at least 2 providers", i, step.Name)
			}
		} else {
			// Normal mode requires either provider or model
			if step.Provider == "" && step.Model == "" {
				return fmt.Errorf("step %d (%s): must specify provider or model", i, step.Name)
			}
		}
	}

	return nil
}

// ValidateProviders checks if providers are supported
func ValidateProviders(providers []string) error {
	supported := map[string]bool{
		"openai":    true,
		"anthropic": true,
		"google":    true,
	}

	for _, provider := range providers {
		if !supported[provider] {
			return fmt.Errorf("unsupported provider: %s", provider)
		}
	}

	return nil
}
