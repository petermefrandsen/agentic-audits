#!/bin/bash
set -euo pipefail

# Inputs:
# SOURCES_CONFIG: Path to the sources configuration file (optional)

SOURCES_CONFIG="${SOURCES_CONFIG:-}"

echo "::group::Configuring sources"

CONFIG_DIR="${HOME}/.config/github-copilot"
CONFIG_FILE="${CONFIG_DIR}/config.json"
mkdir -p "${CONFIG_DIR}"

if [[ -n "${SOURCES_CONFIG}" && -f "${SOURCES_CONFIG}" ]]; then
  echo "Loading sources from: ${SOURCES_CONFIG}"
  
  # Run the node script to parse config
  # Assuming the node script is in the same directory as this script
  SCRIPT_DIR="$(dirname "$(realpath "$0")")"
  PARSE_OUTPUT=$(node "${SCRIPT_DIR}/configure_sources.js" "${SOURCES_CONFIG}")
  
  # Extract values from parsed output
  MCP_JSON=$(echo "${PARSE_OUTPUT}" | node -e "
    const d = JSON.parse(require('fs').readFileSync('/dev/stdin','utf8'));
    console.log(JSON.stringify(d.mcpServers));
  ")

  WEB_SOURCES=$(echo "${PARSE_OUTPUT}" | node -e "
    const d = JSON.parse(require('fs').readFileSync('/dev/stdin','utf8'));
    console.log(d.webSources);
  ")

  # Install MCP packages
  PACKAGES=$(echo "${PARSE_OUTPUT}" | node -e "
    const d = JSON.parse(require('fs').readFileSync('/dev/stdin','utf8'));
    d.mcpPackages.forEach(p => console.log(p));
  ")

  while IFS= read -r pkg; do
    if [[ -n "${pkg}" ]]; then
      echo "Installing MCP package: ${pkg}"
      npm install -g "${pkg}" || echo "::warning::Failed to install ${pkg}"
    fi
  done <<< "${PACKAGES}"

else
  echo "No sources config provided or file not found — running without MCP servers"
  MCP_JSON='{}'
  WEB_SOURCES=""
fi

# Write the Copilot config (empty mcpServers if none configured)
node -e "
  const config = { mcpServers: ${MCP_JSON} };
  require('fs').writeFileSync('${CONFIG_FILE}', JSON.stringify(config, null, 2));
"

echo "Copilot config written to ${CONFIG_FILE}"
cat "${CONFIG_FILE}"

# Export web sources for use in the mission step
if [[ -n "${GITHUB_ENV:-}" ]]; then
  echo "EXTRA_WEB_SOURCES=${WEB_SOURCES}" >> "${GITHUB_ENV}"
else
  # For local testing/non-GitHub Actions env
  export EXTRA_WEB_SOURCES="${WEB_SOURCES}"
fi

echo "::endgroup::"
