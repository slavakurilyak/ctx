#!/usr/bin/env bash

# ctx installer script
# This script builds and installs ctx to make it globally available

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

print_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ first."
        echo "Visit https://go.dev/dl/ for installation instructions."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        print_error "Go version $GO_VERSION is too old. Please upgrade to Go 1.21 or later."
        exit 1
    fi
    
    print_info "Go $GO_VERSION detected ✓"
}

# Detect the install directory
detect_install_dir() {
    # Check if ~/bin exists and is in PATH
    if [[ ":$PATH:" == *":$HOME/bin:"* ]] || [[ ":$PATH:" == *":~/bin:"* ]]; then
        INSTALL_DIR="$HOME/bin"
        NEEDS_SUDO=false
        print_info "Will install to $INSTALL_DIR (user directory)"
    # Check if ~/.local/bin exists and is in PATH
    elif [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
        INSTALL_DIR="$HOME/.local/bin"
        NEEDS_SUDO=false
        print_info "Will install to $INSTALL_DIR (user directory)"
    # Check if /usr/local/bin is in PATH
    elif [[ ":$PATH:" == *":/usr/local/bin:"* ]]; then
        INSTALL_DIR="/usr/local/bin"
        NEEDS_SUDO=true
        print_info "Will install to $INSTALL_DIR (requires sudo)"
    else
        print_warning "No standard directory found in PATH"
        echo "Your PATH: \"$PATH\""
        echo ""
        echo "Where would you like to install ctx?"
        echo "1) $HOME/bin (will be added to PATH)"
        echo "2) $HOME/.local/bin (will be added to PATH)"
        echo "3) /usr/local/bin (requires sudo)"
        echo "4) Custom location"
        read -p "Choose [1-4]: " choice
        
        case $choice in
            1)
                INSTALL_DIR="$HOME/bin"
                NEEDS_SUDO=false
                NEEDS_PATH_UPDATE=true
                ;;
            2)
                INSTALL_DIR="$HOME/.local/bin"
                NEEDS_SUDO=false
                NEEDS_PATH_UPDATE=true
                ;;
            3)
                INSTALL_DIR="/usr/local/bin"
                NEEDS_SUDO=true
                ;;
            4)
                read -p "Enter custom directory: " INSTALL_DIR
                INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}"  # Expand ~
                read -p "Does this directory require sudo? (y/n): " needs_sudo
                NEEDS_SUDO=$([[ "$needs_sudo" == "y" ]] && echo true || echo false)
                ;;
            *)
                print_error "Invalid choice"
                exit 1
                ;;
        esac
    fi
    
    # Create directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        print_info "Creating directory $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
}

# Detect version information
detect_version_info() {
    # Try to get version from git tag
    if command -v git &> /dev/null && [ -d ".git" ]; then
        VERSION=$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "")
        if [ -z "$VERSION" ]; then
            VERSION="dev-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
        fi
        COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    else
        VERSION="dev"
        COMMIT="unknown"
    fi
    
    # Always set the date
    DATE=$(date -u +%Y-%m-%d 2>/dev/null || echo "unknown")
    
    print_info "Version info: $VERSION (commit: $COMMIT, date: $DATE)"
}

# Build ctx
build_ctx() {
    print_info "Building ctx..."
    
    # Download dependencies
    print_info "Downloading dependencies..."
    go mod download
    
    # Detect version information
    detect_version_info
    
    # Build the binary with version information
    print_info "Compiling ctx with version $VERSION..."
    LDFLAGS="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
    go build -ldflags "$LDFLAGS" -o ctx .
    
    if [ ! -f "ctx" ]; then
        print_error "Build failed - ctx binary not created"
        exit 1
    fi
    
    print_info "Build successful ✓"
}

