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
    "job": ("personal-commitments", "See personal commitments through"),
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

    def targets(owner, role):
        return {r["resolution"]["target"] for r in model["references"]
                if r.get("context_owner") == owner["id"] and r.get("role") == role
                and r["resolution"]["status"] == "resolved"}

    requirement = one("Requirement", slug="restart")
    use = one("UseCase", slug="resume-list")
    job = one("Job", slug="personal-commitments")
    availability = one("Requirement", slug="availability", authored_id="keep-recorded-work-available")
    assert availability["id"] in targets(requirement, "refines"), "missing restart refinement"
    assert one("BusinessRule", slug="completion")["id"] in targets(availability, "derived_from")
    assert any(availability["id"] in targets(step, "illustrates") for step in elements
               if step["kind"] == "ScenarioStep" and step["source_document"] == use["source_document"])
    assert job["id"] in targets(use, "supports")
    assert one("DesiredOutcome", owner=job["id"], authored_id="less-reconstruction")["id"] in targets(requirement, "supports")
    sources = {s["id"]: s for s in model["sources"]}
    # This project exercises every supported specification file type.
    root_types = {"System", "Component", "Job", "UseCase", "Concept", "BusinessRule",
                  "Requirement", "Feature", "Capability", "QualityModel", "SystemQualityAttribute"}
    primary = {e["id"]: e for e in elements if e["kind"] in root_types}
    assert {e["kind"] for e in primary.values()} == root_types
    for doc in model["documents"]:
        source = sources[doc["source"]]
        if not source.get("selected"):
            continue
        root_id, = [key for key in doc["elements"] if key in primary]
        kind = doc["declared_type"].removesuffix("-spec")
        expected = "SystemQualityAttribute" if kind == "quality-attribute" else kind.title().replace("-", "")
        assert primary[root_id]["kind"] == expected and primary[root_id]["source_document"] == doc["source"]
        assert source["locator"].endswith("." + kind + ".spec.md")
    for kind in ("System", "Component", "Feature", "Capability", "ConceptUse"):
        subject = one(kind)
        spans = subject["description"]["spans"]
        assert spans and all(span["source"] == subject["source_document"] for span in spans), kind
        data = sources[subject["source_document"]]["markdown"].encode()
        assert all(data[span["start"]:span["end"]].strip() for span in spans), kind
    system = one("System")
    assert targets(use, "subject_system") == {system["id"]}
    attribute = one("SystemQualityAttribute", owner=system["id"], authored_id="durability")
    assert any(r["owner"] == requirement["id"] and r["resolution"].get("target") == attribute["id"]
               for r in model["label_occurrences"])
    assert targets(one("ConceptUse"), "definition") == {one("Concept", slug="checklist")["id"]}
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
    sections = ["Framing", "Purpose", "Goals", "Stakeholders and user needs", "Scope and environment", "System topology",
                "Features", "Concepts", "Quality attributes"]
    assert re.findall(r"^## (.+)$", overview, re.M) == sections
    navigation = overview.split("\n## Framing\n", 1)[0]
    rows = [line for line in navigation.splitlines() if line.startswith("| ")]
    for title, index in (("Features", "features"), ("Capabilities", "capabilities"),
                         ("Use cases", "use-cases"), ("Requirements", "requirements")):
        row, = [line for line in rows if line.split(" | ")[0] in
                (f"| [{title}](<#{index}>)", f"| {title}")]
        assert f"[All {title.lower()}](<{index}.md>)" in row, row
    concepts_row, = [line for line in rows if line.startswith("| [Concepts]")]
    assert "[Concept definitions](<concepts.md>)" in concepts_row and "[Concept uses](<concept-uses.md>)" in concepts_row
    assert "### Behavior" not in overview and "### Obligations" not in overview
    goals = overview.split("\n## Goals\n", 1)[1].split("\n## ", 1)[0]
    assert "Reduce the effort of working out what still needs attention." in goals
    assert "**Supported by**" not in goals and "**Relied on by**" not in overview
    for title, collection, anchor, related, inverse in (
        ("Resume unfinished work", "goals", "resume-unfinished-work", "See personal commitments through", "Supported by"),
        ("Same local storage", "assumptions", "same-local-storage", "Preserve tasks across restart", "Relied on by"),
    ):
        detail = follow(root, title)
        assert detail == root.parent / f"{collection}.md"
        assert f"## {title}\n" in detail.read_text(encoding="utf-8")
        assert f"**{inverse}**\n\n- [{related}]" in detail.read_text(encoding="utf-8")
        assert "**Referenced by**" not in detail.read_text(encoding="utf-8")
        for page in (root, follow(detail, related)):
            assert any(target == detail and fragment == anchor for _, target, fragment, _ in markdown_links(page))
    rules = follow(root, "All business rules")
    rule = follow(rules, "Completion preserves a commitment")
    statement = rule.read_text(encoding="utf-8").split("## Statement\n\n", 1)[1].split("\n\n", 1)[0]
    assert statement in rules.read_text(encoding="utf-8")
    for section in sections:
        assert f"[{section}](<#{section.lower().replace(' ', '-')}>)" in navigation
    assert all(f"[{name}]" in navigation for name in (
        "Jobs", "All business rules", "Components", "Interfaces", "All features", "All capabilities",
        "Concept definitions", "Concept uses", "All use cases", "All requirements", "All quality attributes", "Quality models"))
    topology = overview.split("## System topology\n", 1)[1].split("\n## ", 1)[0]
    assert "```mermaid" in topology and "contains" in topology and "provides interface" in topology
    assert "Task service" in topology and "Task list forms" in topology
    assert "Provides [Local task storage]" in topology and "## Capabilities" not in overview
    concepts = overview.split("\n## Concepts\n", 1)[1].split("\n## ", 1)[0]
    assert "[My checklist]" in concepts and "Uses [Checklist](<concepts/" in concepts
    for name in ("Components", "Interfaces"):
        assert follow(root, name).name == name.lower() + ".md"
    feature = follow(follow(root, "All features"), "Pick up where I left off")
    promise = follow(feature, "Preserve tasks across restart")
    for name, (collection, source) in {
        "Personal task list": (None, "specs/personal-task-list.system.spec.md"),
        "Pick up where I left off": ("All features", "specs/features/pick-up-where-i-left-off.feature.spec.md"),
        "Task service": ("Components", "app/specs/task-service.component.spec.md"),
        "Local task storage": ("All capabilities", "specs/capabilities/local-task-storage.capability.spec.md"),
        "My checklist": ("Concept uses", "specs/personal-task-list.system.spec.md"),
    }.items():
        declaration = (project / source).read_text(encoding="utf-8")
        subject = re.search(r"^#+ " + re.escape(name) + r"\n\n#+ Description\n\n([^\n]+)", declaration, re.M)
        description = subject[1]
        directory = follow(root, collection) if collection else root
        page = follow(directory, name) if collection else root
        text = page.read_text(encoding="utf-8")
        lead = description if name in ("Personal task list", "Task service") else f"## Description\n\n{description}"
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
    assert follow(follow(follow(root, "All capabilities"), "Local task storage"), "Task service") == component
    assert follow(follow(component, "All capabilities"), "Local task storage") == follow(follow(root, "All capabilities"), "Local task storage")
    assert component == project / "specification/components/task-service/task-service.component.spec-report.md"
    component_requirements = follow(component, "All requirements")
    acknowledgement = follow(component_requirements, "Acknowledge after commit")
    assert acknowledgement.parent == component.parent / "requirements"
    assert "## Applicable requirements" in component_requirements.read_text(encoding="utf-8")
    assert "**Applies here**" not in component.read_text(encoding="utf-8")
    assert follow(component, "Acknowledge after commit") == acknowledgement
    statement = acknowledgement.read_text(encoding="utf-8").split("## Statement\n\n", 1)[1].split("\n\n", 1)[0]
    assert component_requirements.read_text(encoding="utf-8").count(statement) == 1
    assert "## Goal" in follow(feature, "Return to my list").read_text(encoding="utf-8")
    qualities_index = follow(root, "All quality attributes")
    qualities = follow(qualities_index, "Durability")
    quality_section = overview.split("\n## Quality attributes\n", 1)[1].split("\n## ", 1)[0]
    attribute_source = project / "specs/quality/durability/durability.quality-attribute.spec.md"
    definition = attribute_source.read_text(encoding="utf-8").split("## Definition\n\n", 1)[1].split("\n\n", 1)[0].strip()
    assert definition in quality_section
    assert "[Acknowledge after commit](" not in quality_section and "[Preserve tasks across restart](" not in quality_section
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
    availability = follow(use_case, "Keep recorded work available")
    assert availability == claims[project / "specs/requirements/availability.requirement.spec.md", "keep-recorded-work-available"]
    assert follow(promise, "Keep recorded work available") == availability
    assert follow(follow(root, "All requirements"), "Keep recorded work available") == availability
    assert claims[project / "specs/concepts/checklist.concept.spec.md", "checklist"] == follow(use_case, "Complete")
    headings = re.findall(r"^## (.+)$", promise.read_text(encoding="utf-8"), re.M)
    assert headings[:3] == ["Statement", "Applies to", "Criterion"]
    text = " ".join(promise.read_text(encoding="utf-8").split())
    for phrase in ("acknowledges a change", "same database", "identity, title and completion state",
                   "Deleted entries must stay absent", "Zero discrepancies across the two views"):
        assert phrase in text, phrase
    assert follow(promise, "restart acceptance test") == project / ACCEPTANCE
    requirements = follow(root, "All requirements")
    assert requirements == root.parent / "requirements.md"
    assert follow(requirements, "Preserve tasks across restart") == promise
    assert follow(requirements, "Task service") == component_requirements
    use_cases = follow(root, "All use cases")
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
        ("job step list", "Revise the plan"),
        ("job outcome list", "Less reconstruction"),
    ]:
        assert detail in cli(command, "personal-commitments", "--path", TODO)
    assert "s1 --> s2" in cli("job map", "personal-commitments", "--format", "mermaid", "--path", TODO)
    assert "Snapshot" in json.loads(cli("model schema"))["$defs"]
    assert "Jobs to be done" in cli("docs")
    assert "examples/todo/specs/jobs/personal-commitments.job.spec.md" in cli("docs", "job")
    catalog = json.loads((ROOT / "docs/cli-handbook/topics.json").read_bytes())
    coverage = (TODO / "COVERAGE.md").read_text(encoding="utf-8")
    exercises = coverage.split("## Handbook and authoring exercises", 1)[1].split("## Field inventory", 1)[0]
    for entry in catalog:
        assert f"`{entry['topic']}`" in exercises, f"assign a review exercise for {entry['topic']}"
        menu = cli("docs", entry["topic"])
        for destination in entry["pages"]:
            page = destination["page"]
            address = page.removeprefix("docs/cli-handbook/").replace(".md", "", 1) if page.startswith("docs/cli-handbook/") else page
            assert f"dotspec docs {address}\n" in menu, address
            assert cli("docs", address), address
            for rule, anchor in destination.get("rules", {}).items():
                target = address.split("#")[0] + "#" + anchor if anchor else address
                assert cli("docs", rule) == cli("docs", target), rule


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
    (draft / "release-plan.md").write_text(
        "# Release plan\n\nThe release lead reviews the durability checks before rollout.\n", encoding="utf-8")
    with target.open("a", encoding="utf-8") as file:
        file.write("\n## Project references\n\nThe [release plan](../release-plan.md) assigns the durability review.\n")
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
        ("relationship", REQUIREMENT, "## Refines", "## Related specifications", "missing restart refinement"),
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
