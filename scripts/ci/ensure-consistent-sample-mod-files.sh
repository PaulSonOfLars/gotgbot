#!/usr/bin/env bash
set -euo pipefail

SAMPLES_DIR="samples"

REPO="github.com/PaulSonOfLars/gotgbot" # Current library import path
GO_VERSION="1.15"                       # Go version we expect our samples to be using
V_MAJOR="v2"                            # Current major version for the library

# We need to do some bash magic for systems that aren't using gnused
sedCmd="sed"
if [[ "$(uname)" == "Darwin" ]]; then
  command -v gsed || (echo "gnu-sed not installed. Please install with 'brew install gnu-sed'" && return 1)
  sedCmd="gsed"
fi

for d in "${SAMPLES_DIR}"/*; do
  if [ ! -d "${d}" ]; then
    # Only check directories.
    continue
  fi

  echo "Checking ${d}"
  # Ensure the module has the right name
  "${sedCmd}" -i "s:^module .*:module ${REPO}/${d}:g" "${d}/go.mod"

  # Ensure the current library version is a dummy version to avoid confusion
  "${sedCmd}" -i "s:${REPO}/${V_MAJOR} ${V_MAJOR}.*$:${REPO}/${V_MAJOR} ${V_MAJOR}.99.99:g" "${d}/go.mod"

  # Ensure go version is consistent
  "${sedCmd}" -i "s:^go .*:go ${GO_VERSION}:g" "${d}/go.mod"
done
