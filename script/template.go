package script

// BashTemplate is the template to render a Bash script.
const BashTemplate = `set -e
{{ range . }}
{{ echo . }}
{{ run . }}
{{ end }}
`
