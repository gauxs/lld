# LLDZen

Personal **low-level design** notes and Go reference code for system-design and OOD interviews. Each problem is broken into **extensions** (baseline first, then optional follow-ups) with requirements, design, and optional implementation.

The browsable site is built with VitePress from this repo. To run it locally or add a new problem, see **[setup.md](setup.md)**.

```bash
npm install
npm run docs:dev
```

Open **`http://localhost:5173/lld/`** (port may differ). If the page is blank after dependency changes, clear the cache: `rm -rf docs/.vitepress/cache` and restart.

## Reference implementations

| Problem | Go package |
| --- | --- |
| Connect Four (baseline) | `problems/connect_four/extensions/baseline/code` — `go run .` from repo root |
| Rate limiter (baseline) | `problems/rate_limiter/extensions/baseline/code` |

## Authoring

- Problem template: [`_template/problem/`](_template/problem/)
- Conventions: [`_template/CONVENTIONS.md`](_template/CONVENTIONS.md)
- Full workflow: [setup.md](setup.md)
