---
title: hugo
description: The `hugo` function provides easy access to Hugo-related data.
keywords: []
categories: [functions]
menu:
  docs:
    parent: functions
toc:
signature: ["hugo"]
relatedfuncs: []
---

`hugo` returns an instance that contains the following functions:

hugo.Generator
: `<meta>` tag for the version of Hugo that generated the site. `hugo.Generator` outputs a *complete* HTML tag; e.g. `<meta name="generator" content="Hugo 0.63.2">`

hugo.Version
: the current version of the Hugo binary you are using e.g. `0.99.1`

hugo.GoVersion
: returns the version of Go that the Hugo binary was built with. {{< new-in "0.101.0" >}}

hugo.Environment
: the current running environment as defined through the `--environment` cli tag

hugo.CommitHash
: the git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`

hugo.BuildDate
: the compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`

hugo.IsExtended
: whether this is the extended Hugo binary.

hugo.IsProduction
: returns true if `hugo.Environment` is set to the production environment

{{% note "Use the Hugo Generator Tag" %}}
We highly recommend using `hugo.Generator` in your website's `<head>`. `hugo.Generator` is included by default in all themes hosted on [themes.gohugo.io](https://themes.gohugo.io). The generator tag allows the Hugo team to track the usage and popularity of Hugo.
{{% /note %}}

hugo.Deps
: See [hugo.Deps](#hugodeps)

## hugo.Deps

{{< new-in "0.92.0" >}}

`hugo.Deps` returns a list of dependencies for a project (either Hugo Modules or local theme components).

Each dependency contains:

Path (string)
: Returns the path to this module. This will either be the module path, e.g. "github.com/gohugoio/myshortcodes", or the path below your /theme folder, e.g. "mytheme".

Version (string)
:  The module version.

Vendor (bool)
: Whether this dependency is vendored.

Time (time.Time)
: Time version was created.

Owner
: In the dependency tree, this is the first module that defines this module as a dependency.

Replace (*Dependency)
: Replaced by this dependency.

An example table listing the dependencies:

```html
 <h2>Dependencies</h2>
<table class="table table-dark">
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">Owner</th>
      <th scope="col">Path</th>
      <th scope="col">Version</th>
      <th scope="col">Time</th>
      <th scope="col">Vendor</th>
    </tr>
  </thead>
  <tbody>
    {{ range $index, $element := hugo.Deps }}
    <tr>
      <th scope="row">{{ add $index 1 }}</th>
      <td>{{ with $element.Owner }}{{.Path }}{{ end }}</td>
      <td>
        {{ $element.Path }}
        {{ with $element.Replace }}
        => {{ .Path }}
        {{ end }}
      </td>
      <td>{{ $element.Version }}</td>
      <td>{{ with $element.Time }}{{ . }}{{ end }}</td>
      <td>{{ $element.Vendor }}</td>
    </tr>
    {{ end }}
  </tbody>
</table>
```
