---
title: base64
description: "`base64Encode` and `base64Decode` let you easily decode content with a base64 encoding and vice versa through pipes."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
relatedfuncs: []
signature: ["base64Decode INPUT", "base64Encode INPUT"]
workson: []
hugoversion:
deprecated: false
draft: false
aliases: []
---

An example:

{{< code file="base64-input.html" >}}
<p>Hello world = {{ "Hello world" | base64Encode }}</p>
<p>SGVsbG8gd29ybGQ = {{ "SGVsbG8gd29ybGQ=" | base64Decode }}</p>
{{< /code >}}

{{< output file="base-64-output.html" >}}
<p>Hello world = SGVsbG8gd29ybGQ=</p>
<p>SGVsbG8gd29ybGQ = Hello world</p>
{{< /output >}}

You can also pass other data types as arguments to the template function which tries to convert them. The following will convert *42* from an integer to a string because both `base64Encode` and `base64Decode` always return a string.

```
{{ 42 | base64Encode | base64Decode }}
=> "42" rather than 42
```

## `base64` with APIs

Using base64 to decode and encode becomes really powerful if we have to handle
responses from APIs.

```
{{ $resp := getJSON "https://api.github.com/repos/gohugoio/hugo/readme"  }}
{{ $resp.content | base64Decode | markdownify }}
```

The response of the GitHub API contains the base64-encoded version of the [README.md](https://github.com/gohugoio/hugo/blob/master/README.md) in the Hugo repository. Now we can decode it and parse the Markdown. The final output will look similar to the rendered version on GitHub.
