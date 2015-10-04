package main

import (
	"html/template"
)

var sourceTemplate = template.Must(template.New("source").Parse(`<html>
	<head>
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.8.0/styles/tomorrow.min.css">
		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.8.0/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>
	</head>

	<body>
		<pre><code>{{printf "%s" .}}</code></pre>
	</body>
</html>
`))

