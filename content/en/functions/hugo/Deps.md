---
title: hugo.Deps
description: Returns a slice of project dependencies, either Hugo Modules or local theme components.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: '[]hugo.Dependency'
  signatures: [hugo.Deps]
---

The `hugo.Deps` function returns a slice of project dependencies, either Hugo Modules or local theme components. Each dependency contains:

Owner
: (`hugo.Dependency`) In the dependency tree, this is the first module that defines this module as a dependency (e.g., `github.com/gohugoio/hugo-mod-bootstrap-scss/v5`).

Path
: (`string`) The module path or the path below your `themes` directory (e.g., `github.com/gohugoio/hugo-mod-jslibs-dist/popperjs/v2`).

Replace
: (`hugo.Dependency`) Replaced by this dependency.

Time
: (`time.Time`) The time that the version was created (e.g., `2022-02-13 15:11:28 +0000 UTC`).

Vendor
: (`bool`) Reports whether the dependency is vendored.

Version
: (`string`) The module version (e.g., `v2.21100.20000`).

An example table listing the dependencies:

```go-html-template
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
      <td>{{ with $element.Owner }}{{ .Path }}{{ end }}</td>
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
