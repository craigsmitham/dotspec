#!/usr/bin/env node
const { spawnSync } = require("node:child_process");
const { join } = require("node:path");

const arch = { x64: "amd64", arm64: "arm64" }[process.arch];
const os = { darwin: "darwin", linux: "linux", win32: "windows" }[process.platform];
if (!os || !arch) {
  console.error(`dotspec does not support ${process.platform}/${process.arch}`);
  process.exit(1);
}
const exe = process.platform === "win32" ? "dotspec.exe" : "dotspec";
const result = spawnSync(join(__dirname, "native", os + "-" + arch, exe), process.argv.slice(2), {
  stdio: "inherit",
});
if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}
if (result.signal) process.kill(process.pid, result.signal);
else process.exitCode = result.status;
