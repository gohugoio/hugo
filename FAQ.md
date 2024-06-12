# FAQ

## General Questions

### What is Hugo?

Hugo is a fast and flexible static site generator written in Go. It is optimized for speed and designed for flexibility, allowing you to build websites quickly and efficiently.

### Who created Hugo?

Hugo was created by [bep], [spf13], and [friends] in [Go].

### Where can I find the documentation?

You can find the documentation on the [Hugo website](https://gohugo.io/documentation).

## Installation

### How do I install Hugo?

You can install Hugo from a [prebuilt binary], package manager, or package repository. Installation instructions for different operating systems are available on the [installation page](https://gohugo.io/installation).

### What are the prerequisites for building Hugo from source?

For the standard edition, you need Go 1.20 or later. For the extended edition, you also need GCC.

### How do I build Hugo from source?

For the standard edition:
```sh
go install github.com/gohugoio/hugo@latest
```

## For the extended edition:
```sh
CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
```

## Usage
### How do I create a new Hugo site?
Use the following command to create a new Hugo site:

```sh
hugo new site <sitename>
```

### How do I add a new content file?
Use the following command to create a new content file:
```sh
hugo new <section>/<filename>.<format>
```
Replace <section>, <filename>, and <format> with your desired section, filename, and format (e.g., markdown).

### How do I start the Hugo server?

```sh
hugo server
```
This will start a local server and you can view your site at http://localhost:1313.



