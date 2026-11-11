"""Relationship contracts and their generated lookup surface.

The schema owns the metadata; this module validates it and presents it without
maintaining another endpoint or semantics table.
"""
from typing import Any
from model_views import field, text, link

CATEGORIES = {
    'declaration-containment': 'Declaration containment',
    'structural-context': 'Structural context',
    'semantic-association': 'Semantic association',
    'contextual-citation': 'Contextual citation',
    'derived-projection': 'Derived projection',
    'provenance': 'Source and inspection provenance',
    'representation': 'Document representation',
}
REQUIRED = ('relationship_category', 'relationship_scope', 'relationship_cycles',
            'relationship_basis', 'relationship_inverse', 'relationship_ordering')


def validate_relationships(catalog: dict[str, Any]) -> None:
    classes = catalog['classes']
    for name, cls in classes.items():
        if 'Assertion' in cls['ancestors'] and not cls['abstract']:
            for endpoint in ('subject', 'object'):
                if not cls['slots'][endpoint]['required']:
                    raise ValueError(f'{name}.{endpoint}: assertion endpoints must be required')
            if not cls['slots'].get('projection_path', {}).get('required'):
                raise ValueError(f'{name}: assertion projection_path must be required')
            if not cls['annotations'].get('projection_of'):
                raise ValueError(f'{name}: assertions must identify canonical projection fields')
        trace_headings = set()
        for prop, slot in cls['slots'].items():
            path = f'{name}.{prop}'
            annotations = slot['annotations']
            heading = annotations.get('trace_heading')
            targets = annotations.get('trace_targets')
            if heading or targets:
                if not heading or not targets or heading in trace_headings:
                    raise ValueError(f'{path}: trace headings and targets must be paired and unique')
                trace_headings.add(heading)
                if 'TraceableElement' not in cls['ancestors'] or not slot['multivalued']:
                    raise ValueError(f'{path}: trace relationships require a traceable source and multiple targets')
                if set(targets.split(' or ')) != set(slot['reference_ranges']):
                    raise ValueError(f'{path}: trace targets must match the canonical endpoint ranges')
                refs = cls['slots'].get(prop + '_references', {})
                if refs.get('range') != 'SourceReference' or not refs.get('multivalued'):
                    raise ValueError(f'{path}: trace relationship needs source-reference occurrences')
                if annotations.get('requires_content') != prop + '_content':
                    raise ValueError(f'{path}: trace relationship needs its complete authored section')
            if annotations.get('requires_content'):
                _, content_slot = field(catalog, name + '.' + annotations['requires_content'])
                if content_slot['range'] != 'RichContent':
                    raise ValueError(f'{path}: requires_content must identify RichContent')
            peer = annotations.get('paired_with')
            if peer and (peer == prop or peer not in cls['slots'] or
                         cls['slots'][peer]['annotations'].get('paired_with') != prop):
                raise ValueError(f'{path}: paired_with must name a reciprocal peer field')
            if not slot['reference_ranges']:
                continue
            for key in REQUIRED:
                if not annotations.get(key):
                    raise ValueError(f'{path}: missing relationship metadata {key}')
            if annotations['relationship_category'] not in CATEGORIES:
                raise ValueError(f'{path}: invalid relationship category')
            if annotations['relationship_cycles'] not in ('forbidden', 'permitted', 'not-applicable'):
                raise ValueError(f'{path}: invalid relationship cycle policy')
            inverse = annotations['relationship_inverse']
            if inverse != 'none':
                for inverse_path in inverse.split(' or '):
                    source, target_slot = field(catalog, inverse_path)
                    if not target_slot['reference_ranges']:
                        raise ValueError(f'{path}: inverse must reference records')
                    if not any(t in classes[source]['ancestors'] for t in slot['reference_ranges']):
                        raise ValueError(f'{path}: inverse source is not an allowed endpoint')
                    # Embedded wrappers carry their owner's members; owned_memberships
                    # validates that owner endpoint instead of treating the wrapper as a record.
                    if 'Record' in cls['ancestors'] and not any(t in cls['ancestors'] for t in target_slot['reference_ranges']):
                        raise ValueError(f'{path}: inverse target does not permit this source')


