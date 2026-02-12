#!/bin/bash
set -euo pipefail

echo "::group::Installing gh-copilot extension"

# Check if gh copilot is already working (could be built-in or pre-installed)
if gh copilot --help >/dev/null 2>&1; then
  echo "gh copilot is already available."
else
  echo "Installing github/gh-copilot extension..."
  gh extension install github/gh-copilot --force
fi

echo "gh-copilot installed/verified successfully"
echo "::endgroup::"
