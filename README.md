[![Netlify Status](https://api.netlify.com/api/v1/badges/e0dbbfc7-34f1-4393-a679-c16e80162705/deploy-status)](https://app.netlify.com/sites/gohugoio/deploys)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://gohugo.io/contribute/documentation/)

# Hugo Docs

Documentation site for [Hugo](https://github.com/gohugoio/hugo), the very fast and flexible static site generator built with love in Go.

## Contributing

We welcome contributions to Hugo of any kind including documentation, suggestions, bug reports, pull requests etc. Also check out our [contribution guide](https://gohugo.io/contribute/documentation/). We would love to hear from you. 

Note that this repository contains solely the documentation for Hugo. For contributions that aren't documentation-related please refer to the [hugo](https://github.com/gohugoio/hugo) repository. 

*Pull requests shall **only** contain changes to the actual documentation. However, changes on the code base of Hugo **and** the documentation shall be a single, atomic pull request in the [hugo](https://github.com/gohugoio/hugo) repository.*

Spelling fixes are most welcomed, and if you want to contribute longer sections to the documentation, it would be great if you had the following criteria in mind when writing:

* Short is good. People go to the library to read novels. If there is more than one way to _do a thing_ in Hugo, describe the current _best practice_ (avoid "… but you can also do …" and "… in older versions of Hugo you had to …".
* For example, try to find short snippets that teaches people about the concept. If the example is also useful as-is (copy and paste), then great. Don't list long and similar examples just so people can use them on their sites.
* Hugo has users from all over the world, so easy to understand and [simple English](https://simple.wikipedia.org/wiki/Basic_English) is good.

## Branches

* The `master` branch is where the site is automatically built from, and is the place to put changes relevant to the current Hugo version.
* The `next` branch is where we store changes that are related to the next Hugo release. This can be previewed here: https://next--gohugoio.netlify.com/

## Build

To view the documentation site locally, you need to clone this repository:

```bash
git clone https://github.com/gohugoio/hugoDocs.git
```

Also note that the documentation version for a given version of Hugo can also be found in the `/docs` sub-folder of the [Hugo source repository](https://github.com/gohugoio/hugo).

Then to view the docs in your browser, run Hugo and open up the link:

```bash
▶ hugo server

Started building sites ...
.
.
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```
