---
title: GitInfo
description: Provides access to commit metadata for a given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '*gitmap.GitInfo'
    signatures: [PAGE.GitInfo]
---

The `GitInfo` method on a `Page` object provides access to commit metadata from your Git history, such as the author's name, the commit hash, and the commit message.

> [!note]
> Hugo's Git integration is performant, but may increase build times for large projects.

## Prerequisites

Install Git, create a repository, and commit your project files.

You must also allow Hugo to access your repository by adding this to your project configuration:

{{< code-toggle file=hugo >}}
enableGitInfo = true
{{< /code-toggle >}}

> [!note]
> When you set [`enableGitInfo`][] to `true`, the last modification date for each content page will automatically be the Author Date of the last commit for that file.
>
> This is configurable. See [details][].

## Scope

Commit metadata is available for content stored in your local repository and for content provided by [modules](g).

### Local content

Hugo retrieves commit metadata for files tracked within your project's local repository. This includes all content files managed by Git in your main project directory.

### Module content

{{< new-in 0.157.0 />}}

Hugo also retrieves commit metadata for content provided by modules. This allows you to display commit data for remote repositories that are mounted as content directories, such as when aggregating documentation from multiple sources.

## Methods

### AbbreviatedHash

(`string`) Returns the seven-character shortened version of the commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .AbbreviatedHash }} → aab9ec0
{{ end }}
```

### AuthorDate

(`time.Time`) Returns the date the author originally created the commit.

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

### AuthorEmail

(`string`) Returns the author's email address, respecting [gitmailmap][].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorEmail }} → jsmith@example.org
{{ end }}
```

### AuthorName

(`string`) Returns the author's name, respecting [gitmailmap][].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorName }} → John Smith
{{ end }}
```

### CommitDate

(`time.Time`) Returns the date the commit was applied to the branch.

```go-html-template
{{ with .GitInfo }}
  {{ .CommitDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

### Hash

(`string`) Returns the full SHA-1 commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .Hash }} → aab9ec0b31ebac916a1468c4c9c305f2bebf78d4
{{ end }}
```

### Subject

(`string`) Returns the first line of the commit message (the summary).

```go-html-template
{{ with .GitInfo }}
  {{ .Subject }} → Add tutorials
{{ end }}
```

### Body

(`string`) Returns the full content of the commit message, excluding the subject line.

```go-html-template
{{ with .GitInfo }}
  {{ .Body }} → Two new pages added.
{{ end }}
```

### Ancestors

(`gitmap.GitInfos`) Returns a list of previous commits for this specific file, ordered from most recent to oldest.

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

(`*gitmap.GitInfo`) Returns the most recent ancestor commit for the file, if any.

## Last modified date

By default, when `enableGitInfo` is `true`, the `Lastmod` method on a `Page` object returns the Git AuthorDate of the last commit that included the file.

You can change this behavior in your [project configuration][].

## Hosting considerations

On a [CI/CD](g) platform, the step that clones your project repository must perform a deep clone. If the clone is shallow, the Git information for a given file may be inaccurate. It might incorrectly reflect the most recent repository commit, rather than the commit that actually modified the file.

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

[`enableGitInfo`]: /configuration/all/#enablegitinfo
[details]: /configuration/front-matter/#dates
[gitmailmap]: https://git-scm.com/docs/gitmailmap
[project configuration]: /configuration/front-matter/
