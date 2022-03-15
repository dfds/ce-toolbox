# Housekeeping

This repository contains various Python modules and utility scripts for housekeeping tasks.

## Development practices

Guidelines for developing and testing the housekeeping scripts.

### Virtual environment

```bash
poetry install
poetry shell
```

### Testing and code coverage

There are currently no test cases, but this is how testing would look like:

```bash
poetry run flake8
```

## Usage

### Set branch protection rules

```bash
poetry install --no-dev
poetry shell
export GITHUB_OAUTH2_TOKEN=<REDACTED>
./manage_branch_protection_rules.py
```
