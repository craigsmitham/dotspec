"""Pack the six checksummed GoReleaser executables into one npm package."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import tarfile
import tempfile
import zipfile

ROOT = Path(__file__).resolve().parents[2]


def pack(version, assets, output):
    checksums = {}
    for line in (assets / "checksums.txt").read_text().splitlines():
        digest, name = line.split()
        checksums[name] = digest
    with tempfile.TemporaryDirectory() as tmp:
        package = Path(tmp)
        for os in ("darwin", "linux", "windows"):
            for arch in ("amd64", "arm64"):
                windows = os == "windows"
                exe = "dotspec.exe" if windows else "dotspec"
                ext = "zip" if windows else "tar.gz"
                archive = assets / f"dotspec_{version}_{os}_{arch}.{ext}"
                if hashlib.sha256(archive.read_bytes()).hexdigest() != checksums[archive.name]:
                    raise ValueError(f"Checksum mismatch for {archive.name}")
                dest = package / "native" / f"{os}-{arch}" / exe
                dest.parent.mkdir(parents=True)
                if windows:
                    with zipfile.ZipFile(archive) as bundle:
                        dest.write_bytes(bundle.read(exe))
                else:
                    with tarfile.open(archive) as bundle:
                        dest.write_bytes(bundle.extractfile(exe).read())
                dest.chmod(0o755)
        shutil.copy(ROOT / "LICENSE", package)
        shutil.copy(Path(__file__).with_name("README.md"), package)
        shutil.copy(Path(__file__).with_name("dotspec.cjs"), package)
        (package / "dotspec.cjs").chmod(0o755)
        manifest = {
            "name": "dotspec", "version": version, "license": "MIT",
            "description": "Write, check, and build software requirements in Markdown.",
            "repository": {"type": "git", "url": "git+https://github.com/dotspec-md/dotspec.git"},
            "bin": {"dotspec": "dotspec.cjs"}, "os": ["darwin", "linux", "win32"],
            "cpu": ["x64", "arm64"],
        }
        (package / "package.json").write_text(json.dumps(manifest, indent=2) + "\n")
        output.mkdir(parents=True, exist_ok=True)
        subprocess.run([shutil.which("npm"), "pack", "--pack-destination", str(output.resolve())],
                       cwd=package, check=True)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("version", help="release version without its leading v")
    parser.add_argument("assets", type=Path)
    parser.add_argument("output", type=Path)
    args = parser.parse_args()
    pack(args.version, args.assets, args.output)
