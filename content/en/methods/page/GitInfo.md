---
title: GitInfo
description: Returns Git information related to the last commit of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '*gitmap.GitInfo'
    signatures: [PAGE.GitInfo]
---

The `GitInfo` method on a `Page` object returns an object with additional methods.

> [!note]
> Hugo's Git integration is performant, but may increase build times on large sites.

## Prerequisites

Install Git, create a repository, and commit your project files.

You must also allow Hugo to access your repository. In your site configuration:

{{< code-toggle file=hugo >}}
enableGitInfo = true
{{< /code-toggle >}}

Alternatively, use the command line flag when building your site:

```sh
hugo --enableGitInfo
```

> [!note]
> When you set `enableGitInfo` to `true`, or enable the feature with the command line flag, the last modification date for each content page will be the Author Date of the last commit for that file.
>
> This is configurable. See&nbsp;[details].

## Methods

### AbbreviatedHash

(`string`) The abbreviated commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .AbbreviatedHash }} → aab9ec0b3
{{ end }}
```

### AuthorDate

(`time.Time`) The author date.

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

### AuthorEmail

(`string`) The author's email address, respecting [gitmailmap].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorEmail }} → jsmith@example.org
{{ end }}
```

### AuthorName

(`string`) The author's name, respecting [gitmailmap].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorName }} → John Smith
{{ end }}
```

### CommitDate

(`time.Time`) The commit date.

```go-html-template
{{ with .GitInfo }}
  {{ .CommitDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

### Hash

(`string`) The commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .Hash }} → aab9ec0b31ebac916a1468c4c9c305f2bebf78d4
{{ end }}
```

### Subject

(`string`) The commit message subject.

```go-html-template
{{ with .GitInfo }}
  {{ .Subject }} → Add tutorials
{{ end }}
```

### Body

(`string`) The commit message body.

```go-html-template
{{ with .GitInfo }}
  {{ .Body }} → - Two new pages added.
{{ end }}
```

### Ancestors

(`gitmap.GitInfos`) A slice of file-filtered ancestor commits, if any, ordered from most recent to least recent.

For example, to list the last 5 commits:

```go-html-template
{{ with .GitInfo }}
  {{ range .Ancestors | first 5 }} 
    {{ .CommitDate.Format "2006-01-02" }}: {{ .Subject }}
  {{ end }}
{{ end }}
```

To reverse the order:

```go-html-template
{{ with .GitInfo }}
  {{ range .Ancestors.Reverse | first 5 }} 
    {{ .CommitDate.Format "2006-01-02" }}: {{ .Subject }}
  {{ end }}
{{ end }}
```

### Parent

(`*gitmap.GitInfo`) The first file-filtered ancestor commit, if any.

## Last modified date

By default, when `enableGitInfo` is `true`, the `Lastmod` method on a `Page` object returns the Git AuthorDate of the last commit that included the file.

You can change this behavior in your [site configuration].

## Hosting considerations

In a [CI/CD](g) environment, the step that clones your project repository must perform a deep clone. If the clone is shallow, the Git information for a given file may be inaccurate. It might incorrectly reflect the most recent repository commit, rather than the commit that actually modified the file.

While some providers perform a deep clone by default, others require you to configure the depth yourself.

Hosting service|Default clone depth|Configurable
:--|:--|:--
AWS Amplify|Deep|N/A
Cloudflare|Shallow|Yes [^1]
DigitalOcean App Platform|Deep|N/A
GitHub Pages|Shallow|Yes [^2]
GitLab Pages|Shallow|Yes [^3]
Netlify|Deep|N/A
Render|Shallow|Yes [^1]
Vercel|Shallow|Yes [^1]

[^1]: To perform a deep clone when hosting on Cloudflare, Render, or Vercel, include this code in the build script after the repository has been cloned:

    ```text
    if [ "$(git rev-parse --is-shallow-repository)" = "true" ]; then
      git fetch --unshallow
    fi
    ```

[^2]: To perform a deep clone when hosting on GitHub Pages, set `fetch-depth: 0` in the `checkout` step of the GitHub Action. See [example](/host-and-deploy/host-on-github-pages/#step-7).

[^3]: To perform a deep clone when hosting on GitLab Pages, set the `GIT_DEPTH` environment variable to `0` in the workflow file. See [example](/host-and-deploy/host-on-gitlab-pages/#configure-gitlab-cicd).

[details]: /configuration/front-matter/#dates
[gitmailmap]: https://git-scm.com/docs/gitmailmap
[site configuration]: /configuration/front-matter/
