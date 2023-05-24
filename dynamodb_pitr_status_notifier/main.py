import subprocess
import json
from pathlib import Path
import boto3


# 1 Run Prowler and generate a report (CSV or JSON) file for each account
# List non-compliant Dynamodb tables along with AWS account ID
# 2 Consume CSV file with the list of capability-accounts with their AWS account ID and Capability members
# 3 For each entry from #1 get the capability members and send an email using email template
# Email template:
# You are receiving this email because you are member of capability <Name here> which contains the below DynamoDB tables that
# don't have Point-in-time recovery setting enabled:
# <List of tables here>
# Please enable the Point-in-time recovery setting before <Date here>. It is crucial to have backups of the databases to prevent data loss that may be caused of physical or logical errors.
# Here you can find instructions on how to enable the Point-in-time recovery setting for DynamoDB tables:
# Using Console: https://amazon-dynamodb-labs.com/hands-on-labs/backups/pitr-backup.html#how-to-enable-pitr
# Using Terraform: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#point_in_time_recovery

def generate_report():
    print("Generating Prowler report...")
    exit_code = subprocess.call('./fetch_backup_stats.sh')
    print("Generating Prowler report. Done.")
    print(exit_code)


def parse_report():
    print("Rendering report content...")
    returned_report_items = []

    directory = './output/'
    files = Path(directory).glob('*.json')

    for file in files:
        temp_account_id = ''
        temp_non_compliant_dynamodb_list = []
        temp_obj = {}  # Object for each file
        with open(file) as json_file:
            json_data = json.load(json_file)
            if len(json_data) != 0:
                for v in json_data:  # For each dynamodb entry
                    if v['Status'] == 'FAIL':
                        temp_table_details = v['ResourceId'] + " in " + v['Region']
                        temp_non_compliant_dynamodb_list.append(temp_table_details)
                        temp_account_id = v['AccountId']
                if temp_non_compliant_dynamodb_list:
                    temp_obj = {
                        'account_id': temp_account_id,
                        'resource_list': temp_non_compliant_dynamodb_list
                    }
                    returned_report_items.append(temp_obj)
    print("Rendering report content. Done.")
    return returned_report_items


def get_capability(capability_list, account_id):
    for c in capability_list:
        if c['awsAccountId'] == account_id:
            return c


def get_account_name(account_id):
    client = boto3.client('organizations')
    response = client.describe_account(AccountId=account_id)
    return response['Account']['Name']

def produce_values_file(report_items, caps_source_file, leg_caps_source_file):

    print("Generating email content...")
    legacy_caps_data = None
    with open(leg_caps_source_file) as json_file:
        legacy_caps_data = json.load(json_file)

    entries = []
    manual_entries = []
    with open(caps_source_file) as json_file:
        json_data = json.load(json_file)
        emails = []
        for report_item in report_items:
            account_id = report_item['account_id']
            cap = get_capability(json_data, account_id)
            if cap is None:
                account_name = get_account_name(account_id)
                for legacy_cap in legacy_caps_data:
                    if legacy_cap['name'] == account_name:
                        emails= legacy_cap['members']
                        break
            else:
                account_name = cap['name']
                emails = cap['emails']
            value_list_item = {
                'name': account_name,
                'emails': emails,
                'values': {
                    'affectedResources': report_item['resource_list']
                }
            }
            if not value_list_item['emails']:
                manual_entries.append(value_list_item)
            else:
                entries.append(value_list_item)

        vars_file = {
            'title': 'Capability [{{ .Vars.Name }}] - DynamoDB tables not compliant!',
            'entries': entries
        }

        with open("vars.json", "w") as outfile:
            json.dump(vars_file, outfile)


        with open("manual_list.json", "w") as outfile:
            json.dump(manual_entries, outfile)

    print("Generating email content. Done.")



def main():

    generate_report()

    dynamodb_list = parse_report()

    produce_values_file(dynamodb_list, './assets/caps.json', './assets/legacy_caps.json')


if __name__ == "__main__":
    main()
