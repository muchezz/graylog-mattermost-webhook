#!/bin/bash

# Bulk Deployment Script for Multiple Graylog Servers
# Deploys the webhook service to multiple servers at once

set -e

echo "=========================================="
echo "Graylog Webhook - Bulk Deploy Script"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
SERVERS_FILE="${1:-servers.txt}"
BINARY_PATH="${2:-./graylog-webhook}"
CONFIG_PATH="${3:-./config.yaml}"

# Check arguments
if [ ! -f "$SERVERS_FILE" ]; then
    echo -e "${RED}Error: servers.txt not found${NC}"
    echo ""
    echo "Usage: $0 [servers_file] [binary_path] [config_path]"
    echo ""
    echo "Example servers.txt:"
    echo "  user@server1.com"
    echo "  user@server2.com"
    echo "  user@server3.com"
    exit 1
fi

if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Error: Binary not found at $BINARY_PATH${NC}"
    echo "Please build first: make build"
    exit 1
fi

if [ ! -f "$CONFIG_PATH" ]; then
    echo -e "${RED}Error: Config not found at $CONFIG_PATH${NC}"
    echo "Please copy config.yaml"
    exit 1
fi

BINARY_SIZE=$(ls -lh "$BINARY_PATH" | awk '{print $5}')
echo -e "${GREEN}Binary: $BINARY_PATH ($BINARY_SIZE)${NC}"
echo -e "${GREEN}Config: $CONFIG_PATH${NC}"
echo -e "${GREEN}Servers: $SERVERS_FILE${NC}"
echo ""

# Read servers from file
mapfile -t SERVERS < "$SERVERS_FILE"
TOTAL_SERVERS=${#SERVERS[@]}

echo -e "${YELLOW}Found $TOTAL_SERVERS servers:${NC}"
for server in "${SERVERS[@]}"; do
    echo "  - $server"
done
echo ""

read -p "Deploy to all servers? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Deployment cancelled."
    exit 0
fi

echo ""
FAILED_SERVERS=()
SUCCEEDED_SERVERS=()

# Deploy to each server
for i in "${!SERVERS[@]}"; do
    SERVER="${SERVERS[$i]}"
    COUNT=$((i + 1))
    
    echo -e "${YELLOW}[$COUNT/$TOTAL_SERVERS] Deploying to $SERVER...${NC}"
    
    # Check SSH connection
    if ! ssh -o ConnectTimeout=5 "$SERVER" "echo ok" &>/dev/null; then
        echo -e "${RED}✗ Cannot connect to $SERVER${NC}"
        FAILED_SERVERS+=("$SERVER")
        continue
    fi
    
    # Copy binary
    if ! scp -q "$BINARY_PATH" "$SERVER":~/ 2>/dev/null; then
        echo -e "${RED}✗ Failed to copy binary to $SERVER${NC}"
        FAILED_SERVERS+=("$SERVER")
        continue
    fi
    
    # Copy config
    if ! scp -q "$CONFIG_PATH" "$SERVER":~/ 2>/dev/null; then
        echo -e "${RED}✗ Failed to copy config to $SERVER${NC}"
        FAILED_SERVERS+=("$SERVER")
        continue
    fi
    
    # Deploy service
    ssh -q "$SERVER" << 'ENDSSH' || {
        echo -e "${RED}✗ Failed to install service on $SERVER${NC}"
        FAILED_SERVERS+=("$SERVER")
        continue
    }
        set -e
        
        # Make binary executable
        chmod +x ~/graylog-webhook
        
        # Install
        sudo cp ~/graylog-webhook /usr/local/bin/
        sudo mkdir -p /etc/graylog-webhook
        sudo cp ~/config.yaml /etc/graylog-webhook/
        
        # Create service user
        sudo useradd -r -s /bin/false graylog-webhook 2>/dev/null || true
        
        # Set permissions
        sudo chown graylog-webhook:graylog-webhook /etc/graylog-webhook -R
        
        # Create systemd service
        sudo tee /etc/systemd/system/graylog-webhook.service > /dev/null << 'EOFSERVICE'
[Unit]
Description=Graylog Webhook Service
After=network.target

[Service]
Type=simple
User=graylog-webhook
Environment="CONFIG_FILE=/etc/graylog-webhook/config.yaml"
ExecStart=/usr/local/bin/graylog-webhook
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOFSERVICE
        
        # Enable and start
        sudo systemctl daemon-reload
        sudo systemctl enable graylog-webhook
        sudo systemctl restart graylog-webhook
        
        # Verify
        sleep 1
        curl -s http://localhost:8080/health | grep -q healthy
ENDSSH
    
    echo -e "${GREEN}✓ $SERVER deployed successfully${NC}"
    SUCCEEDED_SERVERS+=("$SERVER")
done

echo ""
echo "=========================================="
echo -e "${GREEN}Deployment Complete!${NC}"
echo "=========================================="
echo ""
echo "Results:"
echo -e "${GREEN}  ✓ Successful: ${#SUCCEEDED_SERVERS[@]}${NC}"
echo -e "${RED}  ✗ Failed: ${#FAILED_SERVERS[@]}${NC}"
echo ""

if [ ${#SUCCEEDED_SERVERS[@]} -gt 0 ]; then
    echo -e "${GREEN}Successful servers:${NC}"
    for server in "${SUCCEEDED_SERVERS[@]}"; do
        echo "  ✓ $server"
    done
fi

if [ ${#FAILED_SERVERS[@]} -gt 0 ]; then
    echo -e "${RED}Failed servers:${NC}"
    for server in "${FAILED_SERVERS[@]}"; do
        echo "  ✗ $server"
    done
fi

echo ""
echo "Next steps on each server:"
echo "  1. Edit config: sudo nano /etc/graylog-webhook/config.yaml"
echo "  2. Check status: sudo systemctl status graylog-webhook"
echo "  3. View logs: sudo journalctl -u graylog-webhook -f"
echo ""
