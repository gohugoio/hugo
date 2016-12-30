---
lastmod: 2015-01-27
date: 2013-07-09
menu:
  main:
    parent: extras
next: /extras/localfiles
prev: /extras/highlighting
title: Table of Contents
---

Hugo will automatically parse the Markdown for your content and create
a Table of Contents you can use to guide readers to the sections within
your content.

## Usage

Simply create content like you normally would with the appropriate
headers.

Hugo will take this Markdown and create a table of contents stored in the
[content variable](/layout/variables/) `.TableOfContents`.


### Template Example

This is example code of a [single.html template](/layout/content/).

    {{ partial "header.html" . }}
        <div id="toc" class="well col-md-4 col-sm-6">
        {{ .TableOfContents }}
        </div>
        <h1>{{ .Title }}</h1>
        {{ .Content }}
    {{ partial "footer.html" . }}


## Styling your own Table Of Contents

If the automatically generated `.TableOfContents` variable doesn't
meet your needs, you may use the `.TocEntries` array to generate your
own.

### Template Example

This is an example partial for rendering the top two levels of the
`.TocEntries`:

    <div class="toc">
	  <ul>
	    {{ $currentUrl := .URL }}
		{{ range .TocEntries }}
		<li><a href="{{$currentUrl}}#{{ .Id }}">{{ .Text }}</a>
		  {{ if .Contents }}
		  <ul>
		    {{ range .Contents }}
			<li><a href="{{$currentUrl}}#{{ .Id }}">{{ .Text }}</a>
			{{ end }}
		  </ul>
		  {{ end }}
		{{ end}}
      </ul>
	</div>
