# Hugo Docs
This is the Documentation site for [Hugo](https://github.com/gohugoio/hugo), the world’s fastest static website engine built with love in GoLang.

You can view the documentation on the website here: https://gohugo.io/documentation/

## Contributing

We welcome contributions to Hugo of any kind including documentation, suggestions, bug reports, pull requests etc. Also check out our [contribution guide](https://gohugo.io/contribute/documentation/). We would love to hear from you. 

Note that this repository contains solely the documentation for Hugo. For contributions that aren't documentation-related, please refer to the [Hugo repository](https://github.com/gohugoio/hugo). 

*Pull requests shall **only** contain changes to the actual documentation. However, changes on the code base of Hugo **and** the documentation shall be a single, atomic pull request in the [hugo repository](https://github.com/gohugoio/hugo).*

## Branches

* The `master` branch is where the site is automatically built from, and is the place to put changes relevant to the current Hugo version.
* The `next` branch is where we store changes and code base related to the next Hugo release. This can be previewed here: https://next--gohugoio.netlify.com

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
