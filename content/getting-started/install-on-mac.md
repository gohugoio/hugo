---
title: Install on Mac OSX
linktitle: Install on Mac OSX
description:
date: 2016-11-01
publishdate: 2016-11-01
lastmod: 2016-11-01
categories: [getting started]
tags: [install,mac,osx]
weight: 50
draft: false
aliases: []
toc: true
needsreview: true
notesforauthors:
---

## Assumptions

1. You know how to open a terminal window.
2. You're running a modern 64-bit Mac.
3. You will use `~/Sites` as the starting point for your site.

## Pick Your Method

There are three ways to install Hugo on your Mac

1. The [Homebrew][brewlink] `brew` utility
2. Distribution (i.e., tarball)
3. Building from Source

There is no "best" way to install Hugo on your Mac. You should use the method that works best for your use case.

### Pros and Cons

There are pros and cons to each of the aforementioned methods:

1. **Homebrew.** Homebrew is the simplest method and will require the least amount of work to maintain. The drawbacks aren't severe. The default package will be for the most recent release, so it will not have bug fixes until the next release (i.e., unless you install it with the `--HEAD` option). Hugo `brew` releases may lag a few days behind because it has to be coordinated with another team. Nevertheless, `brew` is the recommended installation method if you want to work from a stable, widely used source. Brew works well and is easy to update.

2. **Tarball.** Downloading and installing from the tarball is also easy, although it requires a few more command line skills than does Homebrew. Updates are easy as well: you just repeat the process with the new binary. This gives you the flexibility to have multiple versions on your computer. If you don't want to use `brew`, then the tarball/binary is a good choice.

3. **Building from Source.** Building from source is the most work. The advantage of building from source is that you don't have to wait for a release to add features or bug fixes. The disadvantage is that you need to spend more time managing the setup, which is manageable but requires more time than the preceding two options.

{{% note %}}
Since building from source is appealing to more seasoned command line users, this guide will focus more on installing Hugo via Homebrew or Tarball.
{{% /note %}}

## Installing Hugo with Brew

### Step 1: Install `brew` if you haven't already

Go to the `brew` website, <http://brew.sh/>, and follow the directions there. The most important step is the installation from the command line:

{{% input "install-brew.sh" %}}
```bash
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```
{{% /input %}}

### Step 2: Run the `brew` Command to Install `hugo`

Whenever installing with Homebrew, it's a good idea to update the formulae and Homebrew itself by running the update command:

{{% input "update-brew.sh" %}}
```bash
$ brew update
```
{{% /input %}}

You can then install Hugo using `brew`:

{{% input "install-brew.sh" %}}
```bash
$ brew install hugo
```
{{% /input %}}

If Homebrew is working properly, you should see something similar to the following:

```bash
==> Downloading https://homebrew.bintray.com/bottles/hugo-0.13_1.yosemite.bottle.tar.gz
######################################################################## 100.0%
==> Pouring hugo-0.13_1.yosemite.bottle.tar.gz
🍺  /usr/local/Cellar/hugo/0.13_1: 4 files,  14M
```

{{% note "Installing the Latest Hugo with Brew" %}}
Replace `brew install hugo` with `brew install hugo --HEAD`
if you want the absolute latest version in development.
{{% /note %}}

`brew` should have updated your path to include Hugo. Confirm by opening a new terminal window and running a few commands:

```bash
$ # show the location of the hugo executable
$ which hugo
/usr/local/bin/hugo

$ # show the installed version
$ ls -l $( which hugo )
lrwxr-xr-x  1 mdhender admin  30 Mar 28 22:19 /usr/local/bin/hugo -> ../Cellar/hugo/0.13_1/bin/hugo

$ # verify that hugo runs correctly
$ hugo version
Hugo Static Site Generator v0.13 BuildDate: 2015-03-09T21:34:47-05:00
```

## Installing Hugo from Tarball

### Step 1: Decide on the location

When installing from the tarball, you have to decide if you're going to install the binary in `/usr/local/bin` or in your home directory. There are three camps on this:

1. Install it in `/usr/local/bin` so that all the users on your system have access to it. This is a good idea because it's a fairly standard place for executables. The downside is that you may need elevated privileges to put software into that location. Also, if there are multiple users on your system, they will all run the same version. Sometimes this can be an issue if you want to try out a new release.

2. Install it in `~/bin` so that only you can execute it. This is a good idea because it's easy to do, easy to maintain, and doesn't require elevated privileges. The downside is that only you can run Hugo. If there are other users on your site, they have to maintain their own copies. That can lead to people running different versions. Of course, this does make it easier for you to experiment with different releases.

