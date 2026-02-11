#!/bin/bash
set -euo pipefail

# Local Agent Test Script
# Usage: export GH_TOKEN=your_token && ./local-test.sh

# Default inputs
export INPUT_MISSION="Review this file and point out any issues."
export INPUT_CONTEXT_FILES="README.md"
export INPUT_GITHUB_TOKEN="${GH_TOKEN:-}"
export INPUT_MODEL="gpt-4o-mini" 

if [ -z "$INPUT_GITHUB_TOKEN" ]; then
  echo "Error: GH_TOKEN is not set. Please export GH_TOKEN=<your-token>."
  exit 1
fi

echo "Gathering context..."
CONTEXT_CONTENT=$(cat README.md)
FULL_MISSION="Mission:
$INPUT_MISSION

Context:
$CONTEXT_CONTENT"

echo "Configuring local auth..."
# Fetch username to satisfy gh validation
USERNAME=$(curl -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" https://api.github.com/user | grep '"login":' | awk -F'"' '{print $4}')

if [ -z "$USERNAME" ]; then
    echo "Warning: Could not fetch username. Using fallback."
    USERNAME="headless-agent"
fi

mkdir -p ~/.config/gh
cat > ~/.config/gh/hosts.yml <<EOF
github.com:
    user: "$USERNAME"
    oauth_token: "${INPUT_GITHUB_TOKEN}"
    git_protocol: "https"
EOF
chmod 600 ~/.config/gh/hosts.yml

# Clear config conflict if any
rm -rf ~/.config/github-copilot

# Env vars
export COPILOT_GITHUB_TOKEN="${INPUT_GITHUB_TOKEN}"
export GITHUB_TOKEN="${INPUT_GITHUB_TOKEN}"
export GH_TOKEN="${INPUT_GITHUB_TOKEN}"

echo "Starting Agent..."
gh copilot --allow-all-tools -p "$FULL_MISSION"
