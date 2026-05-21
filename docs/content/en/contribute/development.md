---
title: Development
description: Contribute to the development of Hugo.
categories: []
keywords: []
---

## Introduction

You can contribute to the Hugo project by:

- Answering questions on the [forum]
- Improving the [documentation]
- Monitoring the [issue queue]
- Creating or improving [themes]
- Squashing [bugs]

Please submit documentation issues and pull requests to the [documentation repository].

If you have an idea for an enhancement or new feature, create a new topic on the [forum] in the "Feature" category. This will help you to:

- Determine if the capability already exists
- Measure interest
- Refine the concept

If there is sufficient interest, [create a proposal]. Do not submit a pull request until the project lead accepts the proposal.

For a complete guide to contributing to Hugo, see the [Contribution Guide].

## Prerequisites

To build Hugo from source you must install:

1. Install [Git]
1. Install [Go] version 1.25.0 or later

## GitHub workflow

> [!note]
> This section assumes that you have a working knowledge of Go, Git and GitHub, and are comfortable working on the command line.

Use this workflow to create and submit pull requests.

Step 1
: Fork the [project repository].

Step 2
: Clone your fork.

Step 3
: Create a new branch with a descriptive name that includes the corresponding issue number.

  For a new feature:

  ```sh
  git checkout -b feat/implement-some-feature-99999
  ```

  For a bug fix:

  ```sh
  git checkout -b fix/fix-some-bug-99999
  ```

Step 4
: Make changes.

Step 5
: Build and install.

  To build and install the standard edition:

  ```sh
  CGO_ENABLED=0 go install
  ```

  {{< new-in v0.159.2 />}} To build and install the deploy edition:

  ```sh
  CGO_ENABLED=0 go install -tags withdeploy
  ```

  To build and install the extended edition, first install a C compiler such as [GCC] or [Clang] and then run the following command:

  ```sh
  CGO_ENABLED=1 go install -tags extended
  ```

  To build and install the extended/deploy edition, first install a C compiler such as [GCC] or [Clang] and then run the following command:

  ```sh
  CGO_ENABLED=1 go install -tags extended,withdeploy
  ```

Step 6
: Test your changes:

  ```text
  go test ./...
  ```

Step 7
: Commit your changes with a descriptive commit message:

  - Provide a summary on the first line, typically 50 characters or less, followed by a blank line.
    - Begin the summary with the name of the package, followed by a colon, a space, and a brief description of the change beginning with a capital letter
    - Use imperative present tense
    - See the [commit message guidelines] for requirements
  - Optionally, provide a detailed description where each line is 72 characters or less, followed by a blank line.
  - Add one or more "Fixes" or "Closes" keywords, each on its own line, referencing the [issues] addressed by this change.

  For example:

  ```sh
  git commit -m "tpl/strings: Create wrap function

  The strings.Wrap function wraps a string into one or more lines,
  splitting the string after the given number of characters, but not
  splitting in the middle of a word.

  Fixes #99998
  Closes #99999"
  ```

Step 8
: Push the new branch to your fork of the documentation repository.

Step 9
: Visit the [project repository] and create a pull request (PR).

Step 10
: A project maintainer will review your PR and may request changes. You may delete your branch after the maintainer merges your PR.

[Clang]: https://clang.llvm.org/
[Contribution Guide]: https://github.com/gohugoio/hugo/blob/master/CONTRIBUTING.md
[GCC]: https://gcc.gnu.org/
[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[bugs]: https://github.com/gohugoio/hugo/issues?q=is%3Aopen+is%3Aissue+label%3ABug
[commit message guidelines]: https://github.com/gohugoio/hugo/blob/master/CONTRIBUTING.md#git-commit-message-guidelines
[create a proposal]: https://github.com/gohugoio/hugo/issues/new?labels=Proposal%2C+NeedsTriage&template=feature_request.md
[documentation repository]: https://github.com/gohugoio/hugoDocs
[documentation]: /documentation
[forum]: https://discourse.gohugo.io
[issue queue]: https://github.com/gohugoio/hugo/issues
[issues]: https://github.com/gohugoio/hugo/issues
[project repository]: https://github.com/gohugoio/hugo/
[themes]: https://themes.gohugo.io/
