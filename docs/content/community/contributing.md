---
aliases:
- /doc/contributing/
- /meta/contributing/
lastmod: 2015-02-12
date: 2013-07-01
menu:
  main:
    parent: community
next: /tutorials/automated-deployments
prev: /community/mailing-list
title: Contributing to Hugo
weight: 30
---

All contributions to Hugo are welcome. Whether you want to scratch an itch, or simply contribute to the project, feel free to pick something from the [roadmap]({{< relref "meta/roadmap.md" >}}) or contact [spf13](http://spf13.com/) about what may make sense to do next.

You should fork the project and make your changes.  *We encourage pull requests to discuss code changes.*


When you're ready to create a pull request, be sure to:

  * Have test cases for the new code.  If you have questions about how to do it, please ask in your pull request.
  * Run `go fmt`
  * Squash your commits into a single commit.  `git rebase -i`.  It's okay to force update your pull request.
  * Make sure `go test ./...` passes, and `go build` completes.  Our [Travis CI loop](https://travis-ci.org/spf13/hugo) will catch most things that are missing.  The exception: Windows.  We run on Windows from time to time, but if you have access, please check on a Windows machine too.

## Contribution Overview

We wrote a [detailed guide]({{< relref "tutorials/how-to-contribute-to-hugo.md" >}}) for newcomers that guides you step by step to your first contribution. If you are more experienced read on. You probably know what to do.

# Building from source

## Clone locally (for contributors):

    git clone https://github.com/spf13/hugo
    cd hugo
    go get

Because Go expects all of your libraries to be found in either
`$GOROOT` or `$GOPATH`, it's helpful to symlink the project to one
of the following paths:

 * `ln -s /path/to/your/hugo $GOPATH/src/github.com/spf13/hugo`
 * `ln -s /path/to/your/hugo $GOROOT/src/pkg/github.com/spf13/hugo`

## Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugo
    go run main.go

## Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/


# Showcase additions

You got your new website running and it's powered by Hugo? Great. You can add your website with a few steps to the [showcase](/showcase/).

First, make sure that you created a [fork](https://help.github.com/articles/fork-a-repo/) of Hugo on Github and cloned your fork on your local computer. Next, create a separate branch for your additions:

```
# You can choose a different descriptive branch name if you like
git checkout -b showcase-addition
```

Let's create a new document that contains some metadata of your homepage. Replace `example` in the following examples with something unique like the name of your website. Inside the terminal enter the following commands:

```
cd docs
hugo new showcase/example.md
```

You should find the new file at `content/showcase/example.md`. Open it in an editor. The file should contain a frontmatter with predefined variables like below:

```
---
date: 2016-02-12T21:01:18+01:00
description: ""
license: ""
licenseLink: ""
sitelink: http://spf13.com/
sourceLink: https://github.com/spf13/spf13.com
tags:
- personal
- blog
thumbnail: /img/spf13-tn.jpg
title: example
---
```

Add at least values for `sitelink`, `title`,  `description` and a path for `thumbnail`.

Furthermore, we need to create the thumbnail of your website. **It's important that the thumbnail has the required dimensions of 600px by 400px.** Give your thumbnail a name like `example-tn.png` or `example-tn.jpg`. Save it under `docs/static/img/`.

Check a last time that everything works as expected. Start Hugo's built-in server in order to inspect your local copy of the showcase in the browser:

    hugo server

If everything looks fine, we are ready to commit your additions. For the sake of best practices, please make sure that your commit follows our [code contribution guideline](https://github.com/spf13/hugo#code-contribution-guideline).

    git commit -m"docs: Add example.com to the showcase"

Last but not least, we're ready to create a [pull request](https://github.com/spf13/hugo/compare). 

Don't forget to accept the contributor license agreement. Click on the yellow badge in the automatically added comment in the pull request.