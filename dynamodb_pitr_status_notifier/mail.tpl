You are receiving this email because you are member of capability <Name here> which contains the below DynamoDB tables that
don't have Point-in-time recovery setting enabled:
    <dynamodb_table_list>
Please enable the Point-in-time recovery setting before <Date here>. It is crucial to have backups of the databases to prevent data loss that may be caused of physical or logical errors.
Here you can find instructions on how to enable the Point-in-time recovery setting for DynamoDB tables:
Using Console: https://amazon-dynamodb-labs.com/hands-on-labs/backups/pitr-backup.html#how-to-enable-pitr
Using Terraform: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#point_in_time_recovery