# reconcile

Reconcile a repository or a specified part by resolving unnecessary duplication and unintended divergence in its content.

Use this skill when related documents, instructions, configuration, schemas,
examples, or code repeat the same knowledge or have drifted apart. The concise
[skill](src/SKILL.md) defines the desired outcome and leaves methods to the
caller and model.

## Usage

Invoke this skill directly during maintenance. Automatic model invocation is
disabled through the supported native settings.

Name the scope and any known authority or differences that must remain:

An opportunity question requests findings and recommendations without edits:

```text
What opportunities do we have to $reconcile our API examples and schemas?
```

A direct change request asks for implementation and verification:

```text
$reconcile the setup instructions across README.md, docs/, and example projects
$reconcile our API schema, examples, and implementation against the supported contract
$reconcile repeated repository guidance while preserving package-specific requirements
```

## Philosophy

Reconciliation makes related representations agree on their shared meaning
while reducing unnecessary maintenance. Duplication can be burdensome even
when every copy currently agrees. Divergence can be a defect even when each
representation has a useful reason to exist independently.

Determine whether the content expresses the same underlying fact, rule, or
responsibility. Similar wording or code may serve different purposes, versions,
audiences, or domains. A tutorial example, a reference description, and an
implementation can remain different while agreeing on the contract they share.

Establish intended meaning from applicable authority and evidence. A newer
document can be mistaken; repeated claims can share the same error. Existing
behavior is evidence, but can itself violate a requirement. When the evidence
does not resolve a conflict, expose that uncertainty rather than manufacture
agreement or silently discard an obligation.

Consolidation is useful when it reduces the total burden of keeping content
correct. References, shared definitions, generated representations, and
consistency checks are possible techniques. Choose them according to the
content and its consumers; their coupling and upkeep must justify their cost.
Useful repetition and independent representations can remain.

The result should retain required information, behavior, and supported uses.
Evidence of coherence depends on the scope: it may involve checking claims,
links and discovery paths, generated output, or observable behavior. The skill
prescribes no fixed workflow, tool, reporting template, or automatic execution
of other skills.

## Maintenance companion

`reconcile` remains independently installable after software-maintenance 0.2.0.
It is no longer a member of that pack. Install it directly when this specific
outcome is needed:

```bash
axm install @craigsmitham/skills/reconcile
```

## License

MIT. The instructions and explanatory prose are original; no third-party text
or code is vendored in this package.
