package bash

const tmpl = `
set -e

{{ range . }}
{{ echo . }}
{{ run . }}
{{ end }}
`
