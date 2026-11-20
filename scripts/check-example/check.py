"""Check the runnable reference through the public CLI and the app's forms."""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[2]
TODO = ROOT / "examples/todo"
REQUIREMENT = Path("specs/requirements/restart.requirement.spec.md")
ACCEPTANCE = Path("tests/acceptance.py")
TYPES = {
    "system": ("", "Personal task list"),
    "feature": ("pick-up-where-i-left-off", "Pick up where I left off"),
    "capability": ("local-task-storage", "Local task storage"),
    "business-rule": ("completion", "Completion preserves a commitment"),
    "job": ("resume", "Resume interrupted work"),
    "use-case": ("resume-list", "Return to my list"),
    "concept": ("checklist", "Checklist"),
    "requirement": ("restart", "Preserve tasks across restart"),
    "quality-model": ("personal-data", "Personal data quality"),
}
VIEWS = {
    "overview": "Task list forms",
    "pegs": "A finished commitment stays",
    "trace": "| refines | [Keep recorded work available]",
    "needs": "DesiredOutcome: Less reconstruction",
    "review": "illustrates [Reopen]",
    "docs": "| Reopen | unchanged | unchanged | open |",
}


def run(args, cwd=ROOT, status=0):
    result = subprocess.run(args, cwd=cwd, capture_output=True, text=True, encoding="utf-8")
    if result.returncode != status:
        raise RuntimeError(f"{args} returned {result.returncode}, expected {status}\n"
                           f"{result.stderr}\n{result.stdout[:2500]}")
    return result.stdout


class CLI:
    def __init__(self, binary):
        self.binary = binary
        self.exercised = set()

    def __call__(self, command, *args, cwd=ROOT, status=0):
        self.exercised.add(command)
        return run([self.binary, *command.split(), *args], cwd, status)

    def public_commands(self):
        # Discover using the built binary so a stale generated reference cannot hide a command.
        def walk(command):
            help_text = run([self.binary, *command.split(), "--help"])
            section = re.search(r"^Available Commands:\n((?:  .+\n)+)", help_text, re.M)
            yield command
            if section:
                for child in re.findall(r"^  (\S+)", section[1], re.M):
                    yield from walk((command + " " + child).strip())
        return set(walk(""))


def project_copy(tmp, name):
    # Copy only reference sources. A developer's local database stays in place.
    return shutil.copytree(TODO, tmp / name, ignore=shutil.ignore_patterns("*.sqlite*", "__pycache__"))


def check_project(cli, project, tmp):
    def report(command, *args):
        return json.loads(cli(command, *args, "--path", project, "--json"))

    checked = report("check")
    assert checked["complete"] and checked["scope"]["project"]["name"] == "todo"
    sources = report("sources")
    assert sources["complete"] and Path(sources["scope"]["path"]) == project.resolve()
    files = {Path(f["path"]): f for f in sources["files"]}
    referring = files.get(project / ACCEPTANCE, {}).get("referenced_by", [])
    assert any(Path(link["path"]) == project / REQUIREMENT for link in referring), (
        "restart requirement must link to tests/acceptance.py")

    snapshot = tmp / "snapshot.json"
    snapshot.write_text(cli("model export", "--path", project), encoding="utf-8")
    assert json.loads(cli("model validate", snapshot, "--json"))["valid"]
    model = json.loads(snapshot.read_bytes())
    elements = model["elements"]

    def one(kind, **fields):
        value, = [e for e in elements if e["kind"] == kind and
                  all(e.get(k) == v for k, v in fields.items())]
        return value

    requirement = one("Requirement", slug="restart")
    use = one("UseCase", slug="resume-list")
    job = one("Job", slug="resume")
    clause = one("ObligationClause", owner=use["id"], authored_id="keep-recorded-work-available")
    assert clause["id"] in requirement.get("refines", []), "missing restart refinement of keep-recorded-work-available"
    assert one("BusinessRule", slug="completion")["id"] in clause["derived_from"]
    assert job["id"] in use["supports"]
    assert one("DesiredOutcome", owner=job["id"], authored_id="less-reconstruction")["id"] in requirement["supports"]
    for kind in ("System", "Component", "Feature", "Capability", "ConceptUse"):
        subject = one(kind)
        description = subject["description"]
        assert description["presence"] == "authored" and description["content_nodes"], kind
        bound = {b["content_node"] for b in model["content_bindings"]
                 if b["owner"] == subject["id"] and b["field_path"] == "description"}
        assert set(description["content_nodes"]) <= bound, kind
    system = one("System")
    assert system["id"] == use["subject_system"]
    assert one("SystemQualityAttribute", owner=system["id"], authored_id="durability")["id"] in requirement["quality_attributes"]
    assert one("Concept", slug="checklist")["id"] == one("ConceptUse")["definition"]
    assert one("QualityModel", slug="personal-data")["quality_attributes"]
    for view, marker in VIEWS.items():
        assert marker in cli("model " + view, "--path", project), f"missing {marker} in {view}"


