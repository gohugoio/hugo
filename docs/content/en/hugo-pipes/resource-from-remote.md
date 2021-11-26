---
title: Resource from a remote address
linkTitle: Resource from Remote
description: Hugo Pipes allows the creation of a resource from a remote address.
date: 2021-11-26
publishdate: 2021-11-26
lastmod: 2021-11-26
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 90
weight: 90
sections_weight: 90
draft: false
---

It is possible to create a resource directly from a remote server using `resources.FromRemote`.

This example creates a local resource from a file downloaded from a remote server.

```go-html-template
{{ $script := resources.FromRemote "https://example.com/script.js" }}
<script type="text/javascript" src="{{ $script.Permalink }}"></script>
```

You can process an asset with Hugo Pipes, such as for the minification of a JavaScript file downloaded from a remote server.
```go-html-template
{{ $script := resources.FromRemote "https://example.com/script.js" | resources.Minify }}
<script type="text/javascript" src="{{ $script.Permalink }}"></script>
```

The URL accept [variadic arguments][variadic]:

```
{{ $resource := resources.FromRemote "url" "arg1" "arg2" "arg n" }}
```

All passed arguments will be joined to the final URL:

```
{{ $urlPre := "https://api.github.com" }}
{{ $resource := resources.FromRemote $urlPre "/users/GITHUB_USERNAME/gists" }}
```

This will resolve internally to the following:

```
{{ $resource := resources.FromRemote "https://api.github.com/users/GITHUB_USERNAME/gists" }}
```

### Add HTTP headers

The `resources.FromRemote` function takes an optional map as the last argument, e.g.:

```
{{ $resource := resources.FromRemote "https://example.org/api" (dict "Authorization" "Bearer abcd")  }}
```

If you need multiple values for the same header key, use a slice:

```
{{ $resource := resources.FromRemote "https://example.org/api" (dict "X-List" (slice "a" "b" "c"))  }}
```

[variadic]: https://en.wikipedia.org/wiki/Variadic_function