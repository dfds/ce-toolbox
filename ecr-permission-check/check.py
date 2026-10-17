import boto3
import json
from botocore.exceptions import ClientError

ecr = boto3.client('ecr')

PERFORM_UPDATES = False

DEFAULT_ACTIONS = [
        "ecr:BatchCheckLayerAvailability",
        "ecr:BatchGetImage",
        "ecr:GetDownloadUrlForLayer",
        "ecr:ListImages",   
]

DEFAULT_PRINCIPALS = [
    "arn:aws:iam::123412341234:root",
    "arn:aws:iam::456745674567:root",
]

def ensure_policy(name):
    """Ensure the ECR repository has the default policy applied."""

    try:
        policy = ecr.get_repository_policy(repositoryName=name)
        policy_text = json.loads(policy['policyText'])
        
        updated = False  # are we updating existing policy
        existing = False # do we have to add a new default policy
        
        for s in policy_text.get('Statement', []):
            if s['Effect'] == 'Allow' and s['Sid'] == 'Allow pull from all':
                existing = True

                principal = s.get('Principal', {})
                if isinstance(principal, str):
                    print(f"Principal is a string: {principal}")
                    print(f"\n>> Skipping {name}\n")
                    continue
                principals = principal.get('AWS', [])
                if isinstance(principals, str):
                    principals = [principals]
                principals_match = set(principals) >= set(DEFAULT_PRINCIPALS)
                if not principals_match:
                    s['Principal']['AWS'] = list(set(principals).union(set(DEFAULT_PRINCIPALS)))
                    updated = True
                    
                actions = s.get('Action', [])
                actions_match = set(actions) >= set(DEFAULT_ACTIONS)
                if not actions_match:
                    s['Action'] = list(set(actions).union(set(DEFAULT_ACTIONS)))
                    updated = True

        if not existing:
            print(f"Adding default rule to policy for {name}")
            policy_text['Statement'].append({
                "Sid": "Allow pull from all",
                "Effect": "Allow",
                "Principal": {"AWS": DEFAULT_PRINCIPALS},
                "Action": DEFAULT_ACTIONS
            })
            updated = True

        if updated:
            print(f"Suggested new policy: {json.dumps(policy_text, indent=2)}")
            if PERFORM_UPDATES:
                print(f"Updating policy for {name}")
                ecr.set_repository_policy(
                    repositoryName=name,
                    policyText=json.dumps(policy_text)
                )
        else:
            print(f"Policy for {name} is already up to date")

    except ClientError as e:
        if e.response['Error']['Code'] == 'RepositoryPolicyNotFoundException':
            print(f"RepositoryPolicyNotFoundException: Setting new policy for {name}")
            new_policy = {
                "Version": "2012-10-17",
                "Statement": [{
                    "Sid": "Allow pull from all",
                    "Effect": "Allow",
                    "Principal": {"AWS": DEFAULT_PRINCIPALS},
                    "Action": DEFAULT_ACTIONS
                }]
            }
            
            if PERFORM_UPDATES:
                ecr.set_repository_policy(
                    repositoryName=name,
                    policyText=json.dumps(new_policy)
                )
        else:
            print(f"Error with {name}: {e}")

def main():
    paginator = ecr.get_paginator('describe_repositories')
    for page in paginator.paginate():
        for repo in page['repositories']:
            repo_name = repo['repositoryName']
            ensure_policy(repo_name)

if __name__ == "__main__":
    main()

