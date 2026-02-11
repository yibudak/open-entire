# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
make build          # Build binary to bin/entire
make test           # go test ./... -v
make test-race      # go test -race ./...
make lint           # golangci-lint run ./...
make dev            # fmt → vet → test → build

# Run a single test
go test -run TestParseJSONL ./internal/agent/claude/

# Run a single package's tests
go test -v ./internal/checkpoint/
```

Version info is injected via ldflags (`-X main.version`, `-X main.commit`, `-X main.date`). The entry point is `cmd/entire/main.go`.

## Architecture

This is a Go CLI (`github.com/spf13/cobra`) that captures AI coding agent sessions as Git-stored checkpoints. All state lives in Git or `.entire/` JSON files — no database.

### Data Flow

```
Git Hook (post-commit/pre-push)
  → hooks.Handler dispatches event
    → strategy.Strategy decides whether to checkpoint
      → checkpoint.Store.Create() writes session data
        → git.Repository.CommitOnBranch() persists to orphan branch "entire/checkpoints/v1"
          → git.Repository.AddTrailer() appends Entire-Checkpoint trailer to user commit
```

### Key Design Decisions

- **Git operations use `os/exec`**, not go-git. The `git.Repository` struct wraps all git commands via `run()` helper. Reads and writes both go through exec.
- **Checkpoints live on an orphan branch** (`entire/checkpoints/v1`) to keep them separate from code history. The `CommitOnBranch` method temporarily checks out the target branch, writes files, commits, then returns to the original branch.
- **Sharded storage**: checkpoint ID `a3b2c4d5e6f7` maps to path `a3/b2c4d5e6f7/` on the checkpoints branch.
- **Web assets are embedded** via `//go:embed` in `internal/web/embed.go`. Templates and static files compile into the binary.

### Two Extension Points

**Agents** (`internal/agent/agent.go`): Registry pattern with `Agent` interface. Claude Code is the only implementation (`internal/agent/claude/`). New agents self-register via `init()` calling `agent.Register()`. The interface requires `Detect()` (find active session), `ParseSession()` (parse transcript), and `SessionPaths()`.

**Strategies** (`internal/strategy/strategy.go`): `Strategy` interface with three hooks: `OnCommit`, `OnAgentResponse`, `OnPush`. Two implementations exist: `ManualCommit` (checkpoint only on git commit) and `AutoCommit` (checkpoint on every agent response + commit). Selected via `strategy.New(name)`.

### Config System

4-layer merge, highest priority first:
1. Env vars (`ENTIRE_ENABLED`, `ENTIRE_STRATEGY`, `ENTIRE_LOG_LEVEL`, `ENTIRE_TELEMETRY`)
2. `.entire/settings.local.json` (gitignored)
3. `.entire/settings.json` (committed)
4. `~/.config/entire/settings.json` (global)
5. Built-in defaults in `internal/config/defaults.go`

### Package Dependencies

`cli/` commands use `git.Repository` for git ops, `checkpoint.Store` for reading/writing checkpoints, `session.Store` for local state, and `config.Load()` for configuration. The `hooks.Handler` bridges hook events to strategy implementations. The `web.Server` reads checkpoints via `checkpoint.Store` and serves embedded HTML/CSS/JS through chi router.

### Claude Code Parser

`internal/agent/claude/parser.go` processes JSONL transcripts from `~/.claude/projects/<encoded-path>/<session-id>.jsonl`. Key behaviors: deduplicates streaming assistant events by `requestId` (keeps last for token usage), extracts tool calls from `tool_use` content blocks, aggregates token usage across all requests, and handles a 10MB line buffer for large messages.

### Testing Patterns

Tests use `stretchr/testify` (`assert`/`require`), `t.TempDir()` for filesystem isolation, and write fixture files inline. No mocks — tests either exercise pure functions or create real temp directories with `.git/hooks/` structure.
