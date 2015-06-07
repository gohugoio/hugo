---
author: "Michael Henderson"
date: 2015-03-30
linktitle: Installing on Windows
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/mathjax
prev: /tutorials/installing-on-mac
title: Installing on Windows
weight: 10
---

# Installing Hugo on Windows

This tutorial aims to be a complete guide to installing Hugo on your Windows computer.

## Assumptions

1. You know how to open a command prompt window.
2. You're running a 64-bit version of Windows.
3. Your website is `example.com`.
4. You will use `D:\Hugo\Sites` as the starting point for your site.
5. You will use `D:\Hugo\bin` to store executable files.

## Setup Your Directories

You will need a place to store the Hugo executable, your content (the files that you build), and the generated files (the HTML that Hugo builds for you).

1. Open up Windows Explorer.
2. Create a new folder, `D:\Hugo`.
3. Create a new folder, `D:\Hugo\bin`.
4. Create a new folder, `D:\Hugo\Sites`.

## Download the pre-built Hugo executable for Windows

One advantage of building Hugo in go is that there is just a single binary file to use. You don't need to run an installer to use it. Instead, you need to copy the binary to your hard drive. I'm assuming that you'll put it in `D:\Hugo\bin`. If you chose to place it somewhere else, you'll need to substitute that path in the commands.

1. Open https://github.com/spf13/hugo/releases in your browser.
2. The current version is hugo_0.13_windows_amd64.zip.
3. Download that ZIP file and save it in your `D:\Hugo\bin` folder.
4. Find that ZIP file in Windows Explorer and extract all the files from it.
5. You should see a `hugo_0.13_windows_amd64.exe` file.
6. Rename that file to `hugo.exe`.
7. Verify that the `hugo.exe` file is in the `D:\Hugo\bin` folder. (It's possible that the extract put it in a sub-directory. If it did, use Windows Explorer to move it to `D:\Hugo\bin`.)
8. Add the hugo.exe executable to your PATH with: `D:\Hugo\bin>set PATH=%PATH%;D:\Hugo\bin`

## Verify the executable

Run a few commands to verify that the executable is ready to run and then build a sample site to get started.

1. Open a command prompt window.
2. At the prompt, type `hugo help` and press the Enter key. You should see output that starts with:
```
A Fast and Flexible Static Site Generator built with love by spf13 and friends in Go. Complete documentation is available at http://gohugo.io
```
If you do, then the installation is complete. If you don't, double-check the path that you placed the `hugo.exe` file in and that you typed that path correctly when you added it to your PATH variable. If you're still not getting the output, post a note on the Hugo discussion list (in the `Support` topic) with your command and the output.
3. At the prompt, change your directory to the `Sites` directory.
```
C:\Program Files> cd D:\Hugo\Sites
C:\Program Files> D:
D:\Hugo\Sites>
```
4. Run the command to generate a new site. I'm using `example.com` as the name of the site.
```
D:\Hugo\Sites> hugo new site example.com
```

5. You should now have a directory at D:\Hugo\Sites\example.com.  Change into that directory and list the contents. You should get output similar to the following:
```
D:\Hugo\Sites>cd example.com
D:\Hugo\Sites\example.com>dir
 Directory of D:\hugo\sites\example.com

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

You now have Hugo installed and a site to work with. You need to add a layout (or theme), then create some content. Go to http://gohugo.io/overview/quickstart/ for steps on doing that.

## Troubleshooting

@dhersam has a nice video on common issues at https://www.youtube.com/watch?v=c8fJIRNChmU
