import json
import tempfile
import unittest
from pathlib import Path
from zipfile import ZipFile

from package import ROOT, build, publication_pages


class KnowledgePackageTest(unittest.TestCase):
    def test_catalog_cannot_publish_repository_documents(self):
        for path in ('docs/cli-handbook/AGENTS.md', 'docs/maintainers/releasing.md',
                     'docs/cli-handbook/../maintainers/handbook.md', 'CONTRIBUTING.md'):
            with self.subTest(path=path), self.assertRaisesRegex(ValueError, 'outside'):
                publication_pages([{'description': 'Sample topic', 'pages': [{'page': path}]}])

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
                index = archive.read(root + 'index.md').decode()
                for topic in json.loads((ROOT / 'docs/cli-handbook/topics.json').read_text()):
                    self.assertIn(f'### {topic["title"]}', index)
                    for destination in topic['pages']:
                        page = destination['page']
                        self.assertIn(f']({page})', index)
                        self.assertEqual(archive.namelist().count(root + page.split('#')[0]), 1)
                for page in (ROOT / 'docs/cli-handbook').glob('*/**/*.md'):
                    self.assertIn(root + page.relative_to(ROOT).as_posix(), archive.namelist())
                for name in archive.namelist():
                    self.assertNotIn('/maintainers/', name)
                    self.assertNotIn('/plans/', name)
                    self.assertNotEqual(Path(name).name, 'AGENTS.md')
                source = 'docs/cli-handbook/guides/project-conventions.md'
                sample = (ROOT / source).read_text().split('```markdown\n')[1].split('```')[0]
                matches = [name for name in archive.namelist()
                           if name.startswith(root + 'docs/cli-handbook/guides/samples/project-conventions-')]
                self.assertEqual(len(matches), 1)
                self.assertEqual(archive.read(matches[0]).decode(), sample)
                guide = archive.read(root + source).decode()
                self.assertIn(matches[0].split('guides/')[1], guide)
                todo = 'examples/todo/specs/requirements/restart.requirement.spec.md'
                self.assertIn('Keep recorded work available', archive.read(root + todo).decode())
