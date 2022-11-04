---
title: Install Hugo
linktitle: Install Hugo
description: Install Hugo on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine where the Go compiler tool chain can run.
date: 2016-11-01
publishdate: 2016-11-01
categories: [getting started,fundamentals]
authors: ["Michael Henderson"]
keywords: [install,pc,windows,linux,macos,binary,tarball]
menu:
  docs:
    parent: "getting-started"
    weight: 30
weight: 30
sections_weight: 30
aliases: [/tutorials/installing-on-windows/,/tutorials/installing-on-mac/,/overview/installing/,/getting-started/install,/install/]
toc: true
---

{{% note %}}
There is lots of talk about "Hugo being written in Go", but you don't need to install Go to enjoy Hugo. Just grab a precompiled binary!
{{% /note %}}

Hugo is written in [Go](https://go.dev/) with support for multiple platforms. The latest release can be found at [Hugo Releases][releases].

Hugo currently provides pre-built binaries for the following:

* macOS (Darwin) for x64, i386, and ARM architectures
* Windows
* Linux
* OpenBSD
* FreeBSD

Hugo may also be compiled from source wherever the Go toolchain can run; e.g., on other operating systems such as DragonFly BSD, OpenBSD, Plan&nbsp;9, Solaris, and others. See <https://go.dev/doc/install/source> for the full set of supported combinations of target operating systems and compilation architectures.

## Quick Install

### Binary (Cross-platform)

Download the appropriate version for your platform from [Hugo Releases][releases]. Once downloaded, the binary can be run from anywhere. You don't need to install it into a global location. This works well for shared hosts and other systems where you don't have a privileged account.

Ideally, you should install it somewhere in your `PATH` for easy use. `/usr/local/bin` is the most probable location.

### Docker

We currently do not offer official Hugo images for Docker, but we do recommend these up to date distributions: https://hub.docker.com/r/klakegg/hugo/

### Homebrew (macOS)

If you are on macOS and using [Homebrew][brew], you can install Hugo with the following one-liner:

{{< code file="install-with-homebrew.sh" >}}
brew install hugo
{{< /code >}}

For more detailed explanations, read the installation guides that follow for installing on macOS and Windows.

### MacPorts (macOS)

If you are on macOS and using [MacPorts][macports], you can install Hugo with the following one-liner:

{{< code file="install-with-macports.sh" >}}
port install hugo
{{< /code >}}

### Homebrew (Linux)

If you are using [Homebrew][linuxbrew] on Linux, you can install Hugo with the following one-liner:

{{< code file="install-with-linuxbrew.sh" >}}
brew install hugo
{{< /code >}}

Installation guides for Homebrew on Linux are available on their [website][linuxbrew].

### Chocolatey (Windows)

If you are on a Windows machine and use [Chocolatey][] for package management, you can install Hugo with the following one-liner:

{{< code file="install-with-chocolatey.ps1" >}}
choco install hugo -confirm
{{< /code >}}

Or if you need the “extended” Sass/SCSS version:

{{< code file="install-extended-with-chocolatey.ps1" >}}
choco install hugo-extended -confirm
{{< /code >}}

### Scoop (Windows)

If you are on a Windows machine and use [Scoop][] for package management, you can install Hugo with the following one-liner:

```bash
scoop install hugo
```

Or install the extended version with:

```bash
scoop install hugo-extended
```

### Source

#### Prerequisite Tools

* [Git][installgit]
* [GCC][] (For Windows users only)
* [Go (at least Go 1.11)](https://go.dev/dl/)

#### Fetch from GitHub

Since Hugo 0.48, Hugo uses the Go Modules support built into Go 1.11 to build. The easiest way to get started is to clone Hugo in a directory outside of the GOPATH, as in the following example:

{{< code file="from-gh.sh" >}}
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/gohugoio/hugo.git
cd hugo
go install --tags extended
{{< /code >}}

Remove `--tags extended` if you do not want/need Sass/SCSS support.

{{% note %}}

##### For installation on Windows

* Substitute the `$HOME` environment variable above with `%USERPROFILE%`.
* If you install `--tags extended` version, you may encounter this error `"gcc": executable file not found in %PATH%`
  * Please make sure you have installed `gcc` command and add it to `%PATH%`.
  * "MinGW" is recommended, it has been tested and built successfully

{{% /note %}}

## macOS

### Assumptions

1. You know how to open the macOS terminal.
2. You're running a modern 64-bit Mac.
3. You will use `~/Sites` as the starting point for your site. (`~/Sites` is used for example purposes. If you are familiar enough with the command-line and file system, you should have no issues following along with the instructions.)

### Pick Your Method

There are three ways to install Hugo on your Mac

1. A package manager, like [Homebrew][brew] (`brew`) or [MacPorts][macports] (`port`)
2. Distribution (i.e., tarball)
3. Building from Source

There is no "best" way to install Hugo on your Mac. You should use the method that works best for your use case.

#### Pros and Cons

There are pros and cons to each of the aforementioned methods:

1. **Package Manager.** Using a package manager is the simplest method and will require the least amount of work to maintain. The drawbacks aren't severe. The default package will be for the most recent release, so it will not have bug fixes until the next release (i.e., unless you install it with the `--HEAD` option in Homebrew). Releases may lag a few days behind because it has to be coordinated with another team. Nevertheless, this is the recommended installation method if you want to work from a stable, widely used source. Package managers work well and they are easy to update.

2. **Tarball.** Downloading and installing from the tarball is also easy, although it requires a few more command-line skills than does Homebrew. Updates are easy as well: you just repeat the process with the new binary. This gives you the flexibility to have multiple versions on your computer. If you don't want to use `brew`, then the tarball/binary is a good choice.

3. **Building from Source.** Building from source is the most work. The advantage of building from source is that you don't have to wait for a release to add features or bug fixes. The disadvantage is that you need to spend more time managing the setup, which is manageable but requires more time than the preceding two options.

{{% note %}}
Since building from source is appealing to more seasoned command-line users, this guide will focus more on installing Hugo via Homebrew and Tarball.
{{% /note %}}

### Install Hugo with Brew

{{< youtube WvhCGlLcrF8 >}}

#### Step 1: Install `brew` if you haven't already

Go to the `brew` website, <https://brew.sh>, and follow the directions there. The most important step is the installation from the command line:

{{< code file="install-brew.sh" >}}
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
{{< /code >}}

#### Step 2: Run the `brew` Command to Install `hugo`

Installing Hugo using `brew` is as easy as the following:

{{< code file="install-brew.sh" >}}
brew install hugo
{{< /code >}}

If Homebrew is working properly, you should see something similar to the following:

```txt
==> Downloading https://homebrew.bintray.com/bottles/hugo-0.21.sierra.bottle.tar.gz
######################################################################### 100.0%
==> Pouring hugo-0.21.sierra.bottle.tar.gz
🍺  /usr/local/Cellar/hugo/0.21: 32 files, 17.4MB
```

{{% note "Installing the Latest Hugo with Brew" %}}
Replace `brew install hugo` with `brew install hugo --HEAD` if you want the absolute latest in-development version.
{{% /note %}}

`brew` should have updated your path to include Hugo. You can confirm by opening a new terminal window and running a few commands:

```txt
$ # show the location of the hugo executable
which hugo
/usr/local/bin/hugo

# show the installed version
ls -l $( which hugo )
lrwxr-xr-x  1 mdhender admin  30 Mar 28 22:19 /usr/local/bin/hugo -> ../Cellar/hugo/0.13_1/bin/hugo

# verify that hugo runs correctly
hugo version
Hugo Static Site Generator v0.13 BuildDate: 2015-03-09T21:34:47-05:00
```

### Install Hugo from Tarball

#### Step 1: Decide on the location

When installing from the tarball, you have to decide if you're going to install the binary in `/usr/local/bin` or in your home directory. There are three camps on this:

1. Install it in `/usr/local/bin` so that all the users on your system have access to it. This is a good idea because it's a fairly standard place for executables. The downside is that you may need elevated privileges to put software into that location. Also, if there are multiple users on your system, they will all run the same version. Sometimes this can be an issue if you want to try out a new release.

2. Install it in `~/bin` so that only you can execute it. This is a good idea because it's easy to do, easy to maintain, and doesn't require elevated privileges. The downside is that only you can run Hugo. If there are other users on your site, they have to maintain their own copies. That can lead to people running different versions. Of course, this does make it easier for you to experiment with different releases.

3. Install it in your `Sites` directory. This is not a bad idea if you have only one site that you're building. It keeps every thing in a single place. If you want to try out new releases, you can make a copy of the entire site and update the Hugo executable.

All three locations will work for you. In the interest of brevity, this guide focuses on option #2.

#### Step 2: Download the Tarball

1. Open <https://github.com/gohugoio/hugo/releases> in your browser.

2. Find the current release by scrolling down and looking for the green tag that reads "Latest Release."

3. Download the current tarball for the Mac. The name will be something like `hugo_X.Y_osx-64bit.tgz`, where `X.YY` is the release number.

4. By default, the tarball will be saved to your `~/Downloads` directory. If you choose to use a different location, you'll need to change that in the following steps.

#### Step 3: Confirm your download

Verify that the tarball wasn't corrupted during the download:

```txt
tar tvf ~/Downloads/hugo_X.Y_osx-64bit.tgz
-rwxrwxrwx  0 0      0           0 Feb 22 04:02 hugo_X.Y_osx-64bit/hugo_X.Y_osx-64bit.tgz
-rwxrwxrwx  0 0      0           0 Feb 22 03:24 hugo_X.Y_osx-64bit/README.md
-rwxrwxrwx  0 0      0           0 Jan 30 18:48 hugo_X.Y_osx-64bit/LICENSE.md
```

The `.md` files are documentation for Hugo. The other file is the executable.

#### Step 4: Install Into Your `bin` Directory

```txt
# create the directory if needed
mkdir -p ~/bin

# make it the working directory
cd ~/bin

# extract the tarball
tar -xvzf ~/Downloads/hugo_X.Y_osx-64bit.tgz
Archive:  hugo_X.Y_osx-64bit.tgz
  x ./
  x ./hugo
  x ./LICENSE.md
  x ./README.md

# verify that it runs
./hugo version
Hugo Static Site Generator v0.13 BuildDate: 2015-02-22T04:02:30-06:00
```

You may need to add your bin directory to your `PATH` environment variable. The `which` command will check for us. If it can find `hugo`, it will print the full path to it. Otherwise, it will not print anything.

```txt
# check if hugo is in the path
which hugo
/Users/USERNAME/bin/hugo
```

If `hugo` is not in your `PATH`:

1. Determine your default shell (zsh or bash).

   ```bash
   echo $SHELL
   ```

2. Edit your profile.

   If your default shell is zsh:

    ```zsh
    nano ~/.zprofile
    ```

    If your default shell is bash:

    ```bash
    nano ~/.bash_profile
    ```

3. Insert a line to add `$HOME/bin` to your existing `PATH`.

    ```txt
    export PATH=$PATH:$HOME/bin
    ```

4. Save the file by pressing Control-X, then Y.

5. Close the terminal and open a new terminal to pick up the changes to your profile. Verify the change by running the `which hugo` command again.

You've successfully installed Hugo.

### Build from Source on Mac

If you want to compile Hugo yourself, you'll need to install Go (aka Golang). You can [install Go directly from the Go website](https://go.dev/dl/) or via Homebrew using the following command:

```txt
brew install go
```

#### Step 1: Get the Source

If you want to compile a specific version of Hugo, go to <https://github.com/gohugoio/hugo/releases> and download the source code for the version of your choice. If you want to compile Hugo with all the latest changes (which might include bugs), clone the Hugo repository:

```txt
git clone https://github.com/gohugoio/hugo
```

{{% warning "Sometimes \"Latest\" = \"Bugs\""%}}
Cloning the Hugo repository directly means taking the good with the bad. By using the bleeding-edge version of Hugo, you make your development susceptible to the latest features, as well as the latest bugs. Your feedback is appreciated. If you find a bug in the latest release, [please create an issue on GitHub](https://github.com/gohugoio/hugo/issues/new).
{{% /warning %}}

#### Step 2: Compiling

Make the directory containing the source your working directory and then fetch Hugo's dependencies:

```txt
mkdir -p src/github.com/gohugoio
ln -sf $(pwd) src/github.com/gohugoio/hugo

go get
```

This will fetch the absolute latest version of the dependencies. If Hugo fails to build, it may be the result of a dependency's author introducing a breaking change.

Once you have properly configured your directory, you can compile Hugo using the following command:

```txt
go build -o hugo main.go
```

Then place the `hugo` executable somewhere in your `$PATH`. You're now ready to start using Hugo.

## Windows

The following aims to be a complete guide to installing Hugo on your Windows PC.

{{< youtube G7umPCU-8xc >}}

### Assumptions for Windows

1. You will use `C:\Hugo\Sites` as the starting point for your new project.
2. You will use `C:\Hugo\bin` to store executable files.

### Set up Your Directories

You'll need a place to store the Hugo executable, your [content][], and the generated Hugo website:

1. Open Windows Explorer.
2. Create a new folder: `C:\Hugo`, assuming you want Hugo on your C drive, although this can go anywhere
3. Create a subfolder in the Hugo folder: `C:\Hugo\bin`
4. Create another subfolder in Hugo: `C:\Hugo\Sites`

### Technical Users

1. Download the latest zipped Hugo executable from [Hugo Releases][releases].
2. Extract all contents to your `..\Hugo\bin` folder.
3. Open Windows command-line (cmd, "DOS") to add the `hugo.exe` executable to your PATH
    * do `set PATH=%PATH%;C:\Hugo\bin` to have hugo in PATH for the currently opened cmd box
    * do `setx PATH "%PATH%;C:\Hugo\bin"` to have hugo in PATH for every newly opened cmd box
      * note: "setx", not "set", plus syntax 'key "val"', not 'key=val'

> You may also use "Git CMD" from the [Git for Windows package](https://gitforwindows.org/) for the native Windows commands [set](https://ss64.com/nt/set.html) and [setx](https://ss64.com/nt/setx.html), but not "Git Bash", PowerShell, or any other "CLI" with different commands

### Less-technical Users

1. Go to the [Hugo Releases][releases] page.
2. The latest release is announced on top. Scroll to the bottom of the release announcement to see the downloads. They're all ZIP files.
3. Find the Windows files near the bottom (they're in alphabetical order, so Windows is last) – download either the 32-bit or 64-bit file depending on whether you have 32-bit or 64-bit Windows. (If you don't know, [see here](https://esupport.trendmicro.com/en-us/home/pages/technical-support/1038680.aspx).)
4. Move the ZIP file into your `C:\Hugo\bin` folder.
5. Double-click on the ZIP file and extract its contents. Be sure to extract the contents into the same `C:\Hugo\bin` folder – Windows will do this by default unless you tell it to extract somewhere else.
6. You should now have three new files: The hugo executable (`hugo.exe`), `LICENSE`, and `README.md`.

Now you need to add Hugo to your Windows PATH settings:

#### For Windows 10 Users

* Right click on the **Start** button.
* Click on **System**.
* Click on **Advanced System Settings** on the right.
* Click on the **Environment Variables...** button on the bottom.
* In the User variables section, select the row labeled "Path" and click the **Edit...** button.
* Click the **Browse...** button and select the directory to which `hugo.exe` was extracted, which is `C:\Hugo\bin` if you went by the instructions above. *The path entry should be the folder where Hugo lives and not the binary itself.*
* Click OK at every window to exit.

#### For Windows 7 and 8.x users

Windows 7 and 8.1 do not include the easy path editor included in Windows 10, so non-technical users on those platforms are advised to install a free third-party path editor like [Windows Environment Variables Editor].

### Verify the Executable

Run a few commands to verify that the executable is ready to run, and then build a sample site to get started.

#### 1. Open a Command Prompt

At the prompt, type `hugo help` and press the <kbd>Enter</kbd> key. You should see output that starts with:

```txt
hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at https://gohugo.io/.
```

If you do, then the installation is complete. If you don't, double-check the path that you placed the `hugo.exe` file in and that you typed that path correctly when you added it to your `PATH` variable. If you're still not getting the output, search the [Hugo discussion forum][forum] to see if others have already figured out our problem. If not, add a note---in the "Support" category---and be sure to include your command and the output.

At the prompt, change your directory to the `Sites` directory.

```txt
C:\Program Files> cd C:\Hugo\Sites
C:\Hugo\Sites>
```

#### 2. Run the Command

Run the command to generate a new site. I'm using `example.com` as the name of the site.

```txt
C:\Hugo\Sites> hugo new site example.com
```

You should now have a directory at `C:\Hugo\Sites\example.com`. Change into that directory and list the contents. You should get output similar to the following:

```txt
C:\Hugo\Sites> cd example.com
C:\Hugo\Sites\example.com> dir
Directory of C:\hugo\sites\example.com

04/13/2015  10:44 PM    <DIR>          .
04/13/2015  10:44 PM    <DIR>          ..
04/13/2015  10:44 PM    <DIR>          archetypes
04/13/2015  10:44 PM                83 config.toml
04/13/2015  10:44 PM    <DIR>          content
04/13/2015  10:44 PM    <DIR>          data
04/13/2015  10:44 PM    <DIR>          layouts
04/13/2015  10:44 PM    <DIR>          static
               1 File(s)             83 bytes
               7 Dir(s)   6,273,331,200 bytes free
```

### Troubleshoot Windows Installation

[@dhersam][] has created a nice video on common issues:

{{< youtube c8fJIRNChmU >}}

## Linux

### Snap Package

In any of the [Linux distributions that support snaps][snaps], you may install the "extended" Sass/SCSS version with this command:

```txt
snap install hugo --channel=extended
```

To install the non-extended version without Sass/SCSS support:

```txt
snap install hugo
```

To switch between the two, use either `snap refresh hugo --channel=extended` or `snap refresh hugo --channel=stable`.

{{% note %}}
Hugo installed via Snap can write only inside the user’s `$HOME` directory---and gvfs-mounted directories owned by the user---because of Snaps’ confinement and security model. More information is also available [in this related GitHub issue](https://github.com/gohugoio/hugo/issues/3143).
{{% /note %}}

### Debian and Ubuntu

[@anthonyfok](https://github.com/anthonyfok) and friends in the [Debian Go Packaging Team](https://go-team.pages.debian.net/) maintains an official hugo [Debian package](https://packages.debian.org/hugo) which is shared with [Ubuntu](https://packages.ubuntu.com/hugo) and is installable via `apt-get`:

```txt
sudo apt-get install hugo
```

What this installs depends on your Debian/Ubuntu version. On Ubuntu bionic (18.04), this installs the non-extended version without Sass/SCSS support. On Ubuntu disco (19.04), this installs the extended version with Sass/SCSS support.

This option is not recommended because the Hugo in Linux package managers for Debian and Ubuntu is usually a few versions behind as described [here](https://github.com/gcushen/hugo-academic/issues/703)

### Arch Linux

You can also install Hugo from the Arch Linux [community](https://www.archlinux.org/packages/community/x86_64/hugo/) repository. Applies also to derivatives such as Manjaro.

```txt
sudo pacman -S hugo
```

### Fedora, Red Hat and CentOS

Fedora maintains an [official package for Hugo](https://packages.fedoraproject.org/pkgs/hugo/hugo) which may be installed with:

```txt
sudo dnf install hugo
```

For the latest version, the Hugo package maintained by [@daftaupe](https://github.com/daftaupe) at Fedora Copr is recommended:

* <https://copr.fedorainfracloud.org/coprs/daftaupe/hugo/>

See the [related discussion in the Hugo forums][redhatforum].

### openSUSE Tumbleweed

openSUSE maintains an [official package](https://software.opensuse.org/package/hugo) for the Tumbleweed rolling release distribution, it may be installed with:

````txt
sudo zypper install hugo
````

### Solus

Solus includes Hugo in its package repository, it may be installed with:

```txt
sudo eopkg install hugo
```

## OpenBSD

OpenBSD provides a package for Hugo via `pkg_add`:

```txt
doas pkg_add hugo
```


## Upgrade Hugo

Upgrading Hugo is as easy as downloading and replacing the executable you’ve placed in your `PATH` or run `brew upgrade hugo` if using Homebrew.

## Next Steps

Now that you've installed Hugo, read the [Quick Start guide][quickstart] and explore the rest of the documentation. If you have questions, ask the Hugo community directly by visiting the [Hugo Discussion Forum][forum].

[brew]: https://brew.sh/
[macports]: https://www.macports.org/
[Chocolatey]: https://chocolatey.org/
[content]: /content-management/
[@dhersam]: https://github.com/dhersam
[forum]: https://discourse.gohugo.io
[mage]: https://github.com/magefile/mage
[dep]: https://github.com/golang/dep
[highlight shortcode]: /content-management/shortcodes/#highlight
[installgit]: https://git-scm.com/
[GCC]: http://www.mingw.org/
[installgo]: https://go.dev/dl/
[linuxbrew]: https://docs.brew.sh/Homebrew-on-Linux
[quickstart]: /getting-started/quick-start/
[redhatforum]: https://discourse.gohugo.io/t/solved-fedora-copr-repository-out-of-service/2491
[releases]: https://github.com/gohugoio/hugo/releases
[Scoop]: https://scoop.sh/
[snaps]: https://snapcraft.io/docs/installing-snapd
[windowsarch]: https://esupport.trendmicro.com/en-us/home/pages/technical-support/1038680.aspx
[Windows Environment Variables Editor]: https://eveditor.com/
