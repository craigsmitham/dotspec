"""Validated documentation views over the canonical model catalog."""
from typing import Any
import re

CLASSIFICATIONS = ('core', 'supporting', 'generic')


def validate_groups(catalog: dict[str, Any]) -> None:
    subsets = catalog['subsets']
    if not subsets:
        raise ValueError('model must declare subsets')
    for name, subset in subsets.items():
        if subset['classification'] not in CLASSIFICATIONS:
            raise ValueError(f'{name}: invalid classification {subset["classification"]!r}')
        if not subset['description']:
            raise ValueError(f'{name}: subset needs a description')
    for name, cls in catalog['classes'].items():
        if not cls['subsets']:
            raise ValueError(f'{name}: unclassified class (declare in_subset explicitly)')
        if len(cls['subsets']) != len(set(cls['subsets'])):
            raise ValueError(f'{name}: duplicate subset membership')
        for subset in cls['subsets']:
            if subset not in subsets:
                raise ValueError(f'{name}: unknown subset {subset}')


def validate_declared_memberships(raw_model: dict[str, Any]) -> None:
    subsets = raw_model.get('subsets', {})
    def walk(value: Any, path: str) -> None:
        if isinstance(value, dict):
            for subset in value.get('in_subset', []):
                if subset not in subsets:
                    raise ValueError(f'{path}: unknown subset {subset}')
            for key, child in value.items():
                walk(child, path + '/' + str(key))
        elif isinstance(value, list):
            for index, child in enumerate(value):
                walk(child, path + '/' + str(index))
    walk(raw_model, '')


def field(catalog: dict[str, Any], path: str) -> tuple[str, dict[str, Any]]:
    parts = path.split('.')
    if len(parts) < 2 or parts[0] not in catalog['classes']:
        raise ValueError(f'invalid field path {path}')
    owner = parts[0]
    for i, name in enumerate(parts[1:]):
        slot = catalog['classes'][owner]['slots'].get(name)
        if slot is None:
            raise ValueError(f'{path}: absent field {owner}.{name}')
        if i < len(parts)-2:
            if len(slot['inline_ranges']) != 1 or slot['multivalued']:
                raise ValueError(f'{path}: intermediate field must be one embedded object')
            owner = slot['inline_ranges'][0]
    return parts[0], slot


def validate_view(catalog: dict[str, Any], view: dict[str, Any]) -> None:
    if set(view) != {'title', 'description', 'subsets', 'types', 'diagrams'}:
        raise ValueError('view must declare title, description, subsets, types, and diagrams only')
    if not view['title'] or not view['description']:
        raise ValueError('view needs title and description')
    selected = view['subsets']
    if not selected or len(selected) != len(set(selected)):
        raise ValueError('view needs unique subset selectors')
    for name in selected:
        if name not in catalog['subsets']:
            raise ValueError(f'view: unknown subset {name}')
        if catalog['subsets'][name]['classification'] != 'core':
            raise ValueError(f'core view: subset {name} is not core')
    for name, paths in view['types'].items():
        if name not in catalog['classes']:
            raise ValueError(f'view: unknown class {name}')
        if not set(catalog['classes'][name]['subsets']) & set(selected):
            raise ValueError(f'{name}: outside selected subsets')
        if len(paths) != len(set(paths)):
            raise ValueError(f'{name}: duplicate field selector')
        for path in paths:
            field(catalog, name + '.' + path)
    for diagram in view['diagrams']:
        if set(diagram) - {'direction', 'show_cardinality'} != {'title', 'description', 'nodes', 'paths', 'inheritance'}:
            raise ValueError('diagram must declare title, description, nodes, paths, and inheritance only')
        if diagram.get('direction', 'LR') not in {'LR', 'TB'}:
            raise ValueError('diagram direction must be LR or TB')
        nodes = diagram['nodes']
        if not nodes or len(nodes) != len(set(nodes)):
            raise ValueError('diagram needs unique node selectors')
        if type(diagram.get('show_cardinality', True)) is not bool:
            raise ValueError('diagram show_cardinality must be boolean')
        if type(diagram['inheritance']) is not bool:
            raise ValueError('diagram inheritance must be boolean')
        for node in nodes:
            if node not in catalog['classes'] or not re.fullmatch(r'[A-Za-z][A-Za-z0-9_]*', node):
                raise ValueError(f'diagram: invalid class {node}')
        if len(diagram['paths']) != len(set(diagram['paths'])):
            raise ValueError('diagram has duplicate paths')
        for path in diagram['paths']:
            source, slot = field(catalog, path)
            targets = slot['reference_ranges'] or slot['inline_ranges']
            if not targets or source not in nodes or any(t not in nodes for t in targets):
                raise ValueError(f'{path}: diagram must include source and all declared class targets')


def link(name: str) -> str:
    return f'[{name}](index.md#{name.lower()})'


def text(value: str) -> str:
    return ' '.join(value.split()).replace('|', r'\|')


