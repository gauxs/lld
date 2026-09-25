# LLD problem conventions

Each problem lives under `problems/<problem_name>/` with four authored sections plus code.

| File | Content |
| --- | --- |
| `functional_requirement.md` | `FR-*` items; interview questions inside `<details><summary>Questions to ask</summary>`; NFR table; out-of-scope in collapsible block |
| `design.md` | Entity ↔ responsibilities table and API tables only—no code dumps |
| `code/` | Reference implementation (language-specific package) |
| `followup.md` | `FU-*` follow-up requirements with the same collapsible Q&A pattern |

Site pages under `docs/problems/` mirror these trails for VitePress. The **codebase** page renders files from `problems/<name>/code/` verbatim.
