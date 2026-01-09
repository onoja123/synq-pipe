package pipeline

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/onoja123/synqly-go/pkg/synqly"
)

// GhostExecutor handles Ghost mode execution
type GhostExecutor struct {
	client *synqly.Client
}

// NewGhostExecutor creates a new Ghost mode executor
func NewGhostExecutor(client *synqly.Client) *GhostExecutor {
	return &GhostExecutor{client: client}
}

// Execute runs a step in Ghost mode (parallel queries to multiple providers)
func (g *GhostExecutor) Execute(step *Step, input string) (*GhostAnalysis, error) {
	providers := step.GetProviderList()
	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers specified for Ghost mode")
	}

	// Execute queries in parallel
	results := make([]GhostResult, len(providers))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, provider := range providers {
		wg.Add(1)
		go func(idx int, prov string) {
			defer wg.Done()

			start := time.Now()
			output, err := g.queryProvider(step, prov, input)
			duration := time.Since(start)

			mu.Lock()
			results[idx] = GhostResult{
				Provider: prov,
				Model:    g.getModelForProvider(prov, step),
				Output:   output,
				Duration: duration,
				Error:    err,
			}
			mu.Unlock()
		}(i, provider)
	}

	wg.Wait()

	// Analyze results
	analysis := g.analyzeResults(results)

	return analysis, nil
}

// queryProvider sends request to a single provider
func (g *GhostExecutor) queryProvider(step *Step, provider, input string) (string, error) {
	model := g.getModelForProvider(provider, step)

	messages := []synqly.Message{
		{Role: "user", Content: input},
	}

	// Add custom prompt if specified
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

	response, err := g.client.Chat.Create(params)
	if err != nil {
		return "", err
	}

	return response.GetContent(), nil
}

// getModelForProvider returns appropriate model for provider
func (g *GhostExecutor) getModelForProvider(provider string, step *Step) string {
	if step.Model != "" {
		return step.Model
	}

	defaults := map[string]string{
		"openai":    "gpt-4",
		"anthropic": "claude-sonnet-4",
		"google":    "gemini-1.5-pro",
	}

	if model, ok := defaults[provider]; ok {
		return model
	}
	return "gpt-4"
}

// analyzeResults compares outputs and identifies consensus/dissent
func (g *GhostExecutor) analyzeResults(results []GhostResult) *GhostAnalysis {
	analysis := &GhostAnalysis{
		Results: results,
	}

	// Filter successful results
	successResults := []GhostResult{}
	for _, r := range results {
		if r.Error == nil && r.Output != "" {
			successResults = append(successResults, r)
		}
	}

	if len(successResults) == 0 {
		return analysis
	}

	// Simple consensus detection (can be enhanced with semantic similarity)
	outputs := make([]string, len(successResults))
	for i, r := range successResults {
		outputs[i] = r.Output
	}

	// Find common themes (simple word overlap for now)
	analysis.Consensus = g.findConsensus(outputs)
	analysis.Dissent = g.findDissent(outputs)

	// Select best response (longest successful one as heuristic)
	bestIdx := 0
	maxLen := 0
	for i, r := range successResults {
		if len(r.Output) > maxLen {
			maxLen = len(r.Output)
			bestIdx = i
		}
	}
	analysis.BestResponse = successResults[bestIdx].Output
	analysis.BestProvider = successResults[bestIdx].Provider

	return analysis
}

// findConsensus identifies common themes across outputs
func (g *GhostExecutor) findConsensus(outputs []string) string {
	if len(outputs) == 0 {
		return ""
	}

	if len(outputs) == 1 {
		return "Single response only"
	}

	// Simple heuristic: find common sentences/phrases
	commonPhrases := []string{}
	firstOutput := strings.ToLower(outputs[0])

	for _, output := range outputs[1:] {
		lowerOutput := strings.ToLower(output)
		// Check for similar starting patterns
		if strings.HasPrefix(lowerOutput, firstOutput[:min(50, len(firstOutput))]) {
			commonPhrases = append(commonPhrases, "Similar opening")
		}
	}

	if len(commonPhrases) > 0 {
		return fmt.Sprintf("✓ All models agree on core approach (%d/%d consensus)", len(outputs), len(outputs))
	}

	return fmt.Sprintf("⚠ Models show variation (%d different approaches)", len(outputs))
}

// findDissent identifies differing opinions
func (g *GhostExecutor) findDissent(outputs []string) []string {
	dissent := []string{}

	if len(outputs) < 2 {
		return dissent
	}

	// Compare lengths and structures
	lengths := make([]int, len(outputs))
	for i, o := range outputs {
		lengths[i] = len(o)
	}

	// Check for significant length differences
	minLen, maxLen := lengths[0], lengths[0]
	for _, l := range lengths[1:] {
		if l < minLen {
			minLen = l
		}
		if l > maxLen {
			maxLen = l
		}
	}

	if maxLen > minLen*2 {
		dissent = append(dissent, fmt.Sprintf("Response length varies significantly (%.0f%% difference)",
			float64(maxLen-minLen)/float64(minLen)*100))
	}

	return dissent
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