def markdown_links(page):
    body = page.read_text(encoding="utf-8")
    links = []
    definitions = dict(re.findall(r"^\[([^]\n]+)\]: <([^>\n]+)>$", body, re.M))
    for match in re.finditer(r"\[([^]\n]+)\](?:\((<[^>\n]+>|[^\s)]+)\)|\[([^]\n]+)\])", body):
        dest = match[2] or definitions.get(match[3])
        if dest is None:
            continue
        url = urlsplit(dest.strip("<>"))
        if not url.scheme and not url.netloc:
            target = (page.parent / unquote(url.path)).resolve() if url.path else page
            links.append((match[1], target, unquote(url.fragment), match[0]))
    return links


def follow(page, name):
    matches = {link[1] for link in markdown_links(page) if link[0].casefold() == name.casefold()}
    assert len(matches) == 1, f"{page}: expected one destination for {name!r}, got {matches}"
    target, = matches
    assert target.is_file(), target
    return target


def spec_journey(cli, tmp):
    project = project_copy(tmp, "reader")
    entry = project / "SPEC.md"
    entry.unlink()
    shutil.rmtree(project / "specification")
    cli("spec build", cwd=project)
    root = follow(entry, "Read the system specification")
    assert root == project / "specification/personal-task-list.system.spec-report.md"
    assert len(markdown_links(entry)) == 1
    assert sorted(project.rglob("SPEC.md")) == [entry]
    assert all(target != entry for page in (project / "specification").rglob("*.md")
               for _, target, _, _ in markdown_links(page))
    overview = root.read_text(encoding="utf-8")
    sections = ["Purpose", "Goals", "Business rules", "User needs", "Scope", "System topology", "Requirements",
                "Quality attributes"]
    assert re.findall(r"^## (.+)$", overview, re.M) == sections
    navigation = overview.split("\n## Purpose\n", 1)[0]
    goals = overview.split("\n## Goals\n", 1)[1].split("\n## ", 1)[0]
    assert "**Supported by** · [Resume interrupted work]" in goals
    rules = overview.split("\n## Business rules\n", 1)[1].split("\n## ", 1)[0]
    rule = follow(root, "Completion preserves a commitment")
    statement = rule.read_text(encoding="utf-8").split("## Statement\n\n", 1)[1].split("\n\n", 1)[0]
    assert statement in rules and "[All business rules](<business-rules.md>)" in navigation
    assert "[Business rules](" not in rules
    for section in sections:
        assert f"[{section}](<#{section.lower().replace(' ', '-')}>)" in navigation
    assert all(f"[{name}]" in navigation for name in (
        "Jobs", "Business rules", "Components", "Interfaces", "Features", "Capabilities",
        "Concept definitions", "Concept uses", "Use cases", "All requirements", "Attribute index", "Quality models"))
    topology = overview.split("## System topology\n", 1)[1].split("\n## ", 1)[0]
    assert "```mermaid" in topology and "contains" in topology and "provides interface" in topology
    assert "Task service" in topology and "Task list forms" in topology
    for name in ("Components", "Interfaces"):
        assert follow(root, name).name == name.lower() + ".md"
    feature = follow(follow(root, "Features"), "Pick up where I left off")
    promise = follow(feature, "Preserve tasks across restart")
    for name, (collection, source) in {
        "Personal task list": (None, "specs/personal-task-list.system.spec.md"),
        "Pick up where I left off": ("Features", "specs/features/pick-up-where-i-left-off.feature.spec.md"),
        "Task service": ("Components", "app/specs/task-service.component.spec.md"),
        "Local task storage": ("Capabilities", "specs/capabilities/local-task-storage.capability.spec.md"),
        "My checklist": ("Concept uses", "specs/personal-task-list.system.spec.md"),
    }.items():
        declaration = (project / source).read_text(encoding="utf-8")
        subject = re.search(r"^#+ " + re.escape(name) + r"\n\n#+ Description\n\n([^\n]+)", declaration, re.M)
        description = subject[1]
        directory = follow(root, collection) if collection else root
        page = follow(directory, name) if collection else root
        text = page.read_text(encoding="utf-8")
        lead = description if page == root else f"## Description\n\n{description}"
        assert text.split("---\n", 2)[2].startswith(f"# {name}\n\n{lead}")
        assert text.count(description) == directory.read_text(encoding="utf-8").count(description) == 1
    # Every source in this fixture gets one report documenting its H1. Pages
    # for inline subjects document a deeper heading in that same source.
    claims = {}
    for page in (project / "specification").rglob("*.spec-report.md"):
        header = page.read_text(encoding="utf-8").split("---\n", 2)[1]
        primary, = re.findall(r'^source: "([^"]+)"$', header, re.M)
        uri = urlsplit(primary)
        key = ((page.parent / unquote(uri.path)).resolve(), unquote(uri.fragment))
        assert key[0].is_file() and key[1], page
        assert all(dest != key[0] for _, dest, _, _ in markdown_links(page)), page
        assert key not in claims, (page, claims.get(key))
        claims[key] = page
    for source in project.rglob("*.spec.md"):
        title = re.search(r"^# (.+)$", source.read_text(encoding="utf-8"), re.M)[1]
        anchor = re.sub(r"[^a-z0-9 -]", "", title.lower()).replace(" ", "-")
        assert (source.resolve(), anchor) in claims, source
    component = follow(follow(root, "Components"), "Task service")
    assert follow(follow(follow(root, "Capabilities"), "Local task storage"), "Task service") == component
    assert component == project / "specification/components/task-service/task-service.component.spec-report.md"
    component_requirements = follow(component, "Requirements")
    acknowledgement = follow(component_requirements, "Acknowledge after commit")
    assert acknowledgement.parent == component.parent / "requirements"
    statement = acknowledgement.read_text(encoding="utf-8").split("## Statement\n\n", 1)[1].split("\n\n", 1)[0]
    assert component_requirements.read_text(encoding="utf-8").count(statement) == 1
    assert "## Goal" in follow(feature, "Return to my list").read_text(encoding="utf-8")
    qualities_index = follow(root, "Attribute index")
    qualities = follow(qualities_index, "Durability")
    assert follow(promise, "Durability") == qualities
    for page in (qualities_index, promise):
        assert any(target == qualities and not fragment
                   for _, target, fragment, _ in markdown_links(page))
    assert qualities == project / "specification/quality/durability/durability.quality-attribute.spec-report.md"
    assert promise == project / "specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md"
    assert not (qualities.parent / "restart.requirement.spec-report.md").exists()
    assert follow(qualities, "Acknowledge after commit") == acknowledgement
    assert follow(qualities, "Preserve tasks across restart") == promise
    definitions = follow(qualities, "Source definition")
    assert definitions == project / "specification/quality-models/personal-data-quality.quality-model.spec-report.md"
    assert claims[project / "specs/quality-models/personal-data.quality-model.spec.md", "personal-data-quality"] == definitions
    use_case = follow(feature, "Return to my list")
    passage = use_case.read_text(encoding="utf-8").split("### Keep recorded work available\n", 1)[1]
    assert passage.index("#### Notes") < passage.index("## Related specifications")
    assert claims[project / "specs/concepts/checklist.concept.spec.md", "checklist"] == follow(use_case, "Complete")
    headings = re.findall(r"^## (.+)$", promise.read_text(encoding="utf-8"), re.M)
    assert headings[:3] == ["Statement", "Applies to", "Criterion"]
    text = " ".join(promise.read_text(encoding="utf-8").split())
    for phrase in ("acknowledges a change", "same database", "identity, title and completion state",
                   "Deleted entries must stay absent", "Zero discrepancies across the two views"):
        assert phrase in text, phrase
    assert follow(promise, "restart acceptance test") == project / ACCEPTANCE
    targets = markdown_links(promise)
    assert any("return-to-my-list" in target.name and fragment.endswith("keep-recorded-work-available")
               for _, target, fragment, _ in targets), "missing refinement of the scenario clause"
    requirements = follow(root, "All requirements")
    assert requirements == root.parent / "requirements.md"
    assert follow(requirements, "Preserve tasks across restart") == promise
    assert follow(requirements, "Task service") == component_requirements
    use_cases = follow(root, "Use cases")
    assert use_cases == root.parent / "use-cases.md"
    assert follow(use_cases, "Return to my list") == use_case
    assert not (component.parent / "use-cases.md").exists()
    assert not list((project / "specification").rglob("index.md"))
    cli("check", cwd=project)
    for page, rule in [(promise, "spec.content"), (feature, "spec.navigation"),
                       (qualities, "spec.content"), (root, "spec.navigation"),
                       (requirements, "spec.membership"), (qualities, "spec.document-type"),
                       (promise, "spec.document-type"), (definitions, "spec.document-type")]:
        body = page.read_text(encoding="utf-8")
        if rule == "spec.document-type":
            if page == definitions:
                changed = re.sub(r"^source: .+$", 'source: "wrong.spec.md"', body, flags=re.M)
            elif page == qualities:
                changed = body.replace('durability.quality-attribute.spec.md#durability"', 'personal-task-list.system.spec.md#personal-task-list"')
            else:
                changed = body.replace("type: requirement-spec-report", "type: feature-spec-report")
        elif rule == "spec.membership":
            changed = body.replace("\n- [Preserve tasks across restart]", "\n\nSee [Preserve tasks across restart]")
        elif page == promise:
            changed, count = re.subn(r"same\s+database", "different database", body)
            assert count >= 1
        elif page == qualities:
            changed = body.replace("acknowledged work", "submitted work")
        else:
            changed = body
            destination = qualities_index if page == root else promise
            for label, target, _, link in markdown_links(page):
                if target == destination:
                    changed = changed.replace(link, label)
        assert changed != body
        page.write_text(changed, encoding="utf-8")
        report = json.loads(cli("check", "--json", cwd=project, status=1))
        assert rule in {item["rule"] for item in report["diagnostics"]}, report
        cli("spec build", cwd=project)
        assert page.read_text(encoding="utf-8") == body
    cli("check", cwd=project)


