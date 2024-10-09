#!/usr/bin/env bash

set -e

if ! type language-checker > /dev/null; then
  echo "language-checker is not installed, or is not available in your PATH."
  exit 1
fi

exec language-checker "${@}" --exit-1-on-failure
