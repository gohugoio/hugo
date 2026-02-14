---
title: Partial decorators
description: Use partial decorators to create reusable wrapper components that enclose and compose template content.
categories: []
keywords: [decorator]
weight: 170
---

{{< new-in "0.154.0" />}}

## Overview

{{% glossary-term "partial decorator" %}}

This approach creates a connection between two files. The calling template provides a block of code and the partial decorator determines where that code appears. This allows the partial to wrap around content without needing to know the specific markup or internal logic of the enclosed block.

## Implementation

To use a partial decorator, use a block-style call in your templates. The [`with`][] statement is required to initiate the partial and create a container for the content. This block can include any valid template code including page methods and functions.

```go-html-template {file="layouts/home.html" copy=true}
{{ with partial "components/wrapper.html" . }}
  <p>Everything in this block will be wrapped.</p>
  <p>{{ .Content | transform.Plainify | strings.Truncate 200 }}</p>
{{ end }}
```

Inside the partial template, place the `templates.Inner` function call where the wrapped content should appear.

```go-html-template {file="layouts/_partials/components/wrapper.html" copy=true}
<div class="wrapper-styling">
  {{ templates.Inner . }}
</div>
```

The `with` statement creates a new [scope](g). Variables defined outside of the `with` block are not available inside it. To use external data within the wrapped content, you must ensure it is part of the [context](g) passed in the partial call.

A key feature of the `templates.Inner` function is its ability to accept a context argument. By passing a context to the function, you define what the dot (`.`) represents inside the wrapped block. This ensures that the injected content has access to the correct data even when nested inside multiple layers of wrappers.

## Benefits of composition

Using partial decorators to build wrapper components provides several advantages:

- It eliminates the need to use separate partials for opening and closing tags when encapsulating a block of code.
- It prevents parameter bloat because a standard partial no longer requires an extensive list of arguments to account for every possible variation of the content inside it.
- It enables clean composition where the wrapped block can execute any template logic without the wrapper needing to receive or process that data.

This approach separates container logic from content logic. The wrapper handles structural requirements like specific class hierarchies or CSS grid containers. The calling template retains control over the inner markup and how data is displayed.

## Example

The following templates illustrate how to nest three wrapper components including a section, a column, and a card while passing context through each layer.

The home template initiates the structure by calling the section, column, and card partials as decorators:

```go-html-template {file="layouts/home.html" copy=true}
{{ $ctx := dict
  "page" .
  "label" "Recent Posts"
  "pageCollection" ((site.GetPage "/posts").RegularPages)
}}

{{ with partial "components/section.html" $ctx }}
  <div class="grid-wrapper">
    {{ range .pageCollection }}
      {{ with partial "components/column.html" (dict "page" . "class" "col-half") }}
        {{ with partial "components/card.html" (dict "page" .page "url" .page.RelPermalink "title" .page.LinkTitle) }}
          <p>
            {{ .page.Content | plainify | strings.Truncate 240 }}
          </p>
        {{ end }}
      {{ end }}
    {{ end }}
  </div>
{{ end }}
```

The section component provides a semantic container and an optional heading:

```go-html-template {file="layouts/_partials/components/section.html" copy=true}
<section class="content-section">
  {{ with .label }}
    <h2 class="section-label">{{ . }}</h2>
  {{ end }}
  <div class="section-content">
    {{ templates.Inner . }}
  </div>
</section>
```

The column component manages layout width by applying a CSS class:

```go-html-template {file="layouts/_partials/components/column.html" copy=true}
<div class="{{ .class | default `column-default` }}">
  {{ templates.Inner . }}
</div>
```

The card component defines the visual boundary for the content:

```go-html-template {file="layouts/_partials/components/card.html" copy=true}
<div class="card">
  {{ with .title }}
    <h2 class="card-title">
      {{ if $.url }}
        <a href="{{ $.url }}">{{ . }}</a>
      {{ else }}
        {{ . }}
      {{ end }}
    </h2>
  {{ end }}

  <div class="card-body">
    {{ templates.Inner . }}
  </div>

  {{ with .url }}
    <div class="card-footer">
      <a href="{{ . }}">Read more</a>
    </div>
  {{ end }}
</div>
```

[`with`]: /functions/go-template/with/
