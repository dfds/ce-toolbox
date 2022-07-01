#!/bin/bash

set -euo pipefail

# Get all namespaces, but exclude kube-system and monitoring plus namespaces that have exceptions
NAMESPACES=$(kubectl get ns --no-headers=true | awk '!/kube-system|monitoring|traefik|kyverno|fluentd|flux-system|velero|atlantis|upbound-system|terminal-cargoallocati-zamrl/{print $1}')

# Apply to all namespaces
for NS in $NAMESPACES
do
	echo "Applying limitrange in namespace: $NS"
	kubectl -n $NS apply -f limitrange.yaml
done
