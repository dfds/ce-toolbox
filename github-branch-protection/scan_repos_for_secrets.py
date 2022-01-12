#!/usr/bin/env python3
import os
import subprocess

from github.branch_protection import Repo

GITHUB_ORGANIZATION = 'dfds'


def main():
    """
    Requirements that needs to be pre-installed:
    - git
    - https://github.com/zricethezav/gitleaks
    - https://github.com/kootenpv/gittyleaks
    """

    token: str = os.environ.get('GITHUB_OAUTH2_TOKEN')
    repo: Repo = Repo(token=token, owner=GITHUB_ORGANIZATION)
    repo_list: list = repo.get_private_repos()

    tmp_clone_dir: str = "/tmp/gitleaks/code"
    tmp_report_dir: str = "/tmp/gitleaks/logs"
    os.makedirs(tmp_clone_dir, exist_ok=True)
    os.makedirs(tmp_report_dir, exist_ok=True)

    for repos in repo_list:
        for r in repos:
            name: str = r.get('name')
            ssh_url: str = r.get('ssh_url')
            git_clone: str = f'git clone {ssh_url}'
            subprocess.run(git_clone.split(' '), cwd=tmp_clone_dir)
            git_leaks: str = f'gitleaks detect -r {tmp_report_dir}/{name}.json'
            output = subprocess.run(git_leaks.split(' '), cwd=f'{tmp_clone_dir}/{name}', capture_output=True)
            if output.returncode == 0:  # We only want to store a report when issues are found
                subprocess.run(["rm", "-f", f'{tmp_report_dir}/{name}.json'])


if __name__ == "__main__":
    main()
