#!/usr/bin/env python3
"""One-shot WeKnora -> SemiClaw rebrand for text files.

Preserves external wire values per plan:
  - https://weknora.weixin.qq.com
  - "weknoracloud" provider string literals
"""
from __future__ import annotations

import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

SKIP_DIRS = {
    ".git",
    "node_modules",
    "frontend/dist",
    "frontend/node_modules",
    ".cursor",
    "vendor",
}

SKIP_EXTENSIONS = {
    ".png",
    ".jpg",
    ".jpeg",
    ".gif",
    ".ico",
    ".webp",
    ".woff",
    ".woff2",
    ".ttf",
    ".eot",
    ".pdf",
    ".zip",
    ".tar",
    ".gz",
    ".bin",
    ".exe",
    ".dll",
    ".so",
    ".dylib",
    ".pb",
}

REPLACEMENTS: list[tuple[str, str]] = [
    ("github.com/Tencent/WeKnora", "github.com/vagawind/semiclaw"),
    ("NannaOlympicBroadcast/WeKnoraMCP", "vagawind/semiclaw-mcp"),
    ("Tencent/WeKnora", "vagawind/semiclaw"),
    ("wechatopenai/weknora", "vagawind/semiclaw"),
    ("WeKnoraCloud", "SemiClawCloud"),
    ("WeKnora Lite", "SemiClaw Lite"),
    ("WeKnora", "SemiClaw"),
    ("WeknoraLite", "SemiClawLite"),
    ("weKnora", "semiClaw"),
    ("WEKNORA", "SEMICLAW"),
    ("we_knora", "semi_claw"),
    ("weknora", "semiclaw"),
]

# Placeholders must NOT contain "weknora" or they get corrupted during replacement.
PROTECT_URL = "___WKC_EXTERNAL_URL___"
PROTECT_WIRE = "___WKC_WIRE_PROVIDER___"


def should_process(path: Path) -> bool:
    if path.suffix.lower() in SKIP_EXTENSIONS:
        return False
    for part in path.parts:
        if part in SKIP_DIRS:
            return False
    return path.is_file()


def protect(content: str) -> str:
    content = content.replace("https://weknora.weixin.qq.com", PROTECT_URL)
    content = content.replace('"weknoracloud"', f'"{PROTECT_WIRE}"')
    content = content.replace("'weknoracloud'", f"'{PROTECT_WIRE}'")
    content = content.replace("`weknoracloud`", f"`{PROTECT_WIRE}`")
    return content


def restore(content: str) -> str:
    content = content.replace(PROTECT_URL, "https://weknora.weixin.qq.com")
    content = content.replace(PROTECT_WIRE, "weknoracloud")
    return content


def transform(content: str) -> str:
    content = protect(content)
    for old, new in REPLACEMENTS:
        content = content.replace(old, new)
    return restore(content)


def main() -> int:
    changed = 0
    for dirpath, dirnames, filenames in os.walk(ROOT):
        dirnames[:] = [d for d in dirnames if d not in SKIP_DIRS]
        for name in filenames:
            path = Path(dirpath) / name
            if not should_process(path):
                continue
            try:
                original = path.read_text(encoding="utf-8")
            except (UnicodeDecodeError, OSError):
                continue
            updated = transform(original)
            if updated != original:
                path.write_text(updated, encoding="utf-8")
                changed += 1
                print(f"updated: {path.relative_to(ROOT)}")
    print(f"\nDone. {changed} files updated.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
