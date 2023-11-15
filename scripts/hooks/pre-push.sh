#!/usr/bin/env bash
# git golangci-lint pre-commit hook
#
# To use, store as .git/hooks/pre-push inside your repository and make sure
# it has execute permissions.
set -euo pipefail

echo "Checking generated docs are up to date..."
# Ensure all generated code has up to date.
./scripts/ci/ensure-generated.sh >/dev/null

echo "Linting project..."
# Run golangci-lint to lint against all the code in the lib.
golangci-lint run

# Run golangci-lint against samples
echo "Linting samples..."
for d in ./samples/*; do
  if [[ ! -d "$d" ]]; then
    continue
  fi

  pushd "$d" >/dev/null
  golangci-lint run --config=../../.golangci.yaml --path-prefix="$d"
  popd >/dev/null
done
