---
name: prepare-pr
description: "Create a clean pull request from agent-made changes."
---

# Prepare Pull Request

## When to use

Use this skill after making file changes that should be submitted as a pull
request for review, rather than pushed directly to the default branch.

## Instructions

1. **Create a branch** from the current HEAD:
   ```bash
   BRANCH="agent/$(date +%Y%m%d-%H%M%S)-${MISSION_SLUG}"
   git checkout -b "${BRANCH}"
   ```

2. **Stage changes** selectively — only include files that were intentionally
   modified by the mission. Do NOT stage unrelated files:
   ```bash
   git add <changed-files>
   ```

3. **Commit** using Conventional Commits format:
   ```
   <type>(<scope>): <short summary>

   <body — what changed and why>
   ```
   Types: `fix`, `feat`, `docs`, `refactor`, `chore`, `ci`, `test`.

4. **Push** the branch:
   ```bash
   git push origin "${BRANCH}"
   ```

5. **Open a Pull Request** using `gh pr create`:
   ```bash
   gh pr create \
     --title "<type>(<scope>): <summary>" \
     --body "## Summary
   <what this PR does>

   ## Changes
   - <file 1>: <what changed>
   - <file 2>: <what changed>

   ## Agent Mission
   > <original mission prompt>

   ---
   _This PR was created by an automated agent._" \
     --base main
   ```

## Rules

- Never force-push.
- Never commit secrets, tokens, or credentials.
- Keep commits atomic — one logical change per commit.
- If there are no changes to commit, exit gracefully with a message.
