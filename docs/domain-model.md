# Quill domain model

## Status

This document records the first clean-slate domain model for the Rust rewrite. It is grounded in the
`v0.2.0` Go implementation and retained product/protocol evidence, but it does not prescribe Rust
crates, modules, traits, storage types, or wire formats.

Canonical terms and avoided synonyms live in [CONTEXT.md](../CONTEXT.md).

## Core relationship

Quill exists to make written Repository Policy traceable to executable enforcement without
pretending that automation and human policy are identical.

```text
Repository Policy
|-- expressed through a Style Guide
`-- configured through Repository Profiles (cardinality unresolved)

Resolved Policy
|-- interprets exactly one Style Guide
`-- uses exactly one selected Repository Profile
    |-- selects zero or more available Packs
    |-- supplies Pack Parameters
    |-- defines Scopes, File Sets, Targets, and Tool Pins
    `-- activates Rules through explicit bindings

Pack --owns--> one or more Rules
Pack --may depend on--> zero or more Tools
Rule --has--> exactly one Check
Rule --may have--> zero or one Fix
Active Rule --binds--> one or more Requirements
Active Rule --has--> exactly one Enforcement and one Scope
Requirement <--related to--> zero or more Active Rules

Scope + File Set --evaluation--> Selected Files
Scope + language --selection--> zero or more Targets
Active Rule considered by a Check operation --classified as--> one Rule Outcome
Active Rule selected by a Fix operation --recorded as--> one Fix Outcome
Rule Outcome --carries--> completed Check evidence + Diagnostics OR an Execution Failure
Fix Outcome --carries--> completed Fix evidence OR an Execution Failure
Requirements + declared enforcement relationships --project to--> Automation Coverage
Tool + Tool Pin + optional Lockfile integrity --support--> reproducible execution
```

The Requirement-to-Active-Rule relationship is many-to-many. Neither side can absorb the other
without losing behavior.

## Invariants

1. Requirement identity belongs to the Style Guide. Diagnostic codes belong to their producer. They
   are never interchangeable.
2. Every Active Rule names an available Rule from a selected Pack, belongs to exactly one Pack, and
   binds at least one existing Requirement.
3. Every Rule has one mandatory Check and at most one Fix. A Fix is never independently active.
4. Pack selection and Rule activation are different states.
5. Scope, File Set, and Target remain different concepts: topology, file classification, and
   invocation context.
6. A Rule Outcome carries either completed Check evidence and any Diagnostics, or the corresponding
   Execution Failure. Neither is silently converted into the other.
7. A Fix Outcome carries completed Fix evidence or the corresponding Execution Failure. It does not
   assert compliance; only a subsequent Check determines a Rule Outcome.
8. Repository-level Diagnostics have no file or range. File Diagnostics use clean
   Repository-relative slash paths and valid one-based UTF-8-byte, half-open ranges.
9. Automation Coverage accounts for every Requirement and distinguishes declared machine
   enforcement from human treatment. It never asserts runtime health, successful execution, or
   compliance.
10. Untrusted Repository inputs and Repository-supplied Pack data are fully validated before
    entering the Resolved Policy or trusted Diagnostic model.
11. Executing Repository-supplied Pack code requires an explicit trust and resource contract.
    Containment and resource controls must never be described as sandboxing.
12. Tool installation fails closed when contents or integrity cannot be verified before
    replacement.
13. Unqualified `Result`, `Policy`, and `Coverage` are not accepted domain terms. Name the specific
    concept.

## Scenarios that distinguish terms

### One Requirement, multiple Rules

"Imports are grouped" may be checked by both a formatter Rule and a semantic import Rule. This is
one Requirement with two automation links, not two Requirements and not one merged Check.

### One Rule, multiple Requirements

A broad shell-safety Rule may enforce several independently identified Requirements. The Rule is
the execution and reporting unit; the Requirements remain separately traceable.

### Review-only Requirement

A security-architecture judgment may have no sound automation. It remains a Requirement and appears
in Automation Coverage without inventing a no-op Rule.

### Selected Pack, inactive Rule

A Pack and all its available Rules may be listed while only explicitly bound Rules are active.
Calling every Rule in a selected Pack enabled would erase this distinction.

### Scope versus File Set

An application Scope may contain `cmd`, `src`, and `test`, while a documentation File Set classifies
Markdown files. The File Set may include root `README.md` but exclude an otherwise matching file in
a disjoint tools Scope. Scope determines the region; File Set classifies files within it.

### Default Scope is not a hierarchy

A default Scope is only the invocation fallback. It is not necessarily the widest Scope or the
parent of other Scopes.

### Target versus Scope

One Repository-wide language Rule may select two Targets with different working directories. There
is one applicability region but two invocation contexts.

### Empty Selected Files

A valid Rule and File Set may produce no Selected Files in a disjoint Scope. That is empty input and
normally a clean Check, not an unknown File Set or an Execution Failure.

### Diagnostic versus Execution Failure

A Rule may complete its Check and emit three valid Diagnostics. If Quill cannot complete the Check
because an external capability times out, exits unsuccessfully, exceeds bounds, or returns invalid
data, the same Rule Outcome carries an Execution Failure instead.

### Repository-level Diagnostic

"Required Pack parameter missing" may have no source file. Inventing line 1 would create a false
location.

### Fix Outcome is not compliance

A Fix may complete and report execution evidence, but that does not establish that the Repository
now satisfies the Rule. A subsequent Check produces the Rule Outcome; a failed Fix instead records
an Execution Failure in its Fix Outcome.

### Tool versus External Pack execution

`markdownlint` may be a Tool with a version pin and installation lifecycle. Code supplied by an
External Pack belongs to that Pack's provenance and trust boundary, regardless of which execution
mechanism the Rust product eventually chooses.

### Automation Coverage versus compliance

Final Go called a Requirement automated when it had an Active Rule binding, even if the Tool was
unavailable and the Rule would be blocked. The Rust product may require a stronger reviewed
automation claim, but Automation Coverage cannot mean runtime success, pass rate, or compliance.

## Rejected collapses

### Requirement, Rule, and Check as one term

This cannot describe review-only Requirements, optional Fixes, many-to-many traceability, or the
difference between a selectable Rule and its read-only behavior.

### Scope, File Set, and Target as one selector

This combines Repository topology, file classification, and language invocation. It cannot explain
overlapping Scopes, empty file selection, or multiple Targets for one scoped Rule.

### Repository Policy as only `quill.toml`

This makes human Requirements subordinate to machine configuration. Quill's distinguishing product
idea is traceability between the Style Guide and executable enforcement.

### Tool and External Pack execution as one executable concept

Tools have pinning and installation semantics. External Pack execution has Repository-supplied
provenance and a separate trust contract. A future transport choice does not erase that distinction.

### Plan as ubiquitous language

The final Go `Plan` was an unordered collection of resolved Rules, while operation and Scope
selection happened later. `Resolved Policy` names the actual domain state without implying a
schedule.

### Driver, Template, Job, handler, and registry as domain terms

These names describe possible implementations. None is required to explain Repository Policy to a
user or state a product invariant.

## Unresolved product questions

These questions must be resolved before their answers appear in schemas or architecture.

1. **Compatibility adoption:** Were the version-1 CLI JSON and `quill-pack-v1` protocols consumed
   outside this repository? Publication is established by `v0.2.0`; adoption is not.
2. **Rule identity:** Must every Rule identity be Pack-qualified, such as `<pack>/<rule>`? Final Go
   required global uniqueness but did not require ownership in the spelling.
3. **Automation claim:** Does automated continue to mean an Active Rule binding, or must it carry a
   reviewed claim that the Rule substantially enforces the Requirement?
4. **External Fixes:** Will External Packs remain Check-only, or will a future protocol permit safe
   Fixes?
5. **Profile cardinality:** Is `quill.toml` permanently the one Repository Profile, or is there a
   real use case for selecting among named Profiles?

## Evidence boundary

The final Go implementation at `v0.2.0` is evidence for difficult cases and established behavior.
The Rust rewrite may intentionally change behavior, but it must name the product decision rather
than losing the case accidentally.

Versioned wire spellings such as `policy`, `check_id`, and nested `result` fields do not redefine the
domain language. If compatibility is retained, translate those names at the boundary.
