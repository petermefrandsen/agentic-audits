# Version Bump Recommendation Mission

## Objective
Analyze the changes in the current Pull Request compared to the `main` branch and recommend the appropriate semantic version bump (`patch`, `minor`, or `major`).

## Context
- You have access to the full repository.
- Use `git diff main...HEAD` to see the changes introduced by this PR.

## SemVer Criteria
- **patch**: Backward-compatible bug fixes or minor internal improvements.
- **minor**: New features, non-breaking modifications, or significant internal changes that are backward-compatible.
- **major**: Breaking changes, removed APIs, or fundamental Shifts in architecture that are NOT backward-compatible.

## Instructions
1.  **Analyze the Diff**: Examine all changed files and identify the nature of the modifications.
2.  **Determine Bump Level**: Choose one of: `patch`, `minor`, or `major`.
3.  **Provide Rationale**: Summarize WHY you chose this level (e.g., "Added new public method X in component Y" -> `minor`).
4.  **Write Output**:
    - Write the chosen level (`patch`, `minor`, or `major`) to a file named `bump_level.txt`.
    - Do NOT include any other text in `bump_level.txt`.
    - Write your full analysis and rationale to `bump_analysis.md`.

## Actions
- Use `run_command` to execute `git diff` if needed, although your context window should already contain the relevant files if `context_files` was set appropriately.
- Use `write_to_file` to create `bump_level.txt` and `bump_analysis.md`.
