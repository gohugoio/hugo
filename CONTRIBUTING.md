# Contributing to Hugo

We welcome contributions to Hugo of any kind including documentation, themes, organization, tutorials, blog posts, bug reports, issues, feature requests, feature implementation, pull requests, answering questions on the forum, helping to manage issues, etc. 

The Hugo community and maintainers are very active and helpful and the project benefits greatly from this activity.

[![Throughput Graph](https://graphs.waffle.io/spf13/hugo/throughput.svg)](https://waffle.io/spf13/hugo/metrics)

If you have any questions about how to contribute or what to contribute, please ask on the [forum](http://discuss.gohugo.io).

## Code Contribution Guideline

We welcome your contributions. 
To make the process as seamless as possible, we ask for the following:

* Go ahead and fork the project and make your changes. We encourage pull requests to discuss code changes.
* When you’re ready to create a pull request, be sure to:
     * Sign the [CLA](https://cla-assistant.io/spf13/hugo)
     * Have test cases for the new code. If you have questions about how to do it, please ask in your pull request.
     * Run `go fmt`
     * Squash your commits into a single commit. `git rebase -i`. It’s okay to force update your pull request.
     * **Write a good commit message.** This [blog article](http://chris.beams.io/posts/git-commit/) is a good resource for learning how to write good commit messages, the most important part being that each commit message should have a title/subject in imperative mood starting with a capital letter and no trailing period: *"Return error on wrong use of the Paginator"*, **NOT** *"returning some error."* Also, if your commit references one or more GitHub issues, always end your commit message body with *See #1234* or *Fixes #1234*. Replace *1234* with the GitHub issue ID. The last example will close the issue when the commit is merged into *master*.
     * Make sure `go test ./...` passes, and `go build` completes. Our [Travis CI loop](https://travis-ci.org/spf13/hugo) (Linux and OS&nbsp;X) and [AppVeyor](https://ci.appveyor.com/project/spf13/hugo/branch/master) (Windows) will catch most things that are missing.
