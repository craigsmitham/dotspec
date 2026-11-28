"""Build the optional AXM Knowledge archive from the reviewed handbook sources."""
import argparse
import json
import posixpath
import re
import subprocess
from pathlib import Path, PurePosixPath
from tempfile import NamedTemporaryFile
from urllib.parse import quote, urlsplit
from zipfile import ZIP_DEFLATED, ZipFile, ZipInfo


ROOT = Path(__file__).resolve().parents[2]
HANDBOOK = PurePosixPath("docs/cli-handbook")
PACKAGE = Path(__file__).with_name("knowledge.json")


def bundle_links(markdown, page, included, revision, contents):
    """Preserve local links between packaged pages and pin other sources to Git."""
    def destination(match):
        dest = match[1]
        url = urlsplit(dest.strip("<>"))
        if url.scheme or url.netloc or not url.path:
            return match[0]
        path = posixpath.normpath(posixpath.join(posixpath.dirname(page), url.path))
        if path in included:
            return match[0]
        return '](' + f'https://github.com/dotspec-md/dotspec/blob/{revision}/{quote(path)}' + (
            '#' + url.fragment if url.fragment else '') + ')'

    # AXM indexes bracket links even inside fenced examples. Keep those examples
    # as verbatim text resources so project-relative specimen links stay literal.
    lines = []
    fence = ""
    sample = []
    for number, line in enumerate(markdown.splitlines(keepends=True), 1):
        marker = re.match(r'^ {0,3}(`{3,}|~{3,})(.*)$', line)
        if fence:
            if marker and marker[1].startswith(fence) and not marker[2].strip():
                body = ''.join(sample[1:])
                if '](' in body:
                    asset = f'samples/{PurePosixPath(page).stem}-{number}.txt'
                    contents['src/' + posixpath.join(posixpath.dirname(page), asset)] = body
                    lines.append(f'[Read the literal example]({asset})\n')
                else:
                    lines.extend(sample)
                    lines.append(line)
                fence = ""
            else:
                sample.append(line)
        elif marker:
            fence = marker[1]
            sample = [line]
        else:
            lines.append(re.sub(r'\]\((<[^>]+>|[^()\s]+)\)', destination, line))
    if fence:
        raise ValueError(f"Unclosed fenced example in {page}")
    return ''.join(lines)


def page_form(page):
    path = PurePosixPath(page)
    root = path.parts[0]
    if root in ('examples', 'specs'):
        return 'example'
    if root == 'model' or path.is_relative_to('docs/reference'):
        return 'reference'
    return {'explanations': 'explanation'}.get(path.relative_to(HANDBOOK).parts[0], 'guide')


def publication_pages(catalog):
    public = ROOT / HANDBOOK
    pages = {file.relative_to(ROOT).as_posix(): file.stem.replace('-', ' ')
             for file in sorted(public.glob('*/**/*.md'))}
    pages[str(HANDBOOK / 'README.md')] = 'Choose a specification type to understand or write.'
    for destination in json.loads((ROOT / 'docs/reference/diagnostics.json').read_text()):
        page = destination['page']
        pages[page] = (ROOT / page).read_text().splitlines()[0].removeprefix('# ')
    pages['model/markdown-mapping.md'] = "Markdown mapping profile."
    pages['examples/todo/README.md'] = "Walk through the To-do application."
    for folder in ('examples/todo/specs', 'examples/todo/app/specs'):
        for file in sorted((ROOT / folder).rglob('*.spec.md')):
            pages[file.relative_to(ROOT).as_posix()] = "Worked specification from the To-do app."

        for file in sorted((ROOT / folder).rglob('*.quality-model.md')):
            pages[file.relative_to(ROOT).as_posix()] = "Reusable quality model from the To-do app."

    # Other public specifications are copied when an editor selects them for a topic.
    selectable = {file.relative_to(ROOT).as_posix() for file in (ROOT / 'specs').rglob('*.spec.md')}
    selectable.update(file.relative_to(ROOT).as_posix()
                      for file in (ROOT / 'examples/todo/specification').rglob('*.md'))
    for entry in catalog:
        for destination in entry['pages']:
            page = destination['page'].partition('#')[0]
            if page not in pages and page not in selectable:
                raise ValueError(f"Catalog destination outside the published handbook and examples: {page}")
            pages[page] = entry['description']
    return pages


def build(version, target):
    revision = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    manifest = json.loads(PACKAGE.read_text())
    if version:
        manifest["version"] = version
    contents = {"knowledge.json": json.dumps(manifest, indent=2) + "\n"}
    public = ROOT / HANDBOOK
    catalog = json.loads((public / 'topics.json').read_text())
    pages = publication_pages(catalog)

    # The bundle index lives at src/, the repository navigation at the handbook root.
    navigation = (public / 'README.md').read_text().split("<!-- Generated topic navigation -->", 1)[1]
    navigation = re.sub(r'\]\(([^)]+)\)', lambda m: '](' + posixpath.normpath(
        posixpath.join(str(HANDBOOK), m[1])) + ')', navigation)
    index = ['---', 'okf_version: "0.2"', '---', '', '# dotspec handbook', '',
             'For offline CLI reading instructions, see the [handbook introduction](docs/cli-handbook/README.md).',
             '', navigation.strip()]
    for page, description in pages.items():
        body = (ROOT / page).read_text()
        if page.endswith('.spec.md'):
            # Native spec metadata belongs to the original example. AXM uses OKF frontmatter.
            _, _, body = body.partition('\n---\n')
        form = page_form(page)
        tags = ['dotspec', form, Path(page).stem.split('.')[0].lower()]
        metadata = ('---\ntype: ' + form +
                    '\ndescription: ' + json.dumps(description) +
                    '\ntags: ' + json.dumps(tags) + '\n---\n\n')
        contents["src/" + page] = metadata + bundle_links(body.lstrip(), page, pages, revision, contents)
    contents["src/index.md"] = '\n'.join(index) + '\n'

    target.parent.mkdir(parents=True, exist_ok=True)
    with NamedTemporaryFile(dir=target.parent, delete=False) as output:
        temporary = Path(output.name)
    try:
        with ZipFile(temporary, 'w') as archive:
            for path, content in sorted(contents.items()):
                archive.writestr(ZipInfo('dotspec/' + path), content.encode(), compress_type=ZIP_DEFLATED)
        temporary.replace(target)
    finally:
        temporary.unlink(missing_ok=True)
    print(target)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--version', help='release version for the bundle manifest')
    parser.add_argument('--output', type=Path, default=ROOT / "dist/dotspec-knowledge.zip")
    args = parser.parse_args()
    build(args.version, args.output)
