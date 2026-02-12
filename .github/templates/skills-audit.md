# Skills Audit Mission

## Objective
Review all markdown files in the `**/skills` directory (and subdirectories). Your goal is to refine these skill definitions to be high-quality, safe, and efficient instructions for AI agents, ensuring they are perfectly aligned with the current state of the codebase.

## Review Guidelines
Perform a deep review of each skill file against the following criteria:

1.  **Codebase Alignment & Safety**:
    -   **Analyze the current codebase**: Verify that the skill's instructions (file paths, tool usage, commands) are accurate and reflect current best practices in this repository.
    -   **Identify Unsafe Patterns**: Correct any risky command execution (e.g., destructive commands without safeguards) or hardcoded sensitive data.
    -   **Fix Ambiguity**: Flag or rewrite any instructions that are vague, provide multiple contradictory ways to do things, or could lead to unstable agent behavior.

2.  **Elite Conciseness & Clarity**:
    -   **Strip Meta-Talk**: Remove fluff, introductory filler, or "conversational" instructions (e.g., "In this file, we will...").
    -   **Eliminate Verbosity**: Compress instructions into dense, imperative, and actionable steps. Every word must serve a technical purpose.
    -   **Remove redundant descriptions**: If a tool or process is already standard, don't over-explain it.

## Actions
-   **Refactor In-Place**: Rewrite the skill files (`SKILL.md` and related documentation) to apply these improvements.
-   **Flag Unresolved Issues**: If a skill references a system or tool that no longer exists or is broken, but you don't know the replacement, add a comment: `<!-- ISSUE: [Description] -->`.
-   **Validate**: Ensure the resulting markdown is valid and follows the repository's formatting standards.
