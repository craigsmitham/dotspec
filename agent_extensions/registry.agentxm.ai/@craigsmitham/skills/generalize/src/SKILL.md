---
name: generalize
description: Generalize a codebase or a specified part to support a broader, coherent family of cases.
disable-model-invocation: true
---

Generalize the requested scope to support a broader, coherent family of
cases by expressing their essential common properties and removing
unnecessary assumptions that tie the capability to particular cases.

For opportunity questions, assess candidates and their tradeoffs without
implementing them. For change requests, carry the requested outcome through
implementation and verification within the authorized scope.

Keep case-specific policy in its owning integrations. Shared use alone
does not establish coherent commonality. Build on suitable existing
capabilities, adding only the generality justified by the intended cases.

The result should preserve required behavior for existing cases, remain
simple to use, and serve another intended case without depending on the
original case's domain concepts.

Let the user's intent and intended reuse determine the extent of generality.
