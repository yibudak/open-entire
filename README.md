<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/status-alpha-orange?style=flat-square" alt="Status">
  <img src="https://img.shields.io/badge/agent-Claude_Code-blueviolet?style=flat-square" alt="Claude Code">
</p>

<h1 align="center">entire</h1>

<p align="center">
  <strong>Open-source AI session capture for Git.</strong><br>
  Record every AI coding interaction. Rewind, review, and attribute — all from your terminal.
</p>

<p align="center">
  <code>entire enable</code> &rarr; code with AI &rarr; <code>git commit</code> &rarr; checkpoint captured
</p>

---

## What is Entire?

**Entire** hooks into Git to automatically capture AI coding agent sessions as **checkpoints** — versioned, searchable snapshots that live on a dedicated Git branch. No external database, no cloud dependency. Just Git.

Every checkpoint records:
- Full conversation transcript (prompts, responses, tool calls)
- Token usage (input, output, cache)
- Line attribution (AI vs human contribution %)
- File diffs linked to the session

```
$ entire status
Repository: /Users/dev/myproject
Enabled:    true
Strategy:   manual-commit
Branch:     feat/auth
Active Sessions: 1
  - 8a51...3f56 (claude-code)
Checkpoints: 14
```

---

## Quick Start

```bash
# Install
go install github.com/yibudak/open-entire/cmd/entire@latest

# Enable in your repo
cd your-project
entire enable

# That's it. Code with Claude Code, commit, and checkpoints are captured.
git commit -m "feat: add auth"
# => Entire-Checkpoint: a3b2c4d5e6f7
# => Entire-Attribution: 73% agent (146/200 lines)

# Browse sessions in your browser
entire serve
```

---

## Installation

### Quick Install (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/yibudak/open-entire/main/scripts/install.sh | sh
```

Options:

```bash
# Specific version
curl -fsSL ... | sh -s -- --version v0.1.0

# Custom directory
curl -fsSL ... | sh -s -- --dir ~/.local/bin
```

### From Source

```bash
go install github.com/yibudak/open-entire/cmd/entire@latest
```

### Build Locally

```bash
git clone https://github.com/yibudak/open-entire.git
cd open-entire
make build      # => bin/entire
make install    # => $GOPATH/bin/entire
```

### Homebrew (coming soon)

```bash
brew install yibudak/tap/entire
```

---

## Commands

| Command | Description |
|---------|-------------|
| `entire enable` | Initialize Entire in a Git repo — installs hooks, creates config |
| `entire disable` | Remove hooks (data preserved) |
| `entire status` | Show capture status, active sessions, checkpoint count |
| `entire rewind` | Rewind working tree to a previous checkpoint |
| `entire resume <branch>` | Checkout branch and find associated session |
| `entire explain` | Display transcript, token usage, attribution for a checkpoint |
| `entire serve` | Launch local web viewer to browse all sessions |
| `entire clean` | Remove orphaned shadow branches |
| `entire doctor` | Find and fix stuck sessions |
| `entire reset` | Delete all local Entire state |
| `entire version` | Print build info |

### `entire enable`

```bash
entire enable                          # defaults: manual-commit strategy
entire enable --strategy auto-commit   # checkpoint after every AI response
entire enable --force                  # re-initialize existing setup
```

### `entire rewind`

```bash
entire rewind --list                   # show available checkpoints
entire rewind --to a3b2c4d5e6f7       # restore working tree
entire rewind --to a3b2c4d5 --reset   # hard reset
entire rewind --to a3b2c4d5 --logs-only  # restore session logs only
```

### `entire explain`

```bash
entire explain --checkpoint a3b2c4d5e6f7          # by checkpoint ID
entire explain --commit abc123                     # by commit hash
entire explain --checkpoint a3b2c4d5e6f7 --full    # full transcript
entire explain --checkpoint a3b2c4d5e6f7 --short   # summary only
entire explain --checkpoint a3b2c4d5e6f7 --raw-transcript  # raw JSONL
```

### `entire serve`

```bash
entire serve              # http://localhost:8080
entire serve --port 3000  # custom port
```

The web viewer provides:
- **Dashboard** — recent checkpoints, activity overview
- **Checkpoint list** — filter by branch, view diffs
- **Checkpoint detail** — code diffs, session summaries, attribution
- **Session detail** — full transcript, tool calls, token usage
- **JSON API** — `/api/checkpoints`, `/api/checkpoints/:id`, `/api/checkpoints/:id/sessions/:idx`

---

## How It Works

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────────────────┐
│  Claude Code │────▶│  Git Hooks   │────▶│  entire/checkpoints/v1      │
│  (AI Agent)  │     │  post-commit │     │  (orphan branch)            │
└─────────────┘     │  pre-push    │     │                             │
                    └──────────────┘     │  a3/b2c4d5e6f7/             │
                           │             │  ├── metadata.json           │
                    ┌──────▼──────┐      │  └── 0/                     │
                    │  Strategy   │      │      ├── full.jsonl          │
                    │  Engine     │      │      ├── context.md          │
                    └─────────────┘      │      ├── metadata.json       │
                                         │      ├── prompt.txt          │
                                         │      └── content_hash.txt    │
                                         └─────────────────────────────┘
```

