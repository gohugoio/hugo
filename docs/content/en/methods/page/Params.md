---
title: Params
description: Returns a map of custom parameters as defined in the front matter of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: maps.Params
    signatures: [PAGE.Params]
---

By way of example, consider this front matter:

{{< code-toggle file=content/annual-conference.md fm=true >}}
title = 'Annual conference'
date = 2023-10-17T15:11:37-07:00
[params]
display_related = true
key-with-hyphens = 'must use index function'
[params.author]
  email = 'jsmith@example.org'
  name = 'John Smith'
{{< /code-toggle >}}

The `title` and `date` fields are standard [front matter fields], while the other fields are user-defined.

Access the custom fields by [chaining](g) the [identifiers](g) when needed:

```go-html-template
{{ .Params.display_related }} → true
{{ .Params.author.email }} → jsmith@example.org
{{ .Params.author.name }} → John Smith
```

In the template example above, each of the keys is a valid identifier. For example, none of the keys contains a hyphen. To access a key that is not a valid identifier, use the [`index`] function:

```go-html-template
{{ index .Params "key-with-hyphens" }} → must use index function
```

[`index`]: /functions/collections/indexfunction/
[front matter fields]: /content-management/front-matter/#fields
