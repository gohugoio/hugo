---
title: Get
description: Returns a slice of weighted pages to which the given term has been assigned.
categories: []
keywords: []
action:
  related: []
  returnType: page.WeightedPages
  signatures: [TAXONOMY.Get TERM]
toc: true
---

The `Get` method on a `Taxonomy` object returns a slice of [weighted pages] to which the given [term] has been assigned.

{{% include "methods/taxonomy/_common/get-a-taxonomy-object.md" %}}

## Get the weighted pages

Now that we have captured the "genres" `Taxonomy` object, let's get the weighted pages to which the "suspense" term has been assigned:

```go-html-template
{{ $weightedPages := $taxonomyObject.Get "suspense" }}
```

The above is equivalent to:

```go-html-template
{{ $weightedPages := $taxonomyObject.suspense }}
```

But, if the term is not a valid [identifier], you cannot use the [chaining] syntax. For example, this will throw an error because the identifier contains a hyphen:

```go-html-template
{{ $weightedPages := $taxonomyObject.my-genre }}
```

You could also use the [`index`] function, but the syntax is more verbose:

```go-html-template
{{ $weightedPages := index $taxonomyObject "my-genre" }}
```

To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $weightedPages }}</pre>
```

## Example

With this template:

```go-html-template
{{ $weightedPages := $taxonomyObject.Get "suspense" }}
{{ range $weightedPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

Hugo renders:

```html
<h2><a href="/books/jamaica-inn/">Jamaica inn</a></h2>
<h2><a href="/books/death-on-the-nile/">Death on the nile</a></h2>
<h2><a href="/books/and-then-there-were-none/">And then there were none</a></h2>
```

[chaining]: /getting-started/glossary/#chain
[`index`]: /functions/collections/indexfunction/
[identifier]: /getting-started/glossary/#identifier
[term]: /getting-started/glossary/#term
[weighted pages]: /getting-started/glossary/#weighted-page
