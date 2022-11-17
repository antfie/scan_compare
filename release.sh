#!/usr/bin/env sh

VERSION="v1.3"

env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.AppVersion=$VERSION" -o "dist/scan_compare-mac-arm64" .
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.AppVersion=$VERSION" -o "dist/scan_compare-mac-amd64" .
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.AppVersion=$VERSION" -o "dist/scan_compare-linux-amd64" .
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.AppVersion=$VERSION" -o "dist/scan_compare-win.exe" .