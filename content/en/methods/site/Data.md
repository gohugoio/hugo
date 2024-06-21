---
title: Data
description: Returns a data structure composed from the files in the data directory.
categories: []
keywords: []
action:
  related:
    - functions/collections/IndexFunction
    - functions/transform/Unmarshal
    - functions/collections/Where
    - functions/collections/Sort
  returnType: map
  signatures: [SITE.Data]
---

Use the `Data` method on a `Site` object to access data within the data directory, or within any directory [mounted] to the data directory. Supported data formats include JSON, TOML, YAML, and XML.

[mounted]: /hugo-modules/configuration/#module-configuration-mounts

{{% note %}}
Although Hugo can unmarshal CSV files with the [`transform.Unmarshal`] function, do not place CSV files in the data directory. You cannot access data within CSV files using this method.

[`transform.Unmarshal`]: /functions/transform/unmarshal/
{{% /note %}}

Consider this data directory:

```text
data/
├── books/
│   ├── fiction.yaml
│   └── nonfiction.yaml
├── films.json
├── paintings.xml
└── sculptures.toml
```

And these data files:

{{< code file=data/books/fiction.yaml lang=yaml >}}
- title: The Hunchback of Notre Dame
  author: Victor Hugo
  isbn: 978-0140443530
- title: Les Misérables
  author: Victor Hugo
  isbn: 978-0451419439
{{< /code >}}

{{< code file=data/books/nonfiction.yaml lang=yaml >}}
- title: The Ancien Régime and the Revolution
  author: Alexis de Tocqueville
  isbn: 978-0141441641
- title: Interpreting the French Revolution
  author: François Furet
  isbn: 978-0521280495
{{< /code >}}

Access the data by [chaining] the [identifiers]:

```go-html-template
{{ range $category, $books := .Site.Data.books }}
  <p>{{ $category | title }}</p>
  <ul>
    {{ range $books }}
      <li>{{ .title }} ({{ .isbn }})</li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo renders this to:

```html
<p>Fiction</p>
<ul>
  <li>The Hunchback of Notre Dame (978-0140443530)</li>
  <li>Les Misérables (978-0451419439)</li>
</ul>
<p>Nonfiction</p>
<ul>
  <li>The Ancien Régime and the Revolution (978-0141441641)</li>
  <li>Interpreting the French Revolution (978-0521280495)</li>
</ul>
```

To limit the listing to fiction, and sort by title:

```go-html-template
<ul>
  {{ range sort .Site.Data.books.fiction "title" }}
    <li>{{ .title }} ({{ .author }})</li>
  {{ end }}
</ul>
```

To find a fiction book by ISBN:

```go-html-template
{{ range where .Site.Data.books.fiction "isbn" "978-0140443530" }}
  <li>{{ .title }} ({{ .author }})</li>
{{ end }}
```

In the template examples above, each of the keys is a valid [identifier]. For example, none of the keys contains a hyphen. To access a key that is not a valid identifier, use the [`index`] function. For example:

[identifier]: /getting-started/glossary/#identifier

```go-html-template
{{ index .Site.Data.books "historical-fiction" }}
```

[`index`]: /functions/collections/indexfunction/
[chaining]: /getting-started/glossary/#chain
[identifiers]: /getting-started/glossary/#identifier