def subset_navigation(catalog: dict[str, Any]) -> list[str]:
    lines = ['## Subsets and classifications', '',
             'Classification describes each subset’s role in dotspec. Membership is explicit and may overlap; '
             'it does not classify individual authored requirements or change inheritance.', '',
             'See the [core domain reference](core.md) for focused summaries and diagrams.', '']
    for classification in CLASSIFICATIONS:
        lines += [f'### {classification.title()}', '']
        for name, group in catalog['subsets'].items():
            if group['classification'] != classification:
                continue
            members = [f'[{n}](#{n.lower()})' for n, c in catalog['classes'].items() if name in c['subsets']]
            lines += [f'**{group["title"]}** (`{name}`). {text(group["description"])}', '', ', '.join(members) + '.', '']
    return lines


def core_markdown(catalog: dict[str, Any], view: dict[str, Any]) -> bytes:
    validate_view(catalog, view)
    lines = [f'# {view["title"]}', '',
        'Generated from [the canonical model](../../../model/dotspec.yaml) and '
        '[the view definition](../../../model/views/core-domain.yaml). Do not edit this file.', '',
        text(view['description']), '',
        f'Model version: `{catalog["model_version"]}`. See the [complete reference](index.md) for all fields and constraints.', '',
        '## Core subsets', '', '| Subset | Responsibility | Selected types |', '| --- | --- | --- |']
    for name in view['subsets']:
        group = catalog['subsets'][name]
        names = [link(n) for n in view['types'] if name in catalog['classes'][n]['subsets']]
        lines.append(f'| {group["title"]} (`{group["classification"]}`) | {text(group["description"])} | {", ".join(names)} |')
    lines += ['', 'Subset membership organizes model definitions; it does not infer the level or importance of an authored instance. '
              'Shared types can appear in several subsets. Supporting and generic types appear as context where a diagram needs them.', '',
              '## Relationships', '', 'See the [relationship reference](relationships.md) for the complete edge contracts.', '',
              'Arrows use schema-owned relationship labels; reference tables retain exact canonical field paths. Diamond arrows represent declaration containment, dashed arrows structural context and citations, and plain arrows semantic associations or provenance; '
              'triangle arrows represent declared inheritance. Multiple edges with the same path show alternative target types; cardinality applies to the slot as a whole. '
              'Collection wrappers are omitted from diagrams; exact paths remain in the reference tables. '
              '`system_context` is a directory-derived association, not a use-case boundary or inherited requirement. '
              'Contextual links to business rules, jobs, and concepts do not assert enforcement, fulfillment, or compliance.', '']
    for diagram in view['diagrams']:
        lines += [f'### {diagram["title"]}', '', text(diagram['description']), '', '```mermaid', 'classDiagram', f'    direction {diagram.get("direction", "LR")}']
        for node in diagram['nodes']:
            lines.append(f'    class {node}')
            if catalog['classes'][node]['abstract']:
                lines.append(f'    <<abstract>> {node}')
        if diagram['inheritance']:
            for node in diagram['nodes']:
                parents = [catalog['classes'][node].get('is_a')] + catalog['classes'][node].get('mixins', [])
                for parent in parents:
                    if parent in diagram['nodes']:
                        lines.append(f'    {parent} <|-- {node}')
        for path in diagram['paths']:
            source, slot = field(catalog, path)
            for target in slot['reference_ranges'] or slot['inline_ranges']:
                cardinality = (str(slot.get('minimum_cardinality', 0)) + '..' + str(slot.get('maximum_cardinality', '*'))) if slot['multivalued'] else ('1' if slot['required'] else '0..1')
                marker = f' "{cardinality}"' if diagram.get('show_cardinality', True) else ''
                category = slot['annotations'].get('relationship_category')
                arrow = '*--' if category == 'declaration-containment' else '..>' if category in ('structural-context', 'contextual-citation') else '-->'
                label = slot['annotations'].get('relationship_label', path.split('.', 1)[1].replace('_', ' '))
                lines.append(f'    {source} {arrow}{marker} {target} : {label}')
        lines += ['```', '']
    lines += ['## Selected type reference', '']
    for name, paths in view['types'].items():
        cls = catalog['classes'][name]
        lines += [f'### {name}', '', text(cls['description']), '',
                  'Subsets: ' + ', '.join(f'`{s}`' for s in cls['subsets']) + '. ' + f'{link(name)} contains the complete definition.', '',
                  '| Field | Range | Shape | Meaning |', '| --- | --- | --- | --- |']
        for path in paths:
            _, slot = field(catalog, name + '.' + path)
            ranges = slot['any_of'] or [slot['range']]
            display = ' or '.join(link(n) if n in catalog['classes'] else f'`{n}`' for n in ranges)
            shape = 'required' if slot['required'] else 'optional'
            if slot['multivalued']:
                shape += '; ordered list' if slot['ordered'] else '; list'
            lines.append(f'| `{path}` | {display} | {shape} | {text(slot["description"])} |')
        lines.append('')
    return ('\n'.join(lines).rstrip() + '\n').encode()
