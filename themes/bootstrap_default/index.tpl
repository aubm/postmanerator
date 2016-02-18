<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Name }}</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.1.0/styles/monokai.min.css" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.1.0/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
</head>
<body>
<div class="container">
<h1>{{ .Name }}</h1>

<h2>General notes</h2>

{{ markdown .Description }}

<h2>API details</h2>

{{ range .Folders }}
<div class="endpoints-group">
    <h3>{{ .Name }}</h3>

    <div>{{ markdown .Description }}</div>

    {{ range .Order }}

        {{ with $req := findRequest . }}
        <div class="request">

            <h4>{{ $req.Name }}</h4>

            <div>{{ markdown $req.Description }}</div>

            <h5>Request</h5>

            <table class="table table-bordered">
                <tr><th>Method</th><td>{{ .Method }}</td></tr>
                <tr><th>URL</th><td>{{ .URL }}</td></tr>
            </table>

            {{ with $res := findResponse $req "default" }}

                <h5>Response</h5>

                <table class="table table-bordered">
                    <tr><th>Code</th><td>{{ $res.ResponseCode.Code }}</td></tr>
                    <tr><th>Status</th><td>{{ $res.ResponseCode.Name }}</td></tr>
                </table>

                {{ with $example := $res.Request.Data }}
                    <h6>Example :</h6>

                    <pre><code>{{ $example }}</code></pre>
                {{ end }}

            {{ end }}

        </div>
        {{ end }}

    {{ end }}

</div>
{{ end }}
</div>
</body>
</html>
