#!/bin/bash
set -euxo pipefail

cd $1
env GOOS=linux GOARCH=arm GOARM=5 go build -o $1-arm $1.go
