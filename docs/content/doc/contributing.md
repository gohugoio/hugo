---
title: "Contributing to Hugo"
Pubdate: "2013-07-01"
---

We welcome all contributions. If you want to contribute, all 
that is needed is simply fork Hugo, make changes and submit 
a pull request. If you prefer, pick something from the roadmap
or contact [spf13](http://spf13.com) about what may make sense 
to do next.

## Overview 
1. Fork it from https://github.com/spf13/hugo
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


### Clone locally (for contributors):

    git clone https://github.com/spf13/hugo
    cd hugo
    go get

Because go expects all of your libraries to be found in either 
$GOROOT or $GOPATH, it's helpful to symlink the project to one 
of the following paths:

 * ln -s /path/to/your/hugo $GOPATH/src/github.com/spf13/hugo
 * ln -s /path/to/your/hugo $GOROOT/src/pkg/github.com/spf13/hugo

### Running Hugo

    cd /path/to/hugo
    go install github.com/spf13/hugo/hugolibs
    go run main.go

### Building Hugo

    cd /path/to/hugo
    go build -o hugo main.go
    mv hugo /usr/local/bin/

