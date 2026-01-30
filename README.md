# Graylog Webhook Service

**A lightweight Go service that forwards Graylog alerts to Slack or Mattermost.**

Replaces the unmaintained Java plugin with a simple, 10MB binary that runs as a systemd service.

## Features

✅ **Slack Support** - Send alerts to Slack workspaces  
✅ **Mattermost Support** - Send alerts to Mattermost instances  
✅ **Simple** - Single binary, minimal configuration  
✅ **Lightweight** - 10-15MB binary, ~10MB memory usage  
✅ **Fast** - <100ms per alert, 1000+ alerts/min  
✅ **Easy to deploy** - Copy binary, configure, run  
✅ **Easy to manage** - systemd service file included  
✅ **Production ready** - Error handling, logging, graceful shutdown  
✅ **Configurable** - YAML config + environment variables  
✅ **Severity-based routing** - Route alerts to different channels by severity  

## Quick Start (5 minutes)

### 1. Build the Binary

```bash
git clone https://github.com/muchezz/graylog-webhook.git
cd graylog-webhook
make build
```

Result: `graylog-webhook` binary (~10-15MB)

### 2. Create Slack or Mattermost Webhook

**For Slack:**
1. Go to https://api.slack.com/apps
2. Create New App → From scratch
3. Enable Incoming Webhooks
4. Create New Webhook to Workspace
5. Copy the webhook URL

**For Mattermost:**
1. Go to System Settings → Integrations → Incoming Webhooks
2. Create New Incoming Webhook
3. Select a channel and copy the webhook URL

### 3. Configure

```bash
cp config.example.yaml config.yaml
nano config.yaml
```

Set your webhook URL and platform:

```yaml
destination:
  platform: "slack"  # or "mattermost"
  webhook_url: "https://..."
  channel: "#alerts"
```

### 4. Test

```bash
./graylog-webhook
```

You should see:
```
{"level":"info","msg":"Starting Graylog Webhook Service"}
{"level":"info","msg":"Listening for connections","addr":"0.0.0.0:8080"}
```

Test health:
```bash
curl http://localhost:8080/health
# Response: {"status":"healthy"}
```

### 5. Deploy as Service

```bash
sudo make install
sudo systemctl start graylog-webhook
sudo systemctl enable graylog-webhook
```

### 6. Configure Graylog Alert

In Graylog:
- **Alerts** → **Event Definitions** → Create/Edit
- **Add Notification** → **HTTP Notification**
- **URL**: `http://localhost:8080/webhook`
- **Method**: POST

### 7. Test

Trigger a test alert in Graylog and check Slack/Mattermost!

## Configuration

### Minimal Setup

```yaml
destination:
  platform: "mattermost"
  webhook_url: "https://mattermost.example.com/hooks/xxx"
```

### With Custom Channel

```yaml
destination:
  platform: "slack"
  webhook_url: "https://hooks.slack.com/services/xxx"
  channel: "#incidents"
  username: "Graylog Alerts"
```

### With Severity Routing

```yaml
destination:
  webhook_url: "https://..."
  destinations:
    "0": "#critical-alerts"     # Emergency
    "2": "#error-alerts"        # Error
    "3": "#warning-alerts"      # Warning
    "5": "#info-alerts"         # Info
```

All other severities fall back to the default `channel`.

## Environment Variables

Alternatively, use environment variables:

```bash
PLATFORM="mattermost"                                    # slack or mattermost
WEBHOOK_URL="https://mattermost.example.com/hooks/xxx"
CHANNEL="#alerts"
USERNAME="Graylog"
LISTEN_ADDR="0.0.0.0:8080"
LOG_LEVEL="info"
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
    "event_message": "Test alert",
    "priority": "2",
    "message": "Test message"
  }'
# Response: {"status":"ok"}
```

Message should appear in your Slack/Mattermost channel!

## Monitoring

```bash
# Check service status
sudo systemctl status graylog-webhook

# View logs
sudo journalctl -u graylog-webhook -f

# Check resource usage
ps aux | grep graylog-webhook
```

Expected: ~10-15MB memory, <1% CPU idle

## Troubleshooting

### Service won't start

```bash
sudo journalctl -u graylog-webhook -n 50
```

Most common: Missing `webhook_url` in config.

### Webhook URL error

Set `webhook_url` in `/etc/graylog-webhook/config.yaml` or `WEBHOOK_URL` environment variable.

### Alerts not appearing

1. Test the webhook URL with curl:
   ```bash
   curl -X POST "https://your-webhook-url" \
     -H "Content-Type: application/json" \
     -d '{"text":"Test message"}'
   ```

2. Check webhook is enabled in Slack/Mattermost settings

3. Check firewall between services

4. Review logs: `sudo journalctl -u graylog-webhook -f`

## Performance

| Metric | Value |
|--------|-------|
| Binary size | 10-15MB |
| Memory usage | ~10-15MB |
| CPU (idle) | <1% |
| Alert latency | <100ms |
| Max throughput | 1000+/min |

## Deployment Options

### Same Server as Graylog (Simplest)

```bash
# Copy binary and config
sudo cp graylog-webhook /usr/local/bin/
sudo mkdir -p /etc/graylog-webhook
sudo cp config.yaml /etc/graylog-webhook/

# Run as service
sudo cp graylog-webhook.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable graylog-webhook
sudo systemctl start graylog-webhook

# In Graylog: http://localhost:8080/webhook
```

### Separate Server

```bash
# Same as above, but use server IP in Graylog:
# http://webhook-server-ip:8080/webhook
```

### Multiple Markets

One webhook service handles all Graylog servers:
```
Market 1 Graylog → http://webhook-service:8080/webhook
Market 2 Graylog → http://webhook-service:8080/webhook
Market 3 Graylog → http://webhook-service:8080/webhook
```

All alerts flow to same Slack/Mattermost webhook.

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Build for Multiple Platforms

```bash
make release
# Creates: graylog-webhook-linux-amd64, darwin-amd64, etc.
```

## License

Apache 2.0 - See LICENSE file

## Contributing

Issues? Features? Fork and submit a PR! The code is simple and well-documented.

---

**Ready to get started?**

1. Clone: `git clone https://github.com/muchezz/graylog-webhook.git`
2. Build: `make build`
3. Configure: `cp config.example.yaml config.yaml` + edit
4. Test: `./graylog-webhook`
5. Deploy: `sudo make install`

Enjoy your Graylog alerts in Slack or Mattermost! 🚀
