#!/bin/bash
set -euo pipefail

# Inputs:
# GH_TOKEN: The GitHub token to use for authentication

if [ -z "${GH_TOKEN:-}" ]; then
  echo "::error::GH_TOKEN is not set"
  exit 1
fi

echo "Configuring gh auth manually to bypass scope validation..."
# Fetch username to satisfy gh validation
USERNAME=$(curl -s -H "Authorization: token $GH_TOKEN" https://api.github.com/user | grep '"login":' | awk -F'"' '{print $4}')

if [ -z "$USERNAME" ]; then
  echo "Warning: Could not fetch username. Defaulting to 'headless-agent'."
  USERNAME="headless-agent"
else
  echo "Detected username: $USERNAME"
fi

mkdir -p ~/.config/gh
cat > ~/.config/gh/hosts.yml <<EOF
github.com:
    user: "$USERNAME"
    oauth_token: "$GH_TOKEN"
    git_protocol: "https"
EOF
chmod 600 ~/.config/gh/hosts.yml

# Verify status (warn only on failure)
gh auth status || echo "gh auth status reported issues, but config is present."
