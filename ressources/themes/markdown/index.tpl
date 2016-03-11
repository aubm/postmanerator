# {{ .Name }}

## General notes

{{ .Description }}

{{ with $structures := .Structures }}

## API structures

{{ range $structures }}

### {{ .Name }}

{{ .Description }}

<table class="table table-bordered">
{{ range .Fields }}
<tr>
    <th>{{ .Name }}</th>
    <td>{{ .Type }}</td>
    <td>{{ .Description }}</td>
</tr>
{{ end }}
</table>

{{ end }}

{{ end }}

## API Details

{{ range .Folders }}
### {{ .Name }}

{{ .Description }}

{{ range .Order }}

{{ with $req := findRequest $.Requests . }}

### {{ $req.Name }}

{{ $req.Description }}

#### Request

<table>
    <tr><th>Method</th><td>{{ .Method }}</td></tr>
    <tr><th>URL</th><td>{{ .URL }}</td></tr>
</table>

{{ with $res := findResponse $req "default" }}

#### Response

<table>
    <tr><th>Code</th><td>{{ $res.ResponseCode.Code }}</td></tr>
    <tr><th>Status</th><td>{{ $res.ResponseCode.Name }}</td></tr>
</table>

{{ with $example := $res.Text }}
**Example** :

```
{{ indentJSON $example }}
```
{{ end }}

{{ end }}

{{ end }}

{{ end }}

{{ end }}
