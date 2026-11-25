# What belongs in a specification?

Suppose a task-list owner records a reply they are waiting for, then closes the application. Three artifacts mention durability. The reusable quality model defines a term for many products. The system quality attribute identifies the owner's acknowledged decisions as the information that matters here. The system quality requirement says exactly which saved details must survive a restart on the same database. The latter two record decisions about this product. [Quality model](../../../examples/todo/specs/quality-models/personal-data.quality-model.md#durability) · [System quality attribute](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md#definition) · [System quality requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md#statement).

## Which subject owns the decision?

A specification records product intent for one independently maintained subject. Its `*.spec.md` file declares that subject as its root. A use case owns its steps and extensions. An extension can have a name and links so a reviewer can discuss it precisely. Edit it in the use case whose interaction it changes. A promise shared across interactions has a different lifecycle. Give it one requirement file and link to it wherever it applies. [Use cases](use-cases.md) · [Requirements](requirements.md).

## How do the pieces form a product specification?

The specifications in product scope, their owned parts, and the relationships among them compose the product specification. A system can have open questions while its authors settle individual requirements. The generated report brings the current decisions together and links readers to the declarations where they can change them.

References supply vocabulary or external authority. Evidence records what the team observed, and guidance helps authors decide how to do the work. Ask what would change if the document were revised. Changing a restart requirement changes the product's promise. Recording a failed restart test changes the evidence about whether the implementation keeps it.

For Personal task list, the forms interface adopts the HTTP redirect provision specifically after an accepted form submission. The local declaration gives the edition and applicability. A reference to the form template explains where to see the controls. Following a procedure link gives a reviewer a way to assess the product. An observation from a particular run can support a claim about the revision and conditions assessed. [Task list forms](../../../examples/todo/specs/personal-task-list.system.spec.md#task-list-forms) · [Restart procedure](../../../examples/todo/specs/requirements/restart.requirement.spec.md#verification-procedures)
