#!/bin/bash
set -e
GOOS=linux GOARCH=amd64 go build
upx ./image-resize
