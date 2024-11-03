#!/usr/bin/env bash
set -euo pipefail

SAMPLES_DIR="samples"

REPO="github.com/PaulSonOfLars/gotgbot" # Current library import path
GO_VERSION="1.21"                       # Go version we expect our samples to be using
V_MAJOR="v2"                            # Current major version for the library
V_DUMMY="v2.99.99"                      # dummy version for the library

for d in "${SAMPLES_DIR}"/*; do
  if [ ! -d "${d}" ]; then
    # Only check directories.
    continue
  fi

  goModFile="${d}/go.mod"
  echo "Checking ${goModFile}"

  pushd "${d}" >/dev/null

  # Ensure the following are correct:
  # - module name
  # - go version
  # - library version (set to a dummy version to avoid needing constant updates)
  # - library replace to use the lib in the repo
  go mod edit \
    -module "${REPO}/${d}" \
    -go="${GO_VERSION}" \
    -require="${REPO}/${V_MAJOR}@${V_DUMMY}" \
    -replace="${REPO}/${V_MAJOR}=../../"

  go mod tidy

  popd >/dev/null
done
