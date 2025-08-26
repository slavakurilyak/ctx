# ctx Versioning Strategy

## Overview

ctx uses two independent version numbers:
1. **Software Version**: The version of the ctx tool itself
2. **Schema Version**: The version of the JSON output format

## Software Versioning

ctx follows [Semantic Versioning 2.0.0](https://semver.org/):

**Format**: `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)

- **MAJOR**: Incompatible CLI changes or schema major version changes
- **MINOR**: New functionality in a backward-compatible manner
- **PATCH**: Backward-compatible bug fixes

### Examples:
- `1.0.0` → `1.0.1`: Bug fix (tokenizer counting error)
- `1.0.1` → `1.1.0`: New feature (added `--private` flag)
- `1.1.0` → `2.0.0`: Breaking change (removed deprecated flags)

## Schema Versioning

The JSON output schema uses simplified semantic versioning:

**Format**: `MAJOR.MINOR` (e.g., `1.0`)

- **MAJOR**: Breaking changes to output structure
- **MINOR**: Additions to output (new optional fields)

### Examples:
- `1.0` → `1.1`: Added optional `trace_id` field
- `1.1` → `2.0`: Changed `metadata.tokens` from number to object

## Version Compatibility Matrix

| ctx Version | Schema Version | Notes |
|------------|---------------|--------|
| 1.0.x      | 1.0          | Initial release |
| 1.1.x      | 1.0          | Added features, same output |
| 1.2.x      | 1.1          | Added telemetry fields |
| 2.0.x      | 2.0          | Major schema redesign |

## Implementation Guidelines

### 1. Schema Version in Code

Located in `internal/models/output.go`:
```go
const CurrentSchemaVersion = "1.0"
```

**When to update**:
- Adding optional fields: Increment minor (1.0 → 1.1)
- Changing/removing fields: Increment major (1.1 → 2.0)

### 2. Software Version in Builds

Set during build time via ldflags:
```bash
go build -ldflags "-X main.version=1.2.3"
```

### 3. Version Command

Users can check both versions:
```bash
$ ctx version
ctx version: 1.2.3
schema version: 1.0
commit: abc123f
built: 2024-01-15
```

### 4. Backward Compatibility

When schema version changes:
- **Minor changes**: Old consumers should ignore unknown fields
- **Major changes**: Consider supporting multiple schema versions with a flag:
  ```bash
  ctx --schema-version 1.0 command  # Force old schema
  ```

## Release Process

### 1. Patch Release (e.g., 1.0.0 → 1.0.1)
- Bug fixes only
- No schema changes
- No new features

### 2. Minor Release (e.g., 1.0.1 → 1.1.0)
- New features or improvements
- May include schema minor version bump
- Backward compatible

### 3. Major Release (e.g., 1.2.3 → 2.0.0)
- Breaking changes
- Schema major version changes
- Migration guide required

## Creating a Release

### Prerequisites
- All tests passing (`go test ./...`)
- Version updated if needed
- Changes committed to main branch

### Release Steps

1. **Create an annotated tag:**
```bash
# Format: v{MAJOR}.{MINOR}.{PATCH}
git tag -a v0.1.4 -m "Release v0.1.4"
```

2. **Push the tag to trigger automated release:**
```bash
git push origin v0.1.4
```

3. **Monitor the release:**
```bash
# Watch GitHub Actions progress
gh run watch

# View the release once complete
gh release view v0.1.4
```

### What Happens Next

When you push a tag:
1. GitHub Actions workflow triggers automatically
2. GoReleaser builds binaries for all platforms (Linux, macOS, Windows)
3. Changelog is auto-generated from commit messages
4. GitHub Release is created with all artifacts
5. Installation scripts are updated to use the new version

### Commit Message Convention

For automatic changelog generation, use conventional commits:
- `feat:` New features → appears in "🚀 Features"
- `fix:` Bug fixes → appears in "🐛 Bug Fixes"  
- `docs:` Documentation → appears in "📚 Documentation"
- `perf:` Performance improvements → appears in "⚡ Performance"
- `refactor:` Code refactoring → appears in "♻️ Refactoring"
- `chore:` Maintenance tasks → appears in "🔧 Other"

Example:
```bash
git commit -m "feat: add support for VS Code custom instructions"
git commit -m "fix: correct token counting for large outputs"
```

### Emergency Rollback

If a release has critical issues:
```bash
# Delete the tag locally and remotely
git tag -d v0.1.4
git push origin :v0.1.4

# Fix the issue, then create a new release
```

## Checking Versions

### For Software Version:
```bash
ctx --version  # or ctx version
```

### For Schema Version:
Check any JSON output:
```json
{
  "schema_version": "1.0",
  ...
}
```

## Migration Strategy

When breaking changes are necessary:

1. **Deprecation Warning**: Add warnings in version N
2. **Dual Support**: Support both old and new in version N+1
3. **Removal**: Remove old support in version N+2

Example:
- v1.8.0: Add deprecation warning for `--old-flag`
- v1.9.0: Support both `--old-flag` and `--new-flag`
- v2.0.0: Remove `--old-flag`

## Consumer Guidelines

### For AI Agents:
- Always check `schema_version` in responses
- Handle unknown fields gracefully (for minor updates)
- Fail explicitly on unsupported major versions

### For Scripts:
```bash
# Check schema version
SCHEMA=$(ctx echo test | jq -r .schema_version)
if [[ ! "$SCHEMA" =~ ^1\. ]]; then
  echo "Unsupported schema version: $SCHEMA"
  exit 1
fi
```

## Version History

| Release | Date | ctx Version | Schema | Changes |
|---------|------|------------|--------|---------|
| Beta | 2024-01-15 | 0.1.0 | 0.1 | Initial public beta |
| - | - | - | - | Provider-based tokenization |
| - | - | - | - | Resource limits & controls |
| - | - | - | - | Streaming support |
| (Future) | TBD | 0.2.0 | 0.1 | TBD |
| (Future) | TBD | 1.0.0 | 1.0 | First stable release |