#!/bin/bash

# Graylog Webhook Deployment Script
# This script downloads, builds, and deploys the webhook service
# No git clone required!

set -e

echo "=========================================="
echo "Graylog Webhook Service - Deploy Script"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.21 or later from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go installed: $GO_VERSION${NC}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
esac

echo -e "${GREEN}✓ Detected: $OS / $ARCH${NC}"
echo ""

# Create working directory
WORK_DIR="/tmp/graylog-webhook-build-$$"
mkdir -p "$WORK_DIR"
cd "$WORK_DIR"

echo -e "${YELLOW}Downloading webhook service...${NC}"

# Download source files from GitHub
REPO_URL="https://github.com/muchezz/graylog-mattermost-webhook"
BRANCH="main"

# Download each file individually (no git needed)
FILES=(
    "main.go"
    "config.go"
    "handler.go"
    "graylog.go"
    "notification.go"
    "go.mod"
)

for file in "${FILES[@]}"; do
    echo "  Downloading $file..."
    curl -fsSL "$REPO_URL/raw/$BRANCH/$file" -o "$file"
done

echo -e "${GREEN}✓ Source files downloaded${NC}"
echo ""

# Download go.sum if available (optional)
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod download || true
go mod tidy || true

echo ""
echo -e "${YELLOW}Building binary...${NC}"

# Build the binary
CGO_ENABLED=0 go build -ldflags="-w -s" -o graylog-webhook .

if [ ! -f graylog-webhook ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

BINARY_SIZE=$(ls -lh graylog-webhook | awk '{print $5}')
echo -e "${GREEN}✓ Binary built: graylog-webhook ($BINARY_SIZE)${NC}"
echo ""

# Create default config if not exists
if [ ! -f config.yaml ]; then
    echo -e "${YELLOW}Creating default configuration...${NC}"
    cat > config.yaml << 'EOFCONFIG'
server:
  listen_addr: "0.0.0.0:8080"
  log_level: "info"

destination:
  platform: "mattermost"
  webhook_url: "https://mattermost.example.com/hooks/xxx"
  username: "Graylog"
  icon_emoji: ":clipboard:"
  channel: "#alerts"
  destinations:
    "0": "#critical-alerts"
    "2": "#error-alerts"
    "3": "#warning-alerts"
EOFCONFIG
    echo -e "${GREEN}✓ Config created: config.yaml${NC}"
    echo -e "${YELLOW}⚠️  Edit config.yaml and set your webhook_url!${NC}"
fi

echo ""
echo -e "${YELLOW}Testing binary...${NC}"

# Quick test
if timeout 2 ./graylog-webhook &>/dev/null || true; then
    echo -e "${GREEN}✓ Binary test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Binary test incomplete (expected)${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Build Complete!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Edit the configuration:"
echo "   nano config.yaml"
echo ""
echo "2. Set your Mattermost webhook URL in config.yaml"
echo ""
echo "3. Test the service:"
echo "   CONFIG_FILE=./config.yaml ./graylog-webhook"
echo ""
echo "4. Deploy as a service:"
echo "   sudo cp graylog-webhook /usr/local/bin/"
echo "   sudo mkdir -p /etc/graylog-webhook"
echo "   sudo cp config.yaml /etc/graylog-webhook/"
echo ""
echo "5. Create systemd service file and enable:"
echo "   sudo systemctl enable graylog-webhook"
echo "   sudo systemctl start graylog-webhook"
echo ""
echo "Current directory: $WORK_DIR"
echo ""
