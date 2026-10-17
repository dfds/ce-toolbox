# Check ECR permissions and update

This script requires existing credentials in `~/.aws/credentials` or in environment variables for `shared-prod` in `eu-central-1` in order to run.
It is also highly recommended to run this from a python virtual environment (consider `uv`).

When run, the script will scan all ECR repositories present looking for the rule `Allow pull from all`.
- If the rule does not exist it will be created with `DEFAULT_PRINCIPALS` and `DEFAULT_ACTIONS`.
- If the rule does exist, it will be updated to include `DEFAULT_PRINCIPALS` and `DEFAULT_ACTIONS` while keeping existing rules too.

# Notes

This script is not meant to be run more than once.
A cronjob on the SelfService API will take over, so we don't have to think about this again.

I leave it here as a 'just in cases' measure.
