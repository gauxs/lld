# LLDZen

Low-level design theory, interview problems, and reference implementations.
Problems progress through a baseline and optional extensions, each with
requirements, design, and optional code.

## Run locally

```bash
npm install
npm run docs:dev
```

Open `http://localhost:5173/lld/`.

## Repository

- `theory/` — source Markdown for theory trails
- `problems/` — problem requirements, designs, extensions, and code
- `_template/problem/` — starting point for a new problem
- `docs/` — generated VitePress site and theme

Theory and problems are discovered automatically from their directory
structure. Do not edit generated pages under `docs/learn/` or
`docs/problems/`.

See **[Add study content](_template/setup.md)** for the authoring workflow.

## Build

```bash
npm run docs:build
npm run docs:preview
```
