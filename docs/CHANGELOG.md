# 📋 Changelog

All notable changes to ctx will be documented in this file.

## [0.1.1] - 2025-08-17

### 🚀 Features
- Added Goose AI agent integration support via `ctx setup goose`
- Goose now appears in supported agents list with full documentation links

### 🔧 Improvements
- Migrated to GoReleaser for automated releases and version management
- Version information now properly injected at build time
- Added centralized version package for better maintainability

## [0.1.0] - 2025-08-16

### 🎉 Initial Beta Release

#### Core Features
- Universal command wrapper with structured JSON output
- Token counting support for OpenAI, Anthropic, and Gemini tokenizers
- Precise token awareness for cost optimization

#### AI Assistant Support
- Support for 11 AI coding assistants via `ctx setup` command:
  - Claude Code/Desktop
  - Cursor IDE
  - Aider
  - Windsurf IDE
  - JetBrains AI Assistant
  - Gemini CLI
  - Zed Editor
  - GitHub Copilot
  - Cline
  - Roo Code
  - Kilo Code

#### Safety & Control
- Resource limits for tokens, output bytes, lines, and pipeline stages
- Privacy-aware mode with `--private` flag
- Command history tracking (can be disabled)
- OpenTelemetry support for observability

#### Developer Experience
- Auto-update capabilities with version checking
- Multiple command invocation methods (POSIX separator, explicit subcommand, quoted)
- Streaming support for long-running commands
- Pretty output formatting option

#### Schema
- JSON output with schema version 0.1
- Structured metadata including success status, exit code, and duration
- Telemetry information with trace and span IDs

---

[0.1.1]: https://github.com/slavakurilyak/ctx/releases/tag/v0.1.1
[0.1.0]: https://github.com/slavakurilyak/ctx/releases/tag/v0.1.0