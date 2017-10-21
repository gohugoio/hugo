# Contributing to Hugo

We welcome contributions to Hugo of any kind including documentation, themes,
organization, tutorials, blog posts, bug reports, issues, feature requests,
feature implementations, pull requests, answering questions on the forum,
helping to manage issues, etc.

The Hugo community and maintainers are [very active](https://github.com/gohugoio/hugo/pulse/monthly) and helpful, and the project benefits greatly from this activity. We created a [step by step guide](https://gohugo.io/tutorials/how-to-contribute-to-hugo/) if you're unfamiliar with GitHub or contributing to open source projects in general.

*Note that this repository only contains the actual source code of Hugo. For **only** documentation-related pull requests / issues please refer to the [hugoDocs](https://github.com/gohugoio/hugoDocs) repository.*

*Pull requests that contain changes on the code base **and** related documentation, e.g. for a new feature, shall remain a single, atomic one.*

## Table of Contents

* [Asking Support Questions](#asking-support-questions)
* [Reporting Issues](#reporting-issues)
* [Submitting Patches](#submitting-patches)
  * [Code Contribution Guidelines](#code-contribution-guidelines)
  * [Git Commit Message Guidelines](#git-commit-message-guidelines)
  * [Vendored Dependencies](#vendored-dependencies)
  * [Fetching the Sources From GitHub](#fetching-the-sources-from-github)
  * [Using Git Remotes](#using-git-remotes)
  * [Build Hugo with Your Changes](#build-hugo-with-your-changes)
  * [Updating the Hugo Sources](#updating-the-hugo-sources)

## Asking Support Questions

We have an active [discussion forum](https://discourse.gohugo.io) where users and developers can ask questions.
Please don't use the GitHub issue tracker to ask questions.

## Reporting Issues

If you believe you have found a defect in Hugo or its documentation, use
the GitHub [issue tracker](https://github.com/gohugoio/hugo/issues) to report the problem to the Hugo maintainers.
If you're not sure if it's a bug or not, start by asking in the [discussion forum](https://discourse.gohugo.io).
When reporting the issue, please provide the version of Hugo in use (`hugo version`) and your operating system.

## Submitting Patches

The Hugo project welcomes all contributors and contributions regardless of skill or experience level.
If you are interested in helping with the project, we will help you with your contribution.
Hugo is a very active project with many contributions happening daily.
Because we want to create the best possible product for our users and the best contribution experience for our developers,
we have a set of guidelines which ensure that all contributions are acceptable.
The guidelines are not intended as a filter or barrier to participation.
If you are unfamiliar with the contribution process, the Hugo team will help you and teach you how to bring your contribution in accordance with the guidelines.

### Code Contribution Guidelines

To make the contribution process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes.  We encourage pull requests to allow for review and discussion of code changes.
* When you’re ready to create a pull request, be sure to:
    * Sign the [CLA](https://cla-assistant.io/gohugoio/hugo).
    * Have test cases for the new code. If you have questions about how to do this, please ask in your pull request.
    * Run `go fmt`.
    * Add documentation if you are adding new features or changing functionality.  The docs site lives in `/docs`.
    * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request with `git push -f`.
    * Ensure that `mage check` succeeds. [Travis CI](https://travis-ci.org/gohugoio/hugo) (Linux and macOS) and [AppVeyor](https://ci.appveyor.com/project/gohugoio/hugo/branch/master) (Windows) will fail the build if `mage check` fails.
    * Follow the **Git Commit Message Guidelines** below.

### Git Commit Message Guidelines

This [blog article](http://chris.beams.io/posts/git-commit/) is a good resource for learning how to write good commit messages,
the most important part being that each commit message should have a title/subject in imperative mood starting with a capital letter and no trailing period:
*"Return error on wrong use of the Paginator"*, **NOT** *"returning some error."*

Also, if your commit references one or more GitHub issues, always end your commit message body with *See #1234* or *Fixes #1234*.
Replace *1234* with the GitHub issue ID. The last example will close the issue when the commit is merged into *master*.

Sometimes it makes sense to prefix the commit message with the packagename (or docs folder) all lowercased ending with a colon.
That is fine, but the rest of the rules above apply.
So it is "tpl: Add emojify template func", not "tpl: add emojify template func.", and "docs: Document emoji", not "doc: document emoji."

Please consider to use a short and descriptive branch name, e.g. **NOT** "patch-1". It's very common but creates a naming conflict each time when a submission is pulled for a review.

An example:

```text
tpl: Add custom index function

Add a custom index template function that deviates from the stdlib simply by not
returning an "index out of range" error if an array, slice or string index is
out of range.  Instead, we just return nil values.  This should help make the
new default function more useful for Hugo users.

Fixes #1949
```

### Vendored Dependencies

Hugo uses [Go Dep](https://github.com/golang/dep) to vendor dependencies, but we don't commit the vendored packages themselves to the Hugo git repository.
Therefore, a simple `go get` is not supported since `go get` is not vendor-aware.

You **must use Go Dep** to fetch and manage Hugo's dependencies.

###  Fetch the Sources From GitHub

Due to the way Go handles package imports, the best approach for working on a
Hugo fork is to use Git Remotes.  Here's a simple walk-through for getting
started:

1. Install Go Dep and get the Hugo source:

    ```
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v -d github.com/gohugoio/hugo
	```

1. Change to the Hugo source directory and fetch the dependencies:

    ```
    cd $HOME/go/src/github.com/gohugoio/hugo
	dep ensure
    ```

1. Create a new branch for your changes (the branch name is arbitrary):

    ```
    git checkout -b iss1234
    ```

1. After making your changes, commit them to your new branch:

    ```
    git commit -a -v
    ```

1. Fork Hugo in GitHub.

1. Add your fork as a new remote (the remote name, "fork" in this example, is arbitrary):

    ```
    git remote add fork git://github.com/USERNAME/hugo.git
    ```

1. Push the changes to your new remote:

    ```
    git push --set-upstream fork iss1234
    ```

1. You're now ready to submit a PR based upon the new branch in your forked repository.

### Build Hugo with Your Changes

**Note:** Hugo uses [mage](https://github.com/magefile/mage) to build. To install `mage` run

```bash
go get github.com/magefile/mage
```

`mage -l` lists all available commands with the corresponding description. To build Hugo run

```bash
cd $HOME/go/src/github.com/gohugoio/hugo
mage hugo
# or to install in $HOME/go/bin:
mage install
```

### Updating the Hugo Sources

If you want to stay in sync with the Hugo repository, you can easily pull down
the source changes, but you'll need to keep the vendored packages up-to-date as
well.

```
git pull
mage vendor
```

