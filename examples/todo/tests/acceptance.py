"""Submit the task-list forms to one local database across independent processes.

The actions follow specs/concepts/checklist.concept.spec.md. The cue shown on
each date comes from specs/requirements/follow-up.requirement.spec.md. The
live and restarted list comparisons exercise specs/requirements/availability.requirement.spec.md
and specs/requirements/restart.requirement.spec.md. The failed saves exercise
app/specs/requirements/acknowledge-after-commit.requirement.spec.md.
"""
import argparse
from contextlib import closing, contextmanager
from datetime import date, datetime, timedelta, timezone
from html.parser import HTMLParser
from pathlib import Path
import os
import sqlite3
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
        self.item_text = {}
        self.item = None
        self.feed(body)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == "li":
            self.item = []
        if tag == "a":
            self.links.add(attrs["href"])
        if tag == "form":
            assert attrs["method"].lower() == "post"
            self.form = {}
            action = attrs["action"]
            self.forms[action] = self.form
            if self.item is not None and action.endswith("/rename"):
                self.item_text[action.removesuffix("/rename")] = self.item
        if tag == "input" and self.form is not None:
            self.form[attrs["name"]] = attrs.get("value", "")
        if tag == "p" and attrs.get("class") == "task":
            self.in_title = True

    def handle_data(self, data):
        if self.item is not None:
            self.item.append(data)
        if self.in_title:
            self.titles.append(data)

    def handle_endtag(self, tag):
        if tag == "li":
            self.item = None
        if tag == "form":
            self.form = None
        if tag == "p":
            self.in_title = False

    def text_for(self, entry):
        return " ".join(" ".join(self.item_text[entry]).split())

    def entries(self):
        entries = {action.removesuffix("/rename"): fields["title"]
                   for action, fields in self.forms.items() if action.endswith("/rename")}
        assert sorted(entries.values()) == sorted(self.titles), "visible titles differ from editable entries"
        return entries


@contextmanager
def running(binary, db, *, zone=None):
    env = os.environ.copy()
    if zone:
        env["TZ"] = zone
    with tempfile.TemporaryFile() as errors, (db.parent / "address").open("wb") as address:
        process = subprocess.Popen([binary, "-port", "0", "-db", db], stdout=address, stderr=errors, env=env)
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


def rejected(status, base, current, action, *, message=None, **fields):
    try:
        submit(base, current, action, **fields)
    except HTTPError as error:
        with error:
            assert error.code == status, (action, error.code)
            if message:
                body = error.read().decode()
                assert message in body, (action, body)
            return error.url
    else:
        raise AssertionError(f"{action} accepted a submission that should return {status}")


def add(base, title):
    start = page(base)
    before = start.entries()
    current = submit(base, start, "/tasks", title=title)
    after = current.entries()
    entry, = after.keys() - before.keys()
    assert after == before | {entry: title.strip()}
    return current, entry


def expect_waiting(waiting, entry, reason, follow_up, cue):
    text = waiting.text_for(entry)
    assert reason in text and follow_up in text, (entry, text)
    form = waiting.forms[entry + "/update-follow-up"]
    assert form["waiting_for"] == reason and form["follow_up"] == follow_up, (entry, form)
    assert cue in text, (entry, text)
    other = "Upcoming" if cue == "Follow up now" else "Follow up now"
    assert other not in text, (entry, text)


