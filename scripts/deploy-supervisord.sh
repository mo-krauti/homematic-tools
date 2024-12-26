#!/bin/bash
set -euxo pipefail

rsync -r deployment/supervisor homematic:/usr/local/addons/homematic-tools/
