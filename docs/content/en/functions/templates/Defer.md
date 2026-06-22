---
title: templates.Defer
description: Defer execution of a template until after all sites and output formats have been rendered.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [templates.Defer OPTIONS]
aliases: [/functions/templates.defer]
---

> [!NOTE]
> Do not call this function within a `partialCached` template. This restriction applies transitively: if `partialCached` calls a partial that calls `templates.Defer`, Hugo returns an error. Using this function within shortcode or render hook templates may also lead to unpredictable results.

In some rare use cases, you may need to defer the execution of a template until after all sites and output formats have been rendered. One such example could be [css.TailwindCSS][] using the output of [`hugo_stats.json`][] to determine which classes and other HTML identifiers are being used in the final output:

```go-html-template {file="layouts/baseof.html" copy=true}
<head>
  ...
  {{ with (templates.Defer (dict "key" "global")) }}
    {{ partial "css.html" . }}
  {{ end }}
  ...
</head>
```

```go-html-template {file="layouts/_partials/css.html" copy=true}
{{ with resources.Get "css/main.css" }}
  {{ $opts := dict "minify" (not hugo.IsDevelopment) }}
  {{ with . | css.TailwindCSS $opts }}
    {{ if hugo.IsDevelopment }}
      <link rel="stylesheet" href="{{ .RelPermalink }}">
    {{ else }}
      {{ with . | fingerprint }}
        <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
      {{ end }}
    {{ end }}
  {{ end }}
{{ end }}
```

> [!NOTE]
> This function only works in combination with the `with` keyword.
>
> Variables defined on the outside are not visible on the inside and vice versa. To pass in data, use the `data` [option](#options).

For the above to work well when running the server (or `hugo -w`), you want to have a configuration similar to this:

{{< code-toggle file=hugo >}}
[build]
  [build.buildStats]
    enable = true
  [[build.cachebusters]]
    source = 'assets/notwatching/hugo_stats\.json'
    target = 'css'
  [[build.cachebusters]]
    source = '(postcss|tailwind)\.config\.js'
    target = 'css'
[module]
  [[module.mounts]]
    source = 'assets'
    target = 'assets'
  [[module.mounts]]
    disableWatch = true
    source = 'hugo_stats.json'
    target = 'assets/notwatching/hugo_stats.json'
{{< /code-toggle >}}

## Options

The `templates.Defer` function requires a single argument, a map with the following optional keys:

`key`
: (`string`) The key to use for the deferred template. This will, combined with a hash of the template content, be used as a cache key. If this is not set, Hugo will execute the deferred template on every render. This is not what you want for shared resources like CSS and JavaScript.

`data`
: (`map`) Optional map to pass as data to the deferred template. This will be available in the deferred template as `.` or `$`.

```go-html-template
Language Outside: {{ site.Language.Name }}
Page Outside: {{ .RelPermalink }}
I18n Outside: {{ i18n "hello" }}
{{ $data := (dict "page" . )}}
{{ with (templates.Defer (dict "data" $data )) }}
     Language Inside: {{ site.Language.Name }}
     Page Inside: {{ .page.RelPermalink }}
     I18n Inside: {{ i18n "hello" }}
{{ end }}
```

The [output format][], [site][], and [language][] will be the same, even if the execution is deferred. In the example above, this means that the `site.Language.Name` and `.RelPermalink` will be the same on the inside and the outside of the deferred template.

[`hugo_stats.json`]: /configuration/build/
[css.TailwindCSS]: /functions/css/tailwindcss/
[language]: /methods/site/language/
[output format]: /configuration/output-formats/
[site]: /methods/page/site/
