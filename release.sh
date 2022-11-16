#!/usr/bin/env sh

env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "dist/scan_compare-1.0-mac-arm64" .
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "dist/scan_compare-1.0-mac-amd64" .
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "dist/scan_compare-1.0-linux-amd64" .
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "dist/scan_compare-1.0-win.exe" .