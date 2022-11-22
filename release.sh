#!/usr/bin/env sh

VERSION="1.11"
FLAGS="-X main.AppVersion=$VERSION -s -w"

env GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-arm64" .
env GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-amd64" .
env GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-linux-amd64" .
env GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-win.exe" .