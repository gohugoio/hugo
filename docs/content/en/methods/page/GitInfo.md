---
title: GitInfo
description: Returns Git information related to the last commit of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/CodeOwners
  returnType: source.GitInfo
  signatures: [PAGE.GitInfo]
toc: true
---

The `GitInfo` method on a `Page` object returns an object with additional methods.

{{% note %}}
Hugo's Git integration is performant, but may increase build times on large sites.
{{% /note %}}

## Prerequisites

Install [Git], create a repository, and commit your project files.

You must also allow Hugo to access your repository. In your site configuration:

{{< code-toggle file=hugo >}}
enableGitInfo = true
{{< /code-toggle >}}

Alternatively, use the command line flag when building your site:

```sh
hugo --enableGitInfo
```

{{% note %}}
When you set `enableGitInfo` to `true`, or enable the feature with the command line flag, the last modification date for each content page will be the Author Date of the last commit for that file.

This is configurable. See&nbsp;[details].

[details]: /getting-started/configuration/#configure-dates
{{% /note %}}

## Methods

###### AbbreviatedHash

(`string`) The abbreviated commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .AbbreviatedHash }} → aab9ec0b3
{{ end }}
```

###### AuthorDate

(`time.Time`) The author date.

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

###### AuthorEmail

(`string`) The author's email address, respecting [gitmailmap].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorEmail }} → jsmith@example.org
{{ end }}
```

###### AuthorName

(`string`) The author's name, respecting [gitmailmap].

```go-html-template
{{ with .GitInfo }}
  {{ .AuthorName }} → John Smith
{{ end }}
```

###### CommitDate

(`time.Time`) The commit date.

```go-html-template
{{ with .GitInfo }}
  {{ .CommitDate.Format "2006-01-02" }} → 2023-10-09
{{ end }}
```

###### Hash

(`string`) The commit hash.

```go-html-template
{{ with .GitInfo }}
  {{ .Hash }} → aab9ec0b31ebac916a1468c4c9c305f2bebf78d4
{{ end }}
```

###### Subject

(`string`) The commit message subject.

```go-html-template
{{ with .GitInfo }}
  {{ .Subject }} → Add tutorials
{{ end }}
```

###### Body

(`string`) The commit message body.

```go-html-template
{{ with .GitInfo }}
  {{ .Body }} → - Two new pages added.
{{ end }}
```

## Last modified date

By default, when `enableGitInfo` is `true`, the `Lastmod` method on a `Page` object returns the Git AuthorDate of the last commit that included the file.

You can change this behavior in your [site configuration].

[git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[gitmailmap]: https://git-scm.com/docs/gitmailmap
[site configuration]: /getting-started/configuration/#configure-front-matter

## Hosting considerations

When hosting your site in a CI/CD environment, the step that clones your project repository must perform a deep clone. If the clone is shallow, the Git information for a given file may not be accurate---it may reflect the most recent repository commit, not the commit that last modified the file.

Some providers perform deep clones by default, others allow you to configure the clone depth, and some providers only perform shallow clones.

Hosting service | Default clone depth | Configurable
:-- | :-- | :--
Cloudflare Pages | Shallow | Yes [^CFP]
DigitalOcean App Platform | Deep | N/A
GitHub Pages | Shallow | Yes [^GHP]
GitLab Pages | Shallow | Yes [^GLP]
Netlify | Deep | N/A
Render | Shallow | No
Vercel | Shallow | No

[^CFP]: To configure a Cloudflare Pages site for deep cloning, preface the site's normal Hugo build command with `git fetch --unshallow &&` (*e.g.*, `git fetch --unshallow && hugo`).

[^GHP]: You can configure the GitHub Action to do a deep clone by specifying `fetch-depth: 0` in the applicable "checkout" step of your workflow file, as shown in the Hugo documentation's [example workflow file](/hosting-and-deployment/hosting-on-github/#procedure).

[^GLP]: You can configure the GitLab Runner's clone depth [as explained in the GitLab documentation](https://docs.gitlab.com/ee/ci/large_repositories/#shallow-cloning); see also the Hugo documentation's [example workflow file](/hosting-and-deployment/hosting-on-gitlab/#configure-gitlab-cicd).