def inspect(cli):
    for kind, (slug, title) in TYPES.items():
        args = [slug] if slug else []
        assert title in cli(kind + " list", "--path", TODO)
        assert title in cli(kind + " show", *args, "--path", TODO)
        if kind != "system":
            checked = json.loads(cli(kind + " check", *args, "--path", TODO, "--json"))
            assert checked["complete"] and checked["summary"]["files_checked"] >= 1
    assert "Preserve tasks across restart" in cli("requirement list", "--attribute", "durability", "--project", "todo")
    for command, detail in [
        ("job map", "Recall outstanding work"),
        ("job step list", "Reconsider finished work"),
        ("job outcome list", "Less reconstruction"),
    ]:
        assert detail in cli(command, "resume", "--path", TODO)
    assert "s1 --> s2" in cli("job map", "resume", "--format", "mermaid", "--path", TODO)
    assert "Snapshot" in json.loads(cli("model schema"))["$defs"]
    assert "To-do" in cli("guide")
    assert "examples/todo/specs/personal-task-list.system.spec.md" in cli("guide", "examples")
    catalog = json.loads((ROOT / "docs/handbook/topics.json").read_bytes())
    coverage = (TODO / "COVERAGE.md").read_text(encoding="utf-8")
    exercises = coverage.split("## Handbook and authoring exercises", 1)[1].split("## Field inventory", 1)[0]
    for entry in catalog:
        assert f"`{entry['topic']}`" in exercises, f"assign a review exercise for {entry['topic']}"
        assert cli("guide", entry["topic"]) == cli("guide", entry["page"]), entry["topic"]
        for rule in entry.get("rules", []):
            assert cli("guide", rule) == cli("guide", entry["topic"]), rule


