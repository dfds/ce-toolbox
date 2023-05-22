#!/bin/bash
# Sign into a sufficiently priviledged account and run this script to fetch the backup stats for all accounts in the organization
# Requires prowler to be installed and configured with the saml profile
# Requires jq to be installed
# Requires parallel to be installed
# Requires aws cli to be installed and configured with the saml profile
# Requires the OrgRole role to be created in each account
# example:
#   go-aws-sso -p saml --persist
#   ./fetch_backup_stats.sh

set -e # might need to remove this line, cant remember if prowler fails with a non-zero exit code on recoverable errors


process_account() {
  accountId="$1"
  arnStart='arn:aws:iam::'
  arnEnd=':role/OrgRole'
  echo "prowler aws --profile=saml --role $arnStart$accountId$arnEnd"
  prowler aws --profile=saml --role $arnStart$accountId$arnEnd -c dynamodb_tables_pitr_enabled
}

export -f process_account

# Run the function in parallel for active accounts with a progress bar
aws --profile=saml organizations list-accounts | jq '.Accounts[] | select(.Status == "ACTIVE") | .Id | tonumber' | parallel --bar -j 3 process_account
