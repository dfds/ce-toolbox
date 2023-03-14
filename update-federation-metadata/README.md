# :warning: Not required after 1st May 2023 :warning:

Please note that this tool will no longer be required after 1st May 2023 but may serve as a useful 
reference. It will be removed after this date.


# Update federation metadata
This repository contains a script for bulk update federation metadata for all accounts.

# Quick start
Login with go-aws-sso

```bash
go-aws-sso -p saml
export AWS_PROFILE=saml
```

Execute update script
```bash
bash run.sh
```

It will fail on some legacy accounts, it self and some others.


# If certificate is already expired
Root login with master account in AWS management console and update the federation metadata. Then login with saml2aws and run the script.
