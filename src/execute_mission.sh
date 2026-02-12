#!/bin/bash
set -euo pipefail

# Inputs:
# GH_TOKEN: authentication token
# RESOLVED_MISSION: the mission content (env var)
# CONTEXT_FILES: inputs.context_files
# MODEL: inputs.model
# FALLBACK_MODEL: inputs.fallback_model
# DRY_RUN: inputs.dry_run
# PR_BASE: inputs.pr_base
# PR_BRANCH: inputs.pr_branch
# PR_TITLE: inputs.pr_title
# PR_BODY: inputs.pr_body
# PR_LABELS: inputs.pr_labels
# GITHUB_REPOSITORY: github.repository context

# Explicitly export tokens to ensure Copilot CLI sees them
export COPILOT_GITHUB_TOKEN="${GH_TOKEN}"
export GITHUB_TOKEN="${GH_TOKEN}"

# Get the mission content from the env var set in previous step
AGENT_MISSION="${RESOLVED_MISSION:-}"

if [[ -z "$AGENT_MISSION" ]]; then
  echo "::error::RESOLVED_MISSION env var is empty or unset."
  exit 1
fi

echo "Context:        ${CONTEXT_FILES:-.}"
echo "Model:          ${MODEL:-auto}"
echo "Fallback Model: ${FALLBACK_MODEL:-none}"
echo "---"

# Build the full prompt including web sources
FULL_MISSION="${AGENT_MISSION} (context files: ${CONTEXT_FILES:-.})"
if [[ -n "${EXTRA_WEB_SOURCES:-}" ]]; then
  FULL_MISSION="${FULL_MISSION}. ${EXTRA_WEB_SOURCES}"
fi

# ── Append PR Instructions ──
if [[ "${DRY_RUN:-false}" == "false" ]]; then
  FULL_MISSION="${FULL_MISSION}

  
  ### MANDATORY: Pull Request Creation
  You MUST create a Pull Request for your changes using the \`create_pull_request\` tool from the GitHub MCP server. 
  
  PR Specifications:
  - **Repository**: ${GITHUB_REPOSITORY:-}
  - **Base Branch**: ${PR_BASE:-main}
  - **Branch Name**: ${PR_BRANCH:-$(echo "agent/audit-$(date +%s)")}
  
  - **Title**: ${PR_TITLE:-Use STRICT Conventional Commits format (e.g., refactor(skills): [AI-GENERATED] audit and clarify instructions).}
  
  - **Body**: ${PR_BODY:-You MUST provide a comprehensive, elite-quality description structured as follows:
      ### 🔎 Audit Overview
      Provide a high-level technical summary of what was audited and the general state of the skills.
      
      ### 🛠 Detailed Changes
      Provide a per-skill breakdown of specific technical improvements (e.g., Skill X: Removed 40% verbosity, updated paths to match current source tree).
      
      ### ⚠️ Manual Review Required
      List any specific files where you added <!-- ISSUE --> comments because they require human intervention.}
  
  - **Labels**: ${PR_LABELS:-automated-pr}
  "
else
  FULL_MISSION="${FULL_MISSION}
  
  NOTE: dry_run is set to TRUE. Do NOT create a Pull Request. Just verify the changes and report what you would have done.
  "
fi

# Build model flags
MODEL_FLAGS=""
if [[ -n "${MODEL:-}" ]]; then
  MODEL_FLAGS="--model ${MODEL}"
fi

# ── Attempt with primary model ──
run_agent() {
  local flags="${1:-}"
  # shellcheck disable=SC2086
  gh copilot \
    --allow-all-tools \
    -p "${FULL_MISSION}" \
    ${flags} \
    2>&1
}

if run_agent "${MODEL_FLAGS}"; then
  echo "Agent mission completed successfully."
else
  EXIT_CODE=$?
  echo "::warning::Primary model failed (exit code: ${EXIT_CODE})"

  if [[ -n "${FALLBACK_MODEL:-}" ]]; then
    echo "Retrying with fallback model: ${FALLBACK_MODEL}"
    FALLBACK_FLAGS="--model ${FALLBACK_MODEL}"
    if run_agent "${FALLBACK_FLAGS}"; then
      echo "Agent mission completed with fallback model."
    else
      echo "::error::Agent mission failed with both primary and fallback models"
      exit 1
    fi
  else
    echo "::error::Agent mission failed and no fallback model is configured"
    exit 1
  fi
fi
