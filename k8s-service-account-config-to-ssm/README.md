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

### Running in an environment without a browser

Make sure to authenticate as any role before running the script, otherwise the script will hang silently without redirecting to a browser to complete the authentication.

### Steps

1. `git clone` the repository to your local machine
2. `cd` to the *ded-toolbox/k8s-service-account-config-to-ssm* directory
3. Get and set the local variables information from **ded-aesir** slack channel
4. Execute the `./kube-config-generator.sh` script

The script requires two environment variables set.  
`ROOT_ID`: Can be found from the **capabilityRootId** field from Harald-notify app in **ded-aesir** slack channel.

`ACCOUNT_ID`: Can be found after the **Tax settings for AWS Account** field from Harald-notify app in **ded-aesir** slack channel. 

The syntax for running the script is:  
`ROOT_ID=YOUR_CAPABILITY_ROOT_ID ACCOUNT_ID=YOUR_CAPABILITY_AWS_ACCOUNT_ID ./kube-config-generator.sh`

Sample script execution:
``` bash
ROOT_ID=capabilityplayground-312312 \
ACCOUNT_ID=123456789123 \
./kube-config-generator.sh
```

