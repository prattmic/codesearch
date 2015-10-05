package main

import (
	"html/template"
)

var sourceTemplate = template.Must(template.New("source").Parse(`<html>
	<head>
		<script src="//cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js"></script>
	</head>

	<body>
		<pre class="prettyprint">{{printf "%s" .}}</pre>
	</body>
</html>
`))
