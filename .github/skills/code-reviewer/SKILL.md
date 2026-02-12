---
name: code-reviewer
description: >
  Reviews code for best practices, security vulnerabilities, performance
  issues, and maintainability. Produces structured review comments with
  severity levels and actionable suggestions. Use this skill when the mission
  involves auditing or reviewing code files.
---

# Code Reviewer

## When to use

Use this skill when your mission involves reviewing, auditing, or analysing
code for quality, security, or performance issues.

## Instructions

### Review process

1. **Scan the target files** — read each file in the context to understand
   its purpose and structure.

2. **Check each category** in order of priority:

   | Priority | Category | What to look for |
   |----------|----------|-----------------|
   | 🔴 Critical | **Security** | Hardcoded secrets, injection vulnerabilities, insecure defaults, missing auth checks |
   | 🔴 Critical | **Data safety** | Unvalidated inputs, missing error handling, data leaks in logs |
   | 🟡 Important | **Performance** | N+1 queries, unnecessary allocations, missing caching, blocking I/O |
   | 🟡 Important | **Best practices** | SOLID violations, dead code, duplicated logic, missing types |
   | 🔵 Minor | **Style** | Naming conventions, formatting, comment quality |

3. **Produce structured findings** using this format for each issue:

   ```markdown
   ### [SEVERITY] Short title

   **File:** `path/to/file.ext` (line X-Y)
   **Category:** Security | Performance | Best Practices | Style
   **Description:** What the issue is and why it matters.
   **Suggestion:**
   ```diff
   - current problematic code
   + suggested fix
   ```
   ```

4. **Summarise** at the end with:
   - Total issues by severity
   - Top 3 most impactful improvements
   - Overall assessment (pass / pass with warnings / needs attention)

### If no issues are found

Report a clean bill of health — do not invent issues.

## Rules

- Be constructive, not nitpicky.
- Prioritise issues that affect correctness and security over style.
- Include line numbers and file paths for every finding.
- Suggest fixes, not just problems.
- Never modify files directly during a review — only report findings.
