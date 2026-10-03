# Add study content

Author only in `theory/` and `problems/`. Files under `docs/learn/` and
`docs/problems/` are generated and will be overwritten.

## Run locally

```bash
npm install
npm run docs:dev
```

Open `http://localhost:5173/lld/`. New and edited source files are synced
automatically while the server runs.

## Add theory

Create a Markdown file under a topic:

```text
theory/<topic>/01-introduction.md
theory/<topic>/02-next-concept.md
```

The first `# Heading` becomes the page and sidebar title. Numeric filename
prefixes control order. Subdirectories create nested sidebar groups.

```text
theory/concurrency/03-patterns/01-worker-pool.md
→ /learn/concurrency/03-patterns/01-worker-pool
```

Optional front matter can override generated metadata:

```yaml
---
title: Worker pools
sidebar: Worker pool
description: Bounded concurrent work in Go
---
```

No registration or sidebar edit is required.

## Add a problem

Copy the problem template:

```bash
mkdir -p problems/concurrency
cp -R _template/problem problems/concurrency/meeting_room
```

Choose any path below `problems/`, then:

1. Edit `extensions.json`.
2. Write `extensions/baseline/requirements.md`.
3. Write `extensions/baseline/design.md`.
4. Optionally add source files under `extensions/baseline/code/`.

Directory names use `snake_case`; URL segments use `kebab-case`:

```text
problems/concurrency/meeting_room/
→ /problems/concurrency/meeting-room/
```

The directory hierarchy, overview, sidebar, and extension pages are generated
automatically.

### Writing the two phases

Keep `requirements.md` solution-free:

- Functional requirements numbered `FR-*`
- Out-of-scope behavior
- Non-functional requirements numbered `NFR-*`
- Child extensions as links
- Existing interviewer questions in an **Interview prompts** blockquote; omit
  the block when no prompts exist

Use `design.md` for the solution:

- Start with a brief approach
- Use one `go` design sketch for structs, important fields, interfaces, and
  public method signatures
- Put ownership, responsibilities, and invariants in comments above types and
  methods
- Omit method bodies; keep implementation on the Codebase page
- Follow with main flows, concurrency, complexity, and trade-offs when present

Headings, bullets, blockquotes, fenced code, and Mermaid all render directly.
Avoid tables and custom HTML.

### Adding code

Place text-based source files anywhere under an extension's `code/` directory.
The generated Codebase page preserves the directory tree and selects syntax
highlighting from each filename. Binary and potentially sensitive dotfiles are
skipped; `.gitignore` and `.env.example` are included.

## Add a problem extension

1. Create `extensions/<extension_id>/requirements.md` and `design.md`.
2. Add the extension to `extensions.json` and its `order`.
3. Set `buildsOn` when it extends another extension.
4. Optionally add an `extensions/<extension_id>/code/` directory.

## Add an AI practice guide

Add Markdown under `practice_with_ai/`. Files are published under
`/practice-with-ai/` and added to the third sidebar section automatically.

## Verify

```bash
npm run docs:build
```
