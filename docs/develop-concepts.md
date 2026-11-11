# Develop concepts

Use this guide when you need to explain a coherent piece of product behavior,
reuse that behavior in another setting, or describe how independent concepts work
together. Start with purpose and a characteristic story; add structure where it
helps people make a design decision. The [format reference](../README.md#concept-spec-format)
defines the basic file contract.

ConceptUse and Synchronization are specifications in the same abstract hierarchy
as Concept, even though they are authored inline in a System document. They can
carry shared Notes, Evidence, Related specifications and Validation plan sections,
and optional introductory `**ID:**` values. This classification does not imply a
separate file or an execution guarantee. See [specification kinds](specification-model.md#specification-kinds-and-document-placement).

## Find a purposeful boundary

A concept is a unit of behavior people can understand and use: Shorthand, Trash,
Reservation or Authentication. It can involve several kinds of objects. A domain
noun, class, UI widget, endpoint or delivery story does not by itself define a
concept. For an API, programmers are users too.

Ask what benefit the concept provides on its own and what history of interactions
shows that benefit. Registering credentials alone does not show authentication;
registration followed by a successful or unsuccessful attempt does. If explaining
the benefit requires several independent purposes, consider separate concepts and
a composition rule.

Create a draft with `dotspec concept init shorthand`. Purpose, State, Actions and
Operational principles are conventional sections; unfinished descriptions remain
valid. Model export preserves every section, including unknown details.

## Write the operational principles

Tell one or a few short stories explaining why the design works. Use actions and
observable state, without prescribing buttons or screen layouts. Include enough
history to reveal the benefit, and any conditions on which the story depends.

Use `## Operational principles` once per document. Named subsections distinguish stories without
requiring scenario records. For Trash, recovery and permanent removal can be two
stories. The stories need not all be successful actor-goal tasks.

An operational principle explains the essential design. State and action behavior
answer the remaining questions. Use cases can describe broader interactions across
concepts; they do not become the concept's operational principle automatically.

## Describe state and actions

Include information the concept must remember to support its actions. Sets and
relations often express that information more clearly than implementation classes.
Do not infer a database layout, UI representation or service boundary from the
abstract state.

Under `## State`, optional `### Initial state` and `### Invariants` describe the
starting state and constraints. Other H3 headings identify components. A component
can have `#### Type` (also Value or Values) and `#### Derivation` subsections.
Derived state need not be stored separately.

Use H3 headings under `## Actions` for individual actions. Inputs, Outputs,
Effects, Rejections, Caller assumptions, Guards and Initiation are optional H4
subsections. Caller assumptions and enforced guards make different promises;
rejection behavior states what happens when a request cannot succeed. Initiation
can explain system behavior such as expiration. None of these headings selects an
execution engine or scheduler.

Inputs and outputs can remain prose. When names and types need to be referenced,
use H5 entries:

```markdown
#### Inputs

##### Resource

**Type:** [Resource](#resource-type)

An opaque resource identity.

##### Lifetime

**Primitive:** Duration

A positive duration.
```

The linked type must be declared by this concept, either as a parameter or an
introduced type. Primitive names are authored labels, not a built-in type registry.
Use a `###### Type` description for a compound type expression; dotspec retains it
without interpreting its grammar. Results follow the same pattern under Outputs.
Do not duplicate a signature unless the named structure is useful.

## Separate a definition from a use

Declare generic identity types with `## Type parameters` and H3 headings. Declare
identities the concept introduces with `## Introduced types` and H3 headings.
A reusable definition should describe external identities through parameters rather
than depend on another concept's internals.

In a system specification, explicitly name each inclusion under `## Concept uses`:

```markdown
## Concept uses

### PublicLinks

#### Definition

[Shorthand](concepts/shorthand.concept.spec.md)

#### Bindings

##### Targets

**Parameter:** [Target](concepts/shorthand.concept.spec.md#target)
**Primitive:** URL

### PublicExpiry

#### Definition

[Expiry](concepts/expiry.concept.spec.md)

#### Bindings

##### Resources

**Parameter:** [Resource](concepts/expiry.concept.spec.md#resource)
**Bound type:** [Code](concepts/shorthand.concept.spec.md#code)
**Supplying concept use:** [PublicLinks](#publiclinks)
```

Here PublicExpiry's resource type comes specifically from PublicLinks. An
InternalLinks use of Shorthand would be distinct, even though both use the same
Code definition. Bindings choose a primitive or a Type/Use pair. A parameter can
have at most one binding per use. Missing bindings remain unfinished declarations
and produce advice; references to an incompatible definition or scope are errors.

A use may have a Configuration subsection for application-specific choices. It
may also have Dependencies containing links to other uses in the same system and
an explanation of why those uses are needed. Mutual dependencies are permitted.
These declarations support application scoping and explanation; they do not claim
that one concept calls another. Type bindings and directory placement do not
implicitly create uses or dependencies.

## State coordination explicitly

Add `## Synchronizations` to the system specification and give each rule an H3
heading. When, Where, Then and Semantics are H4 subsections. Start in prose.

To distinguish action references, add H5 entries under When or Then with explicit
Concept use and Action fields:

```markdown
### Start a lifetime

#### When

##### Registration completed

**Concept use:** [PublicLinks](#publiclinks)
**Action:** [Register](concepts/shorthand.concept.spec.md#register)

###### Result

Successful completion returning the newly registered code.

#### Then

##### Allocate lifetime

**Concept use:** [PublicExpiry](#publicexpiry)
**Action:** [Allocate](concepts/expiry.concept.spec.md#allocate)

###### Arguments

Use the returned code as Resource and one hour as Lifetime.

#### Semantics

Event-based delivery with retained completion events and retries using the same
arguments. Allocation is idempotent. No atomic transaction spans the concepts.
```

Every action reference identifies a named use and an action offered by its
definition. Arguments and results remain narrative; Where can describe conditions,
queries or matching. Multiple When entries do not silently define correlation,
ordering or conjunction semantics: state those requirements in the rule.

Explain consequential ordering, atomicity, failures, retries and timing. The
original atomic synchronization model and the newer event-based model make different
promises. Dotspec can document either; omitting Semantics chooses neither.
A heading without Use or Action fields remains supporting content, not a fabricated
invocation. The complete [expiring-links example](../examples/concept-design/specs/system.spec.md)
includes repeated uses, bindings and coordination with explicit delivery limits.

## Inspect and generate views

```sh
dotspec concept check shorthand
dotspec check
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
dotspec model docs --path . --relative-to docs > docs/concepts.md
```

Create the destination directory before redirecting output. Authoring checks handle
structure and ordinary links. Model export additionally diagnoses composition
names, reference kinds, declaration scope and bindings. It returns status 1 for
composition errors while retaining the snapshot. External or unacquired targets
remain unresolved and produce advice; export does not recursively acquire a
concept package. Select a scope containing the definitions needed for resolution.
A saved snapshot can be validated without reading the working tree.

The generated reference presents purpose and operational principles before state
and actions, then system composition, binding tables and dependency diagrams. Only
resolved action references produce coordination edges; unresolved declarations and
findings remain visible in the narrative. `--relative-to` controls source links,
not which files are inspected. [The generated example](reference/concept-design-example.md)
is maintained by `mise run generate`.

Review purpose, independence, adequate state, end-to-end behavior and composition
promises with people who understand the intended use. Structural validation cannot
establish those qualities, execution guarantees or requirement satisfaction.

This guide adapts Daniel Jackson's [concept criteria](https://essenceofsoftware.com/tutorials/concept-basics/criteria/),
[state](https://essenceofsoftware.com/tutorials/concept-basics/state/),
[operational principles](https://essenceofsoftware.com/tutorials/concept-basics/principle/)
and [dependencies](https://essenceofsoftware.com/tutorials/concept-basics/dependency/),
and Meng and Jackson's [WYSIWID paper](https://arxiv.org/html/2508.14511v2).
The Markdown conventions are dotspec choices. See the
[model rationale](specification-model.md#concept-definitions-and-composition) for
how these ideas map to the core domain.
