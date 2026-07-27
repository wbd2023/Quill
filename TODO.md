# Quill maintenance queue

This queue holds small, repository-local tasks that do not alter the ideal MVP
scope. The [MVP roadmap](docs/roadmap.md) owns committed release work.

## Open

- [ ] Audit authored prose for Australian English and replace inconsistent US
  spellings without changing technical identifiers.

## Rules

- Keep each entry independently actionable and small enough for one focused
  change.
- State the expected outcome, not only the symptom.
- Move work that changes MVP scope, acceptance criteria, or phase dependencies
  to [docs/roadmap.md](docs/roadmap.md).
- Move a substantial design decision to an ADR before implementation.
- Remove completed entries rather than retaining a historical changelog.
