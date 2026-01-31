# Graylog Mattermost Webhook

**A production-ready Go service that forwards Graylog alerts to Slack or Mattermost**

Forward your Graylog alerts to Slack or Mattermost with:
- ✅ Rich message formatting with details
- ✅ Severity-based color coding
- ✅ Event metadata (Source, ID, Type)
- ✅ Support for both Slack and Mattermost
- ✅ Lightweight (~10MB binary)
- ✅ Fast deployment
- ✅ No compilation needed on servers (pre-built binaries)

---
<img width="1296" height="344" alt="image" src="https://github.com/user-attachments/assets/01e1252c-2d81-4b3e-8113-66c6d5e33eea" />

## Installation Methods

Choose one of three methods below:

---

## Method 1: One-Command Deploy Script (Easiest) 

**No Git needed! Downloads and builds everything automatically.**

On your Graylog server, run:
```bash
curl -fsSL https://raw.githubusercontent.com/muchezz/graylog-mattermost-webhook/main/deploy.sh | bash
```

The script will:
- ✅ Check for Go installation
- ✅ Download source files from GitHub
- ✅ Build the binary
- ✅ Create config.yaml
- ✅ Test the binary
- ✅ Show deployment instructions

Then follow the on-screen instructions:
```bash
1. Edit the configuration:
   nano config.yaml

2. Set your Mattermost webhook URL in config.yaml

3. Test the service:
   CONFIG_FILE=./config.yaml ./graylog-webhook

4. Deploy as a service:
   sudo cp graylog-webhook /usr/local/bin/
   sudo mkdir -p /etc/graylog-webhook
   sudo cp config.yaml /etc/graylog-webhook/

5. Create systemd service file and enable:
   sudo systemctl enable graylog-webhook
   sudo systemctl start graylog-webhook
```

**Requirements:** Go 1.21+ (check with `go version`)

---

## Method 2: Pre-Built Binary Download (Fastest) 

**No compilation needed! Download ready-to-use binary.**

### 1. Download Pre-Built Binary

No Go installation required!
```bash
# For Linux AMD64 (most servers)
wget https://github.com/muchezz/graylog-mattermost-webhook/releases/download/v1.0.0/graylog-mattermost-webhook-linux-amd64
chmod +x graylog-mattermost-webhook-linux-amd64

# For other platforms, see releases page:
# https://github.com/muchezz/graylog-mattermost-webhook/releases
```

### 2. Create Configuration
```bash
cat > config.yaml << 'EOF'
server:
  listen_addr: "0.0.0.0:8080"
  log_level: "info"

destination:
  platform: "mattermost"
  webhook_url: "https://your-mattermost.com/hooks/xxxxx"
  username: "Graylog"
  icon_emoji: ":clipboard:"
  channel: "#alerts"
EOF
```

Get your Mattermost webhook URL:
1. Go to **System Settings** → **Integrations** → **Incoming Webhooks**
2. **Create New Webhook**
3. Select a channel and copy the webhook URL

### 3. Test It
```bash
CONFIG_FILE=./config.yaml ./graylog-mattermost-webhook-linux-amd64
```

You should see:
```
{"level":"info","msg":"Starting Graylog Webhook Service"}
{"level":"info","msg":"Listening for connections","addr":"0.0.0.0:8080"}
```

### 4. Configure Graylog

In Graylog:
1. **Alerts** → **Event Definitions**
2. Create/Edit an alert
3. **Add Notification** → **HTTP Notification**
4. Set:
   - **URL:** `http://localhost:8080/webhook`
   - **Method:** POST
5. **Test Notification**

Check your Mattermost channel - you should see the alert! 🎉

### 5. Deploy as Service
```bash
# Copy binary
sudo cp graylog-mattermost-webhook-linux-amd64 /usr/local/bin/graylog-webhook

# Create config directory
sudo mkdir -p /etc/graylog-webhook
sudo cp config.yaml /etc/graylog-webhook/

# Create service user
sudo useradd -r -s /bin/false graylog-webhook 2>/dev/null || true

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

# Set permissions
sudo chown graylog-webhook:graylog-webhook /etc/graylog-webhook -R

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable graylog-webhook
sudo systemctl start graylog-webhook

# Verify
sudo systemctl status graylog-webhook
```

---

## Method 3: Build from Source

**Full control - build exactly what you need.**
```bash
# Clone the repository
git clone https://github.com/muchezz/graylog-mattermost-webhook
cd graylog-mattermost-webhook

# Build for your system
make build

# Or build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o graylog-mattermost-webhook-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o graylog-mattermost-webhook-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o graylog-mattermost-webhook-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o graylog-mattermost-webhook-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o graylog-mattermost-webhook-windows-amd64.exe .
```

Then follow steps 2-5 from Method 2 above.

---

## Comparison: Which Method?

| Method | Easiest | Fastest | Customizable |
|--------|---------|---------|--------------|
| **Deploy Script** | ✅ Yes | ❌ No | ✅ Yes |
| **Pre-Built Binary** | ✅ Yes | ✅ Yes | ❌ No |
| **Build from Source** | ❌ No | ❌ No | ✅ Yes |

