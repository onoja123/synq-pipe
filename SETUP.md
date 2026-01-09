# 🚀 Synqly Multiverse CLI - Complete Setup Guide

This guide will take you from zero to running your first multi-model AI workflow in under 5 minutes.

## Prerequisites

- **Go 1.21+** installed ([download here](https://go.dev/dl/))
- **Synqly API Key** ([get one here](https://synqly.com))
- Terminal/Command line access
- Internet connection

## Step 1: Clone the Repository

```bash
git clone https://github.com/onoja123/synqly-multiverse.git
cd synqly-multiverse
```

## Step 2: Set Up Your API Key

### Option A: Environment Variable (Recommended)

```bash
# Linux/Mac
export SYNQLY_API_KEY="your_synqly_api_key_here"

# Windows (PowerShell)
$env:SYNQLY_API_KEY="your_synqly_api_key_here"

# Windows (CMD)
set SYNQLY_API_KEY=your_synqly_api_key_here
```

### Option B: Configuration File

```bash
# Create config file
cp .env.example .env

# Edit with your API key
nano .env  # or use your preferred editor
```

Add your key:
```
SYNQLY_API_KEY=sk_your_actual_key_here
```

## Step 3: Install Dependencies

```bash
# Download all Go dependencies
make deps

# This will:
# - Download Synqly Go SDK
# - Install Bubble Tea (TUI)
# - Install Cobra (CLI framework)
# - Install all other dependencies
```

## Step 4: Build the CLI

```bash
# Build the binary
make build

# This creates: bin/synq-pipe
```

Verify the build:
```bash
./bin/synq-pipe --help
```

You should see:
```
Synqly Multiverse CLI - Multi-model workflow orchestration

A powerful CLI for orchestrating multi-model AI workflows with consensus analysis.

Usage:
  synq-pipe [command]

Available Commands:
  run         Run a workflow from YAML file
  validate    Validate a workflow YAML file
  list        List available providers and models
  help        Help about any command

Flags:
  -h, --help      help for synq-pipe
  -v, --verbose   verbose output
```

## Step 5: Test with Example Workflow

### Run a Basic Workflow

```bash
./bin/synq-pipe run examples/basic-workflow.yaml
```

Expected output:
```
🚀 Starting workflow: Basic AI Pipeline
   A simple workflow demonstrating multi-step AI processing

📍 Step 1/3: Extract Keywords
   🤖 Normal mode: single provider...
✓ Completed in 2.3s

📍 Step 2/3: Generate Summary
   🤖 Normal mode: single provider...
✓ Completed in 3.1s

📍 Step 3/3: Translate to Spanish
   🤖 Normal mode: single provider...
✓ Completed in 2.8s

══════════════════════════════════════════════════
📊 Workflow Summary: Basic AI Pipeline
══════════════════════════════════════════════════
Status: ✅ SUCCESS
Total Time: 8.2s
Steps Completed: 3/3
══════════════════════════════════════════════════
```

### Run with Interactive TUI

```bash
./bin/synq-pipe run examples/basic-workflow.yaml --tui
```

This launches the beautiful terminal UI!

### Run Ghost Mode Demo

```bash
./bin/synq-pipe run examples/ghost-review.yaml --tui
```

## Step 6: Create Your First Workflow

Create a file `my-workflow.yaml`:

```yaml
name: My First Workflow
description: Testing Synqly Multiverse CLI
version: 1.0

steps:
  - name: Say Hello
    model: gpt-4
    provider: openai
    prompt: "Say hello and introduce yourself as an AI assistant"
    temperature: 0.7

  - name: Tell a Joke
    model: claude-sonnet-4
    provider: anthropic
    prompt: "Tell me a joke about programming"
    temperature: 0.8
```

Run it:
```bash
./bin/synq-pipe run my-workflow.yaml --tui
```

## Step 7: Try Advanced Features

### Ghost Mode (Consensus Analysis)

Create `consensus-test.yaml`:

```yaml
name: Ghost Mode Test
description: Testing multi-provider consensus

steps:
  - name: Analyze Topic
    mode: ghost
    providers:
      - openai
      - anthropic
      - google
    prompt: |
      What are the top 3 benefits of using AI in software development?
      Be specific and concise.
    temperature: 0.3
```

Run it:
```bash
./bin/synq-pipe run consensus-test.yaml --tui
```

Watch as it queries all three providers and shows consensus!

### Step Chaining

Create `chain-test.yaml`:

```yaml
name: Chained Workflow
description: Demonstrating step output chaining

steps:
  - name: Generate Idea
    model: gpt-4
    prompt: "Suggest a unique app idea in one sentence"

  - name: Expand Idea
    model: claude-sonnet-4
    input: "{{steps.Generate Idea.output}}"
    prompt: "Expand this app idea into a detailed description"

  - name: Technical Stack
    model: gemini-1.5-pro
    input: "{{steps.Expand Idea.output}}"
    prompt: "Suggest the best technical stack for building this app"
```

Run it:
```bash
./bin/synq-pipe run chain-test.yaml --tui
```

## Troubleshooting

### "API key is required" error

Make sure your API key is set:
```bash
echo $SYNQLY_API_KEY  # Should print your key
```

If not set:
```bash
export SYNQLY_API_KEY="your_key_here"
```

### "Failed to parse workflow" error

Validate your YAML:
```bash
./bin/synq-pipe validate my-workflow.yaml
```

Common YAML issues:
- Indentation must use spaces (not tabs)
- Each step needs a `name` and either `model` or `mode: ghost`
- Quotes are needed for multi-line prompts using `|`

### Build fails

Ensure Go 1.21+ is installed:
```bash
go version  # Should show 1.21 or higher
```

Clear module cache and retry:
```bash
go clean -modcache
make deps
make build
```

### API request fails

1. Check your internet connection
2. Verify your Synqly API key is valid
3. Check Synqly status at [status.synqly.com](https://status.synqly.com)

Enable verbose mode for more details:
```bash
./bin/synq-pipe run workflow.yaml --verbose
```

## Next Steps

### Install Globally (Optional)

```bash
# Install to system PATH
make install

# Now you can run from anywhere:
synq-pipe run workflow.yaml
```

### Explore Examples

Check out the `examples/` directory:
- `basic-workflow.yaml` - Simple linear workflow
- `ghost-review.yaml` - Code review with consensus
- `full-pipeline.yaml` - Complex multi-step demo

### Learn More

- 📖 Read the [README.md](README.md) for full documentation
- 🔮 Learn about [Ghost Mode](docs/ghost-mode.md)
- 📝 See [workflow syntax reference](docs/workflow-syntax.md)
- 💡 Browse [workflow recipes](docs/recipes.md)

## Quick Reference

```bash
# Build
make build

# Run workflow
./bin/synq-pipe run workflow.yaml

# With TUI
./bin/synq-pipe run workflow.yaml --tui

# Force ghost mode
./bin/synq-pipe run workflow.yaml --ghost-enabled

# Validate syntax
./bin/synq-pipe validate workflow.yaml

# List providers
./bin/synq-pipe list

# Show help
./bin/synq-pipe --help
```

## Support

Need help?
- 🐛 Issues: [GitHub Issues](https://github.com/onoja123/synqly-multiverse/issues)

---

🎉 **Congratulations!** You're now ready to build powerful multi-model AI workflows with Synqly Multiverse CLI!