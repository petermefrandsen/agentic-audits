#!/bin/bash
set -euo pipefail

# Mock inputs for local testing
export INPUT_MISSION="Review this file and point out any issues."
export INPUT_CONTEXT_FILES="README.md"
export INPUT_GITHUB_TOKEN="${GH_TOKEN:-}"
export INPUT_MODEL="gpt-4o" # or gpt-5-mini if available
# export GH_DEBUG=api # Optional: debug GH CLI

if [ -z "$INPUT_GITHUB_TOKEN" ]; then
  echo "Error: GH_TOKEN is not set. Please export GH_TOKEN=<your-token>."
  exit 1
fi

# Simulate the Agent Mission step
echo "Gathering context from: $INPUT_CONTEXT_FILES"
CONTEXT_CONTENT=""
FILES=$(find . -path "./$INPUT_CONTEXT_FILES" -type f -not -path '*/.*' | head -n 20)
for f in $FILES; do
  CONTEXT_CONTENT+=$'\n\n--- FILE: '$f' ---\n'
  CONTEXT_CONTENT+=$(cat "$f")
done

FULL_MISSION="Mission:
$INPUT_MISSION

Context:
$CONTEXT_CONTENT"

echo "Full Mission Prompt:"
echo "$FULL_MISSION"

# Export tokens explicitly
export COPILOT_GITHUB_TOKEN="${INPUT_GITHUB_TOKEN}"
export GITHUB_TOKEN="${INPUT_GITHUB_TOKEN}"

echo "Running gh copilot..."
# We use 'gh copilot suggest' for old CLI or 'gh copilot' for new agent?
# Checking what's installed...
gh copilot --version || echo "gh copilot version check failed"

# Run agent
# Note: --allow-all-tools might require interactive approval if not in CI?
# In new CLI, -p is definitely the way.
gh copilot --allow-all-tools -p "$FULL_MISSION"
