# What does a companion add to a specification?

A specification states product intent in Markdown prose. Some of that intent is easier to state precisely, run or explain in another form. Examples include a Quint, Alloy or TLA+ model of the states a task can pass through, a Gherkin scenario the team runs as an example, an acceptance test, and a design note. A **companion** is such a file or resource. The specification lists it in its frontmatter with a **role** that says how the companion relates to the specification. [List a companion](../guides/list-a-companion.md).

## Why declare the relationship?

A link tells a reader that two documents are related. It does not tell them which one to correct when they disagree. Suppose a model of the task list allows an entry to appear in two views, while the requirement says each saved entry appears once. If the model formalizes the requirement, one of them has a defect and someone must decide which. If the model only explains an earlier design, the reader can set it aside.

dotspec does not infer that from a link, a file type or a directory. The same test file can state one team's examples and serve another team as a regression check. The author declares the role, and the report shows it beside each companion.

## What does each role claim?

| Role | The companion | Authority |
| --- | --- | --- |
| `formalizes` | States the specification more precisely, such as a property in a formal model or a schema. It must preserve the specification's meaning. | Normative. A disagreement between the two is a defect. |
| `illustrates` | Gives concrete cases that are part of the specification, such as executable examples. | Normative. Changing a case changes what is specified. |
| `verifies` | Describes a procedure for assessing the specification, such as an automated test or a manual protocol. | A source of evidence. The outcome of a particular run is recorded separately. |
| `explains` | Gives background or rationale. | Informative. |

Each subject's report groups its companions under **Formalized by**, **Illustrated by**, **Verified by** and **Explained by**. The outermost system's report links to an index of every companion in that system and the systems within it.

## Is a test an illustration or a procedure?

Hillel Wayne offers a test for this: “Tests only become specifications when they change iff the requirements do.” A test you would edit during an internal refactoring is tied to the implementation. List it as `verifies`. A scenario that changes exactly when the promise changes can be listed as `illustrates`. If a file serves both purposes, list it once under each role. [Wayne, “What is a ‘Specification’?” (2021)](https://buttondown.com/hillelwayne/archive/what-is-a-specification).

The To-do restart requirement lists the `verify` function in the acceptance test as the procedure that verifies it. The test kills the service and compares each view's identities after it starts again. If someone restructured the test, for example to start the service a different way, the promise would stay the same, so the test verifies the requirement rather than illustrating it. The broader availability requirement lists the README's instructions for running the application, which a reviewer can follow to exercise the forms across the list views. Each requirement's **Evidence** section still describes what its procedure exercises and what it leaves out. [Restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md) · [Restart acceptance test](../../../examples/todo/tests/acceptance.py) · [Availability requirement](../../../examples/todo/specs/requirements/availability.requirement.spec.md) · [To-do companions](../../../examples/todo/specification/companions.md).

## What does a passing check establish?

`dotspec check` establishes that each companion's relative target exists in the repository. It does not fetch an absolute URL. Where the target names a definition in a format dotspec can read, the check confirms that the file declares that name. It does not compare the companion's meaning with the prose.

A model checker that finds no counterexample within its bounds supports a claim about that model at that revision, within the explored bounds. The code can still differ from a verified model. Conformance of the implementation needs its own evidence, such as a test that verifies the requirement against the running product. In the same way, a passing run of a `verifies` companion supports a claim about the product revision and conditions it ran under. Record that observation in the project's verification results rather than in the specification. The requirement changes only when the promise does. Whether the promise serves the person is a separate question of validation. [What would count as evidence?](requirements.md#what-would-count-as-evidence) · [Record results for a procedure](../guides/list-a-companion.md#record-results-for-a-procedure).