# Install ctx
install_ctx() {
    print_info "Installing ctx to $INSTALL_DIR..."
    
    if [ "$NEEDS_SUDO" = true ]; then
        print_warning "sudo required for installation to $INSTALL_DIR"
        sudo cp ctx "$INSTALL_DIR/ctx"
        sudo chmod +x "$INSTALL_DIR/ctx"
    else
        cp ctx "$INSTALL_DIR/ctx"
        chmod +x "$INSTALL_DIR/ctx"
    fi
    
    print_info "Installation successful ✓"
    
    # Mark installation method in config
    print_info "Configuring installation method..."
    "$INSTALL_DIR/ctx" config set-installation install-script 2>/dev/null || {
        print_warning "Could not set installation method - config will be created on first use"
    }
}

# Update PATH if needed
update_path() {
    if [ "${NEEDS_PATH_UPDATE:-false}" = true ]; then
        print_info "Adding $INSTALL_DIR to PATH..."
        
        # Detect shell
        SHELL_NAME=$(basename "$SHELL")
        
        case $SHELL_NAME in
            bash)
                RC_FILE="$HOME/.bashrc"
                [ -f "$HOME/.bash_profile" ] && RC_FILE="$HOME/.bash_profile"
                ;;
            zsh)
                RC_FILE="$HOME/.zshrc"
                ;;
            fish)
                RC_FILE="$HOME/.config/fish/config.fish"
                ;;
            *)
                RC_FILE="$HOME/.profile"
                ;;
        esac
        
        # Add to PATH
        if ! grep -q "$INSTALL_DIR" "$RC_FILE" 2>/dev/null; then
            echo "" >> "$RC_FILE"
            echo "# Added by ctx installer" >> "$RC_FILE"
            if [ "$SHELL_NAME" = "fish" ]; then
                echo "set -gx PATH \"$INSTALL_DIR\" \$PATH" >> "$RC_FILE"
            else
                echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$RC_FILE"
            fi
            print_info "Added $INSTALL_DIR to PATH in $RC_FILE"
            print_warning "Run 'source $RC_FILE' or restart your shell to update PATH"
        fi
    fi
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    # Check if ctx is in PATH
    if command -v ctx &> /dev/null; then
        CTX_PATH=$(which ctx)
        CTX_VERSION=$(ctx --version 2>/dev/null || echo "version unknown")
        print_info "ctx installed successfully at $CTX_PATH"
        print_info "Version: $CTX_VERSION"
        
        # Test ctx
        print_info "Testing ctx..."
        if ctx echo "Hello from ctx!" &> /dev/null; then
            print_info "ctx is working correctly ✓"
            echo ""
            echo -e "${GREEN}Installation complete!${NC}"
            echo ""
            echo "You can now use ctx globally:"
            echo "  ctx <any-command>"
            echo ""
            echo "Examples:"
            echo "  ctx ls -la"
            echo "  ctx git status"
            echo "  ctx docker ps"
            echo ""
            echo -e "${GREEN}🚀 Next Steps - Set up your AI coding assistant:${NC}"
            echo ""
            echo "  ctx setup claude     # Claude Code/Desktop"
            echo "  ctx setup cursor     # Cursor IDE"
            echo "  ctx setup aider      # Aider"
            echo "  ctx setup windsurf   # Windsurf IDE"
            echo "  ctx setup jetbrains  # JetBrains AI Assistant"
            echo "  ctx setup gemini     # Gemini CLI"
            echo ""
            echo "  Or run 'ctx setup' to see all supported tools."
            echo ""
            echo "For more information, visit:"
            echo "  https://github.com/slavakurilyak/ctx"
        else
            print_warning "ctx installed but test failed"
        fi
    else
        print_warning "ctx installed to $INSTALL_DIR but not found in PATH"
        echo "You may need to:"
        echo "1. Run 'source ~/.bashrc' (or your shell's RC file)"
        echo "2. Restart your terminal"
        echo "3. Add \"$INSTALL_DIR\" to your PATH manually"
    fi
}

# Main installation flow
main() {
    echo "======================================"
    echo "     ctx Installer"
    echo "======================================"
    echo ""
    
    # Check prerequisites
    check_go
    
    # Detect installation directory
    detect_install_dir
    
    # Build ctx
    build_ctx
    
    # Install ctx
    install_ctx
    
    # Update PATH if needed
    update_path
    
    # Verify installation
    verify_installation
}

# Run main function
main "$@"