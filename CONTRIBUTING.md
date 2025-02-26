>**Note:** We would appreciate if you hold on with any big refactoring (like renaming deprecated Go packages), mainly because of potential for extra merge work for future coming in in the near future.

# Contributing to Hugo

We welcome contributions to Hugo of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Hugo community and maintainers are [very active](https://github.com/gohugoio/hugo/pulse/monthly) and helpful, and the project benefits greatly from this activity. We created a [step by step guide](https://gohugo.io/tutorials/how-to-contribute-to-hugo/) if you're unfamiliar with GitHub or contributing to open source projects in general.

*Note that this repository only contains the actual source code of Hugo. For **only** documentation-related pull requests / issues please refer to the [hugoDocs](https://github.com/gohugoio/hugoDocs) repository.*

*Changes to the codebase **and** related documentation, e.g. for a new feature, should still use a single pull request.*

## Table of Contents

* [Asking Support Questions](#asking-support-questions)
* [Reporting Issues](#reporting-issues)
* [Submitting Patches](#submitting-patches)
  * [Code Contribution Guidelines](#code-contribution-guidelines)
  * [Git Commit Message Guidelines](#git-commit-message-guidelines)
  * [Fetching the Sources From GitHub](#fetching-the-sources-from-github)
  * [Building Hugo with Your Changes](#building-hugo-with-your-changes)

## Asking Support Questions

We have an active [discussion forum](https://discourse.gohugo.io) where users and developers can ask questions.
Please don't use the GitHub issue tracker to ask questions.

## Reporting Issues

If you believe you have found a defect in Hugo or its documentation, use
the GitHub issue tracker to report
the problem to the Hugo maintainers. If you're not sure if it's a bug or not,
start by asking in the [discussion forum](https://discourse.gohugo.io).
When reporting the issue, please provide the version of Hugo in use (`hugo
version`) and your operating system.

- [Hugo Issues · gohugoio/hugo](https://github.com/gohugoio/hugo/issues)
- [Hugo Documentation Issues · gohugoio/hugoDocs](https://github.com/gohugoio/hugoDocs/issues)
- [Hugo Website Theme Issues · gohugoio/hugoThemesSite](https://github.com/gohugoio/hugoThemesSite/issues)

## Code Contribution

Hugo has become a fully featured static site generator, so any new functionality must:

* be useful to many.
* fit naturally into _what Hugo does best._
* strive not to break existing sites.
* close or update an open [Hugo issue](https://github.com/gohugoio/hugo/issues)

If it is of some complexity, the contributor is expected to maintain and support the new feature in the future (answer questions on the forum, fix any bugs etc.).

Any non-trivial code change needs to update an open [issue](https://github.com/gohugoio/hugo/issues). A non-trivial code change without an issue reference with one of the labels `bug` or `enhancement` will not be merged.

Note that we do not accept new features that require [CGO](https://github.com/golang/go/wiki/cgo).
We have one exception to this rule which is LibSASS.

**Bug fixes are, of course, always welcome.**

## Submitting Patches

The Hugo project welcomes all contributors and contributions regardless of skill or experience level. If you are interested in helping with the project, we will help you with your contribution.

### Code Contribution Guidelines

Because we want to create the best possible product for our users and the best contribution experience for our developers, we have a set of guidelines which ensure that all contributions are acceptable. The guidelines are not intended as a filter or barrier to participation. If you are unfamiliar with the contribution process, the Hugo team will help you and teach you how to bring your contribution in accordance with the guidelines.

To make the contribution process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes.  We encourage pull requests to allow for review and discussion of code changes.
* When you’re ready to create a pull request, be sure to:
    * Sign the [CLA](https://cla-assistant.io/gohugoio/hugo).
    * Have test cases for the new code. If you have questions about how to do this, please ask in your pull request.
    * Run `go fmt`.
    * Add documentation if you are adding new features or changing functionality.  The docs site lives in `/docs`.
    * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request with `git push -f`.
    * Ensure that `mage check` succeeds. [Travis CI](https://travis-ci.org/gohugoio/hugo) (Windows, Linux and macOS) will fail the build if `mage check` fails.
    * Follow the **Git Commit Message Guidelines** below.

### Git Commit Message Guidelines

This [blog article](https://cbea.ms/git-commit/) is a good resource for learning how to write good commit messages,
the most important part being that each commit message should have a title/subject in imperative mood starting with a capital letter and no trailing period:
*"js: Return error when option x is not set"*, **NOT** *"returning some error."*

Most title/subjects should have a lower-cased prefix with a colon and one whitespace. The prefix can be:

* The name of the package where (most of) the changes are made (e.g. `media: Add text/calendar`)
* If the package name is deeply nested/long, try to shorten it from the left side, e.g. `markup/goldmark` is OK, `resources/resource_transformers/js` can be shortened to `js`.
* If this commit touches several packages with a common functional topic, use that as a prefix, e.g. `errors: Resolve correct line numbers`)
* If this commit touches many packages without a common functional topic, prefix with `all:` (e.g. `all: Reformat Go code`)
* If this is a documentation update, prefix with `docs:`.
* If nothing of the above applies, just leave the prefix out.
* Note that the above excludes nouns seen in other repositories, e.g. "chore:".

Also, if your commit references one or more GitHub issues, always end your commit message body with *See #1234* or *Fixes #1234*.
Replace *1234* with the GitHub issue ID. The last example will close the issue when the commit is merged into *master*.

An example:

```text
tpl: Add custom index function

Add a custom index template function that deviates from the stdlib simply by not
returning an "index out of range" error if an array, slice or string index is
out of range.  Instead, we just return nil values.  This should help make the
new default function more useful for Hugo users.

Fixes #1949
```

###  Fetching the Sources From GitHub

Since Hugo 0.48, Hugo uses the Go Modules support built into Go 1.11 to build. The easiest is to clone Hugo in a directory outside of `GOPATH`, as in the following example:

```bash
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/gohugoio/hugo.git
cd hugo
go install
```

For some convenient build and test targets, you also will want to install Mage:

```bash
go install github.com/magefile/mage
```

Now, to make a change to Hugo's source:

1. Create a new branch for your changes (the branch name is arbitrary):

    ```bash
    git checkout -b iss1234
    ```

1. After making your changes, commit them to your new branch:

    ```bash
    git commit -a -v
    ```

1. Fork Hugo in GitHub.

1. Add your fork as a new remote (the remote name, "fork" in this example, is arbitrary):

    ```bash
    git remote add fork git@github.com:USERNAME/hugo.git
    ```

1. Push the changes to your new remote:

    ```bash
    git push --set-upstream fork iss1234
    ```

1. You're now ready to submit a PR based upon the new branch in your forked repository.

### Building Hugo with Your Changes

Hugo uses [mage](https://github.com/magefile/mage) to sync vendor dependencies, build Hugo, run the test suite and other things. You must run mage from the Hugo directory.

```bash
cd $HOME/go/src/github.com/gohugoio/hugo
```

To build Hugo:

```bash
mage hugo
```

To install hugo in `$HOME/go/bin`:

```bash
mage install
```

To run the tests:

```bash
mage hugoRace
mage -v check
```

To list all available commands along with descriptions:

```bash
mage -l
```

**Note:** From Hugo 0.43 we have added a build tag, `extended` that adds **SCSS support**. This needs a C compiler installed to build. You can enable this when building by:

```bash
HUGO_BUILD_TAGS=extended mage install
````
