"""Locate concrete witnesses for the public relationship and authoring references."""
import json
from pathlib import Path
import re
import subprocess
import sys


def record_witnesses(snapshot, classes, expected, evidence, project, label):
    # Locations and IDs in a snapshot belong to that export. Resolve them before
    # reading the next project, retaining a path a reader can actually inspect.
    sources = {s["id"]: Path(s["locator"]) for s in snapshot["sources"]}
    documents = {d["id"]: sources[d["source"]] for d in snapshot["documents"]}

    def context(path):
        if path.is_relative_to(project):
            return f"{label}/{path.relative_to(project).as_posix()}"
        return path.relative_to(Path.cwd()).as_posix() if path.is_relative_to(Path.cwd()) else str(path)

    def visit(record, kind, location):
        concrete = record.get("kind")
        if concrete in classes and kind in classes[concrete]["ancestors"]:
            kind = concrete
        definition = classes[kind]
        source = record.get("source_document") or record.get("source")
        if isinstance(source, str) and source in sources:
            location = context(sources[source])
        document = record.get("document")
        if isinstance(document, str) and document in documents:
            location = context(documents[document])
        for field, slot in definition["slots"].items():
            value = record.get(field)
            if value is None or value == [] or value == {} or value == "":
                continue
            if isinstance(value, dict) and value.get("presence") in {"empty", "absent"}:
                continue
            for ancestor in definition["ancestors"]:
                key = f"{ancestor}.{field}"
                if key in expected:
                    evidence.setdefault(key, location)
            if slot["inline_ranges"]:
                for item in value if isinstance(value, list) else [value]:
                    if isinstance(item, dict):
                        visit(item, slot["range"], location)

    visit(snapshot, "Snapshot", f"model export: {label}")


def document_witnesses(project, classes, evidence):
    names = {cls["slots"]["type"]["equals_string"]: name for name, cls in classes.items()
             if not cls["abstract"] and "SubjectSpecReport" in cls["ancestors"]}

    def witness(kind, path):
        for ancestor in classes[names[kind]]["ancestors"]:
            if "documents" in classes[ancestor]["slots"]:
                evidence.setdefault(f"{ancestor}.documents", path.relative_to(Path.cwd()).as_posix())

    for path in [project / "SPEC.md", *sorted((project / "specification").rglob("*.md"))]:
        if not path.is_file():
            continue
        text = path.read_text(encoding="utf-8")
        header = re.match(r"^---\n(.*?)\n---\n", text, re.S)
        if not header or "generated_by: dotspec" not in header[1].splitlines():
            continue
        kind = re.search(r"^type: ([\w-]+)$", header[1], re.M)
        if kind and kind[1] in names and re.search(r'^source: "[^"]+#[^"]+"$', header[1], re.M):
            witness(kind[1], path)
        for raw in re.findall(r"^<!-- dotspec-doc (.+) -->$", text, re.M):
            doc = json.loads(raw)
            if doc.get("documents"):
                witness(doc["type"], path)


def main(binary, scratch, detailed):
    root = Path.cwd()
    classes = json.loads((root / "model/generated/catalog.json").read_bytes())["classes"]
    reference = (root / "docs/reference/model/relationships.md").read_text(encoding="utf-8")
    relationships = set(re.findall(r"^### (\w+\.\w+)$", reference, re.M))
    # Include rich-content roles and embedded criterion, action and scenario
    # forms alongside the reference's canonical relationship contracts.
    forms = set()
    for kind, definition in classes.items():
        for field, slot in definition["slots"].items():
            if not slot["inline_ranges"] or slot["range"] in {"SourceLocation", "ByteSpan"}:
                continue
            if any(field in classes[parent]["slots"] for parent in definition["ancestors"][1:]):
                continue
            if "DomainElement" in definition["ancestors"] or "Record" not in definition["ancestors"]:
                forms.add(f"{kind}.{field}")
    evidence = {}
    examples = root / "examples"
    todo = examples / "todo"
    projects = [todo, *sorted(p.parent for p in examples.glob("*/dotspec.yaml") if p.parent != todo)]
    for project in projects:
        document_witnesses(project, classes, evidence)
        raw = subprocess.check_output([binary, "model", "export", "--path", project])
        snapshot = json.loads(raw)
        record_witnesses(snapshot, classes, relationships | forms, evidence, project,
                         project.relative_to(root).as_posix())
    record_witnesses(json.loads(scratch.read_bytes()), classes, relationships | forms, evidence,
                     scratch.parent / "author", "scratch authoring and diagnostics")
    assignments = (todo / "COVERAGE.md").read_text(encoding="utf-8").split("### Gaps in the exported examples", 1)[1]
    allowed = set(re.findall(r"`(\w+\.\w+)`", assignments))
    missing = (relationships | forms) - evidence.keys()
    assert not missing - allowed, f"example fields need evidence or review assignments: {sorted(missing - allowed)}"
    assert not allowed - missing, f"retire obsolete field gap assignments: {sorted(allowed - missing)}"
    for title, keys in [("Relationship fields", relationships), ("Structured content fields", forms)]:
        gaps = sorted(keys - evidence.keys())
        print(f"{title}: {len(keys) - len(gaps)}/{len(keys)} populated in example artifacts.")
        if detailed:
            for key in sorted(keys):
                print(f"  {key:<50} {evidence.get(key, 'GAP: no populated example')}")


if __name__ == "__main__":
    main(sys.argv[1], Path(sys.argv[2]), "--details" in sys.argv[3:])