def author(cli, tmp):
    project = project_copy(tmp, "author")
    # Create a second system inside the scratch copy and run the real initializers.
    draft = project / "draft"
    draft.mkdir()
    cli("system init", "--path", draft, cwd=project)
    component = draft / "intake"
    component.mkdir()
    cli("component init", "Intake", "--path", component, "--part-of", draft / "specs/new-system.system.spec.md", cwd=project)
    for kind in TYPES.keys() - {"system"}:
        cli(kind + " init", "sample", "--path", draft, cwd=project)
        assert list((draft / "specs").rglob(f"sample.{kind}.spec.md"))
    for preset, marker in [("odi", "## Steps"),
                           ("job-story", "When [situation]")]:
        cli("job init", preset, "--preset", preset, "--path", draft, cwd=project)
        assert marker in cli("job show", preset, "--path", draft, cwd=project)
        assert json.loads(cli("job check", preset, "--preset", preset, "--json", "--path", draft, cwd=project))["complete"]
    cli("use-case init", "long-form", "--template", "expanded", "--path", draft, cwd=project)
    assert "## Stakeholder interests" in cli("use-case show", "long-form", "--path", draft, cwd=project)
    cli("requirement init", "quality", "--preset", "quality", "--path", draft, cwd=project)
    assert "## Criterion" in cli("requirement show", "quality", "--path", draft, cwd=project)
    target = draft / "specs/new-system.system.spec.md"
    args = ("personal-data", "--system", target, "--attributes", "durability")
    before = target.read_bytes()
    preview = cli("quality-model apply", *args, "--dry-run", cwd=project)
    assert target.read_bytes() == before
    cli("quality-model apply", *args, cwd=project)
    copied = draft / "specs/quality/durability/durability.quality-attribute.spec.md"
    assert copied.read_text(encoding="utf-8") in preview and target.read_bytes() == before
    cli("requirement init", "commit-data", "--attribute", "durability", "--path", draft, cwd=project)
    assert (draft / "specs/requirements/commit-data.requirement.spec.md").is_file()
    cli("requirement init", "data-loss-tolerance", "--preset", "quality", "--attribute", "durability", "--path", draft, cwd=project)
    quality_requirement = draft / "specs/requirements/data-loss-tolerance.requirement.spec.md"
    shown = json.loads(cli("requirement show", quality_requirement, "--json", cwd=project))
    assert shown["requirement"]["quality_attributes"] == ["durability"]
    cli("requirement init", "save-draft", "--attribute", "durability", "--path", component, cwd=project)
    assert (component / "specs/requirements/save-draft.requirement.spec.md").is_file()
    assert "durability" in cli("requirement show", "commit-data", "--path", draft, cwd=project)
    focused = json.loads(cli("requirement check", "commit-data", "--path", draft, "--json", cwd=project))
    assert focused["complete"] and focused["summary"]["errors"] == 0
    cli("quality-model apply", *args, cwd=project, status=2)
    assert target.read_bytes() == before
    cli("quality-model init", "team", "--from", "furps", "--path", draft, cwd=project)
    assert "Supportability" in cli("quality-model show", "team", "--path", draft, cwd=project)
    with (project / REQUIREMENT).open("a", encoding="utf-8") as file:
        file.write("\n[Missing verification file](missing.py)\n")
    broken = json.loads(cli("check", "--json", cwd=project, status=1))
    assert any(d["rule"] == "reference.unresolved" for d in broken["diagnostics"])
    # Inspect the faithful representation of an invalid source alongside the clean exports.
    partial = cli("model export", cwd=project, status=1)
    partial_file = tmp / "author.json"
    partial_file.write_text(partial, encoding="utf-8")
    assert json.loads(cli("model validate", partial_file, "--json"))["valid"]


