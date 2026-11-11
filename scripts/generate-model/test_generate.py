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
        catalog["classes"]["HasJobStep"]["annotations"]["projection_of"] = "Concept.actions.actions"
        with self.assertRaisesRegex(ValueError, "projection violates declared assertion endpoints"):
            generate.validate_annotations(catalog)

    def test_every_induced_slot_is_present(self):
        view = SchemaView(str(generate.MODEL))
        for name in view.all_classes():
            expected = {str(slot.name) for slot in view.class_induced_slots(name)}
            self.assertEqual(set(self.catalog["classes"][name]["slots"]), expected, name)


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


if __name__ == "__main__":
    unittest.main()
