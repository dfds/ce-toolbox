import subprocess
import json
from pathlib import Path
#1 Run Prowler and generate a report (CSV or JSON) file for each account
  # List non-compliant Dynamodb tables along with AWS account ID
#2 Consume CSV file with the list of capability-accounts with their AWS account ID and Capability members
#3 For each entry from #1 get the capability members and send an email using email template
    # Email template:
    # You are receiving this email because you are member of capability <Name here> which contains the below DynamoDB tables that
    # don't have Point-in-time recovery setting enabled:
    # <List of tables here>
    # Please enable the Point-in-time recovery setting before <Date here>. It is crucial to have backups of the databases to prevent data loss that may be caused of physical or logical errors.
    # Here you can find instructions on how to enable the Point-in-time recovery setting for DynamoDB tables:
    # Using Console: https://amazon-dynamodb-labs.com/hands-on-labs/backups/pitr-backup.html#how-to-enable-pitr
    # Using Terraform: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#point_in_time_recovery

def generate_report():
    exit_code = subprocess.call('./fetch_backup_stats.sh')
    print(exit_code)

def parse_report():
    # For each json file:'
    returned_report_items = [] # list of objects ?

    directory = './output/'
    files = Path(directory).glob('*.json')

    for file in files:
        temp_account_id = ''
        temp_non_compliant_dynamodb_list = []
        temp_obj = {} # Object for each file
        with open(file) as json_file:
            json_data = json.load(json_file)
            if len(json_data) != 0:
                for v in json_data: # For each dynamodb entry
                    if v['Status'] == 'FAIL':
                        #print('Table name: ' + v['ResourceId'], v['AccountId'])
                        temp_table_details = v['ResourceId']  + " " + v['Region']
                        temp_non_compliant_dynamodb_list.append(temp_table_details)
                        temp_account_id = v['AccountId']
                if temp_non_compliant_dynamodb_list != []:
                    temp_obj = {
                        'account_id': temp_account_id,
                        'resource_list': temp_non_compliant_dynamodb_list
                    }
                    returned_report_items.append(temp_obj)
    return returned_report_items

def get_capability(capability_list, account_id):
    for c in capability_list:
       if c['awsAccountId'] == account_id:
           return c


def produce_values_file(report_items, caps_source_file):
    entries = []
    non_cap_account = []
    with open(caps_source_file) as json_file:
        json_data = json.load(json_file)
        for it in report_items:
            cap = get_capability(json_data, it['account_id'])
            if cap == None:
                non_cap_account.append(it)
                continue
            value_list_item = {
                'name': cap['name'],
                'emails': cap['emails'],
                'values': {
                    'affectedResources': it['resource_list']
                }
            }
            entries.append(value_list_item)

        vars_file = {
            'title': 'Capability [{{ .Vars.Name }}] - DynamoDB tables not compliant bla!',
            'entries': entries
        }

        with open("vars.json", "w") as outfile:
            json.dump(vars_file, outfile)

        with open("dynamodb-tables-in-non-cap-accounts.json", "w") as outfile:
            json.dump(non_cap_account,outfile)

def main():
    dynamodb_list = parse_report()

    produce_values_file(dynamodb_list,'./assets/caps.json')
if __name__ == "__main__":
    main()