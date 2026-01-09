# 🚀 Synqly Multiverse CLI

**Multi-Model Workflow Engine with Ghost Mode**

A powerful CLI tool that orchestrates AI workflows across multiple providers (OpenAI, Anthropic, Google) using the Synqly unified API. Features include parallel model execution, consensus analysis, and an interactive TUI.

![Demo](https://img.shields.io/badge/status-production-green)
![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

## ✨ Features

- **Multi-Model Orchestration** - Chain AI operations across different providers
- **Ghost Mode** - Query multiple models in parallel and analyze consensus
- **Interactive TUI** - Beautiful terminal UI with real-time updates
- **Fast Execution** - Parallel processing for Ghost mode steps
- **YAML Workflows** - Simple, declarative workflow definitions
- **Provider Agnostic** - Works with OpenAI, Anthropic, Google, and more
- **Zero Config** - No API keys needed (uses Synqly)

## 🎬 Demo

```bash
# Run with interactive TUI
synq-pipe run workflow.yaml --tui

# Run with simple output
synq-pipe run workflow.yaml
```

## 📦 Installation

### Option 1: Install from source

```bash
# Clone repository
git clone https://github.com/onoja123/synq-pipe
cd synq-pipe

# Install dependencies
make deps

# Build and install
make install
```

### Option 2: Download binary

Download pre-built binaries from [Releases](https://github.com/onoja123/synq-pipe/releases)

### Option 3: Go install

```bash
go install github.com/onoja123/synq-pipe/cmd/synq-pipe@latest
```

## 🚀 Quick Start

### 1. Set up API key

```bash
export SYNQLY_API_KEY="your-api-key"
```

### 2. Create a workflow

```yaml
# simple.yaml
name: Translation Pipeline
description: Translate text across models
input: "Hello, world!"

steps:
  - name: Translate to French
    provider: openai
    model: gpt-4

  - name: Translate to Spanish
    provider: anthropic
    model: claude-sonnet-4
```

### 3. Run workflow

```bash
synq-pipe run simple.yaml --tui
```

## 📖 Workflow YAML Format

### Basic Step

```yaml
steps:
  - name: Generate Text
    provider: openai      # Provider: openai, anthropic, google
    model: gpt-4         # Model name
    prompt: "You are a helpful assistant"  # Optional system prompt
    temperature: 0.7     # Optional (0.0 - 1.0)
    max_tokens: 500      # Optional
```

### Ghost Mode Step

```yaml
steps:
  - name: Code Review
    mode: ghost          # Enable Ghost mode
    providers:           # Multiple providers
      - openai
      - anthropic
      - google
    prompt: "Review this code for issues"
```

### Complete Example

```yaml
name: Content Creation Pipeline
description: Multi-stage content generation
input: "Write about AI benefits"

steps:
  - name: Generate Outline
    provider: openai
    model: gpt-4
    temperature: 0.7

  - name: Fact Check (Ghost Mode)
    mode: ghost
    providers: [openai, anthropic, google]
    prompt: "Verify factual accuracy"

  - name: Write Article
    provider: anthropic
    model: claude-sonnet-4
```

## 👻 Ghost Mode

Ghost Mode queries multiple AI models in parallel and provides:

- **Consensus Detection** - Identifies agreement across models
- **Dissent Analysis** - Highlights differing opinions
- **Best Response Selection** - Automatically picks optimal output
- **Comparison View** - Side-by-side model outputs

Use Ghost Mode for:
- Code reviews
- Fact-checking
- Quality assurance
- Security audits
- Critical decisions

## 🎨 TUI Interface

The interactive TUI shows:

- **Left Panel**: Step list with progress
- **Right Panel**: Live model outputs
- **Ghost Mode**: Consensus analysis and provider comparison
- **Progress Bar**: Real-time completion status

### Controls

- `q` or `Ctrl+C` - Quit
- Auto-scrolls to show current step

## 📚 Examples

See `examples/` directory for complete workflows:

- `simple.yaml` - Basic translation pipeline
- `ghost.yaml` - Code review with Ghost mode
- `advanced.yaml` - Multi-stage content creation
- `creative.yaml` - Story writing with consensus
- `qa_testing.yaml` - Test generation and validation

Run examples:

```bash
make examples        # Run all examples
make run            # Run simple example
make run-ghost      # Run Ghost mode example
```

## 🛠️ Development

### Build

```bash
make build          # Development build
make build-prod     # Production build (optimized)
make build-all      # Build for all platforms
```

### Test

```bash
make test           # Run tests
make test-coverage  # Generate coverage report
make lint           # Lint code
```

### Development Mode

```bash
make dev            # Auto-reload on changes (requires air)
```

## 📁 Project Structure

```
synq-pipe/
├── cmd/synq-pipe/      # CLI entry point
├── pkg/
│   ├── pipeline/       # Pipeline execution engine
│   ├── parser/         # YAML parser
│   └── tui/           # Bubble Tea UI
├── examples/           # Example workflows
├── Makefile           # Build automation
└── README.md
```

## 🔧 Configuration

### Environment Variables

```bash
SYNQLY_API_KEY=your-key    # Required: Synqly API key
```

### Supported Providers

| Provider   | Models                    |
|-----------|---------------------------|
| OpenAI    | gpt-4, gpt-3.5-turbo     |
| Anthropic | claude-sonnet-4          |
| Google    | gemini-1.5-pro           |

## 🎯 Use Cases

### Code Review Pipeline

```yaml
name: Multi-Model Code Review
input: |
  def calculate_sum(numbers):
      return sum(numbers)

steps:
  - name: Security Analysis (Ghost)
    mode: ghost
    providers: [openai, anthropic, google]
    prompt: "Check for security vulnerabilities"
  
  - name: Performance Review (Ghost)
    mode: ghost
    providers: [openai, anthropic]
    prompt: "Analyze performance and optimization"
```

### Content Validation

```yaml
name: Fact-Check Pipeline
input: "Article content to verify"

steps:
  - name: Initial Review
    provider: openai
    model: gpt-4
  
  - name: Fact Verification (Ghost)
    mode: ghost
    providers: [openai, anthropic, google]
    prompt: "Verify all factual claims"
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

MIT License - see [LICENSE](LICENSE) file

## 🔗 Links

- [Synqly Website](https://synqly.xyz)
- [Synqly API Docs](https://synqly.xyz/docs)
- [Go SDK](https://github.com/onoja123/synqly-go)

## 💡 Tips

### Optimizing Ghost Mode

```yaml
# Use Ghost mode strategically for:
# 1. Critical decisions
# 2. Quality validation
# 3. Security reviews
# 4. Fact-checking

# Not recommended for:
# - Simple transformations
# - High-frequency operations
# - When speed is critical
```

### Performance

- Ghost mode runs providers in parallel
- Average overhead: 100-200ms for orchestration
- Network latency depends on provider response times

### Best Practices

1. **Chain thoughtfully** - Each step builds on previous output
2. **Use prompts** - Guide models with specific instructions
3. **Ghost strategically** - Use for validation, not every step
4. **Test workflows** - Validate YAML before production use

## 🐛 Troubleshooting

### Common Issues

**API Key Error**
```bash
export SYNQLY_API_KEY="your-key"
```

**Build Fails**
```bash
make clean
make deps
make build
```

**TUI Not Rendering**
- Ensure terminal supports ANSI colors
- Try without `--tui` flag

## 📊 Roadmap

- [ ] Streaming responses
- [ ] Workflow templates
- [ ] Result caching
- [ ] Custom provider integrations
- [ ] Web dashboard
- [ ] Workflow marketplace

## ⭐ Show Your Support

If you find this project useful, please:
- ⭐ Star the repository
- 🐛 Report bugs
- 💡 Suggest features
- 🔀 Submit pull requests

---

**Built with** ❤️ **using** [Synqly](https://synqly.xyz)  by [Onoja](https://iamthecode.xyz)