---
name: docs-writer
description: "Create or update documentation with consistent structure and validated links."
---

# Documentation Writer

When to use

When creating or updating READMEs, guides, API docs, or any markdown.

Create docs

1. Identify audience: end user | developer | operator.
2. Use this structure:
```markdown
# Title

> One-line summary.

## Overview

## Prerequisites

## Steps / Sections

## Examples

## Troubleshooting

## References
```
3. Add TOC if document has more than three top-level sections.
4. Use fenced code blocks with language tags.

Update docs

- Preserve tone and internal comments (`<!-- internal -->`).
- Add a short changelog note describing edits.
- Validate links and cross-references; if uncertain, add `<!-- TODO: verify -->`.

Rules

- Use active voice and concise sentences.
- One sentence per line in source markdown.
- Do not remove content without an explicit note.
