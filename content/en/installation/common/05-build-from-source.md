## Build from source

To build Hugo from source you must:

1. Install [Git]
1. Install [Go] version 1.18 or later
1. Update your PATH environment variable as described in the [Go documentation]

> The install directory is controlled by the GOPATH and GOBIN environment variables. If GOBIN is set, binaries are installed to that directory. If GOPATH is set, binaries are installed to the bin subdirectory of the first directory in the GOPATH list. Otherwise, binaries are installed to the bin subdirectory of the default GOPATH ($HOME/go or %USERPROFILE%\go).

Then build and test:

```sh
go install -tags extended github.com/gohugoio/hugo@latest
hugo version
```

[Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Go]: https://go.dev/doc/install
[Go documentation]: https://go.dev/doc/code#Command
