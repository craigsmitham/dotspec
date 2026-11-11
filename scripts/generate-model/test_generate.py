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
        # Abstract Requirement must accept either concrete requirement subtype.
        subject = self.slot("QualityClassification", "subject")
        self.assertEqual(subject["reference_ranges"], ["Requirement"])
        self.assertTrue(subject["required"])
        self.assertIn("Requirement", self.catalog["classes"]["QualityRequirement"]["ancestors"])
        self.assertIn("Requirement", self.catalog["classes"]["GeneralRequirement"]["ancestors"])
        # A union must not silently reopen the broad DomainElement parent range.
        owner = self.slot("ScenarioStep", "owner")
        self.assertEqual(owner["range"], "DomainElement")
        self.assertEqual(owner["reference_ranges"], ["Scenario", "Extension"])
        self.assertNotIn("DomainElement", owner["reference_ranges"])
        self.assertEqual(self.slot("QualityAttribute", "owner")["reference_ranges"], ["System", "QualityModel"])
        self.assertEqual(self.slot("DeclaresQualityAttribute", "subject")["reference_ranges"], ["System", "QualityModel"])
        self.assertEqual(self.slot("LabelOccurrence", "quality_model_context")["reference_ranges"], ["QualityModel"])

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
        self.assertEqual(classes["HasJobStep"]["annotations"]["projection_of"], "Job.steps.steps")
        self.assertEqual(self.slot("UseCase", "kind")["equals_string"], "UseCase")
        self.assertTrue(self.slot("QualityAttribute", "authored_id")["required"])
        self.assertFalse(self.slot("Job", "authored_id")["required"])
        self.assertEqual(len(classes["RichContent"]["rules"]), 2)
        self.assertIn("json_schema_projection_gap", classes["RichContent"]["annotations"])
        self.assertIn("unresolved", self.slot("Resolution", "status")["enum_values"])

    def test_stale_canonical_membership_paths_cannot_generate(self):
        catalog = copy.deepcopy(self.catalog)
        catalog["classes"]["Job"]["annotations"]["owned_memberships"] = "Job.steps.missing"
        with self.assertRaisesRegex(ValueError, "absent field JobMap.missing"):
            generate.validate_annotations(catalog)
        catalog = copy.deepcopy(self.catalog)
        catalog["classes"]["ActionCollection"]["slots"]["actions"]["annotations"]["relationship_category"] = "semantic-association"
        with self.assertRaisesRegex(ValueError, "membership must declare containment"):
            generate.validate_annotations(catalog)
        catalog = copy.deepcopy(self.catalog)
        catalog["classes"]["HasJobStep"]["annotations"]["projection_of"] = "Concept.actions.actions"
        with self.assertRaisesRegex(ValueError, "projection violates declared assertion endpoints"):
            generate.validate_annotations(catalog)

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
        self.assertIn('projection_path', reference)
        self.assertNotIn('system_context', self.catalog['classes']['System']['slots'])
        self.assertNotIn('system_context', self.catalog['classes']['QualityModel']['slots'])
        for mutate, message in [
            (lambda c: c['classes']['UsesConcept']['slots']['subject'].update(required=False), 'endpoints must be required'),
            (lambda c: c['classes']['ConceptUse']['slots']['definition']['annotations'].pop('relationship_scope'), 'missing relationship metadata'),
            (lambda c: c['classes']['ConceptUse']['slots']['bindings']['annotations'].update(paired_with='missing'), 'reciprocal peer'),
            (lambda c: c['classes']['ConceptUse']['slots']['definition']['annotations'].update(relationship_cycles='guess'), 'cycle policy'),
            (lambda c: c['classes']['SupportsJob']['annotations'].pop('projection_of'), 'canonical projection fields'),
            (lambda c: c['classes']['UseCase']['slots']['involves']['annotations'].update(trace_targets='Concept'), 'canonical endpoint ranges'),
            (lambda c: c['classes']['Requirement']['slots']['refines']['annotations'].pop('trace_heading'), 'paired and unique'),
            (lambda c: c['classes']['UseCase']['slots']['supports']['annotations'].update(trace_heading='Involves'), 'paired and unique'),
            (lambda c: c['classes']['UseCase']['slots']['involves_references'].update(range='Record'), 'source-reference occurrences'),
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
        self.assertEqual(len(self.catalog['subsets']), 7)
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
            (lambda v: v['diagrams'][0]['paths'].append('UseCase.related_jobs'), 'all declared class targets'),
            (lambda v: v['diagrams'][0].update(edges=[['Job', 'UseCase']]), 'diagram must declare'),
        ]:
            v = copy.deepcopy(self.view)
            mutation(v)
            with self.assertRaisesRegex(ValueError, message):
                validate_view(self.catalog, v)

    def test_diagram_edges_follow_model_endpoint_changes(self):
        from model_views import core_markdown
        c = copy.deepcopy(self.catalog)
        c['classes']['BusinessRule']['slots']['system_context']['reference_ranges'] = ['Concept']
        v = copy.deepcopy(self.view)
        v['diagrams'][0]['show_cardinality'] = True
        output = core_markdown(c, v).decode()
        self.assertIn('BusinessRule ..> "0..1" Concept : system_context', output)
        self.assertNotIn('BusinessRule ..> "0..1" System : system_context', output)
        self.assertIn('Requirement <|-- QualityRequirement', output)
        self.assertNotIn('enforces -->', output)

if __name__ == "__main__":
    unittest.main()
