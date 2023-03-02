#!/usr/bin/env sh

VERSION="1.22"
FLAGS="-X main.AppVersion=$VERSION -s -w"

GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-arm64" . && \
GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-amd64" . && \
GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-linux-amd64" . && \
GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-win.exe" . && \

docker build -t antfie/scan_compare:$VERSION . && \
docker build -t antfie/scan_compare . && \
docker push antfie/scan_compare:$VERSION && \
docker push antfie/scan_compare && \

ESCAPE=$'\e'
echo "${ESCAPE}[0;32mSuccess${ESCAPE}[0m"