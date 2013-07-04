{
    "title": "Installing Hugo",
    "Pubdate": "2013-07-01"
}

Installation is very easy. Simply download the appropriate version for your
platform. 

Hugo is written in GoLang with support for Windows, Linux and OSX.

<div class="alert alert-info">
Please make sure that you place the executable in your path. `/usr/local/bin` 
is the most probable location.
</div>


Hugo doesn't have any external dependencies, but can benefit from external
programs.


## Installing from source

Make sure you have a recent version of go installed. Hugo requires go 1.1+.

    git clone https://github.com/spf13/hugo
    cd hugo
    go build -o hugo main.go

