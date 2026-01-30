.PHONY: help build clean test run install

help:
	@echo "Graylog Mattermost Webhook - Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run locally (requires MATTERMOST_WEBHOOK_URL)"
	@echo "  make install     - Install to /usr/local/bin"
	@echo "  make release     - Build for multiple platforms"

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o graylog-mattermost-webhook .
	@echo "Built: ./graylog-mattermost-webhook"

clean:
	rm -f graylog-mattermost-webhook graylog-mattermost-webhook-*
	go clean
	@echo "Cleaned"

test:
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests passed"

run: build
	MATTERMOST_WEBHOOK_URL="http://localhost:9000/test" ./graylog-mattermost-webhook

install: build
	sudo cp graylog-mattermost-webhook /usr/local/bin/
	sudo chmod +x /usr/local/bin/graylog-mattermost-webhook
	sudo cp graylog-mattermost-webhook.service /etc/systemd/system/
	sudo mkdir -p /etc/graylog-mattermost-webhook
	sudo cp config.example.yaml /etc/graylog-mattermost-webhook/config.yaml
	sudo useradd -r -s /bin/false graylog-webhook 2>/dev/null || true
	sudo chown graylog-webhook:graylog-webhook /etc/graylog-mattermost-webhook -R
	sudo systemctl daemon-reload
	@echo "Installed. Configure /etc/graylog-mattermost-webhook/config.yaml and run:"
	@echo "  sudo systemctl enable graylog-mattermost-webhook"
	@echo "  sudo systemctl start graylog-mattermost-webhook"

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o graylog-mattermost-webhook-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o graylog-mattermost-webhook-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o graylog-mattermost-webhook-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o graylog-mattermost-webhook-darwin-arm64 .
	@echo "Built: graylog-mattermost-webhook-*"
