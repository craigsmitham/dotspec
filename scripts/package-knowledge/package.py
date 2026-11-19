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
        return '](' + f'https://github.com/craigsmitham/dotspec/blob/{revision}/{quote(path)}' + (
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


def build(version, target):
    revision = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    manifest = json.loads(PACKAGE.read_text())
    if version:
        manifest["version"] = version
    contents = {"knowledge.json": json.dumps(manifest, indent=2) + "\n"}
    catalog = json.loads((ROOT / "docs/handbook/topics.json").read_text())
    pages = {}
    for entry in catalog:
        pages.setdefault(entry["page"].split("#")[0], entry["description"])
    for file in sorted((ROOT / "examples/todo").rglob("*.spec.md")):
        pages[file.relative_to(ROOT).as_posix()] = "Worked specification from the To-do app."

    index = ['---', 'okf_version: "0.2"', '---', '', '# dotspec handbook', '',
             'Start with the authoring handbook or open a topic for the decision at hand.', '',
             '## Guides', '']
    for entry in catalog:
        index.append(f'- [{entry["topic"]}]({entry["page"]}) — {entry["description"]}')
    index.extend(['', '## Worked examples', ''])
    for page, description in pages.items():
        body = (ROOT / page).read_text()
        example = page.startswith('examples/')
        if example:
            # Native spec metadata belongs to the original example. AXM uses OKF frontmatter.
            _, _, body = body.partition('\n---\n')
            title = re.search(r'^# (.+)$', body, re.MULTILINE)[1]
            index.append(f'- [{title}]({page}) — {description}')
        tags = ['dotspec', 'example' if example else 'authoring', Path(page).stem.split('.')[0].lower()]
        metadata = ('---\ntype: ' + ('example' if example else 'guide') +
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
