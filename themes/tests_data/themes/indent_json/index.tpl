{{ range .Requests -}}
{{ if eq .PayloadType "raw" -}}
------------
{{ .Name }}
{{ indentJSON .PayloadRaw }}
{{- end }}
{{- end }}

{{ range .Folders -}}
{{ range .Requests -}}
{{ if eq .PayloadType "raw" -}}
------------
{{ .Name }}
{{ indentJSON .PayloadRaw }}
{{- end }}
{{- end }}
{{- end }}