---
title: Aliases
description: Returns the aliases defined in front matter as server-relative URLs, resolved according to the current content dimension.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '[]string'
    signatures: [PAGE.Aliases]
---

The `Aliases` method on a `Page` object returns the values defined in the [`aliases`][] front matter field as server-relative URLs, resolved according to the current [content dimension](g).

The `Aliases` method is useful for generating a `_redirects` file, which contains a source URL, a target URL, and an HTTP status code for each alias. You can use a `_redirects` file with hosting services such as Cloudflare, GitLab Pages, and Netlify.

## Redirects

By default, Hugo handles aliases by creating individual HTML files for each alias path. These files contain a `meta http-equiv="refresh"` tag to redirect the visitor via the browser.

While functional, generating a single `_redirects` file allows your hosting provider to handle redirects at the server level. This is more efficient than client-side redirection and improves performance by eliminating the need to load a middle-man HTML page.

> [!tip]
> You can use the same general approach to generate an `.htaccess` file.

## Example

The following example demonstrates how to configure your site and create a template to automate the generation of a `_redirects` file.

### Content structure

The content structure for this multilingual example looks like this:

```text
content/
├── examples/
│   ├── a.de.md   aliases = ['a-old']
│   ├── a.en.md   aliases = ['a-old', 'a-older']
│   ├── b.de.md   aliases = ['b-old']
│   └── b.en.md   aliases = ['b-old', 'b-older']
└── _index.md
```

In the example above, the aliases are [page-relative](g). To specify a [site-relative](g) path, preface the entry with a slash (`/`). Both forms are resolved to [server-relative](g) paths.

Page-relative paths can also include directory traversal:

| Path type | File path | Alias | Server-relative path |
| :--- | :--- | :--- | :--- |
| page-relative | `content/examples/a.en.md` | `a-old` | `/en/examples/a-old/` |
| page-relative | `content/examples/a.en.md` | `../a-old` | `/en/a-old/` |
| site-relative | `content/examples/a.en.md` | `/a-old` | `/en/a-old/` |

### Site configuration

To implement this, you must update your site configuration to:

1. Disable the generation of default HTML redirect files by setting `disableAliases` to `true`.
1. Define a [media type][] named `text/redirects` to handle the file format.
1. Define a custom [output format][] named `redirects` to set the filename to `_redirects` and place it at the root of the published site.
1. Configure the home page [outputs][] to include the `redirects` format in addition to `html`.

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
disableAliases = true

defaultContentLanguage         = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
  languageCode      = 'en-US'
  languageDirection = 'ltr'
  languageName      = 'English'
  weight            = 1
  title             = 'My Site in English'

[languages.de]
  languageCode      = 'de-DE'
  languageDirection = 'ltr'
  languageName      = 'Deutsch'
  weight            = 2
  title             = 'My Site in German'

[mediaTypes]
  [mediaTypes.'text/redirects']
    delimiter = ''

[outputFormats]
  [outputFormats.redirects]
    baseName    = '_redirects'
    isPlainText = true
    mediaType   = 'text/redirects'
    root        = true

[outputs]
  home = ['html', 'redirects']
{{< /code-toggle >}}

### Template implementation

Next, create a home page template specifically for the `redirects` output format. The following template iterates through every page in every language and extracts its aliases.

To ensure the resulting `_redirects` file is valid, the template uses the [`strings.FindRE`][] function to check for whitespace such as tabs or newlines within the alias string. If whitespace is detected, Hugo will throw an error and fail the build to prevent generating an invalid file.

```go-html-template {file="layouts/home.redirects" copy=true}
{{- if not (hugo.Store.Get "has_printed_redirects") -}}
  {{- range .Sites -}}
    {{- range $p := .Pages -}}
      {{- range .Aliases -}}
        {{- if findRE `\s` . -}}
          {{- errorf "One of the front matter aliases in %q contains whitespace" $p.String -}}
        {{- end -}}
        {{- printf "%s %s 301\n" . $p.RelPermalink -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- hugo.Store.Set "has_printed_redirects" true -}}
{{- end -}}
```

### Generated output

Once Hugo processes the template, it produces a clean list of redirect rules. Each line follows the required format: the source URL, the destination URL, and the HTTP status code.

The resulting `_redirects` file looks like this:

```text
/de/examples/a-old /de/examples/a/ 301
/de/examples/b-old /de/examples/b/ 301
/en/examples/b-old /en/examples/b/ 301
/en/examples/b-older /en/examples/b/ 301
/en/examples/a-old /en/examples/a/ 301
/en/examples/a-older /en/examples/a/ 301
```

[`aliases`]: /content-management/front-matter/#aliases
[`strings.FindRE`]: /functions/strings/findre/
[media type]: /configuration/media-types/
[output format]: /configuration/output-formats/
[outputs]: /configuration/outputs/
