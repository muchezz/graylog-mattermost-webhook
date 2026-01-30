.PHONY: help build clean test run install

help:
	@echo "Graylog Webhook Service - Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run locally (requires WEBHOOK_URL)"
	@echo "  make install     - Install to /usr/local/bin"
	@echo "  make release     - Build for multiple platforms"

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o graylog-webhook .
	@echo "Built: ./graylog-webhook (Slack/Mattermost)"

clean:
	rm -f graylog-webhook graylog-webhook-*
	go clean
	@echo "Cleaned"

test:
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests passed"

run: build
	WEBHOOK_URL="http://localhost:9000/test" PLATFORM="mattermost" ./graylog-webhook

install: build
	sudo cp graylog-webhook /usr/local/bin/
	sudo chmod +x /usr/local/bin/graylog-webhook
	sudo cp graylog-webhook.service /etc/systemd/system/
	sudo mkdir -p /etc/graylog-webhook
	sudo cp config.example.yaml /etc/graylog-webhook/config.yaml
	sudo useradd -r -s /bin/false graylog-webhook 2>/dev/null || true
	sudo chown graylog-webhook:graylog-webhook /etc/graylog-webhook -R
	sudo systemctl daemon-reload
	@echo "Installed. Configure /etc/graylog-webhook/config.yaml and run:"
	@echo "  sudo systemctl enable graylog-webhook"
	@echo "  sudo systemctl start graylog-webhook"

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o graylog-webhook-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o graylog-webhook-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o graylog-webhook-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o graylog-webhook-darwin-arm64 .
	@echo "Built: graylog-webhook-*"
