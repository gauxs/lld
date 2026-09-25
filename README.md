# LLDZen

Low-level design notes and Go reference implementations.

## Site

```bash
npm install
npm run docs:dev
```

`sync-problem-docs.mjs` copies `problems/*/`.md into `docs/problems/` before each dev/build.

## Problem layout

Each problem under `problems/<problem_name>/`:

| Path | Purpose |
| --- | --- |
| `functional_requirement.md` | Problem statement, FR-* / NFR-* / out-of-scope with collapsible interview Q&A |
| `design.md` | Entity ↔ responsibilities and API tables |
| `code/` | Reference source (verbatim on the site **Codebase** page) |
| `followup.md` | Follow-up requirements (FU-*) |

Template: [`_template/problem/`](_template/problem/) · Conventions: [`_template/CONVENTIONS.md`](_template/CONVENTIONS.md)

## Packages

| Problem | Code |
| --- | --- |
| Connect Four | `problems/connect_four/code` (`go run .` from repo root) |
| Rate limiter | `rate_limiter/` (WIP) |
