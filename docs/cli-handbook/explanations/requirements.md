# What makes a requirement assessable?

Suppose a completed task comes back after a restart with the right title and a new identity. Has the application kept its promise? A requirement gives that question a determinate answer by stating which result is required and when the promise applies.

The To-do application's restart requirement covers changes the service has acknowledged, on the same database and functioning storage. Each surviving entry must retain its identity, title and state, including the expected reply and follow-up date for a waiting entry. Deleted entries must stay absent. Counting missing, unexpected or changed entries gives a **measure**. Requiring zero discrepancies across all three views gives a **criterion** for satisfaction. A new identity violates the promise even if the screen looks the same. Losing an edit before the service acknowledges it raises a different question, outside this requirement's conditions. [Preserve tasks across restart](../../../examples/todo/specs/requirements/restart.requirement.spec.md).

## Where does the promise live?

The more general To-do requirement says where a person can find saved work on returning to the list. The restart requirement *refines* it by spelling out what must survive a process stop. Giving each commitment its own requirement file lets a use case, a more specific requirement and a test all refer to the same decision. Someone changing the restart test can find out why comparing titles alone would miss a failure. [Keep recorded work available](../../../examples/todo/specs/requirements/availability.requirement.spec.md) · [The restart criterion](../../../examples/todo/specs/requirements/restart.requirement.spec.md#criterion).

ISO/IEC/IEEE 29148 distinguishes the qualities of an individual requirement, including singularity and verifiability, from the completeness and consistency of a set. The one-file convention is dotspec's way to give each requirement an address. A qualitative statement can already say what would violate it. For instance, a deleted task appearing in any view violates the general requirement. A separate numeric target adds value when the decision depends on a measured bound. [ISO/IEC/IEEE 29148:2018, §§5.2.4–5.2.6](https://www.iso.org/standard/72089.html).

## What would count as evidence?

Bertrand Meyer suggests looking for a counterexample when judging verifiability. Wording that rules out an observable result can guide a test before any test has run. The restart test records entries in all three views, restarts the application and compares the entries it gets back. A passing run is verification evidence for those conditions. Validation asks whether preserving that information actually helps a person resume unfinished work. Observing someone return to their list could reveal that they also need a reminder of the next step. The existing test could still pass. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §4.12](https://doi.org/10.1007/978-3-031-06739-6) · [ISO/IEC/IEEE 29148:2018, §§6.5.2–6.5.3](https://www.iso.org/standard/72089.html) · [To-do verification example](../../../examples/todo/specs/requirements/restart.requirement.spec.md#evidence).

To put a commitment into words and try a case on either side of it, follow [Write an assessable requirement](../guides/write-an-assessable-requirement.md).
