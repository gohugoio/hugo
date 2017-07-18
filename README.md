# Hugo Docs

Documentation site for [Hugo](https://github.com/gohugoio/hugo), the very fast and flexible static site generator built with love in GoLang.

## Build

To view the documentation site locally, you need to clone this repository with submodules:

```bash
git clone --recursive https://github.com/gohugoio/hugoDocs.git
```

Or if you already have a clone locally:

```bash
git submodule update --init
```
Also note that the documentation version for a given version of Hugo can also be found in the `/docs` sub-folder of the [Hugo source repository](https://github.com/gohugoio/hugo).

Then to view the docs in your browser, run Hugo and open up the link:
```bash
hugo serve
Started building sites ...
.
.
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```
