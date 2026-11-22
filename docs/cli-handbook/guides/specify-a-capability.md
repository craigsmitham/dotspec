# Specify a capability

Describe what a person or organization needs to be able to do. Try substituting another provider. If replacing “Task service” with a different component leaves the need to retain task changes between sessions intact, **Local task storage** names the ability. Richard Allen uses this distinction to identify capabilities even when a team has yet to select a tool. [Allen, *User Needs Mapping* (1st ed., 2025), ch. 4, step 4](https://www.userneedsmapping.com/moviegoer-walkthrough) · [How capabilities connect needs to providers](../explanations/capabilities.md).

Start in the directory of the system or component that will own the file. For a new capability, use `dotspec capability init local-task-storage`, substituting a lowercase, hyphenated name. If the capability is already specified, edit its `*.capability.spec.md` file. Name the ability in the top-level heading and describe what it makes possible under **Description**.

Under **Supports**, link to the job, desired outcome or other specified need it helps meet. Explain the contribution if it is easy to misread. Put each relationship under a level-two heading and write links relative to the capability file. The system that owns the file may draw on an entirely different provider. [Local task storage’s links](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md#supports).

- When you know which system or component supplies the ability, link to its specification under **Provided by**. Leave the provider open if the team has yet to decide.
- When another *ability* must be available, link to that capability under **Depends on**.
- When a person or procedure helps deliver the ability, describe that contribution under **Enablers**. A dispatcher phoning technicians for available times belongs here.

From the owning directory, run `dotspec capability check local-task-storage` with your slug. If a link fails, correct its path or heading fragment in the capability file and run the check again. Read the result from the performer's perspective. Could the provider deliver the described ability while the performer still struggles? To investigate what else they depend on, [map the need to its supporting capabilities](map-needs-to-capabilities.md).
