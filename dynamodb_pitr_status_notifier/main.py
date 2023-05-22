import subprocess
import json
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


    returned_obj = {}
    with open("./output/prowler-output-<account-id-here>-20230522154659.json") as json_file:
        json_data = json.load(json_file)

        for v in json_data:
            if v['Status'] == 'FAIL':

                print('Table name: ' + v['ResourceId'], v['AccountId'])
                returned_obj = {'AccountID':v['AccountId'], 'TableNames':}


def main():
   # generate_report()
    parse_report()
if __name__ == "__main__":
    main()