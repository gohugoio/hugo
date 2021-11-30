---
title: Hugo Pipes Introduction
description: Hugo Pipes is Hugo's asset processing set of functions.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 20
weight: 01
sections_weight: 01
draft: false
aliases: [/assets/]
---

### Asset directory

Asset files must be stored in the asset directory. This is `/assets` by default, but can be configured via the configuration file's `assetDir` key.

### From file or URL to resource

In order to process an asset with Hugo Pipes, it must be retrieved as a resource using `resources.Get`. The first argument can be the filepath of the file relative to the asset directory or the URL of the file.

```go-html-template
{{ $style := resources.Get "sass/main.scss" }}
{{ $remoteStyle := resources.Get "https://www.example.com/styles.scss" }}
```

#### Request options

When using an URL, the `resources.Get` function takes an optional options map as the last argument, e.g.:

```
{{ $resource := resources.Get "https://example.org/api" (dict "headers" (dict "Authorization" "Bearer abcd"))  }}
```

If you need multiple values for the same header key, use a slice:

```
{{ $resource := resources.Get "https://example.org/api"  (dict "headers" (dict "X-List" (slice "a" "b" "c")))  }}
```

You can also change the request method and set the request body:

```
{{ $postResponse := resources.Get "https://example.org/api"  (dict 
    "method" "post"
    "body" `{"complete": true}` 
    "headers" (dict 
        "Content-Type" "application/json"
    )
)}}
```

### Asset publishing

Assets will only be published (to `/public`) if `.Permalink` or `.RelPermalink` is used.

### Go Pipes

For improved readability, the Hugo Pipes examples of this documentation will be written using [Go Pipes](/templates/introduction/#pipes):
```go-html-template
{{ $style := resources.Get "sass/main.scss" | resources.ToCSS | resources.Minify | resources.Fingerprint }}
<link rel="stylesheet" href="{{ $style.Permalink }}">
```

### Method aliases

Each Hugo Pipes `resources` transformation method uses a __camelCased__ alias (`toCSS` for `resources.ToCSS`).
Non-transformation methods deprived of such aliases are `resources.Get`, `resources.FromString`, `resources.ExecuteAsTemplate` and `resources.Concat`.

The example above can therefore also be written as follows:
```go-html-template
{{ $style := resources.Get "sass/main.scss" | toCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $style.Permalink }}">
```

### Caching

Hugo Pipes invocations are cached based on the entire _pipe chain_.

An example of a pipe chain is:

```go-html-template
{{ $mainJs := resources.Get "js/main.js" | js.Build "main.js" | minify | fingerprint }}
```

The pipe chain is only invoked the first time it is encountered in a site build, and results are otherwise loaded from cache. As such, Hugo Pipes can be used in templates which are executed thousands or millions of times without negatively impacting the build performance.
