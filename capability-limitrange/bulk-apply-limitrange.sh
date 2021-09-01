#!/bin/bash

set -euo pipefail

# Get all namespaces, but exclude kube-system and monitoring plus namespaces that have exceptions
NAMESPACES=$(kubectl get ns --no-headers=true | awk '!/kube-system|monitoring|terminal-cargoallocati-zamrl/{print $1}')

# Apply to all namespaces
for NS in $NAMESPACES
do
	echo "Applying limitrange in namespace: $NS"
	kubectl -n $NS apply -f limitrange.yaml
done

# Apply exceptions
for file in `ls exceptions | grep .yaml`
do
	echo "Applying: $file"
	kubectl apply -f exceptions/$file
done