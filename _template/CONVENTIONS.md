# LLD problem conventions

Each problem lives under `problems/<problem_name>/` with four authored sections plus code.

| File | Content |
| --- | --- |
| `functional_requirement.md` | Problem statement, `FR-*` and `NFR-*` blocks with `lld-reveal` interview prompts; **Out of scope** section (same pattern) |
| `design.md` | Entity ↔ responsibilities table and API tables only—no code dumps |
| `code/` | Reference implementation (language-specific package) |
| `followup.md` | `FU-*` follow-up requirements with the same collapsible Q&A pattern |

Site pages under `docs/problems/` are synced from `problems/` for VitePress. The **codebase** page renders files from `problems/<name>/code/` verbatim.
