# LLDZen

Low-level design notes and Go reference implementations.

## Site

```bash
npm install
npm run docs:dev
```

Open **`http://localhost:5173/lld/`** (port may differ). If you see a blank page after upgrading deps, clear the Vite cache: `rm -rf docs/.vitepress/cache` and restart.

`sync-problem-docs.mjs` copies `problems/*` extensions into `docs/problems/` before each dev/build and writes `docs/.vitepress/sidebar-problems.json`.

### GitHub Pages

The live site must deploy the **VitePress build**, not Jekyll from the `docs/` folder.

1. Repo **Settings → Pages → Build and deployment → Source:** choose **GitHub Actions** (not “Deploy from a branch”).
2. Push to `main` / `master` (workflow: [`.github/workflows/pages.yml`](.github/workflows/pages.yml)).
3. Open **Actions** and confirm “Deploy site to GitHub Pages” succeeded.

If you see only a heading like “lld” at [gauxs.github.io/lld](https://gauxs.github.io/lld/), Pages is still using Jekyll on `/docs`—switch the source to GitHub Actions and re-run the workflow.

## Problem layout

Each problem under `problems/<problem_name>/`:

| Path | Purpose |
| --- | --- |
| `extensions.json` | `order`, `default`, per-extension `title` and optional `buildsOn` |
| `extensions/<id>/requirements.md` | FR/NFR for this extension only + **Extensions from here** |
| `extensions/<id>/design.md` | Entity/API tables (delta from parent when not baseline) |
| `extensions/<id>/code/` | Optional reference source (verbatim on **Codebase** page) |

Template: [`_template/problem/`](_template/problem/) · Conventions: [`_template/CONVENTIONS.md`](_template/CONVENTIONS.md)

## Packages

| Problem | Code |
| --- | --- |
| Connect Four (baseline) | `problems/connect_four/extensions/baseline/code` (`go run .` from repo root) |
| Rate limiter | `rate_limiter/` (WIP) |
