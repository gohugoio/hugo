---
author: "Michael Henderson"
lastmod: 2016-02-10
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

1. We'll call your website `example.com` for the purpose of this tutorial.
2. You will use `C:\Hugo\Sites` as the starting point for your site.
3. You will use `C:\Hugo\bin` to store executable files.

## Setup Your Directories

You'll need a place to store the Hugo executable, your content (the files that you build), and the generated files (the HTML that Hugo builds for you).

1. Open Windows Explorer.
2. Create a new folder: `C:\Hugo` (assuming you want Hugo on your C drive – it can go anywhere.)
3. Create a subfolder in the Hugo folder: `C:\Hugo\bin`.
4. Create another subfolder in Hugo: `C:\Hugo\Sites`.

## Technical users

1. Download the latest zipped Hugo executable from the [Hugo Releases](https://github.com/spf13/hugo/releases) page.
2. Extract all contents to your `..\Hugo\bin` folder.
3. You'll probably want to rename the Hugo executable to something short like `hugo.exe`.
4. In Powershell or your preferred CLI, add the `hugo.exe` executable to your PATH by navigating to `C:\Hugo\bin` (or the location of your hugo.exe file) and use the command `set PATH=%PATH%;C:\Hugo\bin`. If the `hugo` command does not work after a reboot, you may have to run the command prompt as administrator.

## Less technical users

1. Go the [Hugo Releases](https://github.com/spf13/hugo/releases) page.
2. The latest release is announced on top. Scroll to the bottom of the release announcement to see the downloads. They're all ZIP files.
3. Find the Windows files near the bottom (they're in alphabetical order, so Windows is last) – download either the 32-bit or 64-bit file depending on whether you have 32-bit or 64-bit Windows. (If you don't know, [see here](https://esupport.trendmicro.com/en-us/home/pages/technical-support/1038680.aspx).)
4. Move the ZIP file into your `C:\Hugo\bin` folder.
5. Double-click on the ZIP file and extract its contents. Be sure to extract the contents into the same `C:\Hugo\bin` folder – Windows will do this by default unless you tell it to extract somewhere else.
6. You should now have three new files: an .exe file, license.md, and readme.md. (you can delete the ZIP download now.)
7. Rename the .exe file to `hugo.exe`.
8. Now add Hugo to your Windows PATH settings:

#### For Windows 10 users:
- Right click on the **Start** button
- Click on **System**
- Click on **Advanced System Settings** on the left
- Click on the **Environment Variables** button on the bottom
- In the User variables section, find the row that starts with PATH (PATH will be all caps)
- Double-click on **PATH**
- Click the **New** button.
- Type in Hugo's path, which is `C:\Hugo\bin\hugo.exe` if you went by the instructions above. Press Enter when you're done typing.
- Click OK at every window to exit.

(Note that the path editor in Windows 10 was added in the large [November 2015 Update](https://blogs.windows.com/windowsexperience/2015/11/12/first-major-update-for-windows-10-available-today/). You'll need to have that or a later update installed for the above steps to work. You can see what Windows 10 build you have by clicking on the Start button → Settings → System → About. See [here](http://www.howtogeek.com/236195/how-to-find-out-which-build-and-version-of-windows-10-you-have/) for more.)

Windows 7 and 8.1 do not include an easy path editor, so non-technical users on those platforms are advised to install a free third-party path editor like [Windows Environment Variables Editor](http://eveditor.com/) or [Path Editor](https://patheditor2.codeplex.com/).

## Verify the executable

Run a few commands to verify that the executable is ready to run and then build a sample site to get started.

1. Open a command prompt window.

2. At the prompt, type `hugo help` and press the Enter key. You should see output that starts with:

    {{< nohighlight >}}A Fast and Flexible Static Site Generator built with love by spf13 and friends in Go. Complete documentation is available at http://gohugo.io
{{< /nohighlight >}}

    If you do, then the installation is complete. If you don't, double-check the path that you placed the `hugo.exe` file in and that you typed that path correctly when you added it to your PATH variable. If you're still not getting the output, post a note on the Hugo discussion list (in the `Support` topic) with your command and the output.

3. At the prompt, change your directory to the `Sites` directory.

    {{< nohighlight >}}C:\Program Files> cd C:\Hugo\Sites
C:\Hugo\Sites>
{{< /nohighlight >}}

4. Run the command to generate a new site. I'm using `example.com` as the name of the site.

    {{< nohighlight >}}C:\Hugo\Sites> hugo new site example.com
{{< /nohighlight >}}

5. You should now have a directory at `C:\Hugo\Sites\example.com`.  Change into that directory and list the contents. You should get output similar to the following:

    {{< nohighlight >}}C:\Hugo\Sites&gt;cd example.com
C:\Hugo\Sites\example.com&gt;dir
&nbsp;Directory of C:\hugo\sites\example.com
&nbsp;
04/13/2015  10:44 PM    &lt;DIR&gt;          .
04/13/2015  10:44 PM    &lt;DIR&gt;          ..
04/13/2015  10:44 PM    &lt;DIR&gt;          archetypes
04/13/2015  10:44 PM                83 config.toml
04/13/2015  10:44 PM    &lt;DIR&gt;          content
04/13/2015  10:44 PM    &lt;DIR&gt;          data
04/13/2015  10:44 PM    &lt;DIR&gt;          layouts
04/13/2015  10:44 PM    &lt;DIR&gt;          static
               1 File(s)             83 bytes
               7 Dir(s)   6,273,331,200 bytes free
{{< /nohighlight >}}

You now have Hugo installed and a site to work with. You need to add a layout (or theme), then create some content. Go to http://gohugo.io/overview/quickstart/ for steps on doing that.

## Troubleshooting

@dhersam has created a nice video on common issues:

{{< youtube c8fJIRNChmU >}}
