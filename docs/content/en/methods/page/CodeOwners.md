---
title: CodeOwners
description: Returns of slice of code owners for the given page, derived from the CODEOWNERS file in the root of the project directory.
categories: []
keywords: []
action:
  related:
    - methods/page/GitInfo
  returnType: '[]string'
  signatures: [PAGE.CodeOwners]
---

GitHub and GitLab support CODEOWNERS files. This file specifies the users responsible for developing and maintaining software and documentation. This definition can apply to the entire repository, specific directories, or to individual files. To learn more:

- [GitHub CODEOWNERS documentation]
- [GitLab CODEOWNERS documentation]

Use the `CodeOwners` method on a `Page` object to determine the code owners for the given page.

[GitHub CODEOWNERS documentation]: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners
[GitLab CODEOWNERS documentation]: https://docs.gitlab.com/ee/user/project/code_owners.html

To use the `CodeOwners` method you must enable access to your local Git repository:

{{< code-toggle file=hugo >}}
enableGitInfo = true
{{< /code-toggle >}}

Consider this project structure:

```text
my-project/
├── content/
│   ├── books/
│   │   └── les-miserables.md
│   └── films/
│       └── the-hunchback-of-notre-dame.md
└── CODEOWNERS
```

And this CODEOWNERS file:

```text
* @jdoe
/content/books/ @tjones
/content/films/ @mrichards @rsmith
```

The table below shows the slice of code owners returned for each file:

Path|Code owners
:--|:--
`books/les-miserables.md`|`[@tjones]`
`films/the-hunchback-of-notre-dame.md`|`[@mrichards @rsmith]`

Render the code owners for each content page:

```go-html-template
{{ range .CodeOwners }}
  {{ . }}
{{ end }}
```

Combine this method with [`resources.GetRemote`] to retrieve names and avatars from your Git provider by querying their API.

[`resources.GetRemote`]: /functions/resources/getremote/
