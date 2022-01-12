#!/usr/bin/env python3
import logging
import os
import subprocess

from github.branch_protection import Repo

GITHUB_ORGANIZATION = 'dfds'


def main():
    """
    Requirements that needs to be pre-installed:
    - git
    - https://github.com/zricethezav/gitleaks
    """

    token: str = os.environ.get('GITHUB_OAUTH2_TOKEN')
    repo: Repo = Repo(token=token, owner=GITHUB_ORGANIZATION)
    repo_list: list = repo.get_all_repos()

    _clone_dir: str = "/tmp/code"
    _gitleaks_report_dir: str = "/tmp/logs/gitleaks"
    os.makedirs(_clone_dir, exist_ok=True)
    os.makedirs(_gitleaks_report_dir, exist_ok=True)

    for repos in repo_list:
        for r in repos:
            name: str = r.get('name')
            ssh_url: str = r.get('ssh_url')
            logging.info(f'Cloning repository {name}')
            _cmd_git_clone: str = f'git clone {ssh_url}'
            subprocess.run(_cmd_git_clone.split(' '), cwd=_clone_dir)
            logging.info(f'Scanning {name} with gitleaks')
            _cmd_gitleaks: str = f'gitleaks detect -r {_gitleaks_report_dir}/{name}.json'
            _gitleaks_output = subprocess.run(_cmd_gitleaks.split(' '), cwd=f'{_clone_dir}/{name}', capture_output=True)
            if _gitleaks_output.returncode == 0:  # Only store a report when issues are found
                subprocess.run(["rm", "-f", f'{_gitleaks_report_dir}/{name}.json'])


if __name__ == "__main__":
    main()
