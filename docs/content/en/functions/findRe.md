---
title: findRE
description: Returns a list of strings that match the regular expression.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [regex]
signature: ["findRE PATTERN INPUT [LIMIT]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---


By default all matches will be included. The number of matches can be limitted with an optional third parameter.

The example below returns a list of all second level headers (`<h2>`) in the content:

```
{{ findRE "<h2.*?>(.|\n)*?</h2>" .Content }}
```

You can limit the number of matches in the list with a third parameter. The following example shows how to limit the returned value to just one match (or none, if there are no matched substrings):

```
{{ findRE "<h2.*?>(.|\n)*?</h2>" .Content 1 }}
    <!-- returns ["<h2 id="#foo">Foo</h2>"] -->
```

{{% note %}}
Hugo uses Go's [Regular Expression package](https://golang.org/pkg/regexp/), which is the same general syntax used by Perl, Python, and other languages but with a few minor differences for those coming from a background in PCRE. For a full syntax listing, see the [GitHub wiki for re2](https://github.com/google/re2/wiki/Syntax).

If you are just learning RegEx, or at least Go's flavor, you can practice pattern matching in the browser at <https://regex101.com/>.
{{% /note %}}


[partials]: /templates/partials/
[`plainify`]: /functions/plainify/
[toc]: /content-management/toc/
[`urlize`]: /functions/urlize
