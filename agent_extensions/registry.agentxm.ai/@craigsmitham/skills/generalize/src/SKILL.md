---
name: generalize
description: Generalize a codebase or a specified part to support a broader, coherent family of cases.
---

Generalize the requested scope: broaden its applicability by expressing
what related cases have in common and removing assumptions and dependencies
that unnecessarily tie it to particular cases.

Distinguish reusable mechanisms from case-specific policy and behavior.
Shared use does not by itself make code general. Move responsibilities
that belong to particular domains or applications to their owners, keeping
the general part independent of them. Use composition or explicit inputs
where variation is needed; avoid adding options merely to retain unrelated
responsibilities in one abstraction.

The result should depend on the properties those cases actually require,
preserve required behavior for existing cases through their owning
integrations, and remain simple to use. Verify both that existing cases
still work and that the general part can serve another intended case
without depending on the original case's concepts.

Let the user's intent and intended reuse determine the extent of generality.
