---
title: hugo
description: Provides global access to Hugo-related data.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: 
  signatures: [hugo]
relatedFunctions:
  - hugo
  - page
  - site
aliases: [/functions/hugo]
---

`hugo` returns an instance that contains the following functions:

`hugo.BuildDate`
: (`string`) The compile date of the current Hugo binary formatted per [RFC&nbsp;3339](https://datatracker.ietf.org/doc/html/rfc3339) (e.g., `2023-05-23T08:14:20Z`).

`hugo.CommitHash`
: (`string`) The Git commit hash of the Hugo binary (e.g., `0a95d6704a8ac8d41cc5ca8fffaad8c5c7a3754a`).

`hugo.Deps`
: (`[]*hugo.Dependency`) See [hugo.Deps](#hugodeps).

`hugo.Environment`
: (`string`) The current running environment as defined through the `--environment` CLI flag (e.g., `development`, `production`).

`hugo.Generator`
: (`template.HTML`) Renders an HTML `meta` element identifying the software that generated the site (e.g., `<meta name="generator" content="Hugo 0.112.0">`).

`hugo.GoVersion`
: (`string`) The Go version used to compile the Hugo binary (e.g., `go1.20.4`). {{< new-in "0.101.0" >}}

`hugo.IsExtended`
: (`bool`) Returns `true` if the Hugo binary is the extended version.

`hugo.IsProduction`
: (`bool`) Returns `true` if `hugo.Environment` is set to the production environment.

`hugo.Version`
: (`hugo.VersionString`) The current version of the Hugo binary (e.g., `0.112.1`).

`hugo.WorkingDir`
: (`string`) The project working directory (e.g., `/home/user/projects/my-hugo-site`). {{< new-in "0.112.0" >}}

## hugo.Deps

{{< new-in "0.92.0" >}}

`hugo.Deps` returns a list of dependencies for a project (either Hugo Modules or local theme components).

Each dependency contains:

Owner
: (`*hugo.Dependency`) In the dependency tree, this is the first module that defines this module as a dependency (e.g., `github.com/gohugoio/hugo-mod-bootstrap-scss/v5`).

Path
: (`string`) The module path or the path below your `themes` directory (e.g., `github.com/gohugoio/hugo-mod-jslibs-dist/popperjs/v2`).

Replace
: (`*hugo.Dependency`) Replaced by this dependency.

Time
: (`time.Time`) The time that the version was created (e.g., `2022-02-13 15:11:28 +0000 UTC`).

Vendor
: (`bool`) Returns `true` if the dependency is vendored.

Version
: (`string`) The module version (e.g., `v2.21100.20000`).

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
