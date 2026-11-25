# What does a quality model help a team decide?

Suppose a parcel locker comes back online after an outage but has forgotten which compartments hold parcels. The operator who tracks uptime and the courier who needs to retrieve a parcel can give different accounts of how well the locker works. A **quality model** gives them shared terms for the properties they care about. Naming more than one property helps them see what a single observation can miss.

In the ISO/IEC 25010:2023 product quality model, *availability* and *recoverability* are both aspects of reliability. The first concerns readiness for use when needed. The second concerns restoring affected data and the intended system state after disruption. Measuring how soon the locker accepts requests tells the team little about whether it can still associate a parcel with its compartment. This edition organizes product quality under nine characteristics, each with more specific terms that can prompt such questions. [ISO/IEC 25010:2023, §§3.5.2, 3.5.4 and 4.1–4.2](https://www.iso.org/standard/78176.html).

## What can two systems share?

A dotspec model can draw on an international standard or a concern familiar to a few projects. Its ISO starter uses the standard's categories with original prompts for authors. The much smaller *Personal data quality* model names durability for tools that keep a person's work between sessions. Either kind gives a team a vocabulary to select from. The locker team still has to identify which records matter, which interruptions to cover and what recovery it will accept. ISO/IEC 25030 treats selecting a characteristic and specifying a requirement with its measure and acceptance criteria as distinct decisions. [Personal data quality](../../../examples/todo/specs/quality-models/personal-data.quality-model.md) · [ISO/IEC 25030:2019, §7.4.2, steps 2–5](https://www.iso.org/standard/72116.html).

Applying a model in dotspec copies selected attributes and their ancestors into the system as independent declarations. That lets a team tailor a definition to its product. If a later edition broadens the reusable definition, a team can review whether to adopt that change without silently changing the meaning of existing classifications. A requirement draws its labels from its own system's declarations. The acceptable result belongs with that requirement. [System quality attributes](system-quality-attributes.md).

When several projects could reuse the vocabulary, follow [Define a quality model](../guides/define-a-quality-model.md).
