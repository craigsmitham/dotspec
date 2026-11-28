# Where does a product decision live?

A specification can live in Markdown, an executable suite, a schema, a model, or a diagram. Give each decision one governing source. Ask what the team would edit to change the promised behavior.

For example, a feature can promise that a returning owner sees which replies need attention. A Gherkin suite can own the cases that determine when a reply becomes due. The feature introduces the promise in prose and names the suite under `specifications`. Adding another timing case requires one edit in Gherkin. The report lists the named suite under Requirements, alongside requirements written in Markdown.

## When does a test define a requirement?

A suite defines behavior when the team accepts its cases as the governing examples. A test that exercises a promise already specified elsewhere is a verification procedure. In the To-do example the restart requirement owns the preservation promise. Its verification procedure stops the service and compares the recovered identities and fields. The test could change how it starts the service while the requirement stays the same. [Restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md)

An executable specification can itself be run against the product. An observation records what happened for the recorded revision and conditions. A model check describes the model within the explored bounds. A claim about the implementation needs evidence about the implementation.

## Give readers a route

Use an existing dotspec relationship when an authored specification refines another. Use an adoption to identify a provision, edition, and local applicability. Link a worked example or explanation from the passage it helps the reader understand. In the inventory, contextual links, governing sources, and verification procedures remain distinguishable. [Catalog a specification](../guides/catalog-a-specification.md)
