#!/usr/bin/env python3
 
import os
import argparse

# Get some args and stuff
parser = argparse.ArgumentParser(description="Orphan some Terraform resources from subdirectories in the supplied directory")
parser.add_argument('path', action = 'store', type = str, default='.', help = 'The Directory to process. Defaults to current directory')
parser.add_argument('resource', action = 'store', type = str, help = 'The terraform resource to orphan')
args = parser.parse_args()

# Get the supplied directory and strip trailing slash. Note: will not work on / because of this
current_dir = args.path.rstrip("/")
resource = args.resource
print("Current directory: " + current_dir)

sub_directories = []

# Cycle files in current dir and append directories to list
for file in os.listdir(current_dir):
    if os.path.isdir(os.path.join(current_dir, file)):
        #print("Found " + file)
        sub_directories.append(current_dir + "/" + file)

# Cycle through directories
for dir in sub_directories:
    print("Processing: " + dir)
    # Check terragrunt.hcl file exists
    if os.path.exists(dir + "/terragrunt.hcl"):
        # FOR TESTING to avoid inadvertently breaking something. Remove this if statement for final version. Changes are only ran against
        # this directory until removed
        if (dir == current_dir + "/ebc33fe5-48f6-4721-8554-5550378a2c8e"):
            print("HCL found!")
            # Run terragrunt rm on a supplied resource
            stream = os.popen('/usr/local/bin/terragrunt show --terragrunt-working-dir ' + dir + '| grep ' + resource)
            output = stream.readlines()
            if (len(output) > 0):
                print("Found " + resource + "!")
                stream = os.popen('terragrunt state rm ' + resource + '--terragrunt-working-dir ' + dir)
                output = stream.read()
                print(stream)
            else:
                print("No resource " + resource + " found!")



            