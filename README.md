# Agentic Audits

[![CI](https://github.com/petermefrandsen/agentic-audits/actions/workflows/test-action.yml/badge.svg)](https://github.com/petermefrandsen/agentic-audits/actions/workflows/test-action.yml)
[![Go Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/petermefrandsen/19edde6ec0fe2f62bb54778db2eaad40/raw/coverage.json)](https://github.com/petermefrandsen/agentic-audits/actions/workflows/test-action.yml)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Agentic%20Audits-blue.svg)](https://github.com/marketplace/actions/agentic-audits)

Automated governance for AI assets using headless GitHub Copilot agents - CLI agnostic and user-controlled.

## 🛡️ Why CLI Agnostic?

GitHub's native agents are great, but they are often locked into specific ecosystems. **Agentic Audits** puts you in control:

- **CLI Agnostic**: Run any agent, any model, and any MCP server without lock-in.
- **Full Control**: You own the mission, the context, and the execution flow.
- **Headless & Scalable**: Designed for complex, cross-repo CI/CD workflows.

## 🚀 Key Features

- **Headless Execution**: Run complex agentic workflows in CI/CD without human intervention.
- **MCP Native**: Easily install and use MCP servers (like Upstash, GitHub, Slack) via a simple YAML config.
- **Automated PRs**: The agent can automatically refactor code and submit structured, high-quality Pull Requests.
- **Cross-Repo Ready**: Use a central hub of "skills" and templates to audit many repositories.

## 📦 Usage

Add this to your `.github/workflows/agent-audit.yml`:

```yaml
- name: Run Headless Agent
  uses: petermefrandsen/agentic-audits@v0.0.1
  with:
    mission: "Audit the README.md and suggest 3 improvements."
    context_files: "README.md"
    github_token: ${{ secrets.COPILOT_GOV_TOKEN }}
    model: "gpt-5-mini"
    fallback_model: "gpt-4.1"
```

## ⚙️ Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `mission` | ⚠️ | — | The mission prompt. Required if `template` is not used. |
| `template` | ⚠️ | — | Path to a mission template in `.github/templates/`. |
| `github_token` | ✅ | — | GitHub token with Copilot access. |
| `context_files`| | `.` | Files or globs for the agent to consider. |
| `model` | | *(auto)*| Primary Copilot model to use (e.g. `gpt-5-mini`, `gpt-4.1`). |
| `fallback_model` | | *(none)*| Fallback model if the primary model hits a quota or error. |
| `sources_config`| | `.github/sources.yml` | YAML config for MCP servers and web docs. |
| `dry_run` | | `false` | If `true`, skips PR creation. |

## 🛠️ Configuration (`sources.yml`)

Configure external tools and documentation for your agent:

```yaml
sources:
  - name: context7
    type: mcp
    package: "@upstash/context7-mcp"
    enabled: true
```

## 🛡️ Setup Requirements

1. **Secret**: Add `COPILOT_GOV_TOKEN` to your repo secrets.
2. **Permissions**: Ensure your workflow has `contents: write` and `pull-requests: write`.
