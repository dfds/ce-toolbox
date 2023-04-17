#!/bin/bash

# Check if variables is there
if [[ -z $ROOT_ID || -z $ACCOUNT_ID ]]; then
    echo 'required environment variables missing'
    exit 1
fi

# Define variables
CAPABILITY_ROOT_ID=$ROOT_ID
NAMESPACE=$CAPABILITY_ROOT_ID
SERVICE_ACCOUNT_NAME="$CAPABILITY_ROOT_ID-vstsuser"
SECRET_NAME="$CAPABILITY_ROOT_ID-vstsuser-token"
KUBE_ROLE="$CAPABILITY_ROOT_ID-fullaccess"
CAPABILITY_AWS_ACCOUNT_ID=$ACCOUNT_ID
CAPABILITY_AWS_ROLE_SESSION="kube-config-paramstore"
AWS_PROFILE="saml"
SSO_START_URL="https://dfds.awsapps.com/start"
SSO_REGION="eu-west-1"

# Define the SAML role to use and login with that role
SAML_ROLE="CloudAdmin"
SAML_ACCOUNT="738063116313"
go-aws-sso assume --role-name $SAML_ROLE --account-id $SAML_ACCOUNT -p $AWS_PROFILE --start-url $SSO_START_URL --region $SSO_REGION

# Generate kube token and config for service account
kubectl create serviceaccount --namespace kube-system $SERVICE_ACCOUNT_NAME

kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: $SECRET_NAME
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: $SERVICE_ACCOUNT_NAME
type: kubernetes.io/service-account-token
EOF

kubectl create rolebinding $SERVICE_ACCOUNT_NAME --role=$KUBE_ROLE --serviceaccount=kube-system:$SERVICE_ACCOUNT_NAME -n $NAMESPACE
KUBE_TOKEN=$(kubectl -n kube-system get secret $SECRET_NAME -o=jsonpath="{.data.token}" | base64 --decode)
KUBE_CONFIG=$(sed "s/KUBE_TOKEN/${KUBE_TOKEN}/g" config.template | sed "s/NAMESPACE_REPLACE/${NAMESPACE}/g")

# go-aws-sso Connection
SAML_ROLE="CloudAdmin"
SAML_ACCOUNT="454234050858"
go-aws-sso assume --role-name $SAML_ROLE --account-id $SAML_ACCOUNT -p $AWS_PROFILE --start-url $SSO_START_URL --region $SSO_REGION

# Assume role
AWS_ASSUMED_CREDS=( $(aws --profile ${AWS_PROFILE} sts assume-role \
    --role-arn "arn:aws:iam::${CAPABILITY_AWS_ACCOUNT_ID}:role/OrgRole" \
    --role-session-name ${CAPABILITY_AWS_ROLE_SESSION} \
    --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' \
    --output text
) )
AWS_ASSUMED_ACCESS_KEY_ID=${AWS_ASSUMED_CREDS[0]}
AWS_ASSUMED_SECRET_ACCESS_KEY=${AWS_ASSUMED_CREDS[1]}
AWS_ASSUMED_SESSION_TOKEN=${AWS_ASSUMED_CREDS[2]}

# Push to AWS parameter store
AWS_ACCESS_KEY_ID=${AWS_ASSUMED_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${AWS_ASSUMED_SECRET_ACCESS_KEY} AWS_SESSION_TOKEN=${AWS_ASSUMED_SESSION_TOKEN} \
    aws ssm put-parameter --name "/managed/deploy/kube-config" --value "$KUBE_CONFIG" --type "SecureString" --overwrite --region eu-central-1
