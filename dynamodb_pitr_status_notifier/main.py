import subprocess
import json
from pathlib import Path
import boto3


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
                        table_name = v['ResourceId']

                        if table_name != "terraform-locks" or "terraform-locks" not in table_name:
                            temp_table_details = table_name + " in " + v['Region']
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


def enable_pitr(report_items_list):
    sts_client = boto3.client("sts")

    for item in report_items_list:
        account_id = item["account_id"]
        resp = sts_client.assume_role(
            RoleArn="arn:aws:iam::{}:role/OrgRole".format(account_id),
            RoleSessionName="enable-pitr"
        )

        creds = resp["Credentials"]

        for table in item["resource_list"]:
            session = boto3.Session(
                aws_access_key_id=creds['AccessKeyId'],
                aws_secret_access_key=creds['SecretAccessKey'],
                aws_session_token=creds['SessionToken'],
                region_name=table.split()[2]
            )

            dynamodb_client = session.client("dynamodb")

            table_name = table.split()[0]

            pitr_status = dynamodb_client.describe_continuous_backups(
                TableName=table_name
            )["ContinuousBackupsDescription"]["PointInTimeRecoveryDescription"]["PointInTimeRecoveryStatus"]

            if pitr_status == "DISABLED":

                dynamodb_client.update_continuous_backups(
                    TableName=table_name,
                    PointInTimeRecoverySpecification={
                        'PointInTimeRecoveryEnabled': True
                    }
                )
                print("Enabled PITR for {} in {}".format(table, account_id))


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
                        emails = legacy_cap['members']
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
    enable_pitr(dynamodb_list)

    # produce_values_file(dynamodb_list, './assets/caps.json', './assets/legacy_caps.json')


if __name__ == "__main__":
    main()