1. **Enable** — `entire enable` installs Git hooks and saves config to `.entire/`
2. **Detect** — Hooks detect when an AI agent (Claude Code) is active
3. **Capture** — On commit (or agent response), a checkpoint is created
4. **Store** — Session data is committed to the `entire/checkpoints/v1` orphan branch
5. **Link** — Commit trailers (`Entire-Checkpoint`, `Entire-Attribution`) are appended

### Strategies

| Strategy | When Checkpoints Are Created | Best For |
|----------|------------------------------|----------|
| `manual-commit` (default) | On `git commit` | Main branch, clean history |
| `auto-commit` | After each AI response + on commit | Feature branches, granular tracking |

### Data Storage

All data lives in Git — no database required.

```
entire/checkpoints/v1 branch:
  <shard-2>/<remaining-10>/
  ├── metadata.json        # checkpoint ID, commit, branch, author, strategy
  └── 0/                   # session index
      ├── metadata.json    # token usage, attribution, timestamps
      ├── full.jsonl       # complete JSONL transcript
      ├── context.md       # human-readable prompts
      ├── prompt.txt       # raw prompts
      └── content_hash.txt # SHA-256 integrity hash
```

Commit trailers on user commits:
```
feat: Add user authentication

Entire-Checkpoint: a3b2c4d5e6f7
Entire-Attribution: 73% agent (146/200 lines)
```

---

## Configuration

Entire uses a 4-layer config system (highest priority first):

| Layer | Path | Notes |
|-------|------|-------|
| Environment | `ENTIRE_*` vars | Highest priority |
| Local | `.entire/settings.local.json` | Gitignored, per-developer |
| Project | `.entire/settings.json` | Committed, shared |
| Global | `~/.config/entire/settings.json` | User-wide defaults |

```json
{
  "enabled": true,
  "strategy": "manual-commit",
  "log_level": "info",
  "telemetry": false,
  "strategy_options": {
    "summarize": {
      "enabled": false
    }
  }
}
```

### Environment Variables

| Variable | Values | Default |
|----------|--------|---------|
| `ENTIRE_ENABLED` | `true`/`false`, `1`/`0` | `true` |
| `ENTIRE_STRATEGY` | `manual-commit`, `auto-commit` | `manual-commit` |
| `ENTIRE_LOG_LEVEL` | `debug`, `info`, `warn`, `error` | `info` |
| `ENTIRE_TELEMETRY` | `true`/`false` | `false` |

---

## Agent Support

| Agent | Status | Detection |
|-------|--------|-----------|
| Claude Code | Supported | JSONL session files + process detection |
| Gemini CLI | Planned | — |
| GitHub Copilot | Planned | — |
| Custom | Via `Agent` interface | Implement `Detect` + `ParseSession` |

### Claude Code Integration

Entire reads Claude Code session files from:

```
~/.claude/projects/<encoded-repo-path>/<session-id>.jsonl
```

Path encoding: `/Users/dev/myproject` → `-Users-dev-myproject`

Captured data:
- Full JSONL transcript with streaming deduplication
- Token usage aggregation (input, output, cache creation, cache reads)
- Tool calls (Write, Read, Bash, etc.)
- Nested sessions (subagents via Task tool)

---

## Development

```bash
# Build
make build

# Run tests
make test

# Run tests with race detector
make test-race

# Format + vet + test + build
make dev

# Lint (requires golangci-lint)
make lint
```

### Project Structure

```
open-entire/
├── cmd/entire/              # Entry point
├── internal/
│   ├── cli/                 # Cobra commands (11 commands)
│   ├── config/              # 4-layer config system
│   ├── logging/             # Structured logging (slog)
│   ├── git/                 # Git operations (exec-based)
│   ├── hooks/               # Hook templates + installer
│   ├── session/             # Session lifecycle management
│   ├── checkpoint/          # Checkpoint CRUD + sharding
│   ├── strategy/            # manual-commit + auto-commit
│   ├── agent/claude/        # Claude Code JSONL parser
│   ├── attribution/         # AI vs human line tracking
│   └── web/                 # Local viewer (chi + embedded assets)
├── pkg/types/               # Shared types
├── testdata/                # Test fixtures
├── Makefile
├── .goreleaser.yml
└── go.mod
```

### Adding a New Agent

Implement the `Agent` interface and register it:

```go
type Agent interface {
    Name() string
    Detect(repoDir string) (sessionID string, err error)
    ParseSession(sessionID string, repoDir string) (*types.SessionData, error)
    SessionPaths(repoDir string) types.AgentPaths
}

// In your agent's init():
func init() {
    agent.Register(&MyAgent{})
}
```

### Adding a New Strategy

Implement the `Strategy` interface:

```go
type Strategy interface {
    Name() string
    OnAgentResponse(ctx context.Context, event *AgentResponseEvent) error
    OnCommit(ctx context.Context, event *CommitEvent) error
    OnPush(ctx context.Context, event *PushEvent) error
}
```

---

## Comparison with Entire.io

| Feature | open-entire | Entire.io |
|---------|-------------|-----------|
| Session capture | Local Git branch | Cloud + Git branch |
| Web viewer | Local (`entire serve`) | Hosted (entire.io) |
| Auth | None needed | GitHub OAuth |
| Team features | — | Shared dashboards |
| Agents | Claude Code | Claude Code, Gemini CLI |
| Price | Free (MIT) | Freemium |
| Data storage | Your Git repo only | Your repo + Entire cloud |

---

## License

MIT — see [LICENSE](LICENSE).

---

<p align="center">
  <sub>Built with Go. No database. No cloud. Just Git.</sub>
</p>
