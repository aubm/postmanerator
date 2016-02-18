# {{ .Name }}

## General notes

{{ .Description }}

## API Details

{{ range .Folders }}
### {{ .Name }}

{{ .Description }}

{{ range .Order }}

{{ with $req := findRequest . }}

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

{{ with $example := $res.Request.Data }}
**Example** :

```
{{ $example }}
```
{{ end }}

{{ end }}

{{ end }}

{{ end }}

{{ end }}
