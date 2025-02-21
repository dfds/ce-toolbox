# Generate a service connection config
This tool allows for easy creation of a kubernetes configuration file based on a pre-defined template.
After creation of the file the tool automatically provisions the kubernetes configuration inside of the AWS Systems Manager Parameter Store of the AWS account specified as the Root ID in the script

## Prerequisites
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html), version 2.7 or greater
* [Go AWS SSO](https://github.com/theurichde/go-aws-sso), version 1.3.0 or greater

**NOTE:**  
This approach only sets `kubectl` to use the configuration file for the current terminal session you are running. Re-run step five or make sure to export the full path of the `KUBECONFIG` inside your shell of choices rc file (.bashrc / .zshrc)

## How to use

### Steps


1. `git clone` the repository to your local machine
`git clone git@github.com:dfds/ce-toolbox.git`
2. Ensure your KUBECONFIG environment variable is set to the right path, and navigate to the folder containing the kube-config-generator.sh script.
``` bash
export KUBECONFIG=~/.kube/hellman-saml.config
cd ce-toolbox/k8s-service-account-config-to-ssm
```
3. You need to go-aws-sso -p saml login before the next step (as `dfds-oxygen / CloudAdmin`)
``` bash
go-aws-sso -p saml --persist --start-url https://dfds.awsapps.com/start --region eu-west-1
```
4. Execute the command from Harald, looking something like this:
``` bash
ROOT_ID=<name-of-the-new-capability-context> ACCOUNT_ID=xxxxxxxxxxxxxxxx ./kube-config-generator.sh
```
5. How to verify it was successful:
``` bash
export CAPABILITY_AWS_ACCOUNT_ID=xxxxxxxxxxxxxxxx              # -- Where xxxxxxxxxxxxxxxx is the AWS Account ID, as per above
aws sts assume-role --role-arn "arn:aws:iam::${CAPABILITY_AWS_ACCOUNT_ID}:role/OrgRole" --role-session-name "kube-config-check" --duration-seconds 900; \
aws ssm get-parameters --names "/managed/deploy/kube-config" --region eu-central-1
```
