#!/usr/bin/env sh
set -euo pipefail

docker build --target=app -t pgapi:latest .