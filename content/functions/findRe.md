---
title: findre
linktitle: findRE
description: Returns a list of strings that match the regular expression.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [regex]
ns:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

<!-- see https://github.com/spf13/hugo/issues/3188 -->

Returns a list of strings that match the regular expression. By default all matches will be included. The number of matches can be limitted with an optional third parameter.

The example below returns a list of all second level headers (`<h2>`) in the content:

```
{{ findRE "<h2.*?>(.|\n)*?</h2>" .Content }}
```

You can limit the number of matches in the list with a third parameter. The following example shows how to limit the returned value to just one match (or none, if there are no matched substrings):

```golang
{{ findRE "<h2.*?>(.|\n)*?</h2>" .Content 1 }}
    <!-- returns ["<h2 id="#foo">Foo</h2>"] -->
```

## `findRE` Example: Building a Table of Contents

`findRE` allows us to build an automatically generated table of contents that could be used for a simple scrollspy if you don't want to use [Hugo's native .TableOfContents feature][toc]. The following shows how this could be done in a [partial template][partials]:

{{% code file="layouts/partials/toc.html" download="toc.html" %}}
```html
{{ $headers := findRE "<h2.*?>(.|\n)*?</h2>" .Content }}
{{ if ge (len $headers) 1 }}
    <ul>
    {{ range $headers }}
        <li>
            <a href="#{{ . | plainify | urlize }}">
                {{ . | plainify }}
            </a>
        </li>
    {{ end }}
    </ul>
{{ end }}
```
{{% /code %}}

The preceding snippet tries to find all second-level headers and generate a list where at least one header is found. [`plainify`][] strips the HTML and [`urlize`][] converts the header into a valid URL.

[partials]: /templates/partials/
[`plainify`]: /functions/plainify/
[toc]: /content-management/toc/
[`urlize`]: /functions/urlize