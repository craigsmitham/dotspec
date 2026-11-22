# Specify a capability

Name an ability in terms that would still make sense if a different team, tool or procedure provided it. From the owning system or component directory, use `dotspec capability init local-task-storage`, substituting a slug for your ability. Keep the resulting `*.capability.spec.md` as the source for that capability.

Under **Description**, say what the ability makes possible. Under **Supports**, link to the need it helps meet and explain the contribution. Identify a known system or component under **Provided by**. If people or procedures help deliver the ability, describe their part under **Enablers**. When it relies on another ability, link to that capability under **Depends on**. [Local task storage](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).

Read the links from the performer's perspective. Could the provider deliver the described ability while the performer still struggles? The answer can reveal what else the person needs. [How capabilities connect needs to providers](../explanations/capabilities.md) · [Map needs to capabilities](map-needs-to-capabilities.md).
