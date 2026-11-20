"""Exercise task forms across two independent server processes on one SQLite file."""
import argparse
from contextlib import contextmanager
from html.parser import HTMLParser
from pathlib import Path
import os
import subprocess
import tempfile
import time
from urllib.error import HTTPError
from urllib.parse import urlencode
from urllib.request import urlopen


class TaskPage(HTMLParser):
    def __init__(self, body):
        super().__init__()
        self.forms = {}
        self.links = set()
        self.form = None
        self.titles = []
        self.in_title = False
        self.feed(body)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "a":
            self.links.add(attrs["href"])
        if tag == "form":
            assert attrs["method"].lower() == "post"
            self.form = {}
            self.forms[attrs["action"]] = self.form
        if tag == "input" and self.form is not None:
            self.form[attrs["name"]] = attrs.get("value", "")
        if tag == "p" and attrs.get("class") == "task":
            self.in_title = True

    def handle_data(self, data):
        if self.in_title:
            self.titles.append(data)

    def handle_endtag(self, tag):
        if tag == "form":
            self.form = None
        if tag == "p":
            self.in_title = False

    def entries(self):
        entries = {action.removesuffix("/rename"): fields["title"]
                   for action, fields in self.forms.items() if action.endswith("/rename")}
        assert sorted(entries.values()) == sorted(self.titles), "visible titles differ from editable entries"
        return entries


@contextmanager
def running(binary, db):
    with tempfile.TemporaryFile() as errors, (db.parent / "address").open("wb") as address:
        process = subprocess.Popen([binary, "-port", "0", "-db", db], stdout=address, stderr=errors)
        try:
            deadline = time.monotonic() + 30
            while time.monotonic() < deadline:
                if process.poll() is not None:
                    errors.seek(0)
                    raise RuntimeError(errors.read().decode())
                line = (db.parent / "address").read_bytes()
                if line.endswith(b"\n"):
                    yield line.decode().strip()
                    break
                time.sleep(0.025)
            else:
                raise TimeoutError("task list did not start")
        finally:
            process.kill()
            process.wait(timeout=10)


def page(base, path="/", data=None):
    with urlopen(base + path, data, timeout=10) as response:
        return TaskPage(response.read().decode())


def submit(base, current, action, **fields):
    form = current.forms[action].copy()
    assert fields.keys() <= form.keys(), f"missing inputs on {action}"
    form.update(fields)
    return page(base, action, urlencode(form).encode())


def verify(binary, db):
    with running(binary, db) as base:
        current = page(base)
        assert current.entries() == {} and "/completed" in current.links
        try:
            submit(base, current, "/tasks", title="  ")
        except HTTPError as error:
            assert error.code == 400
        else:
            raise AssertionError("blank task accepted")
        current = submit(base, current, "/tasks", title="Call the plumber")
        errand, = current.entries()
        try:
            submit(base, current, errand + "/rename", title=" ")
        except HTTPError as error:
            assert error.code == 400 and page(base).entries() == {errand: "Call the plumber"}
        else:
            raise AssertionError("blank rename accepted")
        current = submit(base, current, errand + "/rename", title="Call Pat & Lee <today>")
        current = submit(base, current, "/tasks", title="Collect keys")
        current = submit(base, current, "/tasks", title="Discard this task")
        discarded, = [key for key, title in current.entries().items() if title == "Discard this task"]
        current = submit(base, current, discarded + "/delete")
        current = submit(base, current, errand + "/complete")
        open_entries = current.entries()
        assert list(open_entries.values()) == ["Collect keys"]
        completed = page(base, "/completed")
        finished = completed.entries()
        assert finished == {errand: "Call Pat & Lee <today>"} and "/" in completed.links

        # Return in a fresh view while the service is still running.
        returned = page(base, "/completed")
        assert returned.entries() == finished and page(base).entries() == open_entries
        submit(base, returned, errand + "/reopen")
        assert page(base).entries() == open_entries | finished
        submit(base, page(base), errand + "/complete")
        assert page(base, "/completed").entries() == finished

    with running(binary, db) as base:
        assert page(base).entries() == open_entries, "restart lost or changed open tasks"
        completed = page(base, "/completed")
        assert completed.entries() == finished, "restart lost or changed completed tasks"
        completed = submit(base, completed, errand + "/reopen")
        assert completed.entries() == {}
        assert page(base).entries() == open_entries | finished, "reopen changed task identity"


def run(app):
    with tempfile.TemporaryDirectory() as tmp:
        binary = Path(tmp) / ("todo.exe" if os.name == "nt" else "todo")
        subprocess.run(["go", "build", "-mod=readonly", "-o", binary, "."], cwd=app, check=True)
        verify(binary, Path(tmp) / "tasks.sqlite")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--app", type=Path, default=Path(__file__).resolve().parents[1] / "app")
    run(parser.parse_args().app)
    print("Task forms and persistence across process restart passed.")
