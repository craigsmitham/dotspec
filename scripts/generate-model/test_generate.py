"""Checks for the generation boundary that supplies runtime graph validation."""

import copy
from pathlib import Path
import tempfile
import unittest
from unittest import mock

import yaml
from linkml_runtime.utils.schemaview import SchemaView

import generate


class CatalogTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        model_bytes = generate.MODEL.read_bytes()
        cls.catalog = generate.make_catalog(
            SchemaView(str(generate.MODEL)), yaml.safe_load(model_bytes), model_bytes
        )

    def slot(self, class_name, slot_name):
        return self.catalog["classes"][class_name]["slots"][slot_name]

    def test_inherited_and_narrowed_reference_ranges(self):
        # Requirements are concrete specification roots and may carry classification labels.
        self.assertFalse(self.catalog["classes"]["Requirement"]["abstract"])
        self.assertIn("SystemContextSpecification", self.catalog["classes"]["Requirement"]["ancestors"])
        self.assertEqual(self.slot("LabelOccurrence", "owner")["reference_ranges"], ["Requirement", "QualityAttribute"])
        self.assertEqual(self.slot("SourceReference", "resolution")["inline_ranges"], ["Resolution"])
        # A union must not silently reopen the broad DomainElement parent range.
        owner = self.slot("ScenarioStep", "owner")
        self.assertEqual(owner["range"], "DomainElement")
        self.assertEqual(owner["reference_ranges"], ["Scenario", "UseCaseExtension"])
        self.assertNotIn("DomainElement", owner["reference_ranges"])
        self.assertEqual(self.slot("QualityModelAttribute", "owner")["reference_ranges"], ["QualityModel"])
        self.assertEqual(self.slot("LabelOccurrence", "governing_quality_model")["reference_ranges"], ["QualityModel"])

    def test_glossary_follows_model_topics_and_group_members(self):
        catalog = copy.deepcopy(self.catalog)
        feature = catalog["classes"]["Feature"]
        feature["subsets"] = ["business-requirements"]
        feature["slots"]["groups"]["reference_ranges"] = ["Requirement"]
        glossary = generate.glossary_markdown(catalog).decode()
        entry = glossary.split("## Feature\n", 1)[1].split("\n## Job", 1)[0]
        self.assertIn("### Topics\n\n- Business requirements", entry)
        self.assertIn("[requirement](#requirement)", entry)
        self.assertNotIn("Software requirements", entry)
        self.assertNotIn("[use case]", entry)
        self.assertNotIn("### Example", entry)

    def test_embedded_values_and_ordered_references_stay_distinct(self):
        scenario = self.slot("UseCase", "main_success_scenario")
        self.assertFalse(scenario["inlined"])
        self.assertEqual(scenario["reference_ranges"], ["Scenario"])
        content = self.slot("Scenario", "content")
        self.assertTrue(content["inlined"])
        self.assertEqual(content["inline_ranges"], ["RichContent"])
        self.assertEqual(content["reference_ranges"], [])
        steps = self.slot("Scenario", "steps")
        self.assertTrue(steps["multivalued"])
        self.assertTrue(steps["ordered"])
        self.assertEqual(steps["reference_ranges"], ["ScenarioStep"])

    def test_annotations_rules_and_discriminators_remain_available(self):
        classes = self.catalog["classes"]
        self.assertEqual(classes["Job"]["annotations"]["owned_memberships"], "Job.steps or Job.outcomes")
        self.assertEqual(self.slot("UseCase", "kind")["equals_string"], "UseCase")
        self.assertTrue(self.slot("QualityAttribute", "authored_id")["required"])
        self.assertFalse(self.slot("Job", "authored_id")["required"])
        self.assertTrue(self.slot("RichContent", "spans")["ordered"])
        self.assertTrue(self.slot("RichContent", "spans")["required"])
        self.assertIn("unresolved", self.slot("Resolution", "status")["enum_values"])

    def test_stale_canonical_membership_paths_cannot_generate(self):
        catalog = copy.deepcopy(self.catalog)
        catalog["classes"]["Job"]["annotations"]["owned_memberships"] = "Job.missing"
        with self.assertRaisesRegex(ValueError, "absent field Job.missing"):
            generate.validate_annotations(catalog)
        catalog = copy.deepcopy(self.catalog)
        catalog["classes"]["Concept"]["slots"]["actions"]["annotations"]["relationship_category"] = "semantic-association"
        with self.assertRaisesRegex(ValueError, "membership must declare containment"):
            generate.validate_annotations(catalog)

    def test_relationship_defaults_complete_the_catalog(self):
        raw = yaml.safe_load(generate.MODEL.read_bytes())
        # The source states only the category for an ordinary citation field.
        source = raw["classes"]["ContextualReferences"]["attributes"]["references"]["annotations"]
        self.assertEqual(source, {"relationship_category": "contextual-citation"})
        completed = self.slot("ContextualReferences", "references")["annotations"]
        defaults = self.catalog["relationship_categories"]["contextual-citation"]["defaults"]
        for key, value in defaults.items():
            self.assertEqual(completed[key], value)
        self.assertEqual(completed["relationship_inverse"], "none")
        self.assertEqual(completed["relationship_label"], "references")
        self.assertEqual(completed["relationship_ordering"], "authored presentation order; no execution order implied")
        self.assertEqual(self.slot("UseCase", "subject_system")["annotations"]["relationship_ordering"], "single target")
        self.assertEqual(self.slot("QualityAttribute", "subattribute_of")["annotations"]["relationship_ordering"], "unordered")
        # A stated value wins over the category default.
        self.assertEqual(self.slot("Requirement", "refines")["annotations"]["relationship_cycles"], "forbidden")
        self.assertEqual(self.slot("Goal", "owner")["annotations"]["relationship_inverse"], "System.goals")
        broken = copy.deepcopy(raw)
        del broken["enums"]["RelationshipCategory"]
        with self.assertRaisesRegex(ValueError, "RelationshipCategory"):
            generate.relationship_categories(broken)

    def test_projection_paths_are_derived_from_relationship_fields(self):
        projections = self.catalog["projections"]
        self.assertEqual(projections, sorted(set(projections)))
        for path in ("TraceableElement.supports", "Requirement.refines", "BusinessRule.refines", "System.goals",
                     "Job.steps", "Capability.provided_by",
                     "System.enclosing_system", "SystemContextSpecification.specification_system"):
            self.assertIn(path, projections)
        # Owner backlinks and provenance are not separate relationship claims.
        for path in ("Goal.owner", "DomainElement.source_document"):
            self.assertNotIn(path, projections)

    def test_every_induced_slot_is_present(self):
        view = SchemaView(str(generate.MODEL))
        for name in view.all_classes():
            expected = {str(slot.name) for slot in view.class_induced_slots(name)}
            self.assertEqual(set(self.catalog["classes"][name]["slots"]), expected, name)


    def test_relationship_contracts_are_complete_and_required(self):
        from relationships import validate_relationships, relationships_markdown
        validate_relationships(self.catalog)
        reference = relationships_markdown(self.catalog).decode()
        self.assertIn('ConceptUse.definition', reference)
        self.assertIn('ContextualReferences.references', reference)
        self.assertIn('TraceableElement.supports', reference)
        for mutate, message in [
            (lambda c: c['classes']['ConceptUse']['slots']['definition']['annotations'].pop('relationship_scope'), 'missing relationship metadata'),
            (lambda c: c['classes']['ConceptUse']['slots']['definition']['annotations'].update(relationship_cycles='guess'), 'cycle policy'),
            (lambda c: c['classes']['ConceptUse']['slots']['definition']['annotations'].update(relationship_category='guess'), 'invalid relationship category'),
            (lambda c: c['classes']['UseCase']['slots']['involves'].update(reference_ranges=[]), 'multivalued link'),
            (lambda c: c['classes']['Requirement']['slots']['refines']['annotations'].pop('occurrence_role'), 'multivalued link'),
            (lambda c: c['classes']['UseCase']['slots']['supports']['annotations'].update(trace_heading='Involves'), 'unique within a class'),
        ]:
            catalog = copy.deepcopy(self.catalog)
            mutate(catalog)
            with self.assertRaisesRegex(ValueError, message):
                validate_relationships(catalog)


