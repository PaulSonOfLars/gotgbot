#!/usr/bin/env bash
# git golangci-lint pre-commit hook
#
# To use, store as .git/hooks/pre-push inside your repository and make sure
# it has execute permissions.
set -euo pipefail

# Run golangci-lint to lint against all go code in the repo. Configuration in .golangci.yaml.
golangci-lint run

# Ensure all generated code has up to date.
./scripts/ci/ensure-generated.sh >/dev/null