def verify(binary, db):
    with running(binary, db) as base:
        current = page(base)
        assert current.entries() == {} and {"/waiting", "/completed"} <= current.links
        rejected(400, base, current, "/tasks", title="  ")
        current, errand = add(base, "Call the plumber")
        rejected(400, base, current, errand + "/rename", title=" ")
        assert page(base).entries() == {errand: "Call the plumber"}
        current = submit(base, current, errand + "/rename", title="  Call Pat & Lee <today>  ")
        submit(base, current, errand + "/complete")
        _, keys = add(base, "  Collect keys  ")
        current, repair = add(base, "Arrange repair")
        today = date.today()
        visit_day = today.isoformat()
        parts_day = (today + timedelta(days=7)).isoformat()
        yesterday = (today - timedelta(days=1)).isoformat()

        # An incomplete handoff leaves the errand actionable.
        for reason, day in [("   ", visit_day), ("Visit time", ""), ("Visit time", "2025-02-29")]:
            rejected(400, base, current, repair + "/wait", waiting_for=reason, follow_up=day)
            assert page(base).entries()[repair] == "Arrange repair" and not page(base, "/waiting").entries()

        after_wait = submit(base, current, repair + "/wait", waiting_for="  Visit time from Lee  ", follow_up=visit_day)
        assert after_wait.entries() == {keys: "Collect keys"}
        # The Complete button in an earlier Open tab belongs to the old state.
        rejected(409, base, current, repair + "/complete")
        waiting = page(base, "/waiting")
        waiting = submit(base, waiting, repair + "/rename", title="Book sink repair")
        expect_waiting(waiting, repair, "Visit time from Lee", visit_day, "Follow up now")

        current, parts = add(base, "Get replacement parts")
        submit(base, current, parts + "/wait", waiting_for="Dispatch notice", follow_up=yesterday)
        waiting = page(base, "/waiting")
        expect_waiting(waiting, parts, "Dispatch notice", yesterday, "Follow up now")

        # The supplier promises a delivery receipt next week. Both details change together.
        update = parts + "/update-follow-up"
        rejected(400, base, waiting, update, waiting_for="Delivery receipt", follow_up="")
        expect_waiting(page(base, "/waiting"), parts, "Dispatch notice", yesterday, "Follow up now")
        waiting = submit(base, waiting, update, waiting_for="Delivery receipt", follow_up=parts_day)
        expect_waiting(waiting, parts, "Delivery receipt", parts_day, "Upcoming")

        # Two key collections share a title. Only the second needs a pickup time.
        current, second_keys = add(base, "Collect keys")
        submit(base, current, second_keys + "/wait", waiting_for="Pickup time", follow_up=yesterday)
        assert page(base).entries() == {keys: "Collect keys"}
        waiting = page(base, "/waiting")
        expect_waiting(waiting, second_keys, "Pickup time", yesterday, "Follow up now")
        submit(base, waiting, second_keys + "/delete")
        current, replacement = add(base, "Collect keys")
        assert replacement != second_keys
        assert current.entries() == {keys: "Collect keys", replacement: "Collect keys"}
        rejected(404, base, waiting, second_keys + "/resume", message="entry is no longer on your list")
        submit(base, current, replacement + "/delete")

        # A return while the service runs, including a follow-up to a settled errand.
        completed = page(base, "/completed")
        finished = {errand: "Call Pat & Lee <today>"}
        assert completed.entries() == finished and "/" in completed.links
        assert submit(base, completed, errand + "/reopen").entries() == {}
        assert page(base).entries() == {keys: "Collect keys"} | finished
        submit(base, page(base), errand + "/complete")

        open_entries = {keys: "Collect keys"}
        waiting_entries = {repair: "Book sink repair", parts: "Get replacement parts"}
        assert page(base).entries() == open_entries
        assert page(base, "/completed").entries() == finished
        waiting = page(base, "/waiting")
        assert waiting.entries() == waiting_entries
        expect_waiting(waiting, repair, "Visit time from Lee", visit_day, "Follow up now")
        expect_waiting(waiting, parts, "Delivery receipt", parts_day, "Upcoming")

    with running(binary, db) as base:
        assert page(base).entries() == open_entries, "restart lost or changed open tasks"
        assert page(base, "/completed").entries() == finished, "restart lost or changed completed tasks"
        waiting = page(base, "/waiting")
        assert waiting.entries() == waiting_entries, "restart lost or changed waiting tasks"
        expect_waiting(waiting, repair, "Visit time from Lee", visit_day, "Follow up now")
        expect_waiting(waiting, parts, "Delivery receipt", parts_day, "Upcoming")

        # Two replies: Lee needs a photo, but the receipt establishes delivery.
        after_resume = submit(base, waiting, repair + "/resume")
        assert after_resume.entries() == {parts: waiting_entries[parts]} and "/" in after_resume.links
        open_page = page(base)
        assert open_page.entries() == open_entries | {repair: waiting_entries[repair]}
        assert "Visit time from Lee" not in open_page.text_for(repair) and visit_day not in open_page.text_for(repair)
        wait_form = open_page.forms[repair + "/wait"]
        assert wait_form["waiting_for"] == wait_form["follow_up"] == ""
        after_complete = submit(base, after_resume, parts + "/complete")
        assert after_complete.entries() == {} and "/completed" in after_complete.links
        completed = page(base, "/completed")
        assert completed.entries() == finished | {parts: waiting_entries[parts]}
        assert "Delivery receipt" not in completed.text_for(parts) and parts_day not in completed.text_for(parts)
        assert page(base, "/waiting").entries() == {}
        submit(base, completed, parts + "/reopen")
        open_page = page(base)
        assert open_page.entries() == open_entries | waiting_entries
        wait_form = open_page.forms[parts + "/wait"]
        assert wait_form["waiting_for"] == wait_form["follow_up"] == ""


