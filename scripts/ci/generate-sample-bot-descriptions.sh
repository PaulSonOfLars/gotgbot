#!/usr/bin/env bash
set -euo pipefail

SAMPLES_DIR="samples"
README_FILE="${SAMPLES_DIR}/README.md"

# Single >, to recreate file

intro="# Sample bots

This is a short description of all the sample bots present in this directory. These are intended to be a source of
inspiration and learning for anyone looking to use this library.
"

echo "$intro" >"${README_FILE}"

for d in "${SAMPLES_DIR}"/*; do
  if [ ! -d "${d}" ]; then
    # Only check directories.
    continue
  fi

  echo "Reading ${d}..."

  intro="
## ${d}
"

  echo "$intro" >>"${README_FILE}"

  # Extract all comments before the main function, and dump them in the readme file.
  description="$(sed -n -e '/package/,/func main/ p' "${d}/main.go" | (grep -E "^//" || true) | sed -E 's:^// ?::g')"
  if [[ -z "${description}" ]]; then
    echo "!!! no doc comments for ${d} - please add some first."
    exit 1
  fi
  echo "${description}" >>"${README_FILE}"
done
