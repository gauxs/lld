# LLD problem conventions

Each problem lives under `problems/<problem_name>/` with **extensions** (modular scopes), not one global requirements list.

| Path | Content |
| --- | --- |
| `extensions.json` | `order` (sidebar sequence, baseline first), `extensions` → `{ title, buildsOn? }`; `default` id |
| `extensions/<id>/requirements.md` | Problem slice, `FR-*` / `NFR-*` for **this extension only**; **Extensions from here** table at the end |
| `extensions/<id>/design.md` | Entity ↔ responsibilities and API tables; **delta from parent** when not root |
| `extensions/<id>/code/` | Optional reference implementation |

Site pages under `docs/problems/` are synced from `problems/` for VitePress. The **codebase** page renders files from the extension’s `code/` directory via `<ProblemCodebase problem="…" extension="…" />`.

## Minimal single-extension problem

One folder `extensions/baseline/` with requirements, design, and code is enough. `extensions.json` lists a single node with empty `children`.

## Naming

Use **extension** (not variation). Siblings in the sidebar; dependency is **`buildsOn`** in `extensions.json` plus **Builds on** / mermaid graph on the problem overview. Document the link at the top of requirements/design and in the parent’s **Extensions from here** section.