def verify_changed_day(binary, db):
    # UTC+14 is 26 hours ahead of UTC-12. A date due in the eastern zone
    # is still in the future in the western zone, even near either midnight.
    east = timezone(timedelta(hours=14))
    day = datetime.now(east).date().isoformat()
    with running(binary, db, zone="Etc/GMT+12") as base:
        current, room = add(base, "Reserve library room")
        submit(base, current, room + "/wait", waiting_for="Room confirmation", follow_up=day)
        expect_waiting(page(base, "/waiting"), room, "Room confirmation", day, "Upcoming")
    with running(binary, db, zone="Etc/GMT-14") as base:
        waiting = page(base, "/waiting")
        assert waiting.entries() == {room: "Reserve library room"}
        expect_waiting(waiting, room, "Room confirmation", day, "Follow up now")


def verify_follow_up_calendar(binary, db):
    with running(binary, db) as base:
        current, visit = add(base, "Arrange a visit")
        submit(base, current, visit + "/wait", waiting_for="Visit time", follow_up="10000-02-29")
        waiting = page(base, "/waiting")
        expect_waiting(waiting, visit, "Visit time", "10000-02-29", "Upcoming")
        rejected(400, base, waiting, visit + "/update-follow-up", follow_up="10100-02-29")
        expect_waiting(page(base, "/waiting"), visit, "Visit time", "10000-02-29", "Upcoming")


def verify_failed_save(binary, db, decision, failing_write):
    day = date.today().isoformat()
    with running(binary, db) as base:
        current, repair = add(base, "Book sink repair")
        submit(base, current, repair + "/wait", waiting_for="Visit time", follow_up=day)

    with closing(sqlite3.connect(db)) as connection:
        connection.execute(f"""CREATE TRIGGER simulated_failure BEFORE {failing_write}
                               BEGIN SELECT RAISE(ABORT, 'simulated write failure'); END""")

    with running(binary, db) as base:
        waiting = page(base, "/waiting")
        action = f"{repair}/{decision}"
        assert rejected(500, base, waiting, action, message="Could not access the task list.") == base + action, \
            f"failed {decision} redirected"
        waiting = page(base, "/waiting")
        assert waiting.entries() == {repair: "Book sink repair"}
        expect_waiting(waiting, repair, "Visit time", day, "Follow up now")
        assert page(base).entries() == page(base, "/completed").entries() == {}


def run(app):
    with tempfile.TemporaryDirectory() as tmp:
        binary = Path(tmp) / ("todo.exe" if os.name == "nt" else "todo")
        subprocess.run(["go", "build", "-mod=readonly", "-o", binary, "."], cwd=app, check=True)
        verify(binary, Path(tmp) / "tasks.sqlite")
        verify_follow_up_calendar(binary, Path(tmp) / "dates.sqlite")
        # Complete fails at its second write, after the waiting details were removed in the same transaction.
        verify_failed_save(binary, Path(tmp) / "failed-completion.sqlite", "complete", "UPDATE OF completed ON tasks")
        # Resume fails at its only write, the removal of the waiting details.
        verify_failed_save(binary, Path(tmp) / "failed-resume.sqlite", "resume", "DELETE ON pending_replies")
        if os.name == "posix":  # Go on Unix reads the process's TZ setting.
            verify_changed_day(binary, Path(tmp) / "calendar.sqlite")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--app", type=Path, default=Path(__file__).resolve().parents[1] / "app")
    run(parser.parse_args().app)
    print("Task forms and persistence across process restart passed.")
