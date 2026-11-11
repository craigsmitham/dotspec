#!/usr/bin/env python3
"""Platform-neutral launcher for the explicitly provisioned model toolchain."""

from pathlib import Path
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[2]
ENVIRONMENT = ROOT / ".venv-model"
PYTHON = ENVIRONMENT / ("Scripts/python.exe" if sys.platform == "win32" else "bin/python")


def main() -> int:
    if len(sys.argv) != 2 or sys.argv[1] not in {"setup", "generate", "check", "test"}:
        raise SystemExit("usage: run.py setup|generate|check|test")
    action = sys.argv[1]
    if action == "setup":
        if sys.version_info[:2] != (3, 11):
            raise SystemExit("model setup requires the pinned Python 3.11 interpreter")
        subprocess.run([sys.executable, "-m", "venv", str(ENVIRONMENT)], check=True)
        subprocess.run([
            str(PYTHON), "-m", "pip", "install", "--disable-pip-version-check",
            "-r", str(ROOT / "scripts/generate-model/requirements.lock"),
        ], check=True)
        return 0
    if not PYTHON.is_file():
        raise SystemExit("model toolchain is not provisioned; run mise run model-setup")
    if action == "test":
        command = [str(PYTHON), "-m", "unittest", "discover", "-s", "scripts/generate-model", "-p", "test_*.py"]
    else:
        command = [str(PYTHON), "scripts/generate-model/generate.py"]
        if action == "check":
            command.append("--check")
    return subprocess.run(command, cwd=ROOT, check=False).returncode


if __name__ == "__main__":
    raise SystemExit(main())
