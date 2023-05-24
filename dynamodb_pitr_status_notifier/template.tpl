Hi,

You are receiving this email because you're a member of the {{index .Vars "Name"}} Capability.

The capability contains the following DynamoDB tables, which don't have point-in-time recovery enabled.

Affected DynamoDB resources:
{{ range .Vars.affectedResources }} - {{ . }}
{{end}}

Please enable the point-in-time recovery setting by Wednesday May 31st 12:00 UTC.
It is important to have backups of the databases to prevent data loss that may be caused of physical or logical errors.

Here you can find instructions on how to enable the Point-in-time recovery setting for DynamoDB tables:
Using Console: https://amazon-dynamodb-labs.com/hands-on-labs/backups/pitr-backup.html#how-to-enable-pitr
Using Terraform: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table#point_in_time_recovery

If you have any questions regarding this, please use the #dev-peer-support channel in Slack or on on Microsoft Teams.

Kind regards,
Aleksandra Fromm, Silviu Calin and Sami @ Cloud Engineering Team