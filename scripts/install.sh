#!/bin/sh
# Open-Entire installer
# Usage: curl -fsSL https://raw.githubusercontent.com/yibudak/open-entire/main/scripts/install.sh | sh
#    or: curl -fsSL ... | sh -s -- --version v0.1.0 --dir /usr/local/bin

set -eu

REPO="yibudak/open-entire"
BINARY="open-entire"
INSTALL_DIR="/usr/local/bin"
VERSION=""

# ── Helpers ──────────────────────────────────────────────────────────────────

log()   { printf "\033[1;34m==>\033[0m %s\n" "$*"; }
ok()    { printf "\033[1;32m==>\033[0m %s\n" "$*"; }
warn()  { printf "\033[1;33m==>\033[0m %s\n" "$*" >&2; }
error() { printf "\033[1;31merror:\033[0m %s\n" "$*" >&2; exit 1; }

# ── Platform detection ───────────────────────────────────────────────────────

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       error "Unsupported OS: $(uname -s). Only Linux and macOS are supported." ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)             error "Unsupported architecture: $(uname -m). Only amd64 and arm64 are supported." ;;
    esac
}

detect_downloader() {
    if command -v curl >/dev/null 2>&1; then
        echo "curl"
    elif command -v wget >/dev/null 2>&1; then
        echo "wget"
    else
        error "Neither curl nor wget found. Please install one and retry."
    fi
}

download() {
    url="$1"
    output="$2"
    downloader="$(detect_downloader)"

    case "$downloader" in
        curl) curl -fsSL -o "$output" "$url" ;;
        wget) wget -qO "$output" "$url" ;;
    esac
}

# ── Version resolution ───────────────────────────────────────────────────────

latest_version() {
    downloader="$(detect_downloader)"
    url="https://api.github.com/repos/${REPO}/releases/latest"

    case "$downloader" in
        curl) resp="$(curl -fsSL "$url")" ;;
        wget) resp="$(wget -qO- "$url")" ;;
    esac

    # Parse tag_name without jq
    echo "$resp" | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/'
}

# ── Checksum verification ────────────────────────────────────────────────────

verify_checksum() {
    archive="$1"
    expected_file="$2"

    if [ ! -f "$expected_file" ]; then
        warn "Checksum file not found, skipping verification"
        return 0
    fi

    expected="$(awk '{print $1}' "$expected_file")"

    if command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$archive" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$archive" | awk '{print $1}')"
    else
        warn "No sha256 tool found, skipping checksum verification"
        return 0
    fi

    if [ "$actual" != "$expected" ]; then
        error "Checksum mismatch!\n  expected: ${expected}\n  got:      ${actual}"
    fi

    ok "Checksum verified"
}

# ── Parse arguments ──────────────────────────────────────────────────────────

while [ $# -gt 0 ]; do
    case "$1" in
        --version|-v)
            VERSION="$2"
            shift 2
            ;;
        --dir|-d)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help|-h)
            cat <<EOF
Open-Entire installer

Usage:
  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/scripts/install.sh | sh
  curl -fsSL ... | sh -s -- [OPTIONS]

Options:
  --version, -v VERSION   Install a specific version (default: latest)
  --dir, -d PATH          Installation directory (default: /usr/local/bin)
  --help, -h              Show this help message
EOF
            exit 0
            ;;
        *)
            error "Unknown option: $1. Use --help for usage."
            ;;
    esac
done

# ── Main ─────────────────────────────────────────────────────────────────────

main() {
    OS="$(detect_os)"
    ARCH="$(detect_arch)"

    log "Detected platform: ${OS}/${ARCH}"

    # Resolve version
    if [ -z "$VERSION" ]; then
        log "Fetching latest version..."
        VERSION="$(latest_version)"
        if [ -z "$VERSION" ]; then
            error "Could not determine latest version. Specify one with --version."
        fi
    fi

    # Strip leading 'v' for archive name
    VERSION_NUM="${VERSION#v}"

    log "Installing open-entire ${VERSION}"

    ARCHIVE="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
    CHECKSUM_URL="${DOWNLOAD_URL}.sha256"

    TMPDIR="$(mktemp -d)"
    trap 'rm -rf "$TMPDIR"' EXIT

    # Download archive
    log "Downloading ${DOWNLOAD_URL}"
    download "$DOWNLOAD_URL" "${TMPDIR}/${ARCHIVE}" || error "Download failed. Check that version ${VERSION} exists at https://github.com/${REPO}/releases"

    # Download and verify checksum
    log "Verifying checksum..."
    download "$CHECKSUM_URL" "${TMPDIR}/${ARCHIVE}.sha256" 2>/dev/null || true
    verify_checksum "${TMPDIR}/${ARCHIVE}" "${TMPDIR}/${ARCHIVE}.sha256"

    # Extract
    log "Extracting..."
    tar xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"

    # Install
    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    else
        log "Elevated permissions required to install to ${INSTALL_DIR}"
        sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
    fi
    chmod +x "${INSTALL_DIR}/${BINARY}"

    ok "Installed open-entire ${VERSION} to ${INSTALL_DIR}/${BINARY}"
    echo ""
    "${INSTALL_DIR}/${BINARY}" version
    echo ""
    ok "Run 'open-entire enable' in a Git repo to get started."
}

main
