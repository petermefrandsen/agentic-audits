# Skills Audit Mission

## Objective
Review all markdown files in the `**/skills` directory (and subdirectories). Your goal is to refine these skill definitions to be high-quality, safe, and efficient instructions for AI agents.

## Review Guidelines
Perform a deep review of each skill file against the following criteria:

1.  **Best Practices & Safety**:
    -   Identify and correct any **unsafe patterns** (e.g., risky command execution without validation, hardcoded secrets).
    -   Ensure instructions align with current codebase conventions.
    -   Flag any instructions that are **ambiguous** or could lead to unintended side effects.

2.  **Conciseness & Clarity**:
    -   **Remove unneeded descriptions**: specific context is good, but "fluff" or conversational filler must be removed.
    -   **Eliminate unnecessary verbosity**: Rewrite instructions to be direct, imperative, and actionable.
    -   Ensure the prompt is strict but provides enough context for high-quality execution.

## Actions
-   **Refactor In-Place**: Rewrite the skill files to apply the improvements found above.
-   **Flag Unresolved Issues**: If a file has issues you cannot fix (e.g., missing context), add a comment: `<!-- ISSUE: [Description] -->`.
