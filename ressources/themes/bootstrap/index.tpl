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
    <style>{{ template "custom.css" }}</style>
</head>
<body data-spy="scroll" data-target=".scrollspy">
<div id="sidebar-wrapper">
    {{ template "menu.tpl" . }}
</div>
<div id="page-content-wrapper">
    <div class="container-fluid">
        <div class="row">
            <div class="col-lg-12">
                <h1>{{ .Name }}</h1>

                <h2 id="doc-general-notes">General notes</h2>

                {{ markdown .Description }}

                {{ with $structures := .Structures }}
                <h2 id="doc-api-structures">API structures</h2>

                {{ range $structures }}

                    <h3 id="struct-{{ .Name }}">{{ .Name }}</h3>

                    <p>{{ .Description }}</p>

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

                <h2 id="doc-api-detail">API detail</h2>

                {{ range .Folders }}
                <div class="endpoints-group">
                    <h3 id="folder-{{ slugify .Name }}">{{ .Name }}</h3>

                    <div>{{ markdown .Description }}</div>

                    {{ range .Order }}

                        {{ with $req := findRequest $.Requests . }}
                        <div class="request">

                            <h4 id="request-{{ slugify $req.Name }}">{{ $req.Name }}</h4>

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

                                {{ with $example := $res.Text }}
                                    {{ $exampleID := randomID }}
                                    <button class="btn btn-default" onclick="showOrHideBlock('example_{{$exampleID}}')">Show example</button>

                                    <pre id="example_{{$exampleID}}" style="display:none;"><code>{{ indentJSON $example }}</code></pre>
                                {{ end }}

                            {{ end }}

                        </div>
                        {{ end }}

                    {{ end }}

                </div>
                {{ end }}
            </div>
        </div>
    </div>
</div>

<script src="https://code.jquery.com/jquery-2.2.2.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>
<script>
function showOrHideBlock(blockId) {
    var block = $('#' + blockId);
    if (block.css('display') === 'none') {
        block.show();
    } else {
        block.hide();
    }
}
</script>
</body>
</html>
