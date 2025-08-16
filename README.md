```
     ██████╗████████╗██╗  ██╗
    ██╔════╝╚══██╔══╝╚██╗██╔╝
    ██║        ██║    ╚███╔╝ 
    ██║        ██║    ██╔██╗ 
    ╚██████╗   ██║   ██╔╝ ██╗
     ╚═════╝   ╚═╝   ╚═╝  ╚═╝
```

# ctx - The Context Engine for AI

> **Slash AI costs by 95%: From $1,125 to just $10/month - the price of ctx Pro itself!**

[![Version](https://img.shields.io/badge/version-0.1.0-orange.svg)](https://github.com/slavakurilyak/ctx/releases)
[![Beta](https://img.shields.io/badge/status-beta-yellow.svg)](docs/VERSIONING.md)
[![Schema](https://img.shields.io/badge/schema-0.1-purple.svg)](docs/VERSIONING.md)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Token Efficient](https://img.shields.io/badge/Token-Efficient-green)](https://docs.anthropic.com/en/docs/build-with-claude/token-efficient-tool-use)
[![Cost Reduction](https://img.shields.io/badge/Cost%20Reduction-$1,125→$10/month-brightgreen)](docs/PRICING.md)
[![ctx Pro](https://img.shields.io/badge/ctx%20Pro-Coming%20Soon-orange.svg)](docs/PRICING.md)

**Beta Software**: ctx is currently in beta (v0.1.0). The API and schema may change before the 1.0 release. See [VERSIONING.md](docs/VERSIONING.md) for details.

Modern AI agents need to interact with external tools, but raw command output is token-expensive, hard to parse, and lacks critical metadata. This wastes context windows and drives up costs, especially with database-heavy workflows.

`ctx` solves this by wrapping any command in a structured JSON "context envelope." It provides precise token counts, execution metadata, and telemetry, enabling agents to work more efficiently and make smarter, cost-aware decisions. `ctx` is how you give your AI agents ground truth.

**Free and Open Source**: Core features including token counting, metadata enrichment, and safety controls are free forever. **[ctx Pro](docs/PRICING.md)** (coming soon) will add intelligent webhook integrations for command analysis, optimization, and security - starting at $10/month.

### Works With Your Favorite Tools

`ctx` has been tested with a wide range of standard CLI tools, including:

- **Databases**: `psql`, `mysql`, `sqlite3`, `redis-cli`, `mongosh`
- **Dev Tools**: `git`, `docker`, `curl`, `jq`, `kubectl`, `terraform`
- **System**: `ls`, `cat`, `grep`, `find`, `wc`, `df`, `uptime`
- **And any multi-stage pipeline you can create.**

## Key Features

### Free & Open Source
- **95% Cost Reduction**: Slash your LLM costs from $1,125 to $10/month - ctx Pro literally pays for itself 100x over!
- **Universal Tool Wrapper**: Works with any CLI tool out-of-the-box.
- **Structured JSON Output**: Enriches raw output with metadata in a clean, predictable format perfect for LLM consumption.
  ```json
  {
    "tokens": 42,
    "output": "...",
    "input": "...",
    "metadata": {
      "success": true,
      "exit_code": 0,
      "duration": 127,
      "bytes": 245,
      "failure_reason": "token_limit_exceeded"  // Machine-readable error codes
    },
    "telemetry": { "trace_id": "...", "span_id": "..." },
    "schema_version": "0.1"
  }
  ```
- **Precise Token Counting**: Uses `tiktoken` and `gemini` tokenizers for providers OpenAI, Anthropic, and Google.
- **Resource Limits & Controls**: Enforce safety limits on output size, line count, tokens, and pipeline complexity to prevent runaway commands and cost overruns.
- **Streaming Support**: Real-time monitoring and termination of long-running commands when limits are exceeded.
- **Intelligent Pipelines**: Chain multiple commands, with `ctx` enriching the output at each stage.
- **Privacy-Aware**: Telemetry and history can be disabled for sensitive operations with a simple `--private` flag.
- **Great DX**: Built with [Charm](https://charm.sh) for a beautiful and usable CLI experience.

### [ctx Pro](docs/PRICING.md) - Intelligent Command Analysis (Coming Soon)
- **Pre-Command Analysis**: AI-powered command analysis before execution with `ALLOW`, `BLOCK`, `MODIFY`, or `WARN` responses
- **Security Protection**: Automatically block dangerous commands like `rm -rf` or unauthorized operations
- **Command Optimization**: Receive suggestions for more efficient commands and best practices
- **Post-Command Insights**: Analyze command results for errors, performance patterns, and improvement opportunities
- **Team Management**: Centralized policies, audit logs, and shared security configurations
- **Custom Webhooks**: Integrate with your own analysis services and monitoring systems

**Pricing (Coming Soon)**: $10/month per individual user • $20/month per team seat

## The `ctx` Advantage: 95% Cost Reduction - Just $10/month Total

`ctx` enables a "measure-then-act" workflow that can reduce token consumption by over 95%, dropping your total AI costs to just the price of ctx Pro itself!

**Before `ctx` (Inefficient & Expensive):**
An agent runs a command and gets a huge, unstructured text blob.
```bash
# Agent runs a query and sends the entire raw output to the LLM
psql -c "SELECT * FROM users" | claude -p 'Summarize users'
# Result: ~25,000 tokens consumed = ~$1,125/month in token costs
```

**With `ctx` (Efficient & Cost-Effective):**
The agent first checks the token count, then refines its query to be more specific.
```bash
# 1. Measure the cost
ctx psql -c "SELECT * FROM users" | jq '.tokens'
# Result: 25000

# 2. Refine and execute
ctx psql -c "SELECT status, COUNT(*) FROM users GROUP BY status" | claude -p 'Analyze user stats'
# Result: ~100 tokens consumed (99.6% reduction) = ~$0/month in token costs
```

**That's right - your total AI costs become just the ctx Pro subscription itself!**

This simple pattern, applied across hundreds of daily operations, can reduce costs from **~$1,125/month** to **~$10/month** - just the price of ctx Pro!

## ctx Pro - Intelligent Command Analysis (Coming Soon)

**[ctx Pro](docs/PRICING.md)** will enhance your command-line experience with AI-powered analysis, security protection, and optimization suggestions.

### Pre-Command Analysis
ctx Pro will analyze every command before execution and can:
- **ALLOW**: Execute the command normally
- **BLOCK**: Prevent dangerous operations for security
- **MODIFY**: Suggest or apply optimizations automatically  
- **WARN**: Flag potentially risky commands but allow execution

```bash
# Dangerous command gets blocked
$ ctx rm -rf /important/data
Error: Command blocked: potentially dangerous rm -rf detected

# Inefficient command gets optimized
$ ctx ls
Note: Command modified: added -la flags for detailed listing
# Executes: ls -la

# Privileged command shows warning
$ ctx sudo systemctl restart nginx  
Warning: Command requires elevated privileges
# Command proceeds with warning
```

### Post-Command Intelligence
After commands execute, ctx Pro will analyze results and provide:
- **Error Analysis**: Detailed explanations of failures with remediation steps
- **Performance Insights**: Identify slow operations and optimization opportunities
- **Learning Recommendations**: Suggest better approaches and best practices
- **Pattern Recognition**: Detect recurring issues and workflow improvements

### Team Features
- **Centralized Policies**: Apply consistent security rules across your team
- **Audit Logs**: Track all command executions and policy decisions
- **Custom Webhooks**: Integrate with your monitoring and security systems
- **Shared Analytics**: Team-wide insights into command patterns and optimization opportunities

### Getting Started with ctx Pro (Coming Soon)

1. **Sign up** at [ctx.click](https://ctx.click) (14-day free trial)
2. **Get your API key** from the dashboard
3. **Authenticate** with ctx:
   ```bash
   ctx login
   # Enter your API key when prompted
   ```
4. **Check your account**:
   ```bash
   ctx account
   # Shows plan, billing, and webhook configuration
   ```

Pro features will activate automatically once authenticated when available. [View upcoming pricing →](docs/PRICING.md)

## Installation

### One-Liner (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/slavakurilyak/ctx/main/scripts/install-remote.sh | bash
```

### Configuration Setup
After installation, set up your environment:
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env to configure your API endpoint (for future ctx Pro features)
# Default endpoint is already configured: https://api.ctx.click
```

### With Go
```bash
# Install latest beta
go install github.com/slavakurilyak/ctx@latest

# Or install specific version
go install github.com/slavakurilyak/ctx@v0.1.0
```

### Binaries
Pre-built binaries for various platforms are available on the [**Releases page**](https://github.com/slavakurilyak/ctx/releases/latest).

After installation, verify with:
```bash
ctx version
# Output:
# Software Version: 0.1.0
# Schema Version:   0.1
```

**Note**: All core features work immediately after installation. [ctx Pro features](docs/PRICING.md) (coming soon) will require a subscription and authentication with `ctx login`.

## Usage

`ctx` is designed to be a simple prefix to any command.

```bash
# Basic command
ctx ls -la

# Database query with a specific token provider
ctx --token-model openai psql -c "SELECT id, email FROM users LIMIT 10"

# Pipeline: filter a large file before counting tokens
ctx cat large.json | jq .items[0:5]

# AI Agent workflow: generate a commit message from a git diff
ctx git diff --staged | claude -p 'Generate a conventional commit message from this diff.'

# ctx Pro (coming soon): Command will be analyzed and potentially optimized before execution
ctx psql -c "SELECT * FROM users"  # Will suggest LIMIT clause or specific columns
```

## Configuration

`ctx` supports multiple configuration sources with the following precedence (highest to lowest):
1. **CLI flags** - Command-line arguments
2. **Environment variables** - System environment
3. **Configuration file** - `~/.config/ctx/config.yaml`
4. **Defaults** - Built-in defaults

### Configuration Options

| Flag | Environment Variable | Description | Default |
|---|---|---|---|
| `--token-model` | `CTX_TOKEN_MODEL` | Set tokenizer provider (`anthropic`, `openai`, `gemini`) | `anthropic` |
| `--max-tokens` | `CTX_MAX_TOKENS` | Maximum tokens allowed in output (0 = unlimited) | `0` |
| `--max-output-bytes` | `CTX_MAX_OUTPUT_BYTES` | Maximum bytes allowed in output (0 = unlimited) | `0` |
| `--max-lines` | `CTX_MAX_LINES` | Maximum lines allowed in output (0 = unlimited) | `0` |
| `--max-pipeline-stages` | `CTX_MAX_PIPELINE_STAGES` | Maximum pipeline stages allowed (0 = unlimited) | `0` |
| `--private` | `CTX_PRIVATE` | Disable history and telemetry | `false` |
| `--no-history` | `CTX_NO_HISTORY` | Disable history recording | `false` |
| `--no-telemetry` | `CTX_NO_TELEMETRY` | Disable OpenTelemetry tracing | `false` |
| `--timeout` | `CTX_TIMEOUT` | Set command timeout (e.g., `30s`, `1m`) | `2m` |
| `--stream` | - | Stream output line by line for long-running commands | `false` |
| - | `CTX_API_ENDPOINT` | API endpoint for ctx Pro features (set in .env) | Coming soon |

### Configuration File

Create `~/.config/ctx/config.yaml` to set persistent defaults:

```yaml
token_model: openai
default_timeout: 5m
limits:
  max_tokens: 10000
  max_output_bytes: 1048576  # 1MB
  max_lines: 1000
  max_pipeline_stages: 5

# ctx Pro configuration (coming soon - requires subscription)
auth:
  # api_endpoint can be set here or via CTX_API_ENDPOINT in .env
webhooks:
  pre_tool_use:
    enabled: true
    timeout: 5s
  post_tool_use:
    enabled: true
    timeout: 10s
```

### View Active Configuration

Use the `config view` command to see your current configuration and where each setting comes from:

```bash
ctx config view

# Output:
# Current ctx Configuration:
# ==========================
# Token Model:        openai (source: config file)
# Default Timeout:    5m0s (source: config file)
# ...
# Limits:
#   Max Tokens:       10000 (source: config file)
#   Max Output Bytes: 1048576 (source: config file)
```

Command history is stored in the `.ctx/` directory in your project root.

## ctx Pro Authentication (WIP - Not Fully Implemented)

**⚠️ Note: Authentication features are currently work in progress and not fully implemented or tested. Coming soon!**

ctx Pro features will require authentication with your subscription. These commands will be available to manage your account:

### Login (WIP)
Authenticate with your ctx Pro API key:
```bash
ctx login  # (WIP - Not fully implemented)
# You'll be prompted to enter your API key securely
# API keys will be stored in your system's secure keychain
```

### Check Account Status (WIP)
View your subscription details and configuration:
```bash
ctx account  # (WIP - Not fully implemented)

# Example output:
# ctx Pro Account Status
# ======================
# Email:       developer@example.com
# Tier:        pro
# Status:      Active
# Valid Until: 2025-02-13
# 
# Billing Information:
#   Plan:         individual
#   Price:        $10.00 USD/monthly
#   Next Billing: 2025-02-13
# 
# Webhook Configuration:
#   Pre-Tool-Use:  Enabled (timeout: 5s)
#   Post-Tool-Use: Enabled (timeout: 10s)
```

### Logout (WIP)
Clear your stored credentials:
```bash
ctx logout  # (WIP - Not fully implemented)
# Will remove API key from keychain and disable Pro features
```

### Available Commands
| Command | Description | Status |
|---------|-------------|--------|
| `ctx login` | Authenticate with your Pro API key | WIP - Not fully implemented |
| `ctx logout` | Remove stored credentials and disable Pro features | WIP - Not fully implemented |
| `ctx account` | View subscription status, billing, and configuration | WIP - Not fully implemented |

Pro features (webhooks, intelligent analysis) will activate automatically when authenticated once fully implemented.

## Resource Limits & Safety Controls

`ctx` provides comprehensive resource limits to prevent runaway commands and control costs. When a limit is exceeded, the command is terminated immediately and a clear error is returned.

### Limit Types

- **Output Size Limit** (`--max-output-bytes`): Prevents commands from producing excessive output
- **Line Count Limit** (`--max-lines`): Useful for sampling the beginning of large outputs
- **Token Limit** (`--max-tokens`): Controls LLM token consumption directly
- **Pipeline Complexity** (`--max-pipeline-stages`): Prevents overly complex command chains

### Examples

```bash
# Sample first 100 lines of a large file
ctx --max-lines 100 cat huge_log_file.txt

# Prevent token explosion in database queries
ctx --max-tokens 5000 psql -c "SELECT * FROM events"

# Limit output size to 1MB
ctx --max-output-bytes 1048576 docker logs container_id

# Restrict pipeline complexity
ctx --max-pipeline-stages 3 -- cat file | grep error | sort | uniq | wc -l
# Error: pipeline stage limit exceeded: 5 stages found, limit is 3

# Streaming mode with limits (terminates on breach)
ctx --stream --max-lines 50 tail -f /var/log/app.log
```

### Limit Exceeded Behavior

When a limit is exceeded:
1. The command is terminated immediately
2. Partial output up to the limit is preserved
3. `metadata.success` is set to `false`
4. `metadata.failure_reason` contains a machine-readable error code:
   - `"line_limit_exceeded"`
   - `"output_limit_exceeded"`
   - `"token_limit_exceeded"`
   - `"pipeline_limit_exceeded"`

This allows AI agents to handle limit breaches intelligently:

```bash
# Agent can check if limit was hit and refine the query
result=$(ctx --max-tokens 1000 psql -c "SELECT * FROM users")
if [ $(echo "$result" | jq -r '.metadata.failure_reason') = "token_limit_exceeded" ]; then
  # Refine query to be more specific
  ctx psql -c "SELECT id, email FROM users LIMIT 100"
fi
```

## Advanced Workflows

### AI Agent Integration

To get the most out of `ctx`, you need to instruct your AI agent on how to use it. We provide templates that teach agents to wrap commands in `ctx` and use the resulting metadata to perform self-optimization.

Create an instruction file (e.g., `CLAUDE.md`, `AGENT.md`) in your project root.

- **[View `CLAUDE.md` Template](https://github.com/slavakurilyak/ctx/blob/main/CLAUDE.md)** - Optimized for Claude AI
- **[View Generic `AGENT.md` Template](https://github.com/slavakurilyak/ctx/blob/main/AGENT.md)** - Works with any AI agent

You can download the templates directly:
```bash
# For Claude AI
curl -L https://raw.githubusercontent.com/slavakurilyak/ctx/main/CLAUDE.md -o CLAUDE.md

# For generic AI agents
curl -L https://raw.githubusercontent.com/slavakurilyak/ctx/main/AGENT.md -o AGENT.md
```

#### Set Up AI Coding Assistants

`ctx` can automatically set up popular AI coding assistants and IDEs with ctx documentation:

```bash
# Set up Cursor IDE (.cursor/rules/ctx.mdc)
ctx setup cursor

# Set up Windsurf IDE (.windsurf/rules/ctx.md)
ctx setup windsurf

# Set up JetBrains AI Assistant (.aiassistant/rules/ctx.md)
ctx setup jetbrains

# Set up Aider (.aider.conf.yml)
ctx setup aider

# Set up Zed editor (.rules)
ctx setup zed

# Set up GitHub Copilot (.github/copilot-instructions.md)
ctx setup github-copilot

# Set up all coding assistants at once (CLAUDE.md + all IDE rules)
ctx setup all

# Force overwrite existing files
ctx setup all --force
```

These commands create the appropriate directory structure and configuration files that enable your IDE's AI assistant to automatically wrap commands with `ctx`. The rules teach the AI about ctx's usage patterns, token optimization workflows, and best practices.

**Configuration:**
- **Cursor**: Rules in `.cursor/rules/ctx.mdc` are automatically applied with `alwaysApply: true`
- **Windsurf**: Rules in `.windsurf/rules/ctx.md` can be set to "Always On" in Windsurf settings
- **JetBrains**: Rules in `.aiassistant/rules/ctx.md` can be set to "Always" type in AI Assistant settings
- **Aider**: Configuration in `.aider.conf.yml` embeds ctx documentation as comments for reference
- **Zed**: Rules in `.rules` file are automatically included in all Agent Panel interactions
- **GitHub Copilot**: Instructions in `.github/copilot-instructions.md` are automatically included for all Copilot interactions
- **Claude Code**: Uses `CLAUDE.md` in the project root automatically

### ctx Pro Workflows (Coming Soon)

**Note: These are preview examples of upcoming ctx Pro features.**

ctx Pro will enhance AI agent interactions with intelligent analysis and optimization:

#### Intelligent Database Queries
```bash
# Agent starts with exploratory query
ctx psql -c "SELECT * FROM users"
# ctx Pro (coming soon): "Warning: Query may return large dataset. Consider adding LIMIT clause."

# Agent refines based on feedback
ctx psql -c "SELECT id, email, created_at FROM users LIMIT 100"
# ctx Pro (coming soon): "ALLOW - Optimized query for efficiency"
```

#### Security-First Development
```bash
# Dangerous operations get blocked automatically
ctx rm -rf node_modules/.git
# ctx Pro (coming soon): "BLOCK - Potentially destructive file operation detected"

# Safe alternatives are suggested
ctx find node_modules -name ".git" -type d -exec rm -rf {} +
# ctx Pro (coming soon): "ALLOW - Safe file operation with proper targeting"
```

#### Team Collaboration
```bash
# Team policies applied consistently
ctx docker run --privileged alpine
# ctx Pro (coming soon): "WARN - Privileged container detected. Team policy requires approval."

# Compliance monitoring
ctx aws s3 rm s3://production-bucket --recursive
# ctx Pro (coming soon): "BLOCK - Production resource modification requires MFA approval"
```

#### Learning and Optimization
```bash
# After command execution, get insights
ctx git log --oneline -n 50
# ctx Pro Post-Analysis (coming soon):
# Pattern Recognition: Frequent small commits detected
# Suggestion: Consider squashing commits before pushing
# Performance: Command completed in 23ms (normal range)
```

## Versioning

ctx uses semantic versioning. Currently in **beta (v0.1.0)**, the API and output schema may change between minor versions. We plan to release v1.0.0 when the API is stable and battle-tested.

- **Current Version**: 0.1.0
- **Schema Version**: 0.1
- **Minimum Go Version**: 1.21

See [VERSIONING.md](docs/VERSIONING.md) for detailed versioning strategy and compatibility information.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

```bash
# Get set up for development
git clone https://github.com/slavakurilyak/ctx.git
cd ctx
go mod tidy
go test ./...
```

### Roadmap
Our roadmap is driven by the community. Check out the [**open issues**](https://github.com/slavakurilyak/ctx/issues) to see what's planned or to suggest new features.

## Links

### ctx
- **Repository**: [https://github.com/slavakurilyak/ctx](https://github.com/slavakurilyak/ctx)
- **Issues**: [https://github.com/slavakurilyak/ctx/issues](https://github.com/slavakurilyak/ctx/issues)
- **Releases**: [https://github.com/slavakurilyak/ctx/releases](https://github.com/slavakurilyak/ctx/releases)
- **Documentation**: [README.md](README.md) | [CLAUDE.md](CLAUDE.md) | [AGENT.md](AGENT.md)
- **Pricing**: [ctx Pro](docs/PRICING.md)

### Coding Agents

| Agent | Links | Compatible with ctx |
|-------|-------|---------------------|
| **Claude Code** | [Official site](https://www.anthropic.com/claude-code) \| [Docs](https://docs.anthropic.com/en/docs/claude-code/overview) | ✓ (via `ctx setup claude`) |
| **Cursor** | [Official site](https://cursor.sh/) \| [Docs](https://cursor.sh/help) | ✓ (via `ctx setup cursor`) |
| **Aider** | [Official site](https://aider.chat/) \| [GitHub](https://github.com/Aider-AI/aider) \| [Docs](https://aider.chat/docs/) | ✓ (via `ctx setup aider`) |
| **Visual Studio Code** | [Official site](https://code.visualstudio.com/) \| [GitHub](https://github.com/microsoft/vscode) \| [Docs](https://code.visualstudio.com/docs) | ✗ |
| **Visual Studio 2022** | [Official site](https://visualstudio.microsoft.com/) \| [Docs](https://learn.microsoft.com/en-us/visualstudio/windows/) | ✗ |
| **JetBrains AI Assistant** | [Official site](https://www.jetbrains.com/ai-assistant/) \| [Docs](https://www.jetbrains.com/help/ai-assistant/) | ✓ (via `ctx setup jetbrains`) |
| **Zed** | [Official site](https://zed.dev/) \| [GitHub](https://github.com/zed-industries/zed) \| [Docs](https://zed.dev/docs) | ✓ (via `ctx setup zed`) |
| **Opencode** | [Official site](https://opencode.ai/) \| [Docs](https://opencode.ai/docs) | ✗ |
| **Augment Code** | [Official site](https://www.augmentcode.com/) \| [Docs](https://www.augmentcode.com/docs/getting-started) | ✗ |
| **GitHub Copilot** | [Official site](https://github.com/features/copilot) \| [Docs](https://docs.github.com/en/copilot) | ✓ (via `ctx setup github-copilot`) |
| **Amazon Q Developer** | [Official site](https://aws.amazon.com/q/developer/) \| [Docs](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/what-is.html) | ✗ |
| **Windsurf** | [Official site](https://windsurf.ai/) \| [GitHub](https://github.com/Windsurf-AI/windsurf) \| [Docs](https://docs.windsurf.com/) | ✓ (via `ctx setup windsurf`) |
| **Trae** | [Official site](https://trae.ai/) \| [Docs](https://docs.trae.ai/) | ✗ |
| **Roo Code** | [Official site](https://roocode.com/) \| [Docs](https://docs.roocode.com/) | ✗ |
| **Zencoder** | [Official site](https://www.zencoder.dev/) \| [GitHub](https://github.com/zencoder-platform/zencoder) | ✗ |
| **Qodo Gen** | [Official site](https://qodo.ai/) \| [Docs](https://docs.qodo.ai/) | ✗ |
| **Kiro** | [Official site](https://kiro.dev/) \| [Docs](https://kiro.dev/docs/intro) | ✗ |
| **Gemini CLI** | [Official site](https://ai.google.dev/docs/gemini_cli) \| [GitHub](https://github.com/google-gemini/gemini-cli) \| [Docs](https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/gcloud.md) | ✗ |
| **Warp Terminal** | [Official site](https://www.warp.dev/) \| [Docs](https://docs.warp.dev/) | ✗ |
| **Cline** | [Official site](https://cline.bot/) \| [GitHub](https://github.com/re-search/cline) | ✗ |
| **Crush** | [Official site](https://charm.land/) \| [GitHub](https://github.com/charmbracelet/crush) | ✗ |
| **Rovo Dev CLI** | [Official site](https://rovodev.com/) \| [GitHub](https://github.com/rovotech/rovodev) | ✗ |
| **LM Studio** | [Official site](https://lmstudio.ai/) \| [Docs](https://lmstudio.ai/docs/welcome) | ✗ |
| **BoltAI** | [Official site](https://boltai.com/) \| [Docs](https://docs.boltai.com/) | ✗ |
| **Perplexity Desktop** | [Official site](https://perplexity.ai/downloads) \| [Docs](https://docs.perplexity.ai/docs) | ✗ |
| **Claude Desktop** | [Official site](https://www.anthropic.com/claude) \| [Docs](https://docs.anthropic.com/en/docs/intro-to-claude) | ✗ |

## Acknowledgments

- [Charm](https://charm.sh) for the amazing Fang CLI framework
- [OpenAI](https://openai.com) and [Google](https://ai.google/) for their tokenizers
- [Anthropic](https://anthropic.com) for pioneering token-efficient tool use
- The creators of countless open-source CLI tools

## License

This project is licensed under the [MIT License](LICENSE).
