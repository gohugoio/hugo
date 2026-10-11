---
_comment: Do not remove front matter.
---

## Prerequisites

Although not required in all cases, [Git][], [Go][], and [Dart Sass][] are commonly used when working with Hugo.

Git is required to:

- Build Hugo from source
- Use [Hugo modules][]
- Install a theme as a Git submodule
- Access [commit information][] from a local Git repository
- Deploy your site to a [CI/CD](g) platform using one of the [hosting guides][]

Go is required to:

- Build Hugo from source
- Use Hugo modules

Dart Sass is required to transpile Sass to CSS when using the latest features of the Sass language.

Please refer to the relevant documentation for installation instructions:

- [Git][git install]
- [Go][go install]
- [Dart Sass][dart sass install]

[Dart Sass]: https://sass-lang.com/dart-sass
[Git]: https://git-scm.com/
[Go]: https://go.dev/
[Hugo modules]: /hugo-modules/
[commit information]: /methods/page/GitInfo/
[dart sass install]: /functions/css/sass/#dart-sass
[git install]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[go install]: https://go.dev/doc/install
[hosting guides]: /host-and-deploy/deployment-methods/#push-to-a-git-repository
