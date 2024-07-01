#!/usr/bin/env python3

import os
import re
from typing import Dict, Optional, Tuple

import requests


def is_stable_version(version: str) -> bool:
    """Check if the version string contains only numbers and dots."""
    return bool(re.match(r"^[\d.]+$", version))


def get_latest_release(repo: str, use_releases: bool) -> Optional[str]:
    """Fetch the latest release version from GitHub."""
    if repo == "golang/go":
        url = "https://go.dev/VERSION?m=text"
        try:
            response = requests.get(url)
            response.raise_for_status()
            version = response.text.split("\n")[0].strip().lstrip("go")
            return version
        except requests.RequestException as e:
            print(f"Error fetching Go version: {e}")
            return None

    base_url = f"https://api.github.com/repos/{repo}"
    url = f"{base_url}/releases/latest" if use_releases else f"{base_url}/tags"

    try:
        response = requests.get(url)
        response.raise_for_status()
        data = response.json()

        if use_releases:
            return data["tag_name"].lstrip("v")
        else:
            stable_tags = [
                tag["name"].lstrip("v")
                for tag in data
                if is_stable_version(tag["name"].lstrip("v"))
            ]
            return stable_tags[0] if stable_tags else None
    except requests.RequestException:
        return None


def update_file(filepath: str, update_func):
    """Update a file using the provided update function."""
    with open(filepath, "r", newline="") as f:
        lines = f.readlines()

    updated_lines = update_func(lines)

    with open(filepath, "w", newline="") as f:
        f.writelines(updated_lines)


def update_versions_md(name: str, latest_version: str):
    """Update the versions.md file with the latest version."""

    def update_func(lines):
        return [
            f"{name}={latest_version}\n" if line.startswith(f"{name}=") else line
            for line in lines
        ]

    update_file("devel/versions.md", update_func)


def update_dockerfiles(name: str, latest_version: str):
    """Update all Dockerfiles with the latest version."""

    def update_func(lines):
        return [
            (
                f'ARG {name}="{latest_version}"\n'
                if line.startswith(f"ARG {name}=")
                else line
            )
            for line in lines
        ]

    for filename in os.listdir("devel"):
        if filename.endswith(".Dockerfile"):
            update_file(os.path.join("devel", filename), update_func)


def update_version(name: str, current_version: str, repo: str, use_releases: bool):
    """Update the version for a specific package."""
    try:
        latest_version = get_latest_release(repo, use_releases)
        if latest_version is None:
            print(f"Couldn't retrieve latest version for {repo}")
            return

        if latest_version != current_version:
            update_versions_md(name, latest_version)
            update_dockerfiles(name, latest_version)
            print(f"Updated {name} from {current_version} to {latest_version}")
        else:
            print(f"{name} already on latest release {current_version}")
    except Exception as e:
        print(f"Error updating {name}: {str(e)}")


def main():
    # Dictionary mapping ARG names to their respective GitHub repos and release type
    repos: Dict[str, Tuple[str, bool]] = {
        "RUSTUP_VERSION": ("rust-lang/rustup", False),
        "RUST_VERSION": ("rust-lang/rust", True),
        "GO_VERSION": ("golang/go", False),
        "DENO_VERSION": ("denoland/deno", True),
        "NVM_VERSION": ("nvm-sh/nvm", True),
        "NODE_VERSION": ("nodejs/node", True),
        "ZIG_VERSION": ("ziglang/zig", True),
        "TINI_VERSION": ("krallin/tini", True),
        "CONTAINERD_VERSION": ("containerd/containerd", True),
        "RUNC_VERSION": ("opencontainers/runc", True),
        "CNI_PLUGINS_VERSION": ("containernetworking/plugins", True),
        "NERDCTL_VERSION": ("containerd/nerdctl", True),
    }

    # Read current versions
    with open("devel/versions.md", "r") as f:
        current_versions = dict(line.strip().split("=") for line in f)

    # Update each version
    for name, (repo, use_releases) in repos.items():
        update_version(name, current_versions[name], repo, use_releases)


if __name__ == "__main__":
    main()
