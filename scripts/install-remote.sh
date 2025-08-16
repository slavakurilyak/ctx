#!/usr/bin/env bash

# One-liner remote installer for ctx
# Usage: curl -sSL https://raw.githubusercontent.com/slavakurilyak/ctx/main/scripts/install-remote.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() { printf "${GREEN}[INFO]${NC} %s\n" "$1"; }
print_error() { printf "${RED}[ERROR]${NC} %s\n" "$1"; }
print_warning() { printf "${YELLOW}[WARNING]${NC} %s\n" "$1"; }

# Check for required tools
if ! command -v git &> /dev/null; then
    print_error "git is required but not installed"
    exit 1
fi

if ! command -v go &> /dev/null; then
    print_error "Go is required but not installed"
    echo "Visit https://go.dev/dl/ for installation instructions"
    exit 1
fi

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf -- '$TEMP_DIR'" EXIT

print_info "Downloading ctx..."

# Clone repository
git clone --quiet --depth 1 https://github.com/slavakurilyak/ctx.git "$TEMP_DIR" 2>/dev/null || {
    print_error "Failed to download ctx"
    exit 1
}

# Change to temp directory
cd "$TEMP_DIR"

# Run installer
if [ -f "install.sh" ]; then
    chmod +x install.sh
    ./install.sh
else
    print_error "Installer not found"
    exit 1
fi