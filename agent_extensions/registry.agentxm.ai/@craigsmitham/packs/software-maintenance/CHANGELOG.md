# Changelog

## 0.2.1

Clarifies migration after verifying installation and upgrade behavior with
AXM 0.51.0. The seven members and their versions are unchanged. A normal update
selects the new membership but can retain former members' native links as
unmanaged entries. The migration below removes the old pack while AXM still
owns those entries, then installs the new pack.

## 0.2.0

Reorganizes the pack around continuing maintenance: prioritization, dependable
behavior and feedback, support lifecycles, removal, simplification, independent
change, and selected quality improvements.

- Adds `assess`, `stabilize`, `refresh`, and `simplify` at 0.1.0.
- Retains `prune`, `modularize`, and `improve`, requiring their strengthened
  0.2.0 contracts.
- Removes `align`, `reconcile`, `reuse`, and `generalize` from pack membership.
  Those packages remain available independently; their 0.1.1 releases update
  companion documentation, recommendations, and opportunity-question handling
  without changing their purpose.
- Adds conditional sequencing, recurring-maintenance guidance, and examples
  distinguishing opportunity questions from implementation requests.
- Preserves standalone installation and explicit-only invocation for every member.

### Retaining former members

Declare any former member you still need directly before replacing the pack:

```bash
axm install @craigsmitham/skills/align
axm install @craigsmitham/skills/reconcile
axm install @craigsmitham/skills/reuse
axm install @craigsmitham/skills/generalize
```

Install only the companions you intend to retain. Direct selections are
independent of membership. For a clean migration from 0.1.0 with AXM 0.51.0,
preview removal, uninstall the old pack, and install the new one:

```bash
axm uninstall @craigsmitham/packs/software-maintenance --preview
axm uninstall @craigsmitham/packs/software-maintenance
axm install @craigsmitham/packs/software-maintenance
axm list
```

Review the preview for local modifications and retained direct selections.
The verified sequence removes pack-only members and preserves directly
installed companions. Version 0.1.0 remains a historical release.

If you already used `axm update`, inspect `axm list` for unmanaged former
members. With AXM 0.51.0, successful lint and sync checks do not establish that
those entries were removed. Preserve local changes and confirm ownership before
removing any remaining native entries; do not assume that uninstalling the
updated pack will clean up entries it no longer owns.

## 0.1.0

Initial seven-member pack: `align`, `reconcile`, `modularize`, `generalize`,
`reuse`, `prune`, and `improve`.
