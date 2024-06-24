---
title: Data sources
description: Use local and remote data sources to augment or create content.
categories: [content management]
keywords: [data,json,toml,yaml,xml]
menu:
  docs:
    parent: content-management
    weight: 280
weight: 280
toc: true
aliases: [/extras/datafiles/,/extras/datadrivencontent/,/doc/datafiles/,/templates/data-templates/]
---

Hugo can access and [unmarshal] local and remote data sources including CSV, JSON, TOML, YAML, and XML. Use this data to augment existing content or to create new content.

[unmarshal]: /getting-started/glossary/#unmarshal

A data source might be a file in the data directory, a [global resource], a [page resource], or a [remote resource].

[global resource]: /getting-started/glossary/#global-resource
[page resource]: /getting-started/glossary/#page-resource
[remote resource]: /getting-started/glossary/#remote-resource

## Data directory

The data directory in the root of your project may contain one or more data files, in either a flat or nested tree. Hugo merges the data files to create a single data structure, accessible with the `Data` method on a `Site` object.

Hugo also merges data directories from themes and modules into this single data structure, where the data directory in the root of your project takes precedence.

{{% note %}}
Hugo reads the combined data structure into memory and keeps it there for the entire build. For data that is infrequently accessed, use global or page resources instead.
{{% /note %}}

Theme and module authors may wish to namespace their data files to prevent collisions. For example:

```text
project/
└── data/
    └── mytheme/
        └── foo.json
```

{{% note %}}
Do not place CSV files in the data directory. Access CSV files as page, global, or remote resources.
{{% /note %}}

See the documentation for the [`Data`] method on a `Site` object for details and examples.

[`Data`]: /methods/site/data/

## Global resources

Use the `resources.Get` and `transform.Unmarshal` functions to access data files that exist as global resources.

See the [`transform.Unmarshal`](/functions/transform/unmarshal/#global-resource) documentation for details and examples.

## Page resources

Use the `Resources.Get` method on a `Page` object combined with the `transform.Unmarshal` function to access data files that exist as page resources.

See the [`transform.Unmarshal`](/functions/transform/unmarshal/#page-resource) documentation for details and examples.

## Remote resources

Use the `resources.GetRemote` and `transform.Unmarshal` functions to access remote data.

See the [`transform.Unmarshal`](/functions/transform/unmarshal/#remote-resource) documentation for details and examples.

## Augment existing content

Use data sources to augment existing content. For example, create a shortcode to render an HTML table from a global CSV resource.

{{< code file=assets/pets.csv >}}
"name","type","breed","age"
"Spot","dog","Collie","3"
"Felix","cat","Malicious","7"
{{< /code >}}

{{< code file=content/example.md lang=text >}}
{{</* csv-to-table "pets.csv" */>}}
{{< /code >}}

{{< code file=layouts/shortcodes/csv-to-table.html >}}
{{ with $file := .Get 0 }}
  {{ with resources.Get $file }}
    {{ with . | transform.Unmarshal }}
      <table>
        <thead>
          <tr>
            {{ range index . 0 }}
              <th>{{ . }}</th>
            {{ end }}
          </tr>
        </thead>
        <tbody>
          {{ range after 1 . }}
            <tr>
              {{ range . }}
                <td>{{ . }}</td>
              {{ end }}
            </tr>
          {{ end }}
        </tbody>
      </table>
    {{ end }}
  {{ else }}
    {{ errorf "The %q shortcode was unable to find %s. See %s" $.Name $file $.Position }}
  {{ end }}
{{ else }}
  {{ errorf "The %q shortcode requires one positional argument, the path to the CSV file relative to the assets directory. See %s" .Name .Position }}
{{ end }}
{{< /code >}}

Hugo renders this to:

name|type|breed|age
:--|:--|:--|:--
Spot|dog|Collie|3
Felix|cat|Malicious|7

## Create new content

Use [content adapters] to create new content.

[content adapters]: /content-management/content-adapters/
