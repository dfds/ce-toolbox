# Capability Limit Range

## Overview

This script should be run to apply a limit range to all namespaces.

## How does it work

The script loops over all namespaces, excluding any exceptions, and applies the limitrange.yaml to each namespace. 

The script then looks for .yaml files in the exceptions folder and applies them.

## Making exceptions

Copy the limitrange.yaml and place it in the exceptions folder as namespace-name-here.yaml. Make sure you 
add the namespace name into the file under metadata

## Usage

Ensure your Kubernetes context is set correctly (logged in as CloudAdmin user), make sure that 
`bulk-apply-limitrange.sh` is executable on your system and run the script.