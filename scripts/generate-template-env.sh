#!/bin/bash
set -euxo pipefail

cat deployment/supervisor/credentials.env | awk -F '=' '{ if ($1) { print $1 "=" } else { print }}' > template.env
