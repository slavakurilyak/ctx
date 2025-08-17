```
     ██████╗████████╗██╗  ██╗
    ██╔════╝╚══██╔══╝╚██╗██╔╝
    ██║        ██║    ╚███╔╝ 
    ██║        ██║    ██╔██╗ 
    ╚██████╗   ██║   ██╔╝ ██╗
     ╚═════╝   ╚═╝   ╚═╝  ╚═╝
```

# Context (ctx) - The Context Engine for AI Agents

> **Slash AI costs by 95%: From $1,125 to just $10/month - the price of the upcoming ctx Pro subscription!**

[![Version](https://img.shields.io/badge/version-0.1.0-orange.svg)](https://github.com/slavakurilyak/ctx/releases)
[![Beta](https://img.shields.io/badge/status-beta-yellow.svg)](docs/VERSIONING.md)
[![Schema](https://img.shields.io/badge/schema-0.1-purple.svg)](docs/VERSIONING.md)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Beta Software**: `ctx` is in beta (v0.1.0). The API and schema may change. See [VERSIONING.md](docs/VERSIONING.md) for details.

Modern AI agents struggle with external tools because raw command output is token-expensive, unstructured, and lacks metadata. This wastes context windows and drives up costs.

`ctx` solves this by wrapping any command in a structured JSON "context envelope." It adds precise token counts, execution metadata, and telemetry, enabling agents to make smarter, cost-aware decisions.

## Supported Coding Agents

`ctx` integrates with leading AI coding assistants and agentic IDEs through simple setup commands. We're constantly expanding support to include more tools based on community needs.

**Want to see your favorite AI tool supported?** [Request a new integration](https://github.com/slavakurilyak/ctx/issues/new?title=Integration%20Request:%20[Tool%20Name]&body=Please%20add%20support%20for%20[Tool%20Name]%0A%0AOfficial%20site:%20%0ADocs:%20%0AHow%20it%20handles%20custom%20instructions:%20) by creating a GitHub issue.

| Agent | Links | Compatible with ctx |
|-------|-------|---------------------|
| **Claude Code** | [Official site](https://www.anthropic.com/claude-code) \| [Docs](https://docs.anthropic.com/en/docs/claude-code/overview) | ✓ (via `ctx setup claude`) |
| **Gemini CLI** | [Official site](https://ai.google.dev/docs/gemini_cli) \| [GitHub](https://github.com/google-gemini/gemini-cli) \| [Docs](https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/gcloud.md) | ✓ (via `ctx setup gemini`) |
| **Kilo Code** | [Official site](https://kilocode.ai/) \| [Docs](https://kilocode.ai/docs) | ✓ (via `ctx setup kilo-code`) |
| **Roo Code** | [Official site](https://roocode.com/) \| [Docs](https://docs.roocode.com/) | ✓ (via `ctx setup roo-code`) |
| **Cursor** | [Official site](https://cursor.sh/) \| [Docs](https://cursor.sh/help) | ✓ (via `ctx setup cursor`) |
| **Aider** | [Official site](https://aider.chat/) \| [GitHub](https://github.com/Aider-AI/aider) \| [Docs](https://aider.chat/docs/) | ✓ (via `ctx setup aider`) |
| **JetBrains AI Assistant** | [Official site](https://www.jetbrains.com/ai-assistant/) \| [Docs](https://www.jetbrains.com/help/ai-assistant/) | ✓ (via `ctx setup jetbrains`) |
| **Zed** | [Official site](https://zed.dev/) \| [GitHub](https://github.com/zed-industries/zed) \| [Docs](https://zed.dev/docs) | ✓ (via `ctx setup zed`) |
| **GitHub Copilot** | [Official site](https://github.com/features/copilot) \| [Docs](https://docs.github.com/en/copilot) | ✓ (via `ctx setup github-copilot`) |
| **Windsurf** | [Official site](https://windsurf.ai/) \| [GitHub](https://github.com/Windsurf-AI/windsurf) \| [Docs](https://docs.windsurf.com/) | ✓ (via `ctx setup windsurf`) |
| **Cline** | [Official site](https://cline.bot/) \| [GitHub](https://github.com/re-search/cline) | ✓ (via `ctx setup cline`) |
| **Visual Studio Code** | [Official site](https://code.visualstudio.com/) \| [GitHub](https://github.com/microsoft/vscode) \| [Docs](https://code.visualstudio.com/docs) | ✗ |
| **Visual Studio 2022** | [Official site](https://visualstudio.microsoft.com/) \| [Docs](https://learn.microsoft.com/en-us/visualstudio/windows/) | ✗ |
| **Opencode** | [Official site](https://opencode.ai/) \| [Docs](https://opencode.ai/docs) | ✗ |
| **Augment Code** | [Official site](https://www.augmentcode.com/) \| [Docs](https://www.augmentcode.com/docs/getting-started) | ✗ |
| **Amazon Q Developer** | [Official site](https://aws.amazon.com/q/developer/) \| [Docs](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/what-is.html) | ✗ |
| **Trae** | [Official site](https://trae.ai/) \| [Docs](https://docs.trae.ai/) | ✗ |
| **Zencoder** | [Official site](https://www.zencoder.dev/) \| [GitHub](https://github.com/zencoder-platform/zencoder) | ✗ |
| **Qodo Gen** | [Official site](https://qodo.ai/) \| [Docs](https://docs.qodo.ai/) | ✗ |
| **Warp Terminal** | [Official site](https://www.warp.dev/) \| [Docs](https://docs.warp.dev/) | ✗ |
| **Crush** | [Official site](https://charm.land/) \| [GitHub](https://github.com/charmbracelet/crush) | ✗ |
| **Rovo Dev CLI** | [Official site](https://rovodev.com/) \| [GitHub](https://github.com/rovotech/rovodev) | ✗ |
| **LM Studio** | [Official site](https://lmstudio.ai/) \| [Docs](https://lmstudio.ai/docs/welcome) | ✗ |
| **BoltAI** | [Official site](https://boltai.com/) \| [Docs](https://docs.boltai.com/) | ✗ |
| **Perplexity Desktop** | [Official site](https://perplexity.ai/downloads) \| [Docs](https://docs.perplexity.ai/docs) | ✗ |
| **Claude Desktop** | [Official site](https://www.anthropic.com/claude) \| [Docs](https://docs.anthropic.com/en/docs/intro-to-claude) | ✗ |

## The `ctx` Advantage: Drastic Cost Reduction

`ctx` enables a "measure-then-act" workflow that can reduce token consumption by over 95%.

**Before `ctx` (Expensive):** An agent gets a huge, raw text blob.
```bash
# Agent sends the entire raw output to the LLM
psql -c "SELECT * FROM users" | llm -p 'Summarize users'
# Result: ~25,000 tokens consumed (~$1,125/month)
```

**With `ctx` (Efficient):** The agent first checks the cost, then refines its query.
```bash
# 1. Measure the token cost first
ctx psql -c "SELECT * FROM users" | jq '.tokens'
# Result: 25000

# 2. Refine the query and execute safely
ctx psql -c "SELECT status, COUNT(*) FROM users GROUP BY status" | llm
# Result: ~100 tokens consumed (99.6% reduction)
```
This simple pattern transforms an expensive operation into a negligible one.

## Key Features

- **Universal Tool Wrapper**: Works with any CLI tool out-of-the-box (`psql`, `git`, `docker`, `kubectl`, `ls`, etc.).
- **Structured JSON Output**: Enriches raw output with vital metadata for LLMs.
  ```json
  {
    "tokens": 42,
    "output": "...",
    "input": "...",
    "metadata": { "success": true, "exit_code": 0, "duration": 127 },
    "telemetry": { "trace_id": "...", "span_id": "..." },
    "schema_version": "0.1"
  }
  ```
- **Precise Token Counting**: Supports OpenAI, Anthropic, and Gemini tokenizers.
- **Safety Controls**: Set limits on tokens, output size, lines, and pipeline stages to prevent runaway costs.
- **Streaming Support**: Terminate long-running commands in real-time when limits are breached.
- **Privacy-Aware**: Instantly disable history and telemetry with a `--private` flag.

## Installation & Updates

### Recommended: Install Script (Full Features)
The recommended method is the one-liner script, which provides full versioning and auto-update capabilities:
```bash
curl -sSL https://raw.githubusercontent.com/slavakurilyak/ctx/main/scripts/install-remote.sh | bash
```

**Features:**
- ✅ Full version information
- ✅ Auto-update notifications  
- ✅ `ctx update` command support
- ✅ Proper build metadata

### Alternative Methods

**Pre-built Releases (Full Features):**
Download binaries from the [**Releases page**](https://github.com/slavakurilyak/ctx/releases/latest).
- ✅ Full version information
- ✅ `ctx update` command support  
- ❌ No auto-update notifications

**Go Install (Limited Features):**
```bash
go install github.com/slavakurilyak/ctx@latest
```
- ❌ Shows "ctx version unknown (built from source)"
- ❌ No auto-update capabilities
- ✅ Automatic Go module updates
- 💡 **Tip:** Run `ctx update` after installing to upgrade to a full-featured version

### Updating ctx

**Automatic Updates (Install Script & Pre-built):**
```bash
ctx update                    # Update to latest stable version
ctx update --pre-release      # Include pre-releases
ctx update --check            # Check for updates without installing
```

**Go Install Users:**
```bash
go install github.com/slavakurilyak/ctx@latest  # Manual update
# Or upgrade to full version:
ctx update  # Replaces go install version with full-featured version
```

Verify installation: `ctx version`

## Usage

Simply prefix any command you want to enrich with `ctx`.

```bash
# Basic command
ctx ls -la

# Database query with a specific tokenizer
ctx --token-model openai psql -c "SELECT id FROM users LIMIT 10"

# Pipeline: filter data before tokenizing
ctx cat large.json | jq .items[0:5]

# AI Agent workflow: generate a commit message
ctx git diff --staged | claude -p 'Generate a conventional commit message.'
```

## AI Assistant Setup

`ctx` supports major AI agents and agentic IDEs, making it easy to integrate token-aware command execution into your existing AI-powered development workflow. With support for 10+ popular tools, you can teach your AI coding assistant to use `ctx` automatically.

```bash
# Default: Set up Claude Code (generates CLAUDE.md)
ctx setup

# Or set up a specific AI tool:
ctx setup claude    # Explicit Claude Code setup
ctx setup cursor    # Cursor IDE
ctx setup aider     # Aider
ctx setup windsurf  # Windsurf IDE
ctx setup jetbrains # JetBrains AI Assistant
ctx setup gemini    # Gemini CLI
```
These commands create local configuration files that instruct your IDE's AI to wrap shell commands with `ctx`, promoting token-efficient workflows.

## Coming Soon: `ctx Pro`

**`ctx Pro`** enhances your ctx experience with powerful analytics and team collaboration features for just **$10/month**.

- **Web Dashboard**: Access a comprehensive dashboard to visualize your command-line usage patterns, most frequently used tools, and command history analytics.
- **Usage Insights**: Discover which tools consume the most tokens, identify optimization opportunities, and track your token savings over time.
- **Team Features**: Centralized policies, audit logs, and shared configurations for teams to maintain consistency across development workflows.
- **Support Open Source**: Your subscription directly supports the continued development and maintenance of the ctx project, ensuring regular updates and new features.

Pro features will activate with `ctx login`. [**Learn More About Pro Pricing**](docs/PRICING.md).

## Configuration

`ctx` is configured via CLI flags, environment variables, or a `~/.config/ctx/config.yaml` file.

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

## Resources & Contributing

- **Announcement Blog**: [Read the announcement of `ctx`](https://slavakurilyak.com/posts/introducing-ctx)
- **Repository**: [https://github.com/slavakurilyak/ctx](https://github.com/slavakurilyak/ctx)
- **Issues**: [https://github.com/slavakurilyak/ctx/issues](https://github.com/slavakurilyak/ctx/issues)
- **Versioning**: [docs/VERSIONING.md](docs/VERSIONING.md)

## Acknowledgments

- [Anthropic](https://anthropic.com) for pioneering token-efficient tool use
- The creators of countless open-source CLI tools

## License

This project is licensed under the [MIT License](LICENSE).