def cardinality(slot: dict[str, Any]) -> str:
    # Required list presence is distinct from minimum list length.
    if slot['multivalued']:
        low = slot.get('minimum_cardinality', 0)
        high = slot.get('maximum_cardinality', '*')
        return f'{low}..{high}' + ('; field required' if slot['required'] else '')
    return '1' if slot['required'] else '0..1'


def declared_relationships(catalog: dict[str, Any]):
    for name, cls in catalog['classes'].items():
        parent = catalog['classes'].get(cls.get('is_a'), {})
        for prop, slot in cls['slots'].items():
            if slot['reference_ranges'] and slot != parent.get('slots', {}).get(prop):
                yield name, prop, slot


def relationships_markdown(catalog: dict[str, Any]) -> bytes:
    lines = ['# Model relationships', '',
        'Generated from [the canonical model](../../../model/dotspec.yaml). Do not edit this file.', '',
        f'Model `{catalog["model_version"]}`. This lookup covers every record-reference field, '
        'including inherited endpoint refinements. Unchanged inherited fields are documented on their base class. '
        'Embedded values and their fields remain in the [complete reference](index.md).', '',
        'Typed fields are authoritative. Assertions are optional derived views, not a complete substitute for those fields. '
        '`projection_path` retains the canonical field and role of each assertion. Reverse navigation is derived; '
        'stored owner backlinks are checked against canonical membership fields.', '',
        'Cardinalities describe snapshots, including unfinished declarations. An optional resolved target does not mean '
        'a complete authored design can omit the definition it needs. Required list presence does not imply a nonempty list. '
        'Ordering preserves authored presentation unless a field explicitly says otherwise; it does not impose execution order. '
        'Cycle policies apply to the named relationship, not the entire graph: owner/member backlinks are intentionally reciprocal.', '',
        'Relationship categories describe edges. They are independent of the core/supporting/generic classifications of model subsets.', '',
        '## Relationship categories', '', '| Category | Scope |', '| --- | --- |',
        '| Declaration containment | A declaration and its owned members; source-local, with checked owner backlinks. |',
        '| Structural context | Directory-derived discovery context; no structural decomposition or semantic inheritance. |',
        '| Semantic association | Explicit domain claims with typed targets and relationship-specific scope. |',
        '| Contextual citation | Authored reference occurrences; citations alone do not establish domain claims. |',
        '| Derived projection | Optional assertion endpoints projected from canonical fields. |',
        '| Source and inspection provenance | Evidence linking model records, source occurrences and inspection results. |',
        '| Document representation | Syntax structure and indexes, rather than domain ownership. |', '',
        '## Fields', '']
    for name, prop, slot in declared_relationships(catalog):
        a = slot['annotations']
        lines += [f'### {name}.{prop}', '',
                  f'{link(name)} → ' + ' or '.join(link(t) for t in slot['reference_ranges']) + '.', '',
                  text(slot['description'] or a['relationship_basis']), '',
                  f'- Category: {CATEGORIES[a["relationship_category"]]}. Cardinality: `{cardinality(slot)}`.',
                  f'- Scope: {text(a["relationship_scope"])}',
                  f'- Basis: {text(a["relationship_basis"])}',
                  f'- Ordering: {text(a["relationship_ordering"])}. Cycles: {a["relationship_cycles"]}.',
                  f'- Stored inverse: `{text(a["relationship_inverse"])}`; other reverse navigation is derived.']
        if a.get('paired_with'):
            lines.append(f'- Declaration/content pair: `{name}.{a["paired_with"]}`; both fields are present or both absent.')
        if a.get('requires_content'):
            lines.append(f'- Complete section content required when this field is present: `{name}.{a["requires_content"]}`.')
        lines.append('')
    lines += ['## Assertion projections', '',
              'Every concrete assertion declares its canonical field paths. An assertion must identify one '
              'of these paths with `projection_path` and agree with that field’s subject, target and provenance.', '',
              '| Assertion | Canonical paths |', '| --- | --- |']
    for name, cls in catalog['classes'].items():
        if cls['annotations'].get('projection_of'):
            lines.append(f'| {link(name)} | `{cls["annotations"]["projection_of"]}` |')
    return ('\n'.join(lines).rstrip()+'\n').encode()
