# Quill Repository Policy

Quill makes a repository's written engineering policy traceable to reproducible Checks and safe
Fixes. This context covers policy authorship, reusable enforcement capabilities, repository
selection, findings, and declared automation coverage.

## Policy

**Repository**:
A local source tree whose engineering policy Quill evaluates.
_Avoid_: Workspace, project when referring to the evaluation root

**Repository Policy**:
The Repository's normative engineering expectations and its declared enforcement approach. It is
expressed through a Style Guide and a Repository Profile; neither file alone is the whole policy.
_Avoid_: Config, Profile, Pack Policy when referring to the whole

**Style Guide**:
The human-readable document containing identified Requirements and explanatory policy prose,
normally `STYLE.md`.
_Avoid_: Profile, config, rules file

**Requirement**:
One stable-ID normative statement in the Style Guide. It is written policy and may be automated,
review-only, or deferred; it is not executable work.
_Avoid_: Rule, Check, lint

**Repository Profile**:
The machine-readable Repository declaration, normally `quill.toml`, that describes topology and
reproducibility inputs, selects Packs, and activates Rules with their Requirements, Enforcement, and
Scope.
_Avoid_: Policy, manifest, config

**Pack Parameters**:
Repository-supplied, Pack-defined values that specialise a selected Pack's Rules. The Pack owns
their meaning and validation.
_Avoid_: Pack Policy, Rule Config, Check Config

**Resolved Policy**:
The fully validated interpretation of one Style Guide, Repository Profile, and available Pack
catalogue, with every reference resolved. It contains Active Rules but is not an operation schedule.
_Avoid_: Plan, compiled config, Effective Profile

## Enforcement capabilities

**Pack**:
A reusable, identified provider of Rules and their supporting definitions, which may include Tool
dependencies, File Set defaults, and Pack Parameters. It is an ownership and distribution concept,
not an executable plugin type.
_Avoid_: Bundle, registry, plugin, manifest

**Shipped Pack**:
A Pack distributed with Quill. Shipped describes provenance only; all Packs expose the same Rule
semantics.
_Avoid_: Built-in package, builtin, bundled checks

**External Pack**:
A Pack supplied by the Repository rather than distributed with Quill. External describes provenance
only, not a required execution or trust mechanism.
_Avoid_: Shipped Pack, Tool

**Rule**:
A stable, Pack-owned enforcement capability with exactly one Check and optionally one Fix. A Rule
becomes active only through an explicit Repository Profile binding.
_Avoid_: Requirement, Check, detector, lint

**Active Rule**:
A Rule activated for the Repository by binding it to one or more Requirements, one Enforcement, and
one Scope. Selecting a Pack makes its Rules available but does not activate them.
_Avoid_: Enabled Pack, Check

**Enforcement**:
The consequence assigned to findings from an Active Rule: blocking (`required`) or reporting-only
(`recommendation`), subject to an explicitly requested strict mode.
_Avoid_: Requirement strength, Check status

**Check**:
The mandatory read-only assessment performed by a Rule.
_Avoid_: Rule, validation, checker identity

**Fix**:
The optional state-changing behaviour of a Rule that applies only changes the Rule declares safe. A
Fix does not weaken the Rule or remove the need to Check the resulting Repository.
_Avoid_: Check, formatter, arbitrary rewrite

**Tool**:
An identified external executable dependency whose command, version detection, installation method,
and Repository pin Quill can manage. An External Pack runtime is not a Tool merely because both are
executables.
_Avoid_: Pack, Rule, arbitrary binary, runtime

**Tool Pin**:
The Repository Profile's required Tool version and execution limits for one Tool.
_Avoid_: Tool, installed Tool, lock entry

**Lockfile**:
A Repository-owned record of resolved integrity data used to reproduce Tool installation, normally
`quill.lock`. It supplements Tool Pins rather than replacing them.
_Avoid_: Profile, dependency manifest

## Repository selection

**Scope**:
A named Repository region defined by one or more roots. Scope determines Rule applicability and
bounds File Set evaluation and Target selection; Scopes may overlap.
_Avoid_: Target, File Set, workspace

**File Set**:
A named, reusable include/exclude definition for classes of Repository files. It is evaluated within
a Scope to produce Selected Files; it is not the concrete list.
_Avoid_: Scope, Target, files

**Selected Files**:
The concrete Repository-relative file list obtained by evaluating a File Set, or the whole eligible
Scope when no File Set is declared, for one Check or Fix.
_Avoid_: File Set, Scope

**Target**:
A named language or ecosystem invocation context associated with a Scope, including its working
directory and operation-specific paths. It says where language-aware work runs, not which files or
Repository region apply.
_Avoid_: Scope, File Set, backend, workspace

## Outcomes and traceability

**Diagnostic**:
One structured policy finding from a completed Check, with a message and optional code, file,
source range, and help link. It is not a Rule Outcome or an operational failure.
_Avoid_: Error, Result, log message

**Rule Outcome**:
The classified disposition of one Active Rule considered by a Check operation: pass, warning,
failure, blocked, or execution error. It carries either completed Check evidence and Diagnostics or
the corresponding Execution Failure; a Rule filtered out before consideration has no outcome.
_Avoid_: Result, Diagnostic, command output

**Fix Outcome**:
The record of attempting one Active Rule's Fix, containing completed execution evidence or the
corresponding Execution Failure. It does not assert policy compliance; a later Check determines the
Rule Outcome.
_Avoid_: Result, Rule Outcome, Diagnostic

**Execution Failure**:
An operational inability to complete a Rule's Check or Fix, such as an unavailable Tool, timeout,
malformed Pack response, or process failure. It is not evidence that Repository Policy was violated.
_Avoid_: Diagnostic, Rule failure, violation

**Automation Coverage**:
An account of how each Requirement relates to declared machine enforcement or explicit human
treatment. The strength required for an automation claim is a product decision; coverage never
denotes runtime success or compliance.
_Avoid_: Coverage, code coverage, pass rate, compliance
