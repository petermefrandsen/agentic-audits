#!/bin/bash
set -euo pipefail

# Inputs:
# INPUT_MISSION: properties.mission
# INPUT_TEMPLATE: properties.template

MISSION="${INPUT_MISSION:-}"
TEMPLATE="${INPUT_TEMPLATE:-}"

if [[ -n "$MISSION" && -n "$TEMPLATE" ]]; then
  echo "::error::Both 'mission' and 'template' inputs are provided. Please use only one."
  exit 1
fi

if [[ -z "$MISSION" && -z "$TEMPLATE" ]]; then
  echo "::error::Either 'mission' or 'template' input must be provided."
  exit 1
fi

if [[ -n "$TEMPLATE" ]]; then
  TEMPLATE_PATH=".github/templates/${TEMPLATE}.md"
  if [[ ! -f "$TEMPLATE_PATH" ]]; then
    echo "::error::Template file not found: $TEMPLATE_PATH"
    exit 1
  fi
  echo "Resolved mission from template: $TEMPLATE_PATH"
  MISSION_CONTENT=$(cat "$TEMPLATE_PATH")
else
  MISSION_CONTENT="$MISSION"
fi

# Export RESOLVED_MISSION safely
if [[ -n "${GITHUB_ENV:-}" ]]; then
{
  echo "RESOLVED_MISSION<<EOF"
  echo "$MISSION_CONTENT"
  echo "EOF"
} >> "$GITHUB_ENV"
else
  export RESOLVED_MISSION="$MISSION_CONTENT"
  echo "Resolved Mission Content:"
  echo "$MISSION_CONTENT"
fi
