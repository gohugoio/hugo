---
title: "Content Templates"
date: "2013-07-01"
---

Content templates are created in a directory matching the name of the content.

Content pages are of the type "page" and have all the [page
variables](/layout/variables/) available to use in the templates.

In the following examples we have created two different content types as well as
a default content type.

    ▾ layouts/
      ▾ post/
          single.html
      ▾ project/
          single.html

Hugo also has support for a default content template to be used in the event
that a specific template has not been provided for that type. The default type
works the same as the other types but the directory must be called "_default".
[Content views](/layout/views) can also be defined in the "_default" directory.


    ▾ layouts/
      ▾ _default/
          single.html




## post/single.html
This content template is used for [spf13.com](http://spf13.com).
It makes use of [chrome templates](/layout/chrome)

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}
    {{ $baseurl := .Site.BaseUrl }}

    <section id="main">
      <h1 id="title">{{ .Title }}</h1>
      <div>
            <article id="content">
               {{ .Content }}
            </article>
      </div>
    </section>

    <aside id="meta">
        <div>
        <section>
          <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
          <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
        </section>
        <ul id="categories">
          {{ range .Params.topics }}
            <li><a href="{{ $baseurl }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        <ul id="tags">
          {{ range .Params.tags }}
            <li> <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        </div>
        <div>
            {{ if .Prev }}
              <a class="previous" href="{{.Prev.Permalink}}"> {{.Prev.Title}}</a>
            {{ end }}
            {{ if .Next }}
              <a class="next" href="{{.Next.Permalink}}"> {{.Next.Title}}</a>
            {{ end }}
        </div>
    </aside>

    {{ template "chrome/disqus.html" . }}
    {{ template "chrome/footer.html" . }}


## project/single.html
This content template is used for [spf13.com](http://spf13.com).
It makes use of [chrome templates](/layout/chrome)


    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}
    {{ $baseurl := .Site.BaseUrl }}

    <section id="main">
      <h1 id="title">{{ .Title }}</h1>
      <div>
            <article id="content">
               {{ .Content }}
            </article>
      </div>
    </section>

    <aside id="meta">
        <div>
        <section>
          <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
          <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
        </section>
        <ul id="categories">
          {{ range .Params.topics }}
          <li><a href="{{ $baseurl }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        <ul id="tags">
          {{ range .Params.tags }}
            <li> <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        </div>
    </aside>

    {{if isset .Params "project_url" }}
    <div id="ribbon">
        <a href="{{ index .Params "project_url" }}" rel="me">Fork me on GitHub</a>
    </div>
    {{ end }}

    {{ template "chrome/footer.html" }}


Notice how the project/single.html template uses an additional parameter unique
to this template. This doesn't need to be defined ahead of time. If the key is
present in the front matter than it can be used in the template.
