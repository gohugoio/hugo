---
title: Configure taxonomies
linkTitle: Taxonomies
description: Configure taxonomies.
categories: []
keywords: []
---

The default configuration defines two [taxonomies](g), `categories` and `tags`.

{{< code-toggle config=taxonomies />}}

When creating a taxonomy:

- Use the singular form for the key (e.g., `category`).
- Use the plural form for the value (e.g., `categories`).

Then use the value as the key in front matter:

{{< code-toggle file=content/example.md fm=true >}}
---
title: Example
categories:
  - vegetarian
  - gluten-free
tags:
  - appetizer
  - main course
{{< /code-toggle >}}

If you do not expect to assign more than one [term](g) from a given taxonomy to a content page, you may use the singular form for both key and value:

{{< code-toggle file=hugo >}}
taxonomies:
  author: author
{{< /code-toggle >}}

Then in front matter:

{{< code-toggle file=content/example.md fm=true >}}
---
title: Example
author:
  - Robert Smith
{{< /code-toggle >}}

The example above illustrates that even with a single term, the value is still provided as an array.

You must explicitly define the default taxonomies to maintain them when adding a new one:

{{< code-toggle file=hugo >}}
taxonomies:
  author: author
  category: categories
  tag: tags
{{< /code-toggle >}}

To disable the taxonomy system, use the [`disableKinds`] setting in the root of your site configuration to disable the `taxonomy` and `term` page [kinds](g).

{{< code-toggle file=hugo >}}
disableKinds = ['taxonomy','term']
{{< /code-toggle >}}

[`disableKinds`]: /configuration/all/#disablekinds

See the [taxonomies] section for more information.

[taxonomies]: /content-management/taxonomies/
