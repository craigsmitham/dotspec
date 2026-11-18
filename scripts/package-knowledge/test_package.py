import json
import tempfile
import unittest
from pathlib import Path
from zipfile import ZipFile

from package import ROOT, build


class KnowledgePackageTest(unittest.TestCase):
    def test_package_uses_reviewed_sources_and_preserves_example_links(self):
        with tempfile.TemporaryDirectory() as folder:
            target = Path(folder) / 'dotspec.zip'
            build('1.2.3', target)
            first = target.read_bytes()
            build('1.2.3', target)
            self.assertEqual(first, target.read_bytes())
            with ZipFile(target) as archive:
                self.assertEqual(json.loads(archive.read('dotspec/knowledge.json'))['version'], '1.2.3')
                root = 'dotspec/src/'
                for topic in json.loads((ROOT / 'docs/handbook/topics.json').read_text()):
                    self.assertIn(root + topic['page'].split('#')[0], archive.namelist())
                source = 'docs/handbook/project-conventions.md'
                sample = (ROOT / source).read_text().split('```markdown\n')[1].split('```')[0]
                matches = [name for name in archive.namelist()
                           if name.startswith(root + 'docs/handbook/samples/project-conventions-')]
                self.assertEqual(len(matches), 1)
                self.assertEqual(archive.read(matches[0]).decode(), sample)
                guide = archive.read(root + source).decode()
                self.assertIn(matches[0].split('handbook/')[1], guide)
                example = 'examples/requirements-review/specs/requirements/recovery-feedback.requirement.spec.md'
                self.assertIn('## Criterion', archive.read(root + example).decode())
                todo = 'examples/todo/specs/requirements/restart.requirement.spec.md'
                self.assertIn('Keep recorded work available', archive.read(root + todo).decode())
