// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"bytes"
	"io"

	"github.com/gohugoio/hugo/transform"
	"github.com/gohugoio/hugo/transform/livereloadinject"
)

var buildErrorTemplate = `<!doctype html>
<html class="no-js" lang="">
	<head>
		<meta charset="utf-8">
		<title>Hugo Server: Error</title>
		<style type="text/css">
		body {
			font-family: "Muli",avenir, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
			font-size: 16px;
			background-color: #2f1e2e;
		}
		main {
			margin: auto;
			width: 95%;
			padding: 1rem;
		}		
		.version {
			color: #ccc;
			padding: 1rem 0;
		}
		.stack {
			margin-top: 4rem;
		}
		pre {
			white-space: pre-wrap;      
			white-space: -moz-pre-wrap;  
			white-space: -pre-wrap;     
			white-space: -o-pre-wrap;    
			word-wrap: break-word;     
		}
		.highlight {
			overflow-x: auto;
			margin-bottom: 1rem;
		}
		a {
			color: #0594cb;
			text-decoration: none;
		}
		a:hover {
			color: #ccc;
		}
		</style>
	</head>
	<body>
		<main>
			{{ highlight .Error "apl" "linenos=false,noclasses=true,style=paraiso-dark" }}
			{{ with .File }}
			{{ $params := printf "noclasses=true,style=paraiso-dark,linenos=table,hl_lines=%d,linenostart=%d" (add .LinesPos 1) (sub .Position.LineNumber .LinesPos) }}
			{{ $lexer := .ChromaLexer | default "go-html-template" }}
			{{  highlight (delimit .Lines "\n") $lexer $params }}
			{{ end }}
			{{ with .StackTrace }}
			{{ highlight . "apl" "noclasses=true,style=paraiso-dark" }}
			{{ end }}
			<p class="version">{{ .Version }}</p>
			<a href="">Reload Page</a>
		</main>
</body>
</html>
`

func injectLiveReloadScript(src io.Reader, port int) string {
	var b bytes.Buffer
	chain := transform.Chain{livereloadinject.New(port)}
	chain.Apply(&b, src)

	return b.String()
}
