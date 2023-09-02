#!/usr/bin/env bash
# git golangci-lint pre-commit hook
#
# To use, store as .git/hooks/pre-push inside your repository and make sure
# it has execute permissions.
set -euo pipefail

# Run golangci-lint to lint against all the code in the lib.
golangci-lint run

# Run golangci-lint against
for d in ./samples/*; do
  if [[ ! -d "$d" ]]; then
    continue
  fi

  pushd "$d"
  golangci-lint run --config=../../.golangci.yaml --path-prefix="$d"
  popd
done

# Ensure all generated code has up to date.
./scripts/ci/ensure-generated.sh >/dev/null
