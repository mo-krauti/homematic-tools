#!/bin/bash
set -euxo pipefail

cd $1
rsync $1-arm homematic:/usr/local/addons/homematic-tools/
