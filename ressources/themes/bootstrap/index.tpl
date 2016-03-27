<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Name }}</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.2.0/styles/darkula.min.css" rel="stylesheet">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.2.0/highlight.min.js"></script>
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

                <h2 id="doc-general-notes">
                    General notes
                    <a href="#doc-general-notes"><i class="glyphicon glyphicon-link"></i></a>
                </h2>

                {{ markdown .Description }}

                {{ with $structures := .Structures }}
                <h2 id="doc-api-structures">
                    API structures
                    <a href="#doc-api-structures"><i class="glyphicon glyphicon-link"></i></a>
                </h2>

                {{ range $structures }}

                    <h3 id="struct-{{ .Name }}">
                        {{ .Name }}
                        <a href="#struct-{{ .Name }}"><i class="glyphicon glyphicon-link"></i></a>
                    </h3>

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

                <h2 id="doc-api-detail">
                    API detail
                    <a href="#doc-api-detail"><i class="glyphicon glyphicon-link"></i></a>
                </h2>

                {{ range .Folders }}
                <div class="endpoints-group">
                    <h3 id="folder-{{ slugify .Name }}">
                        {{ .Name }}
                        <a href="#folder-{{ slugify .Name }}"><i class="glyphicon glyphicon-link"></i></a>
                    </h3>

                    <div>{{ markdown .Description }}</div>

                    {{ range .Order }}

                        {{ with $req := findRequest $.Requests . }}
                        <div class="request">

                            <h4 id="request-{{ slugify $req.Name }}">
                                {{ $req.Name }}
                                <a href="#request-{{ slugify $req.Name }}"><i class="glyphicon glyphicon-link"></i></a>
                            </h4>

                            <div>{{ markdown $req.Description }}</div>

                            <div>
                                <ul class="nav nav-tabs" role="tablist">
                                    <li role="presentation" class="active"><a href="#request-{{ slugify $req.Name }}-example-curl" data-toggle="tab">Curl</a></li>
                                    <li role="presentation"><a href="#request-{{ slugify $req.Name }}-example-http" data-toggle="tab">HTTP</a></li>
                                </ul>
                                <div class="tab-content">
                                    <div class="tab-pane active" id="request-{{ slugify $req.Name }}-example-curl">
                                        <pre><code class="hljs curl">{{ curlSnippet $req }}</code></pre>
                                    </div>
                                    <div class="tab-pane" id="request-{{ slugify $req.Name }}-example-http">
                                        <pre><code class="hljs http">{{ httpSnippet $req }}</code></pre>
                                    </div>
                                </div>
                            </div>

                            {{ with $req.Responses }}
                            <div>
                                <ul class="nav nav-tabs" role="tablist">
                                    {{ range $index, $res := . }}
                                    <li role="presentation"{{ if eq $index 0 }} class="active"{{ end }}>
                                        <a href="#request-{{ slugify $req.Name }}-responses-{{ $res.ID }}" data-toggle="tab">
                                            {{ if eq (len $req.Responses) 1 }}
                                                Response
                                            {{ else}}
                                                {{ $res.Name }}
                                            {{ end }}
                                        </a>
                                    </li>
                                    {{ end }}
                                </ul>
                                <div class="tab-content">
                                    {{ range $index, $res := . }}
                                    <div class="tab-pane{{ if eq $index 0 }} active{{ end }}" id="request-{{ slugify $req.Name }}-responses-{{ $res.ID }}">
                                        <table class="table table-bordered">
                                            <tr><th style="width: 20%;">Status</th><td>{{ $res.ResponseCode.Code }} {{ $res.ResponseCode.Name }}</td></tr>
                                            {{ range $res.Headers }}
                                            <tr><th style="width: 20%;">{{ .Name }}</th><td>{{ .Value }}</td></tr>
                                            {{ end }}
                                            {{ with $example := $res.Text }}
                                            <tr><td class="response-text-sample" colspan="2">
                                                <pre><code>{{ indentJSON $example }}</code></pre>
                                            </td></tr>
                                            {{ end }}
                                        </table>
                                    </div>
                                    {{ end }}
                                </div>
                            </div>
                            {{ end }}

                            <hr>
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
</body>
</html>
