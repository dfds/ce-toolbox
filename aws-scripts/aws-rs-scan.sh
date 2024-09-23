#!/bin/bash

### This script can be used to obtain a list of resources across multiple AWS accounts and regions.
### As of now it is scanning for EC2 instances, EKS clusters, ECS clusters, Lambda functions and IAM roles.
### The script assumes a role in each account and region to scan for resources.
### To run the script you will need to have the AWS CLI installed and configured with the necessary permissions.

# Setting prerequisites:
ROLE_NAME="OrgRole" # This is optional, you can omit this if you are using AWS SSO and Master account role is already assumed
SESSION_NAME="ScanResourcesSession"
OUTPUT_FILE="aws_resource_scan_results.csv" # This name can be changed to something else
REGIONS=("eu-central-1" "eu-west-1") # You can add more regions here, depending on your requirements

# CSV header
echo "AccountId,Region,ResourceType,ResourceId,ResourceDetails" > $OUTPUT_FILE

# Get accounts from AWS Organizations
accounts=$(aws organizations list-accounts --query 'Accounts[].Id' --output text)

# Function to assume role in an account
assume_role() {
  local account_id=$1
  creds=$(aws sts assume-role --role-arn arn:aws:iam::$account_id:role/$ROLE_NAME \
    --role-session-name $SESSION_NAME --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' --output text)
  
  export AWS_ACCESS_KEY_ID=$(echo $creds | awk '{print $1}')
  export AWS_SECRET_ACCESS_KEY=$(echo $creds | awk '{print $2}')
  export AWS_SESSION_TOKEN=$(echo $creds | awk '{print $3}')
}

# Function to write to CSV file
write_to_csv() {
  local account_id=$1
  local region=$2
  local resource_type=$3
  local resource_id=$4
  local resource_details=$5
  echo "$account_id,$region,$resource_type,$resource_id,\"$resource_details\"" >> $OUTPUT_FILE
}

# Function to scan for EC2 instances
scan_ec2() {
  local account_id=$1
  local region=$2
  echo "Scanning EC2 instances in $region..."
  ec2_instances=$(aws ec2 describe-instances --region $region --query 'Reservations[*].Instances[*].[InstanceId,InstanceType,State.Name]' --output json)

  # Loop through each EC2 instance and write to CSV
  echo "$ec2_instances" | jq -c '.[] | .[]' | while read instance; do
    instance_id=$(echo $instance | jq -r '.[0]')
    instance_type=$(echo $instance | jq -r '.[1]')
    instance_state=$(echo $instance | jq -r '.[2]')
    write_to_csv "$account_id" "$region" "EC2" "$instance_id" "$instance_type $instance_state"
  done
}

# Function to scan for EKS clusters
scan_eks() {
  local account_id=$1
  local region=$2
  echo "Scanning EKS clusters in $region..."
  eks_clusters=$(aws eks list-clusters --region $region --output json)

  # Loop through each EKS cluster and write to CSV
  echo "$eks_clusters" | jq -r '.clusters[]' | while read cluster_name; do
    write_to_csv "$account_id" "$region" "EKS" "$cluster_name" ""
  done
}

# Function to scan for ECS clusters
scan_ecs() {
  local account_id=$1
  local region=$2
  echo "Scanning ECS clusters in $region..."
  ecs_clusters=$(aws ecs list-clusters --region $region --output json)

  # Loop through each ECS cluster and write to CSV
  echo "$ecs_clusters" | jq -r '.clusterArns[]' | while read cluster_arn; do
    write_to_csv "$account_id" "$region" "ECS" "$cluster_arn" ""
  done
}

# Function to scan for Lambda
scan_lambda() {
  local account_id=$1
  local region=$2
  echo "Scanning Lambda functions in $region..."
  lambda_functions=$(aws lambda list-functions --region $region --query 'Functions[*].[FunctionName,Runtime,LastModified]' --output json)

  # Loop through each Lambda and write to CSV
  echo "$lambda_functions" | jq -c '.[]' | while read function; do
    function_name=$(echo $function | jq -r '.[0]')
    function_runtime=$(echo $function | jq -r '.[1]')
    function_last_modified=$(echo $function | jq -r '.[2]')
    write_to_csv "$account_id" "$region" "Lambda" "$function_name" "$function_runtime $function_last_modified"
  done
}

# Function to scan for IAM roles (global service, no region)
scan_iam_roles() {
  local account_id=$1
  echo "Scanning IAM roles..."
  iam_roles=$(aws iam list-roles --query 'Roles[*].[RoleName,Arn,CreateDate]' --output json)

  # Loop through each IAM role and write to CSV
  echo "$iam_roles" | jq -c '.[]' | while read role; do
    role_name=$(echo $role | jq -r '.[0]')
    role_arn=$(echo $role | jq -r '.[1]')
    role_create_date=$(echo $role | jq -r '.[2]')
    write_to_csv "$account_id" "Global" "IAM Role" "$role_name" "$role_arn $role_create_date"
  done
}

# Main loop to go through all accounts and regions
for account_id in $accounts; do
  echo "========================================================="
  echo "Scanning resources in account: $account_id"
  echo "========================================================="

  # Assume role in the account
  assume_role $account_id

  # Loop through each region and scan resources
  for region in "${REGIONS[@]}"; do
    # Scan EC2 instances
    scan_ec2 $account_id $region

    # Scan EKS clusters
    scan_eks $account_id $region

    # Scan ECS clusters
    scan_ecs $account_id $region

    # Scan Lambda functions
    scan_lambda $account_id $region
  done

  # Scan IAM roles (global)
  scan_iam_roles $account_id

  # Clear credentials after each account scan
  unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
done

echo "Scan completed. Results saved to $OUTPUT_FILE"