**Recommendation:**
- **Single server:** Use Method 1 (Deploy Script)
- **Multiple servers:** Use Method 2 (Pre-Built Binary)
- **Development:** Use Method 3 (Build from Source)

---

## Configuration

### Basic Setup
```yaml
destination:
  platform: "mattermost"
  webhook_url: "https://your-mattermost.com/hooks/xxxxx"
  channel: "#alerts"
```

### Advanced Setup
```yaml
server:
  listen_addr: "0.0.0.0:8080"
  log_level: "info"

destination:
  platform: "mattermost"          # or "slack"
  webhook_url: "https://..."      # Required
  username: "Graylog"             # Bot name
  icon_emoji: ":clipboard:"       # Bot icon
  channel: "#alerts"              # Default channel
  
  # Route by severity
  destinations:
    "0": "#critical-alerts"       # Emergency
    "1": "#critical-alerts"       # Alert
    "2": "#error-alerts"          # Error
    "3": "#warning-alerts"        # Warning
    "4": "#info-alerts"           # Notice
    "5": "#info-alerts"           # Info
    "6": "#debug-alerts"          # Debug
```

### Environment Variables
```bash
PLATFORM="mattermost"
WEBHOOK_URL="https://..."
CHANNEL="#alerts"
USERNAME="Graylog"
LISTEN_ADDR="0.0.0.0:8080"
LOG_LEVEL="info"
```

---

## Features

### Rich Message Formatting

Messages include:
- **Alert Title** - Main message from Graylog
- **Severity** - Color-coded (Red=Critical, Orange=Error, Yellow=Warning, Blue=Info)
- **Timestamp** - When the alert was triggered
- **Source** - Where the alert came from
- **Event ID** - Graylog event definition ID
- **Trigger ID** - Alert trigger ID
- **Definition Type** - Type of alert

### Severity Colors

| Level | Color | Name |
|-------|-------|------|
| 0-1 | 🔴 Red | Emergency/Alert |
| 2 | 🟠 Orange | Error |
| 3 | 🟡 Yellow | Warning |
| 4-5 | 🔵 Blue | Notice/Info |
| 6 | 🔷 Dark Blue | Debug |


### Both Slack and Mattermost

Use the same binary for both! Just change the config:

**For Mattermost:**
```yaml
destination:
  platform: "mattermost"
  webhook_url: "https://mattermost.example.com/hooks/xxx"
```

**For Slack:**
```yaml
destination:
  platform: "slack"
  webhook_url: "https://hooks.slack.com/services/xxx/yyy/zzz"
```

---

## Testing

### Health Check
```bash
curl http://localhost:8080/health
# Response: {"status":"healthy"}
```

### Manual Alert
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event_definition_id": "test-001",
    "event_timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "event_message": "Test alert message",
    "priority": "2",
    "source": "test-source"
  }'
```

### View Logs
```bash
# Live logs
sudo journalctl -u graylog-webhook -f

# Last 50 lines
sudo journalctl -u graylog-webhook -n 50
```

---

## Troubleshooting

### Service won't start
```bash
sudo systemctl status graylog-webhook
sudo journalctl -u graylog-webhook -n 50
```

Most common: Missing `webhook_url` in config.

### Messages not appearing in Mattermost

1. Test webhook URL directly:
```bash
   curl -X POST "https://your-mattermost.com/hooks/xxx" \
     -H "Content-Type: application/json" \
     -d '{"text":"Test"}'
```

2. Check webhook is enabled in Mattermost settings

3. Check firewall allows outbound HTTPS

4. Review logs: `sudo journalctl -u graylog-webhook -f`

### Connection refused

Make sure the service is running:
```bash
sudo systemctl start graylog-webhook
sudo systemctl status graylog-webhook
```

---

## Performance

- **Binary Size:** 10-15MB
- **Memory Usage:** ~10-15MB
- **CPU (idle):** <1%
- **Latency:** <100ms per alert
- **Throughput:** 1000+ alerts/minute

---

## Deployment Options

### Single Server (Same as Graylog)
```bash
sudo cp graylog-mattermost-webhook-linux-amd64 /usr/local/bin/
# Configure and run
```

Graylog uses: `http://localhost:8080/webhook`

### Multiple Servers

Deploy to all servers:
```bash
for server in server1 server2 server3; do
  scp graylog-mattermost-webhook-linux-amd64 user@$server:~/
  ssh user@$server 'sudo mv graylog-mattermost-webhook-linux-amd64 /usr/local/bin/ && sudo systemctl restart graylog-webhook'
done
```

---

## License

Apache 2.0 - Free to use, modify, and distribute

---

## Support

- 📖 Check documentation above
- 🐛 Report issues on GitHub
- 🚀 Contribute improvements via PR
- 💬 GitHub Discussions for questions

---

## Repository

**GitHub:** https://github.com/muchezz/graylog-mattermost-webhook
```bash
git clone https://github.com/muchezz/graylog-mattermost-webhook
cd graylog-mattermost-webhook
```

---

**Ready to get started?** Choose your installation method above and follow the steps! 🚀