def acceptance(project):
    return run([sys.executable, "-I", project / ACCEPTANCE])


def regressions(cli, tmp):
    for name, path, old, new, reason in [
        ("relationship", REQUIREMENT, "## Refines", "## Notes", "missing restart refinement"),
        ("test-link", REQUIREMENT, "[restart acceptance test](../../tests/acceptance.py)",
         "restart acceptance test", "restart requirement must link"),
        ("persistence", Path("app/main.go"), 'sql.Open("sqlite3", file)',
         'sql.Open("sqlite3", ":memory:")', "restart lost or changed open tasks"),
    ]:
        copy = project_copy(tmp, name)
        file = copy / path
        source = file.read_text(encoding="utf-8")
        assert source.count(old) == 1, f"update the {name} regression fixture"
        file.write_text(source.replace(old, new), encoding="utf-8")
        try:
            if name == "persistence":
                acceptance(copy)
            else:
                cli("spec build", cwd=copy)
                check_project(cli, copy, tmp)
        except (AssertionError, RuntimeError) as error:
            if reason not in str(error):
                raise
            print(f"Rejected {name} regression: {reason}")
        else:
            raise AssertionError(f"{name} regression passed the gate")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--inventory", action="store_true", help="print detailed capability coverage")
    args = parser.parse_args()
    with tempfile.TemporaryDirectory() as folder:
        tmp = Path(folder).resolve()
        binary = tmp / ("dotspec.exe" if os.name == "nt" else "dotspec")
        run(["go", "build", "-mod=readonly", "-o", binary, "."])
        cli = CLI(binary)
        check_project(cli, TODO, tmp)
        inspect(cli)
        spec_journey(cli, tmp)
        author(cli, tmp)
        run(["go", "vet", "./..."], cwd=TODO / "app")
        run(["go", "mod", "tidy", "-diff"], cwd=TODO / "app")
        print(acceptance(TODO).strip())
        regressions(cli, tmp)
        # The project's subprocess tests own shell-specific output assertions.
        events = run(["go", "test", "-json", "-count=1", "./internal/integration", "-run", "^TestCLI$/^help$"])
        assert any(e.get("Test") == "TestCLI/help" and e["Action"] == "pass"
                   for e in map(json.loads, events.splitlines())), "help.txtar did not run"
        tooling = {"man", "completion bash", "completion zsh", "completion fish", "completion powershell"}
        public = cli.public_commands()
        groups = set(TYPES) | {"component","spec", "model", "job outcome", "job step", "completion"}
        for group in groups | {""}:
            assert "Available Commands:" in cli("help", *group.split())
        missing = public - cli.exercised - groups - tooling - {""}
        assert not missing, f"public commands need coverage: {sorted(missing)}"
        inventory = [sys.executable, "-I", Path(__file__).with_name("inventory.py"), binary, tmp / "author.json"]
        print(run(inventory + (["--details"] if args.inventory else [])).strip())
        if args.inventory:
            for cmd in sorted(public):
                via = "help.txtar" if cmd in tooling else "group help" if cmd in groups or not cmd else "example-check"
                print(f"dotspec {cmd:<28} {via}")
        print(f"{len(public)} public command paths accounted for.")


if __name__ == "__main__":
    main()
