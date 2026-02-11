package hooks

const postCommitScript = `#!/bin/sh
# managed by entire
# Post-commit hook: capture checkpoint on commit

# Skip if entire is disabled
if [ "$ENTIRE_ENABLED" = "false" ] || [ "$ENTIRE_ENABLED" = "0" ]; then
    exit 0
fi

# Find entire binary
ENTIRE_BIN=$(command -v entire 2>/dev/null)
if [ -z "$ENTIRE_BIN" ]; then
    # Try common install paths
    for p in /usr/local/bin/entire "$HOME/go/bin/entire" "$HOME/.local/bin/entire"; do
        if [ -x "$p" ]; then
            ENTIRE_BIN="$p"
            break
        fi
    done
fi

if [ -z "$ENTIRE_BIN" ]; then
    exit 0
fi

# Run in background to not block commit
"$ENTIRE_BIN" _hook post-commit &
`

const prePushScript = `#!/bin/sh
# managed by entire
# Pre-push hook: sync checkpoints

# Skip if entire is disabled
if [ "$ENTIRE_ENABLED" = "false" ] || [ "$ENTIRE_ENABLED" = "0" ]; then
    exit 0
fi

ENTIRE_BIN=$(command -v entire 2>/dev/null)
if [ -z "$ENTIRE_BIN" ]; then
    for p in /usr/local/bin/entire "$HOME/go/bin/entire" "$HOME/.local/bin/entire"; do
        if [ -x "$p" ]; then
            ENTIRE_BIN="$p"
            break
        fi
    done
fi

if [ -z "$ENTIRE_BIN" ]; then
    exit 0
fi

"$ENTIRE_BIN" _hook pre-push
`
