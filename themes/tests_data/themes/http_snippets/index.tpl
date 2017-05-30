{{ range .Requests }}
{{ httpSnippet . }}
{{ end }}

{{ range .Folders }}
{{ range .Requests }}
{{ httpSnippet . }}
{{ end }}
{{ end }}