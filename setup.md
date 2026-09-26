# Setup and authoring

How the LLDZen site is built, how theory and problem pages are generated, and how to add content.

## How pages are created

**Theory** and **problems** are authored at the repo root and synced into `docs/` before VitePress runs.

```text
theory/<topic>/sections.json + *.md
        │
        ▼  node scripts/sync-theory-docs.mjs
docs/learn/<topic>/                       (pages with front matter)
docs/.vitepress/sidebar-theory.json

problems/<name>/extensions.json + extensions/…
        │
        ▼  node scripts/sync-problem-docs.mjs
docs/problems/<slug>/ …
docs/.vitepress/sidebar-problems.json
        │
        ▼  vitepress dev docs  |  vitepress build docs
Static site (base path /lld/)
```

| Command | What it does |
| --- | --- |
| `npm install` | Install Node dependencies (once). |
| `npm run docs:dev` | Sync theory + problems, then VitePress dev server. |
| `npm run docs:build` | Sync, production build, add `.nojekyll` for GitHub Pages. |
| `npm run docs:preview` | Serve the production build locally. |

**Do not hand-edit** synced pages under `docs/learn/` or `docs/problems/` except when debugging a generator—they are overwritten on the next sync.

### Theory sync ([`scripts/sync-theory-docs.mjs`](scripts/sync-theory-docs.mjs))

1. Discovers subdirectories of **`theory/`** that contain **`sections.json`** (skip dirs starting with `_`).
2. For each page id in **`order`**, reads **`theory/<topic>/<id>.md`** and writes **`docs/learn/<topic>/<id>.md`** with YAML (`title`, `description`, `prev` / `next` from order or per-page overrides).
3. Removes stale `.md` files in the destination topic folder.
4. Writes **`docs/.vitepress/sidebar-theory.json`** (imported in [`docs/.vitepress/config.ts`](docs/.vitepress/config.ts)).

### Problem sync ([`scripts/sync-problem-docs.mjs`](scripts/sync-problem-docs.mjs))

1. Reads each problem listed in **`SLUGS`** (maps directory name → URL slug, e.g. `connect_four` → `connect-four`).
2. Loads **`extensions.json`**: `order`, `default`, and per-extension `title` and optional **`buildsOn`**.
3. Writes **`docs/problems/<slug>/index.md`**: intro plus a **mermaid** flowchart (child `-->|BuildsOn|` parent).
4. For each extension, copies **`requirements.md`** and **`design.md`** into the docs tree with front matter (`prev` / `next` trail, `pageClass` for requirement styling).
5. If **`extensions/<id>/code/`** is non-empty, generates **`codebase.md`** with `<ProblemCodebase problem="…" extension="…" />` (reads files from disk at build/dev time).
6. Writes **`docs/.vitepress/sidebar-problems.json`** for VitePress to import in [`docs/.vitepress/config.ts`](docs/.vitepress/config.ts).
7. Removes legacy flat problem pages (`functional-requirement.md`, old `variation-*` folders, etc.).

Mermaid diagrams use **`vitepress-plugin-mermaid`** in the VitePress config; fenced ` ```mermaid ` blocks in synced or hand-written markdown render on the site.

---

## Problem directory layout

Each problem is modular: one folder per **extension** (sibling scopes), not one monolithic requirements doc.

```text
problems/my_problem/
├── extensions.json
└── extensions/
    ├── baseline/
    │   ├── requirements.md
    │   ├── design.md
    │   └── code/              optional reference Go (or other) sources
    ├── networked/
    │   ├── requirements.md
    │   └── design.md
    └── ...
```

| File | Purpose |
| --- | --- |
| `extensions.json` | Sidebar **`order`** (put **`baseline`** first), **`default`** extension id, **`extensions`** map with **`title`** and optional **`buildsOn`** (parent extension id). |
| `extensions/<id>/requirements.md` | Interview prompt and **FR-*** / **NFR-*** for **this extension only**; **`## Extensions from here`** table linking child extensions. |
| `extensions/<id>/design.md` | Entities and API tables; for non-baseline extensions, document **delta** from the parent (link at top: **Builds on:**). |
| `extensions/<id>/code/` | Optional; when present, sync adds a **Codebase** trail step and sidebar link. |

Extensions appear as **siblings** in the sidebar; dependency is expressed with **`buildsOn`** in JSON and **Builds on** links in markdown, not nested folders.

Copy [`_template/problem/`](_template/problem/) to `problems/<your_problem>/` and rename placeholders. Naming rules: [`_template/CONVENTIONS.md`](_template/CONVENTIONS.md).

### `extensions.json` example

```json
{
  "default": "baseline",
  "order": ["baseline", "networked"],
  "extensions": {
    "baseline": { "title": "Baseline" },
    "networked": {
      "title": "Networked multiplayer",
      "buildsOn": "baseline"
    }
  }
}
```

Reference: [`problems/connect_four/extensions.json`](problems/connect_four/extensions.json).

### Requirement gathering (per extension)

When drafting `requirements.md`:

1. **Functional:** What must this extension do?
2. **Scope:** What is out of scope here; what belongs in a **child extension** later (YAGNI, but design for extension).
3. **Non-functional:** Assumptions, constraints, scale, failure handling.
4. **Errors:** What can fail and how the API behaves.

Use the HTML patterns in [`_template/problem/extensions/baseline/requirements.md`](_template/problem/extensions/baseline/requirements.md) (`lld-req`, interview reveal blocks) so pages match the rest of the site.

---

## Add a new problem (checklist)

1. **Scaffold** — Copy `_template/problem/` → `problems/<problem_id>/` (snake_case directory name).
2. **Author** — Fill baseline `requirements.md` / `design.md`; add `code/` if you want a codebase page.
3. **Register slug** — Add an entry to **`SLUGS`** in `scripts/sync-problem-docs.mjs`:

   ```js
   const SLUGS = {
     connect_four: "connect-four",
     rate_limiter: "rate-limiter",
   };
   ```

4. **Sidebar** — Extend [`docs/.vitepress/config.ts`](docs/.vitepress/config.ts):
   - Import `sidebar-problems.json` (keys are camelCase problem ids: `connectFour`, `rateLimiter`, …).
   - Add a Problems entry with **Overview** → `/problems/<slug>/` plus the matching array from that JSON.

5. **Run sync** — `npm run docs:dev` and open `/problems/<slug>/`.

6. **Go module** — If you add Go code, use import paths under `github.com/gauxs/lld/problems/...` and wire `main.go` or a small demo as needed.

---

## Add an extension to an existing problem

1. Create `problems/<problem>/extensions/<new_id>/` with `requirements.md` and `design.md`.
2. Add **`buildsOn`** in `extensions.json` if it depends on another extension.
3. Append **`<new_id>`** to **`order`** (siblings: order is display only; baseline stays first).
4. Update the parent extension’s **`## Extensions from here`** table in `requirements.md`.
5. Re-run sync (happens automatically with `docs:dev` / `docs:build`).

---

## GitHub Pages

Deploy the **VitePress build**, not Jekyll from `/docs`.

1. Repo **Settings → Pages → Build and deployment → Source:** **GitHub Actions**.
2. Push to `main` / `master` ([`.github/workflows/pages.yml`](.github/workflows/pages.yml)).
3. Confirm the “Deploy site to GitHub Pages” workflow succeeded.

Live URL uses base path **`/lld/`**. If you only see a bare “lld” heading, Pages is still serving Jekyll—switch to GitHub Actions and redeploy.
