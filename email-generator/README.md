# email-generator

## Build

**linux/macos**: `go build -o email-gen main.go`

**windows**: `go build -o email-gen.exe main.go`

## Usage

You can use one of the predefined templates in ./template or create your own. It uses the [Go templating language](https://pkg.go.dev/text/template).

*vars.json* contains variables that your template can make use of. The program treats vars.json as a map, so anything goes as long as the value is a string.

To use the *kiam_migration_august.tpl* template:

linux/macos: `./email-gen template/kiam_migration_august.tpl`
windows: `email-gen.exe template/kiam_migration_august.tpl`

*kiam_migration_august.tpl* expects two values in vars.json.

```json
{
  "RootId": "sandbox-emcla-pmyxn",
  "Count": "11"
}
```
