{{ range .Requests }}
{{ curlSnippet . }}
{{ end }}

{{ range .Folders }}
{{ range .Requests }}
{{ curlSnippet . }}
{{ end }}
{{ end }}