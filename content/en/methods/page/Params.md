---
title: Params
description: Returns a map of custom parameters as defined in the front matter of the given page.
categories: []
keywords: []
action:
  related:
    - functions/collections/IndexFunction
    - methods/site/Params
    - methods/page/Param
  returnType: maps.Params
  signatures: [PAGE.Params]
---

With this front matter:

{{< code-toggle file=content/news/annual-conference.md >}}
title = 'Annual conference'
date = 2023-10-17T15:11:37-07:00
[params]
display_related = true
[params.author]
  email = 'jsmith@example.org'
  name = 'John Smith'
{{< /code-toggle >}}

The `title` and `date` fields are standard parameters---the other fields are user-defined.

Access the custom parameters by [chaining] the [identifiers]:

```go-html-template
{{ .Params.display_related }} → true
{{ .Params.author.name }} → John Smith
```

In the template example above, each of the keys is a valid identifier. For example, none of the keys contains a hyphen. To access a key that is not a valid identifier, use the [`index`] function:

```go-html-template
{{ index .Params "key-with-hyphens" }} → 2023
```

[`index`]: /functions/collections/indexfunction/
[chaining]: /getting-started/glossary/#chain
[identifiers]: /getting-started/glossary/#identifier