3. Install it in your `sites` directory. This is not a bad idea if you have only one site that you're building. It keeps every thing in a single place. If you want to try out new releases, you can make a copy of the entire site and update the Hugo executable.

All three locations will work for you. In the interest of brevity, this guide focuses on option #2.

### Step 2: Download the Tarball

1. Open <https://github.com/spf13/hugo/releases> in your browser.

2. Find the current release by scrolling down and looking for the green tag that reads "Latest Release."

3. Download the current tarball for the Mac. The name will be something like `hugo_X.Y_osx-64bit.tgz`, where `X.YY` is the release number.

4. By default, the tarball will be saved to your `~/Downloads` directory. If you choose to use a different location, you'll need to change that in the following steps.

### Step 3: Confirm your download

Verify that the tarball wasn't corrupted during the download:

```bash
$ tar tvf ~/Downloads/hugo_X.Y_osx-64bit.tgz
-rwxrwxrwx  0 0      0           0 Feb 22 04:02 hugo_X.Y_osx-64bit/hugo_X.Y_osx-64bit.tgz
-rwxrwxrwx  0 0      0           0 Feb 22 03:24 hugo_X.Y_osx-64bit/README.md
-rwxrwxrwx  0 0      0           0 Jan 30 18:48 hugo_X.Y_osx-64bit/LICENSE.md
```

The `.md` files are documentation for Hugo. The other file is the executable.

### Step 4: Install Into Your `bin` Directory

```bash
$ # create the directory if needed
$ mkdir -p ~/bin

$ # make it the working directory
$ cd ~/bin

$ # extract the tarball
$ tar -xvzf ~/Downloads/hugo_X.Y_osx-64bit.tgz
Archive:  hugo_X.Y_osx-64bit.tgz
  x ./
  x ./hugo
  x ./LICENSE.md
  x ./README.md

$ # verify that it runs
$ ./hugo version
Hugo Static Site Generator v0.13 BuildDate: 2015-02-22T04:02:30-06:00
```

You may need to add your bin directory to your `PATH` variable. The `which` command will check for us. If it can find `hugo`, it will print the full path to it. Otherwise, it will not print anything.

```bash
$ # check if hugo is in the path
$ which hugo
/Users/USERNAME/bin/hugo
```

If `hugo` is not in your `PATH`, add it by updating your `~/.bash_profile` file. First, start up an editor:

```bash
$ nano ~/.bash_profile
```

Add a line to update your `PATH` variable:

```bash
export PATH=$PATH:$HOME/bin
```

Then save the file by pressing Control-X, then Y to save the file and return to the prompt.

Close the terminal and open a new terminal to pick up the changes to your profile. Verify your success by running the `which hugo` command again.

You've successfully installed Hugo.

## Building from Source

If you want to compile Hugo yourself, you'll need to install Go (aka Golang). You can [install Go directly from the Go website][installgo] or via Homebrew using the following command:

```bash
brew install go
```

### Step 1: Get the Source

If you want to compile a specific version of Hugo, go to <https://github.com/spf13/hugo/releases> and download the source code for the version of your choice. If you want to compile Hugo with all the latest changes (which might include bugs), clone the Hugo repository:

```bash
git clone https://github.com/spf13/hugo
```

{{% warning "Sometimes \"Latest\" = \"Bugs\""%}}
Cloning the Hugo repository directly means taking the good with the bad. By using the bleeding-edge version of Hugo, you make your development susceptible to the latest features, as well as the latest bugs. Your feedback is appreciated. If you find a bug in the latest release, [please create an issue on GitHub](https://github.com/spf13/hugo/issues/new).
{{% /warning %}}

### Step 2: Compiling

Make the directory containing the source your working directory and then fetch Hugo's dependencies:

```bash
mkdir -p src/github.com/spf13
ln -sf $(pwd) src/github.com/spf13/hugo

# set the build path for Go
export GOPATH=$(pwd)

go get
```

This will fetch the absolute latest version of the dependencies. If Hugo fails to build, it may be the result of a dependency's author introducing a breaking change.

Once you have properly configured your directory, you can compile Hugo using the following command:

```bash
go build -o hugo main.go
```

Then place the `hugo` executable somewhere in your `$PATH`. You're now ready to start using Hugo.

## Next Steps

Now that you've installed Hugo, read the [Quickstart guide][quickstart] and explore the rest of the documentation. If you have questions, ask the Hugo community directly by visiting the [Hugo Discussion Forum][hugodiscussion].

[brewlink]: https://brew.sh/
[hugodiscussion]: https://discuss.gohugo.io "Visit the Hugo Discussion forum to tap into the community's collective knowledge about Hugo."
[installgo]: https://golang.org/dl/
[quickstart]: /getting-started/quick-start/