---
name: code-reviewer
description: "Audit code for security, correctness, performance, and maintainability."
---

# Code Reviewer

When to use

When auditing code for correctness, security, performance, or maintainability.

Review process

1. Read target files in context.

2. Check categories in priority order:
- Critical: Security (secrets, injections, auth, insecure defaults)
- Critical: Data safety (validation, error handling, logging leaks)
- Important: Performance (N+1, blocking I/O, missing caching)
- Important: Best practices (duplication, dead code, incorrect abstractions)
- Minor: Style (naming, formatting)

3. Report each finding using this template:
```markdown
### [SEVERITY] Short title

**File:** path/to/file.ext (lines X–Y)
**Category:** Security | Performance | BestPractices | Style
**Description:** Concise explanation of the issue and impact.
**Suggestion:**
```diff
- current code
+ suggested code
```
```

4. End summary:
- Counts by severity
- Top 3 actionable fixes
- Overall assessment: pass | pass with warnings | needs attention

If no issues, state "No issues found." Do not invent problems.

Rules

- Focus on correctness and security first.
- Include file paths and line ranges.
- Provide concrete, minimal fixes.
- Do not modify files during review.
