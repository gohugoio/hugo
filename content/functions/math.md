---
title: math
linktitle: Math
description: Hugo provides six mathematical operators in templates.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [math, operators]
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
ns:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

There are 6 basic mathematical operators that can be used in Hugo templates:

| Function | Description              | Example                       |
| -------- | ------------------------ | ----------------------------- |
| `add`    | Adds two integers.       | `{{add 1 2}}` &rarr; 3        |
| `div`    | Divides two integers.    | `{{div 6 3}}` &rarr; 2        |
| `mod`    | Modulus of two integers. | `{{mod 15 3}}` &rarr; 0       |
| `modBool`| Boolean of modulus of two integers. Evaluates to `true` if = 0. | `{{modBool 15 3}}` &rarr; true |
| `mul`    | Multiplies two integers. | `{{mul 2 3}}` &rarr; 6        |
| `sub`    | Subtracts two integers.  | `{{sub 3 2}}` &rarr; 1        |

## Using `add` with Strings

You can also use the `add` function with strings. You may like this functionality in many use cases, including creating new variables by combining page- or site-level variables with other strings.

For example, social media sharing with [Twitter Cards][cards] requires the following `meta` link in your site's `<head>` to display Twitter's ["Summary Card with Large Image"][twtsummary]:

```html
<meta name="twitter:image" content="http://yoursite.com/images/my-twitter-image.jpg">
```

Let's assume you have an `image` field in the front matter of each of your content files:

```yaml
---
title: My Post
image: my-post-image.jpg
---
```

You can then concatenate the `image` value (string) with the path to your `images` directory in `static` and leverage a URL-related templating function for increased flexibility:

{{% code file="partials/head/twitter-card.html" %}}
```html
{{$socialimage := add "images/" .Params.image}}
<meta name="twitter:image" content="{{ $socialimage | absURL }}">
```
{{% /code %}}

{{% note %}}
The `add` example above makes use of the [`absURL` function](/functions/absurl/). `absURL` is a more elegant and future-proofed approach to creating URLs than combining `.Site.BaseURL` with hard-coded strings&mdash;a templating style sometimes seen in older [Hugo themes](/themes). `absURL` works very well, for example, when creating `link` references to stylesheets and other metadata in your rendered site's `<head>`.
{{% /note %}}

[cards]: https://dev.twitter.com/cards/overview
[twtsummary]: https://dev.twitter.com/cards/types/summary-large-image