#!/bin/bash
set -euxo pipefail

cd $1
go run $1.go
