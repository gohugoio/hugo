---
title: urls.URLize
description: Takes a string, sanitizes it for usage in URLs, and converts spaces to hyphens.
categories: []
keywords: []
action:
  aliases: [urlize]
  related:
    - functions/urls/Anchorize
  returnType: string
  signatures: [urls.URLize INPUT]
aliases: [/functions/urlize]
---

The following examples pull from a content file with the following front matter:

{{< code-toggle file="content/blog/greatest-city.md" fm=true >}}
title = "The World's Greatest City"
location = "Chicago IL"
tags = ["pizza","beer","hot dogs"]
{{< /code-toggle >}}

The following might be used as a partial within a [single page template][singletemplate]:

{{< code file="layouts/partials/content-header.html" >}}
<header>
  <h1>{{ .Title }}</h1>
  {{ with .Params.location }}
    <div><a href="/locations/{{ . | urlize }}">{{ . }}</a></div>
  {{ end }}
  <!-- Creates a list of tags for the content and links to each of their pages -->
  {{ with .Params.tags }}
    <ul>
      {{ range .}}
        <li>
          <a href="/tags/{{ . | urlize }}">{{ . }}</a>
        </li>
      {{ end }}
    </ul>
  {{ end }}
</header>
{{< /code >}}

The preceding partial would then output to the rendered page as follows:

```html
<header>
  <h1>The World&#39;s Greatest City</h1>
  <div><a href="/locations/chicago-il">Chicago IL</a></div>
  <ul>
    <li>
      <a href="/tags/pizza">pizza</a>
    </li>
    <li>
      <a href="/tags/beer">beer</a>
    </li>
    <li>
      <a href="/tags/hot-dogs">hot dogs</a>
    </li>
  </ul>
</header>
```

[singletemplate]: /templates/single-page-templates/
