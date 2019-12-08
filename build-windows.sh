#!/bin/bash

GOOS=windows GOARCH=386 go build -o git-sync.exe main.go
zip git-sync_windows-386 git-sync.exe config/windows/config.yaml
rm git-sync.exe