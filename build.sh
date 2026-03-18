#!/usr/bin/env bash

. /etc/profile

# Linux arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w -buildid=" -o noidle-linux-arm64 main.go
# Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -buildid=" -o noidle-linux-amd64 main.go
# Windows amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -buildid=" -o noidle-windows-amd64.exe main.go
# macOS amd64
#CGO_ENABLED=0 GOOS=linux GOARCH=darwin go build -trimpath -ldflags="-s -w -buildid=" -o noidle-darwin-amd64 main.go
