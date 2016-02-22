<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Name }}</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.1.0/styles/solarized_light.min.css" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.1.0/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
    <style>
        pre code { margin: -9.5px; }
    </style>
    <script>
        function showOrHideBlock(blockId) {
            var block = document.getElementById(blockId);
            var styles = window.getComputedStyle(block);
            if (styles.display === 'none') {
                block.style.display = 'block';
            } else {
                block.style.display = 'none';
            }
        }
    </script>
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
                <tr><th style="width: 20%;">Method</th><td>{{ .Method }}</td></tr>
                <tr><th style="width: 20%;">URL</th><td>{{ .URL }}</td></tr>
            </table>

            {{ with $res := findResponse $req "default" }}

                <h5>Response</h5>

                <table class="table table-bordered">
                    <tr><th style="width: 20%;">Code</th><td>{{ $res.ResponseCode.Code }}</td></tr>
                    <tr><th style="width: 20%;">Status</th><td>{{ $res.ResponseCode.Name }}</td></tr>
                </table>

                {{ with $example := $res.Request.Data }}
                    {{ $exampleID := randomID }}
                    <button class="btn btn-default" onclick="showOrHideBlock('example_{{$exampleID}}')">Show example</button>

                    <pre id="example_{{$exampleID}}" style="display:none;"><code>{{ $example }}</code></pre>
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