class GenerationPortabilityTests(unittest.TestCase):
    def test_console_callable_and_windows_stdout_normalize_to_committed_bytes(self):
        schema_path = generate.ROOT / "model/generated/schema.json"
        expected = schema_path.read_bytes()
        with mock.patch.object(generate.subprocess, "check_output", return_value=expected.replace(b"\n", b"\r\n")) as child:
            outputs = generate.generate_outputs()
        command = child.call_args.args[0]
        self.assertEqual(command[:3], [generate.sys.executable, "-c", "from linkml.generators.jsonschemagen import cli; cli()"])
        self.assertEqual(outputs[schema_path], expected)


class CheckModeTests(unittest.TestCase):
    def test_stale_and_missing_outputs_are_reported_without_writes(self):
        with tempfile.TemporaryDirectory() as directory:
            stale = Path(directory) / "existing.json"
            missing = Path(directory) / "missing.json"
            stale.write_bytes(b"old\n")
            outputs = {stale: b"new\n", missing: b"generated\n"}
            self.assertEqual(generate.sync_outputs(outputs, check=True), [stale, missing])
            self.assertEqual(stale.read_bytes(), b"old\n")
            self.assertFalse(missing.exists())
            self.assertEqual(generate.sync_outputs(outputs, check=False), [])
            self.assertEqual(generate.sync_outputs(outputs, check=True), [])



class ModelViewTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        data = generate.MODEL.read_bytes()
        cls.catalog = generate.make_catalog(SchemaView(str(generate.MODEL)), yaml.safe_load(data), data)
        cls.view = yaml.safe_load((generate.ROOT / 'model/views/core-domain.yaml').read_text())

    def test_groups_cover_every_class_and_keep_use_cases_in_software(self):
        from model_views import validate_groups
        validate_groups(self.catalog)
        self.assertEqual(len(self.catalog['subsets']), 8)
        self.assertEqual(self.catalog['classes']['UseCase']['subsets'], ['software-requirements'])
        self.assertEqual(self.catalog['classes']['BusinessRule']['subsets'], ['business-requirements'])
        self.assertEqual(self.catalog['classes']['Job']['subsets'], ['user-requirements'])
        self.assertEqual(set(self.catalog['classes']['System']['subsets']), {'business-requirements', 'software-requirements'})

    def test_invalid_grouping_cannot_generate(self):
        from model_views import validate_groups
        for mutation, message in [
            (lambda c: c['subsets']['business-requirements'].update(classification='technical'), 'invalid classification'),
            (lambda c: c['classes']['Job'].update(subsets=[]), 'unclassified class'),
            (lambda c: c['classes']['Job'].update(subsets=['typo']), 'unknown subset'),
        ]:
            c = copy.deepcopy(self.catalog)
            mutation(c)
            with self.assertRaisesRegex(ValueError, message):
                validate_groups(c)

    def test_generated_document_kind_cannot_shadow_domain_kind(self):
        from model_views import validate_groups
        for kind in ('use-case', 'component', 'concept-action'):
            c = copy.deepcopy(self.catalog)
            c['classes']['ExampleDocument'] = {
                'subsets': ['specification-documentation'], 'abstract': False,
                'ancestors': ['GeneratedDocument'], 'slots': {'type': {'equals_string': kind}},
            }
            with self.assertRaisesRegex(ValueError, 'overlap'):
                validate_groups(c)

    def test_unknown_slot_subset_is_rejected(self):
        from model_views import validate_declared_memberships
        raw = yaml.safe_load(generate.MODEL.read_bytes())
        raw['classes']['Job']['attributes']['statement']['in_subset'] = ['missing']
        with self.assertRaisesRegex(ValueError, 'unknown subset missing'):
            validate_declared_memberships(raw)

    def test_invalid_view_selectors_and_edges_cannot_generate(self):
        from model_views import validate_view
        for mutation, message in [
            (lambda v: v['subsets'].append('unknown'), 'unknown subset'),
            (lambda v: v['subsets'].append('document-content'), 'is not core'),
            (lambda v: v['types'].update(Missing=[]), 'unknown class'),
            (lambda v: v['types']['Job'].append('missing'), 'absent field'),
            (lambda v: v['diagrams'][0]['paths'].append('UseCase.involves'), 'all declared class targets'),
            (lambda v: v['diagrams'][0].update(edges=[['Job', 'UseCase']]), 'diagram must declare'),
        ]:
            v = copy.deepcopy(self.view)
            mutation(v)
            with self.assertRaisesRegex(ValueError, message):
                validate_view(self.catalog, v)

    def test_diagram_edges_follow_model_endpoint_changes(self):
        from model_views import core_markdown
        c = copy.deepcopy(self.catalog)
        c['classes']['BusinessRule']['slots']['specification_system']['reference_ranges'] = ['Concept']
        v = copy.deepcopy(self.view)
        diagram = next(d for d in v['diagrams'] if 'BusinessRule.specification_system' in d['paths'])
        diagram['show_cardinality'] = True
        output = core_markdown(c, v).decode()
        self.assertIn('BusinessRule --> "0..1" Concept : specification system', output)
        self.assertNotIn('BusinessRule --> "0..1" System : specification system', output)
        self.assertIn('SystemContextSpecification <|-- Requirement', output)
        self.assertIn('SystemContextSpecification <|-- Component', output)

if __name__ == "__main__":
    unittest.main()
