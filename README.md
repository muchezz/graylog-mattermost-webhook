# Graylog Mattermost Webhook

A lightweight Go service that forwards Graylog alerts to Mattermost.

**Why not a plugin?** Standalone services are simpler to manage, update, and don't require Graylog restarts.

## Quick Start

### 1. Build the Binary

```bash
git clone https://github.com/muchezz/graylog-mattermost-webhook.git
cd graylog-mattermost-webhook
go build -o graylog-mattermost-webhook
```

### 2. Create Mattermost Webhook

In Mattermost:
1. Go to System Settings → Integrations → Incoming Webhooks
2. Create a new webhook and copy the URL
3. Save it somewhere safe

### 3. Configure the Service

Create `/etc/graylog-mattermost-webhook/config.yaml`:

```yaml
server:
  listen_addr: "0.0.0.0:8080"
  log_level: "info"

mattermost:
  webhook_url: "https://your-mattermost.com/hooks/xxx"
  username: "Graylog"
  icon_emoji: ":clipboard:"
  channel: "#alerts"
  # Optional: Route by severity
  destinations:
    "0": "#critical-alerts"    # Emergency
    "2": "#error-alerts"       # Error
    "3": "#warning-alerts"     # Warning
```

### 4. Deploy as Systemd Service

Create `/etc/systemd/system/graylog-mattermost-webhook.service`:

```ini
[Unit]
Description=Graylog Mattermost Webhook
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=graylog-webhook
Environment="CONFIG_FILE=/etc/graylog-mattermost-webhook/config.yaml"
ExecStart=/usr/local/bin/graylog-mattermost-webhook
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Then:

```bash
# Create user
sudo useradd -r -s /bin/false graylog-webhook

# Copy binary
sudo cp graylog-mattermost-webhook /usr/local/bin/
sudo chmod +x /usr/local/bin/graylog-mattermost-webhook

# Create config directory
sudo mkdir -p /etc/graylog-mattermost-webhook
sudo cp config.yaml /etc/graylog-mattermost-webhook/
sudo chown -R graylog-webhook:graylog-webhook /etc/graylog-mattermost-webhook

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable graylog-mattermost-webhook
sudo systemctl start graylog-mattermost-webhook

# Check status
sudo systemctl status graylog-mattermost-webhook
sudo journalctl -u graylog-mattermost-webhook -f
```

### 5. Configure Graylog Alert

In Graylog, create an alert:

1. Go to **Alerts** → **Event Definitions**
2. Create/edit an event definition
3. Add a **Notification** with:
   - **Type**: HTTP Notification
   - **URL**: `http://localhost:8080/webhook` (or your webhook service address)
   - **Method**: POST
4. Save and test

## Configuration

### Environment Variables

Alternatively, use environment variables:

```bash
MATTERMOST_WEBHOOK_URL="https://..."
MATTERMOST_USERNAME="Graylog"
MATTERMOST_CHANNEL="#alerts"
LISTEN_ADDR="0.0.0.0:8080"
LOG_LEVEL="info"
```

### Severity-Based Routing

Route different alert severities to different channels:

```yaml
mattermost:
  destinations:
    "0": "#critical-alerts"     # Emergency
    "1": "#critical-alerts"     # Alert
    "2": "#error-alerts"        # Error
    "3": "#warning-alerts"      # Warning
    "4": "#info-alerts"         # Notice
    "5": "#info-alerts"         # Info
    "6": "#debug-alerts"        # Debug
```

## Testing

### Health Check

```bash
curl http://localhost:8080/health
# Response: {"status":"healthy"}
```

### Test Alert

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event_definition_id": "test-001",
    "event_timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "event_message": "Test alert from Graylog",
    "priority": "2",
    "message": "This is a test message",
    "full_message": "Full test message"
  }'
```

## Troubleshooting

### Service won't start

```bash
# Check logs
sudo journalctl -u graylog-mattermost-webhook -n 50

# Verify config syntax
./graylog-mattermost-webhook  # Run directly to see errors
```

### Missing MATTERMOST_WEBHOOK_URL

The webhook URL is **required**. Set it in:
- Config file: `mattermost.webhook_url`
- Environment: `MATTERMOST_WEBHOOK_URL`

### Messages not appearing in Mattermost

1. Verify webhook URL is correct (test with curl)
2. Check Mattermost webhook is enabled in System Settings
3. Review logs: `sudo journalctl -u graylog-mattermost-webhook -f`
4. Ensure Graylog can reach the service

## Development

### Build

```bash
go build -o graylog-mattermost-webhook
```

### Test

```bash
go test ./...
```

### Run Locally

```bash
export MATTERMOST_WEBHOOK_URL="your-webhook-url"
./graylog-mattermost-webhook
```

## Architecture

```
Graylog Server
    ↓ (HTTP POST to /webhook)
    ↓
This Service
    ├─ Parse JSON alert
    ├─ Determine severity
    ├─ Select channel
    └─ Format message
    ↓ (HTTP POST to Mattermost)
    ↓
Mattermost Incoming Webhook
    ↓
Message in channel
```

## Performance

- **Memory**: ~10MB
- **CPU**: <5% idle
- **Latency**: <100ms per alert
- **Throughput**: 1000+/min on single instance

## License

Apache 2.0

## Contributing

Fork, modify, and submit a PR! The code is simple and well-documented.
