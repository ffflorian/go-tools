#!/usr/bin/env bash

set -e

for DIR in *; do
  if [ -d "${DIR}" ]; then
    (
      cd "${DIR}" >/dev/null 2>&1 || true
      if ls *.go >/dev/null 2>&1; then
        echo "#### Building \"${DIR}\" ..."
        go build -v .
        echo
      fi
    )
  fi
done
