# Internal Guide

> This document is synced with external documentation using the **On-Demand Doc Sync** workflow.

## Overview

This internal guide summarizes key points from the latest external documentation for GitHub Actions and GitHub Copilot CLI; internal-only notes and the doc-sync marker are preserved below.

### GitHub Actions (summary)

- GitHub Actions is a CI/CD platform for automating build, test, and deployment workflows defined in YAML files under .github/workflows.
- Workflows are triggered by events (push, pull_request, schedule, repository_dispatch, manual, and more) and consist of jobs made up of ordered steps.
- Jobs run on runners (GitHub-hosted virtual machines for Linux/Windows/macOS or self-hosted runners) and can run in parallel or depend on other jobs; matrices allow running jobs across multiple OSes or configurations.
- Actions are reusable building blocks that encapsulate tasks; use Marketplace actions or author custom actions to simplify workflow steps.

### GitHub Copilot CLI (summary)

- GitHub Copilot CLI provides terminal-native AI assistance (build, debug, refactor, and natural-language code interactions) and integrates with GitHub for repo/context access.
- Supported platforms: Linux, macOS, Windows; installation methods include Homebrew, npm, WinGet, or the official install script.
- Authentication: users sign in via the CLI (or supply a fine-grained PAT with Copilot Requests permission via GH_TOKEN/GITHUB_TOKEN); enterprise/org policies may restrict access.
- Features include model selection, experimental mode for early features (e.g., Autopilot), LSP server support (configured via ~/.copilot/lsp-config.json or .github/lsp.json), and explicit preview/consent for actions taken by the agent.

<!-- Content will be populated by the doc-sync agent. -->
