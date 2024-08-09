---
title: Single templates
description: Create a single template to render a single page.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 70
weight: 70
toc: true
aliases: [/layout/content/,/templates/single-page-templates/]
---

The single template below inherits the site's shell from the [base template].

[base template]: /templates/types/

{{< code file=layouts/_default/single.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
{{< /code >}}

Review the [template lookup order] to select a template path that provides the desired level of specificity.

[template lookup order]: /templates/lookup-order/#single-templates

The single template below inherits the site's shell from the base template, and renders the page title, creation date, content, and a list of associated terms in the "tags" taxonomy.

{{< code file=layouts/_default/single.html >}}
{{ define "main" }}
  <section>
    <h1>{{ .Title }}</h1>
    {{ with .Date }}
      {{ $dateMachine := . | time.Format "2006-01-02T15:04:05-07:00" }}
      {{ $dateHuman := . | time.Format ":date_long" }}
      <time datetime="{{ $dateMachine }}">{{ $dateHuman }}</time>
    {{ end }}
    <article>
      {{ .Content }}
    </article>
    <aside>
      {{ with .GetTerms "tags" }}
        <div>{{ (index . 0).Parent.LinkTitle }}</div>
        <ul>
          {{ range . }}
            <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
          {{ end }}
        </ul>
      {{ end }}
    </aside>
  </section>
{{ end }}
{{< /code >}}
