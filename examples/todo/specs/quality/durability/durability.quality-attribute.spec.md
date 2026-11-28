---
type: quality-attribute-spec
specifications:
  - name: Saved follow-up dates across time zones
    defines: requirements
    source: ../../../tests/acceptance.py#verify_changed_day
    verification:
      - name: Run the acceptance suite on Unix
        source: ../../../tests/acceptance.py#run
verification:
  - name: Recover acknowledged changes after restart
    source: ../../../tests/acceptance.py#verify
---

# Durability

## Definition

Preservation of the owner’s acknowledged decisions between local sessions, including the reply they expect for a waiting errand and the day they chose to follow up. The follow-up cue can change overnight even though those decisions survive intact.

The [time-zone cases](#specification.saved-follow-up-dates-across-time-zones) govern what the owner sees when the same saved list is opened on different local days. The suite exercises those cases on Unix.

[Source definition](../../quality-models/personal-data.quality-model.md#durability).

[How SQLite keeps a committed transaction across a crash or power loss](https://www.sqlite.org/atomiccommit.html).
