---
author: "Michael Henderson"
date: 2015-02-22
linktitle: Installing on Mac
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/installing-on-windows
prev: /tutorials/github-pages-blog
title: Installing on a Mac
weight: 10
---

# Installing Hugo on a Mac

This tutorial aims to be a complete guide to installing Hugo on your Mac computer.

## Assumptions

1. You know how to open a terminal window.
2. You're running a modern 64-bit Mac.
3. You will use `~/Sites` as the starting point for your site.

## Pick Your Method

There are three ways to install Hugo on your Mac computer: the `brew` utility, from the distribution, or from source.
There's no "best" way to do this. You should use the method that works best for your use case.

There are pros and cons for each.

1. `Brew` is the simplest and least work to maintain. The drawbacks
   aren't severe. The default package will be for the most recent
   release, so it will not have bug-fixes until the next release
   (unless you install it with the `--HEAD` option). The release to
   `brew` may lag a few days behind because it has to be coordinated
   with another team. Still, I'd recommend `brew` if you want to work
   from a stable, widely used source. It works well and is really easy
   to update.

2. Downloading the tarball and installing from it is also easy. You have to have a few more command line skills. Updates are easy, too. You just repeat the process with the new binary. This gives you the flexibility to have multiple versions on your computer. If you don't want to use `brew`, then the binary is a good choice.

3. Compiling from source is the most work. The advantage is that you don't have to wait for a release to add features or bug fixes. The disadvantage is that you need to spend more time managing the setup. It's not a lot, but it's more than with the other two options.

Since this is a "beginner" how-to, I'm going to cover the first two
options in detail and go over the third more quickly.

## Brew

### Step 1: Install `brew` if you haven't already

Go to the `brew` website, http://brew.sh/, and follow the directions there. The most important step is:

```
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

When I did this, I had some problems with directory permissions. Searches on Google pointed me to pages that walked me through updating permissions on the `/usr/local` directory. Seemed scary, but it's worked well since.

### Step 2: Run the `brew` command to install `hugo`

Replace `brew install hugo` with `brew install hugo --HEAD` if you
want the absolute latest version—there might be bugs!

```
$ brew install hugo
==> Downloading https://homebrew.bintray.com/bottles/hugo-0.13_1.yosemite.bottle.tar.gz
######################################################################## 100.0%
==> Pouring hugo-0.13_1.yosemite.bottle.tar.gz
🍺  /usr/local/Cellar/hugo/0.13_1: 4 files,  14M
```

`Brew` should have updated your path to include Hugo. Confirm by opening a new terminal window and running a few commands:

```
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

### Step 3: You're Done

