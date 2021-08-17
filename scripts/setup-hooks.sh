#!/usr/bin/env bash
# Set all the hooks defined in scripts/hooks as git hooks.
# This allows you to run ci formatting checks before pushing.
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"

# Link all hooks
ln -s "${ROOT_DIR}/scripts/hooks/pre-push.sh" "${ROOT_DIR}/.git/hooks/pre-push"
