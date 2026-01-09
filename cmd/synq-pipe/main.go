package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/onoja123/synq-pipe/pkg/parser"
	"github.com/onoja123/synq-pipe/pkg/pipeline"
	"github.com/onoja123/synq-pipe/pkg/tui"
	"github.com/onoja123/synqly-go/pkg/synqly"
)

const (
	version = "1.0.0"
	banner  = `
╔═══════════════════════════════════════╗
║   Synqly Multiverse CLI v%s       ║
║   Multi-Model Workflow Engine         ║
╚═══════════════════════════════════════╝
`
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		runWorkflow()
	case "version":
		fmt.Printf("synq-pipe v%s\n", version)
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func runWorkflow() {
	if len(os.Args) < 3 {
		fmt.Println("Error: workflow file required")
		fmt.Println("Usage: synq-pipe run <workflow.yaml> [--tui]")
		os.Exit(1)
	}

	workflowFile := os.Args[2]
	useTUI := false

	// Check for --tui flag
	for _, arg := range os.Args[3:] {
		if arg == "--tui" {
			useTUI = true
		}
	}

	// Get API key from environment
	apiKey := os.Getenv("SYNQLY_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: SYNQLY_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Print banner
	fmt.Printf(banner, version)
	fmt.Printf("Loading workflow: %s\n\n", workflowFile)

	// Parse workflow
	workflow, err := parser.ParseWorkflow(workflowFile)
	if err != nil {
		fmt.Printf("Error parsing workflow: %v\n", err)
		os.Exit(1)
	}

	// Create Synqly client
	client := synqly.NewClient(synqly.Config{
		APIKey: apiKey,
	})

	// Convert steps to pointer slice
	steps := make([]*pipeline.Step, len(workflow.Steps))
	for i := range workflow.Steps {
		steps[i] = &workflow.Steps[i]
	}

	// Create pipeline
	p := pipeline.NewPipeline(workflow.Name, workflow.Description, steps, client)

	if useTUI {
		// Run with TUI
		runWithTUI(p, workflow.Input)
	} else {
		// Run without TUI (simple output)
		runSimple(p, workflow.Input)
	}
}

func runWithTUI(p *pipeline.Pipeline, input string) {
	// Start pipeline in goroutine
	go func() {
		if err := p.Execute(input); err != nil {
			fmt.Fprintf(os.Stderr, "Pipeline error: %v\n", err)
		}
	}()

	// Run TUI
	model := tui.NewModel(p)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

func runSimple(p *pipeline.Pipeline, input string) {
	fmt.Println("Executing pipeline...")
	fmt.Println()

	// Execute with simple progress updates
	go func() {
		for update := range p.Updates() {
			step := update.Step
			fmt.Printf("[%d] %s: %s\n",
				update.StepIndex+1,
				step.Name,
				step.Status)

			if step.Status == pipeline.StepStatusComplete {
				if step.IsGhostMode() {
					fmt.Printf("    👻 Ghost mode: %d providers queried\n", len(step.GhostOutputs))
				}
				fmt.Printf("    Duration: %s\n", step.Duration())
			}

			if step.Error != nil {
				fmt.Printf("    Error: %v\n", step.Error)
			}
			fmt.Println()
		}
	}()

	if err := p.Execute(input); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Pipeline completed successfully!")

	// Print final output
	if len(p.Steps) > 0 {
		lastStep := p.Steps[len(p.Steps)-1]
		fmt.Println("\nFinal Output:")
		fmt.Println("─────────────")
		fmt.Println(lastStep.Output)
	}
}

func printHelp() {
	help := `
Synqly Multiverse CLI - Multi-Model Workflow Engine

USAGE:
    synq-pipe <command> [options]

COMMANDS:
    run <workflow.yaml>     Execute a workflow
        --tui               Run with interactive TUI
    version                 Show version
    help                    Show this help

ENVIRONMENT:
    SYNQLY_API_KEY         Your Synqly API key (required)

EXAMPLES:
    # Run workflow with TUI
    synq-pipe run workflow.yaml --tui

    # Run workflow with simple output
    synq-pipe run workflow.yaml

WORKFLOW YAML FORMAT:
    name: My Workflow
    description: Example workflow
    input: Initial prompt text
    steps:
      - name: Step 1
        provider: openai
        model: gpt-4
      
      - name: Ghost Mode Step
        mode: ghost
        providers: [openai, anthropic, google]

For more information, visit: https://synqly.com
`
	fmt.Println(help)
}