You've installed Hugo. Now you need to set up your site. Read the
[Quickstart guide](/overview/quickstart/), explore the rest of the
documentation, and if you still have questions
[just ask!](http://discuss.gohugo.io/ "Discussion forum")

## From Tarball

### Step 1: Decide on the location

When installing from the tarball, you have to decide if you're going to install the binary in `/usr/local/bin` or in your home directory. There are three camps on this:

1. Install it in `/usr/local/bin` so that all the users on your system have access to it. This is a good idea because it's a fairly standard place for exectuables. The downside is that you may need elevated privileges to put software into that location. Also, if there are multiple users on your system, they will all run the same version. Sometimes this can be an issue if you want to try out a new release.

2. Install it in `~/bin` so that only you can execute it. This is a good idea becase it's easy to do, easy to maintain, and doesn't require elevated privileges. The downside is that only you can run Hugo. If there are other users on your site, they have to maintain their own copies. That can lead to people running different versions. of course, this does make it easier for you to experiment with different releases.

3. Install it in your `sites` directory. This is not a bad idea if you have only one site that you're building. It keeps every thing in a single place. If you want to try out new releases, you can just make a copy of the entire site, update the Hugo executable, and have it.

All three locations will work for you. I'm going to document the second option, mostly because I'm comfortable with it.

### Step 2: Download the Tarball

1. Open <https://github.com/spf13/hugo/releases> in your browser.

2. Find the current release by scrolling down and looking for the green tag that reads "Latest Release."

3. Download the current tarball for the Mac. The name will be something like `hugo_X.YY_darwin_amd64.zip`, where `X.YY` is the release number.

4. By default, the tarball will be saved to your `~/Downloads` directory. If you chose to use a different location, you'll need to change that in the following steps.

### Step 3: Confirm your download

Verify that the tarball wasn't corrupted during the download:

```
$ tar tvf ~/Downloads/hugo_0.13_darwin_amd64.zip
-rwxrwxrwx  0 0      0           0 Feb 22 04:02 hugo_0.13_darwin_amd64/hugo_0.13_darwin_amd64
-rwxrwxrwx  0 0      0           0 Feb 22 03:24 hugo_0.13_darwin_amd64/README.md
-rwxrwxrwx  0 0      0           0 Jan 30 18:48 hugo_0.13_darwin_amd64/LICENSE.md
```

The `.md` files are documentation. The other file is the executable.

### Step 4: Install into your bin directory

```
$ # create the directory if needed
$ mkdir -p ~/bin

$ # make it the working directory
$ cd ~/bin

$ # extract the tarball
$ unzip ~/Downloads/hugo_0.13_darwin_amd64.zip
Archive:  hugo_0.13_darwin_amd64.zip
  inflating: hugo_0.13_darwin_amd64/hugo_0.13_darwin_amd64
  inflating: hugo_0.13_darwin_amd64/README.md
  inflating: hugo_0.13_darwin_amd64/LICENSE.md

$ ls -l
total 7704
lrwxr-xr-x  1 mdhender  staff       22 Sep 29 13:34 hugo -> hugo_0.12_darwin_amd/hugo_0.12_darwin_amd64
drwxr-xr-x@ 1 mdhender  staff      102 Sep  1 14:17 hugo_0.12_darwin_amd64
drwxrwxr-x@ 5 mdhender  staff      170 Mar 28 22:46 hugo_0.13_darwin_amd64
-rw-r-----@ 1 mdhender  staff  3942651 Mar 28 22:45 hugo_0.13_darwin_amd64.zip

$ ls -l hugo_0.13_darwin_amd64
total 27560
-rw-r--r--@ 1 mdhender  staff      2707 Jan 30 18:48 LICENSE.md
-rw-r--r--@ 1 mdhender  staff      6748 Feb 22 03:24 README.md
-rwxr-xr-x@ 1 mdhender  staff  14095060 Feb 22 04:02 hugo_0.13_darwin_amd64
```

I was already running Hugo v0.12, so you can see how that was set up. It'll be identical for v0.13 when we're done.

```
$ # create the link to the real executable
$ rm -f hugo
$ ln -s hugo_0.13_darwin_amd64/hugo_0.13_darwin_amd64 hugo
$ ls -l
total 7704
lrwxr-xr-x  1 mdhender  staff       22 Mar 28 22:49 hugo -> hugo_0.13_darwin_amd/hugo_0.12_darwin_amd64
drwxr-xr-x@ 1 mdhender  staff      102 Sep  1 14:17 hugo_0.12_darwin_amd64
drwxrwxr-x@ 5 mdhender  staff      170 Mar 28 22:46 hugo_0.13_darwin_amd64

$ # verify that it runs
$ ./hugo version
Hugo Static Site Generator v0.13 BuildDate: 2015-02-22T04:02:30-06:00
```

You may need to add your bin directory to your `PATH` variable. The `which` command will check for us. If it can find `hugo`, it will print the full path to it. Otherwise, it will not print anything.

```
$ # check if hugo is in the path
$ which hugo
/Users/mdhender/bin/hugo
```

If `hugo` is not in your `PATH`, add it by updating your `~/.bash_profile` file. First, start up an editor:

```
$ nano ~/.bash_profile
```

Add a line to update your `PATH` variable:

```
export PATH=$PATH:$HOME/bin
```

Then save the file by pressing Control-X, then Y to save the file and return to the prompt.

Close the terminal and then open a new terminal to pick up the changes to your profile. Verify by running the `which hugo` command again.

### Step 5: You're Done

You've installed Hugo. Now you need to set up your site. Read the
[Quickstart guide](/overview/quickstart/), explore the rest of the
documentation, and if you still have questions
[just ask!](http://discuss.gohugo.io/ "Discussion forum")

## Building from Source

If you want to compile Hugo yourself, you'll need
[Go](http://golang.org), which is also available from Homebrew: `brew
install go`.

### Step 1: Get the Source

If you want to compile a specific version, go to
<https://github.com/spf13/hugo/releases> and download the source code
for the version of your choice.  If you want to compile Hugo with all
the latest changes (which might include bugs), clone the Hugo
repository:

```
git clone https://github.com/spf13/hugo
```

### Step 2: Compiling

Make the directory containing the source your working directory, then
fetch Hugo's dependencies:

```
mkdir -p src/github.com/spf13
ln -sf $(pwd) src/github.com/spf13/hugo

# set the build path for Go
export GOPATH=$(pwd)

go get
```

This will fetch the absolute latest version of the dependencies, so if
Hugo fails to build it may be because the author of a dependency
introduced a breaking change.

Then compile:

```
go build -o hugo main.go
```

Then place the `hugo` executable somewhere in your `$PATH`.

### Step 3: You're Done

You probably know where to go from here.
