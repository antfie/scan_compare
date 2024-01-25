# !/usr/bin/env sh

ESCAPE=$'\e'

if [[ -z "${VERSION}" ]]; then
    VERSION="0.0"
fi

FLAGS="-X main.AppVersion=$VERSION -s -w"

echo "${ESCAPE}[0;32mBuilding v${VERSION}...${ESCAPE}[0m"

GOOS=darwin GOARCH=arm64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-arm64" . && \
GOOS=darwin GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-mac-amd64" . && \
GOOS=linux GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-linux-amd64" . && \
GOOS=windows GOARCH=amd64 go build -ldflags="$FLAGS" -trimpath -o "dist/scan_compare-win.exe" . && \

echo "${ESCAPE}[0;32mBuild Success${ESCAPE}[0m"